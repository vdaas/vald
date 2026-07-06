package lifecycle

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestLifeCycles() LifeCycles {
	return LifeCycles{
		{Condition: Condition{Type: ConditionWaitForClusterCreate}},
		{Condition: Condition{Type: ConditionWaitForCreateVrs}},
		{Condition: Condition{Type: ConditionCompleted}},
	}
}

func TestLifeCycles_GetIndex(t *testing.T) {
	lcs := newTestLifeCycles()

	tests := []struct {
		name string
		ct   string
		want int
	}{
		{"empty string returns 0", "", 0},
		{"first condition", ConditionWaitForClusterCreate, 0},
		{"middle condition", ConditionWaitForCreateVrs, 1},
		{"last condition", ConditionCompleted, 2},
		{"not found returns -1", "NonExistent", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, lcs.GetIndex(tt.ct))
		})
	}
}

func TestFlow_GetNext(t *testing.T) {
	lcs := newTestLifeCycles()

	tests := []struct {
		name    string
		current int
		want    int
	}{
		{"from first", 0, 1},
		{"from middle", 1, 2},
		{"at last returns -1", 2, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFlow(lcs, tt.current)
			assert.Equal(t, tt.want, f.GetNext())
		})
	}
}
