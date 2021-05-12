package urlopener

import (
	"net/http"
)

type Option func(*urlOpener) error

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
