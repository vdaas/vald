# Vald Mirror APIs

## Overview

Mirror Service is responsible for providing the `Register` interface for the Vald Mirror Gateway.

```rpc
service Mirror {

  rpc Register(payload.v1.Mirror.Targets) returns (payload.v1.Mirror.Targets) {}

}
```

## Register RPC

Register RPC is the method to register other Vald Mirror Gateway targets.

### Input

- the scheme of `payload.v1.Mirror.Targets`

  ```rpc
  message Mirror.Targets {
    repeated Mirror.Target targets = 1;
  }

  message Mirror.Target {
    string host = 1;
    uint32 port = 2;
  }

  ```

  - Mirror.Targets

    | field | type | label | description |
    | :---: | :--- | :---- | :---------- |
    | targets | Mirror.Target | repeated | The multiple target information. |

  - Mirror.Target

    | field | type | label | description |
    | :---: | :--- | :---- | :---------- |
    | host | string |  | The target hostname. |
    | port | uint32 |  | The target port. |
### Output

- the scheme of `payload.v1.Mirror.Targets`

  ```rpc
  message Mirror.Targets {
    repeated Mirror.Target targets = 1;
  }

  message Mirror.Target {
    string host = 1;
    uint32 port = 2;
  }

  ```

  - Mirror.Targets

    | field | type | label | description |
    | :---: | :--- | :---- | :---------- |
    | targets | Mirror.Target | repeated | The multiple target information. |

  - Mirror.Target

    | field | type | label | description |
    | :---: | :--- | :---- | :---------- |
    | host | string |  | The target hostname. |
    | port | uint32 |  | The target port. |

### Status Code

| code | description       |
| :--: | :---------------- |
|  0   | OK                |
|  1   | CANCELLED         |
|  3   | INVALID_ARGUMENT  |
|  4   | DEADLINE_EXCEEDED |
|  13  | INTERNAL          |

Please refer to [Response Status Code](../status.md) for more details.

