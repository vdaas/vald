package urlopener

import (
	"context"
	"net/http"

	"gocloud.dev/blob/gcsblob"
)

type Option func(*urlOpener) error

func WithGoogleAccessID(id string) Option {
	return func(uo *urlOpener) error {
		if len(id) != 0 {
			uo.googleAccessID = id
		}
		return nil
	}
}

func WithPrivateKey(b []byte) Option {
	return func(uo *urlOpener) error {
		if len(b) != 0 {
			uo.privateKey = b
		}
		return nil
	}
}

func WithSignBytes(f func([]byte) ([]byte, error)) Option {
	return func(uo *urlOpener) error {
		if f != nil {
			uo.signBytes = f
		}
		return nil
	}
}

func WithMakeSignBytes(f func(requestCtx context.Context) gcsblob.SignBytesFunc) Option {
	return func(uo *urlOpener) error {
		if f != nil {
			uo.makeSignBytes = f
		}
		return nil
	}
}

func WithHTTPClient(c *http.Client) Option {
	return func(uo *urlOpener) error {
		if c != nil {
			uo.client = c
		}
		return nil
	}
}
