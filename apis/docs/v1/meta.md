# Vald InsertWithMetadata APIs

## Overview

InsertWithMetadata Service is responsible for inserting new vectors into the `vald-agent` and set metadata.

```rpc
service InsertWithMetadata {

  rpc InsertWithMetadata(payload.v1.Insert.Request) returns (payload.v1.Object.Location) {}
  rpc StreamInsertWithMetadata(payload.v1.Insert.Request) returns (payload.v1.Object.StreamLocation) {}
  rpc MultiInsertWithMetadata(payload.v1.Insert.MultiRequest) returns (payload.v1.Object.Locations) {}

}
```

## InsertWithMetadata RPC

InsertWithMetadata RPC is the method to add a new single vector and metadata.

### Input

- the scheme of `payload.v1.Insert.Request`

  ```rpc
  message Insert.Request {
    Object.Vector vector = 1;
    Insert.Config config = 2;
    optional bytes metadata = 3;
  }

  message Object.Vector {
    string id = 1;
    repeated float vector = 2;
    int64 timestamp = 3;
    optional bytes metadata = 4;
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

  ```

  - Insert.Request

    |  field   | type          | label    | description                                    |
    | :------: | :------------ | :------- | :--------------------------------------------- |
    |  vector  | Object.Vector |          | The vector to be inserted.                     |
    |  config  | Insert.Config |          | The configuration of the insert request.       |
    | metadata | bytes         | optional | The metadata is related to the request vector. |

  - Object.Vector

    |   field   | type   | label    | description                                     |
    | :-------: | :----- | :------- | :---------------------------------------------- |
    |    id     | string |          | The vector ID.                                  |
    |  vector   | float  | repeated | The vector.                                     |
    | timestamp | int64  |          | timestamp represents when this vector inserted. |
    | metadata  | bytes  | optional | The metadata is related to the request vector.  |

  - Insert.Config

    |          field          | type          | label | description                                         |
    | :---------------------: | :------------ | :---- | :-------------------------------------------------- |
    | skip_strict_exist_check | bool          |       | A flag to skip exist check during insert operation. |
    |         filters         | Filter.Config |       | Filter configurations.                              |
    |        timestamp        | int64         |       | Insert timestamp.                                   |

  - Filter.Config

    |  field  | type          | label    | description                                |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | description          |
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

    | field | type   | label    | description               |
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
|  5   | NOT_FOUND         |
|  13  | INTERNAL          |

Please refer to [Response Status Code](../status.md) for more details.

### Troubleshooting

The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name              | common reason                                                                                                                                       | how to resolve                                                                           |
| :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
| CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
| INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
| DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
| ALREADY_EXISTS    | Request ID is already inserted.                                                                                                                     | Change request ID.                                                                       |
| INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |

## StreamInsertWithMetadata RPC

StreamInsertWithMetadata RPC is the method to add new multiple vectors and metadata using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
Using the bidirectional streaming RPC, the insert request can be communicated in any order between client and server.
Each Insert request and response are independent.
It's the recommended method to insert a large number of vectors.

### Input

- the scheme of `payload.v1.Insert.Request`

  ```rpc
  message Insert.Request {
    Object.Vector vector = 1;
    Insert.Config config = 2;
    optional bytes metadata = 3;
  }

  message Object.Vector {
    string id = 1;
    repeated float vector = 2;
    int64 timestamp = 3;
    optional bytes metadata = 4;
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

  ```

  - Insert.Request

    |  field   | type          | label    | description                                    |
    | :------: | :------------ | :------- | :--------------------------------------------- |
    |  vector  | Object.Vector |          | The vector to be inserted.                     |
    |  config  | Insert.Config |          | The configuration of the insert request.       |
    | metadata | bytes         | optional | The metadata is related to the request vector. |

  - Object.Vector

    |   field   | type   | label    | description                                     |
    | :-------: | :----- | :------- | :---------------------------------------------- |
    |    id     | string |          | The vector ID.                                  |
    |  vector   | float  | repeated | The vector.                                     |
    | timestamp | int64  |          | timestamp represents when this vector inserted. |
    | metadata  | bytes  | optional | The metadata is related to the request vector.  |

  - Insert.Config

    |          field          | type          | label | description                                         |
    | :---------------------: | :------------ | :---- | :-------------------------------------------------- |
    | skip_strict_exist_check | bool          |       | A flag to skip exist check during insert operation. |
    |         filters         | Filter.Config |       | Filter configurations.                              |
    |        timestamp        | int64         |       | Insert timestamp.                                   |

  - Filter.Config

    |  field  | type          | label    | description                                |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | description          |
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

    |  field   | type              | label | description           |
    | :------: | :---------------- | :---- | :-------------------- |
    | location | Object.Location   |       | The vector location.  |
    |  status  | google.rpc.Status |       | The RPC error status. |

  - Object.Location

    | field | type   | label    | description               |
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
|  10  | ABORTED           |
|  13  | INTERNAL          |

Please refer to [Response Status Code](../status.md) for more details.

### Troubleshooting

The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name              | common reason                                                                                                                                       | how to resolve                                                                           |
| :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
| CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
| INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
| DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
| ALREADY_EXISTS    | Request ID is already inserted.                                                                                                                     | Change request ID.                                                                       |
| INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |

## MultiInsertWithMetadata RPC

MultiInsertWithMetadata RPC is the method to add multiple new vectors and metadata in **1** request.

<div class="notice">
gRPC has a message size limitation.<br>
Please be careful that the size of the request exceeds the limit.
</div>

### Input

- the scheme of `payload.v1.Insert.MultiRequest`

  ```rpc
  message Insert.MultiRequest {
    repeated Insert.Request requests = 1;
  }

  message Insert.Request {
    Object.Vector vector = 1;
    Insert.Config config = 2;
    optional bytes metadata = 3;
  }

  message Object.Vector {
    string id = 1;
    repeated float vector = 2;
    int64 timestamp = 3;
    optional bytes metadata = 4;
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

  ```

  - Insert.MultiRequest

    |  field   | type           | label    | description                                |
    | :------: | :------------- | :------- | :----------------------------------------- |
    | requests | Insert.Request | repeated | Represent multiple insert request content. |

  - Insert.Request

    |  field   | type          | label    | description                                    |
    | :------: | :------------ | :------- | :--------------------------------------------- |
    |  vector  | Object.Vector |          | The vector to be inserted.                     |
    |  config  | Insert.Config |          | The configuration of the insert request.       |
    | metadata | bytes         | optional | The metadata is related to the request vector. |

  - Object.Vector

    |   field   | type   | label    | description                                     |
    | :-------: | :----- | :------- | :---------------------------------------------- |
    |    id     | string |          | The vector ID.                                  |
    |  vector   | float  | repeated | The vector.                                     |
    | timestamp | int64  |          | timestamp represents when this vector inserted. |
    | metadata  | bytes  | optional | The metadata is related to the request vector.  |

  - Insert.Config

    |          field          | type          | label | description                                         |
    | :---------------------: | :------------ | :---- | :-------------------------------------------------- |
    | skip_strict_exist_check | bool          |       | A flag to skip exist check during insert operation. |
    |         filters         | Filter.Config |       | Filter configurations.                              |
    |        timestamp        | int64         |       | Insert timestamp.                                   |

  - Filter.Config

    |  field  | type          | label    | description                                |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | description          |
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

    |   field   | type            | label    | description |
    | :-------: | :-------------- | :------- | :---------- |
    | locations | Object.Location | repeated |             |

  - Object.Location

    | field | type   | label    | description               |
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
|  10  | ABORTED           |
|  13  | INTERNAL          |

Please refer to [Response Status Code](../status.md) for more details.

### Troubleshooting

The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name              | common reason                                                                                                                                       | how to resolve                                                                           |
| :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
| CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
| INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
| DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
| ALREADY_EXISTS    | Request ID is already inserted.                                                                                                                     | Change request ID.                                                                       |
| INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |

