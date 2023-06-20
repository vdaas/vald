# Vald Search APIs

## Overview

Search Service is responsible for searching vectors from `vald-agent` which are similar to user request vector.

```rpc
service Search {
  rpc Search(payload.v1.Search.Request) returns (payload.v1.Search.Response) {}

  rpc SearchByID(payload.v1.Search.IDRequest)
      returns (payload.v1.Search.Response) {}

  rpc StreamSearch(stream payload.v1.Search.Request)
      returns (stream payload.v1.Search.StreamResponse) {}

  rpc StreamSearchByID(stream payload.v1.Search.IDRequest)
      returns (stream payload.v1.Search.StreamResponse) {}

  rpc MultiSearch(payload.v1.Search.MultiRequest)
      returns (payload.v1.Search.Responses) {}

  rpc MultiSearchByID(payload.v1.Search.MultiIDRequest)
      returns (payload.v1.Search.Responses) {}

  rpc LinearSearch(payload.v1.Search.Request) returns (payload.v1.Search.Response) {}

  rpc LinearSearchByID(payload.v1.Search.IDRequest)
      returns (payload.v1.Search.Response) {}

  rpc StreamLinearSearch(stream payload.v1.Search.Request)
      returns (stream payload.v1.Search.StreamResponse) {}

  rpc StreamLinearSearchByID(stream payload.v1.Search.IDRequest)
      returns (stream payload.v1.Search.StreamResponse) {}

  rpc MultiLinearSearch(payload.v1.Search.MultiRequest)
      returns (payload.v1.Search.Responses) {}

  rpc MultiLinearSearchByID(payload.v1.Search.MultiIDRequest)
      returns (payload.v1.Search.Responses) {}
}
```

## Search RPC

Search RPC is the method to search vector(s) similar to request vector.

### Input

- the scheme of `payload.v1.Search.Request`

  ```rpc
  message Search {
    message Request {
      repeated float vector = 1 [ (validate.rules).repeated .min_items = 2 ];
      Config config = 2;
    }

    message Config {
      string request_id = 1;
      uint32 num = 2 [ (validate.rules).uint32.gte = 1 ];
      float radius = 3;
      float epsilon = 4;
      int64 timeout = 5;
      Filter.Config ingress_filters = 6;
      Filter.Config egress_filters = 7;
      uint32 min_num = 8;
      AggregationAlgorithm aggregation_algorithm = 9;
    }
  }

  enum AggregationAlgorithm {
    Unknown = 0;
    ConcurrentQueue = 1;
    SortSlice = 2;
    SortPoolSlice = 3;
    PairingHeap = 4;
  }
  ```

  - Search.Request

    | field  | type   | label                  | required | desc.                                                   |
    | :----: | :----- | :--------------------- | :------: | :------------------------------------------------------ |
    | vector | float  | repeated(Array[float]) |    \*    | the vector data. its dimension is between 2 and 65,536. |
    | config | Config |                        |    \*    | the configuration of the search request                 |

  - Search.Config

    |         field         | type                 | label | required | desc.                                                                        |
    | :-------------------: | :------------------- | :---- | :------: | :--------------------------------------------------------------------------- |
    |      request_id       | string               |       |          | unique request ID                                                            |
    |          num          | uint32               |       |    \*    | the maximum number of result to be returned                                  |
    |        radius         | float                |       |    \*    | the search radius                                                            |
    |        epsilon        | float                |       |    \*    | the search coefficient (default value is `0.1`)                              |
    |        timeout        | int64                |       |          | Search timeout in nanoseconds (default value is `5s`)                        |
    |    ingress_filters    | Filter.Config        |       |          | Ingress Filter configuration                                                 |
    |    egress_filters     | Filter.Config        |       |          | Egress Filter configuration                                                  |
    |        min_num        | uint32               |       |          | the minimum number of result to be returned                                  |
    | aggregation_algorithm | AggregationAlgorithm |       |          | the search aggregation algorithm option (default value is `ConcurrentQueue`) |

### Output

- the scheme of `payload.v1.Search.Response`.

  ```rpc
  message Search {
    message Response {
      string request_id = 1;
      repeated Object.Distance results = 2;
    }
  }

  message Object {
    message Distance {
      string id = 1;
      float distance = 2;
    }
  }
  ```

  - Search.Response

    |   field    | type            | label                            | desc.                 |
    | :--------: | :-------------- | :------------------------------- | :-------------------- |
    | request_id | string          |                                  | the unique request ID |
    |  results   | Object.Distance | repeated(Array[Object.Distance]) | search results        |

  - Object.Distance

    |  field   | type   | label | desc.                                                 |
    | :------: | :----- | :---- | :---------------------------------------------------- |
    |    id    | string |       | the vector ID                                         |
    | distance | float  |       | the distance between result vector and request vector |

### Status Code

| code | desc.             |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  13  | INTERNAL          |

## SearchByID RPC

SearchByID RPC is the method to search similar vectors using a user-defined vector ID.<br>
The vector with the same requested ID should be indexed into the `vald-agent` before searching.

### Input

