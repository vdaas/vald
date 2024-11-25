# Vald Filter APIs

## Overview

Filter Server is responsible for providing insert, update, upsert and search interface for `Vald Filter Gateway`.

Vald Filter Gateway forward user request to user-defined ingress/egress filter components allowing user to run custom logic.

```rpc
service Filter {

  rpc SearchObject(payload.v1.Search.ObjectRequest) returns (payload.v1.Search.Response) {}
  rpc MultiSearchObject(payload.v1.Search.MultiObjectRequest) returns (payload.v1.Search.Responses) {}
  rpc StreamSearchObject(payload.v1.Search.ObjectRequest) returns (payload.v1.Search.StreamResponse) {}
  rpc InsertObject(payload.v1.Insert.ObjectRequest) returns (payload.v1.Object.Location) {}
  rpc StreamInsertObject(payload.v1.Insert.ObjectRequest) returns (payload.v1.Object.StreamLocation) {}
  rpc MultiInsertObject(payload.v1.Insert.MultiObjectRequest) returns (payload.v1.Object.Locations) {}
  rpc UpdateObject(payload.v1.Update.ObjectRequest) returns (payload.v1.Object.Location) {}
  rpc StreamUpdateObject(payload.v1.Update.ObjectRequest) returns (payload.v1.Object.StreamLocation) {}
  rpc MultiUpdateObject(payload.v1.Update.MultiObjectRequest) returns (payload.v1.Object.Locations) {}
  rpc UpsertObject(payload.v1.Upsert.ObjectRequest) returns (payload.v1.Object.Location) {}
  rpc StreamUpsertObject(payload.v1.Upsert.ObjectRequest) returns (payload.v1.Object.StreamLocation) {}
  rpc MultiUpsertObject(payload.v1.Upsert.MultiObjectRequest) returns (payload.v1.Object.Locations) {}

}
```

## SearchObject RPC

SearchObject RPC is the method to search object(s) similar to request object.

### Input

- the scheme of `payload.v1.Search.ObjectRequest`

  ```rpc
  message Search.ObjectRequest {
    bytes object = 1;
    Search.Config config = 2;
    Filter.Target vectorizer = 3;
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



  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }
  ```

  - Search.ObjectRequest

    |   field    | type          | label | desc.                                    |
    | :--------: | :------------ | :---- | :--------------------------------------- |
    |   object   | bytes         |       | The binary object to be searched.        |
    |   config   | Search.Config |       | The configuration of the search request. |
    | vectorizer | Filter.Target |       | Filter configuration.                    |

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

| code | description       |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  6   | ALREADY_EXISTS    |
|  13  | INTERNAL          |

Please refer to [Response Status Code](../status.md) for more details.

## MultiSearchObject RPC

StreamSearchObject RPC is the method to search vectors with multi queries(objects) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
By using the bidirectional streaming RPC, the search request can be communicated in any order between client and server.
Each Search request and response are independent.

### Input

- the scheme of `payload.v1.Search.MultiObjectRequest`

  ```rpc
  message Search.MultiObjectRequest {
    repeated Search.ObjectRequest requests = 1;
  }


  message Search.ObjectRequest {
    bytes object = 1;
    Search.Config config = 2;
    Filter.Target vectorizer = 3;
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



  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }
  ```

  - Search.MultiObjectRequest

    |  field   | type                 | label    | desc.                                                           |
    | :------: | :------------------- | :------- | :-------------------------------------------------------------- |
    | requests | Search.ObjectRequest | repeated | Represent the multiple search by binary object request content. |

  - Search.ObjectRequest

    |   field    | type          | label | desc.                                    |
    | :--------: | :------------ | :---- | :--------------------------------------- |
    |   object   | bytes         |       | The binary object to be searched.        |
    |   config   | Search.Config |       | The configuration of the search request. |
    | vectorizer | Filter.Target |       | Filter configuration.                    |

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

| code | description       |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  6   | ALREADY_EXISTS    |
|  13  | INTERNAL          |

Please refer to [Response Status Code](../status.md) for more details.

## StreamSearchObject RPC

MultiSearchObject RPC is the method to search objects with multiple objects in **1** request.

