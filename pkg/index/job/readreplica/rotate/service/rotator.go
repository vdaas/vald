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
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/client"
	"github.com/vdaas/vald/internal/k8s/vald"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"k8s.io/utils/ptr"
)

const (
	apiName     = "vald/index/job/readreplica/rotate"
	rotateAllID = "rotate-all"
)

// Rotator represents an interface for indexing.
type Rotator interface {
	Start(ctx context.Context) error
}

type rotator struct {
	namespace           string
	volumeName          string
	readReplicaLabelKey string
	subProcesses        []subProcess
}

type subProcess struct {
	listOpts   k8s.ListOptions
	client     client.Client
	volumeName string
}

// New returns Indexer object if no error occurs.
// replicaID must be a comma separated string of replica id or ${rotateAllID} to rotate all read replica at once.
func New(replicaID string, opts ...Option) (Rotator, error) {
	r := new(rotator)

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

	ids, err := r.parseReplicaID(replicaID, c)
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		sub, err := r.newSubprocess(c, id)
		if err != nil {
			return nil, fmt.Errorf("failed to create rotator subprocess: %w", err)
		}

		r.subProcesses = append(r.subProcesses, sub)
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

	eg, ectx := errgroup.New(ctx)
	for _, sub := range r.subProcesses {
		s := sub
		eg.Go(safety.RecoverFunc(func() (err error) {
			if err := s.rotate(ectx); err != nil {
				if span != nil {
					span.RecordError(err)
					span.SetStatus(trace.StatusError, err.Error())
				}
				return err
			}
			return nil
		}))
	}

	return eg.Wait()
}

func (r *rotator) newSubprocess(c client.Client, replicaID string) (subProcess, error) {
	selector, err := c.LabelSelector(r.readReplicaLabelKey, k8s.SelectionOpEquals, []string{replicaID})
	if err != nil {
		return subProcess{}, err
	}
	sub := subProcess{
		client: c,
		listOpts: k8s.ListOptions{
			Namespace:     r.namespace,
			LabelSelector: selector,
		},
		volumeName: r.volumeName,
	}
	return sub, nil
}

func (s *subProcess) rotate(ctx context.Context) error {
	// get deployment here to pass to create methods of snapshot and pvc
	// and put it as owner reference of them so that they will be deleted when the deployment is deleted
	deployment, err := s.getDeployment(ctx)
	if err != nil {
		log.Errorf("failed to get Deployment.")
		return err
	}

	newSnap, oldSnap, err := s.createSnapshot(ctx, deployment)
	if err != nil {
		return err
	}

	newPvc, oldPvc, err := s.createPVC(ctx, newSnap.GetName(), deployment)
	if err != nil {
		log.Errorf("failed to create PVC. removing the new snapshot(%s)...", newSnap.GetName())
		if dserr := s.deleteSnapshot(ctx, newSnap); dserr != nil {
			errors.Join(err, dserr)
		}
		return err
	}

	err = s.updateDeployment(ctx, newPvc.GetName(), deployment, newSnap.CreationTimestamp.Time)
	if err != nil {
		log.Errorf("failed to update Deployment. removing the new snapshot(%s) and pvc(%s)...", newSnap.GetName(), newPvc.GetName())
		if dperr := s.deletePVC(ctx, newPvc); dperr != nil {
			errors.Join(err, dperr)
		}
		if dserr := s.deleteSnapshot(ctx, newSnap); dserr != nil {
			errors.Join(err, dserr)
		}
		return err
	}

	err = s.deleteSnapshot(ctx, oldSnap)
	if err != nil {
		return err
	}

	err = s.deletePVC(ctx, oldPvc)
	if err != nil {
		return err
	}

	return nil
}

