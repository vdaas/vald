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

package k8s

import (
	"reflect"

	corev1 "k8s.io/api/core/v1"
)

type Node struct {
	Type            string
	Properties      map[string]*Node
	Items           *Node
	PreserveUnknown bool
}

var byName = map[string]reflect.Type{
	"affinity":                  reflect.TypeOf(corev1.Affinity{}),
	"tolerations":               reflect.TypeOf(corev1.Toleration{}),
	"topologySpreadConstraints": reflect.TypeOf(corev1.TopologySpreadConstraint{}),
	"resources":                 reflect.TypeOf(corev1.ResourceRequirements{}),
}

var scalarByType = map[string]string{
	"k8s.io/apimachinery/pkg/api/resource.Quantity":   "string",
	"k8s.io/apimachinery/pkg/util/intstr.IntOrString": "string",
	"k8s.io/apimachinery/pkg/apis/meta/v1.Time":       "string",
	"k8s.io/apimachinery/pkg/apis/meta/v1.MicroTime":  "string",
	"k8s.io/apimachinery/pkg/apis/meta/v1.Duration":   "string",
}

func Infer(name string) (n *Node, ok bool) {
	t, ok := byName[name]
	if !ok {
		return nil, false
	}
	return build(t, make(map[reflect.Type]bool)), true
}

func build(t reflect.Type, seen map[reflect.Type]bool) *Node {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.String:
		return &Node{Type: "string"}
	case reflect.Bool:
		return &Node{Type: "boolean"}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return &Node{Type: "integer"}
	case reflect.Float32, reflect.Float64:
		return &Node{Type: "number"}
	case reflect.Slice, reflect.Array:
		if t.Elem().Kind() == reflect.Uint8 {
			return &Node{Type: "string"}
		}
		return &Node{Type: "array", Items: build(t.Elem(), seen)}
	case reflect.Map:
		return &Node{Type: "object", PreserveUnknown: true}
	case reflect.Struct:
		if prim, ok := scalarByType[t.PkgPath()+"."+t.Name()]; ok {
			return &Node{Type: prim}
		}
		if seen[t] {
			return &Node{Type: "object", PreserveUnknown: true}
		}
		seen[t] = true
		defer delete(seen, t)

		props := make(map[string]*Node)
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if f.PkgPath != "" {
				continue
			}
			name, inline, skip := jsonField(f)
			if skip {
				continue
			}
			if name == "" {
				if inline || f.Anonymous {
					child := build(f.Type, seen)
					for k, v := range child.Properties {
						props[k] = v
					}
				}
				continue
			}
			props[name] = build(f.Type, seen)
		}
		if len(props) == 0 {
			return &Node{Type: "object", PreserveUnknown: true}
		}
		return &Node{Type: "object", Properties: props}
	default:
		return &Node{Type: "object", PreserveUnknown: true}
	}
}

func jsonField(f reflect.StructField) (name string, inline, skip bool) {
	tag, ok := f.Tag.Lookup("json")
	if !ok {
		return f.Name, false, false
	}
	if tag == "-" {
		return "", false, true
	}
	name = tag
	if i := indexByte(tag, ','); i >= 0 {
		name = tag[:i]
		opts := tag[i+1:]
		inline = hasOption(opts, "inline")
	}
	return name, inline, false
}

func indexByte(s string, b byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == b {
			return i
		}
	}
	return -1
}

func hasOption(opts, want string) bool {
	for len(opts) > 0 {
		var o string
		if i := indexByte(opts, ','); i >= 0 {
			o, opts = opts[:i], opts[i+1:]
		} else {
			o, opts = opts, ""
		}
		if o == want {
			return true
		}
	}
	return false
}
