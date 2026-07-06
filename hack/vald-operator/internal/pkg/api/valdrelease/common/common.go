// +k8s:deepcopy-gen=package
// +groupName=vald.vdaas.org
package common

var (
	KindTypeDaemonSet   KindType = "DaemonSet"
	KindTypeDeployment  KindType = "Deployment"
	KindTypeStatefulSet KindType = "StatefulSet"
)

type KindType string
