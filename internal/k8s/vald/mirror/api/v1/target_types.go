package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	// GroupVersion is group version used to register these objects
	GroupVersion = schema.GroupVersion{Group: "vald.vdaas.org", Version: "v1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)

func init() {
	SchemeBuilder.Register(
		&ValdMirrorTarget{},
		&ValdMirrorTargetList{},
	)
}

// ValdMirrorTarget is a mirror information.
type ValdMirrorTarget struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   MirrorTargetSpec   `json:"spec,omitempty"`
	Status MirrorTargetStatus `json:"status,omitempty"`
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ValdMirrorTarget) DeepCopyInto(out *ValdMirrorTarget) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ValdMirror.
func (in *ValdMirrorTarget) DeepCopy() *ValdMirrorTarget {
	if in == nil {
		return nil
	}
	out := new(ValdMirrorTarget)
	in.DeepCopyInto(out)
	return out
}

func (in *ValdMirrorTarget) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type MirrorTargetSpec struct {
	Colocation string       `json:"colocation,omitempty"`
	Target     MirrorTarget `json:"target,omitempty"`
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MirrorTargetSpec) DeepCopyInto(out *MirrorTargetSpec) {
	*out = *in
	out.Colocation = in.Colocation
	out.Target = in.Target
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MirrorSpec.
func (in *MirrorTargetSpec) DeepCopy() *MirrorTargetSpec {
	if in == nil {
		return nil
	}
	out := new(MirrorTargetSpec)
	in.DeepCopyInto(out)
	return out
}

type MirrorTarget struct {
	Host string `json:"host,omitempty"`
	Port int    `json:"port,omitempty"`
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MirrorTarget) DeepCopyInto(out *MirrorTarget) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MirrorTarget.
func (in *MirrorTarget) DeepCopy() *MirrorTarget {
	if in == nil {
		return nil
	}
	out := new(MirrorTarget)
	in.DeepCopyInto(out)
	return out
}

// MirrorTargetStatus is status of ValdMirrorTarget
type MirrorTargetStatus struct {
	Phase              MirrorTargetPhase `json:"phase,omitempty"`
	LastTransitionTime string            `json:"last_transition_time,omitempty"`
}

type MirrorTargetPhase string

const (
	MirrorTargetPending      = MirrorTargetPhase("Pending")
	MirrorTargetConnected    = MirrorTargetPhase("Connected")
	MirrorTargetDisconnected = MirrorTargetPhase("Disconnected")
	MirrorTargetUnknown      = MirrorTargetPhase("Unknown")
)

// ValdMirrorList is the whole list of all ValdMirror which have been registered with master.
type ValdMirrorTargetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Items []ValdMirrorTarget `json:"items,omitempty"`
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ValdMirrorTargetList) DeepCopyInto(out *ValdMirrorTargetList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if len(in.Items) != 0 {
		out.Items = make([]ValdMirrorTarget, len(in.Items))
		for i := 0; i < len(in.Items); i++ {
			out.Items[i] = *in.Items[i].DeepCopy()
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ValdMirrorList.
func (in *ValdMirrorTargetList) DeepCopy() *ValdMirrorTargetList {
	if in == nil {
		return nil
	}
	out := new(ValdMirrorTargetList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ValdMirrorTargetList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}
