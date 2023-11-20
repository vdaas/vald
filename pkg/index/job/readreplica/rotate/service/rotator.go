// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
	"fmt"
	"reflect"
	"strings"
	"time"

	snapshotv1 "github.com/kubernetes-csi/external-snapshotter/client/v6/apis/volumesnapshot/v1"
	snapshotclient "github.com/kubernetes-csi/external-snapshotter/client/v6/clientset/versioned"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/sync/errgroup"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

const (
	apiName = "vald/index/job/readreplica/rotate"
)

// Rotator represents an interface for indexing.
type Rotator interface {
	Start(ctx context.Context) error
}

type rotator struct {
	replicaid        int
	namespace        string
	deploymentPrefix string
	snapshotPrefix   string
	pvcPrefix        string
	volumeName       string
	// TODO: この辺はconbenchがマージされたあと、GetClientとかで引っ張ってくる
	clientset  *kubernetes.Clientset
	sClientset *snapshotclient.Clientset
}

// New returns Indexer object if no error occurs.
func New(clientset *kubernetes.Clientset, sClientset *snapshotclient.Clientset, opts ...Option) (Rotator, error) {
	r := new(rotator)

	if clientset == nil {
		return nil, fmt.Errorf("clientset is nil")
	}
	if sClientset == nil {
		return nil, fmt.Errorf("snapshot clientset is nil")
	}

	r.clientset = clientset
	r.sClientset = sClientset

	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(r); err != nil {
			oerr := errors.ErrOptionFailed(err, reflect.ValueOf(opt))
			e := &errors.ErrCriticalOption{}
			if errors.As(oerr, &e) {
				log.Error(err)
				return nil, oerr
			}
			log.Warn(oerr)
		}
	}
	return r, nil
}

// Start starts rotation process.
func (r *rotator) Start(ctx context.Context) error {
	_, span := trace.StartSpan(ctx, apiName+"/service/rotator.Start")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	if err := r.rotate(ctx); err != nil {
		// TODO: better error handling
		return err
	}

	return nil
}

func (r *rotator) rotate(ctx context.Context) error {
	newSnap, oldSnap, err := r.createSnapshot(ctx)
	if err != nil {
		return err
	}

	newPvc, oldPvc, err := r.createPVC(ctx, newSnap.Name)
	if err != nil {
		log.Infof("failed to create PVC. removing the new snapshot(%v)...", newSnap.Name)
		if dserr := r.deleteSnapshot(ctx, newSnap.Name); dserr != nil {
			errors.Join(err, dserr)
		}
		return err
	}

	err = r.updateDeployment(ctx, newPvc.Name)
	if err != nil {
		log.Infof("failed to update Deployment. removing the new snapshot(%v) and pvc(%v)...", newSnap.Name, newPvc.Name)
		if dperr := r.deletePVC(ctx, newPvc.Name); dperr != nil {
			errors.Join(err, dperr)
		}
		if dserr := r.deleteSnapshot(ctx, newSnap.Name); dserr != nil {
			errors.Join(err, dserr)
		}
		return err
	}

	err = r.deleteSnapshot(ctx, oldSnap.Name)
	if err != nil {
		return err
	}

	err = r.deletePVC(ctx, oldPvc.Name)
	if err != nil {
		return err
	}

	return nil
}

func (r *rotator) createSnapshot(ctx context.Context) (new, old *snapshotv1.VolumeSnapshot, err error) {
	snapshotInterface := r.sClientset.SnapshotV1().VolumeSnapshots(r.namespace)

	list, err := snapshotInterface.List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, nil, err
	}

	snapshotname := fmt.Sprintf("%s-%d", r.snapshotPrefix, r.replicaid)
	for _, snap := range list.Items {
		if strings.HasPrefix(snap.Name, snapshotname) {
			old = &snap
			break
		}
	}
	if old == nil {
		return nil, nil, fmt.Errorf("old snapshot not found")
	}

	new = &snapshotv1.VolumeSnapshot{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf(snapshotname+"-%d", time.Now().Unix()),
			Namespace: old.Namespace,
		},
		Spec: old.Spec,
	}

	new, err = snapshotInterface.Create(ctx, new, metav1.CreateOptions{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create snapshot: %w", err)
	}
	return new, old, nil
}

