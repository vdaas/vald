# Vald Update APIs

## Overview

Update service provides ways to update indexed vectors.

```rpc
service Update {

  rpc Update(payload.v1.Update.Request) returns (payload.v1.Object.Location) {}
  rpc StreamUpdate(payload.v1.Update.Request) returns (payload.v1.Object.StreamLocation) {}
  rpc MultiUpdate(payload.v1.Update.MultiRequest) returns (payload.v1.Object.Locations) {}
  rpc UpdateTimestamp(payload.v1.Update.TimestampRequest) returns (payload.v1.Object.Location) {}

}
```

## Update RPC

A method to update an indexed vector.

### Input

- the scheme of `payload.v1.Update.Request`

  ```rpc
  message Update.Request {
    Object.Vector vector = 1;
    Update.Config config = 2;
  }


  message Object.Vector {
    string id = 1;
    repeated float vector = 2;
    int64 timestamp = 3;
  }



  message Update.Config {
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

  - Update.Request

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | vector | Object.Vector |  | The vector to be updated. |
    | config | Update.Config |  | The configuration of the update request. |


  - Object.Vector

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | id | string |  | The vector ID. |
    | vector | float | repeated | The vector. |
    | timestamp | int64 |  | timestamp represents when this vector inserted. |



  - Update.Config

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | skip_strict_exist_check | bool |  | A flag to skip exist check during update operation. |
    | filters | Filter.Config |  | Filter configuration. |
    | timestamp | int64 |  | Update timestamp. |
    | disable_balanced_update | bool |  | A flag to disable balanced update (split remove -> insert operation)
during update operation. |


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

- the scheme of `payload.v1.Object.Location`

  ```rpc
  message Object.Location {
    string name = 1;
    string uuid = 2;
    repeated string ips = 3;
  }
  ```

  - Object.Location

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | name | string |  | The name of the location. |
    | uuid | string |  | The UUID of the vector. |
    | ips | string | repeated | The IP list. |

### Status Code

| code | desc.             |
| :--: | :---------------- |
| 0    | OK                |
| 1    | CANCELLED         |
| 3    | INVALID_ARGUMENT  |
| 4    | DEADLINE_EXCEEDED |
| 5    | NOT_FOUND         |
| 13   | INTERNAL          |

## StreamUpdate RPC

A method to update multiple indexed vectors by bidirectional streaming.

### Input

- the scheme of `payload.v1.Update.Request`

  ```rpc
  message Update.Request {
    Object.Vector vector = 1;
    Update.Config config = 2;
  }


  message Object.Vector {
    string id = 1;
    repeated float vector = 2;
    int64 timestamp = 3;
  }



  message Update.Config {
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

  - Update.Request

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | vector | Object.Vector |  | The vector to be updated. |
    | config | Update.Config |  | The configuration of the update request. |


  - Object.Vector

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | id | string |  | The vector ID. |
    | vector | float | repeated | The vector. |
    | timestamp | int64 |  | timestamp represents when this vector inserted. |



  - Update.Config

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | skip_strict_exist_check | bool |  | A flag to skip exist check during update operation. |
    | filters | Filter.Config |  | Filter configuration. |
    | timestamp | int64 |  | Update timestamp. |
    | disable_balanced_update | bool |  | A flag to disable balanced update (split remove -> insert operation)
during update operation. |


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

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | location | Object.Location |  | The vector location. |
    | status | google.rpc.Status |  | The RPC error status. |


  - Object.Location

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | name | string |  | The name of the location. |
    | uuid | string |  | The UUID of the vector. |
    | ips | string | repeated | The IP list. |

### Status Code

| code | desc.             |
| :--: | :---------------- |
| 0    | OK                |
| 1    | CANCELLED         |
| 3    | INVALID_ARGUMENT  |
| 4    | DEADLINE_EXCEEDED |
| 5    | NOT_FOUND         |
| 13   | INTERNAL          |

## MultiUpdate RPC

A method to update multiple indexed vectors in a single request.

### Input

- the scheme of `payload.v1.Update.MultiRequest`

  ```rpc
  message Update.MultiRequest {
    repeated Update.Request requests = 1;
  }


  message Update.Request {
    Object.Vector vector = 1;
    Update.Config config = 2;
  }


  message Object.Vector {
    string id = 1;
    repeated float vector = 2;
    int64 timestamp = 3;
  }



  message Update.Config {
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

  - Update.MultiRequest

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | requests | Update.Request | repeated | Represent the multiple update request content. |


  - Update.Request

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | vector | Object.Vector |  | The vector to be updated. |
    | config | Update.Config |  | The configuration of the update request. |


  - Object.Vector

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | id | string |  | The vector ID. |
    | vector | float | repeated | The vector. |
    | timestamp | int64 |  | timestamp represents when this vector inserted. |



  - Update.Config

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | skip_strict_exist_check | bool |  | A flag to skip exist check during update operation. |
    | filters | Filter.Config |  | Filter configuration. |
    | timestamp | int64 |  | Update timestamp. |
    | disable_balanced_update | bool |  | A flag to disable balanced update (split remove -> insert operation)
during update operation. |


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

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | locations | Object.Location | repeated |  |


  - Object.Location

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | name | string |  | The name of the location. |
    | uuid | string |  | The UUID of the vector. |
    | ips | string | repeated | The IP list. |

### Status Code

| code | desc.             |
| :--: | :---------------- |
| 0    | OK                |
| 1    | CANCELLED         |
| 3    | INVALID_ARGUMENT  |
| 4    | DEADLINE_EXCEEDED |
| 5    | NOT_FOUND         |
| 13   | INTERNAL          |

## UpdateTimestamp RPC

A method to update timestamp an indexed vector.

### Input

- the scheme of `payload.v1.Update.TimestampRequest`

  ```rpc
  message Update.TimestampRequest {
    string id = 1;
    int64 timestamp = 2;
    bool force = 3;
  }
  ```

  - Update.TimestampRequest

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | id | string |  | The vector ID. |
    | timestamp | int64 |  | timestamp represents when this vector inserted. |
    | force | bool |  | force represents forcefully update the timestamp. |

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

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | name | string |  | The name of the location. |
    | uuid | string |  | The UUID of the vector. |
    | ips | string | repeated | The IP list. |

### Status Code

| code | desc.             |
| :--: | :---------------- |
| 0    | OK                |
| 1    | CANCELLED         |
| 3    | INVALID_ARGUMENT  |
| 4    | DEADLINE_EXCEEDED |
| 5    | NOT_FOUND         |
| 13   | INTERNAL          |

