//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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
package target

import (
	"fmt"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
)

type dummy struct {
	field1 string
}

func payloadObjectAccesses() {
	sc := &payload.Search_Config{}
	_ = sc.Radius                  // want `\QAvoid to access struct fields directly`
	_, _ = sc.Epsilon, "test"      // want `\QAvoid to access struct fields directly`
	_, _ = "test", sc.Timeout      // want `\QAvoid to access struct fields directly`
	_ = sc.GetRadius()             // OK: use function to access the field
	_, _ = sc.GetEpsilon(), "test" // OK: use function to access the field
	_, _ = "test", sc.GetTimeout() // OK: use function to access the field

	loc := &payload.Object_Location{}
	_ = loc.Ips         // want `\QAvoid to access struct fields directly`
	_ = loc.Name        // want `\QAvoid to access struct fields directly`
	fmt.Print(loc.Ips)  // want `\QAvoid to access struct fields directly`
	fmt.Print(loc.Name) // want `\QAvoid to access struct fields directly`
	_ = loc.GetIps()    // OK: use function to access the field
	_ = loc.GetName()   // OK: use function to access the field

	fmt.Printf("%s, %%", "test: ", loc.Uuid) // want `\QAvoid to access struct fields directly`
	_, _, _ = "test", "test", loc.Uuid       // want `\QAvoid to access struct fields directly`
	_, _, _ = "test", "test", loc.GetUuid()  // OK: use function to access the field

	ireq := &payload.Insert_Request{}
	_ = ireq.Config      // want `\QAvoid to access struct fields directly`
	_ = ireq.Vector      // want `\QAvoid to access struct fields directly`
	_ = ireq.GetConfig() // OK: use function to access the field
	_ = ireq.GetVector() // OK: use function to access the field

	loc.Name = "newname" // OK: it is used in LHS
	ic := &payload.Insert_Config{}
	ireq.Config = ic // OK: it is used in LHS

	ireq.Config, _ = ic, "test"    // OK: it is used in LHS
	ireq.GetConfig().Timestamp = 0 // OK: it is used in LHS
	ireq.Config.Timestamp = 0      // want `\QAvoid to access struct fields directly`

	if loc.Name != "" { // want `\QAvoid to access struct fields directly`
	}

	if loc.Name == "" { // want `\QAvoid to access struct fields directly`
	} else {
	}

	if loc != nil && loc.Name != "" { // want `\QAvoid to access struct fields directly`
	}
	if loc != nil && loc.GetName() != "" { // OK: use function to access the field
	}

	if ireq != nil && ireq.Vector.Id != "" { // want `\QAvoid to access struct fields directly` `\QAvoid to access struct fields directly`
	}

	if ireq != nil && ireq.Vector.GetId() != "" { // want `\QAvoid to access struct fields directly`
	}
	if ireq != nil && ireq.GetVector().GetId() != "" { // OK: use function to access the field
	}

	pods := &payload.Info_Pods{}
	if pods.Pods[0].GetNode() != nil { // want `\QAvoid to access struct fields directly`
		pods.Pods[0].GetNode().Pods = nil // want `\QAvoid to access struct fields directly`
	}
	if pods.GetPods()[0].GetNode() != nil { // OK: use function to access the field
		pods.GetPods()[0].GetNode().Pods = nil // OK: use function to access the field
	}

	locs := &payload.Object_Locations{}
	_ = append([]*payload.Object_Location{}, locs.Locations...)          // want `\QAvoid to access struct fields directly`
	_ = append(locs.GetLocations(), locs.Locations...)                   // want `\QAvoid to access struct fields directly`
	locs.Locations = append(locs.GetLocations(), locs.Locations...)      // want `\QAvoid to access struct fields directly`
	_ = append([]*payload.Object_Location{}, locs.GetLocations()...)     // OK: use function to access the field
	locs.Locations = append(locs.GetLocations(), locs.GetLocations()...) // OK: use function to access the field

	dmy := &dummy{}
	_ = dmy.field1 // OK: dummy is not a gRPC payload object
}

func printFmts() {
	fmt.Println("%v", "") // want `\Qfound formatting directive in non-formatting call`
	fmt.Print("%s", "")   // want `\Qfound formatting directive in non-formatting call`

	fmt.Printf("%s", "") // OK: use Printf
}