- the scheme of `payload.v1.Search.IDRequest`

  ```rpc
  message Search {
    message IDRequest {
      string id = 1;
      Config config = 2;
    }

    message Config {
      string request_id = 1;
      uint32 num = 2 [ (validate.rules).uint32.gte = 1 ];
      float radius = 3;
      float epsilon = 4;
      int64 timeout = 5;
      Filter.Config ingress_filters = 6;
      Filter.Config egress_filters = 7;
      uint32 min_num = 8;
      AggregationAlgorithm aggregation_algorithm = 9;
    }
  }

  enum AggregationAlgorithm {
    Unknown = 0;
    ConcurrentQueue = 1;
    SortSlice = 2;
    SortPoolSlice = 3;
    PairingHeap = 4;
  }
  ```

  - Search.IDRequest

    | field  | type   | label | required | desc.                                   |
    | :----: | :----- | :---- | :------: | :-------------------------------------- |
    |   id   | string |       |    \*    | the vector ID to be searched            |
    | config | Config |       |    \*    | the configuration of the search request |

  - Search.Config

    |         field         | type                 | label | required | desc.                                                                        |
    | :-------------------: | :------------------- | :---- | :------: | :--------------------------------------------------------------------------- |
    |      request_id       | string               |       |          | unique request ID                                                            |
    |          num          | uint32               |       |    \*    | the maximum number of result to be returned                                  |
    |        radius         | float                |       |    \*    | the search radius                                                            |
    |        epsilon        | float                |       |    \*    | the search coefficient (default value is `0.1`)                              |
    |        timeout        | int64                |       |          | Search timeout in nanoseconds (default value is `5s`)                        |
    |    ingress_filters    | Filter.Config        |       |          | Ingress Filter configuration                                                 |
    |    egress_filters     | Filter.Config        |       |          | Egress Filter configuration                                                  |
    |        min_num        | uint32               |       |          | the minimum number of result to be returned                                  |
    | aggregation_algorithm | AggregationAlgorithm |       |          | the search aggregation algorithm option (default value is `ConcurrentQueue`) |

### Output

- the scheme of `payload.v1.Search.Response`.

  ```rpc
  message Search {
    message Response {
      string request_id = 1;
      repeated Object.Distance results = 2;
    }
  }

  message Object {
    message Distance {
      string id = 1;
      float distance = 2;
    }
  }
  ```

  - Search.Response

    |   field    | type            | label                            | desc.                 |
    | :--------: | :-------------- | :------------------------------- | :-------------------- |
    | request_id | string          |                                  | the unique request ID |
    |  results   | Object.Distance | repeated(Array[Object.Distance]) | search results        |

  - Object.Distance

    |  field   | type   | label | desc.                                                 |
    | :------: | :----- | :---- | :---------------------------------------------------- |
    |    id    | string |       | the vector ID                                         |
    | distance | float  |       | the distance between result vector and request vector |

### Status Code

| code | desc.             |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  13  | INTERNAL          |

## StreamSearch RPC

StreamSearch RPC is the method to search vectors with multi queries(vectors) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
By using the bidirectional streaming RPC, the search request can be communicated in any order between client and server.
Each Search request and response are independent.

### Input

- the scheme of `payload.v1.Search.Request stream`

  ```rpc
  message Search {
    message Request {
      repeated float vector = 1 [ (validate.rules).repeated .min_items = 2 ];
      Config config = 2;
    }

    message Config {
      string request_id = 1;
      uint32 num = 2 [ (validate.rules).uint32.gte = 1 ];
      float radius = 3;
      float epsilon = 4;
      int64 timeout = 5;
      Filter.Config ingress_filters = 6;
      Filter.Config egress_filters = 7;
      uint32 min_num = 8;
      AggregationAlgorithm aggregation_algorithm = 9;
    }
  }

  enum AggregationAlgorithm {
    Unknown = 0;
    ConcurrentQueue = 1;
    SortSlice = 2;
    SortPoolSlice = 3;
    PairingHeap = 4;
  }
  ```

  - Search.Request

    | field  | type   | label                  | required | desc.                                                   |
    | :----: | :----- | :--------------------- | :------: | :------------------------------------------------------ |
    | vector | float  | repeated(Array[float]) |    \*    | the vector data. its dimension is between 2 and 65,536. |
    | config | Config |                        |    \*    | the configuration of the search request                 |

  - Search.Config

    |         field         | type                 | label | required | desc.                                                                        |
    | :-------------------: | :------------------- | :---- | :------: | :--------------------------------------------------------------------------- |
    |      request_id       | string               |       |          | unique request ID                                                            |
    |          num          | uint32               |       |    \*    | the maximum number of result to be returned                                  |
    |        radius         | float                |       |    \*    | the search radius                                                            |
    |        epsilon        | float                |       |    \*    | the search coefficient (default value is `0.1`)                              |
    |        timeout        | int64                |       |          | Search timeout in nanoseconds (default value is `5s`)                        |
    |    ingress_filters    | Filter.Config        |       |          | Ingress Filter configuration                                                 |
    |    egress_filters     | Filter.Config        |       |          | Egress Filter configuration                                                  |
    |        min_num        | uint32               |       |          | the minimum number of result to be returned                                  |
    | aggregation_algorithm | AggregationAlgorithm |       |          | the search aggregation algorithm option (default value is `ConcurrentQueue`) |

### Output

