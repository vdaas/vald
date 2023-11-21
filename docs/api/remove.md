# Vald Remove APIs

## Overview

Remove Service is responsible for removing vectors indexed in the `vald-agent`.

```rpc
service Remove {

  rpc Remove(payload.v1.Remove.Request) returns (payload.v1.Object.Location) {}

  rpc RemoveByTimestamp(payload.v1.Remove.TimestampRequest) returns (payload.v1.Object.Locations) {}

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

    | field  | type      | label | required | description                              |
    | :----: | :-------- | :---- | :------: | :--------------------------------------- |
    |   id   | Object.ID |       |    \*    | The ID of vector.                        |
    | config | Config    |       |    \*    | The configuration of the remove request. |

  - Remove.Config

    |          field          | type  | label | required | description                                                                                                  |
    | :---------------------: | :---- | :---- | :------: | :----------------------------------------------------------------------------------------------------------- |
    | skip_strict_exist_check | bool  |       |          | Check whether the same vector is already inserted or not.<br>The ID should be unique if the value is `true`. |
    |        timestamp        | int64 |       |          | The timestamp of the vector removed.<br>If it is N/A, the current time will be used.                         |

  - Object.ID

    | field | type   | label | required | description                                                    |
    | :---: | :----- | :---- | :------: | :------------------------------------------------------------- |
    |  id   | string |       |    \*    | The ID of a vector. ID should consist of 1 or more characters. |

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

    | field | type   | label                   | description                                                           |
    | :---: | :----- | :---------------------- | :-------------------------------------------------------------------- |
    | name  | string |                         | The name of vald agent pod where the request vector is removed.       |
    | uuid  | string |                         | The ID of the removed vector. It is the same as an `Object.ID`.       |
    |  ips  | string | repeated(Array[string]) | The IP list of `vald-agent` pods where the request vector is removed. |

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

## RemoveByTimestamp RPC

RemoveByTimestamp RPC is the method to remove vectors based on timestamp.

### Input

- the scheme of `payload.v1.Remove.TimestampRequest`

  ```rpc
  message Remove {
      message TimestampRequest {
          repeated Timestamp timestamps = 1;
      }

      message Timestamp {
          enum Operator {
              Eq = 0;
              Ne = 1;
              Ge = 2;
              Gt = 3;
              Le = 4;
              Lt = 5;
          }
          int64 timestamp = 1;
          Operator operator = 2;
      }
  }

  message Object {
      message ID {
          string id = 1 [ (validate.rules).string.min_len = 1 ];
      }
  }
  ```

  - Remove.TimestampRequest

    | field      | type             | label                             | required | description                                                                                   |
    | :--------: | :--------------- | :-------------------------------- | :------: | :-------------------------------------------------------------------------------------------- |
    | timestamps | Remove.Timestamp | repeated(Array[Remove.Timestamp]) |    \*    | The timestamp comparison list.<br>If more than one is specified, the `AND` search is applied. |

  - Remove.Timestamp

    | field     | type                      | label | required | description                                        |
    | :-------: | :------------------------ | :---- | :------: | :------------------------------------------------- |
    | timestamp | int64                     |       |    \*    | The timestamp.               |
    | operator  | Remove.Timestamp.Operator |       |          | The conditionl operator. (default value is `Eq`). |

  - Remove.Timestamp.Operator

    | value | description            |
    | :---: | :--------------------- |
    |  Eq   | Equal.                 |
    |  Ne   | Not Equal.             |
    |  Ge   | Greater than or Equal. |
    |  Gt   | Greater than.          |
    |  Le   | Less than or Equal.    |
    |  Lt   | Less than.             |

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

    |  field   | type            | label                            | description                    |
    | :------: | :-------------- | :------------------------------- | :----------------------------- |
    | location | Object.Location | repeated(Array[Object.Location]) | The list of `Object.Location`. |

  - Object.Location

    | field | type   | label                   | description                                                           |
    | :---: | :----- | :---------------------- | :-------------------------------------------------------------------- |
    | name  | string |                         | The name of vald agent pod where the request vector is removed.       |
    | uuid  | string |                         | The ID of the removed vector. It is the same as an `Object.ID`.       |
    |  ips  | string | repeated(Array[string]) | The IP list of `vald-agent` pods where the request vector is removed. |

### Status Code

| code | name              |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  13  | INTERNAL          |

Please refer to [Response Status Code](./status.md) for more details.

### Troubleshooting

The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name              | common reason                                                                                   | how to resolve                                                                                                       |
| :---------------- | :---------------------------------------------------------------------------------------------- | :------------------------------------------------------------------------------------------------------------------- |
| CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed.                              |
| DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed.                             |
| NOT_FOUND         | No vectors in the system match the specified timestamp conditions.                              | Check whether vectors matching the specified timestamp conditions exist in the system, and fix conditions if needed. |
| INTERNAL          | Target Vald cluster or network route has some critical error.                                   | Check target Vald cluster first and check network route including ingress as second.                                 |

## StreamRemove RPC

StreamRemove RPC is the method to remove multiple vectors using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
Using the bidirectional streaming RPC, the remove request can be communicated in any order between client and server.
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

    | field  | type      | label | required | description                              |
    | :----: | :-------- | :---- | :------: | :--------------------------------------- |
    |   id   | Object.ID |       |    \*    | The ID of vector.                        |
    | config | Config    |       |    \*    | The configuration of the insert request. |

  - Remove.Config

    |          field          | type  | label | required | description                                                                                                  |
    | :---------------------: | :---- | :---- | :------: | :----------------------------------------------------------------------------------------------------------- |
    | skip_strict_exist_check | bool  |       |          | Check whether the same vector is already inserted or not.<br>The ID should be unique if the value is `true`. |
    |        timestamp        | int64 |       |          | The timestamp of the vector removed.<br>If it is N/A, the current time will be used.                         |

  - Object.ID

    | field | type   | label | required | description                                                    |
    | :---: | :----- | :---- | :------: | :------------------------------------------------------------- |
    |  id   | string |       |    \*    | The ID of a vector. ID should consist of 1 or more characters. |

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

    |  field   | type              | label | description                                |
    | :------: | :---------------- | :---- | :----------------------------------------- |
    | location | Object.Location   |       | The information of `Object.Location` data. |
    |  status  | google.rpc.Status |       | The status of Google RPC                   |

  - Object.Location

    | field | type   | label                   | description                                                           |
    | :---: | :----- | :---------------------- | :-------------------------------------------------------------------- |
    | name  | string |                         | The name of vald agent pod where the request vector is removed.       |
    | uuid  | string |                         | The ID of the removed vector. It is the same as an `Object.ID`.       |
    |  ips  | string | repeated(Array[string]) | The IP list of `vald-agent` pods where the request vector is removed. |

  - [google.rpc.Status](https://github.com/googleapis/googleapis/blob/master/google/rpc/status.proto)

    |  field  | type                | label                | description                             |
    | :-----: | :------------------ | :------------------- | :-------------------------------------- |
    |  code   | int32               |                      | Status code (code list is next section) |
    | message | string              |                      | Error message                           |
    | details | google.protobuf.Any | repeated(Array[any]) | The details error message list          |

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

    |  field   | type           | label                           | required | description      |
    | :------: | :------------- | :------------------------------ | :------: | :--------------- |
    | requests | Remove.Request | repeated(Array[Insert.Request]) |    \*    | the request list |

  - Remove.Request

    | field  | type      | label | required | description                              |
    | :----: | :-------- | :---- | :------: | :--------------------------------------- |
    |   id   | Object.ID |       |    \*    | The ID of vector.                        |
    | config | Config    |       |    \*    | The configuration of the remove request. |

  - Remove.Config

    |          field          | type  | label | required | description                                                                                                  |
    | :---------------------: | :---- | :---- | :------: | :----------------------------------------------------------------------------------------------------------- |
    | skip_strict_exist_check | bool  |       |          | Check whether the same vector is already inserted or not.<br>The ID should be unique if the value is `true`. |
    |        timestamp        | int64 |       |          | The timestamp of the vector removed.<br>If it is N/A, the current time will be used.                         |

  - Object.ID

    | field | type   | label | required | description                                                    |
    | :---: | :----- | :---- | :------: | :------------------------------------------------------------- |
    |  id   | string |       |    \*    | The ID of a vector. ID should consist of 1 or more characters. |

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

    |  field   | type            | label                            | description                    |
    | :------: | :-------------- | :------------------------------- | :----------------------------- |
    | location | Object.Location | repeated(Array[Object.Location]) | The list of `Object.Location`. |

  - Object.Location

    | field | type   | label                   | description                                                           |
    | :---: | :----- | :---------------------- | :-------------------------------------------------------------------- |
    | name  | string |                         | The name of vald agent pod where the request vector is removed.       |
    | uuid  | string |                         | The ID of the removed vector. It is the same as an `Object.ID`.       |
    |  ips  | string | repeated(Array[string]) | The IP list of `vald-agent` pods where the request vector is removed. |

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
