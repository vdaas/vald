# Remove Config

## Vald Remove Service

Vald provides the remove service for deleting indexed vector from the Vald cluster.
To remove the vector, it requires the vector's ID, and optional configuration.

For more details, please refer to [the Remove API document](../api/remove.md).

## Configuration

Here is the current insert config parameters.<BR>
All of parameters are optional.

```rpc
// Represent search configuration.
message Config {
  // Check the same set of vector and ID is already inserted or not.
  bool skip_strict_exist_check = 1;
  // The timestamp when the vector was inserted.
  int64 timestamp = 3;
}
```

### skip_strict_exist_check

`skip_strict_exist_check` is the flag for checking whether the same set of the vector and ID is already inserted or not.
If it is set as `true`, the checking function is available.<BR>
The default value is `false`.

### timestamp

`timestamp` is the timestamp when vector inserted.
When `timestamp` is not set, the current time will be used.

## Example Code

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