- the scheme of `payload.v1.Search.StreamResponse`.

  ```rpc
  message Search {
    message StreamResponse {
      oneof payload {
        Response response = 1;
        google.rpc.Status status = 2;
      }
    }

    message Response {
      string request_id = 1;
      repeated Object.Distance results = 2;
    }
  }

  message Object {
    message Distance {
      string id = 1;
      float distance = 2;
    }
  }
  ```

  - Search.StreamResponse

    |  field   | type              | label | desc.                      |
    | :------: | :---------------- | :---- | :------------------------- |
    | response | Response          |       | the search result response |
    |  status  | google.rpc.Status |       | the status of google RPC   |

  - Search.Response

    |   field    | type            | label                            | desc.                 |
    | :--------: | :-------------- | :------------------------------- | :-------------------- |
    | request_id | string          |                                  | the unique request ID |
    |  results   | Object.Distance | repeated(Array[Object.Distance]) | search results        |

  - Object.Distance

    |  field   | type   | label | desc.                                                 |
    | :------: | :----- | :---- | :---------------------------------------------------- |
    |    id    | string |       | the vector ID                                         |
    | distance | float  |       | the distance between result vector and request vector |

### Status Code

| code | desc.             |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  13  | INTERNAL          |

## StreamSearchByID RPC

StreamSearchByID RPC is the method to search vectors with multi queries(IDs) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
By using the bidirectional streaming RPC, the search request can be communicated in any order between client and server.
Each SearchByID request and response are independent.

### Input

- the scheme of `payload.v1.Search.IDRequest stream`

  ```rpc
  message Search {
    message IDRequest {
      string id = 1;
      Config config = 2;
    }

    message Config {
      string request_id = 1;
      uint32 num = 2 [ (validate.rules).uint32.gte = 1 ];
      float radius = 3;
      float epsilon = 4;
      int64 timeout = 5;
      Filter.Config ingress_filters = 6;
      Filter.Config egress_filters = 7;
      uint32 min_num = 8;
      AggregationAlgorithm aggregation_algorithm = 9;
    }
  }

  enum AggregationAlgorithm {
    Unknown = 0;
    ConcurrentQueue = 1;
    SortSlice = 2;
    SortPoolSlice = 3;
    PairingHeap = 4;
  }
  ```

  - Search.IDRequest

    | field  | type   | label | required | desc.                                   |
    | :----: | :----- | :---- | :------: | :-------------------------------------- |
    |   id   | string |       |    \*    | the vector ID to be searched            |
    | config | Config |       |    \*    | the configuration of the search request |

  - Search.Config

    |         field         | type                 | label | required | desc.                                                                        |
    | :-------------------: | :------------------- | :---- | :------: | :--------------------------------------------------------------------------- |
    |      request_id       | string               |       |          | unique request ID                                                            |
    |          num          | uint32               |       |    \*    | the maximum number of result to be returned                                  |
    |        radius         | float                |       |    \*    | the search radius                                                            |
    |        epsilon        | float                |       |    \*    | the search coefficient (default value is `0.1`)                              |
    |        timeout        | int64                |       |          | Search timeout in nanoseconds (default value is `5s`)                        |
    |    ingress_filters    | Filter.Config        |       |          | Ingress Filter configuration                                                 |
    |    egress_filters     | Filter.Config        |       |          | Egress Filter configuration                                                  |
    |        min_num        | uint32               |       |          | the minimum number of result to be returned                                  |
    | aggregation_algorithm | AggregationAlgorithm |       |          | the search aggregation algorithm option (default value is `ConcurrentQueue`) |

### Output

- the scheme of `payload.v1.Search.StreamResponse`.

  ```rpc
  message Search {
    message StreamResponse {
      oneof payload {
        Response response = 1;
        google.rpc.Status status = 2;
      }
    }

    message Response {
      string request_id = 1;
      repeated Object.Distance results = 2;
    }
  }

  message Object {
    message Distance {
      string id = 1;
      float distance = 2;
    }
  }
  ```

  - Search.StreamResponse

    |  field   | type              | label | desc.                      |
    | :------: | :---------------- | :---- | :------------------------- |
    | response | Response          |       | the search result response |
    |  status  | google.rpc.Status |       | the status of google RPC   |

  - Search.Response

    |   field    | type            | label                            | desc.                 |
    | :--------: | :-------------- | :------------------------------- | :-------------------- |
    | request_id | string          |                                  | the unique request ID |
    |  results   | Object.Distance | repeated(Array[Object.Distance]) | search results        |

  - Object.Distance

    |  field   | type   | label | desc.                                                 |
    | :------: | :----- | :---- | :---------------------------------------------------- |
    |    id    | string |       | the vector ID                                         |
    | distance | float  |       | the distance between result vector and request vector |

### Status Code

| code | desc.             |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  13  | INTERNAL          |

## MultiSearch RPC

MultiSearch RPC is the method to search vectors with multiple vectors in **1** request.

<div class="notice">
gRPC has a message size limitation.<br>
Please be careful that the size of the request exceeds the limit.
</div>

### Input

