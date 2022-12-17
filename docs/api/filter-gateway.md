# Vald Filter Gateway APIs

## Overview

Filter Servie is responsible for providing insert, update, upsert and search interface for `Vald Filter Gateway`.

Vald Filter Gateway forward user request to user-defined ingress/egress filter components allowing user to run custom logic.

## Insert RPC

Insert RPC is the method to insert object through Vald Filter Gateway.

```rpc
service Filter {
  rpc InsertObject(payload.v1.Insert.ObjectRequest)
      returns (payload.v1.Object.Location) {
    option (google.api.http) = {
      post : "/insert/object"
      body : "*"
    };
  }
}
```

### Input

- the scheme of `payload.v1.Insert.ObjectRequest`

  ```rpc
  message Insert {
      message ObjectRequest {
        Object.Blob object = 1;
        Config config = 2;
        Filter.Target vectorizer = 3;
      }

      message Config {
        bool skip_strict_exist_check = 1;
        Filter.Config filters = 2;
        int64 timestamp = 3;
      }
  }

  message Object {
      message Blob {
        string id = 1 [ (validate.rules).string.min_len = 1 ];
        bytes object = 2;
      }
  }

  message Filter {
      message Target {
        string host = 1;
        uint32 port = 2;
      }

      message Config {
        repeated Target targets = 1;
      }
  }
  ```

  - Insert.ObjectRequest

    |   field    | type          | label | required | desc.                                   |
    | :--------: | :------------ | :---- | :------: | :-------------------------------------- |
    |   object   | Object.Blob   |       |    \*    | the binary object to be inserted        |
    |   config   | Config        |       |    \*    | the configuration of the insert request |
    | vectorizer | Filter.Target |       |    \*    | filter target                           |

  - Object.Blob

    | field  | type   | label | required | desc.             |
    | :----: | :----- | :---- | :------: | :---------------- |
    |   id   | string |       |    \*    | the object ID     |
    | object | bytes  |       |    \*    | the binary object |

  - Insert.Config

    |          field          | type          | label | required | desc.                                                                                                |
    | :---------------------: | :------------ | :---- | :------: | :--------------------------------------------------------------------------------------------------- |
    | skip_strict_exist_check | bool          |       |          | check the same vector is already inserted or not.<br>the ID should be unique if the value is `true`. |
    |        timestamp        | int64         |       |          | the timestamp of the vector inserted.<br>if it is N/A, the current time will be used.                |
    |         filters         | Filter.Config |       |          | configuration for filter                                                                             |

  - Filter.Target

    | field | type   | label | required | desc.               |
    | :---: | :----- | :---- | :------: | :------------------ |
    | host  | string |       |    \*    | the target hostname |
    | port  | port   |       |    \*    | the target port     |

  - Filter.Config

    |  field  | type          | label                          | required | desc.                           |
    | :-----: | :------------ | :----------------------------- | :------: | :------------------------------ |
    | targets | Filter.Target | repeated(Array[Filter.Target]) |    \*    | the filter target configuration |

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

    | field | type   | label                   | desc.                                                                  |
    | :---: | :----- | :---------------------- | :--------------------------------------------------------------------- |
    | name  | string |                         | the name of vald agent pod where the request vector is inserted.       |
    | uuid  | string |                         | the ID of the inserted vector. It is the same as an `Object.Vector`.   |
    |  ips  | string | repeated(Array[string]) | the IP list of `vald-agent` pods where the request vector is inserted. |

### Status Code

| code | desc.            |
| :--: | :--------------- |
|  0   | OK               |
|  3   | INVALID_ARGUMENT |
|  6   | ALREADY_EXISTS   |
|  13  | INTERNAL         |

## StreamInsert RPC

StreamInsert RPC is the method to add new multiple object using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).

By using the bidirectional streaming RPC, the insert request can be communicated in any order between client and server.
Each Insert request and response are independent.
It's the recommended method to insert a large number of objects.

```rpc
service Filter {
  rpc StreamInsertObject(stream payload.v1.Insert.ObjectRequest)
      returns (stream payload.v1.Object.StreamLocation) {}
}
```

### Input

