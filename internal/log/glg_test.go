package log

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/kpango/glg"
)

func TestNewGlg(t *testing.T) {
	type test struct {
		name string
		log  *glg.Glg
		want *glglogger
	}

	tests := []test{
		func() test {
			log := glg.Get()

			return test{
				name: "initialize success",
				log:  log,
				want: &glglogger{
					log: log,
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewGlg(tt.log)
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("not equals. want: %v, got: %v", tt.want, got)
			}
		})
	}
}

func TestDefaultGlg(t *testing.T) {
	type test struct {
		name     string
		want     *glglogger
		recovery bool
	}

	tests := []test{
		{
			name: "default glg success",
			want: &glglogger{
				log: glg.Get(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DefaultGlg()
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("not equals. want: %v, got: %v", tt.want, got)
			}
		})
	}
}

func TestOut(t *testing.T) {
	type args struct {
		fn   func(...interface{}) error
		vals []interface{}
	}

	type test struct {
		name      string
		args      args
		checkFunc func() error
		recovery  bool
	}

	tests := []test{
		func() test {
			vals := []interface{}{
				"name",
			}

			fn := func(vals ...interface{}) error {
				return errors.New("aaaa")
			}

			return test{
				name: "output success when fn return nil",
				args: args{
					vals: vals,
					fn:   fn,
				},
				checkFunc: func() error {
					return nil
				},
				recovery: true,
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println(err)
				}
			}()
			out(tt.args.fn, tt.args.vals...)
		})
	}
}

func TestOutf(t *testing.T) {
	type args struct {
		fn     func(string, ...interface{}) error
		format string
		vals   []interface{}
	}

	type test struct {
		name     string
		args     args
		recovery bool
	}

	tests := []test{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outf(tt.args.fn, tt.args.format, tt.args.vals...)
		})
	}
}
