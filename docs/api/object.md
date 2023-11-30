# Vald Object APIs

## Overview

Object Service is responsible for getting inserted vectors and checking whether vectors are inserted into the `vald-agent` or not.

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

Exists RPC is the method to check the a vector exists in the `vald-agent`.

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

    | field | type   | label | required | desc.                                                          |
    | :---: | :----- | :---- | :------: | :------------------------------------------------------------- |
    |  id   | string |       |    \*    | the ID of a vector. ID should consist of 1 or more characters. |

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

    | field | type   | label | desc.                                                          |
    | :---: | :----- | :---- | :------------------------------------------------------------- |
    |  id   | string |       | the ID of a vector. ID should consist of 1 or more characters. |

### Status Code

| code | desc.             |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  13  | INTERNAL          |

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

    |  field  | type          | label | required | desc.                                                          |
    | :-----: | :------------ | :---- | :------: | :------------------------------------------------------------- |
    |   id    | Object.ID     |       |    \*    | the ID of a vector. ID should consist of 1 or more characters. |
    | filters | Filter.Config |       |          | configuration for filter.                                      |

  - Object.ID

    | field | type   | label | required | desc.                                                          |
    | :---: | :----- | :---- | :------: | :------------------------------------------------------------- |
    |  id   | string |       |    \*    | the ID of a vector. ID should consist of 1 or more characters. |

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

    | field  | type   | label                  | desc.                                                          |
    | :----: | :----- | :--------------------- | :------------------------------------------------------------- |
    |   id   | string |                        | the ID of a vector. ID should consist of 1 or more characters. |
    | vector | float  | repeated(Array[float]) | the vector data. its dimension is between 2 and 65,536.        |

### Status Code

| code | desc.             |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  13  | INTERNAL          |

## StreamGetObject RPC

StreamGetObject RPC is the method to get the metadata of multiple existing vectors using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
By using the bidirectional streaming RPC, the GetObject request can be communicated in any order between client and server.
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

    |  field  | type          | label | required | desc.                                                          |
    | :-----: | :------------ | :---- | :------: | :------------------------------------------------------------- |
    |   id    | Object.ID     |       |    \*    | the ID of a vector. ID should consist of 1 or more characters. |
    | filters | Filter.Config |       |          | configuration for filter.                                      |

  - Object.ID

    | field | type   | label | required | desc.                                                          |
    | :---: | :----- | :---- | :------: | :------------------------------------------------------------- |
    |  id   | string |       |    \*    | the ID of a vector. ID should consist of 1 or more characters. |

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

    | field  | type              | label | desc.                                  |
    | :----: | :---------------- | :---- | :------------------------------------- |
    | vector | Vector            |       | the information of Object.Vector data. |
    | status | google.rpc.Status |       | the status of google RPC.              |

  - Object.Vector

    | field  | type   | label                  | desc.                                                          |
    | :----: | :----- | :--------------------- | :------------------------------------------------------------- |
    |   id   | string |                        | the ID of a vector. ID should consist of 1 or more characters. |
    | vector | float  | repeated(Array[float]) | the vector data. its dimension is between 2 and 65,536.        |

### Status Code

| code | desc.             |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  13  | INTERNAL          |
