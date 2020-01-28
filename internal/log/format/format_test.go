package format

import "testing"

func TestString(t *testing.T) {
	type test struct {
		name   string
		format Format
		want   string
	}

	tests := []test{
		{
			name:   "returns raw",
			format: RAW,
			want:   "raw",
		},

		{
			name:   "returns json",
			format: JSON,
			want:   "json",
		},

		{
			name:   "returns ltsv",
			format: LTSV,
			want:   "ltsv",
		},

		{
			name:   "returns unknown",
			format: Format(100),
			want:   "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.format.String()
			if got != tt.want {
				t.Errorf("not equals. want: %v, but got: %v", tt.want, got)
			}
		})
	}
}

func TestAtof(t *testing.T) {
	type test struct {
		name string
		str  string
		want Format
	}

	tests := []test{
		{
			name: "returns RAW when str is raw",
			str:  "raw",
			want: RAW,
		},

		{
			name: "returns RAW when str is RAw",
			str:  "RAw",
			want: RAW,
		},

		{
			name: "returns JSON when str is json",
			str:  "json",
			want: JSON,
		},

		{
			name: "returns JSON when str is JSOn",
			str:  "JSOn",
			want: JSON,
		},

		{
			name: "returns LTSV when str is ltsv",
			str:  "ltsv",
			want: LTSV,
		},

		{
			name: "returns LTSV when str is LTSv",
			str:  "LTSv",
			want: LTSV,
		},

		{
			name: "returns Unknown when str is Vald",
			str:  "Vald",
			want: Unknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Atof(tt.str)
			if got != tt.want {
				t.Errorf("not equals. want: %v, but got: %v", tt.want, got)
			}
		})
	}
}
