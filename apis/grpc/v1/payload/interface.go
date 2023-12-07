// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
package payload

import "google.golang.org/protobuf/reflect/protoreflect"

type Payload interface {
	Reset()
	String() string
	ProtoMessage()
	ProtoReflect() protoreflect.Message
	Descriptor() ([]byte, []int)
	MarshalToSizedBufferVT(dAtA []byte) (int, error)
	MarshalToVT(dAtA []byte) (int, error)
	MarshalVT() (dAtA []byte, err error)
	SizeVT() (n int)
	UnmarshalVT(dAtA []byte) error
}