<div class="notice">
gRPC has a message size limitation.<br>
Please be careful that the size of the request exceeds the limit.
</div>

### Input

- the scheme of `payload.v1.Search.ObjectRequest`

  ```rpc
  message Search.ObjectRequest {
    bytes object = 1;
    Search.Config config = 2;
    Filter.Target vectorizer = 3;
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



  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }
  ```

  - Search.ObjectRequest

    |   field    | type          | label | desc.                                    |
    | :--------: | :------------ | :---- | :--------------------------------------- |
    |   object   | bytes         |       | The binary object to be searched.        |
    |   config   | Search.Config |       | The configuration of the search request. |
    | vectorizer | Filter.Target |       | Filter configuration.                    |

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

| code | description       |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  6   | ALREADY_EXISTS    |
|  13  | INTERNAL          |

Please refer to [Response Status Code](../status.md) for more details.

## InsertObject RPC

InsertObject RPC is the method to insert object through Vald Filter Gateway.

### Input

- the scheme of `payload.v1.Insert.ObjectRequest`

  ```rpc
  message Insert.ObjectRequest {
    Object.Blob object = 1;
    Insert.Config config = 2;
    Filter.Target vectorizer = 3;
  }


  message Object.Blob {
    string id = 1;
    bytes object = 2;
  }



  message Insert.Config {
    bool skip_strict_exist_check = 1;
    Filter.Config filters = 2;
    int64 timestamp = 3;
  }


  message Filter.Config {
    repeated Filter.Target targets = 1;
  }


  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }



  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }
  ```

  - Insert.ObjectRequest

    |   field    | type          | label | desc.                                    |
    | :--------: | :------------ | :---- | :--------------------------------------- |
    |   object   | Object.Blob   |       | The binary object to be inserted.        |
    |   config   | Insert.Config |       | The configuration of the insert request. |
    | vectorizer | Filter.Target |       | Filter configurations.                   |

  - Object.Blob

    | field  | type   | label | desc.              |
    | :----: | :----- | :---- | :----------------- |
    |   id   | string |       | The object ID.     |
    | object | bytes  |       | The binary object. |

  - Insert.Config

    |          field          | type          | label | desc.                                               |
    | :---------------------: | :------------ | :---- | :-------------------------------------------------- |
    | skip_strict_exist_check | bool          |       | A flag to skip exist check during insert operation. |
    |         filters         | Filter.Config |       | Filter configurations.                              |
    |        timestamp        | int64         |       | Insert timestamp.                                   |

  - Filter.Config

    |  field  | type          | label    | desc.                                      |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | desc.                |
    | :---: | :----- | :---- | :------------------- |
    | host  | string |       | The target hostname. |
    | port  | uint32 |       | The target port.     |

  - Filter.Target

    | field | type   | label | desc.                |
    | :---: | :----- | :---- | :------------------- |
    | host  | string |       | The target hostname. |
    | port  | uint32 |       | The target port.     |

### Output

- the scheme of `payload.v1.Object.Location`

  ```rpc
  message Object.Location {
    string name = 1;
    string uuid = 2;
    repeated string ips = 3;
  }
  ```

  - Object.Location

    | field | type   | label    | desc.                     |
    | :---: | :----- | :------- | :------------------------ |
    | name  | string |          | The name of the location. |
    | uuid  | string |          | The UUID of the vector.   |
    |  ips  | string | repeated | The IP list.              |

### Status Code

| code | description       |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  6   | ALREADY_EXISTS    |
|  13  | INTERNAL          |

Please refer to [Response Status Code](../status.md) for more details.

## StreamInsertObject RPC

StreamInsertObject RPC is the method to add new multiple object using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).

By using the bidirectional streaming RPC, the insert request can be communicated in any order between client and server.
Each Insert request and response are independent.
It's the recommended method to insert a large number of objects.

### Input