# Vald ObjectWithMetadata APIs

## Overview

ObjectWithMetadata Service is responsible for getting inserted vectors and metadata, and checking whether vectors are inserted into the `vald-agent`.

```rpc
service ObjectWithMetadata {

  rpc GetObjectWithMetadata(payload.v1.Object.VectorRequest) returns (payload.v1.Object.Vector) {}
  rpc StreamGetObjectWithMetadata(payload.v1.Object.VectorRequest) returns (payload.v1.Object.StreamVector) {}
  rpc StreamListObjectWithMetadata(payload.v1.Object.List.Request) returns (payload.v1.Object.List.Response) {}

}
```

## GetObjectWithMetadata RPC

GetObjectWithMetadata RPC is the method to get the metadata of a vector inserted into the `vald-agent` and metadata.

### Input

- the scheme of `payload.v1.Object.VectorRequest`

  ```rpc
  message Object.VectorRequest {
    Object.ID id = 1;
    Filter.Config filters = 2;
  }

  message Object.ID {
    string id = 1;
  }

  message Filter.Config {
    repeated Filter.Target targets = 1;
  }

  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }

  ```

  - Object.VectorRequest

    |  field  | type          | label | description                  |
    | :-----: | :------------ | :---- | :--------------------------- |
    |   id    | Object.ID     |       | The vector ID to be fetched. |
    | filters | Filter.Config |       | Filter configurations.       |

  - Object.ID

    | field | type   | label | description |
    | :---: | :----- | :---- | :---------- |
    |  id   | string |       |             |

  - Filter.Config

    |  field  | type          | label    | description                                |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | description          |
    | :---: | :----- | :---- | :------------------- |
    | host  | string |       | The target hostname. |
    | port  | uint32 |       | The target port.     |

### Output

- the scheme of `payload.v1.Object.Vector`

  ```rpc
  message Object.Vector {
    string id = 1;
    repeated float vector = 2;
    int64 timestamp = 3;
    optional bytes metadata = 4;
  }

  ```

  - Object.Vector

    |   field   | type   | label    | description                                     |
    | :-------: | :----- | :------- | :---------------------------------------------- |
    |    id     | string |          | The vector ID.                                  |
    |  vector   | float  | repeated | The vector.                                     |
    | timestamp | int64  |          | timestamp represents when this vector inserted. |
    | metadata  | bytes  | optional | The metadata is related to the request vector.  |

### Status Code

| code | description       |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  13  | INTERNAL          |

Please refer to [Response Status Code](../status.md) for more details.

### Troubleshooting

The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name              | common reason                                                                                   | how to resolve                                                                           |
| :---------------- | :---------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
| CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed.  |
| INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                         | Check request payload and fix request payload.                                           |
| DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
| NOT_FOUND         | Requested ID is NOT inserted.                                                                   | Send a request with an ID that is already inserted.                                      |
| INTERNAL          | Target Vald cluster or network route has some critical error.                                   | Check target Vald cluster first and check network route including ingress as second.     |

## StreamGetObjectWithMetadata RPC

StreamGetObjectWithMetadata RPC is the method to get the metadata of multiple existing vectors using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
Using the bidirectional streaming RPC, the GetObject request can be communicated in any order between client and server.
Each Upsert request and response are independent.

### Input

- the scheme of `payload.v1.Object.VectorRequest`

  ```rpc
  message Object.VectorRequest {
    Object.ID id = 1;
    Filter.Config filters = 2;
  }

  message Object.ID {
    string id = 1;
  }

  message Filter.Config {
    repeated Filter.Target targets = 1;
  }

  message Filter.Target {
    string host = 1;
    uint32 port = 2;
  }

  ```

  - Object.VectorRequest

    |  field  | type          | label | description                  |
    | :-----: | :------------ | :---- | :--------------------------- |
    |   id    | Object.ID     |       | The vector ID to be fetched. |
    | filters | Filter.Config |       | Filter configurations.       |

  - Object.ID

    | field | type   | label | description |
    | :---: | :----- | :---- | :---------- |
    |  id   | string |       |             |

  - Filter.Config

    |  field  | type          | label    | description                                |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | description          |
    | :---: | :----- | :---- | :------------------- |
    | host  | string |       | The target hostname. |
    | port  | uint32 |       | The target port.     |

### Output

- the scheme of `payload.v1.Object.StreamVector`

  ```rpc
  message Object.StreamVector {
    Object.Vector vector = 1;
    google.rpc.Status status = 2;
  }

  message Object.Vector {
    string id = 1;
    repeated float vector = 2;
    int64 timestamp = 3;
    optional bytes metadata = 4;
  }

  ```

  - Object.StreamVector

    | field  | type              | label | description           |
    | :----: | :---------------- | :---- | :-------------------- |
    | vector | Object.Vector     |       | The vector.           |
    | status | google.rpc.Status |       | The RPC error status. |

  - Object.Vector

    |   field   | type   | label    | description                                     |
    | :-------: | :----- | :------- | :---------------------------------------------- |
    |    id     | string |          | The vector ID.                                  |
    |  vector   | float  | repeated | The vector.                                     |
    | timestamp | int64  |          | timestamp represents when this vector inserted. |
    | metadata  | bytes  | optional | The metadata is related to the request vector.  |

### Status Code

| code | description       |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  13  | INTERNAL          |

Please refer to [Response Status Code](../status.md) for more details.

### Troubleshooting

The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name              | common reason                                                                                   | how to resolve                                                                           |
| :---------------- | :---------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
| CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed.  |
| INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                         | Check request payload and fix request payload.                                           |
| DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
| NOT_FOUND         | Requested ID is NOT inserted.                                                                   | Send a request with an ID that is already inserted.                                      |
| INTERNAL          | Target Vald cluster or network route has some critical error.                                   | Check target Vald cluster first and check network route including ingress as second.     |

## StreamListObjectWithMetadata RPC

A method to get all the vectors with server streaming

### Input

- the scheme of `payload.v1.Object.List.Request`

  ```rpc
  message Object.List.Request {
    // empty
  }

  ```

  - Object.List.Request

    empty

### Output

- the scheme of `payload.v1.Object.List.Response`

  ```rpc
  message Object.List.Response {
    Object.Vector vector = 1;
    google.rpc.Status status = 2;
  }

  message Object.Vector {
    string id = 1;
    repeated float vector = 2;
    int64 timestamp = 3;
    optional bytes metadata = 4;
  }

  ```

  - Object.List.Response

    | field  | type              | label | description           |
    | :----: | :---------------- | :---- | :-------------------- |
    | vector | Object.Vector     |       | The vector            |
    | status | google.rpc.Status |       | The RPC error status. |

  - Object.Vector

    |   field   | type   | label    | description                                     |
    | :-------: | :----- | :------- | :---------------------------------------------- |
    |    id     | string |          | The vector ID.                                  |
    |  vector   | float  | repeated | The vector.                                     |
    | timestamp | int64  |          | timestamp represents when this vector inserted. |
    | metadata  | bytes  | optional | The metadata is related to the request vector.  |

### Status Code

| code | description |
| :--: | :---------- |

TODO

Please refer to [Response Status Code](../status.md) for more details.

### Troubleshooting

TODO

# Vald RemoveWithMetadata APIs

## Overview

RemoveWithMetadata Service is responsible for removing vectors indexed in the `vald-agent`.

```rpc
service RemoveWithMetadata {

  rpc RemoveWithMetadata(payload.v1.Remove.Request) returns (payload.v1.Object.Location) {}
  rpc RemoveByTimestampWithMetadata(payload.v1.Remove.TimestampRequest) returns (payload.v1.Object.Locations) {}
  rpc StreamRemoveWithMetadata(payload.v1.Remove.Request) returns (payload.v1.Object.StreamLocation) {}
  rpc MultiRemoveWithMetadata(payload.v1.Remove.MultiRequest) returns (payload.v1.Object.Locations) {}

}
```

