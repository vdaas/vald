# Vald Search APIs

## Overview

Search Service is responsible for searching vectors similar to the user request vector from `vald-agent`.

```rpc
service Search {

  rpc Search(payload.v1.Search.Request) returns (payload.v1.Search.Response) {}
  rpc SearchByID(payload.v1.Search.IDRequest) returns (payload.v1.Search.Response) {}
  rpc StreamSearch(payload.v1.Search.Request) returns (payload.v1.Search.StreamResponse) {}
  rpc StreamSearchByID(payload.v1.Search.IDRequest) returns (payload.v1.Search.StreamResponse) {}
  rpc MultiSearch(payload.v1.Search.MultiRequest) returns (payload.v1.Search.Responses) {}
  rpc MultiSearchByID(payload.v1.Search.MultiIDRequest) returns (payload.v1.Search.Responses) {}
  rpc LinearSearch(payload.v1.Search.Request) returns (payload.v1.Search.Response) {}
  rpc LinearSearchByID(payload.v1.Search.IDRequest) returns (payload.v1.Search.Response) {}
  rpc StreamLinearSearch(payload.v1.Search.Request) returns (payload.v1.Search.StreamResponse) {}
  rpc StreamLinearSearchByID(payload.v1.Search.IDRequest) returns (payload.v1.Search.StreamResponse) {}
  rpc MultiLinearSearch(payload.v1.Search.MultiRequest) returns (payload.v1.Search.Responses) {}
  rpc MultiLinearSearchByID(payload.v1.Search.MultiIDRequest) returns (payload.v1.Search.Responses) {}

}
```
## Search RPC

Search RPC is the method to search vector(s) similar to the request vector.


### Input

- the scheme of `payload.v1.Search.Request`

  ```rpc
  message Search.Request {
    repeated float vector = 1;
    Search.Config config = 2;
  }


  message Search.Config {
    string request_id = 1;
    uint32 num = 2;
    float radius = 3;
    float epsilon = 4;
    int64 timeout = 5;
    Filter.Config ingress_filters = 6;
    Filter.Config egress_filters = 7;
    uint32 min_num = 8;
    Search.AggregationAlgorithm aggregation_algorithm = 9;
    google.protobuf.FloatValue ratio = 10;
    uint32 nprobe = 11;
  }


  message Filter.Config {
    repeated Filter.Target targets = 1;
  }


  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }



  message Filter.Config {
    repeated Filter.Target targets = 1;
  }


  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }



  enum  Search.AggregationAlgorithm {
    Unknown = 0;
    ConcurrentQueue = 1;
    SortSlice = 2;
    SortPoolSlice = 3;
    PairingHeap = 4;
  }
  ```

  - Search.Request

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | vector | float | repeated | The vector to be searched. |
    | config | Search.Config |  | The configuration of the search request. |


  - Search.Config

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | request_id | string |  | Unique request ID. |
    | num | uint32 |  | Maximum number of result to be returned. |
    | radius | float |  | Search radius. |
    | epsilon | float |  | Search coefficient. |
    | timeout | int64 |  | Search timeout in nanoseconds. |
    | ingress_filters | Filter.Config |  | Ingress filter configurations. |
    | egress_filters | Filter.Config |  | Egress filter configurations. |
    | min_num | uint32 |  | Minimum number of result to be returned. |
    | aggregation_algorithm | Search.AggregationAlgorithm |  | Aggregation Algorithm |
    | ratio | google.protobuf.FloatValue |  | Search ratio for agent return result number. |
    | nprobe | uint32 |  | Search nprobe. |


  - Filter.Config

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |


  - Filter.Target

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | host | string |  | The target hostname. |
    | port | uint32 |  | The target port. |


  - Filter.Config

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |


  - Filter.Target

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | host | string |  | The target hostname. |
    | port | uint32 |  | The target port. |



### Output

- the scheme of `payload.v1.Search.Response`

  ```rpc
  message Search.Response {
    string request_id = 1;
    repeated Object.Distance results = 2;
  }


  message Object.Distance {
    string id = 1;
    float distance = 2;
  }
  ```

  - Search.Response

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | request_id | string |  | The unique request ID. |
    | results | Object.Distance | repeated | Search results. |


  - Object.Distance

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | id | string |  | The vector ID. |
    | distance | float |  | The distance. |



### Status Code

| code | description       |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  10  | ABORTED           |
|  13  | INTERNAL          |


Please refer to [Response Status Code](../status.md) for more details.



### Troubleshooting

The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name              | common reason                                                                                                   | how to resolve                                                                           |
| :---------------- | :-------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
| CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                 | Check the code, especially around timeout and connection management, and fix if needed.  |
| INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
| DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
| NOT_FOUND         | Search result is empty or insufficient to request result length.                                                | Send a request with another vector or set min_num to a smaller value.                    |
| INTERNAL          | Target Vald cluster or network route has some critical error.                                                   | Check target Vald cluster first and check network route including ingress as second.     |



## SearchByID RPC

SearchByID RPC is the method to search similar vectors using a user-defined vector ID.<br>
The vector with the same requested ID should be indexed into the `vald-agent` before searching.


### Input

