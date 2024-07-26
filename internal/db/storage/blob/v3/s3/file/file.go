// Package file provides temporary file functions for download and upload.
package file

import (
	"os"

	"github.com/vdaas/vald/internal/log"
)

type File struct {
	*os.File
}

func (f *File) Close() (err error) {
	err = f.File.Close()
	if err != nil {
		log.Warnf("failed to close file: %v", err)
	}

	err = os.Remove(f.Name())
	if err != nil {
		return err
	}
	return nil
}

func CreateTemp() (*File, error) {
	f, err := os.CreateTemp("", "")
	if err != nil {
		return nil, err
	}

	return &File{
		File: f,
	}, nil
}
