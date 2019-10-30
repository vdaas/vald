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

package mysql

type Option func(*mysqlClient) error

var (
	defaultOpts = []Option{}
)

func WithDB(db string) Option {
	return func(m *mysqlClient) error {
		if db != "" {
			m.db = db
		}
		return nil
	}
}

func WithHost(host string) Option {
	return func(m *mysqlClient) error {
		if host != "" {
			m.host = host
		}
		return nil
	}
}

func WithPort(port int) Option {
	return func(m *mysqlClient) error {
		m.port = port
		return nil
	}
}

func WithUser(user string) Option {
	return func(m *mysqlClient) error {
		if user != "" {
			m.user = user
		}
		return nil
	}
}

func WithPass(pass string) Option {
	return func(m *mysqlClient) error {
		if pass != "" {
			m.pass = pass
		}
		return nil
	}
}

func WithName(name string) Option {
	return func(m *mysqlClient) error {
		if name != "" {
			m.name = name
		}
		return nil
	}
}
