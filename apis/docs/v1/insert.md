# Vald Insert APIs

## Overview

Insert Service is responsible for inserting new vectors into the `vald-agent`.

```rpc
service Insert {

  rpc Insert(payload.v1.Insert.Request) returns (payload.v1.Object.Location) {}
  rpc StreamInsert(payload.v1.Insert.Request) returns (payload.v1.Object.StreamLocation) {}
  rpc MultiInsert(payload.v1.Insert.MultiRequest) returns (payload.v1.Object.Locations) {}

}
```

## Insert RPC

Inset RPC is the method to add a new single vector.

### Input

- the scheme of `payload.v1.Insert.Request`

  ```rpc
  message Insert.Request {
    Object.Vector vector = 1;
    Insert.Config config = 2;
  }

  message Object.Vector {
    string id = 1;
    repeated float vector = 2;
    int64 timestamp = 3;
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

    | field  | type          | label | description                              |
    | :----: | :------------ | :---- | :--------------------------------------- |
    | vector | Object.Vector |       | The vector to be inserted.               |
    | config | Insert.Config |       | The configuration of the insert request. |

  - Object.Vector

    |   field   | type   | label    | description                                     |
    | :-------: | :----- | :------- | :---------------------------------------------- |
    |    id     | string |          | The vector ID.                                  |
    |  vector   | float  | repeated | The vector.                                     |
    | timestamp | int64  |          | timestamp represents when this vector inserted. |

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

## StreamInsert RPC

StreamInsert RPC is the method to add new multiple vectors using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
Using the bidirectional streaming RPC, the insert request can be communicated in any order between client and server.
Each Insert request and response are independent.
It's the recommended method to insert a large number of vectors.

### Input

- the scheme of `payload.v1.Insert.Request`

  ```rpc
  message Insert.Request {
    Object.Vector vector = 1;
    Insert.Config config = 2;
  }

  message Object.Vector {
    string id = 1;
    repeated float vector = 2;
    int64 timestamp = 3;
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

    | field  | type          | label | description                              |
    | :----: | :------------ | :---- | :--------------------------------------- |
    | vector | Object.Vector |       | The vector to be inserted.               |
    | config | Insert.Config |       | The configuration of the insert request. |

  - Object.Vector

    |   field   | type   | label    | description                                     |
    | :-------: | :----- | :------- | :---------------------------------------------- |
    |    id     | string |          | The vector ID.                                  |
    |  vector   | float  | repeated | The vector.                                     |
    | timestamp | int64  |          | timestamp represents when this vector inserted. |

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

## MultiInsert RPC

MultiInsert RPC is the method to add multiple new vectors in **1** request.

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
  }

  message Object.Vector {
    string id = 1;
    repeated float vector = 2;
    int64 timestamp = 3;
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

    | field  | type          | label | description                              |
    | :----: | :------------ | :---- | :--------------------------------------- |
    | vector | Object.Vector |       | The vector to be inserted.               |
    | config | Insert.Config |       | The configuration of the insert request. |

  - Object.Vector

    |   field   | type   | label    | description                                     |
    | :-------: | :----- | :------- | :---------------------------------------------- |
    |    id     | string |          | The vector ID.                                  |
    |  vector   | float  | repeated | The vector.                                     |
    | timestamp | int64  |          | timestamp represents when this vector inserted. |

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