- the scheme of `payload.v1.Search.IDRequest`

  ```rpc
  message Search.IDRequest {
    string id = 1;
    Search.Config config = 2;
  }


  message Search.Config {
    string request_id = 1;
    uint32 num = 2;
    float radius = 3;
    float epsilon = 4;
    int64 timeout = 5;
    Filter.Config ingress_filters = 6;
    Filter.Config egress_filters = 7;
    uint32 min_num = 8;
    Search.AggregationAlgorithm aggregation_algorithm = 9;
    google.protobuf.FloatValue ratio = 10;
    uint32 nprobe = 11;
  }


  message Filter.Config {
    repeated Filter.Target targets = 1;
  }


  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }



  message Filter.Config {
    repeated Filter.Target targets = 1;
  }


  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }



  enum  Search.AggregationAlgorithm {
    Unknown = 0;
    ConcurrentQueue = 1;
    SortSlice = 2;
    SortPoolSlice = 3;
    PairingHeap = 4;
  }
  ```

  - Search.IDRequest

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | id | string |  | The vector ID to be searched. |
    | config | Search.Config |  | The configuration of the search request. |


  - Search.Config

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | request_id | string |  | Unique request ID. |
    | num | uint32 |  | Maximum number of result to be returned. |
    | radius | float |  | Search radius. |
    | epsilon | float |  | Search coefficient. |
    | timeout | int64 |  | Search timeout in nanoseconds. |
    | ingress_filters | Filter.Config |  | Ingress filter configurations. |
    | egress_filters | Filter.Config |  | Egress filter configurations. |
    | min_num | uint32 |  | Minimum number of result to be returned. |
    | aggregation_algorithm | Search.AggregationAlgorithm |  | Aggregation Algorithm |
    | ratio | google.protobuf.FloatValue |  | Search ratio for agent return result number. |
    | nprobe | uint32 |  | Search nprobe. |


  - Filter.Config

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |


  - Filter.Target

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | host | string |  | The target hostname. |
    | port | uint32 |  | The target port. |


  - Filter.Config

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |


  - Filter.Target

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | host | string |  | The target hostname. |
    | port | uint32 |  | The target port. |



### Output

- the scheme of `payload.v1.Search.Response`

  ```rpc
  message Search.Response {
    string request_id = 1;
    repeated Object.Distance results = 2;
  }


  message Object.Distance {
    string id = 1;
    float distance = 2;
  }
  ```

  - Search.Response

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | request_id | string |  | The unique request ID. |
    | results | Object.Distance | repeated | Search results. |


  - Object.Distance

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | id | string |  | The vector ID. |
    | distance | float |  | The distance. |



### Status Code

| code | description       |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  10  | ABORTED           |
|  13  | INTERNAL          |


Please refer to [Response Status Code](../status.md) for more details.



### Troubleshooting

The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name              | common reason                                                                                                                    | how to resolve                                                                           |
| :---------------- | :------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
| CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                  | Check the code, especially around timeout and connection management, and fix if needed.  |
| INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                                                          | Check request payload and fix request payload.                                           |
| DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                  | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
| NOT_FOUND         | The Requested ID is not inserted on the target Vald cluster, or the search result is insufficient to the required result length. | Send a request with another vector or set min_num to a smaller value.                    |
| INTERNAL          | Target Vald cluster or network route has some critical error.                                                                    | Check target Vald cluster first and check network route including ingress as second.     |



## StreamSearch RPC

StreamSearch RPC is the method to search vectors with multi queries(vectors) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
Using the bidirectional streaming RPC, the search request can be communicated in any order between the client and server.
Each Search request and response are independent.


### Input

- the scheme of `payload.v1.Search.Request`

  ```rpc
  message Search.Request {
    repeated float vector = 1;
    Search.Config config = 2;
  }


  message Search.Config {
    string request_id = 1;
    uint32 num = 2;
    float radius = 3;
    float epsilon = 4;
    int64 timeout = 5;
    Filter.Config ingress_filters = 6;
    Filter.Config egress_filters = 7;
    uint32 min_num = 8;
    Search.AggregationAlgorithm aggregation_algorithm = 9;
    google.protobuf.FloatValue ratio = 10;
    uint32 nprobe = 11;
  }


  message Filter.Config {
    repeated Filter.Target targets = 1;
  }


  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }



  message Filter.Config {
    repeated Filter.Target targets = 1;
  }


  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }



  enum  Search.AggregationAlgorithm {
    Unknown = 0;
    ConcurrentQueue = 1;
    SortSlice = 2;
    SortPoolSlice = 3;
    PairingHeap = 4;
  }
  ```

  - Search.Request

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | vector | float | repeated | The vector to be searched. |
    | config | Search.Config |  | The configuration of the search request. |


  - Search.Config

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | request_id | string |  | Unique request ID. |
    | num | uint32 |  | Maximum number of result to be returned. |
    | radius | float |  | Search radius. |
    | epsilon | float |  | Search coefficient. |
    | timeout | int64 |  | Search timeout in nanoseconds. |
    | ingress_filters | Filter.Config |  | Ingress filter configurations. |
    | egress_filters | Filter.Config |  | Egress filter configurations. |
    | min_num | uint32 |  | Minimum number of result to be returned. |
    | aggregation_algorithm | Search.AggregationAlgorithm |  | Aggregation Algorithm |
    | ratio | google.protobuf.FloatValue |  | Search ratio for agent return result number. |
    | nprobe | uint32 |  | Search nprobe. |


  - Filter.Config

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |


  - Filter.Target

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | host | string |  | The target hostname. |
    | port | uint32 |  | The target port. |


  - Filter.Config

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |


  - Filter.Target

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | host | string |  | The target hostname. |
    | port | uint32 |  | The target port. |



### Output

- the scheme of `payload.v1.Search.StreamResponse`

  ```rpc
  message Search.StreamResponse {
    Search.Response response = 1;
    google.rpc.Status status = 2;
  }


  message Search.Response {
    string request_id = 1;
    repeated Object.Distance results = 2;
  }


  message Object.Distance {
    string id = 1;
    float distance = 2;
  }
  ```

  - Search.StreamResponse

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | response | Search.Response |  | Represent the search response. |
    | status | google.rpc.Status |  | The RPC error status. |


  - Search.Response

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | request_id | string |  | The unique request ID. |
    | results | Object.Distance | repeated | Search results. |


  - Object.Distance

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | id | string |  | The vector ID. |
    | distance | float |  | The distance. |



### Status Code

| code | description       |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  10  | ABORTED           |
|  13  | INTERNAL          |


