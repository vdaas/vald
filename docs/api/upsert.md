# Vald Upsert APIs

## Overview

Upsert Service is responsible for updating existing vectors in the `vald-agent` or inserting new vectors into the `vald-agent` if the vector does not exist.

```rpc
service Upsert {

  rpc Upsert(payload.v1.Upsert.Request)
      returns (payload.v1.Object.Location) {}

  rpc StreamUpsert(stream payload.v1.Upsert.Request)
      returns (stream payload.v1.Object.StreamLocation) {}

  rpc MultiUpsert(payload.v1.Upsert.MultiRequest)
      returns (payload.v1.Object.Locations) {}
}
```

## Upsert RPC

Upsert RPC is the method to update a single vector and add a new single vector.

### Input

- the scheme of `payload.v1.Upsert.Request`

  ```rpc
  message Upsert {
    message Request {
      Object.Vector vector = 1 [ (validate.rules).repeated .min_items = 2 ];
      Config config = 2;
    }

    message Config {
      bool skip_strict_exist_check = 1;
      Filter.Config filters = 2;
      int64 timestamp = 3;
    }
  }

  message Object {
      message Vector {
          string id = 1 [ (validate.rules).string.min_len = 1 ];
          repeated float vector = 2 [ (validate.rules).repeated .min_items = 2 ];
      }
  }
  ```

  - Upsert.Request

    | field  | type          | label | required | desc.                                   |
    | :----: | :------------ | :---- | :------: | :-------------------------------------- |
    | vector | Object.Vector |       |    \*    | the information of vector               |
    | config | Config        |       |    \*    | the configuration of the upsert request |

  - Upsert.Config

    |          field          | type          | label | required | desc.                                                                                                |
    | :---------------------: | :------------ | :---- | :------: | :--------------------------------------------------------------------------------------------------- |
    | skip_strict_exist_check | bool          |       |          | check the same vector is already inserted or not.<br>the ID should be unique if the value is `true`. |
    |        timestamp        | int64         |       |          | the timestamp of the vector updated/inserted.<br>if it is N/A, the current time will be used.        |
    |         filters         | Filter.Config |       |          | configuration for filter                                                                             |

  - Object.Vector

    | field  | type   | label                  | required | desc.                                                          |
    | :----: | :----- | :--------------------- | :------: | :------------------------------------------------------------- |
    |   id   | string |                        |    \*    | the ID of a vector. ID should consist of 1 or more characters. |
    | vector | float  | repeated(Array[float]) |    \*    | the vector data. its dimension is between 2 and 65,536.        |

### Output

- the scheme of `payload.v1.Object.Location`

  ```rpc
  message Object {
      message Location {
        string name = 1;
        string uuid = 2;
        repeated string ips = 3;
      }
  }
  ```

  - Object.Location

    | field | type   | label                   | desc.                                                                          |
    | :---: | :----- | :---------------------- | :----------------------------------------------------------------------------- |
    | name  | string |                         | the name of vald agent pod where the request vector is updated/inserted.       |
    | uuid  | string |                         | the ID of the updated/inserted vector. It is the same as an `Object.Vector`.   |
    |  ips  | string | repeated(Array[string]) | the IP list of `vald-agent` pods where the request vector is updated/inserted. |

### Status Code

| code | desc.            |
| :--: | :--------------- |
|  0   | OK               |
|  3   | INVALID_ARGUMENT |
|  6   | ALREADY_EXISTS   |
|  13  | INTERNAL         |

## StreamUpsert RPC

StreamUpsert RPC is the method to update multiple existing vectors or add new multiple vectors using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
By using the bidirectional streaming RPC, the upsert request can be communicated in any order between client and server.
Each Upsert request and response are independent.
It’s the recommended method to upsert a large number of vectors.

### Input

- the scheme of `payload.v1.Upsert.Request stream`

  ```rpc
  message Upsert {
      message Request {
          Object.Vector vector = 1 [ (validate.rules).repeated .min_items = 2 ];
          Config config = 2;
      }
      message Config {
          bool skip_strict_exist_check = 1;
          Filter.Config filters = 2;
          int64 timestamp = 3;
      }
  }

  message Object {
      message Vector {
          string id = 1 [ (validate.rules).string.min_len = 1 ];
          repeated float vector = 2 [ (validate.rules).repeated .min_items = 2 ];
      }
  }
  ```

  - Upsert.Request

    | field  | type          | label | required | desc.                                   |
    | :----: | :------------ | :---- | :------: | :-------------------------------------- |
    | vector | Object.Vector |       |    \*    | the information of vector               |
    | config | Config        |       |    \*    | the configuration of the upsert request |

  - Upsert.Config

    |          field          | type          | label | required | desc.                                                                                                |
    | :---------------------: | :------------ | :---- | :------: | :--------------------------------------------------------------------------------------------------- |
    | skip_strict_exist_check | bool          |       |          | check the same vector is already inserted or not.<br>the ID should be unique if the value is `true`. |
    |        timestamp        | int64         |       |          | the timestamp of the vector updated/inserted.<br>if it is N/A, the current time will be used.        |
    |         filters         | Filter.Config |       |          | configuration for filter                                                                             |

  - Object.Vector

    | field  | type   | label                  | required | desc.                                                            |
    | :----: | :----- | :--------------------- | :------: | :--------------------------------------------------------------- |
    |   id   | string |                        |    \*    | the ID of the vector. ID should consist of 1 or more characters. |
    | vector | float  | repeated(Array[float]) |    \*    | the vector data. its dimension is between 2 and 65,536.          |

