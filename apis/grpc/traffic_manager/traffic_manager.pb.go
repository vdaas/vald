//
// Copyright (C) 2019-2019 kpango (Yusuke Kato)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//



package traffic_manager

import (
	fmt "fmt"
	_ "github.com/danielvladco/go-proto-gql/pb"
	proto "github.com/gogo/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

func init() { proto.RegisterFile("traffic_manager.proto", fileDescriptor_d948c4d9e0d10939) }

var fileDescriptor_d948c4d9e0d10939 = []byte{
	// 141 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2d, 0x29, 0x4a, 0x4c,
	0x4b, 0xcb, 0x4c, 0x8e, 0xcf, 0x4d, 0xcc, 0x4b, 0x4c, 0x4f, 0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f,
	0xc9, 0x17, 0xe2, 0x47, 0x13, 0x96, 0x92, 0x49, 0xcf, 0xcf, 0x4f, 0xcf, 0x49, 0xd5, 0x4f, 0x2c,
	0xc8, 0xd4, 0x4f, 0xcc, 0xcb, 0xcb, 0x2f, 0x49, 0x2c, 0xc9, 0xcc, 0xcf, 0x2b, 0x86, 0x28, 0x97,
	0xe2, 0x29, 0x48, 0xd2, 0x4f, 0x2f, 0xcc, 0x81, 0xf0, 0x9c, 0x6c, 0x4f, 0x3c, 0x92, 0x63, 0xbc,
	0xf0, 0x48, 0x8e, 0xf1, 0xc1, 0x23, 0x39, 0xc6, 0x28, 0xfd, 0xf4, 0xcc, 0x92, 0x8c, 0xd2, 0x24,
	0xbd, 0xe4, 0xfc, 0x5c, 0xfd, 0xb2, 0x94, 0xc4, 0xc4, 0x62, 0xfd, 0xb2, 0xc4, 0x9c, 0x14, 0x90,
	0x41, 0xc5, 0xfa, 0xe9, 0x45, 0x05, 0xc9, 0xfa, 0x68, 0x56, 0x25, 0xb1, 0x81, 0x4d, 0x31, 0x06,
	0x04, 0x00, 0x00, 0xff, 0xff, 0xc9, 0xcf, 0x88, 0xef, 0x9b, 0x00, 0x00, 0x00,
}