Please refer to [Response Status Code](../status.md) for more details.



### Troubleshooting

The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name              | common reason                                                                                                   | how to resolve                                                                           |
| :---------------- | :-------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
| CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                 | Check the code, especially around timeout and connection management, and fix if needed.  |
| INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
| DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
| NOT_FOUND         | Search result is empty or insufficient to request result length.                                                | Send a request with another vector or set min_num to a smaller value.                    |
| INTERNAL          | Target Vald cluster or network route has some critical error.                                                   | Check target Vald cluster first and check network route including ingress as second.     |



## StreamSearchByID RPC

StreamSearchByID RPC is the method to search vectors with multi queries(IDs) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
Using the bidirectional streaming RPC, the search request can be communicated in any order between the client and server.
Each SearchByID request and response are independent.


### Input

- the scheme of `payload.v1.Search.IDRequest`

  ```rpc
  message Search.IDRequest {
    string id = 1;
    Search.Config config = 2;
  }


  message Search.Config {
    string request_id = 1;
    uint32 num = 2;
    float radius = 3;
    float epsilon = 4;
    int64 timeout = 5;
    Filter.Config ingress_filters = 6;
    Filter.Config egress_filters = 7;
    uint32 min_num = 8;
    Search.AggregationAlgorithm aggregation_algorithm = 9;
    google.protobuf.FloatValue ratio = 10;
    uint32 nprobe = 11;
  }


  message Filter.Config {
    repeated Filter.Target targets = 1;
  }


  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }



  message Filter.Config {
    repeated Filter.Target targets = 1;
  }


  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }



  enum  Search.AggregationAlgorithm {
    Unknown = 0;
    ConcurrentQueue = 1;
    SortSlice = 2;
    SortPoolSlice = 3;
    PairingHeap = 4;
  }
  ```

  - Search.IDRequest

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | id | string |  | The vector ID to be searched. |
    | config | Search.Config |  | The configuration of the search request. |


  - Search.Config

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | request_id | string |  | Unique request ID. |
    | num | uint32 |  | Maximum number of result to be returned. |
    | radius | float |  | Search radius. |
    | epsilon | float |  | Search coefficient. |
    | timeout | int64 |  | Search timeout in nanoseconds. |
    | ingress_filters | Filter.Config |  | Ingress filter configurations. |
    | egress_filters | Filter.Config |  | Egress filter configurations. |
    | min_num | uint32 |  | Minimum number of result to be returned. |
    | aggregation_algorithm | Search.AggregationAlgorithm |  | Aggregation Algorithm |
    | ratio | google.protobuf.FloatValue |  | Search ratio for agent return result number. |
    | nprobe | uint32 |  | Search nprobe. |


  - Filter.Config

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |


  - Filter.Target

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | host | string |  | The target hostname. |
    | port | uint32 |  | The target port. |


  - Filter.Config

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |


  - Filter.Target

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | host | string |  | The target hostname. |
    | port | uint32 |  | The target port. |



### Output

- the scheme of `payload.v1.Search.StreamResponse`

  ```rpc
  message Search.StreamResponse {
    Search.Response response = 1;
    google.rpc.Status status = 2;
  }


  message Search.Response {
    string request_id = 1;
    repeated Object.Distance results = 2;
  }


  message Object.Distance {
    string id = 1;
    float distance = 2;
  }
  ```

  - Search.StreamResponse

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | response | Search.Response |  | Represent the search response. |
    | status | google.rpc.Status |  | The RPC error status. |


  - Search.Response

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | request_id | string |  | The unique request ID. |
    | results | Object.Distance | repeated | Search results. |


  - Object.Distance

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | id | string |  | The vector ID. |
    | distance | float |  | The distance. |



### Status Code

| code | description       |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  10  | ABORTED           |
|  13  | INTERNAL          |


Please refer to [Response Status Code](../status.md) for more details.



### Troubleshooting

The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name              | common reason                                                                                                                    | how to resolve                                                                           |
| :---------------- | :------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
| CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                  | Check the code, especially around timeout and connection management, and fix if needed.  |
| INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                                                          | Check request payload and fix request payload.                                           |
| DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                  | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
| NOT_FOUND         | The Requested ID is not inserted on the target Vald cluster, or the search result is insufficient to the required result length. | Send a request with another vector or set min_num to a smaller value.                    |
| INTERNAL          | Target Vald cluster or network route has some critical error.                                                                    | Check target Vald cluster first and check network route including ingress as second.     |



## MultiSearch RPC

MultiSearch RPC is the method to search vectors with multiple vectors in **1** request.

<div class="notice">
gRPC has a message size limitation.<br>
Please be careful that the size of the request exceeds the limit.
</div>


### Input

