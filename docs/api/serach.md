# Vald Search APIs

## Overview

Search Service is responsible for searching vectors from `vald-agent` which are similar to user request vector.

```bash
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

  ```bash
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
    }
  }
  ```

  - Search.Request
    |field|type|label|required|desc.|
    |:---:|:---|:---|:---:|:---|
    |vector|float| repeated(Array[float]) | \* | the vector data. its dimension is between 2 and 65,536.|
    |config|Config| | \* | the configuration of the search request |

  - Search.Config
    |field|type|label|required|desc.|
    |:---:|:---|:---|:---:|:---|
    |request_id|string| | | unique request ID |
    |num|uint32| | \* | the maximum number of result to be returned |
    |radius|float| | \* | the search radius |
    |epsilon|float| | \* | the search coefficient (default value is `0.1`) |
    |timeout|int64| | | Search timeout in nanoseconds (default value is `5s`) |
    |ingress_filters|Filter.Config| | | Ingress Filter configuration |
    |egress_filters|Filter.Config| | | Egress Filter configuration |
    |min_num| uint32 | | the minimum number of result to be returned |

### Output

- the scheme of `payload.v1.Search.Response`.

  ```bash
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
    |field|type|label|desc.|
    |:---:|:---|:---|:---|
    |request_id|string| | the unique request ID |
    |results|Object.Distance| repeated(Array[Object.Distance]) | search results |

  - Object.Distance
    |field|type|label|desc.|
    |:---:|:---|:---|:---|
    |id|string| | the vector ID |
    |distance|float| | the distance between result vector and request vector |

### Status Code

| code | desc.             |
| :--: | :---------------- |
|  0   | OK                |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |

## SearchByID RPC

SearchByID RPC it the method to search similar vectors using by user defined vector ID.<br>
The vector with the same requested ID should be indexed into the `vald-agent` before searching.

### Input

- the scheme of `payload.v1.Search.IDRequest`

  ```bash
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
    }
  }
  ```

  - Search.IDRequest
    |field|type|label|required|desc.|
    |:---:|:---|:---|:---:|:---|
    |id|string| | \* | the vector ID to be searched |
    |config|Config| | \* | the configuration of the search request |

  - Search.Config
    |field|type|label|required|desc.|
    |:---:|:---|:---|:---:|:---|
    |request_id|string| | | unique request ID |
    |num|uint32| | \* | the maximum number of result to be returned |
    |radius|float| | \* | the search radius |
    |epsilon|float| | \* | the search coefficient (default value is `0.1`) |
    |timeout|int64| | | Search timeout in nanoseconds (default value is `5s`) |
    |ingress_filters|Filter.Config| | | Ingress Filter configuration |
    |egress_filters|Filter.Config| | | Egress Filter configuration |
    |min_num| uint32 | | the minimum number of result to be returned |

### Output

- the scheme of `payload.v1.Search.Response`.

  ```bash
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
    |field|type|label|desc.|
    |:---:|:---|:---|:---|
    |request_id|string| | the unique request ID |
    |results|Object.Distance| repeated(Array[Object.Distance]) | search results |

  - Object.Distance
    |field|type|label|desc.|
    |:---:|:---|:---|:---|
    |id|string| | the vector ID |
    |distance|float| | the distance between result vector and request vector |

### Status Code

| code | desc.             |
| :--: | :---------------- |
|  0   | OK                |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |

## StreamSearch RPC

StreamSearch RPC is the method to search vectors with multi queries(vectors) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
By using the bidirectional streaming RPC, the search request can be communicated in any order between client and server.
Each Search request and response are independent.

### Input

- the scheme of `payload.v1.Search.Request stream`

  ```bash
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
      }
  }
  ```

  - Search.Request
    |field|type|label|required|desc.|
    |:---:|:---|:---|:---:|:---|
    |vector|float| repeated(Array[float]) | \* | the vector data. its dimension is between 2 and 65,536.|
    |config|Config| | \* | the configuration of the search request |

  - Search.Config
    |field|type|label|required|desc.|
    |:---:|:---|:---|:---:|:---|
    |request_id|string| | | unique request ID |
    |num|uint32| | \* | the maximum number of result to be returned |
    |radius|float| | \* | the search radius |
    |epsilon|float| | \* | the search coefficient (default value is `0.1`) |
    |timeout|int64| | | Search timeout in nanoseconds (default value is `5s`) |
    |ingress_filters|Filter.Config| | | Ingress Filter configuration |
    |egress_filters|Filter.Config| | | Egress Filter configuration |
    |min_num| uint32 | | the minimum number of result to be returned |

### Output

- the scheme of `payload.v1.Search.StreamResponse`.

  ```bash
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
    |field|type|label|desc.|
    |:---:|:---|:---|:---|
    |response|Response| | the search result response |
    |status|google.rpc.Status| | the status of google RPC |

  - Search.Response
    |field|type|label|desc.|
    |:---:|:---|:---|:---|
    |request_id|string| | the unique request ID |
    |results|Object.Distance| repeated(Array[Object.Distance]) | search results |

  - Object.Distance
    |field|type|label|desc.|
    |:---:|:---|:---|:---|
    |id|string| | the vector ID |
    |distance|float| | the distance between result vector and request vector |

