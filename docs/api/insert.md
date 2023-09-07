# Vald Insert APIs

## Overview

Insert Service is responsible for inserting new vectors into the `vald-agent`.

```rpc
service Insert {
    rpc Insert(payload.v1.Insert.Request) returns (payload.v1.Object.Location) {}

    rpc StreamInsert(stream payload.v1.Insert.Request) returns (stream payload.v1.Object.Location) {}

    rpc MultiInsert(payload.v1.Insert.MultiRequest) returns (payload.v1.Object.Locations) {}
}
```

## Insert RPC

Inset RPC is the method to add a new single vector.

### Input

- the scheme of `payload.v1.Insert.Request`

  ```rpc
  message Insert {
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

  - Insert.Request

    | field  | type          | label | required | description                             |
    | :----: | :------------ | :---- | :------: | :-------------------------------------- |
    | vector | Object.Vector |       |    \*    | the information of vector               |
    | config | Config        |       |    \*    | the configuration of the insert request |

  - Insert.Config

    |          field          | type          | label | required | description                                                                                          |
    | :---------------------: | :------------ | :---- | :------: | :--------------------------------------------------------------------------------------------------- |
    | skip_strict_exist_check | bool          |       |          | check the same vector is already inserted or not.<br>the ID should be unique if the value is `true`. |
    |        timestamp        | int64         |       |          | the timestamp of the vector inserted.<br>if it is N/A, the current time will be used.                |
    |         filters         | Filter.Config |       |          | configuration for filter                                                                             |

  - Object.Vector

    | field  | type   | label                  | required | description                                                    |
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

    | field | type   | label                   | description                                                            |
    | :---: | :----- | :---------------------- | :--------------------------------------------------------------------- |
    | name  | string |                         | the name of vald agent pod where the request vector is inserted.       |
    | uuid  | string |                         | the ID of the inserted vector. It is the same as an `Object.Vector`.   |
    |  ips  | string | repeated(Array[string]) | the IP list of `vald-agent` pods where the request vector is inserted. |

### Status Code

| code | name              |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  6   | ALREADY_EXISTS    |
|  13  | INTERNAL          |

For more details, please refer to [Response Status Code](./status.md).

### Troubleshooting

When the response code is NOT `0 (OK)`, the request process may not be completed.

Here are some common reason and how to resolve of each error.

| name              | common reason                                                                                            | how to resolve                                            |
| :---------------- | :------------------------------------------------------------------------------------------------------- | :-------------------------------------------------------- |
| CANCELLED         | executed cancel() of rpc from client/server side or something network problems between client and server | Check the code, especially around timeout and connection management and fix if needed.                 |
| INVALID_ARGUMENT  | dimension of request vector is NOT same as Vald Agent's config or requested vector's ID is empty or some request payload is invalid.           | check Agent config and request payload and fix request payload or Agent config |
| DEADLINE_EXCEEDED | RPC timeout setting is too short on client or server side                                                | check timeout setting and fix if needed                   |
| ALREADY_EXISTS    | request ID is already inserted                                                                           | change request ID                                         |
| INTERNAL          | target Vald cluster or network route has some critical error                                                              | check target Vald cluster at first and check network route including ingress as second                       |

## StreamInsert RPC

StreamInsert RPC is the method to add new multiple vectors using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
By using the bidirectional streaming RPC, the insert request can be communicated in any order between client and server.
Each Insert request and response are independent.
It's the recommended method to insert a large number of vectors.

### Input

- the scheme of `payload.v1.Insert.Request stream`

  ```rpc
  message Insert {
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

  - Insert.Request

    | field  | type          | label | required | description                             |
    | :----: | :------------ | :---- | :------: | :-------------------------------------- |
    | vector | Object.Vector |       |    \*    | the information of vector               |
    | config | Config        |       |    \*    | the configuration of the insert request |

  - Insert.Config

    |          field          | type          | label | required | description                                                                                          |
    | :---------------------: | :------------ | :---- | :------: | :--------------------------------------------------------------------------------------------------- |
    | skip_strict_exist_check | bool          |       |          | check the same vector is already inserted or not.<br>the ID should be unique if the value is `true`. |
    |        timestamp        | int64         |       |          | the timestamp of the vector inserted.<br>if it is N/A, the current time will be used.                |
    |         filters         | Filter.Config |       |          | configuration for filter                                                                             |

  - Object.Vector

    | field  | type   | label                  | required | description                                                      |
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

    |  field   | type              | label | description                                |
    | :------: | :---------------- | :---- | :----------------------------------------- |
    | location | Object.Location   |       | the information of `Object.Location` data. |
    |  status  | google.rpc.Status |       | the status of google RPC.                  |

  - Object.Location

    | field | type   | label                   | description                                                            |
    | :---: | :----- | :---------------------- | :--------------------------------------------------------------------- |
    | name  | string |                         | the name of vald agent pod where the request vector is inserted.       |
    | uuid  | string |                         | the ID of the inserted vector. It is the same as an `Object.Vector`.   |
    |  ips  | string | repeated(Array[string]) | the IP list of `vald-agent` pods where the request vector is inserted. |

  - [google.rpc.Status](https://github.com/googleapis/googleapis/blob/master/google/rpc/status.proto)

    |  field  | type                | label                | description                             |
    | :-----: | :------------------ | :------------------- | :-------------------------------------- |
    |  code   | int32               |                      | status code (code list is next section) |
    | message | string              |                      | error message                           |
    | details | google.protobuf.Any | repeated(Array[any]) | the details error message list          |

### Status Code

| code | name              |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  6   | ALREADY_EXISTS    |
|  13  | INTERNAL          |

For more details, please refer to [Response Status Code](./status.md).

### Troubleshooting

When the response code is NOT `0 (OK)`, the request process may not be completed.

Here are some common reason and how to resolve of each error.

| name              | common reason                                                                                            | how to resolve                                                                           |
| :---------------- | :------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
| CANCELLED         | executed cancel() of rpc from client/server side or something network problems between client and server | verify connection between client and server side<BR>check client code and fix if needed. |
| INVALID_ARGUMENT  | dimension of request vector is NOT same as Vald Agent's config or requested vector's ID is ""            | check request config and Agent config, and fix if needed                                 |
| DEADLINE_EXCEEDED | RPC timeout setting is too short on client or server side                                                | check gRPC timeout setting both of client side and server side and fix if needed         |
| ALREADY_EXISTS    | request ID is already inserted                                                                           | change request ID then retry insert                                                      |
| INTERNAL          | target Vald cluster has some critical error                                                              | check target Vald cluster at first                                                       |

## MultiInsert RPC

MultiInsert RPC is the method to add multiple new vectors in **1** request.

<div class="notice">
gRPC has a message size limitation.<br>
Please be careful that the size of the request exceeds the limit.
</div>

### Input

- the scheme of `payload.v1.Insert.MultiRequest`

  ```rpc
  message Insert {
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

  - Insert.MultiRequest

    |  field   | type           | label                           | required | description      |
    | :------: | :------------- | :------------------------------ | :------: | :--------------- |
    | requests | Insert.Request | repeated(Array[Insert.Request]) |    \*    | the request list |

  - Insert.Request

    | field  | type          | label | required | description                             |
    | :----: | :------------ | :---- | :------: | :-------------------------------------- |
    | vector | Object.Vector |       |    \*    | the information of vector               |
    | config | Config        |       |    \*    | the configuration of the insert request |

  - Insert.Config

    |          field          | type          | label | required | description                                                                                          |
    | :---------------------: | :------------ | :---- | :------: | :--------------------------------------------------------------------------------------------------- |
    | skip_strict_exist_check | bool          |       |          | check the same vector is already inserted or not.<br>the ID should be unique if the value is `true`. |
    |        timestamp        | int64         |       |          | the timestamp of the vector inserted.<br>if it is N/A, the current time will be used.                |
    |         filters         | Filter.Config |       |          | configuration for filter                                                                             |

  - Object.Vector

    | field  | type   | label                  | required | description                                                    |
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

    |  field   | type            | label                            | description                    |
    | :------: | :-------------- | :------------------------------- | :----------------------------- |
    | location | Object.Location | repeated(Array[Object.Location]) | the list of `Object.Location`. |

  - Object.Location

    | field | type   | label                   | description                                                            |
    | :---: | :----- | :---------------------- | :--------------------------------------------------------------------- |
    | name  | string |                         | the name of vald agent pod where the request vector is inserted.       |
    | uuid  | string |                         | the ID of the inserted vector. It is the same as an `Object.Vector`.   |
    |  ips  | string | repeated(Array[string]) | the IP list of `vald-agent` pods where the request vector is inserted. |

### Status Code

| code | name              |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  6   | ALREADY_EXISTS    |
|  13  | INTERNAL          |

For more details, please refer to [Response Status Code](./status.md).

### Troubleshooting

When the response code is NOT `0 (OK)`, the request process may not be completed.

Here are some common reason and how to resolve of each error.

| name              | common reason                                                                                            | how to resolve                                                                           |
| :---------------- | :------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
| CANCELLED         | executed cancel() of rpc from client/server side or something network problems between client and server | verify connection between client and server side<BR>check client code and fix if needed. |
| INVALID_ARGUMENT  | dimension of request vector is NOT same as Vald Agent's config or requested vector's ID is ""            | check request config and Agent config, and fix if needed                                 |
| DEADLINE_EXCEEDED | RPC timeout setting is too short on client or server side                                                | check gRPC timeout setting both of client side and server side and fix if needed         |
| ALREADY_EXISTS    | request ID is already inserted                                                                           | change request ID then retry insert                                                      |
| INTERNAL          | target Vald cluster has some critical error                                                              | check target Vald cluster at first                                                       |
