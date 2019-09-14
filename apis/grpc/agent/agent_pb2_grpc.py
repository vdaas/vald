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

import grpc

import payload_pb2 as payload__pb2


class AgentStub(object):
  # missing associated documentation comment in .proto file
  pass

  def __init__(self, channel):
    """Constructor.

    Args:
      channel: A grpc.Channel.
    """
    self.Exists = channel.unary_unary(
        '/agent.Agent/Exists',
        request_serializer=payload__pb2.Object.ID.SerializeToString,
        response_deserializer=payload__pb2.Object.ID.FromString,
        )
    self.Search = channel.unary_unary(
        '/agent.Agent/Search',
        request_serializer=payload__pb2.Search.Request.SerializeToString,
        response_deserializer=payload__pb2.Search.Response.FromString,
        )
    self.SearchByID = channel.unary_unary(
        '/agent.Agent/SearchByID',
        request_serializer=payload__pb2.Search.IDRequest.SerializeToString,
        response_deserializer=payload__pb2.Search.Response.FromString,
        )
    self.StreamSearch = channel.stream_stream(
        '/agent.Agent/StreamSearch',
        request_serializer=payload__pb2.Search.Request.SerializeToString,
        response_deserializer=payload__pb2.Search.Response.FromString,
        )
    self.StreamSearchByID = channel.stream_stream(
        '/agent.Agent/StreamSearchByID',
        request_serializer=payload__pb2.Search.IDRequest.SerializeToString,
        response_deserializer=payload__pb2.Search.Response.FromString,
        )
    self.Insert = channel.unary_unary(
        '/agent.Agent/Insert',
        request_serializer=payload__pb2.Object.Vector.SerializeToString,
        response_deserializer=payload__pb2.Common.Error.FromString,
        )
    self.StreamInsert = channel.stream_stream(
        '/agent.Agent/StreamInsert',
        request_serializer=payload__pb2.Object.Vector.SerializeToString,
        response_deserializer=payload__pb2.Common.Error.FromString,
        )
    self.MultiInsert = channel.unary_unary(
        '/agent.Agent/MultiInsert',
        request_serializer=payload__pb2.Object.Vectors.SerializeToString,
        response_deserializer=payload__pb2.Common.Errors.FromString,
        )
    self.Update = channel.unary_unary(
        '/agent.Agent/Update',
        request_serializer=payload__pb2.Object.Vector.SerializeToString,
        response_deserializer=payload__pb2.Common.Error.FromString,
        )
    self.StreamUpdate = channel.stream_stream(
        '/agent.Agent/StreamUpdate',
        request_serializer=payload__pb2.Object.Vector.SerializeToString,
        response_deserializer=payload__pb2.Common.Error.FromString,
        )
    self.MultiUpdate = channel.unary_unary(
        '/agent.Agent/MultiUpdate',
        request_serializer=payload__pb2.Object.Vectors.SerializeToString,
        response_deserializer=payload__pb2.Common.Errors.FromString,
        )
    self.Remove = channel.unary_unary(
        '/agent.Agent/Remove',
        request_serializer=payload__pb2.Object.ID.SerializeToString,
        response_deserializer=payload__pb2.Common.Error.FromString,
        )
    self.StreamRemove = channel.stream_stream(
        '/agent.Agent/StreamRemove',
        request_serializer=payload__pb2.Object.ID.SerializeToString,
        response_deserializer=payload__pb2.Common.Error.FromString,
        )
    self.MultiRemove = channel.unary_unary(
        '/agent.Agent/MultiRemove',
        request_serializer=payload__pb2.Object.IDs.SerializeToString,
        response_deserializer=payload__pb2.Common.Errors.FromString,
        )
    self.GetObject = channel.unary_unary(
        '/agent.Agent/GetObject',
        request_serializer=payload__pb2.Object.ID.SerializeToString,
        response_deserializer=payload__pb2.Object.Vector.FromString,
        )
    self.StreamGetObject = channel.stream_stream(
        '/agent.Agent/StreamGetObject',
        request_serializer=payload__pb2.Object.ID.SerializeToString,
        response_deserializer=payload__pb2.Object.Vector.FromString,
        )
    self.CreateIndex = channel.unary_unary(
        '/agent.Agent/CreateIndex',
        request_serializer=payload__pb2.Controll.CreateIndexRequest.SerializeToString,
        response_deserializer=payload__pb2.Common.Empty.FromString,
        )
    self.SaveIndex = channel.unary_unary(
        '/agent.Agent/SaveIndex',
        request_serializer=payload__pb2.Common.Empty.SerializeToString,
        response_deserializer=payload__pb2.Common.Empty.FromString,
        )


