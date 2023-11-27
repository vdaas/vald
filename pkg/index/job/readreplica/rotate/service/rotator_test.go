package service

import (
	"testing"
)

func Test_getNewBaseName(t *testing.T) {
	type args struct {
		old string
	}
	type want struct {
		want string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "fist rotation just returns the name",
			args: args{
				old: "vald-agent-ngt-readreplica-0",
			},
			want: want{
				want: "vald-agent-ngt-readreplica-0-",
			},
		},
		{
			name: "successfully remove timestamp",
			args: args{
				old: "vald-agent-ngt-readreplica-0-20220101",
			},
			want: want{
				want: "vald-agent-ngt-readreplica-0-",
			},
		},
		{
			name: "successfully remove timestamp when the name has no dashes",
			args: args{
				old: "vald-1-20220101",
			},
			want: want{
				want: "vald-1-",
			},
		},
		{
			name: "no replica id basename returns empty string",
			args: args{
				old: "vald",
			},
			want: want{
				want: "",
			},
		},
	}
	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := getNewBaseName(tt.args.old); got != tt.want.want {
				t.Errorf("getNewBaseName() = %v, want %v", got, tt.want.want)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestNew(t *testing.T) {
// 	type args struct {
// 		replicaID string
// 		opts      []Option
// 	}
// 	type want struct {
// 		want Rotator
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Rotator, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Rotator, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           replicaID:"",
// 		           opts:nil,
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
// 		           replicaID:"",
// 		           opts:nil,
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
//
// 			got, err := New(test.args.replicaID, test.args.opts...)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_rotator_Start(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		namespace           string
// 		volumeName          string
// 		readReplicaLabelKey string
// 		readReplicaID       string
// 		client              client.Client
// 		listOpts            client.ListOptions
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
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           namespace:"",
// 		           volumeName:"",
// 		           readReplicaLabelKey:"",
// 		           readReplicaID:"",
// 		           client:nil,
// 		           listOpts:nil,
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
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           namespace:"",
// 		           volumeName:"",
// 		           readReplicaLabelKey:"",
// 		           readReplicaID:"",
// 		           client:nil,
// 		           listOpts:nil,
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
// 			r := &rotator{
// 				namespace:           test.fields.namespace,
// 				volumeName:          test.fields.volumeName,
// 				readReplicaLabelKey: test.fields.readReplicaLabelKey,
// 				readReplicaID:       test.fields.readReplicaID,
// 				client:              test.fields.client,
// 				listOpts:            test.fields.listOpts,
// 			}
//
// 			err := r.Start(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
