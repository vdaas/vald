# Vald Mirror Gateway APIs

## Overview

Mirror Service is responsible for providing a register and advertises interface for the `Val Mirror Gateway`.

```rpc
service Mirror {
  rpc Register(payload.v1.Mirror.Targets) returns (payload.v1.Mirror.Targets) {}

  rpc Advertise(payload.v1.Mirror.Targets) returns (payload.v1.Mirror.Targets) {}
}
```

## Register RPC

Register RPC is the method to register other Vald Mirror Gateway targets.

### Input

- the scheme of `payload.v1.Mirror.Targets`.

```rpc
message Mirror {
  message Target {
    string host = 1;
    uint32 port = 2;
  }

  message Targets {
    repeated Target targets = 1;
  }
}
```

- Mirror.Targets

  |  field  | type          | label                          | required | desc.                            |
  | :-----: | :------------ | :----------------------------- | :------: | :------------------------------- |
  | targets | Mirror.Target | repeated(Array[Mirror.Target]) |          | The multiple target information. |

- Mirror.Target

  | field | type   | label | required | desc.                |
  | :---: | :----- | :---- | :------: | :------------------- |
  | host  | string |       |    \*    | The target hostname. |
  | port  | uint32 |       |    \*    | The target port.     |

### Output

- the scheme of `payload.v1.Mirror.Targets`.
- Mirror.Targets

  |  field  | type          | label                          | required | desc.                            |
  | :-----: | :------------ | :----------------------------- | :------: | :------------------------------- |
  | targets | Mirror.Target | repeated(Array[Mirror.Target]) |          | The multiple target information. |

- Mirror.Target

  | field | type   | label | required | desc.                |
  | :---: | :----- | :---- | :------: | :------------------- |
  | host  | string |       |    \*    | The target hostname. |
  | port  | uint32 |       |    \*    | The target port.     |

### Status Code

| code | desc.             |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  13  | INTERNAL          |

## Advertise RPC

Advertise RPC is the method to advertise Vald Mirror Gateway targets.

### Input

- the scheme of `payload.v1.Mirror.Targets`.

```rpc
message Mirror {
  message Target {
    string host = 1;
    uint32 port = 2;
  }

  message Targets {
    repeated Target targets = 1;
  }
}
```

- Mirror.Targets

  |  field  | type          | label                          | required | desc.                            |
  | :-----: | :------------ | :----------------------------- | :------: | :------------------------------- |
  | targets | Mirror.Target | repeated(Array[Mirror.Target]) |          | The multiple target information. |

- Mirror.Target

  | field | type   | label | required | desc.                |
  | :---: | :----- | :---- | :------: | :------------------- |
  | host  | string |       |    \*    | The target hostname. |
  | port  | uint32 |       |    \*    | The target port.     |

### Output

- the scheme of `payload.v1.Mirror.Targets`.
- Mirror.Targets

  |  field  | type          | label                          | required | desc.                            |
  | :-----: | :------------ | :----------------------------- | :------: | :------------------------------- |
  | targets | Mirror.Target | repeated(Array[Mirror.Target]) |          | The multiple target information. |

- Mirror.Target

  | field | type   | label | required | desc.                |
  | :---: | :----- | :---- | :------: | :------------------- |
  | host  | string |       |    \*    | The target hostname. |
  | port  | uint32 |       |    \*    | The target port.     |

### Status Code

| code | desc.             |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  13  | INTERNAL          |
