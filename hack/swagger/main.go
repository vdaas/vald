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

package swagger

import (
	"encoding/json"
	"os"

	"github.com/vdaas/vald/internal/errors"
)

type Swagger struct {
	// TODO: Ignore `json:"protobufAny"`
	// TODO: Ignore `json:"runtimeStreamError"`
	Paths map[string]map[string]struct { // map[path]map[method]info
		EndpointName string `json:"operationId"`
		Parameters   []struct {
			In       string `json:"in"`
			Name     string `json:"name"`
			Required bool   `json:"required"`
			Type     string `json:"type"`
		} `json:"parameters"`
		Responses map[string]struct { // map[code]schema
			Description string `json:"description"`
			Schema      struct {
				Reference string `json:"$ref"`
			} `json:"schema"`
		} `json:"responses"`
		Tags []string `json:"tags"`
	} `json:"paths"`
}

type Route struct {
	Path     string
	Methods  []string
	FuncName string
}

func Parse(path string) (err error) {
	var (
		// f io.ReadCloser
		f *os.File
		d Swagger
	)
	f, err = os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Wrap(err, f.Close().Error())
	}()

	err = json.NewDecoder(f).Decode(&d)
	if err != nil {
		return err
	}

	routes := make([]Route, 0, len(d.Paths))
	for path, data := range d.Paths {
		route := Route{
			Path: path,
		}

		for method, def := range data {
			route.Methods = append(route.Methods, method)
			route.FuncName = def.EndpointName
		}
		routes = append(routes, route)
	}

	return nil
}
