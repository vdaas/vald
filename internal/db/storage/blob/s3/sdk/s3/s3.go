// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package s3

import (
	"github.com/aws/aws-sdk-go/service/s3"
)

type (
	// S3 is type alias for s3.S3.
	S3 = s3.S3
	// GetObjectInput is type alias for s3.GetObjectInput.
	GetObjectInput = s3.GetObjectInput
	// GetObjectOutput is type alias for s3.GetObjectOutput.
	GetObjectOutput = s3.GetObjectOutput
)

const (
	// ErrCodeNoSuchBucket is an alias for s3.ErrCodeNoSuchBucket.
	ErrCodeNoSuchBucket = s3.ErrCodeNoSuchBucket
	// ErrCodeNoSuchKey is an alias for s3.ErrCodeNoSuchKey.
	ErrCodeNoSuchKey = s3.ErrCodeNoSuchKey
)
