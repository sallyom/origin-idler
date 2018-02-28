


package v1alpha2

import (
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!
// Created by "kubebuilder create resource" for you to implement the Idler resource schema definition
// as a go struct

// IdlerSpec defines the desired state of Idler
type IdlerSpec struct {
    // INSERT YOUR CODE HERE - define desired state schema
}

// IdlerStatus defines the observed state of Idler
type IdlerStatus struct {
    // INSERT YOUR CODE HERE - define observed state schema
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Idler
// +k8s:openapi-gen=true
// +resource:path=idlers
type Idler struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`

    Spec   IdlerSpec   `json:"spec,omitempty"`
    Status IdlerStatus `json:"status,omitempty"`
}