- the scheme of `payload.v1.Search.MultiRequest`

  ```rpc
  message Search.MultiRequest {
    repeated Search.Request requests = 1;
  }


  message Search.Request {
    repeated float vector = 1;
    Search.Config config = 2;
  }


  message Search.Config {
    string request_id = 1;
    uint32 num = 2;
    float radius = 3;
    float epsilon = 4;
    int64 timeout = 5;
    Filter.Config ingress_filters = 6;
    Filter.Config egress_filters = 7;
    uint32 min_num = 8;
    Search.AggregationAlgorithm aggregation_algorithm = 9;
    google.protobuf.FloatValue ratio = 10;
    uint32 nprobe = 11;
  }


  message Filter.Config {
    repeated Filter.Target targets = 1;
  }


  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }



  message Filter.Config {
    repeated Filter.Target targets = 1;
  }


  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }



  enum  Search.AggregationAlgorithm {
    Unknown = 0;
    ConcurrentQueue = 1;
    SortSlice = 2;
    SortPoolSlice = 3;
    PairingHeap = 4;
  }
  ```

  - Search.MultiRequest

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | requests | Search.Request | repeated | Represent the multiple search request content. |


  - Search.Request

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | vector | float | repeated | The vector to be searched. |
    | config | Search.Config |  | The configuration of the search request. |


  - Search.Config

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | request_id | string |  | Unique request ID. |
    | num | uint32 |  | Maximum number of result to be returned. |
    | radius | float |  | Search radius. |
    | epsilon | float |  | Search coefficient. |
    | timeout | int64 |  | Search timeout in nanoseconds. |
    | ingress_filters | Filter.Config |  | Ingress filter configurations. |
    | egress_filters | Filter.Config |  | Egress filter configurations. |
    | min_num | uint32 |  | Minimum number of result to be returned. |
    | aggregation_algorithm | Search.AggregationAlgorithm |  | Aggregation Algorithm |
    | ratio | google.protobuf.FloatValue |  | Search ratio for agent return result number. |
    | nprobe | uint32 |  | Search nprobe. |


  - Filter.Config

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |


  - Filter.Target

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | host | string |  | The target hostname. |
    | port | uint32 |  | The target port. |


  - Filter.Config

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |


  - Filter.Target

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | host | string |  | The target hostname. |
    | port | uint32 |  | The target port. |



### Output

- the scheme of `payload.v1.Search.Responses`

  ```rpc
  message Search.Responses {
    repeated Search.Response responses = 1;
  }


  message Search.Response {
    string request_id = 1;
    repeated Object.Distance results = 2;
  }


  message Object.Distance {
    string id = 1;
    float distance = 2;
  }
  ```

  - Search.Responses

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | responses | Search.Response | repeated | Represent the multiple search response content. |


  - Search.Response

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | request_id | string |  | The unique request ID. |
    | results | Object.Distance | repeated | Search results. |


  - Object.Distance

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | id | string |  | The vector ID. |
    | distance | float |  | The distance. |



### Status Code

| code | description       |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  10  | ABORTED           |
|  13  | INTERNAL          |


Please refer to [Response Status Code](../status.md) for more details.



### Troubleshooting

  The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name              | common reason                                                                                                   | how to resolve                                                                           |
| :---------------- | :-------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
| CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                 | Check the code, especially around timeout and connection management, and fix if needed.  |
| INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
| DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
| NOT_FOUND         | Search result is empty or insufficient to request result length.                                                | Send a request with another vector or set min_num to a smaller value.                    |
| INTERNAL          | Target Vald cluster or network route has some critical error.                                                   | Check target Vald cluster first and check network route including ingress as second.     |



## MultiSearchByID RPC

MultiSearchByID RPC is the method to search vectors with multiple IDs in **1** request.

<div class="notice">
gRPC has a message size limitation.<br>
Please be careful that the size of the request exceeds the limit.
</div>


### Input

- the scheme of `payload.v1.Search.MultiIDRequest`

  ```rpc
  message Search.MultiIDRequest {
    repeated Search.IDRequest requests = 1;
  }


  message Search.IDRequest {
    string id = 1;
    Search.Config config = 2;
  }


  message Search.Config {
    string request_id = 1;
    uint32 num = 2;
    float radius = 3;
    float epsilon = 4;
    int64 timeout = 5;
    Filter.Config ingress_filters = 6;
    Filter.Config egress_filters = 7;
    uint32 min_num = 8;
    Search.AggregationAlgorithm aggregation_algorithm = 9;
    google.protobuf.FloatValue ratio = 10;
    uint32 nprobe = 11;
  }


  message Filter.Config {
    repeated Filter.Target targets = 1;
  }


  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }



  message Filter.Config {
    repeated Filter.Target targets = 1;
  }


  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }



  enum  Search.AggregationAlgorithm {
    Unknown = 0;
    ConcurrentQueue = 1;
    SortSlice = 2;
    SortPoolSlice = 3;
    PairingHeap = 4;
  }
  ```

  - Search.MultiIDRequest

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | requests | Search.IDRequest | repeated | Represent the multiple search by ID request content. |


  - Search.IDRequest

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | id | string |  | The vector ID to be searched. |
    | config | Search.Config |  | The configuration of the search request. |


  - Search.Config

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | request_id | string |  | Unique request ID. |
    | num | uint32 |  | Maximum number of result to be returned. |
    | radius | float |  | Search radius. |
    | epsilon | float |  | Search coefficient. |
    | timeout | int64 |  | Search timeout in nanoseconds. |
    | ingress_filters | Filter.Config |  | Ingress filter configurations. |
    | egress_filters | Filter.Config |  | Egress filter configurations. |
    | min_num | uint32 |  | Minimum number of result to be returned. |
    | aggregation_algorithm | Search.AggregationAlgorithm |  | Aggregation Algorithm |
    | ratio | google.protobuf.FloatValue |  | Search ratio for agent return result number. |
    | nprobe | uint32 |  | Search nprobe. |


  - Filter.Config

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |


  - Filter.Target

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | host | string |  | The target hostname. |
    | port | uint32 |  | The target port. |


  - Filter.Config

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |


  - Filter.Target

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | host | string |  | The target hostname. |
    | port | uint32 |  | The target port. |



### Output

- the scheme of `payload.v1.Search.Responses`

  ```rpc
  message Search.Responses {
    repeated Search.Response responses = 1;
  }


  message Search.Response {
    string request_id = 1;
    repeated Object.Distance results = 2;
  }


  message Object.Distance {
    string id = 1;
    float distance = 2;
  }
  ```

  - Search.Responses

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | responses | Search.Response | repeated | Represent the multiple search response content. |


  - Search.Response

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | request_id | string |  | The unique request ID. |
    | results | Object.Distance | repeated | Search results. |


  - Object.Distance

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | id | string |  | The vector ID. |
    | distance | float |  | The distance. |



### Status Code

| code | description       |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  10  | ABORTED           |
|  13  | INTERNAL          |