- the scheme of `payload.v1.Search.MultiRequest`

  ```rpc
  message Search {
    message MultiRequest {
      repeated Request requests = 1;
    }

    message Request {
      repeated float vector = 1 [ (validate.rules).repeated .min_items = 2 ];
      Config config = 2;
    }

    message Config {
      string request_id = 1;
      uint32 num = 2 [ (validate.rules).uint32.gte = 1 ];
      float radius = 3;
      float epsilon = 4;
      int64 timeout = 5;
      Filter.Config ingress_filters = 6;
      Filter.Config egress_filters = 7;
      uint32 min_num = 8;
      AggregationAlgorithm aggregation_algorithm = 9;
    }
  }

  enum AggregationAlgorithm {
    Unknown = 0;
    ConcurrentQueue = 1;
    SortSlice = 2;
    SortPoolSlice = 3;
    PairingHeap = 4;
  }
  ```

  - Search.MultiRequest

    |  field   | type                     | label | required | desc.                   |
    | :------: | :----------------------- | :---- | :------: | :---------------------- |
    | requests | repeated(Array[Request]) |       |    \*    | the search request list |

  - Search.Request

    | field  | type   | label                  | required | desc.                                                   |
    | :----: | :----- | :--------------------- | :------: | :------------------------------------------------------ |
    | vector | float  | repeated(Array[float]) |    \*    | the vector data. its dimension is between 2 and 65,536. |
    | config | Config |                        |    \*    | the configuration of the search request                 |

  - Search.Config

    |         field         | type                 | label | required | desc.                                                                        |
    | :-------------------: | :------------------- | :---- | :------: | :--------------------------------------------------------------------------- |
    |      request_id       | string               |       |          | unique request ID                                                            |
    |          num          | uint32               |       |    \*    | the maximum number of result to be returned                                  |
    |        radius         | float                |       |    \*    | the search radius                                                            |
    |        epsilon        | float                |       |    \*    | the search coefficient (default value is `0.1`)                              |
    |        timeout        | int64                |       |          | Search timeout in nanoseconds (default value is `5s`)                        |
    |    ingress_filters    | Filter.Config        |       |          | Ingress Filter configuration                                                 |
    |    egress_filters     | Filter.Config        |       |          | Egress Filter configuration                                                  |
    |        min_num        | uint32               |       |          | the minimum number of result to be returned                                  |
    | aggregation_algorithm | AggregationAlgorithm |       |          | the search aggregation algorithm option (default value is `ConcurrentQueue`) |

### Output

- the scheme of `payload.v1.Search.Responses`.

  ```rpc
  message Search {
    message Responses {
      repeated Response responses = 1;
    }

    message Response {
      string request_id = 1;
      repeated Object.Distance results = 2;
    }
  }

  message Object {
    message Distance {
      string id = 1;
      float distance = 2;
    }
  }
  ```

  - Search.Responses

    |   field   | type     | label                     | desc.                               |
    | :-------: | :------- | :------------------------ | :---------------------------------- |
    | responses | Response | repeated(Array[Response]) | the list of search results response |

  - Search.Response

    |   field    | type            | label                            | desc.                 |
    | :--------: | :-------------- | :------------------------------- | :-------------------- |
    | request_id | string          |                                  | the unique request ID |
    |  results   | Object.Distance | repeated(Array[Object.Distance]) | search results        |

  - Object.Distance

    |  field   | type   | label | desc.                                                 |
    | :------: | :----- | :---- | :---------------------------------------------------- |
    |    id    | string |       | the vector ID                                         |
    | distance | float  |       | the distance between result vector and request vector |

### Status Code

| code | desc.             |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  13  | INTERNAL          |

### MultiSearchByID RPC

MultiSearchByID RPC is the method to search vectors with multiple IDs in **1** request.

<div class="notice">
gRPC has the message size limitation.<br>
Please be careful that the size of the request exceed the limit.
</div>

### Input

- the scheme of `payload.v1.Search.MultiIDRequest stream`

  ```rpc
  message Search {

    message MultiIDRequest {
        repeated IDRequest requests = 1;
    }

    message IDRequest {
      string id = 1;
      Config config = 2;
    }

    message Config {
      string request_id = 1;
      uint32 num = 2 [ (validate.rules).uint32.gte = 1 ];
      float radius = 3;
      float epsilon = 4;
      int64 timeout = 5;
      Filter.Config ingress_filters = 6;
      Filter.Config egress_filters = 7;
      uint32 min_num = 8;
      AggregationAlgorithm aggregation_algorithm = 9;
    }
  }

  enum AggregationAlgorithm {
    Unknown = 0;
    ConcurrentQueue = 1;
    SortSlice = 2;
    SortPoolSlice = 3;
    PairingHeap = 4;
  }
  ```

  - Search.MultiIDRequest

    |  field   | type      | label                      | required | desc.                       |
    | :------: | :-------- | :------------------------- | :------: | :-------------------------- |
    | requests | IDRequest | repeated(Array[IDRequest]) |    \*    | the searchByID request list |

  - Search.IDRequest

    | field  | type   | label | required | desc.                                   |
    | :----: | :----- | :---- | :------: | :-------------------------------------- |
    |   id   | string |       |    \*    | the vector ID to be searched            |
    | config | Config |       |    \*    | the configuration of the search request |

  - Search.Config

    |         field         | type                 | label | required | desc.                                                                        |
    | :-------------------: | :------------------- | :---- | :------: | :--------------------------------------------------------------------------- |
    |      request_id       | string               |       |          | unique request ID                                                            |
    |          num          | uint32               |       |    \*    | the maximum number of result to be returned                                  |
    |        radius         | float                |       |    \*    | the search radius                                                            |
    |        epsilon        | float                |       |    \*    | the search coefficient (default value is `0.1`)                              |
    |        timeout        | int64                |       |          | Search timeout in nanoseconds (default value is `5s`)                        |
    |    ingress_filters    | Filter.Config        |       |          | Ingress Filter configuration                                                 |
    |    egress_filters     | Filter.Config        |       |          | Egress Filter configuration                                                  |
    |        min_num        | uint32               |       |          | the minimum number of result to be returned                                  |
    | aggregation_algorithm | AggregationAlgorithm |       |          | the search aggregation algorithm option (default value is `ConcurrentQueue`) |

