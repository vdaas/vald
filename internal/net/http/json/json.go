//aaaa

package json

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/francoispqt/gojay"
	"github.com/vdaas/vald/internal/net/http/rest"
)

type RFC7807Error struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
	K8S      struct {
		PodIP    string `json:"pod_ip"`
		PodName  string `json:"pod_name"`
		NodeIP   string `json:"node_ip"`
		NodeName string `json:"node_name"`
	} `json:"k8s"`
	Status int    `json:"status"`
	Error  string `json:"error"`
}

func Encode(w http.ResponseWriter, data interface{}, status int, contentType string) error {
	w.Header().Set(rest.ContentType, contentType)
	w.WriteHeader(status)
	return gojay.NewEncoder(w).Encode(data)
}

func Decode(r *http.Request, data interface{}) (err error) {
	if r != nil && r.Body != nil {
		err = gojay.NewDecoder(r.Body).Decode(data)
		if err != nil {
			return err
		}
		io.Copy(ioutil.Discard, r.Body)
		return r.Body.Close()
	}
	return nil
}

func Handler(w http.ResponseWriter, r *http.Request,
	data interface{}, logic func() (interface{}, error)) (code int, err error) {
	err = Decode(r, &data)
	if err != nil {
		return http.StatusBadRequest, err
	}
	res, err := logic()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	err = Encode(w, res, http.StatusOK, rest.ApplicationJSON+"; "+rest.CharsetUTF8)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func ErrorHandler(w http.ResponseWriter, code int, data RFC7807Error) error {
	return Encode(w, data, code, rest.ProblemJSON+"; "+rest.CharsetUTF8)
}
