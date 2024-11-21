# Vald Insert APIs

## Overview

Insert service provides ways to add new vectors.

```rpc
service Insert {

  rpc Insert(payload.v1.Insert.Request) returns (payload.v1.Object.Location) {}
  rpc StreamInsert(payload.v1.Insert.Request) returns (payload.v1.Object.StreamLocation) {}
  rpc MultiInsert(payload.v1.Insert.MultiRequest) returns (payload.v1.Object.Locations) {}

}
```

## Insert RPC

A method to add a new single vector.

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

    | field  | type          | label | desc.                                    |
    | :----: | :------------ | :---- | :--------------------------------------- |
    | vector | Object.Vector |       | The vector to be inserted.               |
    | config | Insert.Config |       | The configuration of the insert request. |

  - Object.Vector

    |   field   | type   | label    | desc.                                           |
    | :-------: | :----- | :------- | :---------------------------------------------- |
    |    id     | string |          | The vector ID.                                  |
    |  vector   | float  | repeated | The vector.                                     |
    | timestamp | int64  |          | timestamp represents when this vector inserted. |

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

## StreamInsert RPC

A method to add new multiple vectors by bidirectional streaming.

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

    | field  | type          | label | desc.                                    |
    | :----: | :------------ | :---- | :--------------------------------------- |
    | vector | Object.Vector |       | The vector to be inserted.               |
    | config | Insert.Config |       | The configuration of the insert request. |

  - Object.Vector

    |   field   | type   | label    | desc.                                           |
    | :-------: | :----- | :------- | :---------------------------------------------- |
    |    id     | string |          | The vector ID.                                  |
    |  vector   | float  | repeated | The vector.                                     |
    | timestamp | int64  |          | timestamp represents when this vector inserted. |

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

## MultiInsert RPC

A method to add new multiple vectors in a single request.

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

    |  field   | type           | label    | desc.                                      |
    | :------: | :------------- | :------- | :----------------------------------------- |
    | requests | Insert.Request | repeated | Represent multiple insert request content. |

  - Insert.Request

    | field  | type          | label | desc.                                    |
    | :----: | :------------ | :---- | :--------------------------------------- |
    | vector | Object.Vector |       | The vector to be inserted.               |
    | config | Insert.Config |       | The configuration of the insert request. |

  - Object.Vector

    |   field   | type   | label    | desc.                                           |
    | :-------: | :----- | :------- | :---------------------------------------------- |
    |    id     | string |          | The vector ID.                                  |
    |  vector   | float  | repeated | The vector.                                     |
    | timestamp | int64  |          | timestamp represents when this vector inserted. |

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
