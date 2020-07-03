package urlopener

import (
	"context"
	"net/http"

	"github.com/vdaas/vald/internal/errors"
	"gocloud.dev/blob/gcsblob"
	"gocloud.dev/gcp"
)

type URLOpener interface {
	URLOpener() (*gcsblob.URLOpener, error)
}

type urlOpener struct {
	googleAccessID string
	privateKey     []byte
	signBytes      func([]byte) ([]byte, error)
	makeSignBytes  func(requestCtx context.Context) gcsblob.SignBytesFunc

	client *http.Client
}

func New(opts ...Option) (URLOpener, error) {
	uo := new(urlOpener)

	for _, opt := range opts {
		if err := opt(uo); err != nil {
			return nil, errors.Wrap(err, "failed to apply")

		}
	}

	return uo, nil
}

func (uo *urlOpener) URLOpener() (*gcsblob.URLOpener, error) {
	client, err := gcp.NewHTTPClient(uo.client.Transport, nil)
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