### Output

- the scheme of `payload.v1.Search.Responses`.

  ```rpc
  message Search {
    message Responses {
      repeated Response responses = 1;
    }

    message Response {
      string request_id = 1;
      repeated Object.Distance results = 2;
    }
  }

  message Object {
    message Distance {
      string id = 1;
      float distance = 2;
    }
  }
  ```

  - Search.Responses

    |   field   | type     | label                     | desc.                               |
    | :-------: | :------- | :------------------------ | :---------------------------------- |
    | responses | Response | repeated(Array[Response]) | the list of search results response |

  - Search.Response

    |   field    | type            | label                            | desc.                 |
    | :--------: | :-------------- | :------------------------------- | :-------------------- |
    | request_id | string          |                                  | the unique request ID |
    |  results   | Object.Distance | repeated(Array[Object.Distance]) | search results        |

  - Object.Distance

    |  field   | type   | label | desc.                                                 |
    | :------: | :----- | :---- | :---------------------------------------------------- |
    |    id    | string |       | the vector ID                                         |
    | distance | float  |       | the distance between result vector and request vector |

### Status Code

| code | desc.             |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  13  | INTERNAL          |

## LinearSearch RPC

LinearSearch RPC is the method to linear search vector(s) similar to request vector.

### Input

- the scheme of `payload.v1.Search.Request`

  ```rpc
  message Search {
    message Request {
      repeated float vector = 1 [ (validate.rules).repeated .min_items = 2 ];
      Config config = 2;
    }

    message Config {
      string request_id = 1;
      uint32 num = 2 [ (validate.rules).uint32.gte = 1 ];
      int64 timeout = 5;
      Filter.Config ingress_filters = 6;
      Filter.Config egress_filters = 7;
      uint32 min_num = 8;
      AggregationAlgorithm aggregation_algorithm = 9;
    }
  }

  enum AggregationAlgorithm {
    Unknown = 0;
    ConcurrentQueue = 1;
    SortSlice = 2;
    SortPoolSlice = 3;
    PairingHeap = 4;
  }
  ```

  - Search.Request

    | field  | type   | label                  | required | desc.                                                   |
    | :----: | :----- | :--------------------- | :------: | :------------------------------------------------------ |
    | vector | float  | repeated(Array[float]) |    \*    | the vector data. its dimension is between 2 and 65,536. |
    | config | Config |                        |    \*    | the configuration of the search request                 |

  - Search.Config

    |         field         | type                 | label | required | desc.                                                                        |
    | :-------------------: | :------------------- | :---- | :------: | :--------------------------------------------------------------------------- |
    |      request_id       | string               |       |          | unique request ID                                                            |
    |          num          | uint32               |       |    \*    | the maximum number of result to be returned                                  |
    |        timeout        | int64                |       |          | Search timeout in nanoseconds (default value is `5s`)                        |
    |    ingress_filters    | Filter.Config        |       |          | Ingress Filter configuration                                                 |
    |    egress_filters     | Filter.Config        |       |          | Egress Filter configuration                                                  |
    |        min_num        | uint32               |       |          | the minimum number of result to be returned                                  |
    | aggregation_algorithm | AggregationAlgorithm |       |          | the search aggregation algorithm option (default value is `ConcurrentQueue`) |

### Output

- the scheme of `payload.v1.Search.Response`.

  ```rpc
  message Search {
    message Response {
      string request_id = 1;
      repeated Object.Distance results = 2;
    }
  }

  message Object {
    message Distance {
      string id = 1;
      float distance = 2;
    }
  }
  ```

  - Search.Response

    |   field    | type            | label                            | desc.                 |
    | :--------: | :-------------- | :------------------------------- | :-------------------- |
    | request_id | string          |                                  | the unique request ID |
    |  results   | Object.Distance | repeated(Array[Object.Distance]) | search results        |

  - Object.Distance

    |  field   | type   | label | desc.                                                 |
    | :------: | :----- | :---- | :---------------------------------------------------- |
    |    id    | string |       | the vector ID                                         |
    | distance | float  |       | the distance between result vector and request vector |

### Status Code

| code | desc.             |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  13  | INTERNAL          |

## LinearSearchByID RPC

LinearSearchByID RPC is the method to linear search similar vectors using by user defined vector ID.<br>
The vector with the same requested ID should be indexed into the `vald-agent` before searching.
If the vector doesn't be stored, you will get a `NOT_FOUND` error as a result.

### Input

