# Vald Stats APIs

## Overview

Represent the resource stats service for gateway.

```rpc
service Stats {

  rpc ResourceStatsDetail(payload.v1.Empty) returns (payload.v1.Info.Stats.ResourceStatsDetail) {}

}
```

## ResourceStatsDetail RPC

Represent the RPC to get the resource stats for each agents.

### Input

- the scheme of `payload.v1.Empty`

  ```rpc
  message Empty {
    // empty
  }

  ```

  - Empty

    empty

### Output

- the scheme of `payload.v1.Info.Stats.ResourceStatsDetail`

  ```rpc
  message Info.Stats.ResourceStatsDetail {
    repeated Info.Stats.ResourceStatsDetail.DetailsEntry details = 1;
  }

  message Info.Stats.ResourceStatsDetail.DetailsEntry {
    string key = 1;
    Info.Stats.ResourceStats value = 2;
  }

  message Info.Stats.ResourceStats {
    string name = 1;
    string ip = 2;
    Info.Stats.CgroupStats cgroup_stats = 3;
  }

  message Info.Stats.CgroupStats {
    double cpu_limit_cores = 1;
    double cpu_usage_cores = 2;
    uint64 memory_limit_bytes = 3;
    uint64 memory_usage_bytes = 4;
  }

  ```

  - Info.Stats.ResourceStatsDetail

    |  field  | type                                        | label    | description |
    | :-----: | :------------------------------------------ | :------- | :---------- |
    | details | Info.Stats.ResourceStatsDetail.DetailsEntry | repeated |             |

  - Info.Stats.ResourceStatsDetail.DetailsEntry

    | field | type                     | label | description |
    | :---: | :----------------------- | :---- | :---------- |
    |  key  | string                   |       |             |
    | value | Info.Stats.ResourceStats |       |             |

  - Info.Stats.ResourceStats

    |    field     | type                   | label | description                         |
    | :----------: | :--------------------- | :---- | :---------------------------------- |
    |     name     | string                 |       |                                     |
    |      ip      | string                 |       |                                     |
    | cgroup_stats | Info.Stats.CgroupStats |       | Container resource usage statistics |

  - Info.Stats.CgroupStats

    |       field        | type   | label | description                         |
    | :----------------: | :----- | :---- | :---------------------------------- |
    |  cpu_limit_cores   | double |       | CPU cores available                 |
    |  cpu_usage_cores   | double |       | CPU usage in cores (not percentage) |
    | memory_limit_bytes | uint64 |       | Memory limit in bytes               |
    | memory_usage_bytes | uint64 |       | Memory usage in bytes               |