Please refer to [Response Status Code](../status.md) for more details.



### Troubleshooting

The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name              | common reason                                                                                                                    | how to resolve                                                                           |
| :---------------- | :------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
| CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                  | Check the code, especially around timeout and connection management, and fix if needed.  |
| INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                                                          | Check request payload and fix request payload.                                           |
| DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                  | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
| NOT_FOUND         | The Requested ID is not inserted on the target Vald cluster, or the search result is insufficient to the required result length. | Send a request with another vector or set min_num to a smaller value.                    |
| INTERNAL          | Target Vald cluster or network route has some critical error.                                                                    | Check target Vald cluster first and check network route including ingress as second.     |



## LinearSearch RPC

LinearSearch RPC is the method to linear search vector(s) similar to the request vector.


### Input

- the scheme of `payload.v1.Search.Request`

  ```rpc
  message Search.Request {
    repeated float vector = 1;
    Search.Config config = 2;
  }


  message Search.Config {
    string request_id = 1;
    uint32 num = 2;
    float radius = 3;
    float epsilon = 4;
    int64 timeout = 5;
    Filter.Config ingress_filters = 6;
    Filter.Config egress_filters = 7;
    uint32 min_num = 8;
    Search.AggregationAlgorithm aggregation_algorithm = 9;
    google.protobuf.FloatValue ratio = 10;
    uint32 nprobe = 11;
  }


  message Filter.Config {
    repeated Filter.Target targets = 1;
  }


  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }



  message Filter.Config {
    repeated Filter.Target targets = 1;
  }


  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }



  enum  Search.AggregationAlgorithm {
    Unknown = 0;
    ConcurrentQueue = 1;
    SortSlice = 2;
    SortPoolSlice = 3;
    PairingHeap = 4;
  }
  ```

  - Search.Request

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | vector | float | repeated | The vector to be searched. |
    | config | Search.Config |  | The configuration of the search request. |


  - Search.Config

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | request_id | string |  | Unique request ID. |
    | num | uint32 |  | Maximum number of result to be returned. |
    | radius | float |  | Search radius. |
    | epsilon | float |  | Search coefficient. |
    | timeout | int64 |  | Search timeout in nanoseconds. |
    | ingress_filters | Filter.Config |  | Ingress filter configurations. |
    | egress_filters | Filter.Config |  | Egress filter configurations. |
    | min_num | uint32 |  | Minimum number of result to be returned. |
    | aggregation_algorithm | Search.AggregationAlgorithm |  | Aggregation Algorithm |
    | ratio | google.protobuf.FloatValue |  | Search ratio for agent return result number. |
    | nprobe | uint32 |  | Search nprobe. |


  - Filter.Config

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |


  - Filter.Target

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | host | string |  | The target hostname. |
    | port | uint32 |  | The target port. |


  - Filter.Config

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |


  - Filter.Target

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | host | string |  | The target hostname. |
    | port | uint32 |  | The target port. |



### Output

- the scheme of `payload.v1.Search.Response`

  ```rpc
  message Search.Response {
    string request_id = 1;
    repeated Object.Distance results = 2;
  }


  message Object.Distance {
    string id = 1;
    float distance = 2;
  }
  ```

  - Search.Response

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | request_id | string |  | The unique request ID. |
    | results | Object.Distance | repeated | Search results. |


  - Object.Distance

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | id | string |  | The vector ID. |
    | distance | float |  | The distance. |



### Status Code

| code | description       |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  10  | ABORTED           |
|  13  | INTERNAL          |


Please refer to [Response Status Code](../status.md) for more details.



### Troubleshooting

The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name              | common reason                                                                                                   | how to resolve                                                                           |
| :---------------- | :-------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
| CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                 | Check the code, especially around timeout and connection management, and fix if needed.  |
| INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
| DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
| NOT_FOUND         | Search result is empty or insufficient to request result length.                                                | Send a request with another vector or set min_num to a smaller value.                    |
| INTERNAL          | Target Vald cluster or network route has some critical error.                                                   | Check target Vald cluster first and check network route including ingress as second.     |



## LinearSearchByID RPC

LinearSearchByID RPC is the method to linear search similar vectors using a user-defined vector ID.<br>
The vector with the same requested ID should be indexed into the `vald-agent` before searching.
You will get a `NOT_FOUND` error if the vector isn't stored.


### Input

- the scheme of `payload.v1.Search.IDRequest`

  ```rpc
  message Search.IDRequest {
    string id = 1;
    Search.Config config = 2;
  }


  message Search.Config {
    string request_id = 1;
    uint32 num = 2;
    float radius = 3;
    float epsilon = 4;
    int64 timeout = 5;
    Filter.Config ingress_filters = 6;
    Filter.Config egress_filters = 7;
    uint32 min_num = 8;
    Search.AggregationAlgorithm aggregation_algorithm = 9;
    google.protobuf.FloatValue ratio = 10;
    uint32 nprobe = 11;
  }


  message Filter.Config {
    repeated Filter.Target targets = 1;
  }


  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }



  message Filter.Config {
    repeated Filter.Target targets = 1;
  }


  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }



  enum  Search.AggregationAlgorithm {
    Unknown = 0;
    ConcurrentQueue = 1;
    SortSlice = 2;
    SortPoolSlice = 3;
    PairingHeap = 4;
  }
  ```

  - Search.IDRequest

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | id | string |  | The vector ID to be searched. |
    | config | Search.Config |  | The configuration of the search request. |


  - Search.Config

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | request_id | string |  | Unique request ID. |
    | num | uint32 |  | Maximum number of result to be returned. |
    | radius | float |  | Search radius. |
    | epsilon | float |  | Search coefficient. |
    | timeout | int64 |  | Search timeout in nanoseconds. |
    | ingress_filters | Filter.Config |  | Ingress filter configurations. |
    | egress_filters | Filter.Config |  | Egress filter configurations. |
    | min_num | uint32 |  | Minimum number of result to be returned. |
    | aggregation_algorithm | Search.AggregationAlgorithm |  | Aggregation Algorithm |
    | ratio | google.protobuf.FloatValue |  | Search ratio for agent return result number. |
    | nprobe | uint32 |  | Search nprobe. |


  - Filter.Config

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |


  - Filter.Target

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | host | string |  | The target hostname. |
    | port | uint32 |  | The target port. |


  - Filter.Config

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |


  - Filter.Target

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | host | string |  | The target hostname. |
    | port | uint32 |  | The target port. |



