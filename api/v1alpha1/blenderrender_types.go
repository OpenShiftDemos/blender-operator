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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// BlenderRenderSpec defines the desired state of BlenderRender
type BlenderRenderSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// BlendFileLocation is an HTTP(S) URL of a .blend file with the embedded assets
	BlendFileLocation string `json:"blend_file_location"`

	// OutputFilename is the filename that will be created by the blender run
	OutputFilename string `json:"output_filename"`

	// RenderType specifies whether this is an `animation` or a single `image` render
	RenderType string `json:"render_type"`
}

// BlenderRenderStatus defines the observed state of BlenderRender
type BlenderRenderStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	RenderJobStatus string `json:"render_job_status"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// BlenderRender is the Schema for the blenderrenders API
type BlenderRender struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BlenderRenderSpec   `json:"spec,omitempty"`
	Status BlenderRenderStatus `json:"status,omitempty"`
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