## RemoveWithMetadata RPC

RemoveWithMetadata RPC is the method to remove a single vector and metadata.

### Input

- the scheme of `payload.v1.Remove.Request`

  ```rpc
  message Remove.Request {
    Object.ID id = 1;
    Remove.Config config = 2;
  }

  message Object.ID {
    string id = 1;
  }

  message Remove.Config {
    bool skip_strict_exist_check = 1;
    int64 timestamp = 2;
  }

  ```

  - Remove.Request

    | field  | type          | label | description                              |
    | :----: | :------------ | :---- | :--------------------------------------- |
    |   id   | Object.ID     |       | The object ID to be removed.             |
    | config | Remove.Config |       | The configuration of the remove request. |

  - Object.ID

    | field | type   | label | description |
    | :---: | :----- | :---- | :---------- |
    |  id   | string |       |             |

  - Remove.Config

    |          field          | type  | label | description                                         |
    | :---------------------: | :---- | :---- | :-------------------------------------------------- |
    | skip_strict_exist_check | bool  |       | A flag to skip exist check during upsert operation. |
    |        timestamp        | int64 |       | Remove timestamp.                                   |

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

    | field | type   | label    | description               |
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
|  5   | NOT_FOUND         |
|  10  | ABORTED           |
|  13  | INTERNAL          |

Please refer to [Response Status Code](../status.md) for more details.

### Troubleshooting

The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name              | common reason                                                                                   | how to resolve                                                                           |
| :---------------- | :---------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
| CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed.  |
| INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                         | Check request payload and fix request payload.                                           |
| DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
| NOT_FOUND         | Requested ID is NOT inserted.                                                                   | Send a request with an ID that is already inserted.                                      |
| INTERNAL          | Target Vald cluster or network route has some critical error.                                   | Check target Vald cluster first and check network route including ingress as second.     |

## RemoveByTimestampWithMetadata RPC

RemoveByTimestampWithMetadata RPC is the method to remove vectors and metadata based on timestamp.

<div class="notice">
In the TimestampRequest message, the 'timestamps' field is repeated, allowing the inclusion of multiple Timestamp.<br>
When multiple Timestamps are provided, it results in an `AND` condition, enabling the realization of deletions with specified ranges.<br>
This design allows for versatile deletion operations, facilitating tasks such as removing data within a specific time range.
</div>

### Input

- the scheme of `payload.v1.Remove.TimestampRequest`

  ```rpc
  message Remove.TimestampRequest {
    repeated Remove.Timestamp timestamps = 1;
  }

  message Remove.Timestamp {
    int64 timestamp = 1;
    Remove.Timestamp.Operator operator = 2;
  }

  enum Remove.Timestamp.Operator {
    Eq = 0;
    Ne = 1;
    Ge = 2;
    Gt = 3;
    Le = 4;
    Lt = 5;
  }

  ```

  - Remove.TimestampRequest

        | field | type | label | description |
        | :---: | :--- | :---- | :---------- |
        | timestamps | Remove.Timestamp | repeated | The timestamp comparison list. If more than one is specified, the `AND`

    search is applied. |

  - Remove.Timestamp

    |   field   | type                      | label | description               |
    | :-------: | :------------------------ | :---- | :------------------------ |
    | timestamp | int64                     |       | The timestamp.            |
    | operator  | Remove.Timestamp.Operator |       | The conditional operator. |

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

    |   field   | type            | label    | description |
    | :-------: | :-------------- | :------- | :---------- |
    | locations | Object.Location | repeated |             |

  - Object.Location

    | field | type   | label    | description               |
    | :---: | :----- | :------- | :------------------------ |
    | name  | string |          | The name of the location. |
    | uuid  | string |          | The UUID of the vector.   |
    |  ips  | string | repeated | The IP list.              |

### Status Code

| code | description       |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  13  | INTERNAL          |

Please refer to [Response Status Code](../status.md) for more details.

### Troubleshooting

The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name              | common reason                                                                                   | how to resolve                                                                                                       |
| :---------------- | :---------------------------------------------------------------------------------------------- | :------------------------------------------------------------------------------------------------------------------- |
| CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed.                              |
| DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed.                             |
| NOT_FOUND         | No vectors in the system match the specified timestamp conditions.                              | Check whether vectors matching the specified timestamp conditions exist in the system, and fix conditions if needed. |
| INTERNAL          | Target Vald cluster or network route has some critical error.                                   | Check target Vald cluster first and check network route including ingress as second.                                 |

## StreamRemoveWithMetadata RPC

A method to remove multiple with metadata indexed vectors and metadata by bidirectional streaming.

StreamRemoveWithMetadata RPC is the method to remove multiple vectors using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
Using the bidirectional streaming RPC, the remove request can be communicated in any order between client and server.
Each Remove request and response are independent.
It's the recommended method to remove a large number of vectors.

### Input

- the scheme of `payload.v1.Remove.Request`

  ```rpc
  message Remove.Request {
    Object.ID id = 1;
    Remove.Config config = 2;
  }

  message Object.ID {
    string id = 1;
  }

  message Remove.Config {
    bool skip_strict_exist_check = 1;
    int64 timestamp = 2;
  }

  ```

  - Remove.Request

    | field  | type          | label | description                              |
    | :----: | :------------ | :---- | :--------------------------------------- |
    |   id   | Object.ID     |       | The object ID to be removed.             |
    | config | Remove.Config |       | The configuration of the remove request. |

  - Object.ID

    | field | type   | label | description |
    | :---: | :----- | :---- | :---------- |
    |  id   | string |       |             |

  - Remove.Config

    |          field          | type  | label | description                                         |
    | :---------------------: | :---- | :---- | :-------------------------------------------------- |
    | skip_strict_exist_check | bool  |       | A flag to skip exist check during upsert operation. |
    |        timestamp        | int64 |       | Remove timestamp.                                   |

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

    |  field   | type              | label | description           |
    | :------: | :---------------- | :---- | :-------------------- |
    | location | Object.Location   |       | The vector location.  |
    |  status  | google.rpc.Status |       | The RPC error status. |

  - Object.Location

    | field | type   | label    | description               |
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
|  5   | NOT_FOUND         |
|  10  | ABORTED           |
|  13  | INTERNAL          |

Please refer to [Response Status Code](../status.md) for more details.

### Troubleshooting

The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name              | common reason                                                                                   | how to resolve                                                                           |
| :---------------- | :---------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
| CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed.  |
| INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                         | Check request payload and fix request payload.                                           |
| DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
| NOT_FOUND         | Requested ID is NOT inserted.                                                                   | Send a request with an ID that is already inserted.                                      |
| INTERNAL          | Target Vald cluster or network route has some critical error.                                   | Check target Vald cluster first and check network route including ingress as second.     |

## MultiRemoveWithMetadata RPC

MultiRemoveWithMetadata is the method to remove multiple vectors and metadata in **1** request.

<div class="notice">
gRPC has a message size limitation.<br>
Please be careful that the size of the request exceeds the limit.
</div>

### Input

