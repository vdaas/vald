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
	"fmt"
	"reflect"
	"strings"

	errdetails "github.com/gogo/googleapis/google/rpc"
	"github.com/gogo/status"
	"github.com/vdaas/vald/internal/encoding/json"
	"github.com/vdaas/vald/internal/net/grpc/proto"
	"github.com/vdaas/vald/internal/net/grpc/types"
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

type ErrorDetails struct {
	Code    int32           `json:"code,omitempty"`
	Message string          `json:"message,omitempty"`
	Error   error           `json:"error,omitempty"`
	Details []*ErrorDetails `json:"details,omitempty"`
}

const (
	ValdResourceOwner          = "vdaas.org vald team <vald@vdaas.org>"
	ValdGRPCResourceTypePrefix = "github.com/vdaas/vald/apis/grpc/v1"

	typePrefix = "type.googleapis.com/"
)

var (
	debugInfoMessageName                    = new(DebugInfo).XXX_MessageName()
	errorInfoMessageName                    = new(ErrorInfo).XXX_MessageName()
	badRequestMessageName                   = new(BadRequest).XXX_MessageName()
	badRequestFieldViolationMessageName     = new(BadRequestFieldViolation).XXX_MessageName()
	localizedMessageMessageName             = new(LocalizedMessage).XXX_MessageName()
	preconditionFailureMessageName          = new(PreconditionFailure).XXX_MessageName()
	preconditionFailureViolationMessageName = new(PreconditionFailureViolation).XXX_MessageName()
	helpMessageName                         = new(Help).XXX_MessageName()
	helpLinkMessageName                     = new(HelpLink).XXX_MessageName()
	quotaFailureMessageName                 = new(QuotaFailure).XXX_MessageName()
	quotaFailureViolationMessageName        = new(QuotaFailureViolation).XXX_MessageName()
	requestInfoMessageName                  = new(RequestInfo).XXX_MessageName()
	resourceInfoMessageName                 = new(ResourceInfo).XXX_MessageName()
	retryInfoMessageName                    = new(RetryInfo).XXX_MessageName()
)

func (e *ErrorDetails) String() string {
	return fmt.Sprintf(
		"code: %d,\tmessage: %s,\terror: %v,details: %v",
		e.Code,
		e.Message,
		e.Error,
		e.Details,
	)
}

func DecodeErrorDetails(objs ...interface{}) (es []*ErrorDetails) {
	if objs == nil {
		return nil
	}
	es = make([]*ErrorDetails, 0, len(objs))
	for _, obj := range objs {
		v := reflect.ValueOf(obj)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
			iobjs := make([]interface{}, 0, v.Len())
			for i := 0; i < v.Len(); i++ {
				val := v.Index(i).Interface()
				iobjs = append(iobjs, val)
			}
			details := DecodeErrorDetails(iobjs...)
			if details == nil {
				continue
			}
			es = append(es, &ErrorDetails{
				Code:    0,
				Message: v.Kind().String(),
				Details: details,
			})
			continue
		}
		switch v := obj.(type) {
		case status.Status:
			es = append(es, &ErrorDetails{
				Code:    int32(v.Code()),
				Message: v.Message(),
				Error:   v.Err(),
				Details: append(DecodeErrorDetails(v.Details()), DecodeErrorDetails(v.Proto())...),
			})
		case *status.Status:
			es = append(es, &ErrorDetails{
				Code:    int32(v.Code()),
				Message: v.Message(),
				Error:   v.Err(),
				Details: append(DecodeErrorDetails(v.Details()), DecodeErrorDetails(v.Proto())...),
			})
		case errdetails.Status:
			es = append(es, &ErrorDetails{
				Code:    v.GetCode(),
				Message: v.GetMessage(),
				Details: DecodeErrorDetails(v.GetDetails()),
			})
		case *errdetails.Status:
			es = append(es, &ErrorDetails{
				Code:    v.GetCode(),
				Message: v.GetMessage(),
				Details: DecodeErrorDetails(v.GetDetails()),
			})
		case *types.Any:
			d, err := DecodeDetail(v)
			if err != nil {
				es = append(es, &ErrorDetails{
					Code:    0,
					Message: fmt.Sprintf("type: %s\tserialization error: %v,\tpayload: %#v", v.GetTypeUrl(), err, objs),
					Error:   err,
				})
			} else {
				es = append(es, &ErrorDetails{
					Code:    0,
					Message: v.GetTypeUrl(),
					Details: DecodeErrorDetails(d),
				})
			}
		case types.Any:
			d, err := DecodeDetail(&v)
			if err != nil {
				es = append(es, &ErrorDetails{
					Code:    0,
					Message: fmt.Sprintf("type: %s\tserialization error: %v,\tpayload: %#v", v.GetTypeUrl(), err, objs),
					Error:   err,
				})
			} else {
				es = append(es, &ErrorDetails{
					Code:    0,
					Message: v.GetTypeUrl(),
					Details: DecodeErrorDetails(d),
				})
			}
		case proto.Message:
			es = append(es, &ErrorDetails{
				Code:    0,
				Message: v.String(),
			})
		case *proto.Message:
			es = append(es, &ErrorDetails{
				Code:    0,
				Message: (*v).String(),
			})
		case ErrorDetails:
			es = append(es, &v)
		case *ErrorDetails:
			es = append(es, v)
		default:
			b, err := json.Marshal(objs)
			if err != nil {
				es = append(es, &ErrorDetails{
					Code:    0,
					Message: fmt.Sprintf("serialization error: %v,\tpayload: %#v", err, objs),
					Error:   err,
				})
			} else {
				es = append(es, &ErrorDetails{
					Code:    0,
					Message: string(b),
				})
			}
		}
	}
	return es
}

