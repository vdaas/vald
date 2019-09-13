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
  name='vald.proto',
  package='vald',
  syntax='proto3',
  serialized_options=_b('\n\016org.vdaas.valdB\004ValdP\001Z$github.com/vdaas/vald/apis/grpc/vald'),
  serialized_pb=_b('\n\nvald.proto\x12\x04vald\x1a\rpayload.proto\x1a\x1cgoogle/api/annotations.proto\x1a\x0cpb/gql.proto2\xa6\t\n\x04Vald\x12\x46\n\x06\x45xists\x12\x12.payload.Object.ID\x1a\x12.payload.Object.ID\"\x14\x82\xd3\xe4\x93\x02\x0e\x12\x0c/exists/{id}\x12O\n\x06Search\x12\x17.payload.Search.Request\x1a\x18.payload.Search.Response\"\x12\x82\xd3\xe4\x93\x02\x0c\"\x07/search:\x01*\x12U\n\nSearchByID\x12\x19.payload.Search.IDRequest\x1a\x18.payload.Search.Response\"\x12\x82\xd3\xe4\x93\x02\x0c\"\x07/search:\x01*\x12G\n\x0cStreamSearch\x12\x17.payload.Search.Request\x1a\x18.payload.Search.Response\"\x00(\x01\x30\x01\x12M\n\x10StreamSearchByID\x12\x19.payload.Search.IDRequest\x1a\x18.payload.Search.Response\"\x00(\x01\x30\x01\x12O\n\x06Insert\x12\x16.payload.Object.Vector\x1a\x15.payload.Common.Error\"\x16\x82\xd3\xe4\x93\x02\x0c\"\x07/insert:\x01*\xb0\xe0\x1f\x01\x12\x43\n\x0cStreamInsert\x12\x16.payload.Object.Vector\x1a\x15.payload.Common.Error\"\x00(\x01\x30\x01\x12@\n\x0bMultiInsert\x12\x17.payload.Object.Vectors\x1a\x16.payload.Common.Errors\"\x00\x12O\n\x06Update\x12\x16.payload.Object.Vector\x1a\x15.payload.Common.Error\"\x16\x82\xd3\xe4\x93\x02\x0c\"\x07/update:\x01*\xb0\xe0\x1f\x01\x12\x43\n\x0cStreamUpdate\x12\x16.payload.Object.Vector\x1a\x15.payload.Common.Error\"\x00(\x01\x30\x01\x12@\n\x0bMultiUpdate\x12\x17.payload.Object.Vectors\x1a\x16.payload.Common.Errors\"\x00\x12M\n\x06Remove\x12\x12.payload.Object.ID\x1a\x15.payload.Common.Error\"\x18\x82\xd3\xe4\x93\x02\x0e*\x0c/remove/{id}\xb0\xe0\x1f\x01\x12?\n\x0cStreamRemove\x12\x12.payload.Object.ID\x1a\x15.payload.Common.Error\"\x00(\x01\x30\x01\x12<\n\x0bMultiRemove\x12\x13.payload.Object.IDs\x1a\x16.payload.Common.Errors\"\x00\x12M\n\tGetObject\x12\x12.payload.Object.ID\x1a\x16.payload.Object.Vector\"\x14\x82\xd3\xe4\x93\x02\x0e\x12\x0c/object/{id}\x12\x43\n\x0fStreamGetObject\x12\x12.payload.Object.ID\x1a\x16.payload.Object.Vector\"\x00(\x01\x30\x01\x1a\x04\xb0\xe0\x1f\x02\x42>\n\x0eorg.vdaas.valdB\x04ValdP\x01Z$github.com/vdaas/vald/apis/grpc/valdb\x06proto3')
  ,
  dependencies=[payload__pb2.DESCRIPTOR,google_dot_api_dot_annotations__pb2.DESCRIPTOR,pb_dot_gql__pb2.DESCRIPTOR,])



_sym_db.RegisterFileDescriptor(DESCRIPTOR)


DESCRIPTOR._options = None

