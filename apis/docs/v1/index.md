# Vald Index APIs

## Overview

Represent the index manager service.

```rpc
service Index {

  rpc IndexInfo(payload.v1.Empty) returns (payload.v1.Info.Index.Count) {}
  rpc IndexDetail(payload.v1.Empty) returns (payload.v1.Info.Index.Detail) {}
  rpc IndexStatistics(payload.v1.Empty) returns (payload.v1.Info.Index.Statistics) {}
  rpc IndexStatisticsDetail(payload.v1.Empty) returns (payload.v1.Info.Index.StatisticsDetail) {}
  rpc IndexProperty(payload.v1.Empty) returns (payload.v1.Info.Index.PropertyDetail) {}

}
```

## IndexInfo RPC

Represent the RPC to get the index information.

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

    | field | type | label | description |
    | :---: | :--- | :---- | :---------- |
    | stored | uint32 |  | The stored index count. |
    | uncommitted | uint32 |  | The uncommitted index count. |
    | indexing | bool |  | The indexing index count. |
    | saving | bool |  | The saving index count. |

## IndexDetail RPC

Represent the RPC to get the index information for each agents.

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

- the scheme of `payload.v1.Info.Index.Detail`

  ```rpc
  message Info.Index.Detail {
    repeated Info.Index.Detail.CountsEntry counts = 1;
    uint32 replica = 2;
    uint32 live_agents = 3;
  }

  message Info.Index.Detail.CountsEntry {
    string key = 1;
    Info.Index.Count value = 2;
  }

  message Info.Index.Count {
    uint32 stored = 1;
    uint32 uncommitted = 2;
    bool indexing = 3;
    bool saving = 4;
  }

  ```

  - Info.Index.Detail

    | field | type | label | description |
    | :---: | :--- | :---- | :---------- |
    | counts | Info.Index.Detail.CountsEntry | repeated | count infos for each agents |
    | replica | uint32 |  | index replica of vald cluster |
    | live_agents | uint32 |  | live agent replica of vald cluster |

  - Info.Index.Detail.CountsEntry

    | field | type | label | description |
    | :---: | :--- | :---- | :---------- |
    | key | string |  |  |
    | value | Info.Index.Count |  |  |

  - Info.Index.Count

    | field | type | label | description |
    | :---: | :--- | :---- | :---------- |
    | stored | uint32 |  | The stored index count. |
    | uncommitted | uint32 |  | The uncommitted index count. |
    | indexing | bool |  | The indexing index count. |
    | saving | bool |  | The saving index count. |

## IndexStatistics RPC

Represent the RPC to get the index statistics.

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

- the scheme of `payload.v1.Info.Index.Statistics`

  ```rpc
  message Info.Index.Statistics {
    bool valid = 1;
    int32 median_indegree = 2;
    int32 median_outdegree = 3;
    uint64 max_number_of_indegree = 4;
    uint64 max_number_of_outdegree = 5;
    uint64 min_number_of_indegree = 6;
    uint64 min_number_of_outdegree = 7;
    uint64 mode_indegree = 8;
    uint64 mode_outdegree = 9;
    uint64 nodes_skipped_for_10_edges = 10;
    uint64 nodes_skipped_for_indegree_distance = 11;
    uint64 number_of_edges = 12;
    uint64 number_of_indexed_objects = 13;
    uint64 number_of_nodes = 14;
    uint64 number_of_nodes_without_edges = 15;
    uint64 number_of_nodes_without_indegree = 16;
    uint64 number_of_objects = 17;
    uint64 number_of_removed_objects = 18;
    uint64 size_of_object_repository = 19;
    uint64 size_of_refinement_object_repository = 20;
    double variance_of_indegree = 21;
    double variance_of_outdegree = 22;
    double mean_edge_length = 23;
    double mean_edge_length_for_10_edges = 24;
    double mean_indegree_distance_for_10_edges = 25;
    double mean_number_of_edges_per_node = 26;
    double c1_indegree = 27;
    double c5_indegree = 28;
    double c95_outdegree = 29;
    double c99_outdegree = 30;
    repeated int64 indegree_count = 31;
    repeated uint64 outdegree_histogram = 32;
    repeated uint64 indegree_histogram = 33;
  }

  ```

  - Info.Index.Statistics

    | field | type | label | description |
    | :---: | :--- | :---- | :---------- |
    | valid | bool |  |  |
    | median_indegree | int32 |  |  |
    | median_outdegree | int32 |  |  |
    | max_number_of_indegree | uint64 |  |  |
    | max_number_of_outdegree | uint64 |  |  |
    | min_number_of_indegree | uint64 |  |  |
    | min_number_of_outdegree | uint64 |  |  |
    | mode_indegree | uint64 |  |  |
    | mode_outdegree | uint64 |  |  |
    | nodes_skipped_for_10_edges | uint64 |  |  |
    | nodes_skipped_for_indegree_distance | uint64 |  |  |
    | number_of_edges | uint64 |  |  |
    | number_of_indexed_objects | uint64 |  |  |
    | number_of_nodes | uint64 |  |  |
    | number_of_nodes_without_edges | uint64 |  |  |
    | number_of_nodes_without_indegree | uint64 |  |  |
    | number_of_objects | uint64 |  |  |
    | number_of_removed_objects | uint64 |  |  |
    | size_of_object_repository | uint64 |  |  |
    | size_of_refinement_object_repository | uint64 |  |  |
    | variance_of_indegree | double |  |  |
    | variance_of_outdegree | double |  |  |
    | mean_edge_length | double |  |  |
    | mean_edge_length_for_10_edges | double |  |  |
    | mean_indegree_distance_for_10_edges | double |  |  |
    | mean_number_of_edges_per_node | double |  |  |
    | c1_indegree | double |  |  |
    | c5_indegree | double |  |  |
    | c95_outdegree | double |  |  |
    | c99_outdegree | double |  |  |
    | indegree_count | int64 | repeated |  |
    | outdegree_histogram | uint64 | repeated |  |
    | indegree_histogram | uint64 | repeated |  |

