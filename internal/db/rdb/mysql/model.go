//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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

package mysql

import (
	dbr "github.com/gocraft/dbr/v2"
)

// MetaVector is an interface to handle metadata keep in MySQL.
type MetaVector interface {
	GetUUID() string
	GetVector() []byte
	GetMeta() string
	GetIPs() []string
}

type metaVector struct {
	meta   meta
	podIPs []podIP
}

type meta struct {
	ID     int64          `db:"id"`
	UUID   string         `db:"uuid"`
	Vector []byte         `db:"vector"`
	Meta   dbr.NullString `db:"meta"`
}

type podIP struct {
	ID int64  `db:"id"`
	IP string `db:"ip"`
}

// GetUUID returns UUID of metaVector.
func (m *metaVector) GetUUID() string   { return m.meta.UUID }
// GetVector returns Vector of metaVector.
func (m *metaVector) GetVector() []byte { return m.meta.Vector }
// GetMeta returns meta.String of metaVector.
func (m *metaVector) GetMeta() string   { return m.meta.Meta.String }
// GetIPs returns all podIPs which are Vald Agent Pods' IP indexed meta's vector.
func (m *metaVector) GetIPs() []string {
	ips := make([]string, 0, len(m.podIPs))

	for _, ip := range m.podIPs {
		ips = append(ips, ip.IP)
	}

	return ips
}
