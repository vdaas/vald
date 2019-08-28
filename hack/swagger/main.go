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
