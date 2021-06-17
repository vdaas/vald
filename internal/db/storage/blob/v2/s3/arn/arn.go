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
package arn

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
)

func ValidateSupportedARNType(bucket string) (err error) {
	if !arn.IsARN(bucket) {
		return nil
	}

	parsedARN, err := arn.Parse(bucket)
	if err != nil {
		return err
	}

	if parsedARN.Service == "s3-object-lambda" {
		return fmt.Errorf("manager does not support s3-object-lambda service ARNs")
	}

	return nil
}
