# Vald Search APIs

## Overview

Search service provides ways to search indexed vectors.

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

A method to search indexed vectors by a raw vector.

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
  ```

  - Search.Request

    | field  | type          | label    | desc.                                    |
    | :----: | :------------ | :------- | :--------------------------------------- |
    | vector | float         | repeated | The vector to be searched.               |
    | config | Search.Config |          | The configuration of the search request. |

  - Search.Config

    |         field         | type                        | label | desc.                                        |
    | :-------------------: | :-------------------------- | :---- | :------------------------------------------- |
    |      request_id       | string                      |       | Unique request ID.                           |
    |          num          | uint32                      |       | Maximum number of result to be returned.     |
    |        radius         | float                       |       | Search radius.                               |
    |        epsilon        | float                       |       | Search coefficient.                          |
    |        timeout        | int64                       |       | Search timeout in nanoseconds.               |
    |    ingress_filters    | Filter.Config               |       | Ingress filter configurations.               |
    |    egress_filters     | Filter.Config               |       | Egress filter configurations.                |
    |        min_num        | uint32                      |       | Minimum number of result to be returned.     |
    | aggregation_algorithm | Search.AggregationAlgorithm |       | Aggregation Algorithm                        |
    |         ratio         | google.protobuf.FloatValue  |       | Search ratio for agent return result number. |
    |        nprobe         | uint32                      |       | Search nprobe.                               |

  - Filter.Config

    |  field  | type          | label    | desc.                                      |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | desc.                |
    | :---: | :----- | :---- | :------------------- |
    | host  | string |       | The target hostname. |
    | port  | uint32 |       | The target port.     |

  - Filter.Config

    |  field  | type          | label    | desc.                                      |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | desc.                |
    | :---: | :----- | :---- | :------------------- |
    | host  | string |       | The target hostname. |
    | port  | uint32 |       | The target port.     |

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

    |   field    | type            | label    | desc.                  |
    | :--------: | :-------------- | :------- | :--------------------- |
    | request_id | string          |          | The unique request ID. |
    |  results   | Object.Distance | repeated | Search results.        |

  - Object.Distance

    |  field   | type   | label | desc.          |
    | :------: | :----- | :---- | :------------- |
    |    id    | string |       | The vector ID. |
    | distance | float  |       | The distance.  |

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

A method to search indexed vectors by ID.

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
  ```

  - Search.IDRequest

    | field  | type          | label | desc.                                    |
    | :----: | :------------ | :---- | :--------------------------------------- |
    |   id   | string        |       | The vector ID to be searched.            |
    | config | Search.Config |       | The configuration of the search request. |

  - Search.Config

    |         field         | type                        | label | desc.                                        |
    | :-------------------: | :-------------------------- | :---- | :------------------------------------------- |
    |      request_id       | string                      |       | Unique request ID.                           |
    |          num          | uint32                      |       | Maximum number of result to be returned.     |
    |        radius         | float                       |       | Search radius.                               |
    |        epsilon        | float                       |       | Search coefficient.                          |
    |        timeout        | int64                       |       | Search timeout in nanoseconds.               |
    |    ingress_filters    | Filter.Config               |       | Ingress filter configurations.               |
    |    egress_filters     | Filter.Config               |       | Egress filter configurations.                |
    |        min_num        | uint32                      |       | Minimum number of result to be returned.     |
    | aggregation_algorithm | Search.AggregationAlgorithm |       | Aggregation Algorithm                        |
    |         ratio         | google.protobuf.FloatValue  |       | Search ratio for agent return result number. |
    |        nprobe         | uint32                      |       | Search nprobe.                               |

  - Filter.Config

    |  field  | type          | label    | desc.                                      |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | desc.                |
    | :---: | :----- | :---- | :------------------- |
    | host  | string |       | The target hostname. |
    | port  | uint32 |       | The target port.     |

  - Filter.Config

    |  field  | type          | label    | desc.                                      |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | desc.                |
    | :---: | :----- | :---- | :------------------- |
    | host  | string |       | The target hostname. |
    | port  | uint32 |       | The target port.     |

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

    |   field    | type            | label    | desc.                  |
    | :--------: | :-------------- | :------- | :--------------------- |
    | request_id | string          |          | The unique request ID. |
    |  results   | Object.Distance | repeated | Search results.        |

  - Object.Distance

    |  field   | type   | label | desc.          |
    | :------: | :----- | :---- | :------------- |
    |    id    | string |       | The vector ID. |
    | distance | float  |       | The distance.  |

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