### Output

- the scheme of `payload.v1.Search.Response`

  ```rpc
  message Search.Response {
    string request_id = 1;
    repeated Object.Distance results = 2;
  }


  message Object.Distance {
    string id = 1;
    float distance = 2;
  }
  ```

  - Search.Response

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | request_id | string |  | The unique request ID. |
    | results | Object.Distance | repeated | Search results. |


  - Object.Distance

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | id | string |  | The vector ID. |
    | distance | float |  | The distance. |



### Status Code

| code | description       |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  10  | ABORTED           |
|  13  | INTERNAL          |


Please refer to [Response Status Code](../status.md) for more details.



### Troubleshooting

The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name              | common reason                                                                                                                    | how to resolve                                                                           |
| :---------------- | :------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
| CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                  | Check the code, especially around timeout and connection management, and fix if needed.  |
| INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                                                          | Check request payload and fix request payload.                                           |
| DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                  | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
| NOT_FOUND         | The Requested ID is not inserted on the target Vald cluster, or the search result is insufficient to the required result length. | Send a request with another vector or set min_num to a smaller value.                    |
| INTERNAL          | Target Vald cluster or network route has some critical error.                                                                    | Check target Vald cluster first and check network route including ingress as second.     |



## StreamLinearSearch RPC

StreamLinearSearch RPC is the method to linear search vectors with multi queries(vectors) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
Using the bidirectional streaming RPC, the linear search request can be communicated in any order between the client and server.
Each LinearSearch request and response are independent.


### Input

- the scheme of `payload.v1.Search.Request`

  ```rpc
  message Search.Request {
    repeated float vector = 1;
    Search.Config config = 2;
  }


  message Search.Config {
    string request_id = 1;
    uint32 num = 2;
    float radius = 3;
    float epsilon = 4;
    int64 timeout = 5;
    Filter.Config ingress_filters = 6;
    Filter.Config egress_filters = 7;
    uint32 min_num = 8;
    Search.AggregationAlgorithm aggregation_algorithm = 9;
    google.protobuf.FloatValue ratio = 10;
    uint32 nprobe = 11;
  }


  message Filter.Config {
    repeated Filter.Target targets = 1;
  }


  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }



  message Filter.Config {
    repeated Filter.Target targets = 1;
  }


  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }



  enum  Search.AggregationAlgorithm {
    Unknown = 0;
    ConcurrentQueue = 1;
    SortSlice = 2;
    SortPoolSlice = 3;
    PairingHeap = 4;
  }
  ```

  - Search.Request

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | vector | float | repeated | The vector to be searched. |
    | config | Search.Config |  | The configuration of the search request. |


  - Search.Config

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | request_id | string |  | Unique request ID. |
    | num | uint32 |  | Maximum number of result to be returned. |
    | radius | float |  | Search radius. |
    | epsilon | float |  | Search coefficient. |
    | timeout | int64 |  | Search timeout in nanoseconds. |
    | ingress_filters | Filter.Config |  | Ingress filter configurations. |
    | egress_filters | Filter.Config |  | Egress filter configurations. |
    | min_num | uint32 |  | Minimum number of result to be returned. |
    | aggregation_algorithm | Search.AggregationAlgorithm |  | Aggregation Algorithm |
    | ratio | google.protobuf.FloatValue |  | Search ratio for agent return result number. |
    | nprobe | uint32 |  | Search nprobe. |


  - Filter.Config

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |


  - Filter.Target

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | host | string |  | The target hostname. |
    | port | uint32 |  | The target port. |


  - Filter.Config

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |


  - Filter.Target

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | host | string |  | The target hostname. |
    | port | uint32 |  | The target port. |



### Output

- the scheme of `payload.v1.Search.StreamResponse`

  ```rpc
  message Search.StreamResponse {
    Search.Response response = 1;
    google.rpc.Status status = 2;
  }


  message Search.Response {
    string request_id = 1;
    repeated Object.Distance results = 2;
  }


  message Object.Distance {
    string id = 1;
    float distance = 2;
  }
  ```

  - Search.StreamResponse

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | response | Search.Response |  | Represent the search response. |
    | status | google.rpc.Status |  | The RPC error status. |


  - Search.Response

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | request_id | string |  | The unique request ID. |
    | results | Object.Distance | repeated | Search results. |


  - Object.Distance

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | id | string |  | The vector ID. |
    | distance | float |  | The distance. |



### Status Code

| code | description       |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  10  | ABORTED           |
|  13  | INTERNAL          |


Please refer to [Response Status Code](../status.md) for more details.



### Troubleshooting

The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name              | common reason                                                                                                   | how to resolve                                                                           |
| :---------------- | :-------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
| CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                 | Check the code, especially around timeout and connection management, and fix if needed.  |
| INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
| DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
| NOT_FOUND         | Search result is empty or insufficient to request result length.                                                | Send a request with another vector or set min_num to a smaller value.                    |
| INTERNAL          | Target Vald cluster or network route has some critical error.                                                   | Check target Vald cluster first and check network route including ingress as second.     |



## StreamLinearSearchByID RPC

  StreamLinearSearchByID RPC is the method to linear search vectors with multi queries(IDs) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
