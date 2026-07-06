package common

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var scales = []resource.Scale{
	//resource.Giga,
	resource.Mega,
	resource.Kilo,
}

// BuildResources constructs normalized ResourceRequirements from a machine resource list.
// ratio is the fraction of machine resources to allocate.
// requestDiv divides the Requests (e.g., number of pods per node).
// limitDiv divides the Limits (use 1 for no division).
func BuildResources(mc v1.ResourceList, ratio, requestDiv, limitDiv float64) *v1.ResourceRequirements {
	cpu := mc.Cpu().Value()
	memory := mc.Memory().Value()
	res := &v1.ResourceRequirements{
		Requests: v1.ResourceList{
			v1.ResourceCPU:    CalcResource(cpu, ratio, requestDiv),
			v1.ResourceMemory: CalcResource(memory, ratio, requestDiv),
		},
		Limits: v1.ResourceList{
			v1.ResourceCPU:    CalcResource(cpu, ratio, limitDiv),
			v1.ResourceMemory: CalcResource(memory, ratio, limitDiv),
		},
	}
	return &v1.ResourceRequirements{
		Requests: NormalizeResourceList(res.Requests),
		Limits:   NormalizeResourceList(res.Limits),
	}
}

func CalcResource(value int64, ratio float64, divs ...float64) resource.Quantity {

	result := float64(value) * ratio
	for _, d := range divs {
		result /= d
	}

	res := resource.Quantity{}
	res.SetMilli(int64(result * 1000))
	return res
}

func NormalizeResourceList(rl v1.ResourceList) v1.ResourceList {
	normalized := v1.ResourceList{}
	for name, quantity := range rl {
		switch name {
		case v1.ResourceCPU:
			quantity.Format = resource.DecimalSI
			milli := quantity.ScaledValue(resource.Milli)
			if milli >= 1000 && milli%1000 == 0 {
				quantity.SetScaled(milli/1000, 0)
				break
			}
			quantity.SetScaled(milli, resource.Milli)
		case v1.ResourceMemory:
			quantity.RoundUp(GetScale(quantity))
			quantity.Format = resource.DecimalSI
		}
		normalized[name] = quantity.DeepCopy()
	}
	return normalized
}

func GetScale(quantity resource.Quantity) resource.Scale {
	for _, scale := range scales {
		if quantity.ScaledValue(scale) >= 1 {
			return scale
		}
	}
	return resource.Milli
}

// BuildTopologySpreadConstraint returns a TSC that spreads pods by hostname
// using the given app.kubernetes.io/component label value as the selector.
func BuildTopologySpreadConstraint(componentLabel string) v1.TopologySpreadConstraint {
	return v1.TopologySpreadConstraint{
		MaxSkew:           1,
		TopologyKey:       "kubernetes.io/hostname",
		WhenUnsatisfiable: v1.DoNotSchedule,
		LabelSelector: &metav1.LabelSelector{
			MatchLabels: map[string]string{
				"app.kubernetes.io/component": componentLabel,
			},
		},
	}
}