func Serialize(objs ...interface{}) string {
	ed := DecodeErrorDetails(objs...)
	var (
		b   []byte
		err error
	)

	switch len(ed) {
	case 0:
		return fmt.Sprint(objs...)
	case 1:
		b, err = json.Marshal(ed[0])
	default:
		b, err = json.Marshal(ed)
	}
	if err != nil {
		msgs := make([]string, 0, len(ed))
		for _, e := range ed {
			msgs = append(msgs, e.String())
		}
		return strings.Join(msgs, ",\t")
	}
	return string(b)
}

func DecodeDetail(detail *types.Any) (data interface{}, err error) {
	switch strings.TrimPrefix(detail.GetTypeUrl(), typePrefix) {
	case debugInfoMessageName:
		var m DebugInfo
		err = types.UnmarshalAny(detail, &m)
		return m, err
	case errorInfoMessageName:
		var m ErrorInfo
		err = types.UnmarshalAny(detail, &m)
		return m, err
	case badRequestFieldViolationMessageName:
		var m BadRequestFieldViolation
		err = types.UnmarshalAny(detail, &m)
		return m, err
	case badRequestMessageName:
		var m BadRequest
		err = types.UnmarshalAny(detail, &m)
		return m, err
	case localizedMessageMessageName:
		var m LocalizedMessage
		err = types.UnmarshalAny(detail, &m)
		return m, err
	case preconditionFailureViolationMessageName:
		var m PreconditionFailureViolation
		err = types.UnmarshalAny(detail, &m)
		return m, err
	case preconditionFailureMessageName:
		var m PreconditionFailure
		err = types.UnmarshalAny(detail, &m)
		return m, err
	case helpLinkMessageName:
		var m HelpLink
		err = types.UnmarshalAny(detail, &m)
		return m, err
	case helpMessageName:
		var m Help
		err = types.UnmarshalAny(detail, &m)
		return m, err
	case quotaFailureViolationMessageName:
		var m QuotaFailureViolation
		err = types.UnmarshalAny(detail, &m)
		return m, err
	case quotaFailureMessageName:
		var m QuotaFailure
		err = types.UnmarshalAny(detail, &m)
		return m, err
	case requestInfoMessageName:
		var m RequestInfo
		err = types.UnmarshalAny(detail, &m)
		return m, err
	case resourceInfoMessageName:
		var m ResourceInfo
		err = types.UnmarshalAny(detail, &m)
		return m, err
	case retryInfoMessageName:
		var m RetryInfo
		err = types.UnmarshalAny(detail, &m)
		return m, err
	}
	return nil, nil
}

func DecodeDetails(details ...*types.Any) (ds []interface{}, err error) {
	ds = make([]interface{}, 0, len(details))
	for _, detail := range details {
		d, err := DecodeDetail(detail)
		if err != nil {
			return nil, err
		}
		ds = append(ds, d)
	}
	return ds, nil
}

func DecodeDetailsToString(details ...*types.Any) (msg string, err error) {
	ds, err := DecodeDetails(details...)
	if err != nil {
		return "", err
	}
	return Serialize(ds...), nil
}
