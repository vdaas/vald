# Filter Configuration

This page describes how to enable filtering features on the Vald cluster.

Before using the filtering functions, please check [the Vald Filter Gateway document](../overview/component/filter-gateway.md) first for what you can do.

## Requirement

To use any filtering functions with the Vald cluster, you must deploy the ingress and/or egress filter component before deploying the Vald cluster.

The ingress filter can be used for, e.g., converting the object data to the vector, some filtering query vector, as pre-processing.

The egress filter can be used for, e.g., filtering search result from `vald-lb-gateway` by the distance, categories, or any other condition, as post-processing.

The filter component can be deployed anywhere, but it must be able to communicate with the Vald Filter gateway.
Every filter component should meet Vald's filter gRPC interface.

- The ingress RPC definition

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

- The egress RPC definition

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

For more details, please refer to [the Vald Filter Gateway document](../overview/component/filter-gateway.md)

### Official filter components

Vald provides the two types of ingress filter components.
These will help you to implement your original filter component.

Please refer to:

- [Vald ONNX Ingress Filter](https://github.com/vdaas/vald-onnx-ingress-filter)
- [Vald Tensorflow Ingress Filter](https://github.com/vdaas/vald-tensorflow-ingress-filter)

## Configuration

It is easy to enable the filtering feature.

```yaml
---
gateway:
---
filter:
  enabled: true
```

The Vald Filter gateway connects to the filter component specified in the users' request.
So, you can use the filtering function only by setting `gateway.filter.enabled=true` in your Helm chart.

If you want to make more detailed settings, please set the following parameters.

```yaml
gateway:
  filter:
  ...
    gateway_config:
    ...
      # gRPC client config for ingress filter
      ingress_filter:
        # gRPC client for ingress filter (overrides defaults.grpc.client)
        client: {}
        # object ingress vectorize filter targets
        vectorizer: ""
        # search ingress vector filter targets
        search_filters: []
        # insert ingress vector filter targets
        insert_filters: []
        # update ingress vector filter targets
        update_filters: []
        # upsert ingress vector filter targets
        upsert_filters: []
      # gRPC client config for egress filter
      egress_filter:
        # gRPC client config for egress filter (overrides defaults.grpc.client)
        client: {}
        # object egress vector filter targets
        object_filters: []
        # distance egress vector filter targets
        distance_filters: []
```

Those parameters help the Vald filter gateway connect the filter component before getting the users' requests.

## Client configuration

To use the filter function, you must make settings on the client side because the Vald Filter gateway does not automatically create a connection with the filter component.

These sample code describes how to use in client-side when you send the request with the blob data and want to convert it to a vector using the filter component.

```go
package main

import (
	"context"

	"github.com/vdaas/vald-client-go/v1/payload"
	"github.com/vdaas/vald-client-go/v1/vald"
	"github.com/vdaas/vald/internal/log"
	"google.golang.org/grpc"
)

func main() {
	// address of the Vald cluster ingress.
	grpcServerAddr := "vald-ingress-host"

	// ingress filter host and port.
	const ingressHost = "vald-onnx-ingress-filter"
	const ingressPort = 8081

	// object data and its id.
	const id = "object-id"
	var object []byte

	// connect to the Vald cluster
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, grpcServerAddr, grpc.WithInsecure())
	if err != nil {
		log.Error(err)
		return
	}

	// create client
	client := vald.NewFilterClient(conn)
	// config for insert request
	icfg := &payload.Insert_ObjectRequest{
		// object data you'd like to insert into the Vald cluster
		Object: &payload.Object_Blob{
			Id: id,
			Object: object,
		},
		// insert config
		Config: &payload.Insert_Config{
			SkipStrictExistCheck: false,
		},
		// specify vectorizer component location
		Vectorizer: &payload.Filter_Target{
			Host: ingressHost,
			Port: ingressPort,
		},
	}

	// send insertObject request
	res, err := client.InsertObject(ctx, icfg)
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("location: %#v", res.Ips)
}
```

If you are the Vald cluster operator and not a client user, please share this information with client users.
