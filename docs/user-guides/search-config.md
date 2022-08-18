# Search Config

## Vald Search Service

<!-- Describe vald search service: search, searchbyid, linearsearch, lineasearchbyid ... -->

Vald provides four types of search services.

1. Search

   - `Search` is the `ANN(Approximate Nearest Neighbor)` search with query vector.
     It is a fast search even though the vector consists large dimension.
     The search duration is quick but less accurate than `LinearSearch`.
     The search algorithm depends on each core algorithm.

1. SearchById

   - `SearchById` is the `ANN(Approximate Nearest Neighbor)` search with the stored vector's id.
     The id should already exist in the NGT indexes before the search process.
     The search algorithm is the same as `Search`.

1. LinearSearch

   - `LinearSearch` is the primary search algorithm with a query vector.
     It searches all indexed vectors and calculates the distance between the query.
     Its accuracy is exact, but the search time requires more than `Search` (ANN search) and increases the amount of indexed vector.

1. LinearSearchById
   - `LinearSearchById` is the primary search algorithm with the vector's id.
     The id should already exist in the NGT indexes before the search process.
     The search algorithm is the same as `LinearSearch`.

<div class="notice">
Linear Search service is available from Vald v1.4 or later.
</div>

For more details, please refer to [the Search API document](../api/search.md).

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

`request_id` is a unique request ID.
It is **NOT** indexed vector's id.
Users can use it for, e.g., the error handling process.

### num

`num` is the maximum number of search results you'd like to get.
`num` should be a positive integer.

### radius

`radius`, the specific parameter for NGT, specifies the search range centered on the query vector in terms of the radius of a sphere.
The number of search target vectors increases along with the radius is large.
There is a trade-off between accuracy and search speed.
It is hard to set it depending on the dataset in many cases.

The default value is infinity.
When setting a negative number as `radius`, `NGT` applies the radius as infinity.

<div class="notice">
NGT will self-update the radius during the search process.
</div>

### epsilon

`epsilon`, the specific parameter for NGT, specifies the search range's magnification coefficient (epsilon).
NGT will use `radius*(1+epsilon)` as the search range.
The number of search target vectors increases along with the epsilon being large.

The default value is 0.1, and it may work in most cases.
However, the appropriate value may vary depending on the dataset.
While it is desirable to adjust this value within 0 - 0.3, it can also set a negative value (over than -1).

### ingress_filters

`ingress_filters` is required when using the ingress filter component.
It requires the ingress filter component's hostname and port.

### egress_filters

`egress_filters` is required when using the egress filter component.
It requires the egress filter component's hostname and port.

### min_num

`min_num` is the minimum number of search results you'd like to get at least.
It helps you avoid the timeout error when the search process requires more time.
`min_num` should be a positive integer and smaller than `num`.
