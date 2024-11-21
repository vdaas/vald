# Vald Remove APIs

## Overview

Remove service provides ways to remove indexed vectors.

```rpc
service Remove {

  rpc Remove(payload.v1.Remove.Request) returns (payload.v1.Object.Location) {}
  rpc RemoveByTimestamp(payload.v1.Remove.TimestampRequest) returns (payload.v1.Object.Locations) {}
  rpc StreamRemove(payload.v1.Remove.Request) returns (payload.v1.Object.StreamLocation) {}
  rpc MultiRemove(payload.v1.Remove.MultiRequest) returns (payload.v1.Object.Locations) {}

}
```

## Remove RPC

A method to remove an indexed vector.

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

    | field  | type          | label | desc.                                    |
    | :----: | :------------ | :---- | :--------------------------------------- |
    |   id   | Object.ID     |       | The object ID to be removed.             |
    | config | Remove.Config |       | The configuration of the remove request. |

  - Object.ID

    | field | type   | label | desc. |
    | :---: | :----- | :---- | :---- |
    |  id   | string |       |       |

  - Remove.Config

    |          field          | type  | label | desc.                                               |
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

    | field | type   | label    | desc.                     |
    | :---: | :----- | :------- | :------------------------ |
    | name  | string |          | The name of the location. |
    | uuid  | string |          | The UUID of the vector.   |
    |  ips  | string | repeated | The IP list.              |

### Status Code

| code | desc.             |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  13  | INTERNAL          |

## RemoveByTimestamp RPC

A method to remove an indexed vector based on timestamp.

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
  ```

  - Remove.TimestampRequest

        | field | type | label | desc. |
        | :---: | :--- | :---- | :---- |
        | timestamps | Remove.Timestamp | repeated | The timestamp comparison list. If more than one is specified, the `AND`

    search is applied. |

  - Remove.Timestamp

    |   field   | type                      | label | desc.                     |
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

| code | desc.             |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  13  | INTERNAL          |

## StreamRemove RPC

A method to remove multiple indexed vectors by bidirectional streaming.

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

    | field  | type          | label | desc.                                    |
    | :----: | :------------ | :---- | :--------------------------------------- |
    |   id   | Object.ID     |       | The object ID to be removed.             |
    | config | Remove.Config |       | The configuration of the remove request. |

  - Object.ID

    | field | type   | label | desc. |
    | :---: | :----- | :---- | :---- |
    |  id   | string |       |       |

  - Remove.Config

    |          field          | type  | label | desc.                                               |
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

| code | desc.             |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  13  | INTERNAL          |

## MultiRemove RPC

A method to remove multiple indexed vectors in a single request.

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

    |  field   | type           | label    | desc.                                          |
    | :------: | :------------- | :------- | :--------------------------------------------- |
    | requests | Remove.Request | repeated | Represent the multiple remove request content. |

  - Remove.Request

    | field  | type          | label | desc.                                    |
    | :----: | :------------ | :---- | :--------------------------------------- |
    |   id   | Object.ID     |       | The object ID to be removed.             |
    | config | Remove.Config |       | The configuration of the remove request. |

  - Object.ID

    | field | type   | label | desc. |
    | :---: | :----- | :---- | :---- |
    |  id   | string |       |       |

  - Remove.Config

    |          field          | type  | label | desc.                                               |
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

| code | desc.             |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  13  | INTERNAL          |