- the scheme of `payload.v1.Search.IDRequest`

  ```rpc
  message Search {
    message IDRequest {
      string id = 1;
      Config config = 2;
    }

    message Config {
      string request_id = 1;
      uint32 num = 2 [ (validate.rules).uint32.gte = 1 ];
      int64 timeout = 5;
      Filter.Config ingress_filters = 6;
      Filter.Config egress_filters = 7;
      uint32 min_num = 8;
      AggregationAlgorithm aggregation_algorithm = 9;
    }
  }

  enum AggregationAlgorithm {
    Unknown = 0;
    ConcurrentQueue = 1;
    SortSlice = 2;
    SortPoolSlice = 3;
    PairingHeap = 4;
  }
  ```

  - Search.IDRequest

    | field  | type   | label | required | desc.                                   |
    | :----: | :----- | :---- | :------: | :-------------------------------------- |
    |   id   | string |       |    \*    | the vector ID to be searched            |
    | config | Config |       |    \*    | the configuration of the search request |

  - Search.Config

    |         field         | type                 | label | required | desc.                                                                        |
    | :-------------------: | :------------------- | :---- | :------: | :--------------------------------------------------------------------------- |
    |      request_id       | string               |       |          | unique request ID                                                            |
    |          num          | uint32               |       |    \*    | the maximum number of result to be returned                                  |
    |        timeout        | int64                |       |          | Search timeout in nanoseconds (default value is `5s`)                        |
    |    ingress_filters    | Filter.Config        |       |          | Ingress Filter configuration                                                 |
    |    egress_filters     | Filter.Config        |       |          | Egress Filter configuration                                                  |
    |        min_num        | uint32               |       |          | the minimum number of result to be returned                                  |
    | aggregation_algorithm | AggregationAlgorithm |       |          | the search aggregation algorithm option (default value is `ConcurrentQueue`) |

### Output

- the scheme of `payload.v1.Search.Response`.

  ```rpc
  message Search {
    message Response {
      string request_id = 1;
      repeated Object.Distance results = 2;
    }
  }

  message Object {
    message Distance {
      string id = 1;
      float distance = 2;
    }
  }
  ```

  - Search.Response

    |   field    | type            | label                            | desc.                 |
    | :--------: | :-------------- | :------------------------------- | :-------------------- |
    | request_id | string          |                                  | the unique request ID |
    |  results   | Object.Distance | repeated(Array[Object.Distance]) | search results        |

  - Object.Distance

    |  field   | type   | label | desc.                                                 |
    | :------: | :----- | :---- | :---------------------------------------------------- |
    |    id    | string |       | the vector ID                                         |
    | distance | float  |       | the distance between result vector and request vector |

### Status Code

| code | desc.             |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  13  | INTERNAL          |

## StreamLinearSearch RPC

StreamLinearSearch RPC is the method to linear search vectors with multi queries(vectors) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
By using the bidirectional streaming RPC, the linear search request can be communicated in any order between client and server.
Each LinearSearch request and response are independent.

### Input

- the scheme of `payload.v1.Search.Request stream`

  ```rpc
  message Search {
    message Request {
      repeated float vector = 1 [ (validate.rules).repeated .min_items = 2 ];
      Config config = 2;
    }

    message Config {
      string request_id = 1;
      uint32 num = 2 [ (validate.rules).uint32.gte = 1 ];
      int64 timeout = 5;
      Filter.Config ingress_filters = 6;
      Filter.Config egress_filters = 7;
      uint32 min_num = 8;
      AggregationAlgorithm aggregation_algorithm = 9;
    }
  }

  enum AggregationAlgorithm {
    Unknown = 0;
    ConcurrentQueue = 1;
    SortSlice = 2;
    SortPoolSlice = 3;
    PairingHeap = 4;
  }
  ```

  - Search.Request

    | field  | type   | label                  | required | desc.                                                   |
    | :----: | :----- | :--------------------- | :------: | :------------------------------------------------------ |
    | vector | float  | repeated(Array[float]) |    \*    | the vector data. its dimension is between 2 and 65,536. |
    | config | Config |                        |    \*    | the configuration of the search request                 |

  - Search.Config

    |         field         | type                 | label | required | desc.                                                                        |
    | :-------------------: | :------------------- | :---- | :------: | :--------------------------------------------------------------------------- |
    |      request_id       | string               |       |          | unique request ID                                                            |
    |          num          | uint32               |       |    \*    | the maximum number of result to be returned                                  |
    |        timeout        | int64                |       |          | Search timeout in nanoseconds (default value is `5s`)                        |
    |    ingress_filters    | Filter.Config        |       |          | Ingress Filter configuration                                                 |
    |    egress_filters     | Filter.Config        |       |          | Egress Filter configuration                                                  |
    |        min_num        | uint32               |       |          | the minimum number of result to be returned                                  |
    | aggregation_algorithm | AggregationAlgorithm |       |          | the search aggregation algorithm option (default value is `ConcurrentQueue`) |

### Output

- the scheme of `payload.v1.Search.StreamResponse`.

  ```rpc
  message Search {
    message StreamResponse {
      oneof payload {
        Response response = 1;
        google.rpc.Status status = 2;
      }
    }

    message Response {
      string request_id = 1;
      repeated Object.Distance results = 2;
    }
  }

  message Object {
    message Distance {
      string id = 1;
      float distance = 2;
    }
  }
  ```

  - Search.StreamResponse

    |  field   | type              | label | desc.                      |
    | :------: | :---------------- | :---- | :------------------------- |
    | response | Response          |       | the search result response |
    |  status  | google.rpc.Status |       | the status of google RPC   |

  - Search.Response

    |   field    | type            | label                            | desc.                 |
    | :--------: | :-------------- | :------------------------------- | :-------------------- |
    | request_id | string          |                                  | the unique request ID |
    |  results   | Object.Distance | repeated(Array[Object.Distance]) | search results        |

  - Object.Distance

    |  field   | type   | label | desc.                                                 |
    | :------: | :----- | :---- | :---------------------------------------------------- |
    |    id    | string |       | the vector ID                                         |
    | distance | float  |       | the distance between result vector and request vector |

