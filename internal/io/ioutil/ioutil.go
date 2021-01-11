//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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

// Package ioutil provides utility function for I/O
package ioutil

import (
	"bytes"
	"os"

	"github.com/vdaas/vald/internal/safety"
)

func ReadFile(path string) ([]byte, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var n int64 = bytes.MinRead
	if fi, err := f.Stat(); err == nil {
		if size := fi.Size() + bytes.MinRead; size > n {
			n = size
		}
	}

	buf := bytes.NewBuffer(make([]byte, 0, n))

	err = safety.RecoverFunc(func() (err error) {
		_, err = buf.ReadFrom(f)
		return err
	})()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
