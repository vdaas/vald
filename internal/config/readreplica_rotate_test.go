package config

// NOT IMPLEMENTED BELOW
//
// func TestReadreplicaRotate_Bind(t *testing.T) {
// 	type fields struct {
// 		AgentNamespace      string
// 		ReadReplicaLabelKey string
// 		ReadReplicaID       string
// 		VolumeName          string
// 	}
// 	type want struct {
// 		want *ReadreplicaRotate
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *ReadreplicaRotate) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got *ReadreplicaRotate) error {
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
// 		       fields: fields {
// 		           AgentNamespace:"",
// 		           ReadReplicaLabelKey:"",
// 		           ReadReplicaID:"",
// 		           VolumeName:"",
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
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
// 		           fields: fields {
// 		           AgentNamespace:"",
// 		           ReadReplicaLabelKey:"",
// 		           ReadReplicaID:"",
// 		           VolumeName:"",
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
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
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			r := &ReadreplicaRotate{
// 				AgentNamespace:      test.fields.AgentNamespace,
// 				ReadReplicaLabelKey: test.fields.ReadReplicaLabelKey,
// 				ReadReplicaID:       test.fields.ReadReplicaID,
// 				VolumeName:          test.fields.VolumeName,
// 			}
//
// 			got := r.Bind()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