## IndexStatisticsDetail RPC

Represent the RPC to get the index statistics for each agents.

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

- the scheme of `payload.v1.Info.Index.StatisticsDetail`

  ```rpc
  message Info.Index.StatisticsDetail {
    repeated Info.Index.StatisticsDetail.DetailsEntry details = 1;
  }

  message Info.Index.StatisticsDetail.DetailsEntry {
    string key = 1;
    Info.Index.Statistics value = 2;
  }

  message Info.Index.Statistics {
    bool valid = 1;
    int32 median_indegree = 2;
    int32 median_outdegree = 3;
    uint64 max_number_of_indegree = 4;
    uint64 max_number_of_outdegree = 5;
    uint64 min_number_of_indegree = 6;
    uint64 min_number_of_outdegree = 7;
    uint64 mode_indegree = 8;
    uint64 mode_outdegree = 9;
    uint64 nodes_skipped_for_10_edges = 10;
    uint64 nodes_skipped_for_indegree_distance = 11;
    uint64 number_of_edges = 12;
    uint64 number_of_indexed_objects = 13;
    uint64 number_of_nodes = 14;
    uint64 number_of_nodes_without_edges = 15;
    uint64 number_of_nodes_without_indegree = 16;
    uint64 number_of_objects = 17;
    uint64 number_of_removed_objects = 18;
    uint64 size_of_object_repository = 19;
    uint64 size_of_refinement_object_repository = 20;
    double variance_of_indegree = 21;
    double variance_of_outdegree = 22;
    double mean_edge_length = 23;
    double mean_edge_length_for_10_edges = 24;
    double mean_indegree_distance_for_10_edges = 25;
    double mean_number_of_edges_per_node = 26;
    double c1_indegree = 27;
    double c5_indegree = 28;
    double c95_outdegree = 29;
    double c99_outdegree = 30;
    repeated int64 indegree_count = 31;
    repeated uint64 outdegree_histogram = 32;
    repeated uint64 indegree_histogram = 33;
  }

  ```

  - Info.Index.StatisticsDetail

    | field | type | label | description |
    | :---: | :--- | :---- | :---------- |
    | details | Info.Index.StatisticsDetail.DetailsEntry | repeated | count infos for each agents |

  - Info.Index.StatisticsDetail.DetailsEntry

    | field | type | label | description |
    | :---: | :--- | :---- | :---------- |
    | key | string |  |  |
    | value | Info.Index.Statistics |  |  |

  - Info.Index.Statistics

    | field | type | label | description |
    | :---: | :--- | :---- | :---------- |
    | valid | bool |  |  |
    | median_indegree | int32 |  |  |
    | median_outdegree | int32 |  |  |
    | max_number_of_indegree | uint64 |  |  |
    | max_number_of_outdegree | uint64 |  |  |
    | min_number_of_indegree | uint64 |  |  |
    | min_number_of_outdegree | uint64 |  |  |
    | mode_indegree | uint64 |  |  |
    | mode_outdegree | uint64 |  |  |
    | nodes_skipped_for_10_edges | uint64 |  |  |
    | nodes_skipped_for_indegree_distance | uint64 |  |  |
    | number_of_edges | uint64 |  |  |
    | number_of_indexed_objects | uint64 |  |  |
    | number_of_nodes | uint64 |  |  |
    | number_of_nodes_without_edges | uint64 |  |  |
    | number_of_nodes_without_indegree | uint64 |  |  |
    | number_of_objects | uint64 |  |  |
    | number_of_removed_objects | uint64 |  |  |
    | size_of_object_repository | uint64 |  |  |
    | size_of_refinement_object_repository | uint64 |  |  |
    | variance_of_indegree | double |  |  |
    | variance_of_outdegree | double |  |  |
    | mean_edge_length | double |  |  |
    | mean_edge_length_for_10_edges | double |  |  |
    | mean_indegree_distance_for_10_edges | double |  |  |
    | mean_number_of_edges_per_node | double |  |  |
    | c1_indegree | double |  |  |
    | c5_indegree | double |  |  |
    | c95_outdegree | double |  |  |
    | c99_outdegree | double |  |  |
    | indegree_count | int64 | repeated |  |
    | outdegree_histogram | uint64 | repeated |  |
    | indegree_histogram | uint64 | repeated |  |

