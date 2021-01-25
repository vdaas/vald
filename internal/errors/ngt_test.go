package errors

import (
	"math"
	"testing"
)

func TestErrCreateProperty(t *testing.T) {
	type args struct {
		err error
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return wrapped ErrCreateProperty error when err is ngt error",
			args: args{
				err: New("ngt error"),
			},
			want: want{
				want: New("failed to create property: ngt error"),
			},
		},
		{
			name: "return ErrCreateProperty error when err is nil",
			args: args{
				err: nil,
			},
			want: want{
				want: New("failed to create property"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := ErrCreateProperty(test.args.err)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrIndexNotFound(t *testing.T) {
	type want struct {
		want error
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return ErrIndexNotFound error",
			want: want{
				want: New("index file not found"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := ErrIndexNotFound
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrIndexLoadTimeout(t *testing.T) {
	type want struct {
		want error
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return ErrIndexLoadTimeout error",
			want: want{
				want: New("index load timeout"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := ErrIndexLoadTimeout
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrInvalidDimensionSize(t *testing.T) {

}

func TestErrDimensionLimitExceed(t *testing.T) {
	type args struct {
		current int
		limit   int
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return ErrDimensionLimitExceed error when current is 10 and limit is 5",
			args: args{
				current: 10,
				limit:   5,
			},
			want: want{
				want: New("supported dimension limit exceed:\trequired = 10,\tlimit = 5"),
			},
		},

		{
			name: "return ErrDimensionLimitExceed error when current is 0 and limit is 0",
			args: args{
				current: 0,
				limit:   0,
			},
			want: want{
				want: New("supported dimension limit exceed:\trequired = 0,\tlimit = 0"),
			},
		},
		{
			name: "return ErrDimensionLimitExceed error when current and limit are the minimum value of int",
			args: args{
				current: int(math.MinInt64),
				limit:   int(math.MinInt64),
			},
			want: want{
				want: Errorf("supported dimension limit exceed:\trequired = %d,\tlimit = %d", int(math.MinInt64), int(math.MinInt64)),
			},
		},
		{
			name: "return ErrDimensionLimitExceed error when current and limit are the maximum value of int",
			args: args{
				current: int(math.MaxInt64),
				limit:   int(math.MaxInt64),
			},
			want: want{
				want: Errorf("supported dimension limit exceed:\trequired = %d,\tlimit = %d", int(math.MaxInt64), int(math.MaxInt64)),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := ErrDimensionLimitExceed(test.args.current, test.args.limit)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrUnsupportedObjectType(t *testing.T) {
	type want struct {
		want error
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return ErrUnsupportedObjectType error",
			want: want{
				want: New("unsupported ObjectType"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := ErrUnsupportedObjectType
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrUnsupportedDistanceType(t *testing.T) {
	type want struct {
		want error
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return ErrUnsupportedDistanceType error",
			want: want{
				want: New("unsupported DistanceType"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := ErrUnsupportedDistanceType
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrFailedToSetDistanceType(t *testing.T) {
	type args struct {
		err      error
		distance string
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return wrapped ErrFailedToSetDistanceType error when err is ngt error and distance is 'l2'",
			args: args{
				err:      New("ngt error"),
				distance: "l2",
			},
			want: want{
				want: New("failed to set distance type l2: ngt error"),
			},
		},
		{
			name: "return wrapped ErrFailedToSetDistanceType error when err is ngt error and distance is empty",
			args: args{
				err:      New("ngt error"),
				distance: "",
			},
			want: want{
				want: New("failed to set distance type : ngt error"),
			},
		},
		{
			name: "return ErrFailedToSetDistanceType error when err is nil and distance is 'l2'",
			args: args{
				err:      nil,
				distance: "l2",
			},
			want: want{
				want: New("failed to set distance type l2"),
			},
		},
		{
			name: "return ErrFailedToSetDistanceType error when err is nil and distance is empty",
			args: args{
				err:      nil,
				distance: "",
			},
			want: want{
				want: New("failed to set distance type "),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := ErrFailedToSetDistanceType(test.args.err, test.args.distance)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrFailedToSetObjectType(t *testing.T) {
	type args struct {
		err error
		t   string
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return wrapped ErrFailedToSetObjectType error when err is ngt error and t is 'Float'",
			args: args{
				err: New("ngt error"),
				t:   "Float",
			},
			want: want{
				want: New("failed to set object type Float: ngt error"),
			},
		},
		{
			name: "return wrapped ErrFailedToSetObjectType error when err is ngt error and t is empty",
			args: args{
				err: New("ngt error"),
				t:   "",
			},
			want: want{
				want: New("failed to set object type : ngt error"),
			},
		},
		{
			name: "return ErrFailedToSetObjectType error when err is nil and t is 'Float'",
			args: args{
				err: nil,
				t:   "Float",
			},
			want: want{
				want: New("failed to set object type Float"),
			},
		},
		{
			name: "return ErrFailedToSetObjectType error when err is nil and t is empty",
			args: args{
				err: nil,
				t:   "",
			},
			want: want{
				want: New("failed to set object type "),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := ErrFailedToSetObjectType(test.args.err, test.args.t)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrFailedToSetDimension(t *testing.T) {
	type args struct {
		err error
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return wrapped ErrFailedToSetDimension error when err is ngt error",
			args: args{
				err: New("ngt error"),
			},
			want: want{
				want: New("failed to set dimension: ngt error"),
			},
		},
		{
			name: "return ErrFailedToSetDimension error when err is nil",
			args: args{
				err: nil,
			},
			want: want{
				want: New("failed to set dimension"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := ErrFailedToSetDimension(test.args.err)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrFailedToSetCreationEdgeSize(t *testing.T) {
	type args struct {
		err error
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return wrapped ErrFailedToSetCreationEdgeSize error when err is ngt error",
			args: args{
				err: New("ngt error"),
			},
			want: want{
				want: New("failed to set creation edge size: ngt error"),
			},
		},
		{
			name: "return ErrFailedToSetCreationEdgeSize error when err is nil",
			args: args{
				err: nil,
			},
			want: want{
				want: New("failed to set creation edge size"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := ErrFailedToSetCreationEdgeSize(test.args.err)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrFailedToSetSearchEdgeSize(t *testing.T) {
	type args struct {
		err error
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return wrapped ErrFailedToSetSearchEdgeSize error when err is ngt error",
			args: args{
				err: New("ngt error"),
			},
			want: want{
				want: New("failed to set search edge size: ngt error"),
			},
		},
		{
			name: "return ErrFailedToSetSearchEdgeSize error when err is nil",
			args: args{
				err: nil,
			},
			want: want{
				want: New("failed to set search edge size"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := ErrFailedToSetSearchEdgeSize(test.args.err)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrUncommittedIndexExists(t *testing.T) {
	type args struct {
		num uint64
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			return test{
				name: "return ErrUncommittedIndexExists error when num is 100",
				args: args{
					num: 100,
				},
				want: want{
					want: New("100 indexes are not committed"),
				},
			}
		}(),

		func() test {
			return test{
				name: "return ErrUncommittedIndexExists error when num is 0",
				args: args{
					num: 0,
				},
				want: want{
					want: New("0 indexes are not committed"),
				},
			}
		}(),
		func() test {
			var num uint64 = math.MaxUint64
			return test{
				name: "return ErrUncommittedIndexExists error when num is the maximum value of uint64",
				args: args{
					num: num,
				},
				want: want{
					want: Errorf("%d indexes are not committed", num),
				},
			}
		}(),
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := ErrUncommittedIndexExists(test.args.num)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrUncommittedIndexNotFound(t *testing.T) {
	type want struct {
		want error
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return ErrUncommittedIndexNotFound error",
			want: want{
				want: New("uncommitted indexes are not found"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := ErrUncommittedIndexNotFound
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrCAPINotImplemented(t *testing.T) {
	type want struct {
		want error
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return ErrCAPINotImplemented error",
			want: want{
				want: New("not implemented in C API"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := ErrCAPINotImplemented
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrUUIDAlreadyExists(t *testing.T) {
	type args struct {
		uuid string
		oid  uint
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return ErrUUIDAlreadyExists error when uuid is '550e8400-e29b-41d4' and oid is 100",
			args: args{
				uuid: "550e8400-e29b-41d4",
				oid:  100,
			},
			want: want{
				want: New("ngt uuid 550e8400-e29b-41d4 object id 100 already exists "),
			},
		},
		{
			name: "return ErrUUIDAlreadyExists error when uuid is empty and oid is 100",
			args: args{
				uuid: "",
				oid:  100,
			},
			want: want{
				want: New("ngt uuid  object id 100 already exists "),
			},
		},
		{
			name: "return ErrUUIDAlreadyExists error when uuid is '550e8400-e29b-41d4' and oid is the maximum value of uint64",
			args: args{
				uuid: "550e8400-e29b-41d4",
				oid:  uint(math.MaxUint64),
			},
			want: want{
				want: Errorf("ngt uuid 550e8400-e29b-41d4 object id %d already exists ", uint(math.MaxUint64)),
			},
		},
		{
			name: "return ErrUUIDAlreadyExists error when uuid is empty and oid is 0",
			args: args{
				uuid: "",
				oid:  0,
			},
			want: want{
				want: New("ngt uuid  object id 0 already exists "),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := ErrUUIDAlreadyExists(test.args.uuid, test.args.oid)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrUUIDNotFound(t *testing.T) {
	type args struct {
		id uint32
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return ErrUUIDNotFound error when id is 1234",
			args: args{
				id: 1234,
			},
			want: want{
				want: New("ngt object uuid 1234's metadata not found"),
			},
		},
		{
			name: "return ErrUUIDNotFound error when id is the maximum value of uint32",
			args: args{
				id: math.MaxUint32,
			},
			want: want{
				want: Errorf("ngt object uuid %d's metadata not found", math.MaxUint32),
			},
		},
		{
			name: "return ErrUUIDNotFound error when id is 0",
			args: args{
				id: 0,
			},
			want: want{
				want: New("ngt object uuid not found"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := ErrUUIDNotFound(test.args.id)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrObjectIDNotFound(t *testing.T) {
	type args struct {
		uuid string
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return ErrObjectIDNotFound error when uuid is '550e8400-e29b-41d4'",
			args: args{
				uuid: "550e8400-e29b-41d4",
			},
			want: want{
				want: New("ngt uuid 550e8400-e29b-41d4's object id not found"),
			},
		},
		{
			name: "return ErrObjectIDNotFound error when uuid is empty",
			args: args{
				uuid: "",
			},
			want: want{
				want: New("ngt uuid 's object id not found"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := ErrObjectIDNotFound(test.args.uuid)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrObjectNotFound(t *testing.T) {
	type args struct {
		err  error
		uuid string
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return wrapped ErrObjectNotFound error when err is ngt error and uuid is '550e8400-e29b-41d4'",
			args: args{
				err:  New("ngt error"),
				uuid: "550e8400-e29b-41d4",
			},
			want: want{
				want: New("ngt uuid 550e8400-e29b-41d4's object not found: ngt error"),
			},
		},
		{
			name: "return wrapped ErrObjectNotFound error when err is ngt error and uuid is empty",
			args: args{
				err:  New("ngt error"),
				uuid: "",
			},
			want: want{
				want: New("ngt uuid 's object not found: ngt error"),
			},
		},
		{
			name: "return ErrObjectNotFound error when err is nil and uuid is '550e8400-e29b-41d4'",
			args: args{
				err:  nil,
				uuid: "550e8400-e29b-41d4",
			},
			want: want{
				want: New("ngt uuid 550e8400-e29b-41d4's object not found"),
			},
		},
		{
			name: "return ErrObjectNotFound error when err is nil and uuid is empty",
			args: args{
				err:  nil,
				uuid: "",
			},
			want: want{
				want: New("ngt uuid 's object not found"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := ErrObjectNotFound(test.args.err, test.args.uuid)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrRemoveRequestedBeforeIndexing(t *testing.T) {
	type args struct {
		oid uint
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return ErrRemoveRequestedBeforeIndexing error when oid is 100",
			args: args{
				oid: 100,
			},
			want: want{
				want: New("object id 100 is not indexed we cannot remove it"),
			},
		},
		{
			name: "return ErrRemoveRequestedBeforeIndexing error when oid is 0",
			args: args{
				oid: 0,
			},
			want: want{
				want: New("object id 0 is not indexed we cannot remove it"),
			},
		},
		{
			name: "return ErrRemoveRequestedBeforeIndexing error when oid is maximum value of uint",
			args: args{
				oid: uint(math.MaxUint64),
			},
			want: want{
				want: Errorf("object id %d is not indexed we cannot remove it", uint(math.MaxUint64)),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := ErrRemoveRequestedBeforeIndexing(test.args.oid)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
