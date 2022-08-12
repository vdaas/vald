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

// Package hdf5 is load hdf5 file
package hdf5

import (
	"github.com/vdaas/vald/internal/errors"
)

type Option func(d *data) error

var defaultOptions = []Option{
	WithName(FashionMNIST784Euclidean),
	WithFilePath("./data"),
}

func WithNameByString(n string) Option {
	var name DatasetName
	switch n {
	case FashionMNIST784Euclidean.String():
		name = FashionMNIST784Euclidean
	}
	return WithName(name)
}

func WithName(dn DatasetName) Option {
	return func(d *data) error {
		switch dn {
		case FashionMNIST784Euclidean:
			d.name = dn
		default:
			return errors.NewErrInvalidOption("dataname", dn)
		}
		return nil
	}
}

func WithFilePath(f string) Option {
	return func(d *data) error {
		if len(f) != 0 {
			d.path = f
		}
		return nil
	}
}
