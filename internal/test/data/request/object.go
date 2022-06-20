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
package request

import (
	"strconv"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
)

// GenObjectLocations generate ObjectLocations payload with multiple name and ip with generated uuid.
func GenObjectLocations(num int, name string, ipAddr string) *payload.Object_Locations {
	result := &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, num),
	}

	for i := 0; i < num; i++ {
		result.Locations[i] = &payload.Object_Location{
			Name: name,
			Uuid: "uuid-" + strconv.Itoa(i+1),
			Ips:  []string{ipAddr},
		}
	}
	return result
}

// GenObjectStreamLocation generate ObjectStreamLocations payload with multiple name and ip with generated uuid.
func GenObjectStreamLocation(num int, name string, ipAddr string) []*payload.Object_StreamLocation {
	result := make([]*payload.Object_StreamLocation, num)

	for i := 0; i < num; i++ {
		result[i] = &payload.Object_StreamLocation{
			Payload: &payload.Object_StreamLocation_Location{
				Location: &payload.Object_Location{
					Name: name,
					Uuid: "uuid-" + strconv.Itoa(i+1),
					Ips:  []string{ipAddr},
				},
			},
		}
	}
	return result
}