Using the bidirectional streaming RPC, the linear search request can be communicated in any order between the client and server.
Each LinearSearchByID request and response are independent.


### Input

- the scheme of `payload.v1.Search.IDRequest`

  ```rpc
  message Search.IDRequest {
    string id = 1;
    Search.Config config = 2;
  }


  message Search.Config {
    string request_id = 1;
    uint32 num = 2;
    float radius = 3;
    float epsilon = 4;
    int64 timeout = 5;
    Filter.Config ingress_filters = 6;
    Filter.Config egress_filters = 7;
    uint32 min_num = 8;
    Search.AggregationAlgorithm aggregation_algorithm = 9;
    google.protobuf.FloatValue ratio = 10;
    uint32 nprobe = 11;
  }


  message Filter.Config {
    repeated Filter.Target targets = 1;
  }


  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }



  message Filter.Config {
    repeated Filter.Target targets = 1;
  }


  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }



  enum  Search.AggregationAlgorithm {
    Unknown = 0;
    ConcurrentQueue = 1;
    SortSlice = 2;
    SortPoolSlice = 3;
    PairingHeap = 4;
  }
  ```

  - Search.IDRequest

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | id | string |  | The vector ID to be searched. |
    | config | Search.Config |  | The configuration of the search request. |


  - Search.Config

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | request_id | string |  | Unique request ID. |
    | num | uint32 |  | Maximum number of result to be returned. |
    | radius | float |  | Search radius. |
    | epsilon | float |  | Search coefficient. |
    | timeout | int64 |  | Search timeout in nanoseconds. |
    | ingress_filters | Filter.Config |  | Ingress filter configurations. |
    | egress_filters | Filter.Config |  | Egress filter configurations. |
    | min_num | uint32 |  | Minimum number of result to be returned. |
    | aggregation_algorithm | Search.AggregationAlgorithm |  | Aggregation Algorithm |
    | ratio | google.protobuf.FloatValue |  | Search ratio for agent return result number. |
    | nprobe | uint32 |  | Search nprobe. |


  - Filter.Config

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |


  - Filter.Target

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | host | string |  | The target hostname. |
    | port | uint32 |  | The target port. |


  - Filter.Config

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |


  - Filter.Target

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | host | string |  | The target hostname. |
    | port | uint32 |  | The target port. |



### Output

- the scheme of `payload.v1.Search.StreamResponse`

  ```rpc
  message Search.StreamResponse {
    Search.Response response = 1;
    google.rpc.Status status = 2;
  }


  message Search.Response {
    string request_id = 1;
    repeated Object.Distance results = 2;
  }


  message Object.Distance {
    string id = 1;
    float distance = 2;
  }
  ```

  - Search.StreamResponse

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | response | Search.Response |  | Represent the search response. |
    | status | google.rpc.Status |  | The RPC error status. |


  - Search.Response

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | request_id | string |  | The unique request ID. |
    | results | Object.Distance | repeated | Search results. |


  - Object.Distance

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | id | string |  | The vector ID. |
    | distance | float |  | The distance. |



### Status Code

| code | description       |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  10  | ABORTED           |
|  13  | INTERNAL          |


Please refer to [Response Status Code](../status.md) for more details.



### Troubleshooting

The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name              | common reason                                                                                                                    | how to resolve                                                                           |
| :---------------- | :------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
| CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                  | Check the code, especially around timeout and connection management, and fix if needed.  |
| INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                                                          | Check request payload and fix request payload.                                           |
| DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                  | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
| NOT_FOUND         | The Requested ID is not inserted on the target Vald cluster, or the search result is insufficient to the required result length. | Send a request with another vector or set min_num to a smaller value.                    |
| INTERNAL          | Target Vald cluster or network route has some critical error.                                                                    | Check target Vald cluster first and check network route including ingress as second.     |



## MultiLinearSearch RPC

MultiLinearSearch RPC is the method to linear search vectors with multiple vectors in **1** request.

<div class="notice">
gRPC has a message size limitation.<br>
Please be careful that the size of the request exceeds the limit.
</div>


### Input

- the scheme of `payload.v1.Search.MultiRequest`

  ```rpc
  message Search.MultiRequest {
    repeated Search.Request requests = 1;
  }


  message Search.Request {
    repeated float vector = 1;
    Search.Config config = 2;
  }


  message Search.Config {
    string request_id = 1;
    uint32 num = 2;
    float radius = 3;
    float epsilon = 4;
    int64 timeout = 5;
    Filter.Config ingress_filters = 6;
    Filter.Config egress_filters = 7;
    uint32 min_num = 8;
    Search.AggregationAlgorithm aggregation_algorithm = 9;
    google.protobuf.FloatValue ratio = 10;
    uint32 nprobe = 11;
  }


  message Filter.Config {
    repeated Filter.Target targets = 1;
  }


  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }



  message Filter.Config {
    repeated Filter.Target targets = 1;
  }


  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }



  enum  Search.AggregationAlgorithm {
    Unknown = 0;
    ConcurrentQueue = 1;
    SortSlice = 2;
    SortPoolSlice = 3;
    PairingHeap = 4;
  }
  ```

  - Search.MultiRequest

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | requests | Search.Request | repeated | Represent the multiple search request content. |


  - Search.Request

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | vector | float | repeated | The vector to be searched. |
    | config | Search.Config |  | The configuration of the search request. |


  - Search.Config

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | request_id | string |  | Unique request ID. |
    | num | uint32 |  | Maximum number of result to be returned. |
    | radius | float |  | Search radius. |
    | epsilon | float |  | Search coefficient. |
    | timeout | int64 |  | Search timeout in nanoseconds. |
    | ingress_filters | Filter.Config |  | Ingress filter configurations. |
    | egress_filters | Filter.Config |  | Egress filter configurations. |
    | min_num | uint32 |  | Minimum number of result to be returned. |
    | aggregation_algorithm | Search.AggregationAlgorithm |  | Aggregation Algorithm |
    | ratio | google.protobuf.FloatValue |  | Search ratio for agent return result number. |
    | nprobe | uint32 |  | Search nprobe. |


  - Filter.Config

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |


  - Filter.Target

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | host | string |  | The target hostname. |
    | port | uint32 |  | The target port. |


  - Filter.Config

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |


  - Filter.Target

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | host | string |  | The target hostname. |
    | port | uint32 |  | The target port. |