func (s *subProcess) createSnapshot(
	ctx context.Context, deployment *k8s.Deployment,
) (newSnap, oldSnap *k8s.VolumeSnapshot, err error) {
	list := k8s.VolumeSnapshotList{}
	if err := s.client.List(ctx, &list, &s.listOpts); err != nil {
		return nil, nil, fmt.Errorf("failed to get snapshot: %w", err)
	}
	if len(list.Items) == 0 {
		return nil, nil, fmt.Errorf("no snapshot found")
	}

	cur := &list.Items[0]
	oldSnap = cur.DeepCopy()
	newNameBase := getNewBaseName(cur.GetObjectMeta().GetName())
	if newNameBase == "" {
		return nil, nil, fmt.Errorf("the name(%s) doesn't seem to have replica id", cur.GetObjectMeta().GetName())
	}
	newSnap = &k8s.VolumeSnapshot{
		ObjectMeta: k8s.ObjectMeta{
			Name:      fmt.Sprintf("%s%d", newNameBase, time.Now().Unix()),
			Namespace: cur.GetNamespace(),
			Labels:    cur.GetObjectMeta().GetLabels(),
			OwnerReferences: []k8s.OwnerReference{
				{
					APIVersion: "apps/v1",
					Kind:       "Deployment",
					Name:       deployment.GetName(),
					UID:        deployment.GetUID(),
					Controller: ptr.To(true),
				},
			},
		},
		Spec: cur.Spec,
	}

	log.Infof("creating new snapshot(%s)...", newSnap.GetName())
	log.Debugf("snapshot detail: %#v", newSnap)

	err = s.client.Create(ctx, newSnap)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create snapshot: %w", err)
	}

	return newSnap, oldSnap, nil
}

func (s *subProcess) createPVC(
	ctx context.Context, newSnapShot string, deployment *k8s.Deployment,
) (newPvc, oldPvc *k8s.PersistentVolumeClaim, err error) {
	list := k8s.PersistentVolumeClaimList{}
	if err := s.client.List(ctx, &list, &s.listOpts); err != nil {
		return nil, nil, fmt.Errorf("failed to get PVC: %w", err)
	}
	if len(list.Items) == 0 {
		return nil, nil, fmt.Errorf("no PVC found")
	}

	cur := &list.Items[0]
	oldPvc = cur.DeepCopy()
	newNameBase := getNewBaseName(cur.GetObjectMeta().GetName())
	if newNameBase == "" {
		return nil, nil, fmt.Errorf("the name(%s) doesn't seem to have replica id", cur.GetObjectMeta().GetName())
	}

	// remove timestamp from old pvc name
	newPvc = &k8s.PersistentVolumeClaim{
		ObjectMeta: k8s.ObjectMeta{
			Name:      fmt.Sprintf("%s%d", newNameBase, time.Now().Unix()),
			Namespace: cur.GetNamespace(),
			Labels:    cur.GetObjectMeta().GetLabels(),
			OwnerReferences: []k8s.OwnerReference{
				{
					APIVersion: "apps/v1",
					Kind:       "Deployment",
					Name:       deployment.GetName(),
					UID:        deployment.GetUID(),
					Controller: ptr.To(true),
				},
			},
		},
		Spec: k8s.PersistentVolumeClaimSpec{
			AccessModes: cur.Spec.AccessModes,
			Resources:   cur.Spec.Resources,
			DataSource: &k8s.TypedLocalObjectReference{
				Name:     newSnapShot,
				Kind:     cur.Spec.DataSource.Kind,
				APIGroup: cur.Spec.DataSource.APIGroup,
			},
			StorageClassName: cur.Spec.StorageClassName,
		},
	}

	log.Infof("creating new pvc(%s)...", newPvc.GetName())
	log.Debugf("pvc detail: %#v", newPvc)

	if err := s.client.Create(ctx, newPvc); err != nil {
		return nil, nil, fmt.Errorf("failed to create PVC(%s): %w", newPvc.GetName(), err)
	}

	return newPvc, oldPvc, nil
}

func (s *subProcess) getDeployment(ctx context.Context) (*k8s.Deployment, error) {
	list := k8s.DeploymentList{}
	if err := s.client.List(ctx, &list, &s.listOpts); err != nil {
		return nil, fmt.Errorf("failed to get deployment through client: %w", err)
	}
	if len(list.Items) == 0 {
		return nil, fmt.Errorf("no deployment found with the label(%s)", s.listOpts.LabelSelector)
	}

	return &list.Items[0], nil
}

