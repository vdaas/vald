# Vald Object APIs

## Overview

Object Service is responsible for getting inserted vectors and checking whether vectors are inserted into the `vald-agent`.

```rpc
service Object {

  rpc Exists(payload.v1.Object.ID) returns (payload.v1.Object.ID) {}
  rpc GetObject(payload.v1.Object.VectorRequest) returns (payload.v1.Object.Vector) {}
  rpc StreamGetObject(payload.v1.Object.VectorRequest) returns (payload.v1.Object.StreamVector) {}
  rpc StreamListObject(payload.v1.Object.List.Request) returns (payload.v1.Object.List.Response) {}
  rpc GetTimestamp(payload.v1.Object.TimestampRequest) returns (payload.v1.Object.Timestamp) {}

}
```

## Exists RPC

Exists RPC is the method to check that a vector exists in the `vald-agent`.

### Input

- the scheme of `payload.v1.Object.ID`

  ```rpc
  message Object.ID {
    string id = 1;
  }

  ```
  - Object.ID

    | field | type | label | description |
    | :---: | :--- | :---- | :---------- |
    | id | string |  |  |


### Output

- the scheme of `payload.v1.Object.ID`

  ```rpc
  message Object.ID {
    string id = 1;
  }

  ```


  - Object.ID

    | field | type | label | description |
    | :---: | :--- | :---- | :---------- |
    | id | string |  |  |

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
## GetObject RPC

GetObject RPC is the method to get the metadata of a vector inserted into the `vald-agent`.

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

    | field | type | label | description |
    | :---: | :--- | :---- | :---------- |
    | id | Object.ID |  | The vector ID to be fetched. |
    | filters | Filter.Config |  | Filter configurations. |

  - Object.ID

    | field | type | label | description |
    | :---: | :--- | :---- | :---------- |
    | id | string |  |  |

  - Filter.Config

    | field | type | label | description |
    | :---: | :--- | :---- | :---------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type | label | description |
    | :---: | :--- | :---- | :---------- |
    | host | string |  | The target hostname. |
    | port | uint32 |  | The target port. |


### Output

- the scheme of `payload.v1.Object.Vector`

  ```rpc
  message Object.Vector {
    string id = 1;
    repeated float vector = 2;
    int64 timestamp = 3;
  }

  ```


  - Object.Vector

    | field | type | label | description |
    | :---: | :--- | :---- | :---------- |
    | id | string |  | The vector ID. |
    | vector | float | repeated | The vector. |
    | timestamp | int64 |  | timestamp represents when this vector inserted. |

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
## StreamGetObject RPC

StreamGetObject RPC is the method to get the metadata of multiple existing vectors using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
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

    | field | type | label | description |
    | :---: | :--- | :---- | :---------- |
    | id | Object.ID |  | The vector ID to be fetched. |
    | filters | Filter.Config |  | Filter configurations. |

  - Object.ID

    | field | type | label | description |
    | :---: | :--- | :---- | :---------- |
    | id | string |  |  |

  - Filter.Config

    | field | type | label | description |
    | :---: | :--- | :---- | :---------- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |

  - Filter.Target

    | field | type | label | description |
    | :---: | :--- | :---- | :---------- |
    | host | string |  | The target hostname. |
    | port | uint32 |  | The target port. |


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
  }

  ```


  - Object.StreamVector

    | field | type | label | description |
    | :---: | :--- | :---- | :---------- |
    | vector | Object.Vector |  | The vector. |
    | status | google.rpc.Status |  | The RPC error status. |

  - Object.Vector

    | field | type | label | description |
    | :---: | :--- | :---- | :---------- |
    | id | string |  | The vector ID. |
    | vector | float | repeated | The vector. |
    | timestamp | int64 |  | timestamp represents when this vector inserted. |

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
## StreamListObject RPC

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
  }

  ```


  - Object.List.Response

    | field | type | label | description |
    | :---: | :--- | :---- | :---------- |
    | vector | Object.Vector |  | The vector |
    | status | google.rpc.Status |  | The RPC error status. |

  - Object.Vector

    | field | type | label | description |
    | :---: | :--- | :---- | :---------- |
    | id | string |  | The vector ID. |
    | vector | float | repeated | The vector. |
    | timestamp | int64 |  | timestamp represents when this vector inserted. |

### Status Code

| code | description       |
| :--: | :---------------- |
TODO

Please refer to [Response Status Code](../status.md) for more details.



### Troubleshooting

TODO
## GetTimestamp RPC

Represent the RPC to get the vector metadata. This RPC is mainly used for index correction process

### Input

- the scheme of `payload.v1.Object.TimestampRequest`

  ```rpc
  message Object.TimestampRequest {
    Object.ID id = 1;
  }

  message Object.ID {
    string id = 1;
  }

  ```
  - Object.TimestampRequest

    | field | type | label | description |
    | :---: | :--- | :---- | :---------- |
    | id | Object.ID |  | The vector ID to be fetched. |

  - Object.ID

    | field | type | label | description |
    | :---: | :--- | :---- | :---------- |
    | id | string |  |  |


### Output

- the scheme of `payload.v1.Object.Timestamp`

  ```rpc
  message Object.Timestamp {
    string id = 1;
    int64 timestamp = 2;
  }

  ```


  - Object.Timestamp

    | field | type | label | description |
    | :---: | :--- | :---- | :---------- |
    | id | string |  | The vector ID. |
    | timestamp | int64 |  | timestamp represents when this vector inserted. |

### Status Code

| code | description       |
| :--: | :---------------- |
TODO

Please refer to [Response Status Code](../status.md) for more details.



### Troubleshooting

TODO