### Output

- the scheme of `payload.v1.Search.Responses`

  ```rpc
  message Search.Responses {
    repeated Search.Response responses = 1;
  }


  message Search.Response {
    string request_id = 1;
    repeated Object.Distance results = 2;
  }


  message Object.Distance {
    string id = 1;
    float distance = 2;
  }
  ```

  - Search.Responses

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | responses | Search.Response | repeated | Represent the multiple search response content. |


  - Search.Response

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | request_id | string |  | The unique request ID. |
    | results | Object.Distance | repeated | Search results. |


  - Object.Distance

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | id | string |  | The vector ID. |
    | distance | float |  | The distance. |



### Status Code

| code | description       |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  10  | ABORTED           |
|  13  | INTERNAL          |


Please refer to [Response Status Code](../status.md) for more details.



### Troubleshooting

  The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name              | common reason                                                                                                   | how to resolve                                                                           |
| :---------------- | :-------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
| CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                 | Check the code, especially around timeout and connection management, and fix if needed.  |
| INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
| DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
| NOT_FOUND         | Search result is empty or insufficient to request result length.                                                | Send a request with another vector or set min_num to a smaller value.                    |
| INTERNAL          | Target Vald cluster or network route has some critical error.                                                   | Check target Vald cluster first and check network route including ingress as second.     |



## MultiLinearSearchByID RPC

MultiLinearSearchByID RPC is the method to linear search vectors with multiple IDs in **1** request.

<div class="notice">
gRPC has a message size limitation.<br>
Please be careful that the size of the request exceeds the limit.
</div>
// 

### Input

- the scheme of `payload.v1.Search.MultiIDRequest`

  ```rpc
  message Search.MultiIDRequest {
    repeated Search.IDRequest requests = 1;
  }


  message Search.IDRequest {
    string id = 1;
    Search.Config config = 2;
  }


  message Search.Config {
    string request_id = 1;
    uint32 num = 2;
    float radius = 3;
    float epsilon = 4;
    int64 timeout = 5;
    Filter.Config ingress_filters = 6;
    Filter.Config egress_filters = 7;
    uint32 min_num = 8;
    Search.AggregationAlgorithm aggregation_algorithm = 9;
    google.protobuf.FloatValue ratio = 10;
    uint32 nprobe = 11;
  }


  message Filter.Config {
    repeated Filter.Target targets = 1;
  }


  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }



  message Filter.Config {
    repeated Filter.Target targets = 1;
  }


  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }



  enum  Search.AggregationAlgorithm {
    Unknown = 0;
    ConcurrentQueue = 1;
    SortSlice = 2;
    SortPoolSlice = 3;
    PairingHeap = 4;
  }
  ```

  - Search.MultiIDRequest

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | requests | Search.IDRequest | repeated | Represent the multiple search by ID request content. |


  - Search.IDRequest

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | id | string |  | The vector ID to be searched. |
    | config | Search.Config |  | The configuration of the search request. |


  - Search.Config

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | request_id | string |  | Unique request ID. |
    | num | uint32 |  | Maximum number of result to be returned. |
    | radius | float |  | Search radius. |
    | epsilon | float |  | Search coefficient. |
    | timeout | int64 |  | Search timeout in nanoseconds. |
    | ingress_filters | Filter.Config |  | Ingress filter configurations. |
    | egress_filters | Filter.Config |  | Egress filter configurations. |
    | min_num | uint32 |  | Minimum number of result to be returned. |
    | aggregation_algorithm | Search.AggregationAlgorithm |  | Aggregation Algorithm |
    | ratio | google.protobuf.FloatValue |  | Search ratio for agent return result number. |
    | nprobe | uint32 |  | Search nprobe. |


  - Filter.Config

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |


  - Filter.Target

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | host | string |  | The target hostname. |
    | port | uint32 |  | The target port. |


  - Filter.Config

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |


  - Filter.Target

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | host | string |  | The target hostname. |
    | port | uint32 |  | The target port. |



### Output

- the scheme of `payload.v1.Search.Responses`

  ```rpc
  message Search.Responses {
    repeated Search.Response responses = 1;
  }


  message Search.Response {
    string request_id = 1;
    repeated Object.Distance results = 2;
  }


  message Object.Distance {
    string id = 1;
    float distance = 2;
  }
  ```

  - Search.Responses

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | responses | Search.Response | repeated | Represent the multiple search response content. |


  - Search.Response

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | request_id | string |  | The unique request ID. |
    | results | Object.Distance | repeated | Search results. |


  - Object.Distance

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | id | string |  | The vector ID. |
    | distance | float |  | The distance. |



### Status Code

| code | description       |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  10  | ABORTED           |
|  13  | INTERNAL          |


Please refer to [Response Status Code](../status.md) for more details.



### Troubleshooting

The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name              | common reason                                                                                                                    | how to resolve                                                                           |
| :---------------- | :------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
| CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                  | Check the code, especially around timeout and connection management, and fix if needed.  |
| INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                                                          | Check request payload and fix request payload.                                           |
| DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                  | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
| NOT_FOUND         | The Requested ID is not inserted on the target Vald cluster, or the search result is insufficient to the required result length. | Send a request with another vector or set min_num to a smaller value.                    |
| INTERNAL          | Target Vald cluster or network route has some critical error.                                                                    | Check target Vald cluster first and check network route including ingress as second.     |



