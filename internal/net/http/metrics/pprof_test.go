package metrics

import "testing"

func TestNewPProfHandler(t *testing.T) {
	tests := []struct {
		name        string
		initialized bool
	}{
		{
			name:        "initialize success",
			initialized: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewPProfHandler()
			if (got != nil) != tt.initialized {
				t.Error("New() is wrong.")
			}
		})
	}
}
