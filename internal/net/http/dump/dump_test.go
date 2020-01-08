package dump

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/vdaas/vald/internal/errors"
)

func TestRequest(t *testing.T) {
	type args struct {
		values map[string]interface{}
		body   map[string]interface{}
		r      *http.Request
	}

	type test struct {
		name      string
		args      args
		checkFunc func(res interface{}, err error) error
	}

	tests := []test{
		{
			name: "returns object converted to structure",
			args: args{
				r: &http.Request{
					Host:       "hoge",
					RequestURI: "uri",
					URL: &url.URL{
						Scheme: "http",
					},
					Method: http.MethodGet,
					Proto:  "proto",
					Header: http.Header{},
					TransferEncoding: []string{
						"trans1",
					},
					RemoteAddr:    "0.0.0.0",
					ContentLength: 1234,
				},
				body: map[string]interface{}{
					"name": "vald",
				},
				values: map[string]interface{}{
					"version": "1.0.0",
				},
			},
			checkFunc: func(res interface{}, err error) error {
				if err != nil {
					return errors.Errorf("err is not nil. err: %v", err)
				}

				b, err := json.Marshal(res)
				if err != nil {
					return err
				}

				str := `{"host":"hoge","uri":"uri","url":"http:","method":"GET","proto":"proto","header":{},"transfer_encoding":["trans1"],"remote_addr":"0.0.0.0","content_length":1234,"body":{"name":"vald"},"values":{"version":"1.0.0"}}`
				if got, want := string(b), str; got != want {
					return errors.Errorf("response not equals. want: %v, got: %v", want, got)
				}

				return nil
			},
		},
		{
			name: "returns nil and error",
			checkFunc: func(res interface{}, err error) error {
				if got, want := err, errors.ErrInvalidRequest; !errors.Is(got, want) {
					return errors.Errorf("err not equals. want: %v, got: %v", want, got)
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := Request(tt.args.values, tt.args.body, tt.args.r)
			if err := tt.checkFunc(res, err); err != nil {
				t.Error(err)
			}
		})
	}
}
