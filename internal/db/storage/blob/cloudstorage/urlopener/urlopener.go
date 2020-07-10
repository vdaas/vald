package urlopener

import (
	"context"
	"io/ioutil"
	"net/http"
	"reflect"
	"unsafe"

	"github.com/vdaas/vald/internal/errors"
	"gocloud.dev/blob/gcsblob"
	"gocloud.dev/gcp"
	"golang.org/x/oauth2/google"
)

// URLOpener returns gcsblob.URLOpener and error.
type URLOpener interface {
	URLOpener(context.Context) (*gcsblob.URLOpener, error)
}

type urlOpener struct {
	googleAccessID string
	privateKey     []byte
	signBytes      func([]byte) ([]byte, error)
	makeSignBytes  func(requestCtx context.Context) gcsblob.SignBytesFunc

	credentialsFilePath string
	credentialsJSON     string

	client *http.Client
}

// New returns URLOpener implementation.
func New(opts ...Option) (URLOpener, error) {
	uo := new(urlOpener)
	for _, opt := range opts {
		if err := opt(uo); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}
	return uo, nil
}

func (uo *urlOpener) URLOpener(ctx context.Context) (guo *gcsblob.URLOpener, err error) {
	var creds *google.Credentials

	switch {
	case len(uo.credentialsFilePath) != 0:
		data, err := ioutil.ReadFile(uo.credentialsFilePath)
		if err != nil {
			return nil, err
		}
		creds, err = google.CredentialsFromJSON(ctx, data)
		if err != nil {
			return nil, err
		}
	case len(uo.credentialsJSON) != 0:
		data := *(*[]byte)(unsafe.Pointer(&uo.credentialsJSON))
		creds, err = google.CredentialsFromJSON(ctx, data)
		if err != nil {
			return nil, err
		}
	default:
		creds, err = google.FindDefaultCredentials(ctx)
		if err != nil {
			return nil, err
		}
	}

	client, err := gcp.NewHTTPClient(
		uo.client.Transport,
		gcp.CredentialsTokenSource(creds),
	)
	if err != nil {
		return nil, err
	}

	return &gcsblob.URLOpener{
		Client: client,
		Options: gcsblob.Options{
			GoogleAccessID: uo.googleAccessID,
			PrivateKey:     uo.privateKey,
			SignBytes:      uo.signBytes,
			MakeSignBytes:  uo.makeSignBytes,
		},
	}, nil
}