func (s *subProcess) updateDeployment(
	ctx context.Context, newPVC string, deployment *k8s.Deployment, snapshotTime time.Time,
) error {
	if deployment.Spec.Template.ObjectMeta.Annotations == nil {
		deployment.Spec.Template.ObjectMeta.Annotations = map[string]string{}
	}
	now := time.Now().UTC().Format(time.RFC3339)
	deployment.Spec.Template.ObjectMeta.Annotations["kubectl.kubernetes.io/restartedAt"] = now

	if deployment.Annotations == nil {
		deployment.Annotations = map[string]string{}
	}
	deployment.Annotations[vald.LastTimeSnapshotTimestampAnnotationsKey] = snapshotTime.UTC().Format(vald.TimeFormat)

	for _, vol := range deployment.Spec.Template.Spec.Volumes {
		if vol.Name == s.volumeName {
			vol.PersistentVolumeClaim.ClaimName = newPVC
		}
	}

	log.Infof("updating deployment(%s)...", deployment.GetName())
	log.Debugf("deployment detail: %#v", deployment)

	if err := s.client.Update(ctx, deployment); err != nil {
		return fmt.Errorf("failed to update deployment: %w", err)
	}

	return nil
}

func (s *subProcess) deleteSnapshot(ctx context.Context, snapshot *k8s.VolumeSnapshot) error {
	watcher, err := s.client.Watch(ctx,
		&k8s.VolumeSnapshotList{
			Items: []k8s.VolumeSnapshot{*snapshot},
		},
		&s.listOpts,
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
				if event.Type == k8s.WatchDeletedEvent {
					log.Infof("volume snapshot(%s) deleted", snapshot.GetName())
					return nil
				} else {
					log.Debugf("watching volume snapshot(%s) events. event: %v", snapshot.GetName(), event.Type)
				}
			}
		}
	})

	if err := s.client.Delete(ctx, snapshot); err != nil {
		return fmt.Errorf("failed to delete snapshot: %w", err)
	}
	return eg.Wait()
}

func (s *subProcess) deletePVC(ctx context.Context, pvc *k8s.PersistentVolumeClaim) error {
	watcher, err := s.client.Watch(ctx,
		&k8s.PersistentVolumeClaimList{
			Items: []k8s.PersistentVolumeClaim{*pvc},
		},
		&s.listOpts,
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
				if event.Type == k8s.WatchDeletedEvent {
					log.Infof("PVC(%s) deleted", pvc.GetName())
					return nil
				} else {
					log.Debugf("watching PVC(%s) events. event: %v", pvc.GetName(), event.Type)
				}
			}
		}
	})

	if err := s.client.Delete(ctx, pvc); err != nil {
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

func (r *rotator) parseReplicaID(replicaID string, c client.Client) ([]string, error) {
	if replicaID == "" {
		return nil, errors.ErrReadReplicaIDEmpty
	}

	if replicaID == rotateAllID {
		var deploymentList k8s.DeploymentList
		selector, err := c.LabelSelector(r.readReplicaLabelKey, k8s.SelectionOpExists, []string{})
		if err != nil {
			return nil, err
		}
		if err := c.List(context.Background(), &deploymentList, &k8s.ListOptions{
			Namespace:     r.namespace,
			LabelSelector: selector,
		}); err != nil {
			return nil, fmt.Errorf("failed to List deployments in parseReplicaID: %w", err)
		}

		deployments := deploymentList.Items
		if len(deployments) == 0 {
			return nil, fmt.Errorf("no read replica found to rotate")
		}

		var ids []string
		for i := range deployments {
			deployment := &deployments[i]
			ids = append(ids, deployment.Labels[r.readReplicaLabelKey])
		}
		return ids, nil
	}

	return strings.Split(replicaID, ","), nil
}
