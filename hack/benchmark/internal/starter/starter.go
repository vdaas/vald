package starter

import (
	"context"
	"testing"
)

type Starter interface {
	Run(context.Context, testing.TB) func()
}