- the scheme of `payload.v1.Remove.MultiRequest`

  ```rpc
  message Remove.MultiRequest {
    repeated Remove.Request requests = 1;
  }

  message Remove.Request {
    Object.ID id = 1;
    Remove.Config config = 2;
  }

  message Object.ID {
    string id = 1;
  }

  message Remove.Config {
    bool skip_strict_exist_check = 1;
    int64 timestamp = 2;
  }

  ```

  - Remove.MultiRequest

    |  field   | type           | label    | description                                    |
    | :------: | :------------- | :------- | :--------------------------------------------- |
    | requests | Remove.Request | repeated | Represent the multiple remove request content. |

  - Remove.Request

    | field  | type          | label | description                              |
    | :----: | :------------ | :---- | :--------------------------------------- |
    |   id   | Object.ID     |       | The object ID to be removed.             |
    | config | Remove.Config |       | The configuration of the remove request. |

  - Object.ID

    | field | type   | label | description |
    | :---: | :----- | :---- | :---------- |
    |  id   | string |       |             |

  - Remove.Config

    |          field          | type  | label | description                                         |
    | :---------------------: | :---- | :---- | :-------------------------------------------------- |
    | skip_strict_exist_check | bool  |       | A flag to skip exist check during upsert operation. |
    |        timestamp        | int64 |       | Remove timestamp.                                   |

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

    |   field   | type            | label    | description |
    | :-------: | :-------------- | :------- | :---------- |
    | locations | Object.Location | repeated |             |

  - Object.Location

    | field | type   | label    | description               |
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
|  5   | NOT_FOUND         |
|  10  | ABORTED           |
|  13  | INTERNAL          |

Please refer to [Response Status Code](../status.md) for more details.

### Troubleshooting

The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name              | common reason                                                                                   | how to resolve                                                                           |
| :---------------- | :---------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
| CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed.  |
| INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                         | Check request payload and fix request payload.                                           |
| DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
| NOT_FOUND         | Requested ID is NOT inserted.                                                                   | Send a request with an ID that is already inserted.                                      |
| INTERNAL          | Target Vald cluster or network route has some critical error.                                   | Check target Vald cluster first and check network route including ingress as second.     |

# Vald SearchWithMetadata APIs

## Overview

SearchWithMetadata Service is responsible for searching vectors similar to the user request vector and getting metadata from `vald-lb-gateway`.

```rpc
service SearchWithMetadata {

  rpc SearchWithMetadata(payload.v1.Search.Request) returns (payload.v1.Search.Response) {}
  rpc SearchByIDWithMetadata(payload.v1.Search.IDRequest) returns (payload.v1.Search.Response) {}
  rpc StreamSearchWithMetadata(payload.v1.Search.Request) returns (payload.v1.Search.StreamResponse) {}
  rpc StreamSearchByIDWithMetadata(payload.v1.Search.IDRequest) returns (payload.v1.Search.StreamResponse) {}
  rpc MultiSearchWithMetadata(payload.v1.Search.MultiRequest) returns (payload.v1.Search.Responses) {}
  rpc MultiSearchByIDWithMetadata(payload.v1.Search.MultiIDRequest) returns (payload.v1.Search.Responses) {}
  rpc LinearSearchWithMetadata(payload.v1.Search.Request) returns (payload.v1.Search.Response) {}
  rpc LinearSearchByIDWithMetadata(payload.v1.Search.IDRequest) returns (payload.v1.Search.Response) {}
  rpc StreamLinearSearchWithMetadata(payload.v1.Search.Request) returns (payload.v1.Search.StreamResponse) {}
  rpc StreamLinearSearchByIDWithMetadata(payload.v1.Search.IDRequest) returns (payload.v1.Search.StreamResponse) {}
  rpc MultiLinearSearchWithMetadata(payload.v1.Search.MultiRequest) returns (payload.v1.Search.Responses) {}
  rpc MultiLinearSearchByIDWithMetadata(payload.v1.Search.MultiIDRequest) returns (payload.v1.Search.Responses) {}

}
```

## SearchWithMetadata RPC

SearchWithMetadata RPC is the method to search vector(s) similar to the request vector and to get metadata(s).

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

  enum Search.AggregationAlgorithm {
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

  - Search.Request

    | field  | type          | label    | description                              |
    | :----: | :------------ | :------- | :--------------------------------------- |
    | vector | float         | repeated | The vector to be searched.               |
    | config | Search.Config |          | The configuration of the search request. |

  - Search.Config

    |         field         | type                        | label | description                                  |
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

    |  field  | type          | label    | description                                |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | description          |
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
    optional bytes metadata = 3;
  }

  ```

  - Search.Response

    |   field    | type            | label    | description            |
    | :--------: | :-------------- | :------- | :--------------------- |
    | request_id | string          |          | The unique request ID. |
    |  results   | Object.Distance | repeated | Search results.        |

  - Object.Distance

    |  field   | type   | label    | description                                    |
    | :------: | :----- | :------- | :--------------------------------------------- |
    |    id    | string |          | The vector ID.                                 |
    | distance | float  |          | The distance.                                  |
    | metadata | bytes  | optional | The metadata is related to the request vector. |

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

## SearchByIDWithMetadata RPC

SearchByIDWithMetadata RPC is the method to search similar vectors using a user-defined vector ID and to get metadata.<br>
The vector with the same requested ID should be indexed into the `vald-lb-gateway` before searching.

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

  enum Search.AggregationAlgorithm {
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

  - Search.IDRequest

    | field  | type          | label | description                              |
    | :----: | :------------ | :---- | :--------------------------------------- |
    |   id   | string        |       | The vector ID to be searched.            |
    | config | Search.Config |       | The configuration of the search request. |

  - Search.Config

    |         field         | type                        | label | description                                  |
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

    |  field  | type          | label    | description                                |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | description          |
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
    optional bytes metadata = 3;
  }

  ```

  - Search.Response

    |   field    | type            | label    | description            |
    | :--------: | :-------------- | :------- | :--------------------- |
    | request_id | string          |          | The unique request ID. |
    |  results   | Object.Distance | repeated | Search results.        |

  - Object.Distance

    |  field   | type   | label    | description                                    |
    | :------: | :----- | :------- | :--------------------------------------------- |
    |    id    | string |          | The vector ID.                                 |
    | distance | float  |          | The distance.                                  |
    | metadata | bytes  | optional | The metadata is related to the request vector. |

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

## StreamSearchWithMetadata RPC

StreamSearchWithMetadata RPC is the method to search vectors and to get metadata with multi queries(vectors) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
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

  enum Search.AggregationAlgorithm {
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

  - Search.Request

    | field  | type          | label    | description                              |
    | :----: | :------------ | :------- | :--------------------------------------- |
    | vector | float         | repeated | The vector to be searched.               |
    | config | Search.Config |          | The configuration of the search request. |

  - Search.Config

    |         field         | type                        | label | description                                  |
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

    |  field  | type          | label    | description                                |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | description          |
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
    optional bytes metadata = 3;
  }

  ```

  - Search.StreamResponse

    |  field   | type              | label | description                    |
    | :------: | :---------------- | :---- | :----------------------------- |
    | response | Search.Response   |       | Represent the search response. |
    |  status  | google.rpc.Status |       | The RPC error status.          |

  - Search.Response

    |   field    | type            | label    | description            |
    | :--------: | :-------------- | :------- | :--------------------- |
    | request_id | string          |          | The unique request ID. |
    |  results   | Object.Distance | repeated | Search results.        |

  - Object.Distance

    |  field   | type   | label    | description                                    |
    | :------: | :----- | :------- | :--------------------------------------------- |
    |    id    | string |          | The vector ID.                                 |
    | distance | float  |          | The distance.                                  |
    | metadata | bytes  | optional | The metadata is related to the request vector. |

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

## StreamSearchByIDWithMetadata RPC

StreamSearchByIDWithMetadata RPC is the method to search vectors and to get metadata with multi queries(IDs) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
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

  enum Search.AggregationAlgorithm {
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

  - Search.IDRequest

    | field  | type          | label | description                              |
    | :----: | :------------ | :---- | :--------------------------------------- |
    |   id   | string        |       | The vector ID to be searched.            |
    | config | Search.Config |       | The configuration of the search request. |

  - Search.Config

    |         field         | type                        | label | description                                  |
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

    |  field  | type          | label    | description                                |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | description          |
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
    optional bytes metadata = 3;
  }

  ```

  - Search.StreamResponse

    |  field   | type              | label | description                    |
    | :------: | :---------------- | :---- | :----------------------------- |
    | response | Search.Response   |       | Represent the search response. |
    |  status  | google.rpc.Status |       | The RPC error status.          |

  - Search.Response

    |   field    | type            | label    | description            |
    | :--------: | :-------------- | :------- | :--------------------- |
    | request_id | string          |          | The unique request ID. |
    |  results   | Object.Distance | repeated | Search results.        |

  - Object.Distance

    |  field   | type   | label    | description                                    |
    | :------: | :----- | :------- | :--------------------------------------------- |
    |    id    | string |          | The vector ID.                                 |
    | distance | float  |          | The distance.                                  |
    | metadata | bytes  | optional | The metadata is related to the request vector. |

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

