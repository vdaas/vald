package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	snapshotv1 "github.com/kubernetes-csi/external-snapshotter/client/v6/apis/volumesnapshot/v1"
	snapshotclient "github.com/kubernetes-csi/external-snapshotter/client/v6/clientset/versioned"
	"github.com/vdaas/vald/internal/sync/errgroup"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	namespace        = "default"
	replicaid        = 0
	deploymentPrefix = "vald-agent-ngt-readreplica"
	snapshotPrefix   = "vald-agent-ngt-snapshot"
	pvcPrefix        = "vald-agent-ngt-readreplica-pvc"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	kubeconfig := filepath.Join(homeDir, ".kube", "config")

	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}

	client, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}

	sclient, err := snapshotclient.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}

	updater := New(replicaid, client, sclient)

	ctx := context.Background()
	if err := updater.Start(ctx); err != nil {
		panic(err)
	}
}

type readreplicaUpdater struct {
	replicaid        int
	namespace        string
	deploymentPrefix string
	snapshotPrefix   string
	pvcPrefix        string
	clientset        *kubernetes.Clientset
	sClientset       *snapshotclient.Clientset
}

func New(replicaid int, client *kubernetes.Clientset, sclient *snapshotclient.Clientset) *readreplicaUpdater {
	return &readreplicaUpdater{
		replicaid:        replicaid,
		namespace:        namespace,
		deploymentPrefix: fmt.Sprintf("%s-%d", deploymentPrefix, replicaid),
		snapshotPrefix:   fmt.Sprintf("%s-%d", snapshotPrefix, replicaid),
		pvcPrefix:        fmt.Sprintf("%s-%d", pvcPrefix, replicaid),
		clientset:        client,
		sClientset:       sclient,
	}
}

func (r *readreplicaUpdater) Start(ctx context.Context) error {
	newSnapShot, oldSnapShot, err := r.createSnapshot(ctx)
	if err != nil {
		return err
	}

	newPvc, oldPvc, err := r.createPVC(ctx, newSnapShot.Name)
	if err != nil {
		return err
	}

	err = r.updateDeployment(ctx, newPvc.Name)
	if err != nil {
		return err
	}

	err = r.deleteSnapshot(ctx, oldSnapShot.Name)
	if err != nil {
		return err
	}
	err = r.deletePVC(ctx, oldPvc.Name)
	if err != nil {
		return err
	}

	return nil
}

func (r *readreplicaUpdater) createSnapshot(ctx context.Context) (new, old *snapshotv1.VolumeSnapshot, err error) {
	snapshotInterface := r.sClientset.SnapshotV1().VolumeSnapshots(namespace)

	list, err := snapshotInterface.List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, nil, err
	}

	for _, snap := range list.Items {
		if strings.HasPrefix(snap.Name, r.snapshotPrefix) {
			old = &snap
			break
		}
	}
	if old == nil {
		return nil, nil, fmt.Errorf("old snapshot not found")
	}

	new = &snapshotv1.VolumeSnapshot{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf(r.snapshotPrefix+"-%d", time.Now().Unix()),
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

func (r *readreplicaUpdater) createPVC(ctx context.Context, newSnapShot string) (new, old *v1.PersistentVolumeClaim, err error) {
	pvcInterface := r.clientset.CoreV1().PersistentVolumeClaims(namespace)
	list, err := pvcInterface.List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, nil, err
	}

	for _, pvc := range list.Items {
		if strings.HasPrefix(pvc.Name, r.pvcPrefix) {
			old = &pvc
			break
		}
	}
	if old == nil {
		return nil, nil, fmt.Errorf("old pvc not found")
	}

	new = &v1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf(r.pvcPrefix+"-%d", time.Now().Unix()),
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

func (r *readreplicaUpdater) updateDeployment(ctx context.Context, newPVC string) error {
	deploymentInterface := r.clientset.AppsV1().Deployments(namespace)

	deployment, err := deploymentInterface.Get(ctx, r.deploymentPrefix, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if deployment.Spec.Template.ObjectMeta.Annotations == nil {
		deployment.Spec.Template.ObjectMeta.Annotations = map[string]string{}
	}
	deployment.Spec.Template.ObjectMeta.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)

	for _, vol := range deployment.Spec.Template.Spec.Volumes {
		if vol.Name == fmt.Sprintf("%s-clone", r.deploymentPrefix) {
			vol.PersistentVolumeClaim.ClaimName = newPVC
		}
	}

	_, err = deploymentInterface.Update(ctx, deployment, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (r *readreplicaUpdater) deleteSnapshot(ctx context.Context, snapshot string) error {
	snapshotInterface := r.sClientset.SnapshotV1().VolumeSnapshots(namespace)

	watcher, err := snapshotInterface.Watch(ctx, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to watch snapshot: %w", err)
	}
	defer watcher.Stop()

	eg, egctx := errgroup.New(ctx)
	eg.Go(func() error {
		for {
			select {
			case <-egctx.Done():
				return egctx.Err()
			case event := <-watcher.ResultChan():
				if event.Type == watch.Deleted {
					fmt.Println("old volume snapshot deleted.")
					return nil
				} else {
					fmt.Println("event: ", event.Type)
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

func (r *readreplicaUpdater) deletePVC(ctx context.Context, pvc string) error {
	pvcInterface := r.clientset.CoreV1().PersistentVolumeClaims(namespace)

	watcher, err := pvcInterface.Watch(ctx, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to watch PVC: %w", err)
	}
	defer watcher.Stop()

	eg, egctx := errgroup.New(ctx)
	eg.Go(func() error {
		for {
			select {
			case <-egctx.Done():
				return egctx.Err()
			case event := <-watcher.ResultChan():
				if event.Type == watch.Deleted {
					fmt.Println("old PVC deleted.")
					return nil
				} else {
					fmt.Println("event: ", event.Type)
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
