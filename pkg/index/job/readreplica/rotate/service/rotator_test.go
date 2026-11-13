// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package service

import (
	"context"
	"testing"
	"time"

	snapshotv1 "github.com/kubernetes-csi/external-snapshotter/client/v6/apis/volumesnapshot/v1"
	"github.com/stretchr/testify/require"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/client"
	"github.com/vdaas/vald/internal/k8s/resource"
	mock "github.com/vdaas/vald/internal/test/mock/k8s"
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

func Test_parseReplicaID(t *testing.T) {
	labelKey := "foo"
	type args struct {
		c         client.Client
		replicaID string
	}
	type want struct {
		err error
		ids []string
	}
	type test struct {
		args args
		name string
		want want
	}
	tests := []test{
		{
			name: "single replicaID",
			args: args{
				replicaID: "0",
				c:         nil,
			},
			want: want{
				ids: []string{"0"},
				err: nil,
			},
		},
		{
			name: "multiple replicaIDs",
			args: args{
				replicaID: "0,1",
				c:         nil,
			},
			want: want{
				ids: []string{"0", "1"},
				err: nil,
			},
		},
		{
			name: "returns error when replicaID is empty",
			args: args{
				replicaID: "",
				c:         nil,
			},
			want: want{
				ids: nil,
				err: errors.ErrReadReplicaIDEmpty,
			},
		},
		func() test {
			wantID1 := "bar"
			wantID2 := "baz"

			scheme := k8s.NewScheme()
			if err := k8s.AddClientGoScheme(scheme); err != nil {
				t.Fatalf("failed to add scheme: %v", err)
			}
			fakeClient := mock.NewFakeClientBuilder().WithScheme(scheme).WithObjects(
				&k8s.Deployment{Name: "dep1", Labels: map[string]string{labelKey: wantID1}},
				&k8s.Deployment{Name: "dep2", Labels: map[string]string{labelKey: wantID2}},
			).Build()

			mock := &mock.ValdK8sClientMock{}
			mock.WithRaw(fakeClient)
			return test{
				name: "returns all ids when rotate-all option is set",
				args: args{
					replicaID: rotateAllID,
					c:         mock,
				},
				want: want{
					ids: []string{wantID1, wantID2},
					err: nil,
				},
			}
		}(),
	}
	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := &rotator{
				readReplicaLabelKey: labelKey,
			}
			if tt.args.c != nil {
				r.deployments = resource.NewClientOf(tt.args.c, new(k8s.Deployment), new(k8s.DeploymentList))
			}
			ids, err := r.parseReplicaID(tt.args.replicaID)
			require.Equal(t, tt.want.ids, ids)
			require.Equal(t, tt.want.err, err)
		})
	}
}

func Test_rollbackContext(t *testing.T) {
	t.Parallel()
	type ctxKey struct{}
	parent, cancel := context.WithCancel(context.WithValue(t.Context(), ctxKey{}, "v"))
	cancel() // a sibling rotation failed and the errgroup cancelled the shared ctx

	before := time.Now()
	cctx, ccancel := rollbackContext(parent)
	defer ccancel()
	if err := cctx.Err(); err != nil {
		t.Fatalf("rollback ctx inherited the parent cancellation: %v", err)
	}
	dl, ok := cctx.Deadline()
	if !ok {
		t.Fatal("rollback ctx has no deadline")
	}
	if dl.Before(before.Add(rollbackTimeout)) || dl.After(time.Now().Add(rollbackTimeout)) {
		t.Errorf("deadline = %v, want about %v from now", dl, rollbackTimeout)
	}
	if got := cctx.Value(ctxKey{}); got != "v" {
		t.Errorf("rollback ctx value = %v, want %q", got, "v")
	}
}

