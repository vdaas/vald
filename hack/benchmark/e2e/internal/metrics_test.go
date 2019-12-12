package internal

import "testing"

func TestRecall(t *testing.T) {
	x := []string{
		"foo", "bar", "hoge", "huga",
	}
	y := []string{
		"foo", "baz", "hoge", "huga", "moge",
	}

	tests := []struct{
		k int
		recall float64
	}{
		{1,1.0},
		{2,0.5},
		{3,2.0 / 3.0},
		{4,3.0 / 4.0},
	}

	for _, tt := range tests {
		if Recall(x, y, tt.k) != tt.recall {
			t.Errorf("Recall@%d is %f", tt.k, tt.recall)
		}
	}
}