- the scheme of `payload.v1.Insert.ObjectRequest`

  ```rpc
  message Insert.ObjectRequest {
    Object.Blob object = 1;
    Insert.Config config = 2;
    Filter.Target vectorizer = 3;
  }


  message Object.Blob {
    string id = 1;
    bytes object = 2;
  }



  message Insert.Config {
    bool skip_strict_exist_check = 1;
    Filter.Config filters = 2;
    int64 timestamp = 3;
  }


  message Filter.Config {
    repeated Filter.Target targets = 1;
  }


  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }



  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }
  ```

  - Insert.ObjectRequest

    |   field    | type          | label | desc.                                    |
    | :--------: | :------------ | :---- | :--------------------------------------- |
    |   object   | Object.Blob   |       | The binary object to be inserted.        |
    |   config   | Insert.Config |       | The configuration of the insert request. |
    | vectorizer | Filter.Target |       | Filter configurations.                   |

  - Object.Blob

    | field  | type   | label | desc.              |
    | :----: | :----- | :---- | :----------------- |
    |   id   | string |       | The object ID.     |
    | object | bytes  |       | The binary object. |

  - Insert.Config

    |          field          | type          | label | desc.                                               |
    | :---------------------: | :------------ | :---- | :-------------------------------------------------- |
    | skip_strict_exist_check | bool          |       | A flag to skip exist check during insert operation. |
    |         filters         | Filter.Config |       | Filter configurations.                              |
    |        timestamp        | int64         |       | Insert timestamp.                                   |

  - Filter.Config

    |  field  | type          | label    | desc.                                      |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | desc.                |
    | :---: | :----- | :---- | :------------------- |
    | host  | string |       | The target hostname. |
    | port  | uint32 |       | The target port.     |

  - Filter.Target

    | field | type   | label | desc.                |
    | :---: | :----- | :---- | :------------------- |
    | host  | string |       | The target hostname. |
    | port  | uint32 |       | The target port.     |

### Output

- the scheme of `payload.v1.Object.StreamLocation`

  ```rpc
  message Object.StreamLocation {
    Object.Location location = 1;
    google.rpc.Status status = 2;
  }


  message Object.Location {
    string name = 1;
    string uuid = 2;
    repeated string ips = 3;
  }
  ```

  - Object.StreamLocation

    |  field   | type              | label | desc.                 |
    | :------: | :---------------- | :---- | :-------------------- |
    | location | Object.Location   |       | The vector location.  |
    |  status  | google.rpc.Status |       | The RPC error status. |

  - Object.Location

    | field | type   | label    | desc.                     |
    | :---: | :----- | :------- | :------------------------ |
    | name  | string |          | The name of the location. |
    | uuid  | string |          | The UUID of the vector.   |
    |  ips  | string | repeated | The IP list.              |

### Status Code

| code | description       |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  6   | ALREADY_EXISTS    |
|  13  | INTERNAL          |

Please refer to [Response Status Code](../status.md) for more details.

## MultiInsertObject RPC

MultiInsertObject RPC is the method to add multiple new objects in **1** request.

### Input

- the scheme of `payload.v1.Insert.MultiObjectRequest`

  ```rpc
  message Insert.MultiObjectRequest {
    repeated Insert.ObjectRequest requests = 1;
  }


  message Insert.ObjectRequest {
    Object.Blob object = 1;
    Insert.Config config = 2;
    Filter.Target vectorizer = 3;
  }


  message Object.Blob {
    string id = 1;
    bytes object = 2;
  }



  message Insert.Config {
    bool skip_strict_exist_check = 1;
    Filter.Config filters = 2;
    int64 timestamp = 3;
  }


  message Filter.Config {
    repeated Filter.Target targets = 1;
  }


  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }



  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }
  ```

  - Insert.MultiObjectRequest

    |  field   | type                 | label    | desc.                                        |
    | :------: | :------------------- | :------- | :------------------------------------------- |
    | requests | Insert.ObjectRequest | repeated | Represent multiple insert by object content. |

  - Insert.ObjectRequest

    |   field    | type          | label | desc.                                    |
    | :--------: | :------------ | :---- | :--------------------------------------- |
    |   object   | Object.Blob   |       | The binary object to be inserted.        |
    |   config   | Insert.Config |       | The configuration of the insert request. |
    | vectorizer | Filter.Target |       | Filter configurations.                   |

  - Object.Blob

    | field  | type   | label | desc.              |
    | :----: | :----- | :---- | :----------------- |
    |   id   | string |       | The object ID.     |
    | object | bytes  |       | The binary object. |

  - Insert.Config

    |          field          | type          | label | desc.                                               |
    | :---------------------: | :------------ | :---- | :-------------------------------------------------- |
    | skip_strict_exist_check | bool          |       | A flag to skip exist check during insert operation. |
    |         filters         | Filter.Config |       | Filter configurations.                              |
    |        timestamp        | int64         |       | Insert timestamp.                                   |

  - Filter.Config

    |  field  | type          | label    | desc.                                      |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | desc.                |
    | :---: | :----- | :---- | :------------------- |
    | host  | string |       | The target hostname. |
    | port  | uint32 |       | The target port.     |

  - Filter.Target

    | field | type   | label | desc.                |
    | :---: | :----- | :---- | :------------------- |
    | host  | string |       | The target hostname. |
    | port  | uint32 |       | The target port.     |

