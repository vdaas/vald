# Vald Remove APIs

## Overview

Remove Service is responsible for removing vectors that are indexed in the `vald-agent`.

```rpc
service Remove {

  rpc Remove(payload.v1.Remove.Request) returns (payload.v1.Object.Location) {}

  rpc StreamRemove(stream payload.v1.Remove.Request)
      returns (stream payload.v1.Object.StreamLocation) {}

  rpc MultiRemove(payload.v1.Remove.MultiRequest)
      returns (payload.v1.Object.Locations) {}
}
```

## Remove RPC

Remove RPC is the method to remove a single vector.

### Input

- the scheme of `payload.v1.Remove.Request`

  ```rpc
  message Remove {
      message Request {
          Object.ID id = 1;
          Config config = 2;
      }

      message Config {
          bool skip_strict_exist_check = 1;
          int64 timestamp = 3;
      }
  }

  message Object {
      message ID {
          string id = 1 [ (validate.rules).string.min_len = 1 ];
      }
  }
  ```

  - Remove.Request

    | field  | type      | label | required | desc.                                   |
    | :----: | :-------- | :---- | :------: | :-------------------------------------- |
    |   id   | Object.ID |       |    \*    | the id of vector                        |
    | config | Config    |       |    \*    | the configuration of the remove request |

  - Remove.Config

    |          field          | type  | label | required | desc.                                                                                                |
    | :---------------------: | :---- | :---- | :------: | :--------------------------------------------------------------------------------------------------- |
    | skip_strict_exist_check | bool  |       |          | check the same vector is already inserted or not.<br>the ID should be unique if the value is `true`. |
    |        timestamp        | int64 |       |          | the timestamp of the vector removed.<br>if it is N/A, the current time will be used.                 |

  - Object.ID

    | field | type   | label | required | desc.                                                          |
    | :---: | :----- | :---- | :------: | :------------------------------------------------------------- |
    |  id   | string |       |    \*    | the ID of a vector. ID should consist of 1 or more characters. |

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
    | name  | string |                         | the name of vald agent pod where the request vector is removed.       |
    | uuid  | string |                         | the ID of the removed vector. It is the same as an `Object.ID`.       |
    |  ips  | string | repeated(Array[string]) | the IP list of `vald-agent` pods where the request vector is removed. |

### Status Code

| code | desc.             |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  10  | ABORTED           |
|  13  | INTERNAL          |

## StreamRemove RPC

StreamRemove RPC is the method to remove multiple vectors using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
By using the bidirectional streaming RPC, the remove request can be communicated in any order between client and server.
Each Remove request and response are independent.
It's the recommended method to remove a large number of vectors.

### Input

- the scheme of `payload.v1.Remove.Request stream`

  ```rpc
  message Remove {
      message Request {
          Object.ID id = 1;
          Config config = 2;
      }

      message Config {
          bool skip_strict_exist_check = 1;
          int64 timestamp = 3;
      }
  }

  message Object {
      message ID {
          string id = 1 [ (validate.rules).string.min_len = 1 ];
      }
  }
  ```

  - Remove.Request

    | field  | type      | label | required | desc.                                   |
    | :----: | :-------- | :---- | :------: | :-------------------------------------- |
    |   id   | Object.ID |       |    \*    | the id of vector                        |
    | config | Config    |       |    \*    | the configuration of the insert request |

  - Remove.Config

    |          field          | type  | label | required | desc.                                                                                                |
    | :---------------------: | :---- | :---- | :------: | :--------------------------------------------------------------------------------------------------- |
    | skip_strict_exist_check | bool  |       |          | check the same vector is already inserted or not.<br>the ID should be unique if the value is `true`. |
    |        timestamp        | int64 |       |          | the timestamp of the vector removed.<br>if it is N/A, the current time will be used.                 |

  - Object.ID

    | field | type   | label | required | desc.                                                          |
    | :---: | :----- | :---- | :------: | :------------------------------------------------------------- |
    |  id   | string |       |    \*    | the ID of a vector. ID should consist of 1 or more characters. |

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

    |  field   | type              | label | desc.                                      |
    | :------: | :---------------- | :---- | :----------------------------------------- |
    | location | Object.Location   |       | the information of `Object.Location` data. |
    |  status  | google.rpc.Status |       | the status of google RPC                   |

  - Object.Location

    | field | type   | label                   | desc.                                                                 |
    | :---: | :----- | :---------------------- | :-------------------------------------------------------------------- |
    | name  | string |                         | the name of vald agent pod where the request vector is removed.       |
    | uuid  | string |                         | the ID of the removed vector. It is the same as an `Object.ID`.       |
    |  ips  | string | repeated(Array[string]) | the IP list of `vald-agent` pods where the request vector is removed. |

  - [google.rpc.Status](https://github.com/googleapis/googleapis/blob/master/google/rpc/status.proto)

    |  field  | type                | label                | desc.                                   |
    | :-----: | :------------------ | :------------------- | :-------------------------------------- |
    |  code   | int32               |                      | status code (code list is next section) |
    | message | string              |                      | error message                           |
    | details | google.protobuf.Any | repeated(Array[any]) | the details error message list          |

### Status Code

| code | desc.             |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  10  | ABORTED           |
|  13  | INTERNAL          |

## MultiRemove RPC

MultiRemove is the method to remove multiple vectors in **1** request.

<div class="notice">
gRPC has a message size limitation.<br>
Please be careful that the size of the request exceeds the limit.
</div>

## Input

- the scheme of `payload.v1.Remove.MultiRequest`

  ```rpc
  message Remove {
      message MultiRequest {
        repeated Request requests = 1;
      }

      message Request {
          Object.ID id = 1;
          Config config = 2;
      }

      message Config {
          bool skip_strict_exist_check = 1;
          int64 timestamp = 3;
      }
  }

  message Object {
      message ID {
          string id = 1 [ (validate.rules).string.min_len = 1 ];
      }
  }
  ```

  - Remove.MultiRequest

    |  field   | type           | label                           | required | desc.            |
    | :------: | :------------- | :------------------------------ | :------: | :--------------- |
    | requests | Remove.Request | repeated(Array[Insert.Request]) |    \*    | the request list |

  - Remove.Request

    | field  | type      | label | required | desc.                                   |
    | :----: | :-------- | :---- | :------: | :-------------------------------------- |
    |   id   | Object.ID |       |    \*    | the id of vector                        |
    | config | Config    |       |    \*    | the configuration of the remove request |

  - Remove.Config

    |          field          | type  | label | required | desc.                                                                                                |
    | :---------------------: | :---- | :---- | :------: | :--------------------------------------------------------------------------------------------------- |
    | skip_strict_exist_check | bool  |       |          | check the same vector is already inserted or not.<br>the ID should be unique if the value is `true`. |
    |        timestamp        | int64 |       |          | the timestamp of the vector removed.<br>if it is N/A, the current time will be used.                 |

  - Object.ID

    | field | type   | label | required | desc.                                                          |
    | :---: | :----- | :---- | :------: | :------------------------------------------------------------- |
    |  id   | string |       |    \*    | the ID of a vector. ID should consist of 1 or more characters. |

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
    | name  | string |                         | the name of vald agent pod where the request vector is removed.       |
    | uuid  | string |                         | the ID of the removed vector. It is the same as an `Object.Vector`.   |
    |  ips  | string | repeated(Array[string]) | the IP list of `vald-agent` pods where the request vector is removed. |

### Status Code

| code | desc.             |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  10  | ABORTED           |
|  13  | INTERNAL          |