A method to search indexed vectors by multiple vectors.

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
  ```

  - Search.Request

    | field  | type          | label    | desc.                                    |
    | :----: | :------------ | :------- | :--------------------------------------- |
    | vector | float         | repeated | The vector to be searched.               |
    | config | Search.Config |          | The configuration of the search request. |

  - Search.Config

    |         field         | type                        | label | desc.                                        |
    | :-------------------: | :-------------------------- | :---- | :------------------------------------------- |
    |      request_id       | string                      |       | Unique request ID.                           |
    |          num          | uint32                      |       | Maximum number of result to be returned.     |
    |        radius         | float                       |       | Search radius.                               |
    |        epsilon        | float                       |       | Search coefficient.                          |
    |        timeout        | int64                       |       | Search timeout in nanoseconds.               |
    |    ingress_filters    | Filter.Config               |       | Ingress filter configurations.               |
    |    egress_filters     | Filter.Config               |       | Egress filter configurations.                |
    |        min_num        | uint32                      |       | Minimum number of result to be returned.     |
    | aggregation_algorithm | Search.AggregationAlgorithm |       | Aggregation Algorithm                        |
    |         ratio         | google.protobuf.FloatValue  |       | Search ratio for agent return result number. |
    |        nprobe         | uint32                      |       | Search nprobe.                               |

  - Filter.Config

    |  field  | type          | label    | desc.                                      |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | desc.                |
    | :---: | :----- | :---- | :------------------- |
    | host  | string |       | The target hostname. |
    | port  | uint32 |       | The target port.     |

  - Filter.Config

    |  field  | type          | label    | desc.                                      |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | desc.                |
    | :---: | :----- | :---- | :------------------- |
    | host  | string |       | The target hostname. |
    | port  | uint32 |       | The target port.     |

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

    |  field   | type              | label | desc.                          |
    | :------: | :---------------- | :---- | :----------------------------- |
    | response | Search.Response   |       | Represent the search response. |
    |  status  | google.rpc.Status |       | The RPC error status.          |

  - Search.Response

    |   field    | type            | label    | desc.                  |
    | :--------: | :-------------- | :------- | :--------------------- |
    | request_id | string          |          | The unique request ID. |
    |  results   | Object.Distance | repeated | Search results.        |

  - Object.Distance

    |  field   | type   | label | desc.          |
    | :------: | :----- | :---- | :------------- |
    |    id    | string |       | The vector ID. |
    | distance | float  |       | The distance.  |

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

A method to search indexed vectors by multiple IDs.

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
  ```

  - Search.IDRequest

    | field  | type          | label | desc.                                    |
    | :----: | :------------ | :---- | :--------------------------------------- |
    |   id   | string        |       | The vector ID to be searched.            |
    | config | Search.Config |       | The configuration of the search request. |

  - Search.Config

    |         field         | type                        | label | desc.                                        |
    | :-------------------: | :-------------------------- | :---- | :------------------------------------------- |
    |      request_id       | string                      |       | Unique request ID.                           |
    |          num          | uint32                      |       | Maximum number of result to be returned.     |
    |        radius         | float                       |       | Search radius.                               |
    |        epsilon        | float                       |       | Search coefficient.                          |
    |        timeout        | int64                       |       | Search timeout in nanoseconds.               |
    |    ingress_filters    | Filter.Config               |       | Ingress filter configurations.               |
    |    egress_filters     | Filter.Config               |       | Egress filter configurations.                |
    |        min_num        | uint32                      |       | Minimum number of result to be returned.     |
    | aggregation_algorithm | Search.AggregationAlgorithm |       | Aggregation Algorithm                        |
    |         ratio         | google.protobuf.FloatValue  |       | Search ratio for agent return result number. |
    |        nprobe         | uint32                      |       | Search nprobe.                               |

  - Filter.Config

    |  field  | type          | label    | desc.                                      |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | desc.                |
    | :---: | :----- | :---- | :------------------- |
    | host  | string |       | The target hostname. |
    | port  | uint32 |       | The target port.     |

  - Filter.Config

    |  field  | type          | label    | desc.                                      |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | desc.                |
    | :---: | :----- | :---- | :------------------- |
    | host  | string |       | The target hostname. |
    | port  | uint32 |       | The target port.     |

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

    |  field   | type              | label | desc.                          |
    | :------: | :---------------- | :---- | :----------------------------- |
    | response | Search.Response   |       | Represent the search response. |
    |  status  | google.rpc.Status |       | The RPC error status.          |

  - Search.Response

    |   field    | type            | label    | desc.                  |
    | :--------: | :-------------- | :------- | :--------------------- |
    | request_id | string          |          | The unique request ID. |
    |  results   | Object.Distance | repeated | Search results.        |

  - Object.Distance

    |  field   | type   | label | desc.          |
    | :------: | :----- | :---- | :------------- |
    |    id    | string |       | The vector ID. |
    | distance | float  |       | The distance.  |

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

