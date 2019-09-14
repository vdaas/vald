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


from validate import validate_pb2 as validate_dot_validate__pb2


DESCRIPTOR = _descriptor.FileDescriptor(
  name='payload.proto',
  package='payload',
  syntax='proto3',
  serialized_options=_b('\n\026org.vdaas.vald.payloadB\013ValdPayloadP\001Z\'github.com/vdaas/vald/apis/grpc/payload'),
  serialized_pb=_b('\n\rpayload.proto\x12\x07payload\x1a\x17validate/validate.proto\"\xd6\x02\n\x06Search\x1aY\n\x07Request\x12&\n\x06vector\x18\x01 \x01(\x0b\x32\x16.payload.Object.Vector\x12&\n\x06\x63onfig\x18\x02 \x01(\x0b\x32\x16.payload.Search.Config\x1aS\n\tIDRequest\x12\x1e\n\x02id\x18\x01 \x01(\x0b\x32\x12.payload.Object.ID\x12&\n\x06\x63onfig\x18\x02 \x01(\x0b\x32\x16.payload.Search.Config\x1a?\n\x06\x43onfig\x12\x14\n\x03num\x18\x01 \x01(\rB\x07\xfa\x42\x04*\x02(\x01\x12\x0e\n\x06radius\x18\x02 \x01(\x02\x12\x0f\n\x07\x65psilon\x18\x03 \x01(\x02\x1a[\n\x08Response\x12)\n\x07results\x18\x01 \x03(\x0b\x32\x18.payload.Object.Distance\x12$\n\x05\x65rror\x18\x02 \x01(\x0b\x32\x15.payload.Common.Error\"\x81\x02\n\x06Object\x1a<\n\x08\x44istance\x12\x1e\n\x02id\x18\x01 \x01(\x0b\x32\x12.payload.Object.ID\x12\x10\n\x08\x64istance\x18\x02 \x01(\x02\x1a\x19\n\x02ID\x12\x13\n\x02id\x18\x01 \x01(\tB\x07\xfa\x42\x04r\x02\x10\x01\x1a&\n\x03IDs\x12\x1f\n\x03ids\x18\x01 \x03(\x0b\x32\x12.payload.Object.ID\x1a\x42\n\x06Vector\x12\x1e\n\x02id\x18\x01 \x01(\x0b\x32\x12.payload.Object.ID\x12\x18\n\x06vector\x18\x02 \x03(\x01\x42\x08\xfa\x42\x05\x92\x01\x02\x08\x02\x1a\x32\n\x07Vectors\x12\'\n\x07vectors\x18\x01 \x03(\x0b\x32\x16.payload.Object.Vector\"<\n\x08\x43ontroll\x1a\x30\n\x12\x43reateIndexRequest\x12\x1a\n\tpool_size\x18\x01 \x01(\rB\x07\xfa\x42\x04*\x02(\x00\"\xaa\x01\n\x04Info\x1ai\n\x05\x41gent\x12\x13\n\x02ip\x18\x01 \x01(\tB\x07\xfa\x42\x04r\x02x\x01\x12\x16\n\x05\x63ount\x18\x02 \x01(\rB\x07\xfa\x42\x04*\x02(\x00\x12\r\n\x05state\x18\x03 \x01(\t\x12$\n\x05\x65rror\x18\x04 \x01(\x0b\x32\x15.payload.Common.Error\x1a\x37\n\x06\x41gents\x12-\n\x06\x41gents\x18\x01 \x03(\x0b\x32\x13.payload.Info.AgentB\x08\xfa\x42\x05\x92\x01\x02\x08\x01\"\x82\x01\n\x06\x43ommon\x1a\x07\n\x05\x45mpty\x1a>\n\x05\x45rror\x12\x15\n\x04\x63ode\x18\x01 \x01(\rB\x07\xfa\x42\x04*\x02(\x00\x12\x0b\n\x03msg\x18\x02 \x01(\t\x12\x11\n\ttimestamp\x18\x03 \x01(\x03\x1a/\n\x06\x45rrors\x12%\n\x06\x65rrors\x18\x01 \x03(\x0b\x32\x15.payload.Common.ErrorBP\n\x16org.vdaas.vald.payloadB\x0bValdPayloadP\x01Z\'github.com/vdaas/vald/apis/grpc/payloadb\x06proto3')
  ,
  dependencies=[validate_dot_validate__pb2.DESCRIPTOR,])




