//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

// Package config providers configuration type and load configuration logic
package config

import (
	"github.com/vdaas/vald/internal/db/rdb/mysql"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/tls"
)

// MySQL represent the mysql configuration.
type MySQL struct {
	DB                   string `json:"db"                      yaml:"db"`
	Network              string `json:"network"                 yaml:"network"`
	SocketPath           string `json:"socket_path"             yaml:"socket_path"`
	Host                 string `json:"host"                    yaml:"host"`
	Port                 uint16 `json:"port"                    yaml:"port"`
	User                 string `json:"user"                    yaml:"user"`
	Pass                 string `json:"pass"                    yaml:"pass"`
	Name                 string `json:"name"                    yaml:"name"`
	Charset              string `json:"charset"                 yaml:"charset"`
	Timezone             string `json:"timezone"                yaml:"timezone"`
	InitialPingTimeLimit string `json:"initial_ping_time_limit" yaml:"initial_ping_time_limit"`
	InitialPingDuration  string `json:"initial_ping_duration"   yaml:"initial_ping_duration"`
	ConnMaxLifeTime      string `json:"conn_max_life_time"      yaml:"conn_max_life_time"`
	MaxOpenConns         int    `json:"max_open_conns"          yaml:"max_open_conns"`
	MaxIdleConns         int    `json:"max_idle_conns"          yaml:"max_idle_conns"`
	TLS                  *TLS   `json:"tls"                     yaml:"tls"`
	Net                  *Net   `json:"net"                     yaml:"net"`
}

// Bind returns MySQL object whose some string value is filed value or environment value.
func (m *MySQL) Bind() *MySQL {
	if m.TLS != nil {
		m.TLS.Bind()
	} else {
		m.TLS = new(TLS)
	}
	if m.Net != nil {
		m.Net.Bind()
	} else {
		m.Net = new(Net)
	}
	m.DB = GetActualValue(m.DB)
	m.Network = GetActualValue(m.Network)
	m.SocketPath = GetActualValue(m.SocketPath)
	m.Host = GetActualValue(m.Host)
	m.User = GetActualValue(m.User)
	m.Pass = GetActualValue(m.Pass)
	m.Name = GetActualValue(m.Name)
	m.Charset = GetActualValue(m.Charset)
	m.Timezone = GetActualValue(m.Timezone)
	m.ConnMaxLifeTime = GetActualValue(m.ConnMaxLifeTime)
	m.InitialPingTimeLimit = GetActualValue(m.InitialPingTimeLimit)
	m.InitialPingDuration = GetActualValue(m.InitialPingDuration)
	return m
}

// Opts creates and returns the slice with the functional options for the internal mysql package.
// When any errors occur, Opts returns the no functional options and the errors.
func (m *MySQL) Opts() ([]mysql.Option, error) {
	nt := net.NetworkTypeFromString(m.Network)
	if nt == 0 || nt == net.Unknown || strings.EqualFold(nt.String(), net.Unknown.String()) {
		m.Network = net.TCP.String()
	} else {
		m.Network = nt.String()
	}
	opts := []mysql.Option{
		mysql.WithDB(m.DB),
		mysql.WithNetwork(m.Network),
		mysql.WithSocketPath(m.SocketPath),
		mysql.WithHost(m.Host),
		mysql.WithPort(m.Port),
		mysql.WithUser(m.User),
		mysql.WithPass(m.Pass),
		mysql.WithName(m.Name),
		mysql.WithCharset(m.Charset),
		mysql.WithTimezone(m.Timezone),
		mysql.WithInitialPingTimeLimit(m.InitialPingTimeLimit),
		mysql.WithInitialPingDuration(m.InitialPingDuration),
		mysql.WithConnectionLifeTimeLimit(m.ConnMaxLifeTime),
		mysql.WithMaxIdleConns(m.MaxIdleConns),
		mysql.WithMaxOpenConns(m.MaxOpenConns),
	}

	if m.TLS != nil && m.TLS.Enabled {
		tls, err := tls.New(m.TLS.Opts()...)
		if err != nil {
			return nil, err
		}
		opts = append(opts, mysql.WithTLSConfig(tls))
	}

	if m.Net != nil {
		netOpts, err := m.Net.Opts()
		if err != nil {
			return nil, err
		}
		dialer, err := net.NewDialer(netOpts...)
		if err != nil {
			return nil, err
		}
		opts = append(opts, mysql.WithDialer(dialer))
	}

	return opts, nil
}
