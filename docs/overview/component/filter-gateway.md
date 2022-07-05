# Vald Filter Gateway

Vald Filter Gateway is the component to communicate user-defined ingress/egress filter components with requests.

## Responsibility

Vald Filter Gateway is the optional component in the Vald cluster implemented for available user-defined filtering.

The responsibility is bypassed between the Vald cluster and the user-defined ingress filter (or egress filter) component.

## Features

### Ingress Filtering

Ingress filtering is the pre-process before the request is processed by the Vald LB Gateway.

This is useful, e.g., when you want to vectorize object data (when you want to handle object data directly as the request query) or when you want to pad the vector, etc.

Here are the steps of pre-processing with the ingress filter:
1. The user's request flows from the filter gateway to the ingress filter component when using the Ingress filter.
1. Ingress filter performs arbitrary processing for each user and returns the result to the filter gateway.
1. The filter gateway passes the ingress filter result to the Vald LB Gateway.(The Vald LB Gateway migration process is the same as when the ingress filter is not used. Click [here](../../overview/data-flow.md) for a detailed data flow.)

<img src="../../../assets/docs/overview/component/filter-gateway/ingress_filtering.svg" alt="example data flow with ingress filtering" />

There are two types of ingress filtering methods: generate vector and vector filtering.
In generate vector, for example, you can vectorize object data using your favorite model.
The filtering vector allows you to filter out vectors that contain outliers so that they are not indexed.

Vald officially offers two types of ingress filter components.
- [vald-onnx-ingress-filter](https://github.com/vdaas/vald-onnx-ingress-filter)
- [vald-tensorflow-ingress-filter](https://github.com/vdaas/vald-tensorflow-ingress-filter)

Of course, you can use your own ingress filter.
When using your own ingress filter component, make sure to meet the following interface.

```rpc
type FilterClient interface {
    // Represent the RPC to generate the vector.
    GenVector (ctx context.Context, in * payload.Object_Blob, opts ... grpc.CallOption) (* payload.Object_Vector, error)
    // Represent the RPC to filter the vector.
    FilterVector (ctx context.Context, in * payload.Object_Vector, opts ... grpc.CallOption) (* payload.Object_Vector, error)
}
```

### Egress Filtering 

<!-- TODO: describe about egress -->
Egress filtering is mainly used as the post-processing for the search result returned from the Vald LB Gateway.
Egress filtering returns the search results filtered by the egress filter component to the requesting user.

Here are the steps of pre-processing with the egress filter:
1. The filter gateway receives the search results from the Vald LB Gateway and passes them to the egress filter component.
1. The egress filter component filters the search results and returns the filtered search results to the Filter Gateway.
1. The filter gateway returns the filtered search results from the Gateway to the user.

<img src="../../../assets/docs/overview/component/filter-gateway/egress_filtering.svg" alt="search data flow with egress filtering" />

There are two types of egress filtering methods: distance filtering and vector filtering.
In distance filtering, for example, you can add a process to exclude vectors that exceed the upper limit of distance from the search results.
Filtering for vectors allows you to add the process of removing vectors from different categories from the search results. (e.g., excluding kid's T-shirts from the search for adult's T-shirts.)

If you want to use this feature, please deploy your own egress filter component, which meets the following interface.

```rpc
type FilterClient interface {
    // Represent the RPC to filter the distance.
    FilterDistance(ctx context.Context, in *payload.Object_Distance, opts ...grpc.CallOption) (*payload.Object_Distance, error)
    // Represent the RPC to filter the vector.
    FilterVector(ctx context.Context, in *payload.Object_Vector, opts ...grpc.CallOption) (*payload.Object_Vector, error)
}
```