## IndexProperty RPC

Represent the RPC to get the index property.

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

- the scheme of `payload.v1.Info.Index.PropertyDetail`

  ```rpc
  message Info.Index.PropertyDetail {
    repeated Info.Index.PropertyDetail.DetailsEntry details = 1;
  }

  message Info.Index.PropertyDetail.DetailsEntry {
    string key = 1;
    Info.Index.Property value = 2;
  }

  message Info.Index.Property {
    int32 dimension = 1;
    int32 thread_pool_size = 2;
    string object_type = 3;
    string distance_type = 4;
    string index_type = 5;
    string database_type = 6;
    string object_alignment = 7;
    int32 path_adjustment_interval = 8;
    int32 graph_shared_memory_size = 9;
    int32 tree_shared_memory_size = 10;
    int32 object_shared_memory_size = 11;
    int32 prefetch_offset = 12;
    int32 prefetch_size = 13;
    string accuracy_table = 14;
    string search_type = 15;
    float max_magnitude = 16;
    int32 n_of_neighbors_for_insertion_order = 17;
    float epsilon_for_insertion_order = 18;
    string refinement_object_type = 19;
    int32 truncation_threshold = 20;
    int32 edge_size_for_creation = 21;
    int32 edge_size_for_search = 22;
    int32 edge_size_limit_for_creation = 23;
    double insertion_radius_coefficient = 24;
    int32 seed_size = 25;
    string seed_type = 26;
    int32 truncation_thread_pool_size = 27;
    int32 batch_size_for_creation = 28;
    string graph_type = 29;
    int32 dynamic_edge_size_base = 30;
    int32 dynamic_edge_size_rate = 31;
    float build_time_limit = 32;
    int32 outgoing_edge = 33;
    int32 incoming_edge = 34;
  }

  ```

  - Info.Index.PropertyDetail

    | field | type | label | description |
    | :---: | :--- | :---- | :---------- |
    | details | Info.Index.PropertyDetail.DetailsEntry | repeated |  |

  - Info.Index.PropertyDetail.DetailsEntry

    | field | type | label | description |
    | :---: | :--- | :---- | :---------- |
    | key | string |  |  |
    | value | Info.Index.Property |  |  |

  - Info.Index.Property

    | field | type | label | description |
    | :---: | :--- | :---- | :---------- |
    | dimension | int32 |  |  |
    | thread_pool_size | int32 |  |  |
    | object_type | string |  |  |
    | distance_type | string |  |  |
    | index_type | string |  |  |
    | database_type | string |  |  |
    | object_alignment | string |  |  |
    | path_adjustment_interval | int32 |  |  |
    | graph_shared_memory_size | int32 |  |  |
    | tree_shared_memory_size | int32 |  |  |
    | object_shared_memory_size | int32 |  |  |
    | prefetch_offset | int32 |  |  |
    | prefetch_size | int32 |  |  |
    | accuracy_table | string |  |  |
    | search_type | string |  |  |
    | max_magnitude | float |  |  |
    | n_of_neighbors_for_insertion_order | int32 |  |  |
    | epsilon_for_insertion_order | float |  |  |
    | refinement_object_type | string |  |  |
    | truncation_threshold | int32 |  |  |
    | edge_size_for_creation | int32 |  |  |
    | edge_size_for_search | int32 |  |  |
    | edge_size_limit_for_creation | int32 |  |  |
    | insertion_radius_coefficient | double |  |  |
    | seed_size | int32 |  |  |
    | seed_type | string |  |  |
    | truncation_thread_pool_size | int32 |  |  |
    | batch_size_for_creation | int32 |  |  |
    | graph_type | string |  |  |
    | dynamic_edge_size_base | int32 |  |  |
    | dynamic_edge_size_rate | int32 |  |  |
    | build_time_limit | float |  |  |
    | outgoing_edge | int32 |  |  |
    | incoming_edge | int32 |  |  |

