package payload

import proto "github.com/gogo/protobuf/proto"

type Payload interface {
	Reset()
	String() string
	ProtoMessage()
	Descriptor() ([]byte, []int)
	XXX_Unmarshal(b []byte) error
	XXX_Marshal(b []byte, deterministic bool) ([]byte, error)
	XXX_Merge(src proto.Message)
	XXX_Size() int
	XXX_DiscardUnknown()
	XXX_MessageName() string
}
