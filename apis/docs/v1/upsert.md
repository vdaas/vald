# Vald Upsert APIs

## Overview```rpc

service Upsert {

rpc Upsert(payload.v1.Upsert.Request) returns (payload.v1.Object.Location) {}
rpc StreamUpsert(payload.v1.Upsert.Request) returns (payload.v1.Object.StreamLocation) {}
rpc MultiUpsert(payload.v1.Upsert.MultiRequest) returns (payload.v1.Object.Locations) {}

}

````
## Upsert RPC

### Input

- the scheme of `payload.v1.Upsert.Request`

  ```rpc
  message Upsert.Request {
    Object.Vector vector = 1;
    Upsert.Config config = 2;
  }


  message Object.Vector {
    string id = 1;
    repeated float vector = 2;
    int64 timestamp = 3;
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
````

- Upsert.Request

  | field  | type          | label | desc.                                    |
  | :----: | :------------ | :---- | :--------------------------------------- |
  | vector | Object.Vector |       | The vector to be upserted.               |
  | config | Upsert.Config |       | The configuration of the upsert request. |

- Object.Vector

  |   field   | type   | label    | desc.                                           |
  | :-------: | :----- | :------- | :---------------------------------------------- |
  |    id     | string |          | The vector ID.                                  |
  |  vector   | float  | repeated | The vector.                                     |
  | timestamp | int64  |          | timestamp represents when this vector inserted. |

- Upsert.Config

      | field | type | label | desc. |
      | :---: | :--- | :---- | :---- |
      | skip_strict_exist_check | bool |  | A flag to skip exist check during upsert operation. |
      | filters | Filter.Config |  | Filter configuration. |
      | timestamp | int64 |  | Upsert timestamp. |
      | disable_balanced_update | bool |  | A flag to disable balanced update (split remove -> insert operation)

  during update operation. |

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

## StreamUpsert RPC

### Input

- the scheme of `payload.v1.Upsert.Request`

  ```rpc
  message Upsert.Request {
    Object.Vector vector = 1;
    Upsert.Config config = 2;
  }


  message Object.Vector {
    string id = 1;
    repeated float vector = 2;
    int64 timestamp = 3;
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

    | field  | type          | label | desc.                                    |
    | :----: | :------------ | :---- | :--------------------------------------- |
    | vector | Object.Vector |       | The vector to be upserted.               |
    | config | Upsert.Config |       | The configuration of the upsert request. |

  - Object.Vector

    |   field   | type   | label    | desc.                                           |
    | :-------: | :----- | :------- | :---------------------------------------------- |
    |    id     | string |          | The vector ID.                                  |
    |  vector   | float  | repeated | The vector.                                     |
    | timestamp | int64  |          | timestamp represents when this vector inserted. |

  - Upsert.Config

        | field | type | label | desc. |
        | :---: | :--- | :---- | :---- |
        | skip_strict_exist_check | bool |  | A flag to skip exist check during upsert operation. |
        | filters | Filter.Config |  | Filter configuration. |
        | timestamp | int64 |  | Upsert timestamp. |
        | disable_balanced_update | bool |  | A flag to disable balanced update (split remove -> insert operation)

    during update operation. |

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

## MultiUpsert RPC

### Input

- the scheme of `payload.v1.Upsert.MultiRequest`

  ```rpc
  message Upsert.MultiRequest {
    repeated Upsert.Request requests = 1;
  }


  message Upsert.Request {
    Object.Vector vector = 1;
    Upsert.Config config = 2;
  }


  message Object.Vector {
    string id = 1;
    repeated float vector = 2;
    int64 timestamp = 3;
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

    |  field   | type           | label    | desc.                                          |
    | :------: | :------------- | :------- | :--------------------------------------------- |
    | requests | Upsert.Request | repeated | Represent the multiple upsert request content. |

  - Upsert.Request

    | field  | type          | label | desc.                                    |
    | :----: | :------------ | :---- | :--------------------------------------- |
    | vector | Object.Vector |       | The vector to be upserted.               |
    | config | Upsert.Config |       | The configuration of the upsert request. |

  - Object.Vector

    |   field   | type   | label    | desc.                                           |
    | :-------: | :----- | :------- | :---------------------------------------------- |
    |    id     | string |          | The vector ID.                                  |
    |  vector   | float  | repeated | The vector.                                     |
    | timestamp | int64  |          | timestamp represents when this vector inserted. |

  - Upsert.Config

        | field | type | label | desc. |
        | :---: | :--- | :---- | :---- |
        | skip_strict_exist_check | bool |  | A flag to skip exist check during upsert operation. |
        | filters | Filter.Config |  | Filter configuration. |
        | timestamp | int64 |  | Upsert timestamp. |
        | disable_balanced_update | bool |  | A flag to disable balanced update (split remove -> insert operation)

    during update operation. |

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