func Test_subProcess_updateDeployment_volume(t *testing.T) {
	t.Parallel()
	const (
		labelKey = "vald-readreplica-id"
		volName  = "readreplica-pvc"
		newPVC   = "readreplica-pvc-new"
	)
	tests := []struct {
		name      string
		volumes   string
		wantClaim string
		wantErr   bool
	}{
		{
			name:      "PVC-backed volume is rewritten",
			volumes:   `[{"name":"` + volName + `","persistentVolumeClaim":{"claimName":"old"}}]`,
			wantClaim: newPVC,
		},
		{
			name:    "no volume with the configured name is an error",
			volumes: `[{"name":"other","persistentVolumeClaim":{"claimName":"old"}}]`,
			wantErr: true,
		},
		{
			name:    "volume with the configured name but no PVC source is an error",
			volumes: `[{"name":"` + volName + `","emptyDir":{}}]`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dep := new(k8s.Deployment)
			raw := `{"metadata":{"name":"dep","namespace":"default","labels":{"` + labelKey + `":"0"}},` +
				`"spec":{"template":{"spec":{"volumes":` + tt.volumes + `}}}}`
			if err := k8s.YAMLUnmarshal([]byte(raw), dep); err != nil {
				t.Fatalf("failed to build deployment: %v", err)
			}
			scheme := k8s.NewScheme()
			if err := k8s.AddClientGoScheme(scheme); err != nil {
				t.Fatalf("failed to add scheme: %v", err)
			}
			cli := &mock.ValdK8sClientMock{}
			cli.WithRaw(mock.NewFakeClientBuilder().WithScheme(scheme).WithObjects(dep).Build())
			s := &subProcess{
				deployments: resource.NewClientOf(cli, new(k8s.Deployment), new(k8s.DeploymentList)),
				volumeName:  volName,
				listOpts:    []k8s.ListOption{k8s.InNamespace("default"), k8s.MatchingLabels{labelKey: "0"}},
			}
			cur, err := s.getDeployment(t.Context())
			if err != nil {
				t.Fatalf("getDeployment: %v", err)
			}
			err = s.updateDeployment(t.Context(), newPVC, cur, time.Now())
			if tt.wantErr {
				if err == nil {
					t.Fatal("updateDeployment() error = nil, want error")
				}
				got, gerr := s.getDeployment(t.Context())
				if gerr != nil {
					t.Fatalf("getDeployment: %v", gerr)
				}
				if got.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] != "" {
					t.Error("deployment was updated despite the error")
				}
				return
			}
			if err != nil {
				t.Fatalf("updateDeployment() unexpected error: %v", err)
			}
			got, err := s.getDeployment(t.Context())
			if err != nil {
				t.Fatalf("getDeployment: %v", err)
			}
			if claim := got.Spec.Template.Spec.Volumes[0].PersistentVolumeClaim.ClaimName; claim != tt.wantClaim {
				t.Errorf("claimName = %q, want %q", claim, tt.wantClaim)
			}
		})
	}
}

func Test_withReplicaID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		old       string
		replicaID string
		want      string
	}{
		{
			name:      "first rotation of a single-digit replica keeps the id",
			old:       "vald-agent-ngt-readreplica-snapshot-0",
			replicaID: "0",
			want:      "vald-agent-ngt-readreplica-snapshot-0-",
		},
		{
			name:      "first rotation of a two-digit replica keeps the id",
			old:       "vald-agent-ngt-readreplica-snapshot-10",
			replicaID: "10",
			want:      "vald-agent-ngt-readreplica-snapshot-10-",
		},
		{
			name:      "later rotation of a two-digit replica strips only the timestamp",
			old:       "vald-agent-ngt-readreplica-pvc-10-1700000000",
			replicaID: "10",
			want:      "vald-agent-ngt-readreplica-pvc-10-",
		},
		{
			name:      "a name that already lost the replica id is repaired",
			old:       "vald-agent-ngt-readreplica-snapshot-1700000000",
			replicaID: "11",
			want:      "vald-agent-ngt-readreplica-snapshot-11-",
		},
		{
			name:      "a name without a replica id stays empty",
			old:       "vald",
			replicaID: "0",
			want:      "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := withReplicaID(getNewBaseName(tt.old), tt.replicaID); got != tt.want {
				t.Errorf("withReplicaID(getNewBaseName(%q), %q) = %q, want %q", tt.old, tt.replicaID, got, tt.want)
			}
		})
	}
}