A method to search indexed vectors by multiple vectors in a single request.

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
  ```

  - Search.MultiRequest

    |  field   | type           | label    | desc.                                          |
    | :------: | :------------- | :------- | :--------------------------------------------- |
    | requests | Search.Request | repeated | Represent the multiple search request content. |

  - Search.Request

    | field  | type          | label    | desc.                                    |
    | :----: | :------------ | :------- | :--------------------------------------- |
    | vector | float         | repeated | The vector to be searched.               |
    | config | Search.Config |          | The configuration of the search request. |

  - Search.Config

    |         field         | type                        | label | desc.                                        |
    | :-------------------: | :-------------------------- | :---- | :------------------------------------------- |
    |      request_id       | string                      |       | Unique request ID.                           |
    |          num          | uint32                      |       | Maximum number of result to be returned.     |
    |        radius         | float                       |       | Search radius.                               |
    |        epsilon        | float                       |       | Search coefficient.                          |
    |        timeout        | int64                       |       | Search timeout in nanoseconds.               |
    |    ingress_filters    | Filter.Config               |       | Ingress filter configurations.               |
    |    egress_filters     | Filter.Config               |       | Egress filter configurations.                |
    |        min_num        | uint32                      |       | Minimum number of result to be returned.     |
    | aggregation_algorithm | Search.AggregationAlgorithm |       | Aggregation Algorithm                        |
    |         ratio         | google.protobuf.FloatValue  |       | Search ratio for agent return result number. |
    |        nprobe         | uint32                      |       | Search nprobe.                               |

  - Filter.Config

    |  field  | type          | label    | desc.                                      |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | desc.                |
    | :---: | :----- | :---- | :------------------- |
    | host  | string |       | The target hostname. |
    | port  | uint32 |       | The target port.     |

  - Filter.Config

    |  field  | type          | label    | desc.                                      |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | desc.                |
    | :---: | :----- | :---- | :------------------- |
    | host  | string |       | The target hostname. |
    | port  | uint32 |       | The target port.     |

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

    |   field   | type            | label    | desc.                                           |
    | :-------: | :-------------- | :------- | :---------------------------------------------- |
    | responses | Search.Response | repeated | Represent the multiple search response content. |

  - Search.Response

    |   field    | type            | label    | desc.                  |
    | :--------: | :-------------- | :------- | :--------------------- |
    | request_id | string          |          | The unique request ID. |
    |  results   | Object.Distance | repeated | Search results.        |

  - Object.Distance

    |  field   | type   | label | desc.          |
    | :------: | :----- | :---- | :------------- |
    |    id    | string |       | The vector ID. |
    | distance | float  |       | The distance.  |

### Status Code

| code | desc.             |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  13  | INTERNAL          |

## MultiSearchByID RPC

A method to search indexed vectors by multiple IDs in a single request.

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
  ```

  - Search.MultiIDRequest

    |  field   | type             | label    | desc.                                                |
    | :------: | :--------------- | :------- | :--------------------------------------------------- |
    | requests | Search.IDRequest | repeated | Represent the multiple search by ID request content. |

  - Search.IDRequest

    | field  | type          | label | desc.                                    |
    | :----: | :------------ | :---- | :--------------------------------------- |
    |   id   | string        |       | The vector ID to be searched.            |
    | config | Search.Config |       | The configuration of the search request. |

  - Search.Config

    |         field         | type                        | label | desc.                                        |
    | :-------------------: | :-------------------------- | :---- | :------------------------------------------- |
    |      request_id       | string                      |       | Unique request ID.                           |
    |          num          | uint32                      |       | Maximum number of result to be returned.     |
    |        radius         | float                       |       | Search radius.                               |
    |        epsilon        | float                       |       | Search coefficient.                          |
    |        timeout        | int64                       |       | Search timeout in nanoseconds.               |
    |    ingress_filters    | Filter.Config               |       | Ingress filter configurations.               |
    |    egress_filters     | Filter.Config               |       | Egress filter configurations.                |
    |        min_num        | uint32                      |       | Minimum number of result to be returned.     |
    | aggregation_algorithm | Search.AggregationAlgorithm |       | Aggregation Algorithm                        |
    |         ratio         | google.protobuf.FloatValue  |       | Search ratio for agent return result number. |
    |        nprobe         | uint32                      |       | Search nprobe.                               |

  - Filter.Config

    |  field  | type          | label    | desc.                                      |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | desc.                |
    | :---: | :----- | :---- | :------------------- |
    | host  | string |       | The target hostname. |
    | port  | uint32 |       | The target port.     |

  - Filter.Config

    |  field  | type          | label    | desc.                                      |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | desc.                |
    | :---: | :----- | :---- | :------------------- |
    | host  | string |       | The target hostname. |
    | port  | uint32 |       | The target port.     |

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

    |   field   | type            | label    | desc.                                           |
    | :-------: | :-------------- | :------- | :---------------------------------------------- |
    | responses | Search.Response | repeated | Represent the multiple search response content. |

  - Search.Response

    |   field    | type            | label    | desc.                  |
    | :--------: | :-------------- | :------- | :--------------------- |
    | request_id | string          |          | The unique request ID. |
    |  results   | Object.Distance | repeated | Search results.        |

  - Object.Distance

    |  field   | type   | label | desc.          |
    | :------: | :----- | :---- | :------------- |
    |    id    | string |       | The vector ID. |
    | distance | float  |       | The distance.  |

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