### Status Code

| code | desc.             |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  13  | INTERNAL          |

## StreamLinearSearchByID RPC

StreamLinearSearchByID RPC is the method to linear search vectors with multi queries(IDs) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
By using the bidirectional streaming RPC, the linear search request can be communicated in any order between client and server.
Each LinearSearchByID request and response are independent.

### Input

- the scheme of `payload.v1.Search.IDRequest stream`

  ```rpc
  message Search {
    message IDRequest {
      string id = 1;
      Config config = 2;
    }

    message Config {
      string request_id = 1;
      uint32 num = 2 [ (validate.rules).uint32.gte = 1 ];
      int64 timeout = 5;
      Filter.Config ingress_filters = 6;
      Filter.Config egress_filters = 7;
      uint32 min_num = 8;
      AggregationAlgorithm aggregation_algorithm = 9;
    }
  }

  enum AggregationAlgorithm {
    Unknown = 0;
    ConcurrentQueue = 1;
    SortSlice = 2;
    SortPoolSlice = 3;
    PairingHeap = 4;
  }
  ```

  - Search.IDRequest

    | field  | type   | label | required | desc.                                   |
    | :----: | :----- | :---- | :------: | :-------------------------------------- |
    |   id   | string |       |    \*    | the vector ID to be searched            |
    | config | Config |       |    \*    | the configuration of the search request |

  - Search.Config

    |         field         | type                 | label | required | desc.                                                                        |
    | :-------------------: | :------------------- | :---- | :------: | :--------------------------------------------------------------------------- |
    |      request_id       | string               |       |          | unique request ID                                                            |
    |          num          | uint32               |       |    \*    | the maximum number of result to be returned                                  |
    |        timeout        | int64                |       |          | Search timeout in nanoseconds (default value is `5s`)                        |
    |    ingress_filters    | Filter.Config        |       |          | Ingress Filter configuration                                                 |
    |    egress_filters     | Filter.Config        |       |          | Egress Filter configuration                                                  |
    |        min_num        | uint32               |       |          | the minimum number of result to be returned                                  |
    | aggregation_algorithm | AggregationAlgorithm |       |          | the search aggregation algorithm option (default value is `ConcurrentQueue`) |

### Output

- the scheme of `payload.v1.Search.StreamResponse`.

  ```rpc
  message Search {
    message StreamResponse {
      oneof payload {
        Response response = 1;
        google.rpc.Status status = 2;
      }
    }

    message Response {
      string request_id = 1;
      repeated Object.Distance results = 2;
    }
  }

  message Object {
    message Distance {
      string id = 1;
      float distance = 2;
    }
  }
  ```

  - Search.StreamResponse

    |  field   | type              | label | desc.                      |
    | :------: | :---------------- | :---- | :------------------------- |
    | response | Response          |       | the search result response |
    |  status  | google.rpc.Status |       | the status of google RPC   |

  - Search.Response

    |   field    | type            | label                            | desc.                 |
    | :--------: | :-------------- | :------------------------------- | :-------------------- |
    | request_id | string          |                                  | the unique request ID |
    |  results   | Object.Distance | repeated(Array[Object.Distance]) | search results        |

  - Object.Distance

    |  field   | type   | label | desc.                                                 |
    | :------: | :----- | :---- | :---------------------------------------------------- |
    |    id    | string |       | the vector ID                                         |
    | distance | float  |       | the distance between result vector and request vector |

### Status Code

| code | desc.             |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  13  | INTERNAL          |

## MultiLinearSearch RPC

MultiLinearSearch RPC is the method to linear search vectors with multiple vectors in **1** request.

<div class="notice">
gRPC has the message size limitation.<br>
Please be careful that the size of the request exceed the limit.
</div>

### Input

- the scheme of `payload.v1.Search.MultiRequest`

  ```rpc
  message Search {
    message MultiRequest {
      repeated Request requests = 1;
    }

    message Request {
      repeated float vector = 1 [ (validate.rules).repeated .min_items = 2 ];
      Config config = 2;
    }

    message Config {
      string request_id = 1;
      uint32 num = 2 [ (validate.rules).uint32.gte = 1 ];
      int64 timeout = 5;
      Filter.Config ingress_filters = 6;
      Filter.Config egress_filters = 7;
      uint32 min_num = 8;
      AggregationAlgorithm aggregation_algorithm = 9;
    }
  }

  enum AggregationAlgorithm {
    Unknown = 0;
    ConcurrentQueue = 1;
    SortSlice = 2;
    SortPoolSlice = 3;
    PairingHeap = 4;
  }
  ```

  - Search.MultiRequest

    |  field   | type                     | label | required | desc.                   |
    | :------: | :----------------------- | :---- | :------: | :---------------------- |
    | requests | repeated(Array[Request]) |       |    \*    | the search request list |

  - Search.Request

    | field  | type   | label                  | required | desc.                                                   |
    | :----: | :----- | :--------------------- | :------: | :------------------------------------------------------ |
    | vector | float  | repeated(Array[float]) |    \*    | the vector data. its dimension is between 2 and 65,536. |
    | config | Config |                        |    \*    | the configuration of the search request                 |

  - Search.Config

    |         field         | type                 | label | required | desc.                                                                        |
    | :-------------------: | :------------------- | :---- | :------: | :--------------------------------------------------------------------------- |
    |      request_id       | string               |       |          | unique request ID                                                            |
    |          num          | uint32               |       |    \*    | the maximum number of result to be returned                                  |
    |        timeout        | int64                |       |          | Search timeout in nanoseconds (default value is `5s`)                        |
    |    ingress_filters    | Filter.Config        |       |          | Ingress Filter configuration                                                 |
    |    egress_filters     | Filter.Config        |       |          | Egress Filter configuration                                                  |
    |        min_num        | uint32               |       |          | the minimum number of result to be returned                                  |
    | aggregation_algorithm | AggregationAlgorithm |       |          | the search aggregation algorithm option (default value is `ConcurrentQueue`) |

