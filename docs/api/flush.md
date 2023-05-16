# Vald Flush APIs

## Overview

Flush Service is responsible for removing vectors that are indexed in the `vald-agent`.

```rpc
service Flush {

  rpc Flush(payload.v1.Flush.Request) returns (payload.v1.Info.Index.Count) {}

}
```

## Flush RPC

Flush RPC is the method to remove all vectors.

### Input

- the scheme of `payload.v1.Flush.Request`

  ```rpc
  message Flush {
      message Request {

      }
  }
  ```

  - Flush.Request

    | field  | type      | label | required | desc.                                   |
    | :----: | :-------- | :---- | :------: | :-------------------------------------- |


### Output

- the scheme of `payload.v1.Info.Index.Count`

  ```rpc
  message Object {
      message Info_Index_Count {
        int64 stored = 0;
        int64 uncommitted = 0;
        bool indexing = false;
        bool saving = false;
      }
  }
  ```

  - Object.Info_Index_Count

    | field | type  | label                   | desc.                                                                      |
    |:-----:| :----- |:---------------------------------------------------------------------------| :-------------------------------------------------------------------- |
    | stored  | int64 |                         | count of indices.                                                          |
    | uncommitted  | int64 |                         | count of uncommited indices.                                               |
    |  indexing  | bool  | repeated(Array[string]) | the state indicating whether `vald-agent` pods is present in the indexing. |
    |  saving  | bool  | repeated(Array[string]) | the state indicating whether `vald-agent` pods is present in the saving.   |

### Status Code

| code | desc.             |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  13  | INTERNAL          |
