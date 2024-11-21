# Vald Object APIs

## Overview

Object service provides ways to fetch indexed vectors.

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

A method to check whether a specified ID is indexed or not.

### Input

- the scheme of `payload.v1.Object.ID`

  ```rpc
  message Object.ID {
    string id = 1;
  }
  ```

  - Object.ID

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | id | string |  |  |

### Output

- the scheme of `payload.v1.Object.ID`

  ```rpc
  message Object.ID {
    string id = 1;
  }
  ```

  - Object.ID

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | id | string |  |  |

### Status Code

| code | desc.             |
| :--: | :---------------- |
| 0    | OK                |
| 1    | CANCELLED         |
| 3    | INVALID_ARGUMENT  |
| 4    | DEADLINE_EXCEEDED |
| 5    | NOT_FOUND         |
| 13   | INTERNAL          |

## GetObject RPC

A method to fetch a vector.

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

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | id | Object.ID |  | The vector ID to be fetched. |
    | filters | Filter.Config |  | Filter configurations. |


  - Object.ID

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | id | string |  |  |



  - Filter.Config

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |


  - Filter.Target

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
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

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | id | string |  | The vector ID. |
    | vector | float | repeated | The vector. |
    | timestamp | int64 |  | timestamp represents when this vector inserted. |

### Status Code

| code | desc.             |
| :--: | :---------------- |
| 0    | OK                |
| 1    | CANCELLED         |
| 3    | INVALID_ARGUMENT  |
| 4    | DEADLINE_EXCEEDED |
| 5    | NOT_FOUND         |
| 13   | INTERNAL          |

## StreamGetObject RPC

A method to fetch vectors by bidirectional streaming.

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

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | id | Object.ID |  | The vector ID to be fetched. |
    | filters | Filter.Config |  | Filter configurations. |


  - Object.ID

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | id | string |  |  |



  - Filter.Config

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | targets | Filter.Target | repeated | Represent the filter target configuration. |


  - Filter.Target

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
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

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | vector | Object.Vector |  | The vector. |
    | status | google.rpc.Status |  | The RPC error status. |


  - Object.Vector

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | id | string |  | The vector ID. |
    | vector | float | repeated | The vector. |
    | timestamp | int64 |  | timestamp represents when this vector inserted. |

### Status Code

| code | desc.             |
| :--: | :---------------- |
| 0    | OK                |
| 1    | CANCELLED         |
| 3    | INVALID_ARGUMENT  |
| 4    | DEADLINE_EXCEEDED |
| 5    | NOT_FOUND         |
| 13   | INTERNAL          |

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

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | vector | Object.Vector |  | The vector |
    | status | google.rpc.Status |  | The RPC error status. |


  - Object.Vector

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | id | string |  | The vector ID. |
    | vector | float | repeated | The vector. |
    | timestamp | int64 |  | timestamp represents when this vector inserted. |

### Status Code

| code | desc.             |
| :--: | :---------------- |
| 0    | OK                |
| 1    | CANCELLED         |
| 3    | INVALID_ARGUMENT  |
| 4    | DEADLINE_EXCEEDED |
| 5    | NOT_FOUND         |
| 13   | INTERNAL          |

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

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | id | Object.ID |  | The vector ID to be fetched. |


  - Object.ID

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
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

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | id | string |  | The vector ID. |
    | timestamp | int64 |  | timestamp represents when this vector inserted. |

### Status Code

| code | desc.             |
| :--: | :---------------- |
| 0    | OK                |
| 1    | CANCELLED         |
| 3    | INVALID_ARGUMENT  |
| 4    | DEADLINE_EXCEEDED |
| 5    | NOT_FOUND         |
| 13   | INTERNAL          |