### Output

- the scheme of `payload.v1.Search.Responses`.

  ```rpc
  message Search {
    message Responses {
      repeated Response responses = 1;
    }

    message Response {
      string request_id = 1;
      repeated Object.Distance results = 2;
    }
  }

  message Object {
    message Distance {
      string id = 1;
      float distance = 2;
    }
  }
  ```

  - Search.Responses

    |   field   | type     | label                     | desc.                               |
    | :-------: | :------- | :------------------------ | :---------------------------------- |
    | responses | Response | repeated(Array[Response]) | the list of search results response |

  - Search.Response

    |   field    | type            | label                            | desc.                 |
    | :--------: | :-------------- | :------------------------------- | :-------------------- |
    | request_id | string          |                                  | the unique request ID |
    |  results   | Object.Distance | repeated(Array[Object.Distance]) | search results        |

  - Object.Distance

    |  field   | type   | label | desc.                                                 |
    | :------: | :----- | :---- | :---------------------------------------------------- |
    |    id    | string |       | the vector ID                                         |
    | distance | float  |       | the distance between result vector and request vector |

### Status Code

| code | desc.             |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  13  | INTERNAL          |

### MultiLinearSearchByID RPC

MultiLinearSearchByID RPC is the method to linear search vectors with multiple IDs in **1** request.

<div class="notice">
gRPC has the message size limitation.<br>
Please be careful that the size of the request exceed the limit.
</div>

### Input

- the scheme of `payload.v1.Search.MultiIDRequest stream`

  ```rpc
  message Search {

    message MultiIDRequest {
        repeated IDRequest requests = 1;
    }

    message IDRequest {
      string id = 1;
      Config config = 2;
    }

    message Config {
      string request_id = 1;
      uint32 num = 2 [ (validate.rules).uint32.gte = 1 ];
      int64 timeout = 5;
      Filter.Config ingress_filters = 6;
      Filter.Config egress_filters = 7;
      uint32 min_num = 8;
      AggregationAlgorithm aggregation_algorithm = 9;
    }
  }

  enum AggregationAlgorithm {
    Unknown = 0;
    ConcurrentQueue = 1;
    SortSlice = 2;
    SortPoolSlice = 3;
    PairingHeap = 4;
  }
  ```

  - Search.MultiIDRequest

    |  field   | type      | label                      | required | desc.                       |
    | :------: | :-------- | :------------------------- | :------: | :-------------------------- |
    | requests | IDRequest | repeated(Array[IDRequest]) |    \*    | the searchByID request list |

  - Search.IDRequest

    | field  | type   | label | required | desc.                                   |
    | :----: | :----- | :---- | :------: | :-------------------------------------- |
    |   id   | string |       |    \*    | the vector ID to be searched.           |
    | config | Config |       |    \*    | the configuration of the search request |

  - Search.Config

    |         field         | type                 | label | required | desc.                                                                        |
    | :-------------------: | :------------------- | :---- | :------: | :--------------------------------------------------------------------------- |
    |      request_id       | string               |       |          | unique request ID                                                            |
    |          num          | uint32               |       |    \*    | the maximum number of result to be returned                                  |
    |        timeout        | int64                |       |          | Search timeout in nanoseconds (default value is `5s`)                        |
    |    ingress_filters    | Filter.Config        |       |          | Ingress Filter configuration                                                 |
    |    egress_filters     | Filter.Config        |       |          | Egress Filter configuration                                                  |
    |        min_num        | uint32               |       |          | the minimum number of result to be returned                                  |
    | aggregation_algorithm | AggregationAlgorithm |       |          | the search aggregation algorithm option (default value is `ConcurrentQueue`) |

### Output

- the scheme of `payload.v1.Search.Responses`.

  ```rpc
  message Search {
    message Responses {
      repeated Response responses = 1;
    }

    message Response {
      string request_id = 1;
      repeated Object.Distance results = 2;
    }
  }

  message Object {
    message Distance {
      string id = 1;
      float distance = 2;
    }
  }
  ```

  - Search.Responses

    |   field   | type     | label                     | desc.                               |
    | :-------: | :------- | :------------------------ | :---------------------------------- |
    | responses | Response | repeated(Array[Response]) | the list of search results response |

  - Search.Response

    |   field    | type            | label                            | desc.                 |
    | :--------: | :-------------- | :------------------------------- | :-------------------- |
    | request_id | string          |                                  | the unique request ID |
    |  results   | Object.Distance | repeated(Array[Object.Distance]) | search results        |

  - Object.Distance

    |  field   | type   | label | desc.                                                 |
    | :------: | :----- | :---- | :---------------------------------------------------- |
    |    id    | string |       | the vector ID                                         |
    | distance | float  |       | the distance between result vector and request vector |

### Status Code

| code | desc.             |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  13  | INTERNAL          |
