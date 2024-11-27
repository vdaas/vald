# Vald Remove APIs

## Overview

Remove Service is responsible for removing vectors indexed in the `vald-agent`.

```rpc
service Remove {

  rpc Remove(payload.v1.Remove.Request) returns (payload.v1.Object.Location) {}
  rpc RemoveByTimestamp(payload.v1.Remove.TimestampRequest) returns (payload.v1.Object.Locations) {}
  rpc StreamRemove(payload.v1.Remove.Request) returns (payload.v1.Object.StreamLocation) {}
  rpc MultiRemove(payload.v1.Remove.MultiRequest) returns (payload.v1.Object.Locations) {}

}
```

## Remove RPC

Remove RPC is the method to remove a single vector.

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

## RemoveByTimestamp RPC

RemoveByTimestamp RPC is the method to remove vectors based on timestamp.

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

  enum  Remove.Timestamp.Operator {
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

## StreamRemove RPC

A method to remove multiple indexed vectors by bidirectional streaming.

StreamRemove RPC is the method to remove multiple vectors using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
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

## MultiRemove RPC

MultiRemove is the method to remove multiple vectors in **1** request.

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
