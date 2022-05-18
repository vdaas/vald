# Search Config

## Vald Search Service

<!-- Describe vald search service: search, searchbyid, linearsearch, lineasearchbyid ... -->

Vald provides the two types of search service.

1. Search

   - `Search` is just the `ANN(Approximate Nearest Neighbor)` search with query vector.
     It is the fast search even the vector consists large dimension.
     The search duration is fast, but less accurate than `LinearSearch`.
     Search algorithm depends on each core algorithm.

1. SearchById

   - `SearchById` is just the `ANN(Approximate Nearest Neighbor)` search with vector's id.
     The id should be already indexed before search process.
     The sarch algorithm is the same as `Search`.

1. LinearSearch

   - `LinearSearch` is the basic search algorithm with query vector.
     It searches all indexed vectors and calculates the distance between the query.
     Its accuracy is exact, but the search time requires more than `Search` (ANN search) and increases the amount of indexed vector.

1. LinearSearchById
   - `LinearSearchById` is the basic search algorithm with vector's id.
     The id should be already indexed before search process.
     The sarch algorithm is the same as `LinearSearch`.

<div class="notice">
Linear Search service is available from Vald v1.4 or later.
</div>

For more details, please refer to [Search API document](../api/search.md)

## Configuration

<!-- Describe Search parameters: ../../apis/proto/v1/payload -->
<!-- TODO: divide search configuration for each agent algorithm -->

Here is the current search config parameters.

```bash
// Represent search configuration.
message Config {
  // Unique request ID.
  string request_id = 1;
  // Maximum number of result to be returned.
  uint32 num = 2 [ (validate.rules).uint32.gte = 1 ];
  // Search radius.
  float radius = 3;
  // Search coefficient.
  float epsilon = 4;
  // Search timeout in nanoseconds.
  int64 timeout = 5;
  // Ingress filter configurations.
  Filter.Config ingress_filters = 6;
  // Egress filter configurations.
  Filter.Config egress_filters = 7;
  // Minimum number of result to be returned.
  uint32 min_num = 8 [ (validate.rules).uint32.gte = 0 ];
}
```

### request_id

`request_id` is unique request ID.
It is **NOT** indexed vector's id.
It can be used for e.g., error handling process.

### num

`num` is the maximum number of search result which you'd like to get.
`num` should be positive integer.

### radius

`radius`, the specific parameter for NGT, specifies the search range centered on the query vector in terms of the radius of a circle.
The number of search target vectors increases along with the radius is large.
There is the trade-off between accuracy and search speed.
In many cases, it is hard to set it due to depending on dataset.

The default value is infinite circle.
When setting negative number as `radius`, `NGT` applies the radius as infinite circle.

<div class="notice">
NGT will self-update radius during search process.
</div>

### epsilon

`epsilon`, the specific parameter for NGT, specifies the magnification coefficient (epsilon) of the search range.
NGT will use `(1+epsilon)*radius` as the search range.
The number of search target vectors increases along with the epsilon is large.

The default value (recommend value) is `0.1`.
While it is desirable to adjust this value within the range of 0 - 0.3, a negative value (over than `-1`) may also be specified.

### ingress_filters

`ingress_filters` is required when using ingress filter component.
It requires ingress filter component's hostname and port.

### egress_filters

`egress_filters` is required when using egress filter component.
It requires egress filter component's hostname and port.

### min_num

`min_num` is the minimum number of search result which you'd like to get at least.
It helps you to avoid getting the timeout error when the search process requires more time.
`min_num` should be positive integer and smaller than `num`.