### Output

- the scheme of `payload.v1.Object.Locations`

  ```rpc
  message Object.Locations {
    repeated Object.Location locations = 1;
  }


  message Object.Location {
    string name = 1;
    string uuid = 2;
    repeated string ips = 3;
  }
  ```

  - Object.Locations

    |   field   | type            | label    | desc. |
    | :-------: | :-------------- | :------- | :---- |
    | locations | Object.Location | repeated |       |

  - Object.Location

    | field | type   | label    | desc.                     |
    | :---: | :----- | :------- | :------------------------ |
    | name  | string |          | The name of the location. |
    | uuid  | string |          | The UUID of the vector.   |
    |  ips  | string | repeated | The IP list.              |

### Status Code

| code | description       |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  6   | ALREADY_EXISTS    |
|  13  | INTERNAL          |

Please refer to [Response Status Code](../status.md) for more details.

## UpdateObject RPC

UpdateObject RPC is the method to update a single vector.

### Input

- the scheme of `payload.v1.Update.ObjectRequest`

  ```rpc
  message Update.ObjectRequest {
    Object.Blob object = 1;
    Update.Config config = 2;
    Filter.Target vectorizer = 3;
  }


  message Object.Blob {
    string id = 1;
    bytes object = 2;
  }



  message Update.Config {
    bool skip_strict_exist_check = 1;
    Filter.Config filters = 2;
    int64 timestamp = 3;
    bool disable_balanced_update = 4;
  }


  message Filter.Config {
    repeated Filter.Target targets = 1;
  }


  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }



  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }
  ```

  - Update.ObjectRequest

    |   field    | type          | label | desc.                                    |
    | :--------: | :------------ | :---- | :--------------------------------------- |
    |   object   | Object.Blob   |       | The binary object to be updated.         |
    |   config   | Update.Config |       | The configuration of the update request. |
    | vectorizer | Filter.Target |       | Filter target.                           |

  - Object.Blob

    | field  | type   | label | desc.              |
    | :----: | :----- | :---- | :----------------- |
    |   id   | string |       | The object ID.     |
    | object | bytes  |       | The binary object. |

  - Update.Config

        | field | type | label | desc. |
        | :---: | :--- | :---- | :---- |
        | skip_strict_exist_check | bool |  | A flag to skip exist check during update operation. |
        | filters | Filter.Config |  | Filter configuration. |
        | timestamp | int64 |  | Update timestamp. |
        | disable_balanced_update | bool |  | A flag to disable balanced update (split remove -> insert operation)

    during update operation. |

  - Filter.Config

    |  field  | type          | label    | desc.                                      |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | desc.                |
    | :---: | :----- | :---- | :------------------- |
    | host  | string |       | The target hostname. |
    | port  | uint32 |       | The target port.     |

  - Filter.Target

    | field | type   | label | desc.                |
    | :---: | :----- | :---- | :------------------- |
    | host  | string |       | The target hostname. |
    | port  | uint32 |       | The target port.     |

### Output

- the scheme of `payload.v1.Object.Location`

  ```rpc
  message Object.Location {
    string name = 1;
    string uuid = 2;
    repeated string ips = 3;
  }
  ```

  - Object.Location

    | field | type   | label    | desc.                     |
    | :---: | :----- | :------- | :------------------------ |
    | name  | string |          | The name of the location. |
    | uuid  | string |          | The UUID of the vector.   |
    |  ips  | string | repeated | The IP list.              |