A method to linear search indexed vectors by a raw vector.

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
  ```

  - Search.Request

    | field  | type          | label    | desc.                                    |
    | :----: | :------------ | :------- | :--------------------------------------- |
    | vector | float         | repeated | The vector to be searched.               |
    | config | Search.Config |          | The configuration of the search request. |

  - Search.Config

    |         field         | type                        | label | desc.                                        |
    | :-------------------: | :-------------------------- | :---- | :------------------------------------------- |
    |      request_id       | string                      |       | Unique request ID.                           |
    |          num          | uint32                      |       | Maximum number of result to be returned.     |
    |        radius         | float                       |       | Search radius.                               |
    |        epsilon        | float                       |       | Search coefficient.                          |
    |        timeout        | int64                       |       | Search timeout in nanoseconds.               |
    |    ingress_filters    | Filter.Config               |       | Ingress filter configurations.               |
    |    egress_filters     | Filter.Config               |       | Egress filter configurations.                |
    |        min_num        | uint32                      |       | Minimum number of result to be returned.     |
    | aggregation_algorithm | Search.AggregationAlgorithm |       | Aggregation Algorithm                        |
    |         ratio         | google.protobuf.FloatValue  |       | Search ratio for agent return result number. |
    |        nprobe         | uint32                      |       | Search nprobe.                               |

  - Filter.Config

    |  field  | type          | label    | desc.                                      |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | desc.                |
    | :---: | :----- | :---- | :------------------- |
    | host  | string |       | The target hostname. |
    | port  | uint32 |       | The target port.     |

  - Filter.Config

    |  field  | type          | label    | desc.                                      |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | desc.                |
    | :---: | :----- | :---- | :------------------- |
    | host  | string |       | The target hostname. |
    | port  | uint32 |       | The target port.     |

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

    |   field    | type            | label    | desc.                  |
    | :--------: | :-------------- | :------- | :--------------------- |
    | request_id | string          |          | The unique request ID. |
    |  results   | Object.Distance | repeated | Search results.        |

  - Object.Distance

    |  field   | type   | label | desc.          |
    | :------: | :----- | :---- | :------------- |
    |    id    | string |       | The vector ID. |
    | distance | float  |       | The distance.  |

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

A method to linear search indexed vectors by ID.

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
  ```

  - Search.IDRequest

    | field  | type          | label | desc.                                    |
    | :----: | :------------ | :---- | :--------------------------------------- |
    |   id   | string        |       | The vector ID to be searched.            |
    | config | Search.Config |       | The configuration of the search request. |

  - Search.Config

    |         field         | type                        | label | desc.                                        |
    | :-------------------: | :-------------------------- | :---- | :------------------------------------------- |
    |      request_id       | string                      |       | Unique request ID.                           |
    |          num          | uint32                      |       | Maximum number of result to be returned.     |
    |        radius         | float                       |       | Search radius.                               |
    |        epsilon        | float                       |       | Search coefficient.                          |
    |        timeout        | int64                       |       | Search timeout in nanoseconds.               |
    |    ingress_filters    | Filter.Config               |       | Ingress filter configurations.               |
    |    egress_filters     | Filter.Config               |       | Egress filter configurations.                |
    |        min_num        | uint32                      |       | Minimum number of result to be returned.     |
    | aggregation_algorithm | Search.AggregationAlgorithm |       | Aggregation Algorithm                        |
    |         ratio         | google.protobuf.FloatValue  |       | Search ratio for agent return result number. |
    |        nprobe         | uint32                      |       | Search nprobe.                               |

  - Filter.Config

    |  field  | type          | label    | desc.                                      |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | desc.                |
    | :---: | :----- | :---- | :------------------- |
    | host  | string |       | The target hostname. |
    | port  | uint32 |       | The target port.     |

  - Filter.Config

    |  field  | type          | label    | desc.                                      |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | desc.                |
    | :---: | :----- | :---- | :------------------- |
    | host  | string |       | The target hostname. |
    | port  | uint32 |       | The target port.     |

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

    |   field    | type            | label    | desc.                  |
    | :--------: | :-------------- | :------- | :--------------------- |
    | request_id | string          |          | The unique request ID. |
    |  results   | Object.Distance | repeated | Search results.        |

  - Object.Distance

    |  field   | type   | label | desc.          |
    | :------: | :----- | :---- | :------------- |
    |    id    | string |       | The vector ID. |
    | distance | float  |       | The distance.  |

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