_SEARCH_REQUEST = _descriptor.Descriptor(
  name='Request',
  full_name='payload.Search.Request',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='vector', full_name='payload.Search.Request.vector', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='config', full_name='payload.Search.Request.config', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=62,
  serialized_end=151,
)

_SEARCH_IDREQUEST = _descriptor.Descriptor(
  name='IDRequest',
  full_name='payload.Search.IDRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='id', full_name='payload.Search.IDRequest.id', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='config', full_name='payload.Search.IDRequest.config', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=153,
  serialized_end=236,
)

_SEARCH_CONFIG = _descriptor.Descriptor(
  name='Config',
  full_name='payload.Search.Config',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='num', full_name='payload.Search.Config.num', index=0,
      number=1, type=13, cpp_type=3, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=_b('\372B\004*\002(\001'), file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='radius', full_name='payload.Search.Config.radius', index=1,
      number=2, type=2, cpp_type=6, label=1,
      has_default_value=False, default_value=float(0),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='epsilon', full_name='payload.Search.Config.epsilon', index=2,
      number=3, type=2, cpp_type=6, label=1,
      has_default_value=False, default_value=float(0),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=238,
  serialized_end=301,
)

_SEARCH_RESPONSE = _descriptor.Descriptor(
  name='Response',
  full_name='payload.Search.Response',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='results', full_name='payload.Search.Response.results', index=0,
      number=1, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='error', full_name='payload.Search.Response.error', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=303,
  serialized_end=394,
)

_SEARCH = _descriptor.Descriptor(
  name='Search',
  full_name='payload.Search',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
  ],
  extensions=[
  ],
  nested_types=[_SEARCH_REQUEST, _SEARCH_IDREQUEST, _SEARCH_CONFIG, _SEARCH_RESPONSE, ],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=52,
  serialized_end=394,
)


_OBJECT_DISTANCE = _descriptor.Descriptor(
  name='Distance',
  full_name='payload.Object.Distance',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='id', full_name='payload.Object.Distance.id', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='distance', full_name='payload.Object.Distance.distance', index=1,
      number=2, type=2, cpp_type=6, label=1,
      has_default_value=False, default_value=float(0),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=407,
  serialized_end=467,
)

_OBJECT_ID = _descriptor.Descriptor(
  name='ID',
  full_name='payload.Object.ID',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='id', full_name='payload.Object.ID.id', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=_b('\372B\004r\002\020\001'), file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=469,
  serialized_end=494,
)

_OBJECT_IDS = _descriptor.Descriptor(
  name='IDs',
  full_name='payload.Object.IDs',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='ids', full_name='payload.Object.IDs.ids', index=0,
      number=1, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=496,
  serialized_end=534,
)

_OBJECT_VECTOR = _descriptor.Descriptor(
  name='Vector',
  full_name='payload.Object.Vector',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='id', full_name='payload.Object.Vector.id', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='vector', full_name='payload.Object.Vector.vector', index=1,
      number=2, type=1, cpp_type=5, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=_b('\372B\005\222\001\002\010\002'), file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=536,
  serialized_end=602,
)

_OBJECT_VECTORS = _descriptor.Descriptor(
  name='Vectors',
  full_name='payload.Object.Vectors',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='vectors', full_name='payload.Object.Vectors.vectors', index=0,
      number=1, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=604,
  serialized_end=654,
)

_OBJECT = _descriptor.Descriptor(
  name='Object',
  full_name='payload.Object',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
  ],
  extensions=[
  ],
  nested_types=[_OBJECT_DISTANCE, _OBJECT_ID, _OBJECT_IDS, _OBJECT_VECTOR, _OBJECT_VECTORS, ],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=397,
  serialized_end=654,
)


_CONTROLL_CREATEINDEXREQUEST = _descriptor.Descriptor(
  name='CreateIndexRequest',
  full_name='payload.Controll.CreateIndexRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='pool_size', full_name='payload.Controll.CreateIndexRequest.pool_size', index=0,
      number=1, type=13, cpp_type=3, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=_b('\372B\004*\002(\000'), file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=668,
  serialized_end=716,
)

