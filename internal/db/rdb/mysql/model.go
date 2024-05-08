//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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

package mysql

// Vector is an interface to handle vector keep in MySQL.
type Vector interface {
	GetUUID() string
	GetVector() []byte
	GetIPs() []string
}

type vector struct {
	data   data
	podIPs []podIP
}

type data struct {
	ID     int64  `db:"id"`
	UUID   string `db:"uuid"`
	Vector []byte `db:"vector"`
}

type podIP struct {
	ID int64  `db:"id"`
	IP string `db:"ip"`
}

// GetUUID returns UUID of Vector.
func (v *vector) GetUUID() string { return v.data.UUID }

// GetVector returns Vector of Vector.
func (v *vector) GetVector() []byte { return v.data.Vector }

// GetIPs returns all podIPs which are Vald Agent Pods' IP indexed vector's vector.
func (v *vector) GetIPs() []string {
	ips := make([]string, 0, len(v.podIPs))

	for _, ip := range v.podIPs {
		ips = append(ips, ip.IP)
	}

	return ips
}
