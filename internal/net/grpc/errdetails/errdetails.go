//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

// Package errdetails provides error detail for gRPC status
package errdetails

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/vdaas/vald/apis/grpc/v1/rpc/errdetails"
	"github.com/vdaas/vald/internal/encoding/json"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc/proto"
	"github.com/vdaas/vald/internal/net/grpc/types"
	"github.com/vdaas/vald/internal/strings"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/status"
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

const (
	ValdResourceOwner          = "vdaas.org vald team <vald@vdaas.org>"
	ValdGRPCResourceTypePrefix = "github.com/vdaas/vald/apis/grpc/v1"

	typePrefix = "type.googleapis.com/google.rpc."
)

var (
	debugInfoMessageName                    = new(DebugInfo).ProtoReflect().Descriptor().FullName().Name()
	errorInfoMessageName                    = new(ErrorInfo).ProtoReflect().Descriptor().FullName().Name()
	badRequestMessageName                   = new(BadRequest).ProtoReflect().Descriptor().FullName().Name()
	badRequestFieldViolationMessageName     = new(BadRequestFieldViolation).ProtoReflect().Descriptor().FullName().Name()
	localizedMessageMessageName             = new(LocalizedMessage).ProtoReflect().Descriptor().FullName().Name()
	preconditionFailureMessageName          = new(PreconditionFailure).ProtoReflect().Descriptor().FullName().Name()
	preconditionFailureViolationMessageName = new(PreconditionFailureViolation).ProtoReflect().Descriptor().FullName().Name()
	helpMessageName                         = new(Help).ProtoReflect().Descriptor().FullName().Name()
	helpLinkMessageName                     = new(HelpLink).ProtoReflect().Descriptor().FullName().Name()
	quotaFailureMessageName                 = new(QuotaFailure).ProtoReflect().Descriptor().FullName().Name()
	quotaFailureViolationMessageName        = new(QuotaFailureViolation).ProtoReflect().Descriptor().FullName().Name()
	requestInfoMessageName                  = new(RequestInfo).ProtoReflect().Descriptor().FullName().Name()
	resourceInfoMessageName                 = new(ResourceInfo).ProtoReflect().Descriptor().FullName().Name()
	retryInfoMessageName                    = new(RetryInfo).ProtoReflect().Descriptor().FullName().Name()
)

type Detail struct {
	TypeURL string        `json:"type_url,omitempty" yaml:"type_url"`
	Message proto.Message `json:"message,omitempty"  yaml:"message"`
}

func decodeDetails(objs ...interface{}) (details []Detail) {
	if objs == nil {
		return nil
	}
	details = make([]Detail, 0, len(objs))
	for _, obj := range objs {
		if obj == nil {
			continue
		}
		v := reflect.ValueOf(obj)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
			iobjs := make([]interface{}, 0, v.Len())
			for i := 0; i < v.Len(); i++ {
				var val interface{}
				if v.Index(i).Kind() == reflect.Ptr {
					val = v.Index(i).Elem().Interface()
				} else {
					val = v.Index(i).Interface()
				}
				if val != nil {
					iobjs = append(iobjs, val)
				}
			}
			if len(iobjs) == 0 {
				continue
			}
			details = append(details, decodeDetails(iobjs...)...)
			continue
		}
		switch v := obj.(type) {
		case *spb.Status:
			if v != nil {
				details = append(details, Detail{
					TypeURL: string((*v).ProtoReflect().Descriptor().FullName()),
					Message: v,
				})
			}
		case spb.Status:
			details = append(details, Detail{
				TypeURL: string(v.ProtoReflect().Descriptor().FullName()),
				Message: &v,
			})
		case *status.Status:
			if v != nil {
				details = append(details, append([]Detail{
					{
						TypeURL: string(v.Proto().ProtoReflect().Descriptor().FullName()),
						Message: &spb.Status{
							Code:    v.Proto().GetCode(),
							Message: v.Message(),
						},
					},
				}, decodeDetails(v.Proto().Details)...)...)
			}
		case status.Status:
			details = append(details, append([]Detail{
				{
					TypeURL: string(v.Proto().ProtoReflect().Descriptor().FullName()),
					Message: &spb.Status{
						Code:    v.Proto().GetCode(),
						Message: v.Message(),
					},
				},
			}, decodeDetails(v.Proto().Details)...)...)
		case *Detail:
			if v != nil {
				details = append(details, *v)
			}
		case Detail:
			details = append(details, v)
		case *info.Detail:
			if v != nil {
				di := DebugInfoFromInfoDetail(v)
				details = append(details, Detail{
					TypeURL: string(di.ProtoReflect().Descriptor().FullName()),
					Message: di,
				})
			}
		case info.Detail:
			di := DebugInfoFromInfoDetail(&v)
			details = append(details, Detail{
				TypeURL: string(di.ProtoReflect().Descriptor().FullName()),
				Message: di,
			})
		case *types.Any:
			if v != nil {
				details = append(details, Detail{
					TypeURL: v.GetTypeUrl(),
					Message: AnyToErrorDetail(v),
				})
			}
		case types.Any:
			details = append(details, Detail{
				TypeURL: v.GetTypeUrl(),
				Message: AnyToErrorDetail(&v),
			})
		case *proto.Message:
			if v != nil {
				details = append(details, Detail{
					TypeURL: string((*v).ProtoReflect().Descriptor().FullName()),
					Message: *v,
				})
			}
		case proto.Message:
			details = append(details, Detail{
				TypeURL: string(v.ProtoReflect().Descriptor().FullName()),
				Message: v,
			})
		case *proto.MessageV1:
			if v != nil {
				v2 := proto.ToMessageV2(*v)
				details = append(details, Detail{
					TypeURL: string(v2.ProtoReflect().Descriptor().FullName()),
					Message: v2,
				})
			}
		case proto.MessageV1:
			v2 := proto.ToMessageV2(v)
			details = append(details, Detail{
				TypeURL: string(v2.ProtoReflect().Descriptor().FullName()),
				Message: v2,
			})
		}
	}
	return details
}

