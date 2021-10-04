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

package v1alpha1

import (
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// BlenderRenderSpec defines the desired state of BlenderRender
type BlenderRenderSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// BlendLocation is an HTTP(S) URL of a .blend file with the embedded assets
	BlendLocation string `json:"blend_location"`

	// RenderType specifies whether this is an `animation` or a single `image` render
	RenderType string `json:"render_type"`

	// S3Endpoint is the API endpoint for the S3 storage
	S3Endpoint string `json:"s3_endpoint"`

	// S3Key is the API username/key for the S3 storage
	S3Key string `json:"s3_key"`

	// S3Secret is the API secret for the S3 storage
	S3Secret string `json:"s3_secret"`
}

// BlenderRenderStatus defines the observed state of BlenderRender
type BlenderRenderJobStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Conditions []batchv1.JobCondition `json:"conditions"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// BlenderRender is the Schema for the blenderrenders API
type BlenderRender struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BlenderRenderSpec      `json:"spec,omitempty"`
	Status BlenderRenderJobStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// BlenderRenderList contains a list of BlenderRender
type BlenderRenderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BlenderRender `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BlenderRender{}, &BlenderRenderList{})
}