_CONTROLL = _descriptor.Descriptor(
  name='Controll',
  full_name='payload.Controll',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
  ],
  extensions=[
  ],
  nested_types=[_CONTROLL_CREATEINDEXREQUEST, ],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=656,
  serialized_end=716,
)


_INFO_AGENT = _descriptor.Descriptor(
  name='Agent',
  full_name='payload.Info.Agent',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='ip', full_name='payload.Info.Agent.ip', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=_b('\372B\004r\002x\001'), file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='count', full_name='payload.Info.Agent.count', index=1,
      number=2, type=13, cpp_type=3, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=_b('\372B\004*\002(\000'), file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='state', full_name='payload.Info.Agent.state', index=2,
      number=3, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='error', full_name='payload.Info.Agent.error', index=3,
      number=4, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=727,
  serialized_end=832,
)

_INFO_AGENTS = _descriptor.Descriptor(
  name='Agents',
  full_name='payload.Info.Agents',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='Agents', full_name='payload.Info.Agents.Agents', index=0,
      number=1, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=_b('\372B\005\222\001\002\010\001'), file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=834,
  serialized_end=889,
)

_INFO = _descriptor.Descriptor(
  name='Info',
  full_name='payload.Info',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
  ],
  extensions=[
  ],
  nested_types=[_INFO_AGENT, _INFO_AGENTS, ],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=719,
  serialized_end=889,
)


_COMMON_EMPTY = _descriptor.Descriptor(
  name='Empty',
  full_name='payload.Common.Empty',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=902,
  serialized_end=909,
)

_COMMON_ERROR = _descriptor.Descriptor(
  name='Error',
  full_name='payload.Common.Error',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='code', full_name='payload.Common.Error.code', index=0,
      number=1, type=13, cpp_type=3, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=_b('\372B\004*\002(\000'), file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='msg', full_name='payload.Common.Error.msg', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='timestamp', full_name='payload.Common.Error.timestamp', index=2,
      number=3, type=3, cpp_type=2, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=911,
  serialized_end=973,
)

_COMMON_ERRORS = _descriptor.Descriptor(
  name='Errors',
  full_name='payload.Common.Errors',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='errors', full_name='payload.Common.Errors.errors', index=0,
      number=1, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=975,
  serialized_end=1022,
)

_COMMON = _descriptor.Descriptor(
  name='Common',
  full_name='payload.Common',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
  ],
  extensions=[
  ],
  nested_types=[_COMMON_EMPTY, _COMMON_ERROR, _COMMON_ERRORS, ],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=892,
  serialized_end=1022,
)

_SEARCH_REQUEST.fields_by_name['vector'].message_type = _OBJECT_VECTOR
_SEARCH_REQUEST.fields_by_name['config'].message_type = _SEARCH_CONFIG
_SEARCH_REQUEST.containing_type = _SEARCH
_SEARCH_IDREQUEST.fields_by_name['id'].message_type = _OBJECT_ID
_SEARCH_IDREQUEST.fields_by_name['config'].message_type = _SEARCH_CONFIG
_SEARCH_IDREQUEST.containing_type = _SEARCH
_SEARCH_CONFIG.containing_type = _SEARCH
_SEARCH_RESPONSE.fields_by_name['results'].message_type = _OBJECT_DISTANCE
_SEARCH_RESPONSE.fields_by_name['error'].message_type = _COMMON_ERROR
_SEARCH_RESPONSE.containing_type = _SEARCH
_OBJECT_DISTANCE.fields_by_name['id'].message_type = _OBJECT_ID
_OBJECT_DISTANCE.containing_type = _OBJECT
_OBJECT_ID.containing_type = _OBJECT
_OBJECT_IDS.fields_by_name['ids'].message_type = _OBJECT_ID
_OBJECT_IDS.containing_type = _OBJECT
_OBJECT_VECTOR.fields_by_name['id'].message_type = _OBJECT_ID
_OBJECT_VECTOR.containing_type = _OBJECT
_OBJECT_VECTORS.fields_by_name['vectors'].message_type = _OBJECT_VECTOR
_OBJECT_VECTORS.containing_type = _OBJECT
_CONTROLL_CREATEINDEXREQUEST.containing_type = _CONTROLL
_INFO_AGENT.fields_by_name['error'].message_type = _COMMON_ERROR
_INFO_AGENT.containing_type = _INFO
_INFO_AGENTS.fields_by_name['Agents'].message_type = _INFO_AGENT
_INFO_AGENTS.containing_type = _INFO
_COMMON_EMPTY.containing_type = _COMMON
_COMMON_ERROR.containing_type = _COMMON
_COMMON_ERRORS.fields_by_name['errors'].message_type = _COMMON_ERROR
_COMMON_ERRORS.containing_type = _COMMON
DESCRIPTOR.message_types_by_name['Search'] = _SEARCH
DESCRIPTOR.message_types_by_name['Object'] = _OBJECT
DESCRIPTOR.message_types_by_name['Controll'] = _CONTROLL
DESCRIPTOR.message_types_by_name['Info'] = _INFO
DESCRIPTOR.message_types_by_name['Common'] = _COMMON
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