- the scheme of `payload.v1.Insert.ObjectRequest`

  ```rpc
  message Insert {
      message ObjectRequest {
        Object.Blob object = 1;
        Config config = 2;
        Filter.Target vectorizer = 3;
      }

      message Config {
        bool skip_strict_exist_check = 1;
        Filter.Config filters = 2;
        int64 timestamp = 3;
      }
  }

  message Object {
      message Blob {
        string id = 1 [ (validate.rules).string.min_len = 1 ];
        bytes object = 2;
      }
  }

  message Filter {
      message Target {
        string host = 1;
        uint32 port = 2;
      }

      message Config {
        repeated Target targets = 1;
      }
  }
  ```

  - Insert.ObjectRequest

    |   field    | type          | label | required | desc.                                   |
    | :--------: | :------------ | :---- | :------: | :-------------------------------------- |
    |   object   | Object.Blob   |       |    \*    | the binary object to be inserted        |
    |   config   | Config        |       |    \*    | the configuration of the insert request |
    | vectorizer | Filter.Target |       |    \*    | filter configurations                   |

  - Object.Blob

    | field  | type   | label | required | desc.             |
    | :----: | :----- | :---- | :------: | :---------------- |
    |   id   | string |       |    \*    | the object ID     |
    | object | bytes  |       |    \*    | the binary object |

  - Insert.Config

    |          field          | type          | label | required | desc.                                                                                                |
    | :---------------------: | :------------ | :---- | :------: | :--------------------------------------------------------------------------------------------------- |
    | skip_strict_exist_check | bool          |       |          | check the same vector is already inserted or not.<br>the ID should be unique if the value is `true`. |
    |        timestamp        | int64         |       |          | the timestamp of the vector inserted.<br>if it is N/A, the current time will be used.                |
    |         filters         | Filter.Config |       |          | configuration for filter                                                                             |

  - Filter.Target

    | field | type   | label | required | desc.               |
    | :---: | :----- | :---- | :------: | :------------------ |
    | host  | string |       |    \*    | the target hostname |
    | port  | port   |       |    \*    | the target port     |

  - Filter.Config

    |  field  | type          | label                          | required | desc.                           |
    | :-----: | :------------ | :----------------------------- | :------: | :------------------------------ |
    | targets | Filter.Target | repeated(Array[Filter.Target]) |    \*    | the filter target configuration |

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
    |  status  | google.rpc.Status |       | the status of google RPC.                  |

  - Object.Location

    | field | type   | label                   | desc.                                                                  |
    | :---: | :----- | :---------------------- | :--------------------------------------------------------------------- |
    | name  | string |                         | the name of vald agent pod where the request vector is inserted.       |
    | uuid  | string |                         | the ID of the inserted vector. It is the same as an `Object.Vector`.   |
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
|  6   | ALREADY_EXISTS   |
|  13  | INTERNAL         |

## MultiInsert RPC

MultiInsert RPC is the method to add multiple new objects in **1** request.

```rpc
service Filter {
  rpc MultiInsertObject(payload.v1.Insert.MultiObjectRequest)
      returns (payload.v1.Object.Locations) {
    option (google.api.http) = {
      post : "/insert/object/multiple"
      body : "*"
    };
  }
}
```

### Input

- the scheme of `payload.v1.Insert.MultiObjectRequest`

  ```rpc
  message Insert {
      message MultiObjectRequest {
        repeated ObjectRequest requests = 1;
      }

      message ObjectRequest {
        Object.Blob object = 1;
        Config config = 2;
        Filter.Target vectorizer = 3;
      }

      message Config {
        bool skip_strict_exist_check = 1;
        Filter.Config filters = 2;
        int64 timestamp = 3;
      }
  }

  message Object {
      message Blob {
        string id = 1 [ (validate.rules).string.min_len = 1 ];
        bytes object = 2;
      }
  }

  message Filter {
      message Target {
        string host = 1;
        uint32 port = 2;
      }

      message Config {
        repeated Target targets = 1;
      }
  }
  ```

  - Insert.MultiObjectRequest

    |  field   | type          | label                                 | required | desc.                                                |
    | :------: | :------------ | :------------------------------------ | :------: | :--------------------------------------------------- |
    | requests | ObjectRequest | repeated(Array[Insert.ObjectRequest]) |    \*    | the multiple search by binary object request content |

  - Insert.ObjectRequest

    |   field    | type          | label | required | desc.                                   |
    | :--------: | :------------ | :---- | :------: | :-------------------------------------- |
    |   object   | Object.Blob   |       |    \*    | the binary object to be inserted        |
    |   config   | Config        |       |    \*    | the configuration of the insert request |
    | vectorizer | Filter.Target |       |    \*    | filter configurations                   |

  - Object.Blob

    | field  | type   | label | required | desc.             |
    | :----: | :----- | :---- | :------: | :---------------- |
    |   id   | string |       |    \*    | the object ID     |
    | object | bytes  |       |    \*    | the binary object |

  - Insert.Config

    |          field          | type          | label | required | desc.                                                                                                |
    | :---------------------: | :------------ | :---- | :------: | :--------------------------------------------------------------------------------------------------- |
    | skip_strict_exist_check | bool          |       |          | check the same vector is already inserted or not.<br>the ID should be unique if the value is `true`. |
    |        timestamp        | int64         |       |          | the timestamp of the vector inserted.<br>if it is N/A, the current time will be used.                |
    |         filters         | Filter.Config |       |          | configuration for filter                                                                             |

  - Filter.Target

    | field | type   | label | required | desc.               |
    | :---: | :----- | :---- | :------: | :------------------ |
    | host  | string |       |    \*    | the target hostname |
    | port  | port   |       |    \*    | the target port     |

  - Filter.Config

    |  field  | type          | label                          | required | desc.                           |
    | :-----: | :------------ | :----------------------------- | :------: | :------------------------------ |
    | targets | Filter.Target | repeated(Array[Filter.Target]) |    \*    | the filter target configuration |

### Output

- the scheme of `payload.v1.Object.Locations`

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

    | field | type   | label                   | desc.                                                                  |
    | :---: | :----- | :---------------------- | :--------------------------------------------------------------------- |
    | name  | string |                         | the name of vald agent pod where the request vector is inserted.       |
    | uuid  | string |                         | the ID of the inserted vector. It is the same as an `Object.Vector`.   |
    |  ips  | string | repeated(Array[string]) | the IP list of `vald-agent` pods where the request vector is inserted. |

