package urlopener

import (
	"context"
	"net/http"
	"unsafe"

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

func WithPrivateKey(str string) Option {
	return func(uo *urlOpener) error {
		if len(str) != 0 {
			uo.privateKey = *(*[]byte)(unsafe.Pointer(&str))
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

func WithCredentialsFile(path string) Option {
	return func(uo *urlOpener) error {
		if len(path) != 0 {
			uo.credentialsFilePath = path
		}
		return nil
	}
}

func WithCredentialsJSON(str string) Option {
	return func(uo *urlOpener) error {
		if len(str) != 0 {
			uo.credentialsJSON = str
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