### Status Code

| code | description       |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  6   | ALREADY_EXISTS    |
|  13  | INTERNAL          |

Please refer to [Response Status Code](../status.md) for more details.

## StreamUpdateObject RPC

StreamUpdateObject RPC is the method to update multiple objects using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
By using the bidirectional streaming RPC, the update request can be communicated in any order between client and server.
Each Update request and response are independent.
It's the recommended method to update the large amount of objects.

### Input

- the scheme of `payload.v1.Update.ObjectRequest`

  ```rpc
  message Update.ObjectRequest {
    Object.Blob object = 1;
    Update.Config config = 2;
    Filter.Target vectorizer = 3;
  }


  message Object.Blob {
    string id = 1;
    bytes object = 2;
  }



  message Update.Config {
    bool skip_strict_exist_check = 1;
    Filter.Config filters = 2;
    int64 timestamp = 3;
    bool disable_balanced_update = 4;
  }


  message Filter.Config {
    repeated Filter.Target targets = 1;
  }


  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }



  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }
  ```

  - Update.ObjectRequest

    |   field    | type          | label | desc.                                    |
    | :--------: | :------------ | :---- | :--------------------------------------- |
    |   object   | Object.Blob   |       | The binary object to be updated.         |
    |   config   | Update.Config |       | The configuration of the update request. |
    | vectorizer | Filter.Target |       | Filter target.                           |

  - Object.Blob

    | field  | type   | label | desc.              |
    | :----: | :----- | :---- | :----------------- |
    |   id   | string |       | The object ID.     |
    | object | bytes  |       | The binary object. |

  - Update.Config

        | field | type | label | desc. |
        | :---: | :--- | :---- | :---- |
        | skip_strict_exist_check | bool |  | A flag to skip exist check during update operation. |
        | filters | Filter.Config |  | Filter configuration. |
        | timestamp | int64 |  | Update timestamp. |
        | disable_balanced_update | bool |  | A flag to disable balanced update (split remove -> insert operation)

    during update operation. |

  - Filter.Config

    |  field  | type          | label    | desc.                                      |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | desc.                |
    | :---: | :----- | :---- | :------------------- |
    | host  | string |       | The target hostname. |
    | port  | uint32 |       | The target port.     |

  - Filter.Target

    | field | type   | label | desc.                |
    | :---: | :----- | :---- | :------------------- |
    | host  | string |       | The target hostname. |
    | port  | uint32 |       | The target port.     |

### Output

- the scheme of `payload.v1.Object.StreamLocation`

  ```rpc
  message Object.StreamLocation {
    Object.Location location = 1;
    google.rpc.Status status = 2;
  }


  message Object.Location {
    string name = 1;
    string uuid = 2;
    repeated string ips = 3;
  }
  ```

  - Object.StreamLocation

    |  field   | type              | label | desc.                 |
    | :------: | :---------------- | :---- | :-------------------- |
    | location | Object.Location   |       | The vector location.  |
    |  status  | google.rpc.Status |       | The RPC error status. |

  - Object.Location

    | field | type   | label    | desc.                     |
    | :---: | :----- | :------- | :------------------------ |
    | name  | string |          | The name of the location. |
    | uuid  | string |          | The UUID of the vector.   |
    |  ips  | string | repeated | The IP list.              |

### Status Code

| code | description       |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  6   | ALREADY_EXISTS    |
|  13  | INTERNAL          |

Please refer to [Response Status Code](../status.md) for more details.

## MultiUpdateObject RPC

MultiUpdateObject is the method to update multiple objects in **1** request.

<div class="notice">
gRPC has the message size limitation.<br>
Please be careful that the size of the request exceed the limit.
</div>

### Input

