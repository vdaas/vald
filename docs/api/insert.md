# Vald Insert APIs

## Overview

Insert Service is responsible for inserting new vectors into the `vald-agent`.

```bash
service Insert {
    rpc Insert(payload.v1.Insert.Request) returns (payload.v1.Object.Location) {}

    rpc StreamInsert(stream payload.v1.Insert.Request) returns (stream payload.v1.Object.Location) {}

    rpc MultiInsert(payload.v1.Insert.MultiRequest) returns (payload.v1.Object.Location) {}
    }
```

## Insert RPC

Inset RPC is the method to ad a new single vector.

### Input

- the scheme of `payload.v1.Insert.Request`

  ```bash
  message Insert {
      message Request {
          Object.Vector vector = 1 [ (validate.rules).repeated.min_items = 2 ];
          Config config =2;
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

  - Insert.Request
    |field|type|label|required|desc.|
    |:---:|:---|:---|:---:|:---|
    |vector|Object.Vector| | \* | the information of vector |
    |config|Config| | \* | configuration for inserting vector |

  - Insert.Config
    |field|type|label|required|desc.|
    |:---:|:---|:---|:---:|:---|
    |skip_strict_exist_check|bool| | | check the same vector is already inserted or not.<br>the same ID is not indexed if the value is `true`|
    |timestamp|int64| | | it shows the time of vector is inserted.<br>if it is N/A, Vald will use unix timestamp.
    |filters|Filter.Config| | | configuration for filter |

  - Object.Vector
    |field|type|label|required|desc.|
    |:---:|:---|:---|:---:|:---|
    |id|string| | \* | the ID of a vector. ID consists 1 or more strings. |
    |vector|float| repeated(Array[float]) | \* | the vector data. its dimension is between 2 and 65,536.|

### Output

- the scheme of `payload.v1.Object.Location`

  ```bash
  message Object {
      message Location {
        string name = 1;
        string uuid = 2;
        repeated string ips = 3;
      }
  }
  ```

  - Object.Location
    |field|type|label|desc.|
    |:---:|:---|:---|:---|
    |name|string| | the name of vald agent pod which has been inserted the request vector. |
    |uuid|string| | the ID of a inserted vector. it is same as Object.Vector |
    |ips|string| repeated(Array[string]) | the ip list of `vald-agent` pods which has been inserted the request vector. |

### Status Code

| code | desc.            |
| :--: | :--------------- |
|  0   | OK               |
|  3   | INVALID_ARGUMENT |
|  6   | ALREADY_EXISTS   |
|  13  | INTERNAL         |

## StreamInsert RPC<sup>recommended</sup>

StreamInset RPC is the method  to add new multiple vectors by the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
It can communicate between client and server in any order.
Each Insert request and response are independent.
It's a recommended method when a large amount of vector should be inserted.

### Input

- the scheme of `payload.v1.Insert.Request stream`

  ```bash
  message Insert {
      message Request {
          Object.Vector vector = 1 [ (validate.rules).repeated.min_items = 2 ];
          Config config =2;
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

  - Insert.Request
    |field|type|label|required|desc.|
    |:---:|:---|:---|:---:|:---|
    |vector|Object.Vector| | \* | the information of vector |
    |config|Config| | \* | configuration for inserting vector |

  - Insert.Config
    |field|type|label|required|desc.|
    |:---:|:---|:---|:---:|:---|
    |skip_strict_exist_check|bool| | | check the same vector is already inserted or not.<br>the same ID is not indexed if the value is `true`|
    |timestamp|int64| | | it shows the time of vector is inserted.<br>if it is N/A, Vald will use unix timestamp.
    |filters|Filter.Config| | | configuration for filter |

  - Object.Vector
    |field|type|label|required:|desc.|
    |:---:|:---|:---|:---:|:---|
    |id|string| | \* | the ID of a vector. ID consists 1 or more strings. |
    |vector|float| repeated(Array[float]) | \* | the vector data. its dimension is between 2 and 65,536. |

### Output

- the scheme of `payload.v1.Object.StreamLocation`

  ```bash
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
    |field|type|label|desc.|
    |:---:|:---|:---|:---|
    |location|Object.Location| | the information of Object.Location data. |
    |status|google.rpc.Status| | the status of google RPC. |

  - Object.Location
    |field|type|label|desc.|
    |:---:|:---|:---|:---|
    |name|string| | the name of vald agent pod which has been inserted the request vector. |
    |uuid|string| | the ID of a inserted vector. it is same as Object.Vector |
    |ips|string| repeated(Array[string]) | the ip list of `vald-agent` pods which has been inserted the request vector. |

  - [google.rpc.Status](https://github.com/googleapis/googleapis/blob/master/google/rpc/status.proto)
    |field|type|label|desc.|
    |:---:|:---|:---|:---|
    |code|int32| | status code (code list is next section)|
    |message|string| | error message |
    |details|google.protobuf.Any| repeated(Array[any]) | the details error message list|

### Status Code

| code | desc.            |
| :--: | :--------------- |
|  0   | OK               |
|  3   | INVALID_ARGUMENT |
|  6   | ALREADY_EXISTS   |
|  13  | INTERNAL         |

## MultiInsert RPC

MultiInsert is the method to add new multiple vectors in **1** request.

<div class="card-note">
gRPC has the limitation message size.<br>
Please check the request is smaller than it.
</div>

### Input

- the scheme of `payload.v1.Insert.MultiRequest`

  ```bash
  message Insert {
      message MultiRequest { repeated Request requests = 1; }

      message Request {
          Object.Vector vector = 1 [ (validate.rules).repeated.min_items = 2 ];
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

  - Insert.MultiRequest
    |field|type|label|required|desc.|
    |:---:|:---|:---|:---:|:---|
    |requests|Insert.Request| repeated(Array[Insert.Request]) | \* | the request list |

  - Insert.Request
    |field|type|label|required|desc.|
    |:---:|:---|:---|:---:|:---|
    |vector|Object.Vector| | \* | the information of vector |
    |config|Config| | \* | configuration for inserting vector |

  - Insert.Config
    |field|type|label|required|desc.|
    |:---:|:---|:---|:---:|:---|
    |skip_strict_exist_check|bool| | | check the same vector is already inserted or not.<br>the same ID is not indexed if the value is `true`|
    |timestamp|int64| | | it shows the time of vector is inserted.<br>if it is N/A, Vald will use unix timestamp.
    |filters|Filter.Config| | | configuration for filter |

  - Object.Vector
    |field|type|label|required|desc.|
    |:---:|:---|:---|:---:|:---|
    |id|string| | \* | the ID of a vector. ID consists 1 or more strings. |
    |vector|float| repeated(Array[float]) | \* | the vector data. its dimension is between 2 and 65,536.|

### Output

- the scheme of `payload.v1.Object.Locations`.

  ```bash
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
    |field|type|label|desc.|
    |:---:|:---|:---|:---|
    |location|Object.Location| repeated(Array[Object.Location]) | the list of Object.Location. |

  - Object.Location
    |field|type|label|desc.|
    |:---:|:---|:---|:---|
    |name|string| | the name of vald agent pod which has been inserted the request vector. |
    |uuid|string| | the ID of a inserted vector. it is same as Object.Vector |
    |ips|string| repeated(Array[string]) | the ip list of `vald-agent` pods which has been inserted the request vector. |

### Status Code

| code | desc.            |
| :--: | :--------------- |
|  0   | OK               |
|  3   | INVALID_ARGUMENT |
|  6   | ALREADY_EXISTS   |
|  13  | INTERNAL         |
