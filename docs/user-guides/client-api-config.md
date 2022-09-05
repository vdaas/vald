# Client API Config

Vald provides client libraries for indexing vectors, searching approximate nearest neighbor vectors, and removing indexing vectors.
Each request allows setting configuration for user demands.
Please select the section which you need.

Please also see each API document.

## Insert Service

`Insert` is inserting new vectors into the Vald cluster.
It requires the vector, its ID (specific ID for the vector), and optional configuration.

### Configuration

```rpc
// Represent search configuration.
message Config {
  // Check the same set of vector and ID is already inserted or not.
  bool skip_strict_exist_check = 1;
  // Configuration for filter if your Vald cluster uses filter.
  Filter.Config filters = 2;
  // The timestamp when the vector was inserted.
  int64 timestamp = 3;
}
```

<details><summary>Insert Configuration Sample (Go)</summary><br>

```go
package main

import (
	"context"
	"time"

	"github.com/vdaas/vald-client-go/v1/payload"
	"github.com/vdaas/vald-client-go/v1/vald"
	"google.golang.org/grpc"
)

func main() {
	// Create connection
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	target := "localhost:8080"
	conn, err := grpc.DialContext(ctx, target)
	if err != nil {
		panic(err)
	}
	// Init vald client
	client := vald.NewValdClient(conn)

	// Insert sample
	location, err := client.Insert(ctx, &payload.Insert_Request{
		// Vector information (mandatory)
		Vector: &payload.Object_Vector{
			// Set the vector
			Vector: []float32{0, 1, 2},
			// Set the specific ID for the vector
                    // The ID must be not indexed.
			Id:     "sample",
		},
		// Insert configuration (optional)
		Config: &payload.Insert_Config{
			SkipStrictExistCheck: true,
			Filters: &payload.Filter_Config{
				Targets: []*payload.Filter_Target{
					{
						Host: "vald-ingress-filter",
						Port: 8081,
					},
				},
			},
			Timestamp: time.Now().UnixMilli(),
		},
	})
	if err != nil {
		panic(err)
	}
	...
}
```

</details>

#### skip_strict_exist_check

`skip_strict_exist_check` is the flag for checking whether the same set of the vector and ID is already inserted or not.
If it is set as `true`, the checking function is available.<BR>
The default value is `false`.

#### filters

`filters` is the configuration when using filter functions.
In the `Insert` section, it is popular for using ingress filtering.

The detail configuration is following.

```rpc
// Filter related messages.
message Filter {

  // Represent the target filter server.
  message Target {
    // The target filter component hostname.
    string host = 1;
    // The target filter component port.
    uint32 port = 2;
  }

  // Represent filter configuration.
  message Config {
    // Represent the filter target configuration.
    repeated Target targets = 1;
  }
}
```

#### timestamp

`timestamp` is the timestamp when vector inserted.
When `timestamp` is not set, the current time will be used.

## Update Service

`Update` is updating vectors which are already inserted in the `vald-agent` component.
It requires the new vector, its ID (the target ID already indexed), and optional configuration.

### Configuration

```rpc
// Represent search configuration.
message Config {
  // Check the same set of vector and ID is already inserted or not.
  bool skip_strict_exist_check = 1;
  // Configuration for filter if your Vald cluster uses filter.
  Filter.Config filters = 2;
  // The timestamp when the vector was inserted.
  int64 timestamp = 3;
}
```

<details><summary>Update Configuration Sample (Go)</summary><br>

```go
package main

import (
	"context"
	"time"

	"github.com/vdaas/vald-client-go/v1/payload"
	"github.com/vdaas/vald-client-go/v1/vald"
	"google.golang.org/grpc"
)

func example() {
	// Create connection
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	target := "localhost:8080"
	conn, err := grpc.DialContext(ctx, target)
	if err != nil {
		panic(err)
	}
	// Init vald client
	client := vald.NewValdClient(conn)

	// Update sample
	location, err := client.Update(ctx, &payload.Update_Request{
		// Vector information (mandatory)
		Vector: &payload.Object_Vector{
			// Set the vector
			Vector: []float32{0, 1, 2},
			// Set the specific ID for the vector
                        // The ID must be already indexed.
			Id:     "sample",
		},
		// Insert configuration (optional)
		Config: &payload.Update_Config{
			SkipStrictExistCheck: true,
			Filters: &payload.Filter_Config{
				Targets: []*payload.Filter_Target{
					{
						Host: "vald-ingress-filter",
						Port: 8081,
					},
				},
			},
			Timestamp: time.Now().UnixMilli(),
		},
	})
	if err != nil {
		panic(err)
	}
	...
}
```

