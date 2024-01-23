package pogreb

import (
	"context"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestNew(t *testing.T) {
	t.Parallel()
	type args struct {
		opts []Option
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		func() test {
			return test{
				name: "Succeeds to create pogreb instance with new path when the path does not exist",
				args: args{
					opts: []Option{
						WithPath(t.TempDir()),
						WithBackgroundSyncInterval("0s"),
					},
				},
			}
		}(),
		func() test {
			opts := []Option{
				WithPath(t.TempDir()),
				WithBackgroundSyncInterval("0s"),
			}
			return test{
				name: " Succeeds to restart the pogres instance",
				args: args{
					opts: opts,
				},
				beforeFunc: func(t *testing.T, _ args) {
					t.Helper()
					db, err := New(opts...)
					if err != nil {
						t.Fatal(err)
					}
					if err := db.Close(false); err != nil {
						t.Fatal(err)
					}
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			_, err := New(test.args.opts...)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOTE: To test Get, data needs to be inserted beforehand, so we test Set and Get together.
// If the test function name is changed, it is regenerated by gotests, so the function name is kept the same.
func Test_db_Get(t *testing.T) {
	t.Parallel()
	type args struct {
		opts []Option
		key  string
	}
	type want struct {
		want  []byte
		want1 bool
		err   error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, []byte, bool, error) error
		beforeFunc func(*testing.T, Pogreb, args)
		afterFunc  func(*testing.T, Pogreb, args)
	}
	defaultCheckFunc := func(w want, got []byte, got1 bool, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		if !reflect.DeepEqual(got1, w.want1) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got1, w.want1)
		}
		return nil
	}
	tests := []test{
		func() test {
			var (
				key = "key"
				val = []byte("val")
			)
			return test{
				name: "Succeeds to get the value associated with the key",
				args: args{
					opts: []Option{
						WithPath(t.TempDir()),
						WithBackgroundSyncInterval("0s"),
					},
					key: key,
				},
				want: want{
					want:  val,
					want1: true,
				},
				beforeFunc: func(t *testing.T, d Pogreb, args args) {
					t.Helper()
					if err := d.Set(key, val); err != nil {
						t.Fatal(err)
					}
				},
				afterFunc: func(t *testing.T, d Pogreb, args args) {
					t.Helper()
					if err := d.Close(true); err != nil {
						t.Fatal(err)
					}
				},
			}
		}(),
		func() test {
			var (
				key = "key"
				val = []byte("val")
			)
			return test{
				name: "Fails to get the value associated with the key if it does not exist",
				args: args{
					opts: []Option{
						WithPath(t.TempDir()),
						WithBackgroundSyncInterval("0s"),
					},
					key: "not-exist",
				},
				want: want{
					want1: false,
				},
				beforeFunc: func(t *testing.T, d Pogreb, args args) {
					t.Helper()
					if err := d.Set(key, val); err != nil {
						t.Fatal(err)
					}
				},
				afterFunc: func(t *testing.T, d Pogreb, args args) {
					t.Helper()
					if err := d.Close(true); err != nil {
						t.Fatal(err)
					}
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())

			d, err := New(test.args.opts...)
			if err != nil {
				t.Fatal(err)
			}
			if test.beforeFunc != nil {
				test.beforeFunc(tt, d, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, d, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got, got1, err := d.Get(test.args.key)
			if err := checkFunc(test.want, got, got1, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOTE: To test Delete, data needs to be inserted beforehand, so we test Set and Delete together.
// If the test function name is changed, it is regenerated by gotests, so the function name is kept the same.
func Test_db_Delete(t *testing.T) {
	t.Parallel()
	type args struct {
		opts []Option
		key  string
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Pogreb, error) error
		beforeFunc func(*testing.T, Pogreb, args)
		afterFunc  func(*testing.T, Pogreb, args)
	}
	defaultCheckFunc := func(w want, _ Pogreb, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		func() test {
			var (
				key = "key"
				val = []byte("val")
			)
			return test{
				name: "Succeeds to delete the value associated with the key",
				args: args{
					opts: []Option{
						WithPath(t.TempDir()),
						WithBackgroundSyncInterval("0s"),
					},
					key: key,
				},
				checkFunc: func(w want, d Pogreb, err error) error {
					if err := defaultCheckFunc(w, d, err); err != nil {
						return err
					}
					_, ok, err := d.Get(key)
					if err != nil {
						return err
					}
					if ok {
						return errors.New("key exists")
					}
					return nil
				},
				beforeFunc: func(t *testing.T, d Pogreb, args args) {
					t.Helper()
					if err := d.Set(key, val); err != nil {
						t.Fatal(err)
					}
				},
				afterFunc: func(t *testing.T, d Pogreb, args args) {
					t.Helper()
					if err := d.Close(true); err != nil {
						t.Fatal(err)
					}
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())

			d, err := New(test.args.opts...)
			if err != nil {
				t.Fatal(err)
			}
			if test.beforeFunc != nil {
				test.beforeFunc(tt, d, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, d, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			err = d.Delete(test.args.key)
			if err := checkFunc(test.want, d, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// skipcq: GO-R1005
func Test_db_Range(t *testing.T) {
	t.Parallel()
	type args struct {
		opts []Option
		ctx  context.Context
		f    func(key string, val []byte) bool
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(*testing.T, Pogreb, args)
		afterFunc  func(*testing.T, Pogreb, args)
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		func() test {
			data := map[string][]byte{
				"key-1": []byte("val-1"),
				"key-2": []byte("val-2"),
			}
			got := make(map[string][]byte)
			return test{
				name: "Succeeds to get all keys",
				args: args{
					opts: []Option{
						WithPath(t.TempDir()),
						WithBackgroundSyncInterval("0s"),
					},
					ctx: context.Background(),
					f: func(key string, val []byte) bool {
						got[key] = val
						return true
					},
				},
				checkFunc: func(w want, err error) error {
					if err := defaultCheckFunc(w, err); err != nil {
						return err
					}
					if !reflect.DeepEqual(got, data) {
						return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, data)
					}
					return nil
				},
				beforeFunc: func(t *testing.T, d Pogreb, args args) {
					t.Helper()
					for key, val := range data {
						if err := d.Set(key, val); err != nil {
							t.Fatal(err)
						}
					}
				},
				afterFunc: func(t *testing.T, d Pogreb, args args) {
					t.Helper()
					if err := d.Close(true); err != nil {
						t.Fatal(err)
					}
				},
			}
		}(),
		func() test {
			got := make(map[string][]byte)
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			return test{
				name: "Fails to get all keys when the context is canceled",
				args: args{
					opts: []Option{
						WithPath(t.TempDir()),
						WithBackgroundSyncInterval("0s"),
					},
					ctx: ctx,
					f: func(key string, val []byte) bool {
						got[key] = val
						return true
					},
				},
				checkFunc: func(w want, err error) error {
					if err := defaultCheckFunc(w, err); err != nil {
						return err
					}
					if data := make(map[string][]byte); !reflect.DeepEqual(got, data) {
						return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, data)
					}
					return nil
				},
				beforeFunc: func(t *testing.T, d Pogreb, args args) {
					t.Helper()
					data := map[string][]byte{
						"key-1": []byte("val-1"),
						"key-2": []byte("val-2"),
					}
					for key, val := range data {
						if err := d.Set(key, val); err != nil {
							t.Fatal(err)
						}
					}
				},
				afterFunc: func(t *testing.T, d Pogreb, args args) {
					t.Helper()
					if err := d.Close(true); err != nil {
						t.Fatal(err)
					}
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())

			d, err := New(test.args.opts...)
			if err != nil {
				t.Fatal(err)
			}
			if test.beforeFunc != nil {
				test.beforeFunc(tt, d, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, d, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			err = d.Range(test.args.ctx, test.args.f)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_db_Len(t *testing.T) {
	t.Parallel()
	type args struct {
		opts []Option
	}
	type want struct {
		want uint32
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, uint32) error
		beforeFunc func(*testing.T, Pogreb, args)
		afterFunc  func(*testing.T, Pogreb, args)
	}
	defaultCheckFunc := func(w want, got uint32) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			data := map[string][]byte{
				"key-1": []byte("val-1"),
				"key-2": []byte("val-2"),
			}
			return test{
				name: "Succeeds to get the number of keys",
				args: args{
					opts: []Option{
						WithPath(t.TempDir()),
						WithBackgroundSyncInterval("0s"),
					},
				},
				want: want{
					want: uint32(len(data)),
				},
				beforeFunc: func(t *testing.T, d Pogreb, args args) {
					t.Helper()
					for key, val := range data {
						if err := d.Set(key, val); err != nil {
							t.Fatal(err)
						}
					}
				},
				afterFunc: func(t *testing.T, d Pogreb, args args) {
					t.Helper()
					if err := d.Close(true); err != nil {
						t.Fatal(err)
					}
				},
			}
		}(),
		func() test {
			return test{
				name: "Succeeds to get the number of keys when key does not exist",
				args: args{
					opts: []Option{
						WithPath(t.TempDir()),
						WithBackgroundSyncInterval("0s"),
					},
				},
				afterFunc: func(t *testing.T, d Pogreb, args args) {
					t.Helper()
					if err := d.Close(true); err != nil {
						t.Fatal(err)
					}
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())

			d, err := New(test.args.opts...)
			if err != nil {
				t.Fatal(err)
			}
			if test.beforeFunc != nil {
				test.beforeFunc(tt, d, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, d, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := d.Len()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
//
// func Test_db_Set(t *testing.T) {
// 	type args struct {
// 		key string
// 		val []byte
// 	}
// 	type fields struct {
// 		db   *pogreb.DB
// 		opts *pogreb.Options
// 		path string
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           key:"",
// 		           val:nil,
// 		       },
// 		       fields: fields {
// 		           db:nil,
// 		           opts:nil,
// 		           path:"",
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           key:"",
// 		           val:nil,
// 		           },
// 		           fields: fields {
// 		           db:nil,
// 		           opts:nil,
// 		           path:"",
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			d := &db{
// 				db:   test.fields.db,
// 				opts: test.fields.opts,
// 				path: test.fields.path,
// 			}
//
// 			err := d.Set(test.args.key, test.args.val)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_db_Close(t *testing.T) {
// 	type args struct {
// 		remove bool
// 	}
// 	type fields struct {
// 		db   *pogreb.DB
// 		opts *pogreb.Options
// 		path string
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           remove:false,
// 		       },
// 		       fields: fields {
// 		           db:nil,
// 		           opts:nil,
// 		           path:"",
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           remove:false,
// 		           },
// 		           fields: fields {
// 		           db:nil,
// 		           opts:nil,
// 		           path:"",
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			d := &db{
// 				db:   test.fields.db,
// 				opts: test.fields.opts,
// 				path: test.fields.path,
// 			}
//
// 			err := d.Close(test.args.remove)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
