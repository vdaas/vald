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


from google.api import annotations_pb2 as google_dot_api_dot_annotations__pb2
from pb import gql_pb2 as pb_dot_gql__pb2


DESCRIPTOR = _descriptor.FileDescriptor(
  name='ingress_filter.proto',
  package='ingress_filter',
  syntax='proto3',
  serialized_options=_b('Z.github.com/vdaas/vald/apis/grpc/ingress_filter'),
  serialized_pb=_b('\n\x14ingress_filter.proto\x12\x0eingress_filter\x1a\x1cgoogle/api/annotations.proto\x1a\x0cpb/gql.protoB0Z.github.com/vdaas/vald/apis/grpc/ingress_filterb\x06proto3')
  ,
  dependencies=[google_dot_api_dot_annotations__pb2.DESCRIPTOR,pb_dot_gql__pb2.DESCRIPTOR,])



_sym_db.RegisterFileDescriptor(DESCRIPTOR)


DESCRIPTOR._options = None
# @@protoc_insertion_point(module_scope)
