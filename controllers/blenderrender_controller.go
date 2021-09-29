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
		log.Info("ConfigMap not found - will create")
	} else if err != nil {
		log.Error(err, "Failed to get ConfigMap")
		return ctrl.Result{}, err
	}

	// Check if the job already exists and, if not, create a new one
	job_found := &batchv1.Job{}
	err = r.Get(ctx, types.NamespacedName{Name: blender_render.Name, Namespace: blender_render.Namespace}, job_found)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		log.Info("Creating a new Job")
		return ctrl.Result{}, nil
		//dep := r.deploymentForMemcached(memcached)
		//log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		//err = r.Create(ctx, dep)
		//if err != nil {
		//	log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		//	return ctrl.Result{}, err
		//}
		//// Deployment created successfully - return and requeue
		//return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *BlenderRenderReconciler) createConfigMap(blender_render *blenderv1alpha1.BlenderRender) *corev1.ConfigMap {
	// labels := labelsForBlenderRender(blender_render.Name)

	config_map := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      blender_render.Name,
			Namespace: blender_render.Namespace,
		},
		Data: {
			blend_location: blender_render.Spec.BlendLocation
		}
	}
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
