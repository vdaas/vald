package grpc

import (
	"context"
	"fmt"
	"testing"

	"github.com/vdaas/vald/internal/errors"
)

func TestHoge(t *testing.T) {
	err := context.Canceled
	// err = errors.Wrap(context.Canceled, "fdsfasf")
	fmt.Println(errors.Is(errors.ErrRPCCallFailed("aaa", context.Canceled), err))
	fmt.Println(errors.Is(err, errors.ErrRPCCallFailed("aaa", context.Canceled)))
}