### Status Code

| code | desc.             |
| :--: | :---------------- |
|  0   | OK                |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |

## StreamSearchByID RPC

StreamSearchByID RPC is the method to search vectors with multi queries(IDs) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
By using the bidirectional streaming RPC, the search request can be communicated in any order between client and server.
Each SearchByID request and response are independent.

### Input

- the scheme of `payload.v1.Search.IDRequest stream`

  ```bash
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
    }
  }
  ```

  - Search.IDRequest
    |field|type|label|required|desc.|
    |:---:|:---|:---|:---:|:---|
    |id|string| | \* | the vector ID to be searched |
    |config|Config| | \* | the configuration of the search request |

  - Search.Config
    |field|type|label|required|desc.|
    |:---:|:---|:---|:---:|:---|
    |request_id|string| | | unique request ID |
    |num|uint32| | \* | the maximum number of result to be returned |
    |radius|float| | \* | the search radius |
    |epsilon|float| | \* | the search coefficient (default value is `0.1`) |
    |timeout|int64| | | Search timeout in nanoseconds (default value is `5s`) |
    |ingress_filters|Filter.Config| | | Ingress Filter configuration |
    |egress_filters|Filter.Config| | | Egress Filter configuration |
    |min_num| uint32 | | the minimum number of result to be returned |

### Output

- the scheme of `payload.v1.Search.StreamResponse`.

  ```bash
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
    |field|type|label|desc.|
    |:---:|:---|:---|:---|
    |response|Response| | the search result response |
    |status|google.rpc.Status| | the status of google RPC |

  - Search.Response
    |field|type|label|desc.|
    |:---:|:---|:---|:---|
    |request_id|string| | the unique request ID |
    |results|Object.Distance| repeated(Array[Object.Distance]) | search results |

  - Object.Distance
    |field|type|label|desc.|
    |:---:|:---|:---|:---|
    |id|string| | the vector ID |
    |distance|float| | the distance between result vector and request vector |

### Status Code

| code | desc.             |
| :--: | :---------------- |
|  0   | OK                |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |

## MultiSearch RPC

MultiSearch RPC is the method to search vectors with multiple vectors in **1** request.

<div class="card-note">
gRPC has the message size limitation.<br>
Please be careful that the size of the request exceed the limit.
</div>

### Input

- the scheme of `payload.v1.Search.MultiRequest`

  ```bash
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
    }
  }
  ```

  - Search.MultiRequest
    |field|type|label|required|desc.|
    |:---:|:---|:---|:---:|:---|
    |requests| repeated(Array[Request]) | | \* | the search request list |

  - Search.Request
    |field|type|label|required|desc.|
    |:---:|:---|:---|:---:|:---|
    |vector|float| repeated(Array[float]) | \* | the vector data. its dimension is between 2 and 65,536.|
    |config|Config| | \* | the configuration of the search request |

  - Search.Config
    |field|type|label|required|desc.|
    |:---:|:---|:---|:---:|:---|
    |request_id|string| | | unique request ID |
    |num|uint32| | \* | the maximum number of result to be returned |
    |radius|float| | \* | the search radius |
    |epsilon|float| | \* | the search coefficient (default value is `0.1`) |
    |timeout|int64| | | Search timeout in nanoseconds (default value is `5s`) |
    |ingress_filters|Filter.Config| | | Ingress Filter configuration |
    |egress_filters|Filter.Config| | | Egress Filter configuration |
    |min_num| uint32 | | the minimum number of result to be returned |

### Output

- the scheme of `payload.v1.Search.Responses`.

  ```bash
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
    |field|type|label|desc.|
    |:---:|:---|:---|:---|
    |responses| Response | repeated(Array[Response]) | the list of search results response |

  - Search.Response
    |field|type|label|desc.|
    |:---:|:---|:---|:---|
    |request_id|string| | the unique request ID |
    |results|Object.Distance| repeated(Array[Object.Distance]) | search results |

  - Object.Distance
    |field|type|label|desc.|
    |:---:|:---|:---|:---|
    |id|string| | the vector ID |
    |distance|float| | the distance between result vector and request vector |

