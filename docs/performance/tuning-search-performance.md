# Tuning Search Performance

ANN is fast, but sometimes it can be improved more by tuning parameters.
This page shows the essence for improve ANN search by the Vald cluster.

## Tuning Guideline

When the search results do NOT satisfy the expected result, it can be improved by tuning parameters.

First of all, we recommend tuning by following the steps below without doing it blindly.

<div class="mermaid">
flowchart TD
    A[Perform Linear Search API]
    B{Is satisfies?}
    C[Tuning Parameters to improve precision]
    D[Tuning Embedding Models]
    E[Perform Search API]
    F[Tuning Parameters to improve latency]
    A-->B
    B-- Yes -->C
    B-- No -->D
    C--> E
    E--> C
    E--> F
    F--> E
</div>

The best practice is:

1. Measure Linear Search performance and use it as a baseline for Search API
1. Repeat tuning to improve precision and measure Search API until the conditions are met
1. Repeat tuning to improve latency and measure Search API until the conditions are met

<div class="notice">
When the results are not good by Linear Search API, it may need to rethink the embedding model for vectorization.
</div>

## Tuning parameters

There are two viewpoints, client-side and cluster-side, for improving search performance.

### Client side

On the client side, parameters of `Search.Config` will affect the search result.

|             | description                                       | how does it affect?                                                                                                                                                                                                    | memo                                                          |
| :---------- | :------------------------------------------------ | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :------------------------------------------------------------ |
| radius      | the search radius for NGT                         | Define the search range when NGT searches the nearest neighbors                                                                                                                                                        | -1 is recommended.                                            |
| epsilon     | the search coefficient for NGT                    | Expansion factor of the NGT search range.<BR>Search operation time increases when the epsilon is big.                                                                                                                  | recommended value range: `0.01 ~ 0.5`<BR>default value: `0.1` |
| timeout(ns) | max time duration until receiving search results. | An error will be returned if the set `num` search results cannot be obtained within the set time.<BR>By setting `min_num`, the search results will be returned if more than `min_num` can be searched within the time. | default value: `3,000,000,000ns`                              |

### Cluster-side

On the cluster side, these parameters can be set by `values.yaml`, affect the search result.

|                              | description                                                   | how does it affect?                                                                                                                                                                        | Memo                |
| :--------------------------- | :------------------------------------------------------------ | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :------------------ |
| agent.ngt.creation_edge_size | Number of nodes connected to one node                         | It helps reduce unreachable edges.<BR>The larger it is, the denser the graph structure will be, but the memory usage, search speed, and index construction time will increase accordingly. | default value: `20` |
| agent.ngt.search_edge_size   | Number of nodes to search from the origin node when searching | The number of nodes to search will increase.<BR>Accuracy will be higher, but speed will be lower.<BR>Adjust if adjusting the radius and epsilon does not improve the situation.            | default value: `10` |
