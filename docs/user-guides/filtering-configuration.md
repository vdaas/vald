# Filter Configuration

This page describes how enabled filtering features for your Vald cluster.

## Requirement

To use any filtering functions with the Vald cluster, you must deploy the ingress or egress filter component before deploying the Vald cluster.

The filter component can be deployed anywhere, but it must be able to communicate with the Vald Filter gateway.
Every filter component which you'd like to use as the Vald filter component must apply Vald's filter gRPC interface.

- The ingress rpc definition

    ```rpc
    service Filter {
      // Represent the RPC to generate the vector.
      rpc GenVector(payload.v1.Object.Blob) returns (payload.v1.Object.Vector) {
        option (google.api.http) = {
          post : "/filter/ingress/object"
          body : "*"
        };
      }
    
      // Represent the RPC to filter the vector.
      rpc FilterVector(payload.v1.Object.Vector)
          returns (payload.v1.Object.Vector) {
        option (google.api.http) = {
          post : "/filter/ingress/vector"
          body : "*"
        };
      }
    }
    ```

- The egress rpc definition

    ```rpc
    service Filter {
    
      // Represent the RPC to filter the distance.
      rpc FilterDistance(payload.v1.Object.Distance)
          returns (payload.v1.Object.Distance) {
        option (google.api.http) = {
          post : "/filter/egress/distance"
          body : "*"
        };
      }
    
      // Represent the RPC to filter the vector.
      rpc FilterVector(payload.v1.Object.Vector)
          returns (payload.v1.Object.Vector) {
        option (google.api.http) = {
          post : "/filter/egress/vector"
          body : "*"
        };
      }
    }
    ```

For more details, please refer to [the Vald Filter Gateway document](../overview/component/filter-gatewya.md)

### Official filter components

Vald provides the two types of ingress filter components.
These will help you to implement your original filter component.

Please refer to:
- [Vald ONNX Ingress Filter](https://github.com/vdaas/vald-onnx-ingress-filter)
- [Vald Tensorflow Ingress Filter](https://github.com/vdaas/vald-tensorflow-ingress-filter)

## Configuration

It is easy to enable the filtering feature.

```yaml
...
gateway:
...
    filter:
        enabled: true
...
```

It is because the Vald Filter gateway connects the filter component specified in the users' request.

If you want to make more detailed settings, please set the following parameters.

<!-- TODO: parameter list -->

```yaml
gateway:
  filter:
  ...
    gateway_config:
    ...
      # @schema {"name": "gateway.filter.gateway_config.ingress_filter", "type": "object"}
      # gateway.filter.gateway_config.ingress_filter -- gRPC client config for ingress filter
      ingress_filter:
        # @schema {"name": "gateway.filter.gateway_config.ingress_filter.client", "alias": "grpc.client"}
        # gateway.filter.gateway_config.ingress_filter.client -- gRPC client for ingress filter (overrides defaults.grpc.client)
        client: {}
        # @schema {"name": "gateway.filter.gateway_config.ingress_filter.vectorizer", "type": "string"}
        # gateway.filter.gateway_config.ingress_filter.vectorizer -- object ingress vectorize filter targets
        vectorizer: ""
        # @schema {"name": "gateway.filter.gateway_config.ingress_filter.search_filters", "type": "array", "items": {"type": "string"}}
        # gateway.filter.gateway_config.ingress_filter.search_filters -- search ingress vector filter targets
        search_filters: []
        # @schema {"name": "gateway.filter.gateway_config.ingress_filter.insert_filters", "type": "array", "items": {"type": "string"}}
        # gateway.filter.gateway_config.ingress_filter.insert_filters -- insert ingress vector filter targets
        insert_filters: []
        # @schema {"name": "gateway.filter.gateway_config.ingress_filter.update_filters", "type": "array", "items": {"type": "string"}}
        # gateway.filter.gateway_config.ingress_filter.update_filters -- update ingress vector filter targets
        update_filters: []
        # @schema {"name": "gateway.filter.gateway_config.ingress_filter.upsert_filters", "type": "array", "items": {"type": "string"}}
        # gateway.filter.gateway_config.ingress_filter.upsert_filters -- upsert ingress vector filter targets
        upsert_filters: []
      # @schema {"name": "gateway.filter.gateway_config.egress_filter", "type": "object"}
      # gateway.filter.gateway_config.egress_filter -- gRPC client config for egress filter
      egress_filter:
        # @schema {"name": "gateway.filter.gateway_config.egress_filter.client", "alias": "grpc.client"}
        # gateway.filter.gateway_config.egress_filter.client -- gRPC client config for egress filter (overrides defaults.grpc.client)
        client: {}
        # @schema {"name": "gateway.filter.gateway_config.egress_filter.object_filters", "type": "array", "items": {"type": "string"}}
        # gateway.filter.gateway_config.egress_filter.object_filters -- object egress vector filter targets
        object_filters: []
        # @schema {"name": "gateway.filter.gateway_config.egress_filter.distance_filters", "type": "array", "items": {"type": "string"}}
        # gateway.filter.gateway_config.egress_filter.distance_filters -- distance egress vector filter targets
        distance_filters: []
```

Those parameters help the Vald filter gateway connect the filter component before getting the users' requests.

## Client configuration

To use the filter function, you must make settings on the client side because the Vald Filter gateway does not automatically create a connection with the filter component.

For example, when you send the request with the blob data and want to convert it to a vector using the filter component, the client configuration should be:

```python
import grpc
import numpy as np
import tensorflow_hub as hub
import tensorflow_text as text
from vald.v1.payload import payload_pb2
from vald.v1.vald import (
    filter_pb2_grpc,
    search_pb2_grpc,
)

# preprocess
preprocess = hub.load('https://tfhub.dev/tensorflow/bert_en_uncased_preprocess/1')
token = preprocess(["TF Hub makes BERT easy!"])
sample = np.vstack([i for i in token.values()])

channel = grpc.insecure_channel("localhost:8081")

# Insert
stub = filter_pb2_grpc.FilterStub(channel)
resize_vector = payload_pb2.Object.ReshapeVector(
    object=sample.tobytes(),
    shape=[3, 128],
)
resize_vector = resize_vector.SerializeToString()

req = payload_pb2.Insert.ObjectRequest(
    object=payload_pb2.Object.Blob(
        id="0",
        object=resize_vector
    ),
    config=payload_pb2.Insert.Config(skip_strict_exist_check=False),
    vectorizer=payload_pb2.Filter.Target(
        host="vald-tensorflow-ingress-filter",
        port=8081,
    )
)
stub.InsertObject(req)
```

If you are the Vald cluster operator and not a client user, please share this information with client users.
