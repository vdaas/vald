// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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
	"github.com/vdaas/vald/internal/errors"
	client "github.com/vdaas/vald/internal/k8s/client"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/sync/errgroup"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

const (
	apiName = "vald/index/job/readreplica/rotate"
)

// Rotator represents an interface for indexing.
type Rotator interface {
	Start(ctx context.Context) error
}

type rotator struct {
	namespace           string
	volumeName          string
	readReplicaLabelKey string
	readReplicaID       string
	client              client.Client
	listOpts            client.ListOptions
}

// New returns Indexer object if no error occurs.
func New(replicaID string, opts ...Option) (Rotator, error) {
	r := new(rotator)

	if replicaID == "" {
		return nil, fmt.Errorf("readreplica id is empty. it should be set via MY_TARGET_REPLICA_ID env var")
	}
	r.readReplicaID = replicaID

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

	c, err := client.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}
	r.client = c

	selector, err := c.LabelSelector(r.readReplicaLabelKey, client.SelectionOpEquals, []string{r.readReplicaID})
	if err != nil {
		return nil, err
	}
	r.listOpts = client.ListOptions{
		Namespace:     r.namespace,
		LabelSelector: selector,
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
		if span != nil {
			span.RecordError(err)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return err
	}

	return nil
}

func (r *rotator) rotate(ctx context.Context) error {
	// get deployment here to pass to create methods of snapshot and pvc
	// and put it as owner reference of them so that they will be deleted when the deployment is deleted
	deployment, err := r.getDeployment(ctx)
	if err != nil {
		log.Errorf("failed to get Deployment.")
		return err
	}

	newSnap, oldSnap, err := r.createSnapshot(ctx, deployment)
	if err != nil {
		return err
	}

	newPvc, oldPvc, err := r.createPVC(ctx, newSnap.GetName(), deployment)
	if err != nil {
		log.Errorf("failed to create PVC. removing the new snapshot(%s)...", newSnap.GetName())
		if dserr := r.deleteSnapshot(ctx, newSnap); dserr != nil {
			errors.Join(err, dserr)
		}
		return err
	}

	err = r.updateDeployment(ctx, newPvc.GetName(), deployment)
	if err != nil {
		log.Errorf("failed to update Deployment. removing the new snapshot(%s) and pvc(%s)...", newSnap.GetName(), newPvc.GetName())
		if dperr := r.deletePVC(ctx, newPvc); dperr != nil {
			errors.Join(err, dperr)
		}
		if dserr := r.deleteSnapshot(ctx, newSnap); dserr != nil {
			errors.Join(err, dserr)
		}
		return err
	}

	err = r.deleteSnapshot(ctx, oldSnap)
	if err != nil {
		return err
	}

	err = r.deletePVC(ctx, oldPvc)
	if err != nil {
		return err
	}

	return nil
}

func (r *rotator) createSnapshot(ctx context.Context, deployment appsv1.Deployment) (newSnap, oldSnap *client.VolumeSnapshot, err error) {
	list := snapshotv1.VolumeSnapshotList{}
	if err := r.client.List(ctx, &list, &r.listOpts); err != nil {
		return nil, nil, fmt.Errorf("failed to get snapshot: %w", err)
	}
	if len(list.Items) == 0 {
		return nil, nil, fmt.Errorf("no snapshot found")
	}

	cur := &list.Items[0]
	oldSnap = cur.DeepCopy()
	newNameBase := getNewBaseName(cur.GetObjectMeta().GetName())
	if newNameBase == "" {
		return nil, nil, fmt.Errorf("the name(%s) doesn't seem to have replicaid", cur.GetObjectMeta().GetName())
	}
	newSnap = &client.VolumeSnapshot{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s%d", newNameBase, time.Now().Unix()),
			Namespace: cur.GetNamespace(),
			Labels:    cur.GetObjectMeta().GetLabels(),
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: "apps/v1",
					Kind:       "Deployment",
					Name:       deployment.GetName(),
					UID:        deployment.GetUID(),
					Controller: pointer.Bool(true),
				},
			},
		},
		Spec: cur.Spec,
	}

	log.Infof("creating new snapshot(%s)...", newSnap.GetName())
	log.Debugf("snapshot detail: %#v", newSnap)

	err = r.client.Create(ctx, newSnap)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create snapshot: %w", err)
	}

	return newSnap, oldSnap, nil
}