## MultiSearchWithMetadata RPC

MultiSearchWithMetadata RPC is the method to search vectors and to get metadata with multiple vectors in **1** request.

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

  enum Search.AggregationAlgorithm {
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

  - Search.MultiRequest

    |  field   | type           | label    | description                                    |
    | :------: | :------------- | :------- | :--------------------------------------------- |
    | requests | Search.Request | repeated | Represent the multiple search request content. |

  - Search.Request

    | field  | type          | label    | description                              |
    | :----: | :------------ | :------- | :--------------------------------------- |
    | vector | float         | repeated | The vector to be searched.               |
    | config | Search.Config |          | The configuration of the search request. |

  - Search.Config

    |         field         | type                        | label | description                                  |
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

    |  field  | type          | label    | description                                |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | description          |
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
    optional bytes metadata = 3;
  }

  ```

  - Search.Responses

    |   field   | type            | label    | description                                     |
    | :-------: | :-------------- | :------- | :---------------------------------------------- |
    | responses | Search.Response | repeated | Represent the multiple search response content. |

  - Search.Response

    |   field    | type            | label    | description            |
    | :--------: | :-------------- | :------- | :--------------------- |
    | request_id | string          |          | The unique request ID. |
    |  results   | Object.Distance | repeated | Search results.        |

  - Object.Distance

    |  field   | type   | label    | description                                    |
    | :------: | :----- | :------- | :--------------------------------------------- |
    |    id    | string |          | The vector ID.                                 |
    | distance | float  |          | The distance.                                  |
    | metadata | bytes  | optional | The metadata is related to the request vector. |

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

## MultiSearchByIDWithMetadata RPC

MultiSearchByIDWithMetadata RPC is the method to search vectors and to get metadata with multiple IDs in **1** request.

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

  enum Search.AggregationAlgorithm {
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

  - Search.MultiIDRequest

    |  field   | type             | label    | description                                          |
    | :------: | :--------------- | :------- | :--------------------------------------------------- |
    | requests | Search.IDRequest | repeated | Represent the multiple search by ID request content. |

  - Search.IDRequest

    | field  | type          | label | description                              |
    | :----: | :------------ | :---- | :--------------------------------------- |
    |   id   | string        |       | The vector ID to be searched.            |
    | config | Search.Config |       | The configuration of the search request. |

  - Search.Config

    |         field         | type                        | label | description                                  |
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

    |  field  | type          | label    | description                                |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | description          |
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
    optional bytes metadata = 3;
  }

  ```

  - Search.Responses

    |   field   | type            | label    | description                                     |
    | :-------: | :-------------- | :------- | :---------------------------------------------- |
    | responses | Search.Response | repeated | Represent the multiple search response content. |

  - Search.Response

    |   field    | type            | label    | description            |
    | :--------: | :-------------- | :------- | :--------------------- |
    | request_id | string          |          | The unique request ID. |
    |  results   | Object.Distance | repeated | Search results.        |

  - Object.Distance

    |  field   | type   | label    | description                                    |
    | :------: | :----- | :------- | :--------------------------------------------- |
    |    id    | string |          | The vector ID.                                 |
    | distance | float  |          | The distance.                                  |
    | metadata | bytes  | optional | The metadata is related to the request vector. |

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

## LinearSearchWithMetadata RPC

LinearSearchWithMetadata RPC is the method to linear search vector(s) similar to the request vector and to get metadata.

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

  enum Search.AggregationAlgorithm {
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

  - Search.Request

    | field  | type          | label    | description                              |
    | :----: | :------------ | :------- | :--------------------------------------- |
    | vector | float         | repeated | The vector to be searched.               |
    | config | Search.Config |          | The configuration of the search request. |

  - Search.Config

    |         field         | type                        | label | description                                  |
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

    |  field  | type          | label    | description                                |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | description          |
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
    optional bytes metadata = 3;
  }

  ```

  - Search.Response

    |   field    | type            | label    | description            |
    | :--------: | :-------------- | :------- | :--------------------- |
    | request_id | string          |          | The unique request ID. |
    |  results   | Object.Distance | repeated | Search results.        |

  - Object.Distance

    |  field   | type   | label    | description                                    |
    | :------: | :----- | :------- | :--------------------------------------------- |
    |    id    | string |          | The vector ID.                                 |
    | distance | float  |          | The distance.                                  |
    | metadata | bytes  | optional | The metadata is related to the request vector. |

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

## LinearSearchByIDWithMetadata RPC

LinearSearchByIDWithMetadata RPC is the method to linear search similar vectors using a user-defined vector ID and to get metadata.<br>
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

  enum Search.AggregationAlgorithm {
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

  - Search.IDRequest

    | field  | type          | label | description                              |
    | :----: | :------------ | :---- | :--------------------------------------- |
    |   id   | string        |       | The vector ID to be searched.            |
    | config | Search.Config |       | The configuration of the search request. |

  - Search.Config

    |         field         | type                        | label | description                                  |
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

    |  field  | type          | label    | description                                |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | description          |
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
    optional bytes metadata = 3;
  }

  ```

  - Search.Response

    |   field    | type            | label    | description            |
    | :--------: | :-------------- | :------- | :--------------------- |
    | request_id | string          |          | The unique request ID. |
    |  results   | Object.Distance | repeated | Search results.        |

  - Object.Distance

    |  field   | type   | label    | description                                    |
    | :------: | :----- | :------- | :--------------------------------------------- |
    |    id    | string |          | The vector ID.                                 |
    | distance | float  |          | The distance.                                  |
    | metadata | bytes  | optional | The metadata is related to the request vector. |

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

## StreamLinearSearchWithMetadata RPC