### Status Code

| code | desc.            |
| :--: | :--------------- |
|  0   | OK               |
|  3   | INVALID_ARGUMENT |
|  6   | ALREADY_EXISTS   |
|  13  | INTERNAL         |

## Update RPC

Update RPC is the method to update a single vector.

```rpc
service Filter {
  rpc UpdateObject(payload.v1.Update.ObjectRequest)
      returns (payload.v1.Object.Location) {
    option (google.api.http) = {
      post : "/update/object"
      body : "*"
    };
  }
}
```

### Input

- the scheme of `payload.v1.Update.ObjectRequest`

  ```rpc
  message Update {
      message ObjectRequest {
        Object.Blob object = 1;
        Config config = 2;
        Filter.Target vectorizer = 3;
      }

      message Config {
        bool skip_strict_exist_check = 1;
        Filter.Config filters = 2;
        int64 timestamp = 3;
      }
  }

  message Object {
      message Blob {
        string id = 1 [ (validate.rules).string.min_len = 1 ];
        bytes object = 2;
      }
  }

  message Filter {
      message Target {
        string host = 1;
        uint32 port = 2;
      }

      message Config {
        repeated Target targets = 1;
      }
  }
  ```

  - Update.ObjectRequest

    |   field    | type          | label | required | desc.                                   |
    | :--------: | :------------ | :---- | :------: | :-------------------------------------- |
    |   object   | Object.Blob   |       |    \*    | the binary object to be updated         |
    |   config   | Config        |       |    \*    | the configuration of the update request |
    | vectorizer | Filter.Target |       |    \*    | filter target                           |

  - Object.Blob

    | field  | type   | label | required | desc.             |
    | :----: | :----- | :---- | :------: | :---------------- |
    |   id   | string |       |    \*    | the object ID     |
    | object | bytes  |       |    \*    | the binary object |

  - Update.Config

    |          field          | type          | label | required | desc.                                                                                               |
    | :---------------------: | :------------ | :---- | :------: | :-------------------------------------------------------------------------------------------------- |
    | skip_strict_exist_check | bool          |       |          | check the same vector is already updated or not.<br>the ID should be unique if the value is `true`. |
    |        timestamp        | int64         |       |          | the timestamp of the vector updated.<br>if it is N/A, the current time will be used.                |
    |         filters         | Filter.Config |       |          | configuration for filter                                                                            |

  - Filter.Target

    | field | type   | label | required | desc.               |
    | :---: | :----- | :---- | :------: | :------------------ |
    | host  | string |       |    \*    | the target hostname |
    | port  | port   |       |    \*    | the target port     |

  - Filter.Config

    |  field  | type          | label                          | required | desc.                           |
    | :-----: | :------------ | :----------------------------- | :------: | :------------------------------ |
    | targets | Filter.Target | repeated(Array[Filter.Target]) |    \*    | the filter target configuration |

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
|  6   | ALREADY_EXISTS   |
|  13  | INTERNAL         |

## Stream Update RPC

StreamUpdate RPC is the method to update multiple objects using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
By using the bidirectional streaming RPC, the update request can be communicated in any order between client and server.
Each Update request and response are independent.
It's the recommended method to update the large amount of objects.

```rpc
service Filter {
  rpc StreamUpdateObject(stream payload.v1.Update.ObjectRequest)
      returns (stream payload.v1.Object.StreamLocation) {}
}
```

### Input

- the scheme of `payload.v1.Update.ObjectRequest stream`

  ```rpc
  message Update {
      message ObjectRequest {
        Object.Blob object = 1;
        Config config = 2;
        Filter.Target vectorizer = 3;
      }

      message Config {
        bool skip_strict_exist_check = 1;
        Filter.Config filters = 2;
        int64 timestamp = 3;
      }
  }

  message Object {
      message Blob {
        string id = 1 [ (validate.rules).string.min_len = 1 ];
        bytes object = 2;
      }
  }

  message Filter {
      message Target {
        string host = 1;
        uint32 port = 2;
      }

      message Config {
        repeated Target targets = 1;
      }
  }
  ```

  - Update.ObjectRequest

    |   field    | type          | label | required | desc.                                   |
    | :--------: | :------------ | :---- | :------: | :-------------------------------------- |
    |   object   | Object.Blob   |       |    \*    | the binary object to be updated         |
    |   config   | Config        |       |    \*    | the configuration of the update request |
    | vectorizer | Filter.Target |       |    \*    | filter target                           |

  - Object.Blob

    | field  | type   | label | required | desc.             |
    | :----: | :----- | :---- | :------: | :---------------- |
    |   id   | string |       |    \*    | the object ID     |
    | object | bytes  |       |    \*    | the binary object |

  - Update.Config

    |          field          | type          | label | required | desc.                                                                                               |
    | :---------------------: | :------------ | :---- | :------: | :-------------------------------------------------------------------------------------------------- |
    | skip_strict_exist_check | bool          |       |          | check the same vector is already updated or not.<br>the ID should be unique if the value is `true`. |
    |        timestamp        | int64         |       |          | the timestamp of the vector updated.<br>if it is N/A, the current time will be used.                |
    |         filters         | Filter.Config |       |          | configuration for filter                                                                            |

  - Filter.Target

    | field | type   | label | required | desc.               |
    | :---: | :----- | :---- | :------: | :------------------ |
    | host  | string |       |    \*    | the target hostname |
    | port  | port   |       |    \*    | the target port     |

  - Filter.Config

    |  field  | type          | label                          | required | desc.                           |
    | :-----: | :------------ | :----------------------------- | :------: | :------------------------------ |
    | targets | Filter.Target | repeated(Array[Filter.Target]) |    \*    | the filter target configuration |

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
    |  status  | google.rpc.Status |       | the status of google RPC.                  |

  - Object.Location

    | field | type   | label                   | desc.                                                                 |
    | :---: | :----- | :---------------------- | :-------------------------------------------------------------------- |
    | name  | string |                         | the name of vald agent pod where the request vector is updated .      |
    | uuid  | string |                         | the ID of the updated vector. It is the same as an `Object.Vector`.   |
    |  ips  | string | repeated(Array[string]) | the IP list of `vald-agent` pods where the request vector is updated. |

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