func Serialize(objs ...interface{}) string {
	var (
		b   []byte
		err error
	)
	msgs := decodeDetails(objs...)
	switch len(msgs) {
	case 0:
		return fmt.Sprint(objs...)
	case 1:
		b, err = json.Marshal(msgs[0])
	default:
		b, err = json.Marshal(msgs)
	}
	if err != nil {
		return fmt.Sprint(objs...)
	}
	return string(b)
}

func AnyToErrorDetail(a *types.Any) proto.Message {
	if a == nil {
		return nil
	}
	var err error
	switch proto.Name(strings.TrimPrefix(a.GetTypeUrl(), typePrefix)) {
	case debugInfoMessageName:
		var m DebugInfo
		err = types.UnmarshalAny(a, &m)
		if err == nil {
			return &m
		}
	case errorInfoMessageName:
		var m ErrorInfo
		err = types.UnmarshalAny(a, &m)
		if err == nil {
			return &m
		}
	case badRequestFieldViolationMessageName:
		var m BadRequestFieldViolation
		err = types.UnmarshalAny(a, &m)
		if err == nil {
			return &m
		}
	case badRequestMessageName:
		var m BadRequest
		err = types.UnmarshalAny(a, &m)
		if err == nil {
			return &m
		}
	case localizedMessageMessageName:
		var m LocalizedMessage
		err = types.UnmarshalAny(a, &m)
		if err == nil {
			return &m
		}
	case preconditionFailureViolationMessageName:
		var m PreconditionFailureViolation
		err = types.UnmarshalAny(a, &m)
		if err == nil {
			return &m
		}
	case preconditionFailureMessageName:
		var m PreconditionFailure
		err = types.UnmarshalAny(a, &m)
		if err == nil {
			return &m
		}
	case helpLinkMessageName:
		var m HelpLink
		err = types.UnmarshalAny(a, &m)
		if err == nil {
			return &m
		}
	case helpMessageName:
		var m Help
		err = types.UnmarshalAny(a, &m)
		if err == nil {
			return &m
		}
	case quotaFailureViolationMessageName:
		var m QuotaFailureViolation
		err = types.UnmarshalAny(a, &m)
		if err == nil {
			return &m
		}
	case quotaFailureMessageName:
		var m QuotaFailure
		err = types.UnmarshalAny(a, &m)
		if err == nil {
			return &m
		}
	case requestInfoMessageName:
		var m RequestInfo
		err = types.UnmarshalAny(a, &m)
		if err == nil {
			return &m
		}
	case resourceInfoMessageName:
		var m ResourceInfo
		err = types.UnmarshalAny(a, &m)
		if err == nil {
			return &m
		}
	case retryInfoMessageName:
		var m RetryInfo
		err = types.UnmarshalAny(a, &m)
		if err == nil {
			return &m
		}
	}
	if err != nil {
		log.Warn(err)
	}
	return a.ProtoReflect().Interface()
}

func DebugInfoFromInfoDetail(v *info.Detail) *DebugInfo {
	debug := &DebugInfo{
		Detail: strings.Join(append(append([]string{
			"Version:", v.Version, ",",
			"Name:", v.ServerName, ",",
			"GitCommit:", v.GitCommit, ",",
			"BuildTime:", v.BuildTime, ",",
			"NGT_Version:", v.NGTVersion, ",",
			"Go_Version:", v.GoVersion, ",",
			"GOARCH:", v.GoArch, ",",
			"GOOS:", v.GoOS, ",",
			"CGO_Enabled:", v.CGOEnabled, ",",
			"BuildCPUInfo: [",
		}, v.BuildCPUInfoFlags...), "]"), " "),
	}
	if debug.GetStackEntries() == nil {
		debug.StackEntries = make([]string, 0, len(v.StackTrace))
	}
	for i, stack := range v.StackTrace {
		debug.StackEntries = append(debug.GetStackEntries(), strings.Join([]string{
			"id:",
			strconv.Itoa(i),
			"stack_trace:",
			stack.String(),
		}, " "))
	}
	return debug
}