StreamLinearSearchWithMetadata RPC is the method to linear search vectors and to get metadata with multi queries(vectors) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
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

  enum Search.AggregationAlgorithm {
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

  - Search.Request

    | field  | type          | label    | description                              |
    | :----: | :------------ | :------- | :--------------------------------------- |
    | vector | float         | repeated | The vector to be searched.               |
    | config | Search.Config |          | The configuration of the search request. |

  - Search.Config

    |         field         | type                        | label | description                                  |
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

    |  field  | type          | label    | description                                |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | description          |
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
    optional bytes metadata = 3;
  }

  ```

  - Search.StreamResponse

    |  field   | type              | label | description                    |
    | :------: | :---------------- | :---- | :----------------------------- |
    | response | Search.Response   |       | Represent the search response. |
    |  status  | google.rpc.Status |       | The RPC error status.          |

  - Search.Response

    |   field    | type            | label    | description            |
    | :--------: | :-------------- | :------- | :--------------------- |
    | request_id | string          |          | The unique request ID. |
    |  results   | Object.Distance | repeated | Search results.        |

  - Object.Distance

    |  field   | type   | label    | description                                    |
    | :------: | :----- | :------- | :--------------------------------------------- |
    |    id    | string |          | The vector ID.                                 |
    | distance | float  |          | The distance.                                  |
    | metadata | bytes  | optional | The metadata is related to the request vector. |

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

## StreamLinearSearchByIDWithMetadata RPC

StreamLinearSearchByIDWithMetadata RPC is the method to linear search vectors and to get metadata with multi queries(IDs) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
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

  enum Search.AggregationAlgorithm {
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

  - Search.IDRequest

    | field  | type          | label | description                              |
    | :----: | :------------ | :---- | :--------------------------------------- |
    |   id   | string        |       | The vector ID to be searched.            |
    | config | Search.Config |       | The configuration of the search request. |

  - Search.Config

    |         field         | type                        | label | description                                  |
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

    |  field  | type          | label    | description                                |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | description          |
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
    optional bytes metadata = 3;
  }

  ```

  - Search.StreamResponse

    |  field   | type              | label | description                    |
    | :------: | :---------------- | :---- | :----------------------------- |
    | response | Search.Response   |       | Represent the search response. |
    |  status  | google.rpc.Status |       | The RPC error status.          |

  - Search.Response

    |   field    | type            | label    | description            |
    | :--------: | :-------------- | :------- | :--------------------- |
    | request_id | string          |          | The unique request ID. |
    |  results   | Object.Distance | repeated | Search results.        |

  - Object.Distance

    |  field   | type   | label    | description                                    |
    | :------: | :----- | :------- | :--------------------------------------------- |
    |    id    | string |          | The vector ID.                                 |
    | distance | float  |          | The distance.                                  |
    | metadata | bytes  | optional | The metadata is related to the request vector. |

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

## MultiLinearSearchWithMetadata RPC

MultiLinearSearchWithMetadata RPC is the method to linear search vectors and to get metadata with multiple vectors in **1** request.

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

  enum Search.AggregationAlgorithm {
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

  - Search.MultiRequest

    |  field   | type           | label    | description                                    |
    | :------: | :------------- | :------- | :--------------------------------------------- |
    | requests | Search.Request | repeated | Represent the multiple search request content. |

  - Search.Request

    | field  | type          | label    | description                              |
    | :----: | :------------ | :------- | :--------------------------------------- |
    | vector | float         | repeated | The vector to be searched.               |
    | config | Search.Config |          | The configuration of the search request. |

  - Search.Config

    |         field         | type                        | label | description                                  |
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

    |  field  | type          | label    | description                                |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | description          |
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
    optional bytes metadata = 3;
  }

  ```

  - Search.Responses

    |   field   | type            | label    | description                                     |
    | :-------: | :-------------- | :------- | :---------------------------------------------- |
    | responses | Search.Response | repeated | Represent the multiple search response content. |

  - Search.Response

    |   field    | type            | label    | description            |
    | :--------: | :-------------- | :------- | :--------------------- |
    | request_id | string          |          | The unique request ID. |
    |  results   | Object.Distance | repeated | Search results.        |

  - Object.Distance

    |  field   | type   | label    | description                                    |
    | :------: | :----- | :------- | :--------------------------------------------- |
    |    id    | string |          | The vector ID.                                 |
    | distance | float  |          | The distance.                                  |
    | metadata | bytes  | optional | The metadata is related to the request vector. |

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

## MultiLinearSearchByIDWithMetadata RPC

MultiLinearSearchByIDWithMetadata RPC is the method to linear search vectors and to get metadata with multiple IDs in **1** request.

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

  enum Search.AggregationAlgorithm {
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

  - Search.MultiIDRequest

    |  field   | type             | label    | description                                          |
    | :------: | :--------------- | :------- | :--------------------------------------------------- |
    | requests | Search.IDRequest | repeated | Represent the multiple search by ID request content. |

  - Search.IDRequest

    | field  | type          | label | description                              |
    | :----: | :------------ | :---- | :--------------------------------------- |
    |   id   | string        |       | The vector ID to be searched.            |
    | config | Search.Config |       | The configuration of the search request. |

  - Search.Config

    |         field         | type                        | label | description                                  |
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

    |  field  | type          | label    | description                                |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | description          |
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
    optional bytes metadata = 3;
  }

  ```

  - Search.Responses

    |   field   | type            | label    | description                                     |
    | :-------: | :-------------- | :------- | :---------------------------------------------- |
    | responses | Search.Response | repeated | Represent the multiple search response content. |

  - Search.Response

    |   field    | type            | label    | description            |
    | :--------: | :-------------- | :------- | :--------------------- |
    | request_id | string          |          | The unique request ID. |
    |  results   | Object.Distance | repeated | Search results.        |

  - Object.Distance

    |  field   | type   | label    | description                                    |
    | :------: | :----- | :------- | :--------------------------------------------- |
    |    id    | string |          | The vector ID.                                 |
    | distance | float  |          | The distance.                                  |
    | metadata | bytes  | optional | The metadata is related to the request vector. |

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

# Vald UpdateWithMetadata APIs

## Overview

UpdateWithMetadata Service updates to new vector from inserted vector in the `vald-agent` components and updates to new metadata from set metadata.

```rpc
service UpdateWithMetadata {

  rpc UpdateWithMetadata(payload.v1.Update.Request) returns (payload.v1.Object.Location) {}
  rpc StreamUpdateWithMetadata(payload.v1.Update.Request) returns (payload.v1.Object.StreamLocation) {}
  rpc MultiUpdateWithMetadata(payload.v1.Update.MultiRequest) returns (payload.v1.Object.Locations) {}
  rpc UpdateTimestampWithMetadata(payload.v1.Update.TimestampRequest) returns (payload.v1.Object.Location) {}

}
```

## UpdateWithMetadata RPC

UpdateWithMetadata RPC is the method to update a single vector.

### Input

- the scheme of `payload.v1.Update.Request`

  ```rpc
  message Update.Request {
    Object.Vector vector = 1;
    Update.Config config = 2;
    optional bytes metadata = 3;
  }

  message Object.Vector {
    string id = 1;
    repeated float vector = 2;
    int64 timestamp = 3;
    optional bytes metadata = 4;
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

  ```

  - Update.Request

    |  field   | type          | label    | description                                    |
    | :------: | :------------ | :------- | :--------------------------------------------- |
    |  vector  | Object.Vector |          | The vector to be updated.                      |
    |  config  | Update.Config |          | The configuration of the update request.       |
    | metadata | bytes         | optional | The metadata is related to the request vector. |

  - Object.Vector

    |   field   | type   | label    | description                                     |
    | :-------: | :----- | :------- | :---------------------------------------------- |
    |    id     | string |          | The vector ID.                                  |
    |  vector   | float  | repeated | The vector.                                     |
    | timestamp | int64  |          | timestamp represents when this vector inserted. |
    | metadata  | bytes  | optional | The metadata is related to the request vector.  |

  - Update.Config

        | field | type | label | description |
        | :---: | :--- | :---- | :---------- |
        | skip_strict_exist_check | bool |  | A flag to skip exist check during update operation. |
        | filters | Filter.Config |  | Filter configuration. |
        | timestamp | int64 |  | Update timestamp. |
        | disable_balanced_update | bool |  | A flag to disable balanced update (split remove -> insert operation)

    during update operation. |

  - Filter.Config

    |  field  | type          | label    | description                                |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | description          |
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

    | field | type   | label    | description               |
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
|  5   | NOT_FOUND         |
|  6   | ALREADY_EXISTS    |
|  10  | ABORTED           |
|  13  | INTERNAL          |

Please refer to [Response Status Code](../status.md) for more details.

### Troubleshooting

The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name              | common reason                                                                                                                                       | how to resolve                                                                           |
| :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
| CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
| INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
| DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
| NOT_FOUND         | Requested ID is NOT inserted.                                                                                                                       | Send a request with an ID that is already inserted.                                      |
| ALREADY_EXISTS    | Request pair of ID and vector is already inserted.                                                                                                  | Change request ID.                                                                       |
| INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |

## StreamUpdateWithMetadata RPC

