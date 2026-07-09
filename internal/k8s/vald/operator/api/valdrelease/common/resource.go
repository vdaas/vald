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
package common

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var scales = []resource.Scale{
	resource.Mega,
	resource.Kilo,
}

func CalcResource(value int64, ratio float64, divs ...float64) resource.Quantity {
	result := float64(value) * ratio
	for _, d := range divs {
		if d == 0 {
			// Skip zero divisors to avoid division by zero (+Inf);
			// a zero divisor means "no division" (same as 1).
			continue
		}
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
			quantity.RoundUp(getScale(quantity))
			quantity.Format = resource.DecimalSI
		}
		normalized[name] = quantity.DeepCopy()
	}
	return normalized
}

func getScale(quantity resource.Quantity) resource.Scale {
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