## MultiUpdate RPC

MultiUpdate is the method to update multiple objects in **1** request.

<div class="notice">
gRPC has the message size limitation.<br>
Please be careful that the size of the request exceed the limit.
</div>

```rpc
service Filter {
  rpc MultiUpdateObject(payload.v1.Update.MultiObjectRequest)
      returns (payload.v1.Object.Locations) {
    option (google.api.http) = {
      post : "/update/object/multiple"
      body : "*"
    };
  }
}
```

### Input

- the scheme of `payload.v1.Update.MultiObjectRequest`

  ```rpc
  message Update {
      message MultiObjectRequest {
        repeated ObjectRequest requests = 1;
      }

      message ObjectRequest {
        Object.Blob object = 1;
        Config config = 2;
        Filter.Target vectorizer = 3;
      }

      message Config {
        bool skip_strict_exist_check = 1;
        Filter.Config filters = 2;
        int64 timestamp = 3;
      }
  }

  message Object {
      message Blob {
        string id = 1 [ (validate.rules).string.min_len = 1 ];
        bytes object = 2;
      }
  }

  message Filter {
      message Target {
        string host = 1;
        uint32 port = 2;
      }

      message Config {
        repeated Target targets = 1;
      }
  }
  ```

  - Update.MultiObjectRequest

    |  field   | type                 | label                                 | required | desc.            |
    | :------: | :------------------- | :------------------------------------ | :------: | :--------------- |
    | requests | Insert.ObjectRequest | repeated(Array[Insert.ObjectRequest]) |    \*    | the request list |

  - Update.ObjectRequest

    |   field    | type          | label | required | desc.                                   |
    | :--------: | :------------ | :---- | :------: | :-------------------------------------- |
    |   object   | Object.Blob   |       |    \*    | the binary object to be updated         |
    |   config   | Config        |       |    \*    | the configuration of the update request |
    | vectorizer | Filter.Target |       |    \*    | filter target                           |

  - Object.Blob

    | field  | type   | label | required | desc.             |
    | :----: | :----- | :---- | :------: | :---------------- |
    |   id   | string |       |    \*    | the object ID     |
    | object | bytes  |       |    \*    | the binary object |

  - Update.Config

    |          field          | type          | label | required | desc.                                                                                               |
    | :---------------------: | :------------ | :---- | :------: | :-------------------------------------------------------------------------------------------------- |
    | skip_strict_exist_check | bool          |       |          | check the same vector is already updated or not.<br>the ID should be unique if the value is `true`. |
    |        timestamp        | int64         |       |          | the timestamp of the vector updated.<br>if it is N/A, the current time will be used.                |
    |         filters         | Filter.Config |       |          | configuration for filter                                                                            |

  - Filter.Target

    | field | type   | label | required | desc.               |
    | :---: | :----- | :---- | :------: | :------------------ |
    | host  | string |       |    \*    | the target hostname |
    | port  | port   |       |    \*    | the target port     |

  - Filter.Config

    |  field  | type          | label                          | required | desc.                           |
    | :-----: | :------------ | :----------------------------- | :------: | :------------------------------ |
    | targets | Filter.Target | repeated(Array[Filter.Target]) |    \*    | the filter target configuration |

### Output

- the scheme of `payload.v1.Object.Locations`

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
    |  ips  | string | repeated(Array[string]) | the IP list of `vald-agent` pods where the request vector is updated. |

### Status Code

| code | desc.            |
| :--: | :--------------- |
|  0   | OK               |
|  3   | INVALID_ARGUMENT |
|  6   | ALREADY_EXISTS   |
|  13  | INTERNAL         |

## Upsert RPC

Upsert RPC is the method to update a single object and add a new single object.

```rpc
service Filter {
  rpc UpsertObject(payload.v1.Upsert.ObjectRequest)
      returns (payload.v1.Object.Location) {
    option (google.api.http) = {
      post : "/upsert/object"
      body : "*"
    };
  }
}
```

### Input

