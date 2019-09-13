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


class DiscovererStub(object):
  # missing associated documentation comment in .proto file
  pass

  def __init__(self, channel):
    """Constructor.

    Args:
      channel: A grpc.Channel.
    """
    self.Discover = channel.unary_unary(
        '/discoverer.Discoverer/Discover',
        request_serializer=payload__pb2.Common.Empty.SerializeToString,
        response_deserializer=payload__pb2.Info.Agents.FromString,
        )


class DiscovererServicer(object):
  # missing associated documentation comment in .proto file
  pass

  def Discover(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')


def add_DiscovererServicer_to_server(servicer, server):
  rpc_method_handlers = {
      'Discover': grpc.unary_unary_rpc_method_handler(
          servicer.Discover,
          request_deserializer=payload__pb2.Common.Empty.FromString,
          response_serializer=payload__pb2.Info.Agents.SerializeToString,
      ),
  }
  generic_handler = grpc.method_handlers_generic_handler(
      'discoverer.Discoverer', rpc_method_handlers)
  server.add_generic_rpc_handlers((generic_handler,))
