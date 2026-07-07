//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package controller

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	controllerv1 "github.com/vdaas/vald/hack/vald-operator/api/v1"
	"github.com/vdaas/vald/hack/vald-operator/internal/infrastructure/config"
)

var _ = Describe("ValdOperatorRelease Controller", func() {
	Context("When reconciling a resource", func() {
		const resourceName = "test-resource"

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: "default", // TODO(user):Modify as needed
		}
		valdoperatorrelease := &controllerv1.ValdOperatorRelease{}

		BeforeEach(func() {
			By("creating the custom resource for the Kind ValdOperatorRelease")
			err := k8sClient.Get(ctx, typeNamespacedName, valdoperatorrelease)
			if err != nil && errors.IsNotFound(err) {
				resource := &controllerv1.ValdOperatorRelease{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: "default",
					},
					Spec: controllerv1.ValdOperatorReleaseSpec{
						Infrastructure: []controllerv1.ValdOperatorReleaseInfra{
							{
								Role:   "green",
								Type:   "local",
								Active: true,
								Clusters: []controllerv1.DestClusters{
									{
										ID:   "12345678-1234-1234-1234-123456789012",
										Name: "lpokijhbgaev",
									},
									{
										ID:   "87654321-4321-4321-4321-210987654321",
										Name: "bhnjmktyuioi",
									},
								},
								NodePools: controllerv1.NodePools{
									controllerv1.NodePoolTypeGeneral: {
										Name: "node-pool-1-general",
										MachineResource: controllerv1.MachineResource{
											Name:    "c4.large",
											Cpu:     "1",
											Memory:  "2Gi",
											Storage: "10Gi",
										},
										Replicas: 1,
									},
									controllerv1.NodePoolTypeValdAgent: {
										Name: "agent-pool-1-agent",
										MachineResource: controllerv1.MachineResource{
											Name:    "highcpu",
											Cpu:     "1",
											Memory:  "2Gi",
											Storage: "10Gi",
										},
										Replicas: 1,
									},
								},
							},
						},
						VectorEngine: controllerv1.VectorEngine{
							Name: "vald",
							Vald: controllerv1.Vald{
								Defaults: controllerv1.ValdDefaults{
									LogLevel: "info",
								},
								Agent: controllerv1.Agent{
									Ngt: controllerv1.Ngt{
										CreationEdgeSize: 10,
										SearchEdgeSize:   40,
										Dimension:        128,
										DistanceType:     "l2",
										ObjectType:       "float",
									},
									PersistentVolume: &controllerv1.AgentPersistentVolume{
										Enabled: true,
									},
								},
								Indexer: controllerv1.Indexer{
									Manager:       false,
									IndexDuration: "24h",
									IndexSchedule: "0 2 * * *",
									SaveDuration:  "1h",
									SaveSchedule:  "0 * * * *",
									Concurrency:   2,
								},
								Gateway: controllerv1.Gateway{
									IndexReplica: 2,
								},
								Discoverer: controllerv1.Discoverer{
									Kind: "Deployment",
								},
							},
						},
					},
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
		})

		AfterEach(func() {
			// TODO(user): Cleanup logic after each test, like removing the resource instance.
			resource := &controllerv1.ValdOperatorRelease{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			Expect(err).NotTo(HaveOccurred())

			By("Cleanup the specific resource instance ValdOperatorRelease")
			Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
		})
		It("should successfully reconcile the resource", func() {
			By("Reconciling the created resource")
			controllerReconciler := &ValdOperatorReleaseReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
				Config: &config.Config{},
				Syncer: NewResourceSyncer(k8sClient, k8sClient.Scheme()),
			}

			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())
			// TODO(user): Add more specific assertions depending on your controller's reconciliation logic.
			// Example: If you expect a certain status condition after reconciliation, verify it here.
		})
	})
})