func (r *rotator) createPVC(ctx context.Context, newSnapShot string, deployment appsv1.Deployment) (newPvc, oldPvc *v1.PersistentVolumeClaim, err error) {
	list := v1.PersistentVolumeClaimList{}
	if err := r.client.List(ctx, &list, &r.listOpts); err != nil {
		return nil, nil, fmt.Errorf("failed to get PVC: %w", err)
	}
	if len(list.Items) == 0 {
		return nil, nil, fmt.Errorf("no PVC found")
	}

	cur := &list.Items[0]
	oldPvc = cur.DeepCopy()
	newNameBase := getNewBaseName(cur.GetObjectMeta().GetName())
	if newNameBase == "" {
		return nil, nil, fmt.Errorf("the name(%s) doesn't seem to have replicaid", cur.GetObjectMeta().GetName())
	}

	// remove timestamp from old pvc name
	newPvc = &v1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s%d", newNameBase, time.Now().Unix()),
			Namespace: cur.GetNamespace(),
			Labels:    cur.GetObjectMeta().GetLabels(),
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: "apps/v1",
					Kind:       "Deployment",
					Name:       deployment.GetName(),
					UID:        deployment.GetUID(),
					Controller: pointer.Bool(true),
				},
			},
		},
		Spec: v1.PersistentVolumeClaimSpec{
			AccessModes: cur.Spec.AccessModes,
			Resources:   cur.Spec.Resources,
			DataSource: &v1.TypedLocalObjectReference{
				Name:     newSnapShot,
				Kind:     cur.Spec.DataSource.Kind,
				APIGroup: cur.Spec.DataSource.APIGroup,
			},
		},
	}

	log.Infof("creating new pvc(%s)...", newPvc.GetName())
	log.Debugf("pvc detail: %#v", newPvc)

	if err := r.client.Create(ctx, newPvc); err != nil {
		return nil, nil, fmt.Errorf("failed to create PVC(%s): %w", newPvc.GetName(), err)
	}

	return newPvc, oldPvc, nil
}

func (r *rotator) getDeployment(ctx context.Context) (appsv1.Deployment, error) {
	list := appsv1.DeploymentList{}
	if err := r.client.List(ctx, &list, &r.listOpts); err != nil {
		return appsv1.Deployment{}, fmt.Errorf("failed to get deployment through client: %w", err)
	}
	if len(list.Items) == 0 {
		return appsv1.Deployment{}, fmt.Errorf("no deployment found")
	}

	return list.Items[0], nil
}

func (r *rotator) updateDeployment(ctx context.Context, newPVC string, deployment appsv1.Deployment) error {
	if deployment.Spec.Template.ObjectMeta.Annotations == nil {
		deployment.Spec.Template.ObjectMeta.Annotations = map[string]string{}
	}
	deployment.Spec.Template.ObjectMeta.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)

	for _, vol := range deployment.Spec.Template.Spec.Volumes {
		if vol.Name == r.volumeName {
			vol.PersistentVolumeClaim.ClaimName = newPVC
		}
	}

	log.Infof("updating deployment(%s)...", deployment.GetName())
	log.Debugf("deployment detail: %#v", deployment)

	if err := r.client.Update(ctx, &deployment); err != nil {
		return fmt.Errorf("failed to update deployment: %w", err)
	}

	return nil
}

func (r *rotator) deleteSnapshot(ctx context.Context, snapshot *snapshotv1.VolumeSnapshot) error {
	watcher, err := r.client.Watch(ctx,
		&snapshotv1.VolumeSnapshotList{
			Items: []snapshotv1.VolumeSnapshot{*snapshot},
		},
		&r.listOpts,
	)
	if err != nil {
		return fmt.Errorf("failed to watch snapshot(%s): %w", snapshot.GetName(), err)
	}
	defer watcher.Stop()

	eg, egctx := errgroup.New(ctx)
	eg.Go(func() error {
		log.Infof("deleting volume snapshot(%s)...", snapshot.GetName())
		log.Debugf("volume snapshot detail: %#v", snapshot)
		for {
			select {
			case <-egctx.Done():
				return egctx.Err()
			case event := <-watcher.ResultChan():
				if event.Type == client.WatchDeletedEvent {
					log.Infof("volume snapshot(%s) deleted", snapshot.GetName())
					return nil
				} else {
					log.Debugf("watching volume snapshot(%s) events. event: %v", snapshot.GetName(), event.Type)
				}
			}
		}
	})

	if err := r.client.Delete(ctx, snapshot); err != nil {
		return fmt.Errorf("failed to delete snapshot: %w", err)
	}
	return eg.Wait()
}

func (r *rotator) deletePVC(ctx context.Context, pvc *v1.PersistentVolumeClaim) error {
	watcher, err := r.client.Watch(ctx,
		&v1.PersistentVolumeClaimList{
			Items: []v1.PersistentVolumeClaim{*pvc},
		},
		&r.listOpts,
	)
	if err != nil {
		return fmt.Errorf("failed to watch PVC: %w", err)
	}
	defer watcher.Stop()

	eg, egctx := errgroup.New(ctx)
	eg.Go(func() error {
		log.Infof("deleting PVC(%s)...", pvc.GetName())
		log.Debugf("PVC detail: %#v", pvc)
		for {
			select {
			case <-egctx.Done():
				return egctx.Err()
			case event := <-watcher.ResultChan():
				if event.Type == client.WatchDeletedEvent {
					log.Infof("PVC(%s) deleted", pvc.GetName())
					return nil
				} else {
					log.Debugf("watching PVC(%s) events. event: %v", pvc.GetName(), event.Type)
				}
			}
		}
	})

	if err := r.client.Delete(ctx, pvc); err != nil {
		return fmt.Errorf("failed to delete PVC(%s): %w", pvc.GetName(), err)
	}

	return eg.Wait()
}

func getNewBaseName(old string) string {
	splits := strings.SplitAfter(old, "-")
	newNameBase := old + "-"
	// if this is not the first rotation, remove timestamp from the name
	// e.g. vald-agent-ngt-readreplica-0 -> the last element will be "0" which len is 1
	// so this means this is the first rotation
	if len(splits[len(splits)-1]) != 1 {
		newNameBase = strings.Join(splits[:len(splits)-1], "")
	}
	return newNameBase
}