### Output

- the scheme of `payload.v1.Object.StreamLocation`

  ```rpc
  message Object {
      message StreamLocation {
        oneof payload {
            Location location = 1;
            google.rpc.Status status = 2;
        }
      }

      message Location {
        string name = 1;
        string uuid = 2;
        repeated string ips = 3;
      }
  }
  ```

  - Object.StreamLocation

    |  field   | type              | label | desc.                                    |
    | :------: | :---------------- | :---- | :--------------------------------------- |
    | location | Object.Location   |       | the information of Object.Location data. |
    |  status  | google.rpc.Status |       | the status of google RPC.                |

  - Object.Location

    | field | type   | label                   | desc.                                                                          |
    | :---: | :----- | :---------------------- | :----------------------------------------------------------------------------- |
    | name  | string |                         | the name of vald agent pod where the request vector is updated/inserted.       |
    | uuid  | string |                         | the ID of the updated/inserted vector. It is the same as an `Object.Vector`.   |
    |  ips  | string | repeated(Array[string]) | the IP list of `vald-agent` pods where the request vector is updated/inserted. |

  - [google.rpc.Status](https://github.com/googleapis/googleapis/blob/master/google/rpc/status.proto)

    |  field  | type                | label                | desc.                                   |
    | :-----: | :------------------ | :------------------- | :-------------------------------------- |
    |  code   | int32               |                      | status code (code list is next section) |
    | message | string              |                      | error message                           |
    | details | google.protobuf.Any | repeated(Array[any]) | the details error message list          |

### Status Code

| code | desc.            |
| :--: | :--------------- |
|  0   | OK               |
|  3   | INVALID_ARGUMENT |
|  6   | ALREADY_EXISTS   |
|  13  | INTERNAL         |

## MultiUpsert RPC

MultiUpsert is the method to update existing multiple vectors and add new multiple vectors in **1** request.

<div class="notice">
gRPC has a message size limitation.<br>
Please be careful that the size of the request exceeds the limit.
</div>

### Input

- the scheme of `payload.v1.Upsert.MultiRequest`

  ```rpc
  message Upsert {
      message MultiRequest { repeated Request requests = 1; }

      message Request {
          Object.Vector vector = 1 [ (validate.rules).repeated .min_items = 2 ];
          Config config = 2;
      }

      message Config {
          bool skip_strict_exist_check = 1;
          Filter.Config filters = 2;
          int64 timestamp = 3;
      }
  }

  message Object {
      message Vector {
          string id = 1 [ (validate.rules).string.min_len = 1 ];
          repeated float vector = 2 [ (validate.rules).repeated .min_items = 2 ];
      }
  }
  ```

  - Upsert.MultiRequest

    |  field   | type           | label                           | required | desc.            |
    | :------: | :------------- | :------------------------------ | :------: | :--------------- |
    | requests | Upsert.Request | repeated(Array[Insert.Request]) |    \*    | the request list |

  - Upsert.Request

    | field  | type          | label | required | desc.                                   |
    | :----: | :------------ | :---- | :------: | :-------------------------------------- |
    | vector | Object.Vector |       |    \*    | the information of vector               |
    | config | Config        |       |    \*    | the configuration of the upsert request |

  - Upsert.Config

    |          field          | type          | label | required | desc.                                                                                                        |
    | :---------------------: | :------------ | :---- | :------: | :----------------------------------------------------------------------------------------------------------- |
    | skip_strict_exist_check | bool          |       |          | check the same vector is already updated/inserted or not.<br>the ID should be unique if the value is `true`. |
    |        timestamp        | int64         |       |          | the timestamp of the vector updated/inserted.<br>if it is N/A, the current time will be used.                |
    |         filters         | Filter.Config |       |          | configuration for filter                                                                                     |

  - Object.Vector

    | field  | type   | label                  | required | desc.                                                          |
    | :----: | :----- | :--------------------- | :------: | :------------------------------------------------------------- |
    |   id   | string |                        |    \*    | the ID of a vector. ID should consist of 1 or more characters. |
    | vector | float  | repeated(Array[float]) |    \*    | the vector data. its dimension is between 2 and 65,536.        |

### Output

- the scheme of `payload.v1.Object.Locations`.

  ```rpc
  message Object {
      message Locations { repeated Location locations = 1; }

      message Location {
        string name = 1;
        string uuid = 2;
        repeated string ips = 3;
      }
  }
  ```

  - Object.Locations

    |  field   | type            | label                            | desc.                         |
    | :------: | :-------------- | :------------------------------- | :---------------------------- |
    | location | Object.Location | repeated(Array[Object.Location]) | the list of `Object.Location` |

  - Object.Location

    | field | type   | label                   | desc.                                                                          |
    | :---: | :----- | :---------------------- | :----------------------------------------------------------------------------- |
    | name  | string |                         | the name of vald agent pod where the request vector is updated/inserted.       |
    | uuid  | string |                         | the ID of the updated/inserted vector. It is the same as an `Object.Vector`.   |
    |  ips  | string | repeated(Array[string]) | the IP list of `vald-agent` pods where the request vector is updated/inserted. |

### Status Code

| code | desc.            |
| :--: | :--------------- |
|  0   | OK               |
|  3   | INVALID_ARGUMENT |
|  6   | ALREADY_EXISTS   |
|  13  | INTERNAL         |
