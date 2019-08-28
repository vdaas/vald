// MIT License
//
// Copyright (c) 2019 kpango (Yusuke Kato)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package config providers configuration type and load configuration logic
package config

type Debug struct {
	// Profile represent profiling the server
	Profile struct {
		Enable bool    `yaml:"enable" json:"enable"`
		Server *Server `yaml:"server" json:"server"`
	} `yaml:"profile" json:"profile"`

	// Log represent the server enable debug log or not.
	Log struct {
		Level string `yaml:"level" json:"level"`
		Mode  string `yaml:"mode" json:"mode"`
	} `yaml:"log" json:"log"`
}

func (d *Debug) Bind() *Debug {
	if d.Profile.Server != nil {
		d.Profile.Server.Bind()
	}
	d.Log.Level = GetActualValue(d.Log.Level)
	d.Log.Mode = GetActualValue(d.Log.Mode)
	return d
}
