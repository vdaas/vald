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

Upsert RPC is the method to update the inserted vector to a new single vector or add a new single vector if not inserted before.

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

    | field  | type          | label | required | description                              |
    | :----: | :------------ | :---- | :------: | :--------------------------------------- |
    | vector | Object.Vector |       |    \*    | The information of vector.               |
    | config | Config        |       |    \*    | The configuration of the update request. |

  - Upsert.Config

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

    | field | type   | label                   | description                                                                    |
    | :---: | :----- | :---------------------- | :----------------------------------------------------------------------------- |
    | name  | string |                         | The name of vald agent pod where the request vector is updated/inserted.       |
    | uuid  | string |                         | The ID of the updated/inserted vector. It is the same as an `Object.Vector`.   |
    |  ips  | string | repeated(Array[string]) | The IP list of `vald-agent` pods where the request vector is updated/inserted. |

### Status Code

| code | name              |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
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
| ALREADY_EXISTS    | Requested pair of ID and vector is already inserted                                                                                                 | Change request payload or nothing to do if update is unnecessary.                        |
| INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |

## StreamUpsert RPC

StreamUpsert RPC is the method to update multiple existing vectors or add new multiple vectors using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
Using the bidirectional streaming RPC, the upsert request can be communicated in any order between the client and server.
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

    | field  | type          | label | required | description                              |
    | :----: | :------------ | :---- | :------: | :--------------------------------------- |
    | vector | Object.Vector |       |    \*    | The information of vector.               |
    | config | Config        |       |    \*    | The configuration of the update request. |

  - Upsert.Config

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

    | field | type   | label                   | description                                                                    |
    | :---: | :----- | :---------------------- | :----------------------------------------------------------------------------- |
    | name  | string |                         | The name of vald agent pod where the request vector is updated/inserted.       |
    | uuid  | string |                         | The ID of the updated/inserted vector. It is the same as an `Object.Vector`.   |
    |  ips  | string | repeated(Array[string]) | The IP list of `vald-agent` pods where the request vector is updated/inserted. |

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
| ALREADY_EXISTS    | Requested pair of ID and vector is already inserted                                                                                                 | Change request payload or nothing to do if update is unnecessary.                        |
| INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |

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

    |  field   | type           | label                           | required | description       |
    | :------: | :------------- | :------------------------------ | :------: | :---------------- |
    | requests | Upsert.Request | repeated(Array[Insert.Request]) |    \*    | The request list. |

  - Upsert.Request

    | field  | type          | label | required | description                              |
    | :----: | :------------ | :---- | :------: | :--------------------------------------- |
    | vector | Object.Vector |       |    \*    | The information of vector.               |
    | config | Config        |       |    \*    | The configuration of the update request. |

  - Upsert.Config

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

    | field | type   | label                   | description                                                                    |
    | :---: | :----- | :---------------------- | :----------------------------------------------------------------------------- |
    | name  | string |                         | The name of vald agent pod where the request vector is updated/inserted.       |
    | uuid  | string |                         | The ID of the updated/inserted vector. It is the same as an `Object.Vector`.   |
    |  ips  | string | repeated(Array[string]) | The IP list of `vald-agent` pods where the request vector is updated/inserted. |

### Status Code

| code | name              |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
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
| ALREADY_EXISTS    | Requested pair of ID and vector is already inserted                                                                                                 | Change request payload or nothing to do if update is unnecessary.                        |
| INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |
