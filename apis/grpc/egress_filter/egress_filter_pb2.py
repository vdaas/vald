# -*- coding: utf-8 -*-

#
# Copyright (C) 2019-2019 kpango (Yusuke Kato)
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
  name='egress_filter.proto',
  package='egress_filter',
  syntax='proto3',
  serialized_options=_b('Z-github.com/vdaas/vald/apis/grpc/egress_filter'),
  serialized_pb=_b('\n\x13\x65gress_filter.proto\x12\regress_filter\x1a\rpayload.proto\x1a\x1cgoogle/api/annotations.proto\x1a\x0cpb/gql.proto2\x9e\x01\n\x0c\x45gressFilter\x12>\n\x06\x46ilter\x12\x18.payload.Search.Response\x1a\x18.payload.Search.Response\"\x00\x12H\n\x0cStreamFilter\x12\x18.payload.Object.Distance\x1a\x18.payload.Object.Distance\"\x00(\x01\x30\x01\x1a\x04\xb0\xe0\x1f\x02\x42/Z-github.com/vdaas/vald/apis/grpc/egress_filterb\x06proto3')
  ,
  dependencies=[payload__pb2.DESCRIPTOR,google_dot_api_dot_annotations__pb2.DESCRIPTOR,pb_dot_gql__pb2.DESCRIPTOR,])



_sym_db.RegisterFileDescriptor(DESCRIPTOR)


DESCRIPTOR._options = None

_EGRESSFILTER = _descriptor.ServiceDescriptor(
  name='EgressFilter',
  full_name='egress_filter.EgressFilter',
  file=DESCRIPTOR,
  index=0,
  serialized_options=_b('\260\340\037\002'),
  serialized_start=98,
  serialized_end=256,
  methods=[
  _descriptor.MethodDescriptor(
    name='Filter',
    full_name='egress_filter.EgressFilter.Filter',
    index=0,
    containing_service=None,
    input_type=payload__pb2._SEARCH_RESPONSE,
    output_type=payload__pb2._SEARCH_RESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='StreamFilter',
    full_name='egress_filter.EgressFilter.StreamFilter',
    index=1,
    containing_service=None,
    input_type=payload__pb2._OBJECT_DISTANCE,
    output_type=payload__pb2._OBJECT_DISTANCE,
    serialized_options=None,
  ),
])
_sym_db.RegisterServiceDescriptor(_EGRESSFILTER)

DESCRIPTOR.services_by_name['EgressFilter'] = _EGRESSFILTER

# @@protoc_insertion_point(module_scope)