func (r *rotator) createPVC(ctx context.Context, newSnapShot string) (new, old *v1.PersistentVolumeClaim, err error) {
	pvcInterface := r.clientset.CoreV1().PersistentVolumeClaims(r.namespace)
	list, err := pvcInterface.List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, nil, err
	}

	pvcname := fmt.Sprintf("%s-%d", r.pvcPrefix, r.replicaid)
	for _, pvc := range list.Items {
		if strings.HasPrefix(pvc.Name, pvcname) {
			old = &pvc
			break
		}
	}
	if old == nil {
		return nil, nil, fmt.Errorf("old pvc not found")
	}

	new = &v1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf(pvcname+"-%d", time.Now().Unix()),
			Namespace: old.Namespace,
		},
		Spec: v1.PersistentVolumeClaimSpec{
			AccessModes: old.Spec.AccessModes,
			Resources:   old.Spec.Resources,
			DataSource: &v1.TypedLocalObjectReference{
				Name:     newSnapShot,
				Kind:     old.Spec.DataSource.Kind,
				APIGroup: old.Spec.DataSource.APIGroup,
			},
		},
	}

	new, err = pvcInterface.Create(ctx, new, metav1.CreateOptions{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create pvc: %w", err)
	}
	return new, old, nil
}

func (r *rotator) updateDeployment(ctx context.Context, newPVC string) error {
	deploymentInterface := r.clientset.AppsV1().Deployments(r.namespace)

	deploymentname := fmt.Sprintf("%s-%d", r.deploymentPrefix, r.replicaid)
	deployment, err := deploymentInterface.Get(ctx, deploymentname, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if deployment.Spec.Template.ObjectMeta.Annotations == nil {
		deployment.Spec.Template.ObjectMeta.Annotations = map[string]string{}
	}
	deployment.Spec.Template.ObjectMeta.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)

	for _, vol := range deployment.Spec.Template.Spec.Volumes {
		if vol.Name == r.volumeName {
			vol.PersistentVolumeClaim.ClaimName = newPVC
		}
	}

	_, err = deploymentInterface.Update(ctx, deployment, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (r *rotator) deleteSnapshot(ctx context.Context, snapshot string) error {
	snapshotInterface := r.sClientset.SnapshotV1().VolumeSnapshots(r.namespace)

	watcher, err := snapshotInterface.Watch(ctx, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to watch snapshot: %w", err)
	}
	defer watcher.Stop()

	eg, egctx := errgroup.New(ctx)
	eg.Go(func() error {
		log.Infof("deleting volume snapshot(%v)...", snapshot)
		for {
			select {
			case <-egctx.Done():
				return egctx.Err()
			case event := <-watcher.ResultChan():
				if event.Type == watch.Deleted {
					log.Infof("volume snapshot(%v) deleted", snapshot)
					return nil
				} else {
					log.Debugf("waching volume snapshot events. event: ", event.Type)
				}
			}
		}
	})

	err = snapshotInterface.Delete(ctx, snapshot, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete snapshot: %w", err)
	}

	return eg.Wait()
}

func (r *rotator) deletePVC(ctx context.Context, pvc string) error {
	pvcInterface := r.clientset.CoreV1().PersistentVolumeClaims(r.namespace)

	watcher, err := pvcInterface.Watch(ctx, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to watch PVC: %w", err)
	}
	defer watcher.Stop()

	eg, egctx := errgroup.New(ctx)
	eg.Go(func() error {
		log.Infof("deleting PVC(%v)...", pvc)
		for {
			select {
			case <-egctx.Done():
				return egctx.Err()
			case event := <-watcher.ResultChan():
				if event.Type == watch.Deleted {
					log.Infof("PVC(%v) deleted", pvc)
					return nil
				} else {
					log.Debugf("waching PVC events. event: %v", event.Type)
				}
			}
		}
	})

	err = pvcInterface.Delete(ctx, pvc, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete PVC: %w", err)
	}

	return eg.Wait()
}