_VALD = _descriptor.ServiceDescriptor(
  name='Vald',
  full_name='vald.Vald',
  file=DESCRIPTOR,
  index=0,
  serialized_options=_b('\260\340\037\002'),
  serialized_start=80,
  serialized_end=1270,
  methods=[
  _descriptor.MethodDescriptor(
    name='Exists',
    full_name='vald.Vald.Exists',
    index=0,
    containing_service=None,
    input_type=payload__pb2._OBJECT_ID,
    output_type=payload__pb2._OBJECT_ID,
    serialized_options=_b('\202\323\344\223\002\016\022\014/exists/{id}'),
  ),
  _descriptor.MethodDescriptor(
    name='Search',
    full_name='vald.Vald.Search',
    index=1,
    containing_service=None,
    input_type=payload__pb2._SEARCH_REQUEST,
    output_type=payload__pb2._SEARCH_RESPONSE,
    serialized_options=_b('\202\323\344\223\002\014\"\007/search:\001*'),
  ),
  _descriptor.MethodDescriptor(
    name='SearchByID',
    full_name='vald.Vald.SearchByID',
    index=2,
    containing_service=None,
    input_type=payload__pb2._SEARCH_IDREQUEST,
    output_type=payload__pb2._SEARCH_RESPONSE,
    serialized_options=_b('\202\323\344\223\002\014\"\007/search:\001*'),
  ),
  _descriptor.MethodDescriptor(
    name='StreamSearch',
    full_name='vald.Vald.StreamSearch',
    index=3,
    containing_service=None,
    input_type=payload__pb2._SEARCH_REQUEST,
    output_type=payload__pb2._SEARCH_RESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='StreamSearchByID',
    full_name='vald.Vald.StreamSearchByID',
    index=4,
    containing_service=None,
    input_type=payload__pb2._SEARCH_IDREQUEST,
    output_type=payload__pb2._SEARCH_RESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='Insert',
    full_name='vald.Vald.Insert',
    index=5,
    containing_service=None,
    input_type=payload__pb2._OBJECT_VECTOR,
    output_type=payload__pb2._COMMON_ERROR,
    serialized_options=_b('\202\323\344\223\002\014\"\007/insert:\001*\260\340\037\001'),
  ),
  _descriptor.MethodDescriptor(
    name='StreamInsert',
    full_name='vald.Vald.StreamInsert',
    index=6,
    containing_service=None,
    input_type=payload__pb2._OBJECT_VECTOR,
    output_type=payload__pb2._COMMON_ERROR,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='MultiInsert',
    full_name='vald.Vald.MultiInsert',
    index=7,
    containing_service=None,
    input_type=payload__pb2._OBJECT_VECTORS,
    output_type=payload__pb2._COMMON_ERRORS,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='Update',
    full_name='vald.Vald.Update',
    index=8,
    containing_service=None,
    input_type=payload__pb2._OBJECT_VECTOR,
    output_type=payload__pb2._COMMON_ERROR,
    serialized_options=_b('\202\323\344\223\002\014\"\007/update:\001*\260\340\037\001'),
  ),
  _descriptor.MethodDescriptor(
    name='StreamUpdate',
    full_name='vald.Vald.StreamUpdate',
    index=9,
    containing_service=None,
    input_type=payload__pb2._OBJECT_VECTOR,
    output_type=payload__pb2._COMMON_ERROR,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='MultiUpdate',
    full_name='vald.Vald.MultiUpdate',
    index=10,
    containing_service=None,
    input_type=payload__pb2._OBJECT_VECTORS,
    output_type=payload__pb2._COMMON_ERRORS,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='Remove',
    full_name='vald.Vald.Remove',
    index=11,
    containing_service=None,
    input_type=payload__pb2._OBJECT_ID,
    output_type=payload__pb2._COMMON_ERROR,
    serialized_options=_b('\202\323\344\223\002\016*\014/remove/{id}\260\340\037\001'),
  ),
  _descriptor.MethodDescriptor(
    name='StreamRemove',
    full_name='vald.Vald.StreamRemove',
    index=12,
    containing_service=None,
    input_type=payload__pb2._OBJECT_ID,
    output_type=payload__pb2._COMMON_ERROR,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='MultiRemove',
    full_name='vald.Vald.MultiRemove',
    index=13,
    containing_service=None,
    input_type=payload__pb2._OBJECT_IDS,
    output_type=payload__pb2._COMMON_ERRORS,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='GetObject',
    full_name='vald.Vald.GetObject',
    index=14,
    containing_service=None,
    input_type=payload__pb2._OBJECT_ID,
    output_type=payload__pb2._OBJECT_VECTOR,
    serialized_options=_b('\202\323\344\223\002\016\022\014/object/{id}'),
  ),
  _descriptor.MethodDescriptor(
    name='StreamGetObject',
    full_name='vald.Vald.StreamGetObject',
    index=15,
    containing_service=None,
    input_type=payload__pb2._OBJECT_ID,
    output_type=payload__pb2._OBJECT_VECTOR,
    serialized_options=None,
  ),
])
_sym_db.RegisterServiceDescriptor(_VALD)

DESCRIPTOR.services_by_name['Vald'] = _VALD

# @@protoc_insertion_point(module_scope)
