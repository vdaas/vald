# -*- coding: utf-8 -*-

#
# Copyright (C) 2019 kpango (Yusuke Kato)
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
import sys
_b=sys.version_info[0]<3 and (lambda x:x) or (lambda x:x.encode('latin1'))
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


import payload_pb2 as payload__pb2
from google.api import annotations_pb2 as google_dot_api_dot_annotations__pb2
from pb import gql_pb2 as pb_dot_gql__pb2


DESCRIPTOR = _descriptor.FileDescriptor(
  name='discoverer.proto',
  package='discoverer',
  syntax='proto3',
  serialized_options=_b('\n\031org.vdaas.vald.discovererB\nDiscovererP\001Z*github.com/vdaas/vald/apis/grpc/discoverer'),
  serialized_pb=_b('\n\x10\x64iscoverer.proto\x12\ndiscoverer\x1a\rpayload.proto\x1a\x1cgoogle/api/annotations.proto\x1a\x0cpb/gql.proto2^\n\nDiscoverer\x12J\n\x08\x44iscover\x12\x15.payload.Common.Empty\x1a\x14.payload.Info.Agents\"\x11\x82\xd3\xe4\x93\x02\x0b\x12\t/discover\x1a\x04\xb0\xe0\x1f\x02\x42U\n\x19org.vdaas.vald.discovererB\nDiscovererP\x01Z*github.com/vdaas/vald/apis/grpc/discovererb\x06proto3')
  ,
  dependencies=[payload__pb2.DESCRIPTOR,google_dot_api_dot_annotations__pb2.DESCRIPTOR,pb_dot_gql__pb2.DESCRIPTOR,])



_sym_db.RegisterFileDescriptor(DESCRIPTOR)


DESCRIPTOR._options = None

_DISCOVERER = _descriptor.ServiceDescriptor(
  name='Discoverer',
  full_name='discoverer.Discoverer',
  file=DESCRIPTOR,
  index=0,
  serialized_options=_b('\260\340\037\002'),
  serialized_start=91,
  serialized_end=185,
  methods=[
  _descriptor.MethodDescriptor(
    name='Discover',
    full_name='discoverer.Discoverer.Discover',
    index=0,
    containing_service=None,
    input_type=payload__pb2._COMMON_EMPTY,
    output_type=payload__pb2._INFO_AGENTS,
    serialized_options=_b('\202\323\344\223\002\013\022\t/discover'),
  ),
])
_sym_db.RegisterServiceDescriptor(_DISCOVERER)

DESCRIPTOR.services_by_name['Discoverer'] = _DISCOVERER

# @@protoc_insertion_point(module_scope)
