package glg

import "testing"

func Test_level_String(t *testing.T) {
	type test struct {
		name  string
		level level
		want  string
	}

	tests := []test{
		{
			name:  "returns Debug",
			level: DEBUG,
			want:  "Debug",
		},

		{
			name:  "returns Info",
			level: INFO,
			want:  "Info",
		},

		{
			name:  "returns Warn",
			level: WARN,
			want:  "Warn",
		},

		{
			name:  "returns Error",
			level: ERROR,
			want:  "Error",
		},

		{
			name:  "returns Fatal",
			level: FATAL,
			want:  "Fatal",
		},

		{
			name:  "returns Unknown",
			level: level(100),
			want:  "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.level.String()
			if got != tt.want {
				t.Errorf("not equals. want: %v, but got: %v", tt.want, got)
			}
		})
	}
}

func Test_toLevel(t *testing.T) {
	type test struct {
		name string
		lv   string
		want level
	}

	tests := []test{
		{
			name: "returns DEBUG",
			lv:   "Debug",
			want: DEBUG,
		},

		{
			name: "returns INFO",
			lv:   "Info",
			want: INFO,
		},

		{
			name: "returns WARN",
			lv:   "Warn",
			want: WARN,
		},

		{
			name: "returns ERROR",
			lv:   "Error",
			want: ERROR,
		},

		{
			name: "returns FATAL",
			lv:   "Fatal",
			want: FATAL,
		},

		{
			name: "returns INFO",
			lv:   "Unknown",
			want: INFO,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toLevel(tt.lv)
			if got != tt.want {
				t.Errorf("not equals. want: %v, but got: %v", tt.want, got)
			}
		})
	}
}