### Status Code

| code | desc.             |
| :--: | :---------------- |
|  0   | OK                |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |

### MultiSearchByID RPC

MultiSearchByID RPC is the method to search vectors with multiple IDs in **1** request.

<div class="card-note">
gRPC has the message size limitation.<br>
Please be careful that the size of the request exceed the limit.
</div>

### Input

- the scheme of `payload.v1.Search.MultiIDRequest stream`

  ```bash
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
    }
  }
  ```

  - Search.MultiIDRequest
    |field|type|label|required|desc.|
    |:---:|:---|:---|:---:|:---|
    |requests|IDRequest| repeated(Array[IDRequest]) | \* | the searchByID request list |

  - Search.IDRequest
    |field|type|label|required|desc.|
    |:---:|:---|:---|:---:|:---|
    |id|string| | \* | the vector ID to be searched |
    |config|Config| | \* | the configuration of the search request |

  - Search.Config
    |field|type|label|required|desc.|
    |:---:|:---|:---|:---:|:---|
    |request_id|string| | | unique request ID |
    |num|uint32| | \* | the maximum number of result to be returned |
    |radius|float| | \* | the search radius |
    |epsilon|float| | \* | the search coefficient (default value is `0.1`) |
    |timeout|int64| | | Search timeout in nanoseconds (default value is `5s`) |
    |ingress_filters|Filter.Config| | | Ingress Filter configuration |
    |egress_filters|Filter.Config| | | Egress Filter configuration |
    |min_num| uint32 | | the minimum number of result to be returned |

### Output

- the scheme of `payload.v1.Search.Responses`.

  ```bash
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
    |field|type|label|desc.|
    |:---:|:---|:---|:---|
    |responses| Response | repeated(Array[Response]) | the list of search results response |

  - Search.Response
    |field|type|label|desc.|
    |:---:|:---|:---|:---|
    |request_id|string| | the unique request ID |
    |results|Object.Distance| repeated(Array[Object.Distance]) | search results |

  - Object.Distance
    |field|type|label|desc.|
    |:---:|:---|:---|:---|
    |id|string| | the vector ID |
    |distance|float| | the distance between result vector and request vector |

### Status Code

| code | desc.             |
| :--: | :---------------- |
|  0   | OK                |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |

## LinearSearch RPC

LinearSearch RPC is the method to linear search vector(s) similar to request vector.

### Input

- the scheme of `payload.v1.Search.Request`

  ```bash
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
    }
  }
  ```

  - Search.Request
    |field|type|label|required|desc.|
    |:---:|:---|:---|:---:|:---|
    |vector|float| repeated(Array[float]) | \* | the vector data. its dimension is between 2 and 65,536.|
    |config|Config| | \* | the configuration of the search request |

  - Search.Config
    |field|type|label|required|desc.|
    |:---:|:---|:---|:---:|:---|
    |request_id|string| | | unique request ID |
    |num|uint32| | \* | the maximum number of result to be returned |
    |timeout|int64| | | Search timeout in nanoseconds (default value is `5s`) |
    |ingress_filters|Filter.Config| | | Ingress Filter configuration |
    |egress_filters|Filter.Config| | | Egress Filter configuration |
    |min_num| uint32 | | the minimum number of result to be returned |

### Output

- the scheme of `payload.v1.Search.Response`.

  ```bash
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
    |field|type|label|desc.|
    |:---:|:---|:---|:---|
    |request_id|string| | the unique request ID |
    |results|Object.Distance| repeated(Array[Object.Distance]) | search results |

  - Object.Distance
    |field|type|label|desc.|
    |:---:|:---|:---|:---|
    |id|string| | the vector ID |
    |distance|float| | the distance between result vector and request vector |

### Status Code

| code | desc.             |
| :--: | :---------------- |
|  0   | OK                |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |

## LinearSearchByID RPC

LinearSearchByID RPC is the method to linear search similar vectors using by user defined vector ID.<br>
The vector with the same requested ID should be indexed into the `vald-agent` before searching.
If the vector doesn't be stored, you will get a `NOT_FOUND` error as a result.

### Input