- the scheme of `payload.v1.Update.MultiObjectRequest`

  ```rpc
  message Update.MultiObjectRequest {
    repeated Update.ObjectRequest requests = 1;
  }


  message Update.ObjectRequest {
    Object.Blob object = 1;
    Update.Config config = 2;
    Filter.Target vectorizer = 3;
  }


  message Object.Blob {
    string id = 1;
    bytes object = 2;
  }



  message Update.Config {
    bool skip_strict_exist_check = 1;
    Filter.Config filters = 2;
    int64 timestamp = 3;
    bool disable_balanced_update = 4;
  }


  message Filter.Config {
    repeated Filter.Target targets = 1;
  }


  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }



  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }
  ```

  - Update.MultiObjectRequest

    |  field   | type                 | label    | desc.                                                 |
    | :------: | :------------------- | :------- | :---------------------------------------------------- |
    | requests | Update.ObjectRequest | repeated | Represent the multiple update object request content. |

  - Update.ObjectRequest

    |   field    | type          | label | desc.                                    |
    | :--------: | :------------ | :---- | :--------------------------------------- |
    |   object   | Object.Blob   |       | The binary object to be updated.         |
    |   config   | Update.Config |       | The configuration of the update request. |
    | vectorizer | Filter.Target |       | Filter target.                           |

  - Object.Blob

    | field  | type   | label | desc.              |
    | :----: | :----- | :---- | :----------------- |
    |   id   | string |       | The object ID.     |
    | object | bytes  |       | The binary object. |

  - Update.Config

        | field | type | label | desc. |
        | :---: | :--- | :---- | :---- |
        | skip_strict_exist_check | bool |  | A flag to skip exist check during update operation. |
        | filters | Filter.Config |  | Filter configuration. |
        | timestamp | int64 |  | Update timestamp. |
        | disable_balanced_update | bool |  | A flag to disable balanced update (split remove -> insert operation)

    during update operation. |

  - Filter.Config

    |  field  | type          | label    | desc.                                      |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | desc.                |
    | :---: | :----- | :---- | :------------------- |
    | host  | string |       | The target hostname. |
    | port  | uint32 |       | The target port.     |

  - Filter.Target

    | field | type   | label | desc.                |
    | :---: | :----- | :---- | :------------------- |
    | host  | string |       | The target hostname. |
    | port  | uint32 |       | The target port.     |

### Output

- the scheme of `payload.v1.Object.Locations`

  ```rpc
  message Object.Locations {
    repeated Object.Location locations = 1;
  }


  message Object.Location {
    string name = 1;
    string uuid = 2;
    repeated string ips = 3;
  }
  ```

  - Object.Locations

    |   field   | type            | label    | desc. |
    | :-------: | :-------------- | :------- | :---- |
    | locations | Object.Location | repeated |       |

  - Object.Location

    | field | type   | label    | desc.                     |
    | :---: | :----- | :------- | :------------------------ |
    | name  | string |          | The name of the location. |
    | uuid  | string |          | The UUID of the vector.   |
    |  ips  | string | repeated | The IP list.              |

### Status Code

| code | description |
| :--: | :---------- |

| 0 | OK |
| 1 | CANCELLED |
| 3 | INVALID_ARGUMENT |
| 4 | DEADLINE_EXCEEDED |
| 6 | ALREADY_EXISTS |
| 13 | INTERNAL |

Please refer to [Response Status Code](../status.md) for more details.

## UpsertObject RPC

UpsertObject RPC is the method to update a single object and add a new single object.

### Input

- the scheme of `payload.v1.Upsert.ObjectRequest`

  ```rpc
  message Upsert.ObjectRequest {
    Object.Blob object = 1;
    Upsert.Config config = 2;
    Filter.Target vectorizer = 3;
  }


  message Object.Blob {
    string id = 1;
    bytes object = 2;
  }



  message Upsert.Config {
    bool skip_strict_exist_check = 1;
    Filter.Config filters = 2;
    int64 timestamp = 3;
    bool disable_balanced_update = 4;
  }


  message Filter.Config {
    repeated Filter.Target targets = 1;
  }


  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }



  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }
  ```

  - Upsert.ObjectRequest

    |   field    | type          | label | desc.                                    |
    | :--------: | :------------ | :---- | :--------------------------------------- |
    |   object   | Object.Blob   |       | The binary object to be upserted.        |
    |   config   | Upsert.Config |       | The configuration of the upsert request. |
    | vectorizer | Filter.Target |       | Filter target.                           |

  - Object.Blob

    | field  | type   | label | desc.              |
    | :----: | :----- | :---- | :----------------- |
    |   id   | string |       | The object ID.     |
    | object | bytes  |       | The binary object. |

  - Upsert.Config

        | field | type | label | desc. |
        | :---: | :--- | :---- | :---- |
        | skip_strict_exist_check | bool |  | A flag to skip exist check during upsert operation. |
        | filters | Filter.Config |  | Filter configuration. |
        | timestamp | int64 |  | Upsert timestamp. |
        | disable_balanced_update | bool |  | A flag to disable balanced update (split remove -> insert operation)

    during update operation. |

  - Filter.Config

    |  field  | type          | label    | desc.                                      |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | desc.                |
    | :---: | :----- | :---- | :------------------- |
    | host  | string |       | The target hostname. |
    | port  | uint32 |       | The target port.     |

  - Filter.Target

    | field | type   | label | desc.                |
    | :---: | :----- | :---- | :------------------- |
    | host  | string |       | The target hostname. |
    | port  | uint32 |       | The target port.     |

