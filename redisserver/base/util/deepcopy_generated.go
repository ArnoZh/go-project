// +build !ignore_autogenerated

/*
Copyright The Kubernetes Authors.

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

// Code generated by deepcopy-gen. DO NOT EDIT.

package util

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (d *Dirty) DeepCopyInto(out *Dirty) {
	*out = *d
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Dirty.
func (d *Dirty) DeepCopy() *Dirty {
	if d == nil {
		return nil
	}
	out := new(Dirty)
	d.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (kvp *Int32KVPairArray) DeepCopyInto(out *Int32KVPairArray) {
	*out = *kvp
	if kvp.Array != nil {
		in, out := &kvp.Array, &out.Array
		*out = make([]*Int32KVpair, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(Int32KVpair)
				**out = **in
			}
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Int32KVPairArray.
func (kvp *Int32KVPairArray) DeepCopy() *Int32KVPairArray {
	if kvp == nil {
		return nil
	}
	out := new(Int32KVPairArray)
	kvp.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Int32KVpair) DeepCopyInto(out *Int32KVpair) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Int32KVpair.
func (in *Int32KVpair) DeepCopy() *Int32KVpair {
	if in == nil {
		return nil
	}
	out := new(Int32KVpair)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StringKVpair) DeepCopyInto(out *StringKVpair) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StringKVpair.
func (in *StringKVpair) DeepCopy() *StringKVpair {
	if in == nil {
		return nil
	}
	out := new(StringKVpair)
	in.DeepCopyInto(out)
	return out
}