// Test_subProcess_deleteOld pins that the post-switch cleanup survives the
// sibling cancellation of the shared ctx and always attempts both deletes.
func Test_subProcess_deleteOld(t *testing.T) {
	t.Parallel()
	const (
		ns       = "default"
		labelKey = "vald-readreplica-id"
	)
	tests := []struct {
		name     string
		withSnap bool
		wantErr  bool
	}{
		{
			name:     "both old resources are deleted although ctx was cancelled",
			withSnap: true,
		},
		{
			name:    "the old PVC is deleted although the snapshot delete fails",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			labels := map[string]string{labelKey: "0"}
			snap := &k8s.VolumeSnapshot{}
			snap.Name, snap.Namespace, snap.Labels = "snap-old", ns, labels
			pvc := &k8s.PersistentVolumeClaim{}
			pvc.Name, pvc.Namespace, pvc.Labels = "pvc-old", ns, labels

			scheme := k8s.NewScheme()
			if err := k8s.AddClientGoScheme(scheme); err != nil {
				t.Fatalf("failed to add scheme: %v", err)
			}
			if err := snapshotv1.AddToScheme(scheme); err != nil {
				t.Fatalf("failed to add snapshot scheme: %v", err)
			}
			objs := []k8s.Object{pvc.DeepCopy()}
			if tt.withSnap {
				objs = append(objs, snap.DeepCopy())
			}
			// The fake client is passed as is (not through ValdK8sClientMock) because
			// DeleteAndWait needs its Watch support.
			fc := mock.NewFakeClientBuilder().WithScheme(scheme).WithObjects(objs...).Build()
			s := &subProcess{
				snapshots: resource.NewClient(fc, new(k8s.VolumeSnapshot), new(k8s.VolumeSnapshotList)),
				pvcs:      resource.NewClient(fc, new(k8s.PersistentVolumeClaim), new(k8s.PersistentVolumeClaimList)),
				listOpts:  []k8s.ListOption{k8s.InNamespace(ns), k8s.MatchingLabels(labels)},
			}

			// a sibling rotation failed and the errgroup cancelled the shared ctx
			ctx, cancel := context.WithCancel(t.Context())
			cancel()
			done := make(chan error, 1)
			go func() { done <- s.deleteOld(ctx, snap, pvc) }()
			var err error
			select {
			case err = <-done:
			case <-time.After(30 * time.Second):
				t.Fatal("deleteOld did not return")
			}
			if tt.wantErr != (err != nil) {
				t.Errorf("deleteOld() error = %v, wantErr %v", err, tt.wantErr)
			}
			if _, gerr := s.pvcs.Get(t.Context(), pvc.Name, ns); !k8s.IsNotFound(gerr) {
				t.Errorf("old PVC still exists or lookup failed: %v", gerr)
			}
			if _, gerr := s.snapshots.Get(t.Context(), snap.Name, ns); !k8s.IsNotFound(gerr) {
				t.Errorf("old snapshot still exists or lookup failed: %v", gerr)
			}
		})
	}
}

