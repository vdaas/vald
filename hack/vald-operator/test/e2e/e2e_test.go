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

package e2e

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/vdaas/vald/hack/vald-operator/test/utils"
)

// namespace where the project is deployed in
const namespace = "mvaldrelease-system"
const sampleMvrsName = "mvaldrelease-sample"
const reconcileNamespace = "mvaldrelease-reconcile-e2e"
const valdNamespace = "mvaldrelease-vald-e2e"
const singleClusterMvrsName = "mvaldrelease-single-cluster"
const valdReleaseCRDPath = "../../charts/vald-helm-operator/crds/valdrelease.yaml"
const waitForValdTimeout = 10 * time.Minute

// serviceAccountName created for the project
const serviceAccountName = "mvaldrelease-controller-manager"

// metricsServiceName is the name of the metrics service of the project
const metricsServiceName = "mvaldrelease-controller-manager-metrics-service"

// metricsRoleBindingName is the name of the RBAC that will be created to allow get the metrics data
const metricsRoleBindingName = "mvaldrelease-metrics-binding"

var _ = Describe("Manager", Ordered, func() {
	var controllerPodName string

	// Before running the tests, set up the environment by creating the namespace,
	// enforce the restricted security policy to the namespace, installing CRDs,
	// and deploying the controller.
	BeforeAll(func() {
		By("creating manager namespace")
		cmd := exec.Command("kubectl", "create", "ns", namespace)
		_, err := utils.Run(cmd)
		Expect(err).NotTo(HaveOccurred(), "Failed to create namespace")

		By("labeling the namespace to enforce the restricted security policy")
		cmd = exec.Command("kubectl", "label", "--overwrite", "ns", namespace,
			"pod-security.kubernetes.io/enforce=restricted")
		_, err = utils.Run(cmd)
		Expect(err).NotTo(HaveOccurred(), "Failed to label namespace with restricted policy")

		By("installing CRDs")
		cmd = exec.Command("make", "install")
		_, err = utils.Run(cmd)
		Expect(err).NotTo(HaveOccurred(), "Failed to install CRDs")

		By("deploying the controller-manager")
		cmd = exec.Command("make", "deploy", fmt.Sprintf("IMG=%s", projectImage))
		_, err = utils.Run(cmd)
		Expect(err).NotTo(HaveOccurred(), "Failed to deploy the controller-manager")
	})

	// After all tests have been executed, clean up by undeploying the controller, uninstalling CRDs,
	// and deleting the namespace.
	AfterAll(func() {
		By("cleaning up the curl pod for metrics")
		cmd := exec.Command("kubectl", "delete", "pod", "curl-metrics", "-n", namespace)
		_, _ = utils.Run(cmd)

		By("undeploying the controller-manager")
		cmd = exec.Command("make", "undeploy")
		_, _ = utils.Run(cmd)

		By("uninstalling CRDs")
		cmd = exec.Command("make", "uninstall")
		_, _ = utils.Run(cmd)

		By("removing manager namespace")
		cmd = exec.Command("kubectl", "delete", "ns", namespace)
		_, _ = utils.Run(cmd)
	})

	// After each test, check for failures and collect logs, events,
	// and pod descriptions for debugging.
	AfterEach(func() {
		specReport := CurrentSpecReport()
		if specReport.Failed() {
			By("Fetching controller manager pod logs")
			cmd := exec.Command("kubectl", "logs", controllerPodName, "-n", namespace)
			controllerLogs, err := utils.Run(cmd)
			if err == nil {
				_, _ = fmt.Fprintf(GinkgoWriter, "Controller logs:\n %s", controllerLogs)
			} else {
				_, _ = fmt.Fprintf(GinkgoWriter, "Failed to get Controller logs: %s", err)
			}

			By("Fetching Kubernetes events")
			cmd = exec.Command("kubectl", "get", "events", "-n", namespace, "--sort-by=.lastTimestamp")
			eventsOutput, err := utils.Run(cmd)
			if err == nil {
				_, _ = fmt.Fprintf(GinkgoWriter, "Kubernetes events:\n%s", eventsOutput)
			} else {
				_, _ = fmt.Fprintf(GinkgoWriter, "Failed to get Kubernetes events: %s", err)
			}

			By("Fetching curl-metrics logs")
			cmd = exec.Command("kubectl", "logs", "curl-metrics", "-n", namespace)
			metricsOutput, err := utils.Run(cmd)
			if err == nil {
				_, _ = fmt.Fprintf(GinkgoWriter, "Metrics logs:\n %s", metricsOutput)
			} else {
				_, _ = fmt.Fprintf(GinkgoWriter, "Failed to get curl-metrics logs: %s", err)
			}

			By("Fetching controller manager pod description")
			cmd = exec.Command("kubectl", "describe", "pod", controllerPodName, "-n", namespace)
			podDescription, err := utils.Run(cmd)
			if err == nil {
				fmt.Println("Pod description:\n", podDescription)
			} else {
				fmt.Println("Failed to describe controller pod")
			}
		}
	})

	SetDefaultEventuallyTimeout(2 * time.Minute)
	SetDefaultEventuallyPollingInterval(time.Second)

	Context("Manager", func() {
		It("should run successfully", func() {
			By("validating that the controller-manager pod is running as expected")
			verifyControllerUp := func(g Gomega) {
				// Get the name of the controller-manager pod
				cmd := exec.Command("kubectl", "get",
					"pods", "-l", "control-plane=controller-manager",
					"-o", "go-template={{ range .items }}"+
						"{{ if not .metadata.deletionTimestamp }}"+
						"{{ .metadata.name }}"+
						"{{ \"\\n\" }}{{ end }}{{ end }}",
					"-n", namespace,
				)

				podOutput, err := utils.Run(cmd)
				g.Expect(err).NotTo(HaveOccurred(), "Failed to retrieve controller-manager pod information")
				podNames := utils.GetNonEmptyLines(podOutput)
				g.Expect(podNames).To(HaveLen(1), "expected 1 controller pod running")
				controllerPodName = podNames[0]
				g.Expect(controllerPodName).To(ContainSubstring("controller-manager"))

				// Validate the pod's status
				cmd = exec.Command("kubectl", "get",
					"pods", controllerPodName, "-o", "jsonpath={.status.phase}",
					"-n", namespace,
				)
				output, err := utils.Run(cmd)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(output).To(Equal("Running"), "Incorrect controller-manager pod status")
			}
			Eventually(verifyControllerUp).Should(Succeed())
		})

		It("should ensure the metrics endpoint is serving metrics", func() {
			By("creating a ClusterRoleBinding for the service account to allow access to metrics")
			cmd := exec.Command("kubectl", "create", "clusterrolebinding", metricsRoleBindingName,
				"--clusterrole=mvaldrelease-metrics-reader",
				fmt.Sprintf("--serviceaccount=%s:%s", namespace, serviceAccountName),
			)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred(), "Failed to create ClusterRoleBinding")

			By("validating that the metrics service is available")
			cmd = exec.Command("kubectl", "get", "service", metricsServiceName, "-n", namespace)
			_, err = utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred(), "Metrics service should exist")

			By("getting the service account token")
			token, err := serviceAccountToken()
			Expect(err).NotTo(HaveOccurred())
			Expect(token).NotTo(BeEmpty())

			By("waiting for the metrics endpoint to be ready")
			verifyMetricsEndpointReady := func(g Gomega) {
				cmd := exec.Command("kubectl", "get", "endpoints", metricsServiceName, "-n", namespace)
				output, err := utils.Run(cmd)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(output).To(ContainSubstring("8443"), "Metrics endpoint is not ready")
			}
			Eventually(verifyMetricsEndpointReady).Should(Succeed())

			By("verifying that the controller manager is serving the metrics server")
			verifyMetricsServerStarted := func(g Gomega) {
				cmd := exec.Command("kubectl", "logs", controllerPodName, "-n", namespace)
				output, err := utils.Run(cmd)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(output).To(ContainSubstring("controller-runtime.metrics\tServing metrics server"),
					"Metrics server not yet started")
			}
			Eventually(verifyMetricsServerStarted).Should(Succeed())

			By("creating the curl-metrics pod to access the metrics endpoint")
			cmd = exec.Command("kubectl", "run", "curl-metrics", "--restart=Never",
				"--namespace", namespace,
				"--image=curlimages/curl:latest",
				"--overrides",
				fmt.Sprintf(`{
					"spec": {
						"containers": [{
							"name": "curl",
							"image": "curlimages/curl:latest",
							"command": ["/bin/sh", "-c"],
							"args": ["curl -v -k -H 'Authorization: Bearer %s' https://%s.%s.svc.cluster.local:8443/metrics"],
							"securityContext": {
								"allowPrivilegeEscalation": false,
								"capabilities": {
									"drop": ["ALL"]
								},
								"runAsNonRoot": true,
								"runAsUser": 1000,
								"seccompProfile": {
									"type": "RuntimeDefault"
								}
							}
						}],
						"serviceAccount": "%s"
					}
				}`, token, metricsServiceName, namespace, serviceAccountName))
			_, err = utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred(), "Failed to create curl-metrics pod")

			By("waiting for the curl-metrics pod to complete.")
			verifyCurlUp := func(g Gomega) {
				cmd := exec.Command("kubectl", "get", "pods", "curl-metrics",
					"-o", "jsonpath={.status.phase}",
					"-n", namespace)
				output, err := utils.Run(cmd)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(output).To(Equal("Succeeded"), "curl pod in wrong status")
			}
			Eventually(verifyCurlUp, 5*time.Minute).Should(Succeed())

			By("getting the metrics by checking curl-metrics logs")
			metricsOutput := getMetricsOutput()
			Expect(metricsOutput).To(ContainSubstring(
				"controller_runtime_reconcile_total",
			))
		})
	})

	Context("Mvaldrelease reconciliation", Ordered, func() {
		BeforeAll(func() {
			By("installing the ValdRelease CRD required for generated resources")
			cmd := exec.Command("kubectl", "apply", "-f", valdReleaseCRDPath)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred(), "Failed to install ValdRelease CRD")

			By("creating a dedicated namespace for reconcile verification")
			Expect(createNamespace(reconcileNamespace)).To(Succeed())

			By("applying the sample Mvaldrelease resource")
			Expect(applyManifest(reconcileNamespace, "config/samples/controller_v1_mvaldrelease.yaml")).To(Succeed())
		})

		AfterAll(func() {
			By("deleting the sample Mvaldrelease resource")
			cmd := exec.Command("kubectl", "delete", "mvaldrelease", sampleMvrsName, "-n", reconcileNamespace, "--ignore-not-found=true")
			_, _ = utils.Run(cmd)

			By("removing the reconcile verification namespace")
			Expect(deleteNamespace(reconcileNamespace)).To(Succeed())
		})

		It("should create ValdRelease resources for active clusters", func() {
			verifyValdReleasesCreated := func(g Gomega) {
				names, err := listValdReleaseNames(reconcileNamespace)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(names).To(HaveLen(2), "expected one ValdRelease per active cluster")
			}
			Eventually(verifyValdReleasesCreated).Should(Succeed())
		})

		It("should advance the Mvaldrelease phase to Completed", func() {
			verifyCompleted := func(g Gomega) {
				status, err := getMvaldreleaseStatus(reconcileNamespace, sampleMvrsName)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(status.Status.Phase).To(Equal("Completed"))
				g.Expect(hasCondition(status.Status.Conditions, "Completed", "True")).To(BeTrue(),
					"expected Completed=True condition")
			}
			Eventually(verifyCompleted).Should(Succeed())
		})

		It("should recreate a deleted ValdRelease", func() {
			names, err := listValdReleaseNames(reconcileNamespace)
			Expect(err).NotTo(HaveOccurred())
			Expect(names).NotTo(BeEmpty())

			deleted := names[0]

			By("deleting one generated ValdRelease")
			cmd := exec.Command("kubectl", "delete", "valdrelease", deleted, "-n", reconcileNamespace)
			_, err = utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred(), "Failed to delete generated ValdRelease")

			By("waiting for the controller to recreate it")
			verifyRecreated := func(g Gomega) {
				recreated, err := listValdReleaseNames(reconcileNamespace)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(recreated).To(HaveLen(2))
				g.Expect(recreated).To(ContainElement(deleted))
			}
			Eventually(verifyRecreated).Should(Succeed())
		})
	})

	Context("Vald deployment", Ordered, func() {
		var manifestPath string

		BeforeAll(func() {
			By("creating a dedicated namespace for Vald deployment")
			Expect(createNamespace(valdNamespace)).To(Succeed())

			By("deploying the Vald Helm Operator")
			cmd := exec.Command("make", "-C", "../..", "k8s/vald-helm-operator/deploy")
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred(), "Failed to deploy vald-helm-operator")

			By("writing a single-cluster Mvaldrelease manifest to avoid resource name collisions")
			manifestPath, err = writeSingleClusterManifest(valdNamespace)
			Expect(err).NotTo(HaveOccurred(), "Failed to create single-cluster Mvaldrelease manifest")

			By("applying the single-cluster Mvaldrelease resource")
			Expect(applyManifest(valdNamespace, manifestPath)).To(Succeed())
		})

		AfterAll(func() {
			By("deleting the single-cluster Mvaldrelease resource")
			cmd := exec.Command("kubectl", "delete", "mvaldrelease", singleClusterMvrsName, "-n", valdNamespace, "--ignore-not-found=true")
			_, _ = utils.Run(cmd)

			By("removing the Vald deployment namespace")
			Expect(deleteNamespace(valdNamespace)).To(Succeed())

			By("undeploying the Vald Helm Operator")
			cmd = exec.Command("make", "-C", "../..", "k8s/vald-helm-operator/delete")
			_, _ = utils.Run(cmd)

			if manifestPath != "" {
				Expect(os.Remove(manifestPath)).To(Succeed())
			}
		})

		It("should reconcile a single-cluster Mvaldrelease to Completed", func() {
			verifyCompleted := func(g Gomega) {
				status, err := getMvaldreleaseStatus(valdNamespace, singleClusterMvrsName)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(status.Status.Phase).To(Equal("Completed"))
				g.Expect(hasCondition(status.Status.Conditions, "Completed", "True")).To(BeTrue())

				names, err := listValdReleaseNames(valdNamespace)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(names).To(HaveLen(1), "expected a single generated ValdRelease")
			}
			Eventually(verifyCompleted).Should(Succeed())
		})

		It("should deploy core Vald workloads", func() {
			By("waiting for the gateway service to exist")
			waitForResource("service", "vald-lb-gateway", valdNamespace, waitForValdTimeout)

			By("waiting for discoverer deployment to become available")
			waitForResourceCondition("deployment", "vald-discoverer", valdNamespace, "Available", waitForValdTimeout)

			By("waiting for gateway deployment to become available")
			waitForResourceCondition("deployment", "vald-lb-gateway", valdNamespace, "Available", waitForValdTimeout)

			By("waiting for agent statefulset to become ready")
			waitForResourceJSONPath("statefulset", "vald-agent", valdNamespace, "{.status.readyReplicas}", "1", waitForValdTimeout)
		})
	})
})