- the scheme of `payload.v1.Upsert.ObjectRequest`

  ```rpc
  message Upsert {
      message ObjectRequest {
        Object.Blob object = 1;
        Config config = 2;
        Filter.Target vectorizer = 3;
      }

      message Config {
        bool skip_strict_exist_check = 1;
        Filter.Config filters = 2;
        int64 timestamp = 3;
      }
  }

  message Object {
      message Blob {
        string id = 1 [ (validate.rules).string.min_len = 1 ];
        bytes object = 2;
      }
  }

  message Filter {
      message Target {
        string host = 1;
        uint32 port = 2;
      }

      message Config {
        repeated Target targets = 1;
      }
  }
  ```

  - Upsert.ObjectRequest

    |   field    | type          | label | required | desc.                                   |
    | :--------: | :------------ | :---- | :------: | :-------------------------------------- |
    |   object   | Object.Blob   |       |    \*    | the binary object to be upserted        |
    |   config   | Config        |       |    \*    | the configuration of the upsert request |
    | vectorizer | Filter.Target |       |    \*    | filter target                           |

  - Object.Blob

    | field  | type   | label | required | desc.             |
    | :----: | :----- | :---- | :------: | :---------------- |
    |   id   | string |       |    \*    | the object ID     |
    | object | bytes  |       |    \*    | the binary object |

  - Update.Config

    |          field          | type          | label | required | desc.                                                                                                |
    | :---------------------: | :------------ | :---- | :------: | :--------------------------------------------------------------------------------------------------- |
    | skip_strict_exist_check | bool          |       |          | check the same vector is already upserted or not.<br>the ID should be unique if the value is `true`. |
    |        timestamp        | int64         |       |          | the timestamp of the vector upserted.<br>if it is N/A, the current time will be used.                |
    |         filters         | Filter.Config |       |          | configuration for filter                                                                             |

  - Filter.Target

    | field | type   | label | required | desc.               |
    | :---: | :----- | :---- | :------: | :------------------ |
    | host  | string |       |    \*    | the target hostname |
    | port  | port   |       |    \*    | the target port     |

  - Filter.Config

    |  field  | type          | label                          | required | desc.                           |
    | :-----: | :------------ | :----------------------------- | :------: | :------------------------------ |
    | targets | Filter.Target | repeated(Array[Filter.Target]) |    \*    | the filter target configuration |

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

    | field | type   | label                   | desc.                                                                  |
    | :---: | :----- | :---------------------- | :--------------------------------------------------------------------- |
    | name  | string |                         | the name of vald agent pod where the request vector is upserted.       |
    | uuid  | string |                         | the ID of the upserted vector. It is the same as an `Object.Vector`.   |
    |  ips  | string | repeated(Array[string]) | the IP list of `vald-agent` pods where the request vector is upserted. |

### Status Code

| code | desc.            |
| :--: | :--------------- |
|  0   | OK               |
|  3   | INVALID_ARGUMENT |
|  6   | ALREADY_EXISTS   |
|  13  | INTERNAL         |

## StreamUpsert RPC

Upsert RPC is the method to update a single object and add a new single object.

```rpc
service Filter {
  rpc StreamUpsertObject(stream payload.v1.Upsert.ObjectRequest)
      returns (stream payload.v1.Object.StreamLocation) {}
}
```

### Input

- the scheme of `payload.v1.Upsert.ObjectRequest stream`

  ```rpc
  message Upsert {
      message ObjectRequest {
        Object.Blob object = 1;
        Config config = 2;
        Filter.Target vectorizer = 3;
      }

      message Config {
        bool skip_strict_exist_check = 1;
        Filter.Config filters = 2;
        int64 timestamp = 3;
      }
  }

  message Object {
      message Blob {
        string id = 1 [ (validate.rules).string.min_len = 1 ];
        bytes object = 2;
      }
  }

  message Filter {
      message Target {
        string host = 1;
        uint32 port = 2;
      }

      message Config {
        repeated Target targets = 1;
      }
  }
  ```

  - Upsert.ObjectRequest

    |   field    | type          | label | required | desc.                                   |
    | :--------: | :------------ | :---- | :------: | :-------------------------------------- |
    |   object   | Object.Blob   |       |    \*    | the binary object to be upserted        |
    |   config   | Config        |       |    \*    | the configuration of the upsert request |
    | vectorizer | Filter.Target |       |    \*    | filter target                           |

  - Object.Blob

    | field  | type   | label | required | desc.             |
    | :----: | :----- | :---- | :------: | :---------------- |
    |   id   | string |       |    \*    | the object ID     |
    | object | bytes  |       |    \*    | the binary object |

  - Update.Config

    |          field          | type          | label | required | desc.                                                                                                |
    | :---------------------: | :------------ | :---- | :------: | :--------------------------------------------------------------------------------------------------- |
    | skip_strict_exist_check | bool          |       |          | check the same vector is already upserted or not.<br>the ID should be unique if the value is `true`. |
    |        timestamp        | int64         |       |          | the timestamp of the vector upserted.<br>if it is N/A, the current time will be used.                |
    |         filters         | Filter.Config |       |          | configuration for filter                                                                             |

  - Filter.Target

    | field | type   | label | required | desc.               |
    | :---: | :----- | :---- | :------: | :------------------ |
    | host  | string |       |    \*    | the target hostname |
    | port  | port   |       |    \*    | the target port     |

  - Filter.Config

    |  field  | type          | label                          | required | desc.                           |
    | :-----: | :------------ | :----------------------------- | :------: | :------------------------------ |
    | targets | Filter.Target | repeated(Array[Filter.Target]) |    \*    | the filter target configuration |

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
    |  status  | google.rpc.Status |       | the status of google RPC.                  |

  - Object.Location

    | field | type   | label                   | desc.                                                                  |
    | :---: | :----- | :---------------------- | :--------------------------------------------------------------------- |
    | name  | string |                         | the name of vald agent pod where the request vector is upserted.       |
    | uuid  | string |                         | the ID of the upserted vector. It is the same as an `Object.Vector`.   |
    |  ips  | string | repeated(Array[string]) | the IP list of `vald-agent` pods where the request vector is upserted. |

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