</details>

#### skip_strict_exist_check

`skip_strict_exist_check` is the flag for checking whether the same set of the vector and ID is already inserted or not.
If it is set as `true`, the checking function is available.<BR>
The default value is `false`.

#### filters

`filters` is the configuration when using filter functions.
In the `Insert` section, it is popular for using ingress filtering.

The detail configuration is following.

```rpc
// Filter related messages.
message Filter {

  // Represent the target filter server.
  message Target {
    // The target filter component hostname.
    string host = 1;
    // The target filter component port.
    uint32 port = 2;
  }

  // Represent filter configuration.
  message Config {
    // Represent the filter target configuration.
    repeated Target targets = 1;
  }
}
```

#### timestamp

`timestamp` is the timestamp when vector inserted.
When `timestamp` is not set, the current time will be used.

## Upsert Service

`Upsert` is updating existing vectors in the `vald-agent` or inserting new vectors into the `vald-agent` if the vector does not exist.
It requires the vector, its ID (specific ID for the vector), and optional configuration.

### Configuration

```rpc
// Represent search configuration.
message Config {
  // Check the same set of vector and ID is already inserted or not.
  bool skip_strict_exist_check = 1;
  // Configuration for filter if your Vald cluster uses filter.
  Filter.Config filters = 2;
  // The timestamp when the vector was inserted.
  int64 timestamp = 3;
}
```

<details><summary>Upsert Configuration Sample (Go)</summary><br>

```go
package main

import (
	"context"
	"time"

	"github.com/vdaas/vald-client-go/v1/payload"
	"github.com/vdaas/vald-client-go/v1/vald"
	"google.golang.org/grpc"
)

func example() {
	// Create connection
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	target := "localhost:8080"
	conn, err := grpc.DialContext(ctx, target)
	if err != nil {
		panic(err)
	}
	// Init vald client
	client := vald.NewValdClient(conn)

	// Update sample
	location, err := client.Upsert(ctx, &payload.Upsert_Request{
		// Vector information (mandatory)
		Vector: &payload.Object_Vector{
			// Set the vector
			Vector: []float32{0, 1, 2},
			// Set the specific ID for the vector
			Id:     "sample",
		},
		// Insert configuration (optional)
		Config: &payload.Upsert_Config{
			SkipStrictExistCheck: true,
			Filters: &payload.Filter_Config{
				Targets: []*payload.Filter_Target{
					{
						Host: "vald-ingress-filter",
						Port: 8081,
					},
				},
			},
			Timestamp: time.Now().UnixMilli(),
		},
	})
	if err != nil {
		panic(err)
	}
        ...
}
```

</details>


## Search Service

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

### Configuration

```rpc
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

<details><summary>Search Configuration Sample (Go)</summary><br>

```go
package main

import (
	"context"
	"time"

	"github.com/vdaas/vald-client-go/v1/payload"
	"github.com/vdaas/vald-client-go/v1/vald"
	"google.golang.org/grpc"
)