// Test_subProcess_createPVC_mountedClaim pins that the PVC replaced (and later
// deleted) is the one the Deployment mounts, not just the first labeled PVC.
func Test_subProcess_createPVC_mountedClaim(t *testing.T) {
	t.Parallel()
	const (
		ns       = "default"
		labelKey = "vald-readreplica-id"
		volName  = "readreplica-pvc"
		stale    = "pvc-0-100"
		mounted  = "pvc-0-200"
	)
	pvcYAML := func(name string, terminating bool) string {
		meta := `"name":"` + name + `","namespace":"` + ns + `","labels":{"` + labelKey + `":"0"}`
		if terminating {
			meta += `,"deletionTimestamp":"2026-01-01T00:00:00Z","finalizers":["kubernetes.io/pvc-protection"]`
		}
		return `{"metadata":{` + meta + `},"spec":{"dataSource":{"kind":"VolumeSnapshot","name":"snap"}}}`
	}
	tests := []struct {
		name    string
		volumes string
		want    string
		pvcs    []string
	}{
		{
			name:    "the mounted PVC is chosen over an older labeled one",
			volumes: `[{"name":"` + volName + `","persistentVolumeClaim":{"claimName":"` + mounted + `"}}]`,
			pvcs:    []string{pvcYAML(stale, false), pvcYAML(mounted, false)},
			want:    mounted,
		},
		{
			name:    "a terminating PVC with the mounted name is skipped",
			volumes: `[{"name":"` + volName + `","persistentVolumeClaim":{"claimName":"` + mounted + `"}}]`,
			pvcs:    []string{pvcYAML(stale, false), pvcYAML(mounted, true)},
			want:    stale,
		},
		{
			name:    "without a mounted claim the first PVC is used",
			volumes: `[{"name":"other","persistentVolumeClaim":{"claimName":"` + mounted + `"}}]`,
			pvcs:    []string{pvcYAML(stale, false), pvcYAML(mounted, false)},
			want:    stale,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dep := new(k8s.Deployment)
			raw := `{"metadata":{"name":"dep","namespace":"` + ns + `"},"spec":{"template":{"spec":{"volumes":` + tt.volumes + `}}}}`
			if err := k8s.YAMLUnmarshal([]byte(raw), dep); err != nil {
				t.Fatalf("failed to build deployment: %v", err)
			}
			objs := make([]k8s.Object, 0, len(tt.pvcs))
			for _, p := range tt.pvcs {
				pvc := new(k8s.PersistentVolumeClaim)
				if err := k8s.YAMLUnmarshal([]byte(p), pvc); err != nil {
					t.Fatalf("failed to build pvc: %v", err)
				}
				objs = append(objs, pvc)
			}
			scheme := k8s.NewScheme()
			if err := k8s.AddClientGoScheme(scheme); err != nil {
				t.Fatalf("failed to add scheme: %v", err)
			}
			fc := mock.NewFakeClientBuilder().WithScheme(scheme).WithObjects(objs...).Build()
			s := &subProcess{
				pvcs:       resource.NewClient(fc, new(k8s.PersistentVolumeClaim), new(k8s.PersistentVolumeClaimList)),
				listOpts:   []k8s.ListOption{k8s.InNamespace(ns), k8s.MatchingLabels{labelKey: "0"}},
				replicaID:  "0",
				volumeName: volName,
			}
			newPvc, oldPvc, err := s.createPVC(t.Context(), "snap-new", dep)
			if err != nil {
				t.Fatalf("createPVC() unexpected error: %v", err)
			}
			if got := oldPvc.GetName(); got != tt.want {
				t.Errorf("old PVC = %q, want %q", got, tt.want)
			}
			if got := newPvc.Spec.DataSource.Name; got != "snap-new" {
				t.Errorf("new PVC data source = %q, want %q", got, "snap-new")
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
// 		})
// 	}
// }
//
// func Test_rotator_Start(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		snapshots           *resource.Client[*k8s.VolumeSnapshot, *k8s.VolumeSnapshotList]
// 		pvcs                *resource.Client[*k8s.PersistentVolumeClaim, *k8s.PersistentVolumeClaimList]
// 		deployments         *resource.Client[*k8s.Deployment, *k8s.DeploymentList]
// 		namespace           string
// 		volumeName          string
// 		readReplicaLabelKey string
// 		subProcesses        []subProcess
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
// 		           snapshots:nil,
// 		           pvcs:nil,
// 		           deployments:nil,
// 		           namespace:"",
// 		           volumeName:"",
// 		           readReplicaLabelKey:"",
// 		           subProcesses:nil,
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
// 		           snapshots:nil,
// 		           pvcs:nil,
// 		           deployments:nil,
// 		           namespace:"",
// 		           volumeName:"",
// 		           readReplicaLabelKey:"",
// 		           subProcesses:nil,
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
// 				snapshots:           test.fields.snapshots,
// 				pvcs:                test.fields.pvcs,
// 				deployments:         test.fields.deployments,
// 				namespace:           test.fields.namespace,
// 				volumeName:          test.fields.volumeName,
// 				readReplicaLabelKey: test.fields.readReplicaLabelKey,
// 				subProcesses:        test.fields.subProcesses,
// 			}
//
// 			err := r.Start(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_rotator_newSubprocess(t *testing.T) {
// 	type args struct {
// 		replicaID string
// 	}
// 	type fields struct {
// 		snapshots           *resource.Client[*k8s.VolumeSnapshot, *k8s.VolumeSnapshotList]
// 		pvcs                *resource.Client[*k8s.PersistentVolumeClaim, *k8s.PersistentVolumeClaimList]
// 		deployments         *resource.Client[*k8s.Deployment, *k8s.DeploymentList]
// 		namespace           string
// 		volumeName          string
// 		readReplicaLabelKey string
// 		subProcesses        []subProcess
// 	}
// 	type want struct {
// 		want subProcess
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, subProcess) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got subProcess) error {
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
// 		       },
// 		       fields: fields {
// 		           snapshots:nil,
// 		           pvcs:nil,
// 		           deployments:nil,
// 		           namespace:"",
// 		           volumeName:"",
// 		           readReplicaLabelKey:"",
// 		           subProcesses:nil,
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
// 		           },
// 		           fields: fields {
// 		           snapshots:nil,
// 		           pvcs:nil,
// 		           deployments:nil,
// 		           namespace:"",
// 		           volumeName:"",
// 		           readReplicaLabelKey:"",
// 		           subProcesses:nil,
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
// 				snapshots:           test.fields.snapshots,
// 				pvcs:                test.fields.pvcs,
// 				deployments:         test.fields.deployments,
// 				namespace:           test.fields.namespace,
// 				volumeName:          test.fields.volumeName,
// 				readReplicaLabelKey: test.fields.readReplicaLabelKey,
// 				subProcesses:        test.fields.subProcesses,
// 			}
//
// 			got := r.newSubprocess(test.args.replicaID)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_subProcess_rotate(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		snapshots   *resource.Client[*k8s.VolumeSnapshot, *k8s.VolumeSnapshotList]
// 		pvcs        *resource.Client[*k8s.PersistentVolumeClaim, *k8s.PersistentVolumeClaimList]
// 		deployments *resource.Client[*k8s.Deployment, *k8s.DeploymentList]
// 		selector    string
// 		volumeName  string
// 		listOpts    []k8s.ListOption
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
// 		           snapshots:nil,
// 		           pvcs:nil,
// 		           deployments:nil,
// 		           selector:"",
// 		           volumeName:"",
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
// 		           snapshots:nil,
// 		           pvcs:nil,
// 		           deployments:nil,
// 		           selector:"",
// 		           volumeName:"",
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
// 			s := &subProcess{
// 				snapshots:   test.fields.snapshots,
// 				pvcs:        test.fields.pvcs,
// 				deployments: test.fields.deployments,
// 				selector:    test.fields.selector,
// 				volumeName:  test.fields.volumeName,
// 				listOpts:    test.fields.listOpts,
// 			}
//
// 			err := s.rotate(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_subProcess_createSnapshot(t *testing.T) {
// 	type args struct {
// 		ctx        context.Context
// 		deployment *k8s.Deployment
// 	}
// 	type fields struct {
// 		snapshots   *resource.Client[*k8s.VolumeSnapshot, *k8s.VolumeSnapshotList]
// 		pvcs        *resource.Client[*k8s.PersistentVolumeClaim, *k8s.PersistentVolumeClaimList]
// 		deployments *resource.Client[*k8s.Deployment, *k8s.DeploymentList]
// 		selector    string
// 		volumeName  string
// 		listOpts    []k8s.ListOption
// 	}
// 	type want struct {
// 		wantNewSnap *k8s.VolumeSnapshot
// 		wantOldSnap *k8s.VolumeSnapshot
// 		err         error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *k8s.VolumeSnapshot, *k8s.VolumeSnapshot, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotNewSnap *k8s.VolumeSnapshot, gotOldSnap *k8s.VolumeSnapshot, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotNewSnap, w.wantNewSnap) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotNewSnap, w.wantNewSnap)
// 		}
// 		if !reflect.DeepEqual(gotOldSnap, w.wantOldSnap) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotOldSnap, w.wantOldSnap)
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
// 		           deployment:nil,
// 		       },
// 		       fields: fields {
// 		           snapshots:nil,
// 		           pvcs:nil,
// 		           deployments:nil,
// 		           selector:"",
// 		           volumeName:"",
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
// 		           deployment:nil,
// 		           },
// 		           fields: fields {
// 		           snapshots:nil,
// 		           pvcs:nil,
// 		           deployments:nil,
// 		           selector:"",
// 		           volumeName:"",
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
// 			s := &subProcess{
// 				snapshots:   test.fields.snapshots,
// 				pvcs:        test.fields.pvcs,
// 				deployments: test.fields.deployments,
// 				selector:    test.fields.selector,
// 				volumeName:  test.fields.volumeName,
// 				listOpts:    test.fields.listOpts,
// 			}
//
// 			gotNewSnap, gotOldSnap, err := s.createSnapshot(test.args.ctx, test.args.deployment)
// 			if err := checkFunc(test.want, gotNewSnap, gotOldSnap, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_subProcess_createPVC(t *testing.T) {
// 	type args struct {
// 		ctx         context.Context
// 		newSnapShot string
// 		deployment  *k8s.Deployment
// 	}
// 	type fields struct {
// 		snapshots   *resource.Client[*k8s.VolumeSnapshot, *k8s.VolumeSnapshotList]
// 		pvcs        *resource.Client[*k8s.PersistentVolumeClaim, *k8s.PersistentVolumeClaimList]
// 		deployments *resource.Client[*k8s.Deployment, *k8s.DeploymentList]
// 		selector    string
// 		volumeName  string
// 		listOpts    []k8s.ListOption
// 	}
// 	type want struct {
// 		wantNewPvc *k8s.PersistentVolumeClaim
// 		wantOldPvc *k8s.PersistentVolumeClaim
// 		err        error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *k8s.PersistentVolumeClaim, *k8s.PersistentVolumeClaim, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotNewPvc *k8s.PersistentVolumeClaim, gotOldPvc *k8s.PersistentVolumeClaim, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotNewPvc, w.wantNewPvc) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotNewPvc, w.wantNewPvc)
// 		}
// 		if !reflect.DeepEqual(gotOldPvc, w.wantOldPvc) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotOldPvc, w.wantOldPvc)
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
// 		           newSnapShot:"",
// 		           deployment:nil,
// 		       },
// 		       fields: fields {
// 		           snapshots:nil,
// 		           pvcs:nil,
// 		           deployments:nil,
// 		           selector:"",
// 		           volumeName:"",
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
// 		           newSnapShot:"",
// 		           deployment:nil,
// 		           },
// 		           fields: fields {
// 		           snapshots:nil,
// 		           pvcs:nil,
// 		           deployments:nil,
// 		           selector:"",
// 		           volumeName:"",
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
// 			s := &subProcess{
// 				snapshots:   test.fields.snapshots,
// 				pvcs:        test.fields.pvcs,
// 				deployments: test.fields.deployments,
// 				selector:    test.fields.selector,
// 				volumeName:  test.fields.volumeName,
// 				listOpts:    test.fields.listOpts,
// 			}
//
// 			gotNewPvc, gotOldPvc, err := s.createPVC(test.args.ctx, test.args.newSnapShot, test.args.deployment)
// 			if err := checkFunc(test.want, gotNewPvc, gotOldPvc, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_subProcess_getDeployment(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		snapshots   *resource.Client[*k8s.VolumeSnapshot, *k8s.VolumeSnapshotList]
// 		pvcs        *resource.Client[*k8s.PersistentVolumeClaim, *k8s.PersistentVolumeClaimList]
// 		deployments *resource.Client[*k8s.Deployment, *k8s.DeploymentList]
// 		selector    string
// 		volumeName  string
// 		listOpts    []k8s.ListOption
// 	}
// 	type want struct {
// 		want *k8s.Deployment
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *k8s.Deployment, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got *k8s.Deployment, err error) error {
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
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           snapshots:nil,
// 		           pvcs:nil,
// 		           deployments:nil,
// 		           selector:"",
// 		           volumeName:"",
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
// 		           snapshots:nil,
// 		           pvcs:nil,
// 		           deployments:nil,
// 		           selector:"",
// 		           volumeName:"",
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
// 			s := &subProcess{
// 				snapshots:   test.fields.snapshots,
// 				pvcs:        test.fields.pvcs,
// 				deployments: test.fields.deployments,
// 				selector:    test.fields.selector,
// 				volumeName:  test.fields.volumeName,
// 				listOpts:    test.fields.listOpts,
// 			}
//
// 			got, err := s.getDeployment(test.args.ctx)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_subProcess_updateDeployment(t *testing.T) {
// 	type args struct {
// 		ctx          context.Context
// 		newPVC       string
// 		deployment   *k8s.Deployment
// 		snapshotTime time.Time
// 	}
// 	type fields struct {
// 		snapshots   *resource.Client[*k8s.VolumeSnapshot, *k8s.VolumeSnapshotList]
// 		pvcs        *resource.Client[*k8s.PersistentVolumeClaim, *k8s.PersistentVolumeClaimList]
// 		deployments *resource.Client[*k8s.Deployment, *k8s.DeploymentList]
// 		selector    string
// 		volumeName  string
// 		listOpts    []k8s.ListOption
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
// 		           newPVC:"",
// 		           deployment:nil,
// 		           snapshotTime:time.Time{},
// 		       },
// 		       fields: fields {
// 		           snapshots:nil,
// 		           pvcs:nil,
// 		           deployments:nil,
// 		           selector:"",
// 		           volumeName:"",
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
// 		           newPVC:"",
// 		           deployment:nil,
// 		           snapshotTime:time.Time{},
// 		           },
// 		           fields: fields {
// 		           snapshots:nil,
// 		           pvcs:nil,
// 		           deployments:nil,
// 		           selector:"",
// 		           volumeName:"",
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
// 			s := &subProcess{
// 				snapshots:   test.fields.snapshots,
// 				pvcs:        test.fields.pvcs,
// 				deployments: test.fields.deployments,
// 				selector:    test.fields.selector,
// 				volumeName:  test.fields.volumeName,
// 				listOpts:    test.fields.listOpts,
// 			}
//
// 			err := s.updateDeployment(test.args.ctx, test.args.newPVC, test.args.deployment, test.args.snapshotTime)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_rotator_parseReplicaID(t *testing.T) {
// 	type args struct {
// 		replicaID string
// 	}
// 	type fields struct {
// 		snapshots           *resource.Client[*k8s.VolumeSnapshot, *k8s.VolumeSnapshotList]
// 		pvcs                *resource.Client[*k8s.PersistentVolumeClaim, *k8s.PersistentVolumeClaimList]
// 		deployments         *resource.Client[*k8s.Deployment, *k8s.DeploymentList]
// 		namespace           string
// 		volumeName          string
// 		readReplicaLabelKey string
// 		subProcesses        []subProcess
// 	}
// 	type want struct {
// 		want []string
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, []string, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got []string, err error) error {
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
// 		       },
// 		       fields: fields {
// 		           snapshots:nil,
// 		           pvcs:nil,
// 		           deployments:nil,
// 		           namespace:"",
// 		           volumeName:"",
// 		           readReplicaLabelKey:"",
// 		           subProcesses:nil,
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
// 		           },
// 		           fields: fields {
// 		           snapshots:nil,
// 		           pvcs:nil,
// 		           deployments:nil,
// 		           namespace:"",
// 		           volumeName:"",
// 		           readReplicaLabelKey:"",
// 		           subProcesses:nil,
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
// 				snapshots:           test.fields.snapshots,
// 				pvcs:                test.fields.pvcs,
// 				deployments:         test.fields.deployments,
// 				namespace:           test.fields.namespace,
// 				volumeName:          test.fields.volumeName,
// 				readReplicaLabelKey: test.fields.readReplicaLabelKey,
// 				subProcesses:        test.fields.subProcesses,
// 			}
//
// 			got, err := r.parseReplicaID(test.args.replicaID)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