MultiUpsert is the method to update existing multiple objects and add new multiple objects in **1** request.

<div class="notice">
gRPC has a message size limitation.<br>
Please be careful that the size of the request exceeds the limit.
</div>

```rpc
service Filter {
  rpc MultiUpsertObject(payload.v1.Upsert.MultiObjectRequest)
      returns (payload.v1.Object.Locations) {
    option (google.api.http) = {
      post : "/upsert/object/multiple"
      body : "*"
    };
  }
}
```

### Input

- the scheme of `payload.v1.Upsert.MultiObjectRequest`

  ```rpc
  message Upsert {
      message MultiObjectRequest {
        repeated ObjectRequest requests = 1;
      }

      message ObjectRequest {
        Object.Blob object = 1;
        Config config = 2;
        Filter.Target vectorizer = 3;
      }

      message Config {
        bool skip_strict_exist_check = 1;
        Filter.Config filters = 2;
        int64 timestamp = 3;
      }
  }

  message Object {
      message Blob {
        string id = 1 [ (validate.rules).string.min_len = 1 ];
        bytes object = 2;
      }
  }

  message Filter {
      message Target {
        string host = 1;
        uint32 port = 2;
      }

      message Config {
        repeated Target targets = 1;
      }
  }
  ```

  - Upsert.MultiObjectRequest

    |  field   | type                 | label                                 | required | desc.            |
    | :------: | :------------------- | :------------------------------------ | :------: | :--------------- |
    | requests | Upsert.ObjectRequest | repeated(Array[Upsert.ObjectRequest]) |    \*    | the request list |

  - Upsert.ObjectRequest

    |   field    | type          | label | required | desc.                                   |
    | :--------: | :------------ | :---- | :------: | :-------------------------------------- |
    |   object   | Object.Blob   |       |    \*    | the binary object to be upserted        |
    |   config   | Config        |       |    \*    | the configuration of the upsert request |
    | vectorizer | Filter.Target |       |    \*    | filter target                           |

  - Object.Blob

    | field  | type   | label | required | desc.             |
    | :----: | :----- | :---- | :------: | :---------------- |
    |   id   | string |       |    \*    | the object ID     |
    | object | bytes  |       |    \*    | the binary object |

  - Update.Config

    |          field          | type          | label | required | desc.                                                                                                |
    | :---------------------: | :------------ | :---- | :------: | :--------------------------------------------------------------------------------------------------- |
    | skip_strict_exist_check | bool          |       |          | check the same vector is already upserted or not.<br>the ID should be unique if the value is `true`. |
    |        timestamp        | int64         |       |          | the timestamp of the vector upserted.<br>if it is N/A, the current time will be used.                |
    |         filters         | Filter.Config |       |          | configuration for filter                                                                             |

  - Filter.Target

    | field | type   | label | required | desc.               |
    | :---: | :----- | :---- | :------: | :------------------ |
    | host  | string |       |    \*    | the target hostname |
    | port  | port   |       |    \*    | the target port     |

  - Filter.Config

    |  field  | type          | label                          | required | desc.                           |
    | :-----: | :------------ | :----------------------------- | :------: | :------------------------------ |
    | targets | Filter.Target | repeated(Array[Filter.Target]) |    \*    | the filter target configuration |

### Output

- the scheme of `payload.v1.Object.Locations`

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

    | field | type   | label                   | desc.                                                                  |
    | :---: | :----- | :---------------------- | :--------------------------------------------------------------------- |
    | name  | string |                         | the name of vald agent pod where the request vector is upserted.       |
    | uuid  | string |                         | the ID of the upserted vector. It is the same as an `Object.Vector`.   |
    |  ips  | string | repeated(Array[string]) | the IP list of `vald-agent` pods where the request vector is upserted. |

### Status Code

| code | desc.            |
| :--: | :--------------- |
|  0   | OK               |
|  3   | INVALID_ARGUMENT |
|  6   | ALREADY_EXISTS   |
|  13  | INTERNAL         |

## Search RPC

Search RPC is the method to search object(s) similar to request object.

```rpc
service Filter {
  rpc SearchObject(payload.v1.Search.ObjectRequest)
      returns (payload.v1.Search.Response) {
    option (google.api.http) = {
      post : "/search/object"
      body : "*"
    };
  }
}
```

### Input