// serviceAccountToken returns a token for the specified service account in the given namespace.
// It uses the Kubernetes TokenRequest API to generate a token by directly sending a request
// and parsing the resulting token from the API response.
func serviceAccountToken() (string, error) {
	const tokenRequestRawString = `{
		"apiVersion": "authentication.k8s.io/v1",
		"kind": "TokenRequest"
	}`

	// Temporary file to store the token request
	secretName := fmt.Sprintf("%s-token-request", serviceAccountName)
	tokenRequestFile := filepath.Join(os.TempDir(), secretName)
	err := os.WriteFile(tokenRequestFile, []byte(tokenRequestRawString), os.FileMode(0o644))
	if err != nil {
		return "", err
	}

	var out string
	verifyTokenCreation := func(g Gomega) {
		// Execute kubectl command to create the token
		cmd := exec.Command("kubectl", "create", "--raw", fmt.Sprintf(
			"/api/v1/namespaces/%s/serviceaccounts/%s/token",
			namespace,
			serviceAccountName,
		), "-f", tokenRequestFile)

		output, err := cmd.CombinedOutput()
		g.Expect(err).NotTo(HaveOccurred())

		// Parse the JSON output to extract the token
		var token tokenRequest
		err = json.Unmarshal(output, &token)
		g.Expect(err).NotTo(HaveOccurred())

		out = token.Status.Token
	}
	Eventually(verifyTokenCreation).Should(Succeed())

	return out, err
}

