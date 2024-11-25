# Vald Flush APIs

## Overview```rpc
service Flush {

  rpc Flush(payload.v1.Flush.Request) returns (payload.v1.Info.Index.Count) {}

}
```
## Flush RPC

### Input

- the scheme of `payload.v1.Flush.Request`

  ```rpc
  message Flush.Request {
    // empty
  }
  ```

  - Flush.Request

    empty

### Output

- the scheme of `payload.v1.Info.Index.Count`

  ```rpc
  message Info.Index.Count {
    uint32 stored = 1;
    uint32 uncommitted = 2;
    bool indexing = 3;
    bool saving = 4;
  }
  ```

  - Info.Index.Count

    | field | type | label | desc. |
    | :---: | :--- | :---- | :---- |
    | stored | uint32 |  | The stored index count. |
    | uncommitted | uint32 |  | The uncommitted index count. |
    | indexing | bool |  | The indexing index count. |
    | saving | bool |  | The saving index count. |



