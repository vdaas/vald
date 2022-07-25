# Vald Update APIs

## Overview

Update Service is responsible for updating vectors which are already inserted in the `vald-agent` components.

```rpc
service Update {
    rpc Update(payload.v1.Update.Request) returns (payload.v1.Object.Location) {}

    rpc StreamUpdate(stream payload.v1.Update.Request) returns (stream payload.v1.Object.Location) {}

    rpc MultiUpdate(payload.v1.Update.MultiRequest) returns (stream payload.v1.Object.Locations) {}
}
```

## Update RPC

Update RPC is the method to update a single vector.

### Input

- the scheme of `payload.v1.Update.Request`

  ```rpc
  message Update {
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
          repeated float vector = 2 [ (validate.rules).repeated.min_items = 2 ];
      }
  }
  ```

  - Update.Request

    | field  | type          | label | required | desc.                                   |
    | :----: | :------------ | :---- | :------: | :-------------------------------------- |
    | vector | Object.Vector |       |    \*    | the information of vector               |
    | config | Config        |       |    \*    | the configuration of the update request |

  - Update.Config

    |          field          | type          | label | required | desc.                                                                                                |
    | :---------------------: | :------------ | :---- | :------: | :--------------------------------------------------------------------------------------------------- |
    | skip_strict_exist_check | bool          |       |          | check the same vector is already inserted or not.<br>the ID should be unique if the value is `true`. |
    |        timestamp        | int64         |       |          | the timestamp of the vector updated.<br>if it is N/A, the current time will be used.                 |
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

    | field | type   | label                   | desc.                                                                 |
    | :---: | :----- | :---------------------- | :-------------------------------------------------------------------- |
    | name  | string |                         | the name of vald agent pod where the request vector is updated.       |
    | uuid  | string |                         | the ID of the updated vector. It is the same as an `Object.Vector`.   |
    |  ips  | string | repeated(Array[string]) | the IP list of `vald-agent` pods where the request vector is updated. |

### Status Code

| code | desc.            |
| :--: | :--------------- |
|  0   | OK               |
|  3   | INVALID_ARGUMENT |
|  5   | NOT_FOUND        |
|  6   | ALREADY_EXISTS   |
|  13  | INTERNAL         |

## StreamUpdate RPC

StreamUpdate RPC is the method to update multiple vectors using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
By using the bidirectional streaming RPC, the update request can be communicated in any order between client and server.
Each Update request and response are independent.
It's the recommended method to update the large amount of vectors.

### Input

- the scheme of `payload.v1.Update.Request stream`

  ```rpc
  message Update {
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

  - Update.Request

    | field  | type          | label | required | desc.                                   |
    | :----: | :------------ | :---- | :------: | :-------------------------------------- |
    | vector | Object.Vector |       |    \*    | the information of vector               |
    | config | Config        |       |    \*    | the configuration of the update request |

  - Update.Config

    |          field          | type          | label | required | desc.                                                                                                |
    | :---------------------: | :------------ | :---- | :------: | :--------------------------------------------------------------------------------------------------- |
    | skip_strict_exist_check | bool          |       |          | check the same vector is already inserted or not.<br>the ID should be unique if the value is `true`. |
    |        timestamp        | int64         |       |          | the timestamp of the vector updated.<br>if it is N/A, the current time will be used.                 |
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

    | field | type   | label                   | desc.                                                                  |
    | :---: | :----- | :---------------------- | :--------------------------------------------------------------------- |
    | name  | string |                         | the name of vald agent pod where the request vector is updated.        |
    | uuid  | string |                         | the ID of the exists vector. It is the same as an `Object.Vector`.     |
    |  ips  | string | repeated(Array[string]) | the IP list of `vald-agent` pods where the request vector is inserted. |

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
|  5   | NOT_FOUND        |
|  6   | ALREADY_EXISTS   |
|  13  | INTERNAL         |

## MultiUpdate RPC

MultiUpdate is the method to update multiple vectors in **1** request.

<div class="card-note">
gRPC has the message size limitation.<br>
Please be careful that the size of the request exceed the limit.
</div>

### Input

- the scheme of `payload.v1.Update.MultiRequest`

  ```rpc
  message Update {
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

  - Update.MultiRequest

    |  field   | type           | label                           | required | desc.            |
    | :------: | :------------- | :------------------------------ | :------: | :--------------- |
    | requests | Insert.Request | repeated(Array[Insert.Request]) |    \*    | the request list |

  - Update.Request

    | field  | type          | label | required | desc.                                   |
    | :----: | :------------ | :---- | :------: | :-------------------------------------- |
    | vector | Object.Vector |       |    \*    | the information of vector               |
    | config | Config        |       |    \*    | the configuration of the update request |

  - Update.Config

    |          field          | type          | label | required | desc.                                                                                                |
    | :---------------------: | :------------ | :---- | :------: | :--------------------------------------------------------------------------------------------------- |
    | skip_strict_exist_check | bool          |       |          | check the same vector is already inserted or not.<br>the ID should be unique if the value is `true`. |
    |        timestamp        | int64         |       |          | the timestamp of the vector updated.<br>if it is N/A, the current time will be used.                 |
    |         filters         | Filter.Config |       |          | configuration for filter                                                                             |

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

    |  field   | type            | label                            | desc.                          |
    | :------: | :-------------- | :------------------------------- | :----------------------------- |
    | location | Object.Location | repeated(Array[Object.Location]) | the list of `Object.Location`. |

  - Object.Location

    | field | type   | label                   | desc.                                                                 |
    | :---: | :----- | :---------------------- | :-------------------------------------------------------------------- |
    | name  | string |                         | the name of vald agent pod where the request vector is updated.       |
    | uuid  | string |                         | the ID of the updated vector. It is the same as an `Object.Vector`.   |
    |  ips  | string | repeated(Array[string]) | the ip list of `vald-agent` pods where the request vector is updated. |

### Status Code

| code | desc.            |
| :--: | :--------------- |
|  0   | OK               |
|  3   | INVALID_ARGUMENT |
|  5   | NOT_FOUND        |
|  6   | ALREADY_EXISTS   |
|  13  | INTERNAL         |
