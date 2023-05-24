# Vald Flush APIs

## Overview

Flush Service is responsible for removing all vectors that are indexed and uncommitted in the `vald-agent`.

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

    empty


### Output

- the scheme of `payload.v1.Info.Index.Count`

  ```rpc
  message Object {
      message Info_Index_Count {
        uint32 stored = 0;
        uint32 uncommitted = 0;
        bool indexing = false;
        bool saving = false;
      }
  }
  ```

   Object.Info_Index_Count

    | field       | type   | label | desc.                                                                      |
    | :---------: |:------ | :---- | :------------------------------------------------------------------------- |
    | stored      | uint32 |       | count of indices.                                                          | 
    | uncommitted | uint32 |       | count of uncommitted indices.                                              |
    | indexing    |  bool  |       | the state indicating whether `vald-agent` pods is present in the indexing. |
    | saving      |  bool  |       | the state indicating whether `vald-agent` pods is present in the saving.   |

### Status Code

| code | desc.             |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  5   | NOT_FOUND         |
|  13  | INTERNAL          |