Search = _reflection.GeneratedProtocolMessageType('Search', (_message.Message,), {

  'Request' : _reflection.GeneratedProtocolMessageType('Request', (_message.Message,), {
    'DESCRIPTOR' : _SEARCH_REQUEST,
    '__module__' : 'payload_pb2'
    # @@protoc_insertion_point(class_scope:payload.Search.Request)
    })
  ,

  'IDRequest' : _reflection.GeneratedProtocolMessageType('IDRequest', (_message.Message,), {
    'DESCRIPTOR' : _SEARCH_IDREQUEST,
    '__module__' : 'payload_pb2'
    # @@protoc_insertion_point(class_scope:payload.Search.IDRequest)
    })
  ,

  'Config' : _reflection.GeneratedProtocolMessageType('Config', (_message.Message,), {
    'DESCRIPTOR' : _SEARCH_CONFIG,
    '__module__' : 'payload_pb2'
    # @@protoc_insertion_point(class_scope:payload.Search.Config)
    })
  ,

  'Response' : _reflection.GeneratedProtocolMessageType('Response', (_message.Message,), {
    'DESCRIPTOR' : _SEARCH_RESPONSE,
    '__module__' : 'payload_pb2'
    # @@protoc_insertion_point(class_scope:payload.Search.Response)
    })
  ,
  'DESCRIPTOR' : _SEARCH,
  '__module__' : 'payload_pb2'
  # @@protoc_insertion_point(class_scope:payload.Search)
  })
_sym_db.RegisterMessage(Search)
_sym_db.RegisterMessage(Search.Request)
_sym_db.RegisterMessage(Search.IDRequest)
_sym_db.RegisterMessage(Search.Config)
_sym_db.RegisterMessage(Search.Response)

Object = _reflection.GeneratedProtocolMessageType('Object', (_message.Message,), {

  'Distance' : _reflection.GeneratedProtocolMessageType('Distance', (_message.Message,), {
    'DESCRIPTOR' : _OBJECT_DISTANCE,
    '__module__' : 'payload_pb2'
    # @@protoc_insertion_point(class_scope:payload.Object.Distance)
    })
  ,

  'ID' : _reflection.GeneratedProtocolMessageType('ID', (_message.Message,), {
    'DESCRIPTOR' : _OBJECT_ID,
    '__module__' : 'payload_pb2'
    # @@protoc_insertion_point(class_scope:payload.Object.ID)
    })
  ,

  'IDs' : _reflection.GeneratedProtocolMessageType('IDs', (_message.Message,), {
    'DESCRIPTOR' : _OBJECT_IDS,
    '__module__' : 'payload_pb2'
    # @@protoc_insertion_point(class_scope:payload.Object.IDs)
    })
  ,

  'Vector' : _reflection.GeneratedProtocolMessageType('Vector', (_message.Message,), {
    'DESCRIPTOR' : _OBJECT_VECTOR,
    '__module__' : 'payload_pb2'
    # @@protoc_insertion_point(class_scope:payload.Object.Vector)
    })
  ,

  'Vectors' : _reflection.GeneratedProtocolMessageType('Vectors', (_message.Message,), {
    'DESCRIPTOR' : _OBJECT_VECTORS,
    '__module__' : 'payload_pb2'
    # @@protoc_insertion_point(class_scope:payload.Object.Vectors)
    })
  ,
  'DESCRIPTOR' : _OBJECT,
  '__module__' : 'payload_pb2'
  # @@protoc_insertion_point(class_scope:payload.Object)
  })
