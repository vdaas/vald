package timeutil

import (
	"reflect"
	"testing"
	"time"
)

func BenchmarkParse(b *testing.B) {
	type args struct {
		t string
	}
	type want struct {
		want time.Duration
		err  error
	}
	type test struct {
		name string
		args args
		want want
	}

	tests := []test{
		{
			name: "when t is 10ms",
			args: args{
				t: "10ms",
			},
			want: want{
				want: 10 * time.Millisecond,
				err:  nil,
			},
		},
		{
			name: "when t is 100ms",
			args: args{
				t: "100ms",
			},
			want: want{
				want: 100 * time.Millisecond,
				err:  nil,
			},
		},
		{
			name: "when t is 1s",
			args: args{
				t: "1s",
			},
			want: want{
				want: time.Second,
				err:  nil,
			},
		},
		{
			name: "when t is 10s",
			args: args{
				t: "10s",
			},
			want: want{
				want: 10 * time.Second,
				err:  nil,
			},
		},
		{
			name: "when t is 100s",
			args: args{
				t: "100s",
			},
			want: want{
				want: 100 * time.Second,
				err:  nil,
			},
		},
	}

	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					got, err := Parse(test.args.t)
					if err != nil {
						b.Error(err)
					}
					if !reflect.DeepEqual(test.want.want, got) {
						b.Errorf("want: %v, but got: %v", test.want.want, got)
					}
				}
			})
		})
	}
}

func BenchmarkParseWithDefault(b *testing.B) {
	type args struct {
		t string
		d time.Duration
	}
	type want struct {
		want time.Duration
	}
	type test struct {
		name string
		args args
		want want
	}

	tests := []test{
		{
			name: "when t is 10second",
			args: args{
				t: "10second",
				d: 50 * time.Millisecond,
			},
			want: want{
				want: 50 * time.Millisecond,
			},
		},
		{
			name: "when t is 100second",
			args: args{
				t: "100second",
				d: 50 * time.Millisecond,
			},
			want: want{
				want: 50 * time.Millisecond,
			},
		},
		{
			name: "when t is 1000second",
			args: args{
				t: "1000second",
				d: 50 * time.Millisecond,
			},
			want: want{
				want: 50 * time.Millisecond,
			},
		},
		{
			name: "when t is 10000second",
			args: args{
				t: "1000second",
				d: 50 * time.Millisecond,
			},
			want: want{
				want: 50 * time.Millisecond,
			},
		},
		{
			name: "when t is 100000second",
			args: args{
				t: "10000second",
				d: 50 * time.Millisecond,
			},
			want: want{
				want: 50 * time.Millisecond,
			},
		},
	}

	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					got := ParseWithDefault(test.args.t, test.args.d)
					if !reflect.DeepEqual(test.want.want, got) {
						b.Errorf("want: %v, but got: %v", test.want.want, got)
					}
				}
			})
		})
	}
}