A method to linear search indexed vectors by multiple vectors.

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
  ```

  - Search.Request

    | field  | type          | label    | desc.                                    |
    | :----: | :------------ | :------- | :--------------------------------------- |
    | vector | float         | repeated | The vector to be searched.               |
    | config | Search.Config |          | The configuration of the search request. |

  - Search.Config

    |         field         | type                        | label | desc.                                        |
    | :-------------------: | :-------------------------- | :---- | :------------------------------------------- |
    |      request_id       | string                      |       | Unique request ID.                           |
    |          num          | uint32                      |       | Maximum number of result to be returned.     |
    |        radius         | float                       |       | Search radius.                               |
    |        epsilon        | float                       |       | Search coefficient.                          |
    |        timeout        | int64                       |       | Search timeout in nanoseconds.               |
    |    ingress_filters    | Filter.Config               |       | Ingress filter configurations.               |
    |    egress_filters     | Filter.Config               |       | Egress filter configurations.                |
    |        min_num        | uint32                      |       | Minimum number of result to be returned.     |
    | aggregation_algorithm | Search.AggregationAlgorithm |       | Aggregation Algorithm                        |
    |         ratio         | google.protobuf.FloatValue  |       | Search ratio for agent return result number. |
    |        nprobe         | uint32                      |       | Search nprobe.                               |

  - Filter.Config

    |  field  | type          | label    | desc.                                      |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | desc.                |
    | :---: | :----- | :---- | :------------------- |
    | host  | string |       | The target hostname. |
    | port  | uint32 |       | The target port.     |

  - Filter.Config

    |  field  | type          | label    | desc.                                      |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | desc.                |
    | :---: | :----- | :---- | :------------------- |
    | host  | string |       | The target hostname. |
    | port  | uint32 |       | The target port.     |

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

    |  field   | type              | label | desc.                          |
    | :------: | :---------------- | :---- | :----------------------------- |
    | response | Search.Response   |       | Represent the search response. |
    |  status  | google.rpc.Status |       | The RPC error status.          |

  - Search.Response

    |   field    | type            | label    | desc.                  |
    | :--------: | :-------------- | :------- | :--------------------- |
    | request_id | string          |          | The unique request ID. |
    |  results   | Object.Distance | repeated | Search results.        |

  - Object.Distance

    |  field   | type   | label | desc.          |
    | :------: | :----- | :---- | :------------- |
    |    id    | string |       | The vector ID. |
    | distance | float  |       | The distance.  |

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

A method to linear search indexed vectors by multiple IDs.

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
  ```

  - Search.IDRequest

    | field  | type          | label | desc.                                    |
    | :----: | :------------ | :---- | :--------------------------------------- |
    |   id   | string        |       | The vector ID to be searched.            |
    | config | Search.Config |       | The configuration of the search request. |

  - Search.Config

    |         field         | type                        | label | desc.                                        |
    | :-------------------: | :-------------------------- | :---- | :------------------------------------------- |
    |      request_id       | string                      |       | Unique request ID.                           |
    |          num          | uint32                      |       | Maximum number of result to be returned.     |
    |        radius         | float                       |       | Search radius.                               |
    |        epsilon        | float                       |       | Search coefficient.                          |
    |        timeout        | int64                       |       | Search timeout in nanoseconds.               |
    |    ingress_filters    | Filter.Config               |       | Ingress filter configurations.               |
    |    egress_filters     | Filter.Config               |       | Egress filter configurations.                |
    |        min_num        | uint32                      |       | Minimum number of result to be returned.     |
    | aggregation_algorithm | Search.AggregationAlgorithm |       | Aggregation Algorithm                        |
    |         ratio         | google.protobuf.FloatValue  |       | Search ratio for agent return result number. |
    |        nprobe         | uint32                      |       | Search nprobe.                               |

  - Filter.Config

    |  field  | type          | label    | desc.                                      |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | desc.                |
    | :---: | :----- | :---- | :------------------- |
    | host  | string |       | The target hostname. |
    | port  | uint32 |       | The target port.     |

  - Filter.Config

    |  field  | type          | label    | desc.                                      |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | desc.                |
    | :---: | :----- | :---- | :------------------- |
    | host  | string |       | The target hostname. |
    | port  | uint32 |       | The target port.     |

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

    |  field   | type              | label | desc.                          |
    | :------: | :---------------- | :---- | :----------------------------- |
    | response | Search.Response   |       | Represent the search response. |
    |  status  | google.rpc.Status |       | The RPC error status.          |

  - Search.Response

    |   field    | type            | label    | desc.                  |
    | :--------: | :-------------- | :------- | :--------------------- |
    | request_id | string          |          | The unique request ID. |
    |  results   | Object.Distance | repeated | Search results.        |

  - Object.Distance

    |  field   | type   | label | desc.          |
    | :------: | :----- | :---- | :------------- |
    |    id    | string |       | The vector ID. |
    | distance | float  |       | The distance.  |

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

