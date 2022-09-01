# Insert / Update / Upsert Config

## Vald Insert / Update / Upsert Service

Vald provides the insert service, the update service, and the upsert service for indexing vector on the Vald cluster.

1. Insert

   - `Insert` is inserting new vectors into the Vald cluster.
     It requires the vector, its ID (specific ID for the vector), and optional configuration.

1. Update

   - `Update` is updating vectors which are already inserted in the `vald-agent` component.
     It requires the new vector, its ID (the target ID already indexed), and optional configuration.

1. Upsert

   - `Upsert` is updating existing vectors in the `vald-agent` or inserting new vectors into the `vald-agent` if the vector does not exist.
     It requires the vector, its ID (specific ID for the vector), and optional configuration.

For more details, please refer to each documents:

- [the Insert API document](../api/insert.md).
- [the Update API document](../api/update.md).
- [the Upsert API document](../api/upsert.md).

## Configuration

Here is the current insert config parameters.

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

- Example Code

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

### skip_strict_exist_check

`skip_strict_exist_check` is the flag for checking whether the same set of the vector and ID is already inserted or not.
If it is set as `true`, the checking function is available.<BR>
The default value is `false`.

### filters

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

### timestamp

`timestamp` is the timestamp when vector inserted.
When `timestamp` is not set, the current time will be used.
