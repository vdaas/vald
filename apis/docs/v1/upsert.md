# Vald Upsert APIs

## Overview

Upsert Service is responsible for updating existing vectors in the `vald-agent` or inserting new vectors into the `vald-agent` if the vector does not exist.

```rpc
service Upsert {

  rpc Upsert(payload.v1.Upsert.Request) returns (payload.v1.Object.Location) {}
  rpc StreamUpsert(payload.v1.Upsert.Request) returns (payload.v1.Object.StreamLocation) {}
  rpc MultiUpsert(payload.v1.Upsert.MultiRequest) returns (payload.v1.Object.Locations) {}

}
```

## Upsert RPC

Upsert RPC is the method to update the inserted vector to a new single vector or add a new single vector if not inserted before.

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

## StreamUpsert RPC

StreamUpsert RPC is the method to update multiple existing vectors or add new multiple vectors using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
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

## MultiUpsert RPC

MultiUpsert is the method to update existing multiple vectors and add new multiple vectors in **1** request.

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