_sym_db.RegisterMessage(Object)
_sym_db.RegisterMessage(Object.Distance)
_sym_db.RegisterMessage(Object.ID)
_sym_db.RegisterMessage(Object.IDs)
_sym_db.RegisterMessage(Object.Vector)
_sym_db.RegisterMessage(Object.Vectors)

Controll = _reflection.GeneratedProtocolMessageType('Controll', (_message.Message,), {

  'CreateIndexRequest' : _reflection.GeneratedProtocolMessageType('CreateIndexRequest', (_message.Message,), {
    'DESCRIPTOR' : _CONTROLL_CREATEINDEXREQUEST,
    '__module__' : 'payload_pb2'
    # @@protoc_insertion_point(class_scope:payload.Controll.CreateIndexRequest)
    })
  ,
  'DESCRIPTOR' : _CONTROLL,
  '__module__' : 'payload_pb2'
  # @@protoc_insertion_point(class_scope:payload.Controll)
  })
_sym_db.RegisterMessage(Controll)
_sym_db.RegisterMessage(Controll.CreateIndexRequest)

Info = _reflection.GeneratedProtocolMessageType('Info', (_message.Message,), {

  'Agent' : _reflection.GeneratedProtocolMessageType('Agent', (_message.Message,), {
    'DESCRIPTOR' : _INFO_AGENT,
    '__module__' : 'payload_pb2'
    # @@protoc_insertion_point(class_scope:payload.Info.Agent)
    })
  ,

  'Agents' : _reflection.GeneratedProtocolMessageType('Agents', (_message.Message,), {
    'DESCRIPTOR' : _INFO_AGENTS,
    '__module__' : 'payload_pb2'
    # @@protoc_insertion_point(class_scope:payload.Info.Agents)
    })
  ,
  'DESCRIPTOR' : _INFO,
  '__module__' : 'payload_pb2'
  # @@protoc_insertion_point(class_scope:payload.Info)
  })
_sym_db.RegisterMessage(Info)
_sym_db.RegisterMessage(Info.Agent)
_sym_db.RegisterMessage(Info.Agents)

Common = _reflection.GeneratedProtocolMessageType('Common', (_message.Message,), {

  'Empty' : _reflection.GeneratedProtocolMessageType('Empty', (_message.Message,), {
    'DESCRIPTOR' : _COMMON_EMPTY,
    '__module__' : 'payload_pb2'
    # @@protoc_insertion_point(class_scope:payload.Common.Empty)
    })
  ,

  'Error' : _reflection.GeneratedProtocolMessageType('Error', (_message.Message,), {
    'DESCRIPTOR' : _COMMON_ERROR,
    '__module__' : 'payload_pb2'
    # @@protoc_insertion_point(class_scope:payload.Common.Error)
    })
  ,

  'Errors' : _reflection.GeneratedProtocolMessageType('Errors', (_message.Message,), {
    'DESCRIPTOR' : _COMMON_ERRORS,
    '__module__' : 'payload_pb2'
    # @@protoc_insertion_point(class_scope:payload.Common.Errors)
    })
  ,
  'DESCRIPTOR' : _COMMON,
  '__module__' : 'payload_pb2'
  # @@protoc_insertion_point(class_scope:payload.Common)
  })
_sym_db.RegisterMessage(Common)
_sym_db.RegisterMessage(Common.Empty)
_sym_db.RegisterMessage(Common.Error)
_sym_db.RegisterMessage(Common.Errors)


DESCRIPTOR._options = None
_SEARCH_CONFIG.fields_by_name['num']._options = None
_OBJECT_ID.fields_by_name['id']._options = None
_OBJECT_VECTOR.fields_by_name['vector']._options = None
_CONTROLL_CREATEINDEXREQUEST.fields_by_name['pool_size']._options = None
_INFO_AGENT.fields_by_name['ip']._options = None
_INFO_AGENT.fields_by_name['count']._options = None
_INFO_AGENTS.fields_by_name['Agents']._options = None
_COMMON_ERROR.fields_by_name['code']._options = None
# @@protoc_insertion_point(module_scope)