class AgentServicer(object):
  # missing associated documentation comment in .proto file
  pass

  def Exists(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def Search(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def SearchByID(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def StreamSearch(self, request_iterator, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def StreamSearchByID(self, request_iterator, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def Insert(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def StreamInsert(self, request_iterator, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def MultiInsert(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def Update(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def StreamUpdate(self, request_iterator, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def MultiUpdate(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def Remove(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def StreamRemove(self, request_iterator, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def MultiRemove(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def GetObject(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def StreamGetObject(self, request_iterator, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def CreateIndex(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def SaveIndex(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')


def add_AgentServicer_to_server(servicer, server):
  rpc_method_handlers = {
      'Exists': grpc.unary_unary_rpc_method_handler(
          servicer.Exists,
          request_deserializer=payload__pb2.Object.ID.FromString,
          response_serializer=payload__pb2.Object.ID.SerializeToString,
      ),
      'Search': grpc.unary_unary_rpc_method_handler(
          servicer.Search,
          request_deserializer=payload__pb2.Search.Request.FromString,
          response_serializer=payload__pb2.Search.Response.SerializeToString,
      ),
      'SearchByID': grpc.unary_unary_rpc_method_handler(
          servicer.SearchByID,
          request_deserializer=payload__pb2.Search.IDRequest.FromString,
          response_serializer=payload__pb2.Search.Response.SerializeToString,
      ),
      'StreamSearch': grpc.stream_stream_rpc_method_handler(
          servicer.StreamSearch,
          request_deserializer=payload__pb2.Search.Request.FromString,
          response_serializer=payload__pb2.Search.Response.SerializeToString,
      ),
      'StreamSearchByID': grpc.stream_stream_rpc_method_handler(
          servicer.StreamSearchByID,
          request_deserializer=payload__pb2.Search.IDRequest.FromString,
          response_serializer=payload__pb2.Search.Response.SerializeToString,
      ),
      'Insert': grpc.unary_unary_rpc_method_handler(
          servicer.Insert,
          request_deserializer=payload__pb2.Object.Vector.FromString,
          response_serializer=payload__pb2.Common.Error.SerializeToString,
      ),
      'StreamInsert': grpc.stream_stream_rpc_method_handler(
          servicer.StreamInsert,
          request_deserializer=payload__pb2.Object.Vector.FromString,
          response_serializer=payload__pb2.Common.Error.SerializeToString,
      ),
      'MultiInsert': grpc.unary_unary_rpc_method_handler(
          servicer.MultiInsert,
          request_deserializer=payload__pb2.Object.Vectors.FromString,
          response_serializer=payload__pb2.Common.Errors.SerializeToString,
      ),
      'Update': grpc.unary_unary_rpc_method_handler(
          servicer.Update,
          request_deserializer=payload__pb2.Object.Vector.FromString,
          response_serializer=payload__pb2.Common.Error.SerializeToString,
      ),
      'StreamUpdate': grpc.stream_stream_rpc_method_handler(
          servicer.StreamUpdate,
          request_deserializer=payload__pb2.Object.Vector.FromString,
          response_serializer=payload__pb2.Common.Error.SerializeToString,
      ),
      'MultiUpdate': grpc.unary_unary_rpc_method_handler(
          servicer.MultiUpdate,
          request_deserializer=payload__pb2.Object.Vectors.FromString,
          response_serializer=payload__pb2.Common.Errors.SerializeToString,
      ),
      'Remove': grpc.unary_unary_rpc_method_handler(
          servicer.Remove,
          request_deserializer=payload__pb2.Object.ID.FromString,
          response_serializer=payload__pb2.Common.Error.SerializeToString,
      ),
      'StreamRemove': grpc.stream_stream_rpc_method_handler(
          servicer.StreamRemove,
          request_deserializer=payload__pb2.Object.ID.FromString,
          response_serializer=payload__pb2.Common.Error.SerializeToString,
      ),
      'MultiRemove': grpc.unary_unary_rpc_method_handler(
          servicer.MultiRemove,
          request_deserializer=payload__pb2.Object.IDs.FromString,
          response_serializer=payload__pb2.Common.Errors.SerializeToString,
      ),
      'GetObject': grpc.unary_unary_rpc_method_handler(
          servicer.GetObject,
          request_deserializer=payload__pb2.Object.ID.FromString,
          response_serializer=payload__pb2.Object.Vector.SerializeToString,
      ),
      'StreamGetObject': grpc.stream_stream_rpc_method_handler(
          servicer.StreamGetObject,
          request_deserializer=payload__pb2.Object.ID.FromString,
          response_serializer=payload__pb2.Object.Vector.SerializeToString,
      ),
      'CreateIndex': grpc.unary_unary_rpc_method_handler(
          servicer.CreateIndex,
          request_deserializer=payload__pb2.Controll.CreateIndexRequest.FromString,
          response_serializer=payload__pb2.Common.Empty.SerializeToString,
      ),
      'SaveIndex': grpc.unary_unary_rpc_method_handler(
          servicer.SaveIndex,
          request_deserializer=payload__pb2.Common.Empty.FromString,
          response_serializer=payload__pb2.Common.Empty.SerializeToString,
      ),
  }
  generic_handler = grpc.method_handlers_generic_handler(
      'agent.Agent', rpc_method_handlers)
  server.add_generic_rpc_handlers((generic_handler,))