- the scheme of `payload.v1.Search.ObjectRequest`

  ```rpc
  message Search {
      message ObjectRequest {
        bytes object = 1;
        Config config = 2;
        Filter.Target vectorizer = 3;
      }

      message Config {
        string request_id = 1;
        uint32 num = 2 [ (validate.rules).uint32.gte = 1 ];
        float radius = 3;
        float epsilon = 4;
        int64 timeout = 5;
        Filter.Config ingress_filters = 6;
        Filter.Config egress_filters = 7;
        uint32 min_num = 8 [ (validate.rules).uint32.gte = 0 ];
      }
  }

  message Filter {
      message Target {
        string host = 1;
        uint32 port = 2;
      }

      message Config {
        repeated Target targets = 1;
      }
  }
  ```

  - Search.ObjectRequest

    |   field    | type          | label | required | desc.                                   |
    | :--------: | :------------ | :---- | :------: | :-------------------------------------- |
    |   object   | bytes         |       |    \*    | the binary object to be searched        |
    |   config   | Config        |       |    \*    | the configuration of the search request |
    | vectorizer | Filter.Target |       |    \*    | filter target                           |

  - Search.Config

    |      field      | type          | label | required | desc.                                                 |
    | :-------------: | :------------ | :---- | :------: | :---------------------------------------------------- |
    |   request_id    | string        |       |          | unique request ID                                     |
    |       num       | uint32        |       |    \*    | the maximum number of result to be returned           |
    |     radius      | float         |       |    \*    | the search radius                                     |
    |     epsilon     | float         |       |    \*    | the search coefficient (default value is `0.1`)       |
    |     timeout     | int64         |       |          | Search timeout in nanoseconds (default value is `5s`) |
    | ingress_filters | Filter.Config |       |          | Ingress Filter configuration                          |
    | egress_filters  | Filter.Config |       |          | Egress Filter configuration                           |
    |     min_num     | uint32        |       |          | the minimum number of result to be returned           |

### Output

- the scheme of `payload.v1.Search.Response`.

  ```rpc
  message Search {
    message Response {
      string request_id = 1;
      repeated Object.Distance results = 2;
    }
  }

  message Object {
    message Distance {
      string id = 1;
      float distance = 2;
    }
  }
  ```

  - Search.Response

    |   field    | type            | label                            | desc.                 |
    | :--------: | :-------------- | :------------------------------- | :-------------------- |
    | request_id | string          |                                  | the unique request ID |
    |  results   | Object.Distance | repeated(Array[Object.Distance]) | search results        |

  - Object.Distance

    |  field   | type   | label | desc.                                                 |
    | :------: | :----- | :---- | :---------------------------------------------------- |
    |    id    | string |       | the vector ID                                         |
    | distance | float  |       | the distance between result vector and request vector |

### Status Code

| code | desc.             |
| :--: | :---------------- |
|  0   | OK                |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |

## StreamSearch RPC

StreamSearch RPC is the method to search vectors with multi queries(objects) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
By using the bidirectional streaming RPC, the search request can be communicated in any order between client and server.
Each Search request and response are independent.

```rpc
service Filter {
  rpc StreamSearchObject(stream payload.v1.Search.ObjectRequest)
      returns (stream payload.v1.Search.StreamResponse) {}
}
```

### Input

- the scheme of `payload.v1.Search.ObjectRequest stream`

  ```rpc
  message Search {
      message ObjectRequest {
        bytes object = 1;
        Config config = 2;
        Filter.Target vectorizer = 3;
      }

      message Config {
        string request_id = 1;
        uint32 num = 2 [ (validate.rules).uint32.gte = 1 ];
        float radius = 3;
        float epsilon = 4;
        int64 timeout = 5;
        Filter.Config ingress_filters = 6;
        Filter.Config egress_filters = 7;
        uint32 min_num = 8 [ (validate.rules).uint32.gte = 0 ];
      }
  }

  message Filter {
      message Target {
        string host = 1;
        uint32 port = 2;
      }

      message Config {
        repeated Target targets = 1;
      }
  }
  ```

  - Search.ObjectRequest

    |   field    | type          | label | required | desc.                                   |
    | :--------: | :------------ | :---- | :------: | :-------------------------------------- |
    |   object   | bytes         |       |    \*    | the binary object to be searched        |
    |   config   | Config        |       |    \*    | the configuration of the search request |
    | vectorizer | Filter.Target |       |    \*    | filter target                           |

  - Search.Config

    |      field      | type          | label | required | desc.                                                 |
    | :-------------: | :------------ | :---- | :------: | :---------------------------------------------------- |
    |   request_id    | string        |       |          | unique request ID                                     |
    |       num       | uint32        |       |    \*    | the maximum number of result to be returned           |
    |     radius      | float         |       |    \*    | the search radius                                     |
    |     epsilon     | float         |       |    \*    | the search coefficient (default value is `0.1`)       |
    |     timeout     | int64         |       |          | Search timeout in nanoseconds (default value is `5s`) |
    | ingress_filters | Filter.Config |       |          | Ingress Filter configuration                          |
    | egress_filters  | Filter.Config |       |          | Egress Filter configuration                           |
    |     min_num     | uint32        |       |          | the minimum number of result to be returned           |