- the scheme of `payload.v1.Search.IDRequest`

  ```bash
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
    }
  }
  ```

  - Search.IDRequest
    |field|type|label|required|desc.|
    |:---:|:---|:---|:---:|:---|
    |id|string| | \* | the vector ID to be searched |
    |config|Config| | \* | the configuration of the search request |

  - Search.Config
    |field|type|label|required|desc.|
    |:---:|:---|:---|:---:|:---|
    |request_id|string| | | unique request ID |
    |num|uint32| | \* | the maximum number of result to be returned |
    |timeout|int64| | | Search timeout in nanoseconds (default value is `5s`) |
    |ingress_filters|Filter.Config| | | Ingress Filter configuration |
    |egress_filters|Filter.Config| | | Egress Filter configuration |
    |min_num| uint32 | | the minimum number of result to be returned |

### Output

- the scheme of `payload.v1.Search.Response`.

  ```bash
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
    |field|type|label|desc.|
    |:---:|:---|:---|:---|
    |request_id|string| | the unique request ID |
    |results|Object.Distance| repeated(Array[Object.Distance]) | search results |

  - Object.Distance
    |field|type|label|desc.|
    |:---:|:---|:---|:---|
    |id|string| | the vector ID |
    |distance|float| | the distance between result vector and request vector |

### Status Code

| code | desc.             |
| :--: | :---------------- |
|  0   | OK                |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |

## StreamLinearSearch RPC

StreamLinearSearch RPC is the method to linear search vectors with multi queries(vectors) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
By using the bidirectional streaming RPC, the linear search request can be communicated in any order between client and server.
Each LinearSearch request and response are independent.

### Input

- the scheme of `payload.v1.Search.Request stream`

  ```bash
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
      }
  }
  ```

  - Search.Request
    |field|type|label|required|desc.|
    |:---:|:---|:---|:---:|:---|
    |vector|float| repeated(Array[float]) | \* | the vector data. its dimension is between 2 and 65,536.|
    |config|Config| | \* | the configuration of the search request |

  - Search.Config
    |field|type|label|required|desc.|
    |:---:|:---|:---|:---:|:---|
    |request_id|string| | | unique request ID |
    |num|uint32| | \* | the maximum number of result to be returned |
    |timeout|int64| | | Search timeout in nanoseconds (default value is `5s`) |
    |ingress_filters|Filter.Config| | | Ingress Filter configuration |
    |egress_filters|Filter.Config| | | Egress Filter configuration |
    |min_num| uint32 | | the minimum number of result to be returned |

### Output

- the scheme of `payload.v1.Search.StreamResponse`.

  ```bash
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
    |field|type|label|desc.|
    |:---:|:---|:---|:---|
    |response|Response| | the search result response |
    |status|google.rpc.Status| | the status of google RPC |

  - Search.Response
    |field|type|label|desc.|
    |:---:|:---|:---|:---|
    |request_id|string| | the unique request ID |
    |results|Object.Distance| repeated(Array[Object.Distance]) | search results |

  - Object.Distance
    |field|type|label|desc.|
    |:---:|:---|:---|:---|
    |id|string| | the vector ID |
    |distance|float| | the distance between result vector and request vector |

### Status Code

| code | desc.             |
| :--: | :---------------- |
|  0   | OK                |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |

## StreamLinearSearchByID RPC

StreamLinearSearchByID RPC is the method to linear search vectors with multi queries(IDs) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
By using the bidirectional streaming RPC, the linear search request can be communicated in any order between client and server.
Each LinearSearchByID request and response are independent.

### Input

- the scheme of `payload.v1.Search.IDRequest stream`

  ```bash
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
    }
  }
  ```

  - Search.IDRequest
    |field|type|label|required|desc.|
    |:---:|:---|:---|:---:|:---|
    |id|string| | \* | the vector ID to be searched |
    |config|Config| | \* | the configuration of the search request |

  - Search.Config
    |field|type|label|required|desc.|
    |:---:|:---|:---|:---:|:---|
    |request_id|string| | | unique request ID |
    |num|uint32| | \* | the maximum number of result to be returned |
    |timeout|int64| | | Search timeout in nanoseconds (default value is `5s`) |
    |ingress_filters|Filter.Config| | | Ingress Filter configuration |
    |egress_filters|Filter.Config| | | Egress Filter configuration |
    |min_num| uint32 | | the minimum number of result to be returned |

### Output

