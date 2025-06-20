//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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

package config

// Client represents the Client configurations.
type Client struct {
	Net       *Net       `json:"net"       yaml:"net"`
	Transport *Transport `json:"transport" yaml:"transport"`
}

// Bind binds the actual data from the Client receiver field.
func (c *Client) Bind() *Client {
	if c.Net == nil {
		c.Net = new(Net)
	}
	if c.Net != nil {
		c.Net.Bind()
	}

	if c.Transport == nil {
		c.Transport = new(Transport)
	}
	if c.Transport != nil {
		c.Transport.Bind()
	}
	return c
}
