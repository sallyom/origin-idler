// +build !ignore_autogenerated

// This file was autogenerated by deepcopy-gen. Do not edit it manually!

package v1alpha2

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CrossGroupObjectReference) DeepCopyInto(out *CrossGroupObjectReference) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CrossGroupObjectReference.
func (in *CrossGroupObjectReference) DeepCopy() *CrossGroupObjectReference {
	if in == nil {
		return nil
	}
	out := new(CrossGroupObjectReference)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Idler) DeepCopyInto(out *Idler) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Idler.
func (in *Idler) DeepCopy() *Idler {
	if in == nil {
		return nil
	}
	out := new(Idler)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Idler) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IdlerList) DeepCopyInto(out *IdlerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Idler, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IdlerList.
func (in *IdlerList) DeepCopy() *IdlerList {
	if in == nil {
		return nil
	}
	out := new(IdlerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *IdlerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IdlerSpec) DeepCopyInto(out *IdlerSpec) {
	*out = *in
	if in.TargetScalables != nil {
		in, out := &in.TargetScalables, &out.TargetScalables
		*out = make([]CrossGroupObjectReference, len(*in))
		copy(*out, *in)
	}
	if in.TriggerServiceNames != nil {
		in, out := &in.TriggerServiceNames, &out.TriggerServiceNames
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IdlerSpec.
func (in *IdlerSpec) DeepCopy() *IdlerSpec {
	if in == nil {
		return nil
	}
	out := new(IdlerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IdlerStatus) DeepCopyInto(out *IdlerStatus) {
	*out = *in
	if in.UnidledScales != nil {
		in, out := &in.UnidledScales, &out.UnidledScales
		*out = make([]UnidleInfo, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IdlerStatus.
func (in *IdlerStatus) DeepCopy() *IdlerStatus {
	if in == nil {
		return nil
	}
	out := new(IdlerStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UnidleInfo) DeepCopyInto(out *UnidleInfo) {
	*out = *in
	out.CrossGroupObjectReference = in.CrossGroupObjectReference
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UnidleInfo.
func (in *UnidleInfo) DeepCopy() *UnidleInfo {
	if in == nil {
		return nil
	}
	out := new(UnidleInfo)
	in.DeepCopyInto(out)
	return out
}
