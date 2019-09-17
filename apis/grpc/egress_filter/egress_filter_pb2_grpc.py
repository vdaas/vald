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
import grpc

import payload_pb2 as payload__pb2


class EgressFilterStub(object):
  # missing associated documentation comment in .proto file
  pass

  def __init__(self, channel):
    """Constructor.

    Args:
      channel: A grpc.Channel.
    """
    self.Filter = channel.unary_unary(
        '/egress_filter.EgressFilter/Filter',
        request_serializer=payload__pb2.Search.Response.SerializeToString,
        response_deserializer=payload__pb2.Search.Response.FromString,
        )
    self.StreamFilter = channel.stream_stream(
        '/egress_filter.EgressFilter/StreamFilter',
        request_serializer=payload__pb2.Object.Distance.SerializeToString,
        response_deserializer=payload__pb2.Object.Distance.FromString,
        )


class EgressFilterServicer(object):
  # missing associated documentation comment in .proto file
  pass

  def Filter(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def StreamFilter(self, request_iterator, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')


def add_EgressFilterServicer_to_server(servicer, server):
  rpc_method_handlers = {
      'Filter': grpc.unary_unary_rpc_method_handler(
          servicer.Filter,
          request_deserializer=payload__pb2.Search.Response.FromString,
          response_serializer=payload__pb2.Search.Response.SerializeToString,
      ),
      'StreamFilter': grpc.stream_stream_rpc_method_handler(
          servicer.StreamFilter,
          request_deserializer=payload__pb2.Object.Distance.FromString,
          response_serializer=payload__pb2.Object.Distance.SerializeToString,
      ),
  }
  generic_handler = grpc.method_handlers_generic_handler(
      'egress_filter.EgressFilter', rpc_method_handlers)
  server.add_generic_rpc_handlers((generic_handler,))
