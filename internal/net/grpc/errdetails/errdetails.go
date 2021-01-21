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

// Package errdetails provides error detail for grpc status
package errdetails

import (
	errdetails "github.com/gogo/googleapis/google/rpc"
)

type (
	DebugInfo                    = errdetails.DebugInfo
	ErrorInfo                    = errdetails.ErrorInfo
	BadRequest                   = errdetails.BadRequest
	BadRequestFieldViolation     = errdetails.BadRequest_FieldViolation
	LocalizedMessage             = errdetails.LocalizedMessage
	PreconditionFailure          = errdetails.PreconditionFailure
	PreconditionFailureViolation = errdetails.PreconditionFailure_Violation
	Help                         = errdetails.Help
	HelpLink                     = errdetails.Help_Link
	QuotaFailure                 = errdetails.QuotaFailure
	QuotaFailureViolation        = errdetails.QuotaFailure_Violation
	RequestInfo                  = errdetails.RequestInfo
	ResourceInfo                 = errdetails.ResourceInfo
	RetryInfo                    = errdetails.RetryInfo
)
