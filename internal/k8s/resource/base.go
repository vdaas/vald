// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package resource

import (
	"unsafe"

	"github.com/vdaas/vald/internal/log"
	k8sresource "k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
)

// Quantity is an alias for k8s.io/apimachinery/pkg/api/resource.Quantity.
type Quantity = k8sresource.Quantity

// BinarySI is the binary SI format for Quantity values (Ki, Mi, Gi, …).
const BinarySI = k8sresource.BinarySI

// NewQuantity constructs a Quantity from a value and format, delegating to
// k8s.io/apimachinery/pkg/api/resource.NewQuantity.
func NewQuantity(value int64, format k8sresource.Format) *Quantity {
	return k8sresource.NewQuantity(value, format)
}

// MustParse parses str into a Quantity, panicking if str is invalid.
// Delegates to k8s.io/apimachinery/pkg/api/resource.MustParse.
func MustParse(str string) Quantity {
	return k8sresource.MustParse(str)
}

// DeepCopyIntoer constrains PT to a pointer to T that provides DeepCopyInto,
// the single method each type still has to implement by hand (field
// enumeration cannot be derived generically without reflection).
type DeepCopyIntoer[T any] interface {
	*T
	DeepCopyInto(*T)
}

// Base provides the DeepCopy and DeepCopyObject boilerplate for the type T
// that embeds it, replacing the controller-gen generated wrappers. The PT
// type parameter forces T to implement DeepCopyInto at compile time.
//
// Contract: Base MUST be embedded as the FIRST field of T. Base is
// zero-sized, so the address of the embedded Base equals the address of the
// outer T, which lets the promoted methods recover the outer object without
// reflection:
//
//	type ValdOperatorRelease struct {
//		resource.Base[ValdOperatorRelease, *ValdOperatorRelease]
//		metav1.TypeMeta   `json:",inline"`
//		...
//	}
//
//	func (in *ValdOperatorRelease) DeepCopyInto(out *ValdOperatorRelease) { ... }
type Base[T any, PT DeepCopyIntoer[T]] struct{}

// self recovers the outer object from the embedded zero-sized Base pointer.
// This relies on the first-field embedding contract documented on Base.
func (b *Base[T, PT]) self() PT {
	return PT(unsafe.Pointer(b)) //nolint:gosec // skipcq: GSC-G103 zero-size first-field embedding contract documented on Base.
}

// Base deliberately does NOT provide a DeepCopyInto method: if it did, a type
// that embeds Base but forgets its hand-written DeepCopyInto would still
// satisfy DeepCopyIntoer through the promoted method and recurse infinitely at
// runtime (self().DeepCopyInto -> promoted method -> self().DeepCopyInto).
// Without it, that mistake is a compile error: "*T does not satisfy
// DeepCopyIntoer[T] (missing method DeepCopyInto)".

func (b *Base[T, PT]) DeepCopy() *T {
	if b == nil {
		return nil
	}
	out := new(T)
	b.self().DeepCopyInto(out)
	return out
}

// DeepCopyObject returns a deep copy of the receiver's outer object as a
// runtime.Object. Embed Base only into scheme-registered API types when
// relying on this method; it returns nil for non API types.
func (b *Base[T, PT]) DeepCopyObject() runtime.Object {
	if b == nil {
		return nil
	}
	if obj, ok := any(b.DeepCopy()).(runtime.Object); ok {
		return obj
	}
	log.Errorf("resource.Base: %T does not implement runtime.Object", b.self())
	return nil
}

// CopyPtr returns a copy of the pointee. Use it for pointers to value types
// that own no references.
func CopyPtr[T any](in *T) *T {
	if in == nil {
		return nil
	}
	out := *in
	return &out
}

// CopyPtrInto returns a deep copy of the pointee via its own DeepCopyInto.
// Use it for pointers to types that own references and implement
// DeepCopyInto.
func CopyPtrInto[T any, PT DeepCopyIntoer[T]](in PT) PT {
	if in == nil {
		return nil
	}
	out := new(T)
	in.DeepCopyInto(out)
	return out
}

// CopySliceInto returns a copy of in with every element deep-copied via its
// own DeepCopyInto. Use it for slices of elements that own references and
// implement DeepCopyInto.
func CopySliceInto[T any, PT DeepCopyIntoer[T]](in []T) []T {
	if in == nil {
		return nil
	}
	out := make([]T, len(in))
	for i := range in {
		PT(&in[i]).DeepCopyInto(&out[i])
	}
	return out
}

// CopyPtrSliceInto returns a copy of the pointed-to slice with every element
// deep-copied via its own DeepCopyInto. Use it for pointers to slices of
// elements that own references and implement DeepCopyInto.
func CopyPtrSliceInto[T any, PT DeepCopyIntoer[T]](in *[]T) *[]T {
	if in == nil {
		return nil
	}
	out := CopySliceInto[T, PT](*in)
	return &out
}

// CopySliceFunc returns a copy of in with every element copied by cp. Use it
// for slices of elements that own references.
func CopySliceFunc[S ~[]E, E any](in S, cp func(E) E) S {
	if in == nil {
		return nil
	}
	out := make(S, len(in))
	for i := range in {
		out[i] = cp(in[i])
	}
	return out
}