A method to linear search indexed vectors by multiple vectors in a single
request.

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
  ```

  - Search.MultiRequest

    |  field   | type           | label    | desc.                                          |
    | :------: | :------------- | :------- | :--------------------------------------------- |
    | requests | Search.Request | repeated | Represent the multiple search request content. |

  - Search.Request

    | field  | type          | label    | desc.                                    |
    | :----: | :------------ | :------- | :--------------------------------------- |
    | vector | float         | repeated | The vector to be searched.               |
    | config | Search.Config |          | The configuration of the search request. |

  - Search.Config

    |         field         | type                        | label | desc.                                        |
    | :-------------------: | :-------------------------- | :---- | :------------------------------------------- |
    |      request_id       | string                      |       | Unique request ID.                           |
    |          num          | uint32                      |       | Maximum number of result to be returned.     |
    |        radius         | float                       |       | Search radius.                               |
    |        epsilon        | float                       |       | Search coefficient.                          |
    |        timeout        | int64                       |       | Search timeout in nanoseconds.               |
    |    ingress_filters    | Filter.Config               |       | Ingress filter configurations.               |
    |    egress_filters     | Filter.Config               |       | Egress filter configurations.                |
    |        min_num        | uint32                      |       | Minimum number of result to be returned.     |
    | aggregation_algorithm | Search.AggregationAlgorithm |       | Aggregation Algorithm                        |
    |         ratio         | google.protobuf.FloatValue  |       | Search ratio for agent return result number. |
    |        nprobe         | uint32                      |       | Search nprobe.                               |

  - Filter.Config

    |  field  | type          | label    | desc.                                      |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | desc.                |
    | :---: | :----- | :---- | :------------------- |
    | host  | string |       | The target hostname. |
    | port  | uint32 |       | The target port.     |

  - Filter.Config

    |  field  | type          | label    | desc.                                      |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | desc.                |
    | :---: | :----- | :---- | :------------------- |
    | host  | string |       | The target hostname. |
    | port  | uint32 |       | The target port.     |

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

    |   field   | type            | label    | desc.                                           |
    | :-------: | :-------------- | :------- | :---------------------------------------------- |
    | responses | Search.Response | repeated | Represent the multiple search response content. |

  - Search.Response

    |   field    | type            | label    | desc.                  |
    | :--------: | :-------------- | :------- | :--------------------- |
    | request_id | string          |          | The unique request ID. |
    |  results   | Object.Distance | repeated | Search results.        |

  - Object.Distance

    |  field   | type   | label | desc.          |
    | :------: | :----- | :---- | :------------- |
    |    id    | string |       | The vector ID. |
    | distance | float  |       | The distance.  |

### Status Code

| code | desc.             |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  13  | INTERNAL          |

## MultiLinearSearchByID RPC

A method to linear search indexed vectors by multiple IDs in a single
request.

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
  ```

  - Search.MultiIDRequest

    |  field   | type             | label    | desc.                                                |
    | :------: | :--------------- | :------- | :--------------------------------------------------- |
    | requests | Search.IDRequest | repeated | Represent the multiple search by ID request content. |

  - Search.IDRequest

    | field  | type          | label | desc.                                    |
    | :----: | :------------ | :---- | :--------------------------------------- |
    |   id   | string        |       | The vector ID to be searched.            |
    | config | Search.Config |       | The configuration of the search request. |

  - Search.Config

    |         field         | type                        | label | desc.                                        |
    | :-------------------: | :-------------------------- | :---- | :------------------------------------------- |
    |      request_id       | string                      |       | Unique request ID.                           |
    |          num          | uint32                      |       | Maximum number of result to be returned.     |
    |        radius         | float                       |       | Search radius.                               |
    |        epsilon        | float                       |       | Search coefficient.                          |
    |        timeout        | int64                       |       | Search timeout in nanoseconds.               |
    |    ingress_filters    | Filter.Config               |       | Ingress filter configurations.               |
    |    egress_filters     | Filter.Config               |       | Egress filter configurations.                |
    |        min_num        | uint32                      |       | Minimum number of result to be returned.     |
    | aggregation_algorithm | Search.AggregationAlgorithm |       | Aggregation Algorithm                        |
    |         ratio         | google.protobuf.FloatValue  |       | Search ratio for agent return result number. |
    |        nprobe         | uint32                      |       | Search nprobe.                               |

  - Filter.Config

    |  field  | type          | label    | desc.                                      |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | desc.                |
    | :---: | :----- | :---- | :------------------- |
    | host  | string |       | The target hostname. |
    | port  | uint32 |       | The target port.     |

  - Filter.Config

    |  field  | type          | label    | desc.                                      |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | desc.                |
    | :---: | :----- | :---- | :------------------- |
    | host  | string |       | The target hostname. |
    | port  | uint32 |       | The target port.     |

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

    |   field   | type            | label    | desc.                                           |
    | :-------: | :-------------- | :------- | :---------------------------------------------- |
    | responses | Search.Response | repeated | Represent the multiple search response content. |

  - Search.Response

    |   field    | type            | label    | desc.                  |
    | :--------: | :-------------- | :------- | :--------------------- |
    | request_id | string          |          | The unique request ID. |
    |  results   | Object.Distance | repeated | Search results.        |

  - Object.Distance

    |  field   | type   | label | desc.          |
    | :------: | :----- | :---- | :------------- |
    |    id    | string |       | The vector ID. |
    | distance | float  |       | The distance.  |

### Status Code

| code | desc.             |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  13  | INTERNAL          |