// getMetricsOutput retrieves and returns the logs from the curl pod used to access the metrics endpoint.
func getMetricsOutput() string {
	By("getting the curl-metrics logs")
	cmd := exec.Command("kubectl", "logs", "curl-metrics", "-n", namespace)
	metricsOutput, err := utils.Run(cmd)
	Expect(err).NotTo(HaveOccurred(), "Failed to retrieve logs from curl pod")
	Expect(metricsOutput).To(ContainSubstring("< HTTP/1.1 200 OK"))
	return metricsOutput
}

// tokenRequest is a simplified representation of the Kubernetes TokenRequest API response,
// containing only the token field that we need to extract.
type tokenRequest struct {
	Status struct {
		Token string `json:"token"`
	} `json:"status"`
}

type mvaldreleaseStatusResponse struct {
	Status struct {
		Phase      string `json:"phase"`
		Conditions []struct {
			Type   string `json:"type"`
			Status string `json:"status"`
		} `json:"conditions"`
	} `json:"status"`
}

func createNamespace(name string) error {
	cmd := exec.Command("kubectl", "create", "ns", name)
	_, err := utils.Run(cmd)
	return err
}

func deleteNamespace(name string) error {
	cmd := exec.Command("kubectl", "delete", "ns", name, "--ignore-not-found=true")
	_, err := utils.Run(cmd)
	return err
}

func applyManifest(namespace, manifestPath string) error {
	cmd := exec.Command("kubectl", "apply", "-n", namespace, "-f", manifestPath)
	_, err := utils.Run(cmd)
	return err
}

func listValdReleaseNames(namespace string) ([]string, error) {
	cmd := exec.Command("kubectl", "get", "valdrelease", "-n", namespace,
		"-o", `jsonpath={range .items[*]}{.metadata.name}{"\n"}{end}`)
	output, err := utils.Run(cmd)
	if err != nil {
		return nil, err
	}
	return utils.GetNonEmptyLines(output), nil
}

func getMvaldreleaseStatus(namespace, name string) (*mvaldreleaseStatusResponse, error) {
	cmd := exec.Command("kubectl", "get", "mvaldrelease", name, "-n", namespace, "-o", "json")
	output, err := utils.Run(cmd)
	if err != nil {
		return nil, err
	}

	var status mvaldreleaseStatusResponse
	if err := json.Unmarshal([]byte(output), &status); err != nil {
		return nil, err
	}
	return &status, nil
}

func hasCondition(conditions []struct {
	Type   string `json:"type"`
	Status string `json:"status"`
}, condType, condStatus string) bool {
	for _, cond := range conditions {
		if cond.Type == condType && cond.Status == condStatus {
			return true
		}
	}
	return false
}