StreamUpdateWithMetadata RPC is the method to update multiple vectors and metadata using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
Using the bidirectional streaming RPC, the update request can be communicated in any order between client and server.
Each Update request and response are independent.
It's the recommended method to update the large amount of vectors.

### Input

- the scheme of `payload.v1.Update.Request`

  ```rpc
  message Update.Request {
    Object.Vector vector = 1;
    Update.Config config = 2;
    optional bytes metadata = 3;
  }

  message Object.Vector {
    string id = 1;
    repeated float vector = 2;
    int64 timestamp = 3;
    optional bytes metadata = 4;
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

  ```

  - Update.Request

    |  field   | type          | label    | description                                    |
    | :------: | :------------ | :------- | :--------------------------------------------- |
    |  vector  | Object.Vector |          | The vector to be updated.                      |
    |  config  | Update.Config |          | The configuration of the update request.       |
    | metadata | bytes         | optional | The metadata is related to the request vector. |

  - Object.Vector

    |   field   | type   | label    | description                                     |
    | :-------: | :----- | :------- | :---------------------------------------------- |
    |    id     | string |          | The vector ID.                                  |
    |  vector   | float  | repeated | The vector.                                     |
    | timestamp | int64  |          | timestamp represents when this vector inserted. |
    | metadata  | bytes  | optional | The metadata is related to the request vector.  |

  - Update.Config

        | field | type | label | description |
        | :---: | :--- | :---- | :---------- |
        | skip_strict_exist_check | bool |  | A flag to skip exist check during update operation. |
        | filters | Filter.Config |  | Filter configuration. |
        | timestamp | int64 |  | Update timestamp. |
        | disable_balanced_update | bool |  | A flag to disable balanced update (split remove -> insert operation)

    during update operation. |

  - Filter.Config

    |  field  | type          | label    | description                                |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | description          |
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

    |  field   | type              | label | description           |
    | :------: | :---------------- | :---- | :-------------------- |
    | location | Object.Location   |       | The vector location.  |
    |  status  | google.rpc.Status |       | The RPC error status. |

  - Object.Location

    | field | type   | label    | description               |
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
|  5   | NOT_FOUND         |
|  6   | ALREADY_EXISTS    |
|  10  | ABORTED           |
|  13  | INTERNAL          |

Please refer to [Response Status Code](../status.md) for more details.

### Troubleshooting

The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name              | common reason                                                                                                                                       | how to resolve                                                                           |
| :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
| CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
| INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
| DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
| NOT_FOUND         | Requested ID is NOT inserted.                                                                                                                       | Send a request with an ID that is already inserted.                                      |
| ALREADY_EXISTS    | Request pair of ID and vector is already inserted.                                                                                                  | Change request ID.                                                                       |
| INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |

## MultiUpdateWithMetadata RPC

MultiUpdateWithMetadata is the method to update multiple vectors and metadata in **1** request.

<div class="notice">
gRPC has a message size limitation.<br>
Please be careful that the size of the request exceeds the limit.
</div>

### Input

- the scheme of `payload.v1.Update.MultiRequest`

  ```rpc
  message Update.MultiRequest {
    repeated Update.Request requests = 1;
  }

  message Update.Request {
    Object.Vector vector = 1;
    Update.Config config = 2;
    optional bytes metadata = 3;
  }

  message Object.Vector {
    string id = 1;
    repeated float vector = 2;
    int64 timestamp = 3;
    optional bytes metadata = 4;
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

  ```

  - Update.MultiRequest

    |  field   | type           | label    | description                                    |
    | :------: | :------------- | :------- | :--------------------------------------------- |
    | requests | Update.Request | repeated | Represent the multiple update request content. |

  - Update.Request

    |  field   | type          | label    | description                                    |
    | :------: | :------------ | :------- | :--------------------------------------------- |
    |  vector  | Object.Vector |          | The vector to be updated.                      |
    |  config  | Update.Config |          | The configuration of the update request.       |
    | metadata | bytes         | optional | The metadata is related to the request vector. |

  - Object.Vector

    |   field   | type   | label    | description                                     |
    | :-------: | :----- | :------- | :---------------------------------------------- |
    |    id     | string |          | The vector ID.                                  |
    |  vector   | float  | repeated | The vector.                                     |
    | timestamp | int64  |          | timestamp represents when this vector inserted. |
    | metadata  | bytes  | optional | The metadata is related to the request vector.  |

  - Update.Config

        | field | type | label | description |
        | :---: | :--- | :---- | :---------- |
        | skip_strict_exist_check | bool |  | A flag to skip exist check during update operation. |
        | filters | Filter.Config |  | Filter configuration. |
        | timestamp | int64 |  | Update timestamp. |
        | disable_balanced_update | bool |  | A flag to disable balanced update (split remove -> insert operation)

    during update operation. |

  - Filter.Config

    |  field  | type          | label    | description                                |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | description          |
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

    |   field   | type            | label    | description |
    | :-------: | :-------------- | :------- | :---------- |
    | locations | Object.Location | repeated |             |

  - Object.Location

    | field | type   | label    | description               |
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
|  5   | NOT_FOUND         |
|  6   | ALREADY_EXISTS    |
|  10  | ABORTED           |
|  13  | INTERNAL          |

Please refer to [Response Status Code](../status.md) for more details.

### Troubleshooting

The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name              | common reason                                                                                                                                       | how to resolve                                                                           |
| :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
| CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
| INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
| DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
| NOT_FOUND         | Requested ID is NOT inserted.                                                                                                                       | Send a request with an ID that is already inserted.                                      |
| ALREADY_EXISTS    | Request pair of ID and vector is already inserted.                                                                                                  | Change request ID.                                                                       |
| INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |

## UpdateTimestampWithMetadata RPC

A method to update timestamp an indexed vector and metadata.

### Input

- the scheme of `payload.v1.Update.TimestampRequest`

  ```rpc
  message Update.TimestampRequest {
    string id = 1;
    int64 timestamp = 2;
    bool force = 3;
    optional bytes metadata = 4;
  }

  ```

  - Update.TimestampRequest

    |   field   | type   | label    | description                                       |
    | :-------: | :----- | :------- | :------------------------------------------------ |
    |    id     | string |          | The vector ID.                                    |
    | timestamp | int64  |          | timestamp represents when this vector inserted.   |
    |   force   | bool   |          | force represents forcefully update the timestamp. |
    | metadata  | bytes  | optional | The metadata is related to the request vector.    |

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

    | field | type   | label    | description               |
    | :---: | :----- | :------- | :------------------------ |
    | name  | string |          | The name of the location. |
    | uuid  | string |          | The UUID of the vector.   |
    |  ips  | string | repeated | The IP list.              |

### Status Code

| code | description |
| :--: | :---------- |

TODO

Please refer to [Response Status Code](../status.md) for more details.

### Troubleshooting

TODO

# Vald UpsertWithMetadata APIs

## Overview

UpsertWithMetadata Service is responsible for updating existing vectors and metadata in the `vald-agent` or inserting new vectors into the `vald-agent` and metadata if the vector does not exist.

```rpc
service UpsertWithMetadata {

  rpc UpsertWithMetadata(payload.v1.Upsert.Request) returns (payload.v1.Object.Location) {}
  rpc StreamUpsertWithMetadata(payload.v1.Upsert.Request) returns (payload.v1.Object.StreamLocation) {}
  rpc MultiUpsertWithMetadata(payload.v1.Upsert.MultiRequest) returns (payload.v1.Object.Locations) {}

}
```

## UpsertWithMetadata RPC

UpsertWithMetadata RPC is the method to update the inserted vector and metadata to a new single vector and metadata or add a new single vector and metadata if not inserted before.

### Input