### Output

- the scheme of `payload.v1.Search.StreamResponse`.

  ```rpc
  message Search {
    message StreamResponse {
      oneof payload {
        Response response = 1;
        google.rpc.Status status = 2;
      }
    }

    message Response {
      string request_id = 1;
      repeated Object.Distance results = 2;
    }
  }

  message Object {
    message Distance {
      string id = 1;
      float distance = 2;
    }
  }
  ```

  - Search.StreamResponse

    |  field   | type              | label | desc.                      |
    | :------: | :---------------- | :---- | :------------------------- |
    | response | Response          |       | the search result response |
    |  status  | google.rpc.Status |       | the status of google RPC   |

  - Search.Response

    |   field    | type            | label                            | desc.                 |
    | :--------: | :-------------- | :------------------------------- | :-------------------- |
    | request_id | string          |                                  | the unique request ID |
    |  results   | Object.Distance | repeated(Array[Object.Distance]) | search results        |

  - Object.Distance

    |  field   | type   | label | desc.                                                 |
    | :------: | :----- | :---- | :---------------------------------------------------- |
    |    id    | string |       | the vector ID                                         |
    | distance | float  |       | the distance between result vector and request vector |

## MultiSearch RPC

MultiSearch RPC is the method to search objects with multiple objects in **1** request.

<div class="notice">
gRPC has a message size limitation.<br>
Please be careful that the size of the request exceeds the limit.
</div>

```rpc
service Filter {
  rpc MultiSearchObject(payload.v1.Search.MultiObjectRequest)
      returns (payload.v1.Search.Responses) {
    option (google.api.http) = {
      post : "/search/object/multiple"
      body : "*"
    };
  }
}
```

### Input

- the scheme of `payload.v1.Search.MultiObjectRequest`

  ```rpc
  message Search {
      message MultiObjectRequest {
        repeated ObjectRequest requests = 1;
      }

      message ObjectRequest {
        bytes object = 1;
        Config config = 2;
        Filter.Target vectorizer = 3;
      }

      message Config {
        string request_id = 1;
        uint32 num = 2 [ (validate.rules).uint32.gte = 1 ];
        float radius = 3;
        float epsilon = 4;
        int64 timeout = 5;
        Filter.Config ingress_filters = 6;
        Filter.Config egress_filters = 7;
        uint32 min_num = 8 [ (validate.rules).uint32.gte = 0 ];
      }
  }

  message Filter {
      message Target {
        string host = 1;
        uint32 port = 2;
      }

      message Config {
        repeated Target targets = 1;
      }
  }
  ```

  - Search.MultiObjectRequest

    |  field   | type                                | label | required | desc.                   |
    | :------: | :---------------------------------- | :---- | :------: | :---------------------- |
    | requests | repeated(Array[MultiObjectRequest]) |       |    \*    | the search request list |

  - Search.ObjectRequest

    |   field    | type          | label | required | desc.                                   |
    | :--------: | :------------ | :---- | :------: | :-------------------------------------- |
    |   object   | bytes         |       |    \*    | the binary object to be searched        |
    |   config   | Config        |       |    \*    | the configuration of the search request |
    | vectorizer | Filter.Target |       |    \*    | filter target                           |

  - Search.Config

    |      field      | type          | label | required | desc.                                                 |
    | :-------------: | :------------ | :---- | :------: | :---------------------------------------------------- |
    |   request_id    | string        |       |          | unique request ID                                     |
    |       num       | uint32        |       |    \*    | the maximum number of result to be returned           |
    |     radius      | float         |       |    \*    | the search radius                                     |
    |     epsilon     | float         |       |    \*    | the search coefficient (default value is `0.1`)       |
    |     timeout     | int64         |       |          | Search timeout in nanoseconds (default value is `5s`) |
    | ingress_filters | Filter.Config |       |          | Ingress Filter configuration                          |
    | egress_filters  | Filter.Config |       |          | Egress Filter configuration                           |
    |     min_num     | uint32        |       |          | the minimum number of result to be returned           |

### Output

- the scheme of `payload.v1.Search.Responses`.

  ```rpc
  message Search {
    message Responses {
      repeated Response responses = 1;
    }

    message Response {
      string request_id = 1;
      repeated Object.Distance results = 2;
    }
  }

  message Object {
    message Distance {
      string id = 1;
      float distance = 2;
    }
  }
  ```

  - Search.Responses

    |   field   | type     | label                     | desc.                               |
    | :-------: | :------- | :------------------------ | :---------------------------------- |
    | responses | Response | repeated(Array[Response]) | the list of search results response |

  - Search.Response

    |   field    | type            | label                            | desc.                 |
    | :--------: | :-------------- | :------------------------------- | :-------------------- |
    | request_id | string          |                                  | the unique request ID |
    |  results   | Object.Distance | repeated(Array[Object.Distance]) | search results        |

  - Object.Distance

    |  field   | type   | label | desc.                                                 |
    | :------: | :----- | :---- | :---------------------------------------------------- |
    |    id    | string |       | the vector ID                                         |
    | distance | float  |       | the distance between result vector and request vector |

### Status Code

| code | desc.             |
| :--: | :---------------- |
|  0   | OK                |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
