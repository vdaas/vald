# Vald Update APIs

## Overview

Update Service updates to new vector from inserted vector in the `vald-agent` components.

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

    | field  | type          | label | required | description                              |
    | :----: | :------------ | :---- | :------: | :--------------------------------------- |
    | vector | Object.Vector |       |    \*    | The information of vector.               |
    | config | Config        |       |    \*    | The configuration of the update request. |

  - Update.Config

    |          field          | type          | label | required | description                                                                                                   |
    | :---------------------: | :------------ | :---- | :------: | :------------------------------------------------------------------------------------------------------------ |
    | skip_strict_exist_check | bool          |       |          | Check whether the same vector is already inserted or not.<br> The ID should be unique if the value is `true`. |
    |        timestamp        | int64         |       |          | The timestamp of the vector inserted.<br>If it is N/A, the current time will be used.                         |
    |         filters         | Filter.Config |       |          | Configuration for filter.                                                                                     |
    | disable_balanced_update | bool          |       |          | A flag to disable balanced update (split remove -&gt; insert operation) during update operation.              |

  - Object.Vector

    | field  | type   | label                  | required | description                                                    |
    | :----: | :----- | :--------------------- | :------: | :------------------------------------------------------------- |
    |   id   | string |                        |    \*    | The ID of a vector. ID should consist of 1 or more characters. |
    | vector | float  | repeated(Array[float]) |    \*    | The vector data. Its dimension is between 2 and 65,536.        |

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
    | name  | string |                         | The name of vald agent pod where the request vector is updated.       |
    | uuid  | string |                         | The ID of the updated vector. It is the same as an `Object.Vector`.   |
    |  ips  | string | repeated(Array[string]) | The IP list of `vald-agent` pods where the request vector is updated. |

### Status Code

| code | name              |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  6   | ALREADY_EXISTS    |
|  10  | ABORTED           |
|  13  | INTERNAL          |

Please refer to [Response Status Code](./status.md) for more details.

### Troubleshooting

The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name              | common reason                                                                                                                                       | how to resolve                                                                           |
| :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
| CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
| INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
| DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
| NOT_FOUND         | Requested ID is NOT inserted.                                                                                                                       | Send a request with an ID that is already inserted.                                      |
| ALREADY_EXISTS    | Request pair of ID and vector is already inserted.                                                                                                  | Change request ID.                                                                       |
| INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |

## StreamUpdate RPC

StreamUpdate RPC is the method to update multiple vectors using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
Using the bidirectional streaming RPC, the update request can be communicated in any order between client and server.
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

    | field  | type          | label | required | description                              |
    | :----: | :------------ | :---- | :------: | :--------------------------------------- |
    | vector | Object.Vector |       |    \*    | The information of vector.               |
    | config | Config        |       |    \*    | The configuration of the update request. |

  - Update.Config

    |          field          | type          | label | required | description                                                                                                   |
    | :---------------------: | :------------ | :---- | :------: | :------------------------------------------------------------------------------------------------------------ |
    | skip_strict_exist_check | bool          |       |          | Check whether the same vector is already inserted or not.<br> The ID should be unique if the value is `true`. |
    |        timestamp        | int64         |       |          | The timestamp of the vector inserted.<br>If it is N/A, the current time will be used.                         |
    |         filters         | Filter.Config |       |          | Configuration for filter.                                                                                     |
    | disable_balanced_update | bool          |       |          | A flag to disable balanced update (split remove -&gt; insert operation) during update operation.              |

  - Object.Vector

    | field  | type   | label                  | required | description                                                    |
    | :----: | :----- | :--------------------- | :------: | :------------------------------------------------------------- |
    |   id   | string |                        |    \*    | The ID of a vector. ID should consist of 1 or more characters. |
    | vector | float  | repeated(Array[float]) |    \*    | The vector data. Its dimension is between 2 and 65,536.        |

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

    |  field   | type              | label | description                              |
    | :------: | :---------------- | :---- | :--------------------------------------- |
    | location | Object.Location   |       | The information of Object.Location data. |
    |  status  | google.rpc.Status |       | The status of Google RPC.                |

  - Object.Location

    | field | type   | label                   | description                                                           |
    | :---: | :----- | :---------------------- | :-------------------------------------------------------------------- |
    | name  | string |                         | The name of vald agent pod where the request vector is updated.       |
    | uuid  | string |                         | The ID of the updated vector. It is the same as an `Object.Vector`.   |
    |  ips  | string | repeated(Array[string]) | The IP list of `vald-agent` pods where the request vector is updated. |

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
|  6   | ALREADY_EXISTS    |
|  10  | ABORTED           |
|  13  | INTERNAL          |

Please refer to [Response Status Code](./status.md) for more details.

### Troubleshooting

The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name              | common reason                                                                                                                                       | how to resolve                                                                           |
| :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
| CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
| INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
| DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
| NOT_FOUND         | Requested ID is NOT inserted.                                                                                                                       | Send a request with an ID that is already inserted.                                      |
| ALREADY_EXISTS    | Request pair of ID and vector is already inserted.                                                                                                  | Change request ID.                                                                       |
| INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |

## MultiUpdate RPC

MultiUpdate is the method to update multiple vectors in **1** request.

<div class="notice">
gRPC has a message size limitation.<br>
Please be careful that the size of the request exceeds the limit.
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

    |  field   | type           | label                           | required | description       |
    | :------: | :------------- | :------------------------------ | :------: | :---------------- |
    | requests | Insert.Request | repeated(Array[Insert.Request]) |    \*    | The request list. |

  - Update.Request

    | field  | type          | label | required | description                              |
    | :----: | :------------ | :---- | :------: | :--------------------------------------- |
    | vector | Object.Vector |       |    \*    | The information of vector.               |
    | config | Config        |       |    \*    | The configuration of the update request. |

  - Update.Config

    |          field          | type          | label | required | description                                                                                                   |
    | :---------------------: | :------------ | :---- | :------: | :------------------------------------------------------------------------------------------------------------ |
    | skip_strict_exist_check | bool          |       |          | Check whether the same vector is already inserted or not.<br> The ID should be unique if the value is `true`. |
    |        timestamp        | int64         |       |          | The timestamp of the vector inserted.<br>If it is N/A, the current time will be used.                         |
    |         filters         | Filter.Config |       |          | Configuration for filter.                                                                                     |
    | disable_balanced_update | bool          |       |          | A flag to disable balanced update (split remove -&gt; insert operation) during update operation.              |

  - Object.Vector

    | field  | type   | label                  | required | description                                                    |
    | :----: | :----- | :--------------------- | :------: | :------------------------------------------------------------- |
    |   id   | string |                        |    \*    | The ID of a vector. ID should consist of 1 or more characters. |
    | vector | float  | repeated(Array[float]) |    \*    | The vector data. Its dimension is between 2 and 65,536.        |

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
    | name  | string |                         | The name of vald agent pod where the request vector is updated.       |
    | uuid  | string |                         | The ID of the updated vector. It is the same as an `Object.Vector`.   |
    |  ips  | string | repeated(Array[string]) | The IP list of `vald-agent` pods where the request vector is updated. |

### Status Code

| code | name              |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  6   | ALREADY_EXISTS    |
|  10  | ABORTED           |
|  13  | INTERNAL          |

Please refer to [Response Status Code](./status.md) for more details.

### Troubleshooting

The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name              | common reason                                                                                                                                       | how to resolve                                                                           |
| :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
| CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
| INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
| DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
| NOT_FOUND         | Requested ID is NOT inserted.                                                                                                                       | Send a request with an ID that is already inserted.                                      |
| ALREADY_EXISTS    | Request pair of ID and vector is already inserted.                                                                                                  | Change request ID.                                                                       |
| INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |
