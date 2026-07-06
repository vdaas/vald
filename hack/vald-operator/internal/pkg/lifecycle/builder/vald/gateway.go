package vald

import (
	controllerv1 "github.com/vdaas/vald/hack/vald-operator/api/v1"
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/api/valdrelease/gateway"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
)

func (b *VrsBuilder) buildGateway() gateway.Gateway {
	return gateway.Gateway{
		Lb: *b.buildLb(),
	}
}

func (b *VrsBuilder) buildLb() *gateway.Lb {
	inputGw := b.CR.Spec.VectorEngine.Vald.Gateway

	return &gateway.Lb{
		Logging: b.buildLogging(inputGw.LogLevel),
		Hpa: &gateway.Hpa{
			TargetCPUUtilizationPercentage: 80,
		},
		GatewayConfig: gateway.GatewayConfig{
			IndexReplica: inputGw.IndexReplica,
		},
		ServiceType: b.getGatewayServiceType(inputGw.ServiceType),
		Ingress:     b.buildIngress(inputGw.Ingress),
	}
}

func (b *VrsBuilder) buildIngress(in *controllerv1.GatewayIngress) gateway.Ingress {
	base := gateway.Ingress{
		DefaultBackend: gateway.DefaultBackend{Enabled: false},
		PathType:       v1beta1.PathTypePrefix,
		ServicePort:    "grpc",
	}
	if in == nil || !in.Enabled {
		return base
	}
	base.Enabled = true
	base.Host = in.Host
	return base
}

func (b *VrsBuilder) getGatewayServiceType(st string) corev1.ServiceType {
	switch st {
	case string(corev1.ServiceTypeClusterIP):
		return corev1.ServiceTypeClusterIP
	case string(corev1.ServiceTypeLoadBalancer):
		return corev1.ServiceTypeLoadBalancer
	default:
		return corev1.ServiceTypeNodePort
	}
}
