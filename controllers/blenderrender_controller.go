/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"reflect"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"

	blenderv1alpha1 "github.com/openshiftdemos/blender-operator/api/v1alpha1"
)

// BlenderRenderReconciler reconciles a BlenderRender object
type BlenderRenderReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=blender.openshiftdemos.github.io,resources=blenderrenders,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=blender.openshiftdemos.github.io,resources=blenderrenders/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=blender.openshiftdemos.github.io,resources=blenderrenders/finalizers,verbs=update
//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the BlenderRender object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.9.2/pkg/reconcile
func (r *BlenderRenderReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrllog.FromContext(ctx)
	log.Info("Begin reconcile loop")

	// Lookup the BlenderRender instance for this reconcile request
	blender_render := &blenderv1alpha1.BlenderRender{}
	err := r.Get(ctx, req.NamespacedName, blender_render)

	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("BlenderRender resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get BlenderRender")
		return ctrl.Result{}, err
	}

	// Find if the configmap already exists and, if not, create one
	cf_found := &corev1.ConfigMap{}
	err = r.Get(ctx, types.NamespacedName{Name: blender_render.Name, Namespace: blender_render.Namespace}, cf_found)
	if err != nil && errors.IsNotFound(err) {
		config_map := r.createConfigMap(blender_render)
		log.Info("Creating a new Configmap", "ConfigMap.Namespace", config_map.Namespace, "ConfigMap.Name", config_map.Name)
		err = r.Create(ctx, config_map)
		if err != nil {
			log.Error(err, "Failed to create new ConfigMap", "ConfigMap.Namespace", config_map.Namespace, "ConfigMap.Name", config_map.Name)
			return ctrl.Result{}, err
		}
		// Deployment created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get ConfigMap")
		return ctrl.Result{}, err
	}

	// Check if the job already exists and, if not, create a new one
	job_found := &batchv1.Job{}
	err = r.Get(ctx, types.NamespacedName{Name: blender_render.Name, Namespace: blender_render.Namespace}, job_found)
	if err != nil && errors.IsNotFound(err) {
		// Define a new Job
		job := r.createJob(blender_render)
		log.Info("Creating a new Job", "Job.Namespace", job.Namespace, "Job.Name", job.Name)
		err = r.Create(ctx, job)
		if err != nil {
			log.Info("Failed to create new Job", "Job.Namespace", job.Namespace, "Job.Name", job.Name)
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Job")
		return ctrl.Result{}, err
	}

	// Update the BlenderRender status with the Job conditions
	if !reflect.DeepEqual(job_found.Status.Conditions, blender_render.Status.Conditions) {
		blender_render.Status.Conditions = job_found.Status.Conditions
		err := r.Status().Update(ctx, blender_render)
		if err != nil {
			log.Error(err, "Failed to update BlenderRender status")
		}
	}

	return ctrl.Result{}, nil
}

func (r *BlenderRenderReconciler) createConfigMap(blender_render *blenderv1alpha1.BlenderRender) *corev1.ConfigMap {
	labels := labelsForBlenderRender(blender_render.Name)

	config_map := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      blender_render.Name,
			Namespace: blender_render.Namespace,
			Labels:    labels,
		},
		Data: map[string]string{
			"blend_location": blender_render.Spec.BlendLocation,
			"render_type":    blender_render.Spec.RenderType,
			"s3_endpoint":    blender_render.Spec.S3Endpoint,
			"s3_key":         blender_render.Spec.S3Key,
			"s3_secret":      blender_render.Spec.S3Secret,
		},
	}

	// set BlenderRender instance as the owner and controller
	ctrl.SetControllerReference(blender_render, config_map, r.Scheme)
	return config_map
}

func (r *BlenderRenderReconciler) createJob(blender_render *blenderv1alpha1.BlenderRender) *batchv1.Job {
	labels := labelsForBlenderRender(blender_render.Name)
	back_off_limit := int32(1)

	blend_location_source_keyselector := corev1.ConfigMapKeySelector{
		LocalObjectReference: corev1.LocalObjectReference{Name: blender_render.Name},
		Key:                  "blend_location",
	}

	blend_location_source := corev1.EnvVarSource{
		ConfigMapKeyRef: &blend_location_source_keyselector,
	}

	render_type_source_keyselector := corev1.ConfigMapKeySelector{
		LocalObjectReference: corev1.LocalObjectReference{Name: blender_render.Name},
		Key:                  "render_type",
	}

	render_type_source := corev1.EnvVarSource{
		ConfigMapKeyRef: &render_type_source_keyselector,
	}

	s3_endpoint_source_keyselector := corev1.ConfigMapKeySelector{
		LocalObjectReference: corev1.LocalObjectReference{Name: blender_render.Name},
		Key:                  "s3_endpoint",
	}

	s3_endpoint_source := corev1.EnvVarSource{
		ConfigMapKeyRef: &s3_endpoint_source_keyselector,
	}

	s3_key_source_keyselector := corev1.ConfigMapKeySelector{
		LocalObjectReference: corev1.LocalObjectReference{Name: blender_render.Name},
		Key:                  "s3_key",
	}

	s3_key_source := corev1.EnvVarSource{
		ConfigMapKeyRef: &s3_key_source_keyselector,
	}

	s3_secret_source_keyselector := corev1.ConfigMapKeySelector{
		LocalObjectReference: corev1.LocalObjectReference{Name: blender_render.Name},
		Key:                  "s3_secret",
	}

	s3_secret_source := corev1.EnvVarSource{
		ConfigMapKeyRef: &s3_secret_source_keyselector,
	}

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      blender_render.Name,
			Namespace: blender_render.Namespace,
			Labels:    labels,
		},
		Spec: batchv1.JobSpec{
			BackoffLimit: &back_off_limit,
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: "Never",
					Containers: []corev1.Container{{
						Name:  "blender",
						Image: "quay.io/openshiftdemos/blender-remote:latest",
						Env: []corev1.EnvVar{
							{
								Name:      "BLEND_LOCATION",
								ValueFrom: &blend_location_source,
							},
							{
								Name:      "RENDER_TYPE",
								ValueFrom: &render_type_source,
							},
							{
								Name:      "S3_ENDPOINT",
								ValueFrom: &s3_endpoint_source,
							},
							{
								Name:      "S3_KEY",
								ValueFrom: &s3_key_source,
							},
							{
								Name:      "S3_SECRET",
								ValueFrom: &s3_secret_source,
							},
						},
					}},
				},
			},
		},
	}

	// set BlenderRender instance as the owner and controller
	ctrl.SetControllerReference(blender_render, job, r.Scheme)
	return job
}

// labelsForBlenderRender returns the labels for selecting the resources
// belonging to the given blender render CR name.
func labelsForBlenderRender(name string) map[string]string {
	return map[string]string{"app": "blenderrender", "blenderrender_cr": name}
}

// SetupWithManager sets up the controller with the Manager.
func (r *BlenderRenderReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&blenderv1alpha1.BlenderRender{}).
		Owns(&batchv1.Job{}).
		Complete(r)
}