- the scheme of `payload.v1.Search.StreamResponse`.

  ```bash
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
    |field|type|label|desc.|
    |:---:|:---|:---|:---|
    |response|Response| | the search result response |
    |status|google.rpc.Status| | the status of google RPC |

  - Search.Response
    |field|type|label|desc.|
    |:---:|:---|:---|:---|
    |request_id|string| | the unique request ID |
    |results|Object.Distance| repeated(Array[Object.Distance]) | search results |

  - Object.Distance
    |field|type|label|desc.|
    |:---:|:---|:---|:---|
    |id|string| | the vector ID |
    |distance|float| | the distance between result vector and request vector |

### Status Code

| code | desc.             |
| :--: | :---------------- |
|  0   | OK                |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |

## MultiLinearSearch RPC

MultiLinearSearch RPC is the method to linear search vectors with multiple vectors in **1** request.

<div class="card-note">
gRPC has the message size limitation.<br>
Please be careful that the size of the request exceed the limit.
</div>

### Input

- the scheme of `payload.v1.Search.MultiRequest`

  ```bash
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
    }
  }
  ```

  - Search.MultiRequest
    |field|type|label|required|desc.|
    |:---:|:---|:---|:---:|:---|
    |requests| repeated(Array[Request]) | | \* | the search request list |

  - Search.Request
    |field|type|label|required|desc.|
    |:---:|:---|:---|:---:|:---|
    |vector|float| repeated(Array[float]) | \* | the vector data. its dimension is between 2 and 65,536.|
    |config|Config| | \* | the configuration of the search request |

  - Search.Config
    |field|type|label|required|desc.|
    |:---:|:---|:---|:---:|:---|
    |request_id|string| | | unique request ID |
    |num|uint32| | \* | the maximum number of result to be returned |
    |timeout|int64| | | Search timeout in nanoseconds (default value is `5s`) |
    |ingress_filters|Filter.Config| | | Ingress Filter configuration |
    |egress_filters|Filter.Config| | | Egress Filter configuration |
    |min_num| uint32 | | the minimum number of result to be returned |

### Output

- the scheme of `payload.v1.Search.Responses`.

  ```bash
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
    |field|type|label|desc.|
    |:---:|:---|:---|:---|
    |responses| Response | repeated(Array[Response]) | the list of search results response |

  - Search.Response
    |field|type|label|desc.|
    |:---:|:---|:---|:---|
    |request_id|string| | the unique request ID |
    |results|Object.Distance| repeated(Array[Object.Distance]) | search results |

  - Object.Distance
    |field|type|label|desc.|
    |:---:|:---|:---|:---|
    |id|string| | the vector ID |
    |distance|float| | the distance between result vector and request vector |

### Status Code

| code | desc.             |
| :--: | :---------------- |
|  0   | OK                |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |

### MultiLinearSearchByID RPC

MultiLinearSearchByID RPC is the method to linear search vectors with multiple IDs in **1** request.

<div class="card-note">
gRPC has the message size limitation.<br>
Please be careful that the size of the request exceed the limit.
</div>

### Input

- the scheme of `payload.v1.Search.MultiIDRequest stream`

  ```bash
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
    }
  }
  ```

  - Search.MultiIDRequest
    |field|type|label|required|desc.|
    |:---:|:---|:---|:---:|:---|
    |requests|IDRequest| repeated(Array[IDRequest]) | \* | the searchByID request list |

  - Search.IDRequest
    |field|type|label|required|desc.|
    |:---:|:---|:---|:---:|:---|
    |id|string| | \* | the vector ID to be searched.|
    |config|Config| | \* | the configuration of the search request |

  - Search.Config
    |field|type|label|required|desc.|
    |:---:|:---|:---|:---:|:---|
    |request_id|string| | | unique request ID |
    |num|uint32| | \* | the maximum number of result to be returned |
    |timeout|int64| | | Search timeout in nanoseconds (default value is `5s`) |
    |ingress_filters|Filter.Config| | | Ingress Filter configuration |
    |egress_filters|Filter.Config| | | Egress Filter configuration |
    |min_num| uint32 | | the minimum number of result to be returned |

### Output

- the scheme of `payload.v1.Search.Responses`.

  ```bash
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
    |field|type|label|desc.|
    |:---:|:---|:---|:---|
    |responses| Response | repeated(Array[Response]) | the list of search results response |

  - Search.Response
    |field|type|label|desc.|
    |:---:|:---|:---|:---|
    |request_id|string| | the unique request ID |
    |results|Object.Distance| repeated(Array[Object.Distance]) | search results |

  - Object.Distance
    |field|type|label|desc.|
    |:---:|:---|:---|:---|
    |id|string| | the vector ID |
    |distance|float| | the distance between result vector and request vector |

### Status Code

| code | desc.             |
| :--: | :---------------- |
|  0   | OK                |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
