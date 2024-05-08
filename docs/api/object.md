# Vald Object APIs

## Overview

Object Service is responsible for getting inserted vectors and checking whether vectors are inserted into the `vald-agent`.

```rpc
service Object {
  rpc Exists(payload.v1.Object.ID) returns (payload.v1.Object.ID) {}

  rpc GetObject(payload.v1.Object.VectorRequest)
      returns (payload.v1.Object.Vector) {}

  rpc StreamGetObject(stream payload.v1.Object.VectorRequest)
      returns (stream payload.v1.Object.StreamVector) {}
}
```

## Exists RPC

Exists RPC is the method to check that a vector exists in the `vald-agent`.

### Input

- the scheme of `payload.v1.Object.ID`

  ```rpc
  message Object {
      message ID {
          string id = 1 [ (validate.rules).string.min_len = 1 ];
      }
  }
  ```

  - Object.ID

    | field | type   | label | required | description                                                    |
    | :---: | :----- | :---- | :------: | :------------------------------------------------------------- |
    |  id   | string |       |    \*    | The ID of a vector. ID should consist of 1 or more characters. |

### Output

- the scheme of `payload.v1.Object.ID`

  ```rpc
  message Object {
      message ID {
          string id = 1 [ (validate.rules).string.min_len = 1 ];
      }
  }
  ```

  - Object.ID

    | field | type   | label | description                                                    |
    | :---: | :----- | :---- | :------------------------------------------------------------- |
    |  id   | string |       | The ID of a vector. ID should consist of 1 or more characters. |

### Status Code

| code | name              |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  13  | INTERNAL          |

Please refer to [Response Status Code](./status.md) for more details.

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
  message Object {
      message VectorRequest {
        ID id = 1 [ (validate.rules).repeated .min_items = 2 ];
        Filter.Config filters = 2;
      }

      message ID {
          string id = 1 [ (validate.rules).string.min_len = 1 ];
      }
  }
  ```

  - Object.VectorRequest

    |  field  | type          | label | required | description                                                    |
    | :-----: | :------------ | :---- | :------: | :------------------------------------------------------------- |
    |   id    | Object.ID     |       |    \*    | The ID of a vector. ID should consist of 1 or more characters. |
    | filters | Filter.Config |       |          | Configuration for filter.                                      |

  - Object.ID

    | field | type   | label | required | description                                                    |
    | :---: | :----- | :---- | :------: | :------------------------------------------------------------- |
    |  id   | string |       |    \*    | The ID of a vector. ID should consist of 1 or more characters. |

### Output

- the scheme of `payload.v1.Object.Vector`

  ```rpc
  message Object {
      message Vector {
          string id = 1 [ (validate.rules).string.min_len = 1 ];
          repeated float vector = 2 [ (validate.rules).repeated .min_items = 2 ];
      }
  }
  ```

  - Object.Vector

    | field  | type   | label                  | description                                                    |
    | :----: | :----- | :--------------------- | :------------------------------------------------------------- |
    |   id   | string |                        | The ID of a vector. ID should consist of 1 or more characters. |
    | vector | float  | repeated(Array[float]) | The vector data. Its dimension is between 2 and 65,536.        |

### Status Code

| code | name              |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  13  | INTERNAL          |

Please refer to [Response Status Code](./status.md) for more details.

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

- the scheme of `payload.v1.Object.VectorRequest stream`

  ```rpc
  message Object {
      message VectorRequest {
        ID id = 1 [ (validate.rules).repeated .min_items = 2 ];
        Filter.Config filters = 2;
      }

      message ID {
          string id = 1 [ (validate.rules).string.min_len = 1 ];
      }
  }
  ```

  - Object.VectorRequest

    |  field  | type          | label | required | description                                                    |
    | :-----: | :------------ | :---- | :------: | :------------------------------------------------------------- |
    |   id    | Object.ID     |       |    \*    | The ID of a vector. ID should consist of 1 or more characters. |
    | filters | Filter.Config |       |          | Configuration for the filter targets.                          |

  - Object.ID

    | field | type   | label | required | description                                                    |
    | :---: | :----- | :---- | :------: | :------------------------------------------------------------- |
    |  id   | string |       |    \*    | The ID of a vector. ID should consist of 1 or more characters. |

### Output

- the scheme of `payload.v1.Object.StreamVector`

  ```rpc
  message Object {
      message StreamVector {
        oneof payload {
            Vector vector = 1;
            google.rpc.Status status = 2;
        }
      }
      message Vector {
          string id = 1 [ (validate.rules).string.min_len = 1 ];
          repeated float vector = 2 [ (validate.rules).repeated .min_items = 2 ];
      }
  }
  ```

  - Object.StreamVector

    | field  | type              | label | description                            |
    | :----: | :---------------- | :---- | :------------------------------------- |
    | vector | Vector            |       | The information of Object.Vector data. |
    | status | google.rpc.Status |       | The status of Google RPC.              |

  - Object.Vector

    | field  | type   | label                  | description                                                    |
    | :----: | :----- | :--------------------- | :------------------------------------------------------------- |
    |   id   | string |                        | The ID of a vector. ID should consist of 1 or more characters. |
    | vector | float  | repeated(Array[float]) | The vector data. Its dimension is between 2 and 65,536.        |

### Status Code

| code | name              |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  13  | INTERNAL          |

Please refer to [Response Status Code](./status.md) for more details.

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
