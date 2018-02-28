package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!
// Created by "kubebuilder create resource" for you to implement the Idler resource schema definition
// as a go struct

// IdlerSpec defines the desired state of Idler
type IdlerSpec struct {
	// WantIdle represents the desired state of idling
	WantIdle bool
	// TargetScalables contains the collection of scalables that
	// are idled/unidled together.
	TargetScalables []CrossGroupObjectReference
	// TriggerServiceNames contains the collection of services that shold
	// trigger unidling.  Their corresponding endpoints objects will be
	// used to determine whether or not unidling is successful.
	TriggerServiceNames []string
}

// IdlerStatus defines the observed state of Idler
type IdlerStatus struct {
	// Idled represents the current state of idling
	Idled bool
	// UnidleScales contains the previous scales of idled scalables
	UnidledScales []UnidleInfo
}

// UnidleInfo represents the information needed to restore an idled object
// to its unidled state.
type UnidleInfo struct {
	CrossGroupObjectReference
	// PreviousScale represents the replica count of this object before it
	// was idled.
	PreviousScale int32
}

// CrossGroupObjectReference references an object in the same namespace as
// the current "context", but potentially in a different API group.
type CrossGroupObjectReference struct {
	// TODO(directxman12): ask deads/liggitt if we're still
	// going to fight the Group vs APIVersion battle...

	// Group is the API group that the given resource belongs to.
	Group string
	// Resource is the type of resource that this references.
	Resource string
	// Name is the name of the object that we're referencing.
	Name string
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