- the scheme of `payload.v1.Upsert.Request`

  ```rpc
  message Upsert.Request {
    Object.Vector vector = 1;
    Upsert.Config config = 2;
    optional bytes metadata = 3;
  }

  message Object.Vector {
    string id = 1;
    repeated float vector = 2;
    int64 timestamp = 3;
    optional bytes metadata = 4;
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

  ```

  - Upsert.Request

    |  field   | type          | label    | description                                    |
    | :------: | :------------ | :------- | :--------------------------------------------- |
    |  vector  | Object.Vector |          | The vector to be upserted.                     |
    |  config  | Upsert.Config |          | The configuration of the upsert request.       |
    | metadata | bytes         | optional | The metadata is related to the request vector. |

  - Object.Vector

    |   field   | type   | label    | description                                     |
    | :-------: | :----- | :------- | :---------------------------------------------- |
    |    id     | string |          | The vector ID.                                  |
    |  vector   | float  | repeated | The vector.                                     |
    | timestamp | int64  |          | timestamp represents when this vector inserted. |
    | metadata  | bytes  | optional | The metadata is related to the request vector.  |

  - Upsert.Config

        | field | type | label | description |
        | :---: | :--- | :---- | :---------- |
        | skip_strict_exist_check | bool |  | A flag to skip exist check during upsert operation. |
        | filters | Filter.Config |  | Filter configuration. |
        | timestamp | int64 |  | Upsert timestamp. |
        | disable_balanced_update | bool |  | A flag to disable balanced update (split remove -> insert operation)

    during update operation. |

  - Filter.Config

    |  field  | type          | label    | description                                |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | description          |
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

    | field | type   | label    | description               |
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
|  5   | NOT_FOUND         |
|  6   | ALREADY_EXISTS    |
|  10  | ABORTED           |
|  13  | INTERNAL          |

Please refer to [Response Status Code](../status.md) for more details.

### Troubleshooting

The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name              | common reason                                                                                                                                       | how to resolve                                                                           |
| :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
| CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
| INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
| DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
| ALREADY_EXISTS    | Requested pair of ID and vector is already inserted                                                                                                 | Change request payload or nothing to do if update is unnecessary.                        |
| INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |

## StreamUpsertWithMetadata RPC

StreamUpsertWithMetadata RPC is the method to update multiple existing vectors and metadata or add new multiple vectors and metadata using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
Using the bidirectional streaming RPC, the upsert request can be communicated in any order between the client and server.
Each Upsert request and response are independent.
It’s the recommended method to upsert a large number of vectors.

### Input

- the scheme of `payload.v1.Upsert.Request`

  ```rpc
  message Upsert.Request {
    Object.Vector vector = 1;
    Upsert.Config config = 2;
    optional bytes metadata = 3;
  }

  message Object.Vector {
    string id = 1;
    repeated float vector = 2;
    int64 timestamp = 3;
    optional bytes metadata = 4;
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

  ```

  - Upsert.Request

    |  field   | type          | label    | description                                    |
    | :------: | :------------ | :------- | :--------------------------------------------- |
    |  vector  | Object.Vector |          | The vector to be upserted.                     |
    |  config  | Upsert.Config |          | The configuration of the upsert request.       |
    | metadata | bytes         | optional | The metadata is related to the request vector. |

  - Object.Vector

    |   field   | type   | label    | description                                     |
    | :-------: | :----- | :------- | :---------------------------------------------- |
    |    id     | string |          | The vector ID.                                  |
    |  vector   | float  | repeated | The vector.                                     |
    | timestamp | int64  |          | timestamp represents when this vector inserted. |
    | metadata  | bytes  | optional | The metadata is related to the request vector.  |

  - Upsert.Config

        | field | type | label | description |
        | :---: | :--- | :---- | :---------- |
        | skip_strict_exist_check | bool |  | A flag to skip exist check during upsert operation. |
        | filters | Filter.Config |  | Filter configuration. |
        | timestamp | int64 |  | Upsert timestamp. |
        | disable_balanced_update | bool |  | A flag to disable balanced update (split remove -> insert operation)

    during update operation. |

  - Filter.Config

    |  field  | type          | label    | description                                |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | description          |
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

    |  field   | type              | label | description           |
    | :------: | :---------------- | :---- | :-------------------- |
    | location | Object.Location   |       | The vector location.  |
    |  status  | google.rpc.Status |       | The RPC error status. |

  - Object.Location

    | field | type   | label    | description               |
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
|  5   | NOT_FOUND         |
|  6   | ALREADY_EXISTS    |
|  10  | ABORTED           |
|  13  | INTERNAL          |

Please refer to [Response Status Code](../status.md) for more details.

### Troubleshooting

The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name              | common reason                                                                                                                                       | how to resolve                                                                           |
| :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
| CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
| INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
| DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
| ALREADY_EXISTS    | Requested pair of ID and vector is already inserted                                                                                                 | Change request payload or nothing to do if update is unnecessary.                        |
| INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |

## MultiUpsertWithMetadata RPC

MultiUpsertWithMetadata is the method to update existing multiple vectors and metadata and add new multiple vectors and metadata in **1** request.

<div class="notice">
gRPC has a message size limitation.<br>
Please be careful that the size of the request exceeds the limit.
</div>

### Input

- the scheme of `payload.v1.Upsert.MultiRequest`

  ```rpc
  message Upsert.MultiRequest {
    repeated Upsert.Request requests = 1;
  }

  message Upsert.Request {
    Object.Vector vector = 1;
    Upsert.Config config = 2;
    optional bytes metadata = 3;
  }

  message Object.Vector {
    string id = 1;
    repeated float vector = 2;
    int64 timestamp = 3;
    optional bytes metadata = 4;
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

  ```

  - Upsert.MultiRequest

    |  field   | type           | label    | description                                    |
    | :------: | :------------- | :------- | :--------------------------------------------- |
    | requests | Upsert.Request | repeated | Represent the multiple upsert request content. |

  - Upsert.Request

    |  field   | type          | label    | description                                    |
    | :------: | :------------ | :------- | :--------------------------------------------- |
    |  vector  | Object.Vector |          | The vector to be upserted.                     |
    |  config  | Upsert.Config |          | The configuration of the upsert request.       |
    | metadata | bytes         | optional | The metadata is related to the request vector. |

  - Object.Vector

    |   field   | type   | label    | description                                     |
    | :-------: | :----- | :------- | :---------------------------------------------- |
    |    id     | string |          | The vector ID.                                  |
    |  vector   | float  | repeated | The vector.                                     |
    | timestamp | int64  |          | timestamp represents when this vector inserted. |
    | metadata  | bytes  | optional | The metadata is related to the request vector.  |

  - Upsert.Config

        | field | type | label | description |
        | :---: | :--- | :---- | :---------- |
        | skip_strict_exist_check | bool |  | A flag to skip exist check during upsert operation. |
        | filters | Filter.Config |  | Filter configuration. |
        | timestamp | int64 |  | Upsert timestamp. |
        | disable_balanced_update | bool |  | A flag to disable balanced update (split remove -> insert operation)

    during update operation. |

  - Filter.Config

    |  field  | type          | label    | description                                |
    | :-----: | :------------ | :------- | :----------------------------------------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type   | label | description          |
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

    |   field   | type            | label    | description |
    | :-------: | :-------------- | :------- | :---------- |
    | locations | Object.Location | repeated |             |

  - Object.Location

    | field | type   | label    | description               |
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
|  5   | NOT_FOUND         |
|  6   | ALREADY_EXISTS    |
|  10  | ABORTED           |
|  13  | INTERNAL          |

Please refer to [Response Status Code](../status.md) for more details.

### Troubleshooting

The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name              | common reason                                                                                                                                       | how to resolve                                                                           |
| :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
| CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
| INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
| DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
| ALREADY_EXISTS    | Requested pair of ID and vector is already inserted                                                                                                 | Change request payload or nothing to do if update is unnecessary.                        |
| INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |
