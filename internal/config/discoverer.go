//
// Copyright (C) 2019 kpango (Yusuke Kato)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package config providers configuration type and load configuration logic
package config

import "fmt"

type Discoverer struct {
	Addr     string
	Host     string
	Port     int
	Duration string
}

func (d *Discoverer) Bind() *Discoverer {
	d.Host = GetActualValue(d.Host)
	d.Addr = GetActualValue(d.Addr)
	if len(d.Addr) == 0 && len(d.Host) != 0 {
		d.Addr = fmt.Sprintf("%s:%d", d.Host, d.Port)
	}
	return d
}