func waitForResource(kind, name, namespace string, timeout time.Duration) {
	verifyExists := func(g Gomega) {
		cmd := exec.Command("kubectl", "get", kind, name, "-n", namespace)
		_, err := utils.Run(cmd)
		g.Expect(err).NotTo(HaveOccurred())
	}
	Eventually(verifyExists, timeout).Should(Succeed())
}

func waitForResourceCondition(kind, name, namespace, condition string, timeout time.Duration) {
	verifyAvailable := func(g Gomega) {
		cmd := exec.Command("kubectl", "wait",
			fmt.Sprintf("%s/%s", kind, name),
			fmt.Sprintf("--for=condition=%s", condition),
			fmt.Sprintf("--timeout=%s", 30*time.Second),
			"-n", namespace,
		)
		_, err := utils.Run(cmd)
		g.Expect(err).NotTo(HaveOccurred())
	}
	Eventually(verifyAvailable, timeout).Should(Succeed())
}

func waitForResourceJSONPath(kind, name, namespace, jsonPath, expected string, timeout time.Duration) {
	verifyField := func(g Gomega) {
		cmd := exec.Command("kubectl", "get", kind, name, "-n", namespace, "-o", "jsonpath="+jsonPath)
		output, err := utils.Run(cmd)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(output).To(Equal(expected))
	}
	Eventually(verifyField, timeout).Should(Succeed())
}

func writeSingleClusterManifest(namespace string) (string, error) {
	manifest := fmt.Sprintf(`apiVersion: vald.vdaas.org/v1
kind: Mvaldrelease
metadata:
  name: %s
  namespace: %s
spec:
  infrastructure:
    - role: "green"
      type: "kubernetes"
      active: true
      clusters:
        - id: "12345678-1234-1234-1234-123456789012"
          name: "single-cluster"
      nodePools:
        general:
          name: "node-pool-1-general"
          machineResource:
            name: "highmemory"
            cpu: "1"
            memory: "2Gi"
            storage: "10Gi"
          replicas: 1
        agent:
          name: "agent-pool-1-agent"
          machineResource:
            name: "highcpu"
            cpu: "1"
            memory: "2Gi"
            storage: "10Gi"
          replicas: 1
  vectorEngine:
    name: "vald"
    vald:
      defaults:
        logLevel: "info"
      agent:
        ngt:
          creationEdgeSize: 10
          searchEdgeSize: 40
          dimension: 2
          distanceType: "l2"
          objectType: float
        persistentVolume:
          enabled: false
      indexer:
        manager: false
        indexDuration: "24h"
        indexSchedule: "0 2 * * *"
        saveDuration: "1h"
        saveSchedule: "0 * * * *"
        concurrency: 1
      gateway:
        indexReplica: 1
      discoverer:
        kind: "Deployment"
`, singleClusterMvrsName, namespace)

	path := filepath.Join(os.TempDir(), fmt.Sprintf("%s.yaml", singleClusterMvrsName))
	if err := os.WriteFile(path, []byte(manifest), 0o600); err != nil {
		return "", err
	}
	return path, nil
}