### Output

- the scheme of `payload.v1.Object.Location`

  ```rpc
  message Object.Location {
    string name = 1;
    string uuid = 2;
    repeated string ips = 3;
  }
  ```

  - Object.Location

    | field | type   | label    | desc.                     |
    | :---: | :----- | :------- | :------------------------ |
    | name  | string |          | The name of the location. |
    | uuid  | string |          | The UUID of the vector.   |
    |  ips  | string | repeated | The IP list.              |

### Status Code

| code | description       |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  6   | ALREADY_EXISTS    |
|  13  | INTERNAL          |

Please refer to [Response Status Code](../status.md) for more details.

## StreamUpsertObject RPC

UpsertObject RPC is the method to update a single object and add a new single object.

### Input

- the scheme of `payload.v1.Upsert.ObjectRequest`

  ```rpc
  message Upsert.ObjectRequest {
    Object.Blob object = 1;
    Upsert.Config config = 2;
    Filter.Target vectorizer = 3;
  }


  message Object.Blob {
    string id = 1;
    bytes object = 2;
  }



  message Upsert.Config {
    bool skip_strict_exist_check = 1;
    Filter.Config filters = 2;
    int64 timestamp = 3;
    bool disable_balanced_update = 4;
  }


  message Filter.Config {
    repeated Filter.Target targets = 1;
  }


  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }



  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }
  ```

  - Upsert.ObjectRequest

    |   field    | type          | label | desc.                                    |
    | :--------: | :------------ | :---- | :--------------------------------------- |
    |   object   | Object.Blob   |       | The binary object to be upserted.        |
    |   config   | Upsert.Config |       | The configuration of the upsert request. |
    | vectorizer | Filter.Target |       | Filter target.                           |

  - Object.Blob

    | field  | type   | label | desc.              |
    | :----: | :----- | :---- | :----------------- |
    |   id   | string |       | The object ID.     |
    | object | bytes  |       | The binary object. |

  - Upsert.Config

        | field | type | label | desc. |
        | :---: | :--- | :---- | :---- |
        | skip_strict_exist_check | bool |  | A flag to skip exist check during upsert operation. |
        | filters | Filter.Config |  | Filter configuration. |
        | timestamp | int64 |  | Upsert timestamp. |
        | disable_balanced_update | bool |  | A flag to disable balanced update (split remove -> insert operation)

    during update operation. |

  - Filter.Config

    |  field  | type          | label    | desc.                                      |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | desc.                |
    | :---: | :----- | :---- | :------------------- |
    | host  | string |       | The target hostname. |
    | port  | uint32 |       | The target port.     |

  - Filter.Target

    | field | type   | label | desc.                |
    | :---: | :----- | :---- | :------------------- |
    | host  | string |       | The target hostname. |
    | port  | uint32 |       | The target port.     |

### Output

