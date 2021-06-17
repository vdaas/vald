//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package errors provides error types and function
package errors

import "fmt"

var (
	ErrInvalidStorageType = New("invalid storage type")

	ErrStorageReaderNotOpened = New("reader not opened")
	ErrStorageWriterNotOpened = New("writer not opened")

	ErrBucketNotOpened = New("bucket not opened")
)

func NewErrS3ReadingBody(err error) error {
	return &ErrS3ReadingBody{
		err: err,
	}
}

type ErrS3ReadingBody struct {
	err error
}

func (e *ErrS3ReadingBody) Error() string {
	return fmt.Sprintf("failed to read part body: %v", e.err)
}

func (e *ErrS3ReadingBody) Unwrap() error {
	return e.err
}

func NewErrMultiUpload(uploadID string, err error) error {
	return &ErrMultiUpload{
		err:      err,
		uploadID: uploadID,
	}
}

type ErrMultiUpload struct {
	err      error
	uploadID string
}

func (m *ErrMultiUpload) Error() string {
	var extra string
	if m.err != nil {
		extra = fmt.Sprintf(", cause: %s", m.err.Error())
	}
	return fmt.Sprintf("upload multipart failed, upload id: %s%s", m.uploadID, extra)
}

func (m *ErrMultiUpload) Unwrap() error {
	return m.err
}

// UploadID returns the id of the S3 upload which failed.
func (m *ErrMultiUpload) UploadID() string {
	return m.uploadID
}
