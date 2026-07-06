package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJoinStr(t *testing.T) {
	tests := []struct {
		name string
		sep  string
		args []string
		want string
	}{
		{"slash separated", "/", []string{"a", "b", "c"}, "a/b/c"},
		{"dot separated", ".", []string{"ns", "name"}, "ns.name"},
		{"single element", "-", []string{"only"}, "only"},
		{"empty sep", "", []string{"foo", "bar"}, "foobar"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, JoinStr(tt.sep, tt.args...))
		})
	}
}