func main() {
	// Create connection
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	target := "localhost:8080"
	conn, err := grpc.DialContext(ctx, target)
	if err != nil {
		panic(err)
	}
	// Init vald client
	client := vald.NewValdClient(conn)

	// Search sample
	res, err := client.Search(ctx, &payload.Search_Request{
		// Query vector
		Vector: []float32{0, 1, 2},
		// Search configuration
		Config: &payload.Search_Config{
			RequestId: "unique-request-id",
			// The number of search result to be returned.
			Num: 10,
			// The minimum number of search result to be returned.
			// It prevent the timeout error when the number of result does NOT satisfy Num.
			MinNum: 5,
			// The space of search candidate redius for NN vectors.
			Radius: -1,
			// Epsilon is used to determine how much to expand from search candidate radius.
			Epsilon: 0.1,
			// Search timeout setting.
			Timeout: 100000000,
			IngressFilters: &payload.Filter_Config{
				Targets: []*payload.Filter_Target{
					{
						Host: "vald-ingress-filter",
						Port: 8081,
					},
				},
			},
			EgressFilters: &payload.Filter_Config{
				Targets: []*payload.Filter_Target{
					{
						Host: "vald-egress-filter",
						Port: 8081,
					},
				},
			},
		},
	})

	if err != nil {
		panic(err)
	}
	...
}
```

</details>

#### request_id

`request_id` is a unique request ID.
It is **NOT** indexed vector's id.
Users can use it for, e.g., the error handling process.

#### num

`num` is the maximum number of search results you'd like to get.
`num` should be a positive integer.

#### radius

`radius`, the specific parameter for NGT, specifies the search range centered on the query vector in terms of the radius of a sphere.
The number of search target vectors increases along with the radius is large.
There is a trade-off between accuracy and search speed.
It is hard to set it depending on the dataset in many cases.

The default value is infinity.
When setting a negative number as `radius`, `NGT` applies the radius as infinity.

<div class="notice">
NGT will self-update the radius during the search process.
</div>

#### epsilon

`epsilon`, the specific parameter for NGT, specifies the search range's magnification coefficient (epsilon).
NGT will use `radius*(1+epsilon)` as the search range.
The number of search target vectors increases along with the epsilon being large.

The default value is 0.1, and it may work in most cases.
However, the appropriate value may vary depending on the dataset.
While it is desirable to adjust this value within 0 - 0.3, it can also set a negative value (over than -1).

#### ingress_filters

`ingress_filters` is required when using the ingress filter component.
It requires the ingress filter component's hostname and port.

```rpc
// Filter related messages.
message Filter {

  // Represent the target filter server.
  message Target {
    // The target filter component hostname.
    string host = 1;
    // The target filter component port.
    uint32 port = 2;
  }

  // Represent filter configuration.
  message Config {
    // Represent the filter target configuration.
    repeated Target targets = 1;
  }
}
```

#### egress_filters

`egress_filters` is required when using the egress filter component.
It requires the egress filter component's hostname and port.

```rpc
// Filter related messages.
message Filter {

  // Represent the target filter server.
  message Target {
    // The target filter component hostname.
    string host = 1;
    // The target filter component port.
    uint32 port = 2;
  }

  // Represent filter configuration.
  message Config {
    // Represent the filter target configuration.
    repeated Target targets = 1;
  }
}
```

#### min_num

`min_num` is the minimum number of search results you'd like to get at least.
It helps you avoid the timeout error when the search process requires more time.
`min_num` should be a positive integer and smaller than `num`.

## Remove Service

`Remove` is deleting indexed vector from the Vald cluster.
To remove the vector, it requires the vector's ID, and optional configuration.

For more details, please refer to [the Remove API document](../api/remove.md).

### Configuration

```rpc
// Represent search configuration.
message Config {
  // Check the same set of vector and ID is already inserted or not.
  bool skip_strict_exist_check = 1;
  // The timestamp when the vector was inserted.
  int64 timestamp = 3;
}
```

<details><summary>Remove Configuration Sample (Go)</summary><br>

```go
package main

import (
	"context"
	"time"

	"github.com/vdaas/vald-client-go/v1/payload"
	"github.com/vdaas/vald-client-go/v1/vald"
	"google.golang.org/grpc"
)

func example() {
	// Create connection
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	target := "localhost:8080"
	conn, err := grpc.DialContext(ctx, target)
	if err != nil {
		panic(err)
	}
	// Init vald client
	client := vald.NewValdClient(conn)

	// Remove sample
	location, err := client.Remove(ctx, &payload.Remove_Request{
		// Vector ID (mandatory)
                // The ID must be already indexed.
		Id: &payload.Object_ID{
			Id: "sample",
		},
		// Insert configuration (optional)
		Config: &payload.Remove_Config{
			SkipStrictExistCheck: true,
			Timestamp: time.Now().UnixMilli(),
		},
	})
	if err != nil {
		panic(err)
	}
        ...
}
```

</details>

### skip_strict_exist_check

`skip_strict_exist_check` is the flag for checking whether the same set of the vector and ID is already inserted or not.
If it is set as `true`, the checking function is available.<BR>
The default value is `false`.

### timestamp

`timestamp` is the timestamp when vector inserted.
When `timestamp` is not set, the current time will be used.