- the scheme of `payload.v1.Object.StreamLocation`

  ```rpc
  message Object.StreamLocation {
    Object.Location location = 1;
    google.rpc.Status status = 2;
  }


  message Object.Location {
    string name = 1;
    string uuid = 2;
    repeated string ips = 3;
  }
  ```

  - Object.StreamLocation

    |  field   | type              | label | desc.                 |
    | :------: | :---------------- | :---- | :-------------------- |
    | location | Object.Location   |       | The vector location.  |
    |  status  | google.rpc.Status |       | The RPC error status. |

  - Object.Location

    | field | type   | label    | desc.                     |
    | :---: | :----- | :------- | :------------------------ |
    | name  | string |          | The name of the location. |
    | uuid  | string |          | The UUID of the vector.   |
    |  ips  | string | repeated | The IP list.              |

### Status Code

| code | description       |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  6   | ALREADY_EXISTS    |
|  13  | INTERNAL          |

Please refer to [Response Status Code](../status.md) for more details.

## MultiUpsertObject RPC

MultiUpsertObject is the method to update existing multiple objects and add new multiple objects in **1** request.

<div class="notice">
gRPC has a message size limitation.<br>
Please be careful that the size of the request exceeds the limit.
</div>

### Input

- the scheme of `payload.v1.Upsert.MultiObjectRequest`

  ```rpc
  message Upsert.MultiObjectRequest {
    repeated Upsert.ObjectRequest requests = 1;
  }


  message Upsert.ObjectRequest {
    Object.Blob object = 1;
    Upsert.Config config = 2;
    Filter.Target vectorizer = 3;
  }


  message Object.Blob {
    string id = 1;
    bytes object = 2;
  }



  message Upsert.Config {
    bool skip_strict_exist_check = 1;
    Filter.Config filters = 2;
    int64 timestamp = 3;
    bool disable_balanced_update = 4;
  }


  message Filter.Config {
    repeated Filter.Target targets = 1;
  }


  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }



  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }
  ```

  - Upsert.MultiObjectRequest

    |  field   | type                 | label    | desc.                                                 |
    | :------: | :------------------- | :------- | :---------------------------------------------------- |
    | requests | Upsert.ObjectRequest | repeated | Represent the multiple upsert object request content. |

  - Upsert.ObjectRequest

    |   field    | type          | label | desc.                                    |
    | :--------: | :------------ | :---- | :--------------------------------------- |
    |   object   | Object.Blob   |       | The binary object to be upserted.        |
    |   config   | Upsert.Config |       | The configuration of the upsert request. |
    | vectorizer | Filter.Target |       | Filter target.                           |

  - Object.Blob

    | field  | type   | label | desc.              |
    | :----: | :----- | :---- | :----------------- |
    |   id   | string |       | The object ID.     |
    | object | bytes  |       | The binary object. |

  - Upsert.Config

        | field | type | label | desc. |
        | :---: | :--- | :---- | :---- |
        | skip_strict_exist_check | bool |  | A flag to skip exist check during upsert operation. |
        | filters | Filter.Config |  | Filter configuration. |
        | timestamp | int64 |  | Upsert timestamp. |
        | disable_balanced_update | bool |  | A flag to disable balanced update (split remove -> insert operation)

    during update operation. |

  - Filter.Config

    |  field  | type          | label    | desc.                                      |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | desc.                |
    | :---: | :----- | :---- | :------------------- |
    | host  | string |       | The target hostname. |
    | port  | uint32 |       | The target port.     |

  - Filter.Target

    | field | type   | label | desc.                |
    | :---: | :----- | :---- | :------------------- |
    | host  | string |       | The target hostname. |
    | port  | uint32 |       | The target port.     |

### Output

- the scheme of `payload.v1.Object.Locations`

  ```rpc
  message Object.Locations {
    repeated Object.Location locations = 1;
  }


  message Object.Location {
    string name = 1;
    string uuid = 2;
    repeated string ips = 3;
  }
  ```

  - Object.Locations

    |   field   | type            | label    | desc. |
    | :-------: | :-------------- | :------- | :---- |
    | locations | Object.Location | repeated |       |

  - Object.Location

    | field | type   | label    | desc.                     |
    | :---: | :----- | :------- | :------------------------ |
    | name  | string |          | The name of the location. |
    | uuid  | string |          | The UUID of the vector.   |
    |  ips  | string | repeated | The IP list.              |

### Status Code

| code | description       |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  6   | ALREADY_EXISTS    |
|  13  | INTERNAL          |

Please refer to [Response Status Code](../status.md) for more details.
