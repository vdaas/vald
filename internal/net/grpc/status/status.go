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

// Package status provides statuses and errors returned by grpc handler functions
package status

import (
	"cmp"
	"context"
	"os"
	"slices"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/proto"
	"github.com/vdaas/vald/internal/net/grpc/types"
	"github.com/vdaas/vald/internal/strings"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/status"
)

type (
	Status = status.Status
	Code   = codes.Code
)

func New(c codes.Code, msg string) *Status {
	return status.New(c, msg)
}

func newStatus(code codes.Code, msg string, err error, details ...any) (st *Status) {
	st = New(code, msg)
	return withDetails(st, err, details...)
}

func WrapWithCanceled(msg string, err error, details ...any) error {
	return newStatus(codes.Canceled, msg, err, details...).Err()
}

func WrapWithUnknown(msg string, err error, details ...any) error {
	return newStatus(codes.Unknown, msg, err, details...).Err()
}

func WrapWithInvalidArgument(msg string, err error, details ...any) error {
	return newStatus(codes.InvalidArgument, msg, err, details...).Err()
}

func WrapWithDeadlineExceeded(msg string, err error, details ...any) error {
	return newStatus(codes.DeadlineExceeded, msg, err, details...).Err()
}

func WrapWithNotFound(msg string, err error, details ...any) error {
	return newStatus(codes.NotFound, msg, err, details...).Err()
}

func WrapWithAlreadyExists(msg string, err error, details ...any) error {
	return newStatus(codes.AlreadyExists, msg, err, details...).Err()
}

func WrapWithPermissionDenied(msg string, err error, details ...any) error {
	return newStatus(codes.PermissionDenied, msg, err, details...).Err()
}

func WrapWithResourceExhausted(msg string, err error, details ...any) error {
	return newStatus(codes.ResourceExhausted, msg, err, details...).Err()
}

func WrapWithFailedPrecondition(msg string, err error, details ...any) error {
	return newStatus(codes.FailedPrecondition, msg, err, details...).Err()
}

func WrapWithAborted(msg string, err error, details ...any) error {
	return newStatus(codes.Aborted, msg, err, details...).Err()
}

func WrapWithOutOfRange(msg string, err error, details ...any) error {
	return newStatus(codes.OutOfRange, msg, err, details...).Err()
}

func WrapWithUnimplemented(msg string, err error, details ...any) error {
	return newStatus(codes.Unimplemented, msg, err, details...).Err()
}

func WrapWithInternal(msg string, err error, details ...any) error {
	return newStatus(codes.Internal, msg, err, details...).Err()
}

func WrapWithUnavailable(msg string, err error, details ...any) error {
	return newStatus(codes.Unavailable, msg, err, details...).Err()
}

func WrapWithDataLoss(msg string, err error, details ...any) error {
	return newStatus(codes.DataLoss, msg, err, details...).Err()
}

func WrapWithUnauthenticated(msg string, err error, details ...any) error {
	return newStatus(codes.Unauthenticated, msg, err, details...).Err()
}

func CreateWithNotFound(msg string, err error, details ...any) *Status {
	return newStatus(codes.NotFound, msg, err, details...)
}

func Error(code codes.Code, msg string) error {
	return status.Error(code, msg)
}

func Errorf(code codes.Code, format string, args ...any) error {
	return status.Errorf(code, format, args...)
}

func ParseError(
	err error, defaultCode codes.Code, defaultMsg string, details ...any,
) (st *Status, msg string, rerr error) {
	if err == nil {
		st = newStatus(codes.OK, "", nil, details...)
		msg = st.Message()
		return st, msg, nil
	}
	var ok bool
	st, ok = FromError(err)
	if !ok || st == nil {
		if defaultCode == 0 {
			defaultCode = codes.Internal
		}
		if len(defaultMsg) == 0 {
			defaultMsg = "failed to parse gRPC status from error"
		}
		st = newStatus(defaultCode, defaultMsg, err, details...)
		if st != nil {
			if st.Message() != "" {
				msg = st.Message()
			} else if e := st.Err(); e != nil {
				msg = e.Error()
			} else {
				msg = defaultMsg
			}
			err = st.Err()
		} else {
			msg = defaultMsg
		}
		return st, msg, err
	}

	sst := withDetails(st, err, details...)
	if sst != nil {
		return sst, sst.Message(), sst.Err()
	}
	return st, st.Message(), st.Err()
}

func FromError(err error) (st *Status, ok bool) {
	if err == nil {
		return nil, false
	}
	root := err
	possibleStatus := func() (st *Status, ok bool) {
		switch {
		case errors.Is(root, context.DeadlineExceeded):
			st = newStatus(codes.DeadlineExceeded, root.Error(), errors.Unwrap(root))
			return st, true
		case errors.Is(root, context.Canceled):
			st = newStatus(codes.Canceled, root.Error(), errors.Unwrap(root))
			return st, true
		}
		st = newStatus(codes.Unknown, root.Error(), errors.Unwrap(root))
		return st, false
	}
	for {
		if st, ok = status.FromError(err); ok && st != nil {
			if st.Code() == codes.Unknown {
				switch x := err.(type) {
				case interface{ Unwrap() error }:
					err = x.Unwrap()
					if err == nil {
						return st, true
					}
					sst, ok := FromError(err)
					if ok && sst != nil {
						return sst, true
					}
				case interface{ Unwrap() []error }:
					for _, err = range x.Unwrap() {
						if sst, ok := FromError(err); ok && sst != nil {
							if sst.Code() != codes.Unknown {
								return sst, true
							}
						}
					}
				}
			}
			return st, true
		}

		switch x := err.(type) {
		case interface{ Unwrap() error }:
			err = x.Unwrap()
			if err == nil {
				return possibleStatus()
			}
		case interface{ Unwrap() []error }:
			errs := x.Unwrap()
			if errs == nil {
				return possibleStatus()
			}
			var prev *Status
			for _, err = range errs {
				if st, ok = FromError(err); ok && st != nil {
					if st.Code() != codes.Unknown {
						return st, true
					}
					prev = st
				}
			}
			if prev != nil {
				return prev, true
			}
			return possibleStatus()
		default:
			return possibleStatus()

		}
	}
}

func withDetails(st *Status, err error, details ...any) *Status {
	if st != nil {
		details = append(st.Details(), details...)
	}
	dmap := make(map[string][]proto.Message, len(details)+1)
	if err != nil {
		typeName := errdetails.ErrorInfoMessageName
		dmap[typeName] = []proto.Message{&errdetails.ErrorInfo{
			Reason: err.Error(),
			Domain: func() (hostname string) {
				var err error
				hostname, err = os.Hostname()
				if err != nil {
					log.Warn("failed to fetch hostname:", err)
				}
				return hostname
			}(),
		}}
	}
	for _, detail := range details {
		if detail == nil {
			continue
		}
		var (
			typeName string
			msg      proto.Message
		)
		switch v := detail.(type) {
		case *spb.Status:
			if v != nil {
				for _, d := range v.GetDetails() {
					typeName = d.GetTypeUrl()
					if typeName != "" {
						msg = errdetails.AnyToErrorDetail(d)
						if msg != nil {
							dm, ok := dmap[typeName]
							if ok && dm != nil {
								dmap[typeName] = append(dm, msg)
							} else {
								dmap[typeName] = []proto.Message{msg}
							}
						}
					}
				}
			}
		case spb.Status:
			for _, d := range v.GetDetails() {
				typeName = d.GetTypeUrl()
				if typeName != "" {
					msg = errdetails.AnyToErrorDetail(d)
					if msg != nil {
						d, ok := dmap[typeName]
						if ok && d != nil {
							dmap[typeName] = append(d, msg)
						} else {
							dmap[typeName] = []proto.Message{msg}
						}
					}
				}
			}
		case *status.Status:
			if v != nil {
				for _, d := range v.Proto().GetDetails() {
					typeName = d.GetTypeUrl()
					if typeName != "" {
						msg = errdetails.AnyToErrorDetail(d)
						if msg != nil {
							d, ok := dmap[typeName]
							if ok && d != nil {
								dmap[typeName] = append(d, msg)
							} else {
								dmap[typeName] = []proto.Message{msg}
							}
						}
					}
				}
			}
		case status.Status:
			for _, d := range v.Proto().GetDetails() {
				typeName = d.GetTypeUrl()
				if typeName != "" {
					msg = errdetails.AnyToErrorDetail(d)
					if msg != nil {
						d, ok := dmap[typeName]
						if ok && d != nil {
							dmap[typeName] = append(d, msg)
						} else {
							dmap[typeName] = []proto.Message{msg}
						}
					}
				}
			}
		case *info.Detail:
			if v != nil {
				typeName = errdetails.DebugInfoMessageName
				msg = errdetails.DebugInfoFromInfoDetail(v)
				if msg != nil {
					d, ok := dmap[typeName]
					if ok && d != nil {
						dmap[typeName] = append(d, msg)
					} else {
						dmap[typeName] = []proto.Message{msg}
					}
				}
			}
		case info.Detail:
			typeName = errdetails.DebugInfoMessageName
			msg = errdetails.DebugInfoFromInfoDetail(&v)
			if msg != nil {
				d, ok := dmap[typeName]
				if ok && d != nil {
					dmap[typeName] = append(d, msg)
				} else {
					dmap[typeName] = []proto.Message{msg}
				}
			}
		case *types.Any:
			if v != nil {
				typeName = v.GetTypeUrl()
				msg = errdetails.AnyToErrorDetail(v)
				if msg != nil {
					d, ok := dmap[typeName]
					if ok && d != nil {
						dmap[typeName] = append(d, msg)
					} else {
						dmap[typeName] = []proto.Message{msg}
					}
				}
			}
		case types.Any:
			typeName = v.GetTypeUrl()
			msg = errdetails.AnyToErrorDetail(&v)
			if msg != nil {
				d, ok := dmap[typeName]
				if ok && d != nil {
					dmap[typeName] = append(d, msg)
				} else {
					dmap[typeName] = []proto.Message{msg}
				}
			}
		case *proto.Message:
			if v != nil {
				typeName = typeURL(*v)
				d, ok := dmap[typeName]
				if ok && d != nil {
					dmap[typeName] = append(d, *v)
				} else {
					dmap[typeName] = []proto.Message{*v}
				}
			}
		case proto.Message:
			typeName = typeURL(v)
			d, ok := dmap[typeName]
			if ok && d != nil {
				dmap[typeName] = append(d, v)
			} else {
				dmap[typeName] = []proto.Message{v}
			}
		case *proto.MessageV1:
			if v != nil {
				msg = proto.ToMessageV2(*v)
				typeName = typeURL(msg)
				if msg != nil {
					d, ok := dmap[typeName]
					if ok && d != nil {
						dmap[typeName] = append(d, msg)
					} else {
						dmap[typeName] = []proto.Message{msg}
					}
				}
			}
		case proto.MessageV1:
			msg = proto.ToMessageV2(v)
			typeName = typeURL(msg)
			if msg != nil {
				d, ok := dmap[typeName]
				if ok && d != nil {
					dmap[typeName] = append(d, msg)
				} else {
					dmap[typeName] = []proto.Message{msg}
				}
			}
		}
	}
	msgs := make([]proto.MessageV1, 0, len(dmap))
	visited := make(map[string]bool, len(dmap))
	for typeName, ds := range dmap {
		switch typeName {
		case errdetails.DebugInfoMessageName:
			m := new(errdetails.DebugInfo)
			for _, msg := range ds {
				d, ok := msg.(*errdetails.DebugInfo)
				if ok && d != nil && !visited[d.String()] {
					visited[d.String()] = true
					if m.GetDetail() == "" {
						m.Detail = d.GetDetail()
					} else if m.GetDetail() != d.GetDetail() && !strings.Contains(m.GetDetail(), d.GetDetail()) {
						m.Detail += "\t" + d.GetDetail()
					}
					if len(m.GetStackEntries()) < len(d.GetStackEntries()) {
						m.StackEntries = d.GetStackEntries()
					}
				}
			}
			m.Detail = removeDuplicatesFromTSVLine(m.GetDetail())
			msgs = append(msgs, m)
		case errdetails.ErrorInfoMessageName:
			m := new(errdetails.ErrorInfo)
			for _, msg := range ds {
				e, ok := msg.(*errdetails.ErrorInfo)
				if ok && e != nil && !visited[e.String()] && !visited[e.GetReason()] {
					visited[e.String()] = true
					visited[e.GetReason()] = true
					if m.GetDomain() == "" {
						m.Domain = e.GetDomain()
					} else if m.GetDomain() != e.GetDomain() && !strings.Contains(m.GetDomain(), e.GetDomain()) {
						m.Domain += "\t" + e.GetDomain()
					}
					if m.GetReason() == "" {
						m.Reason += e.GetReason()
					} else if m.GetReason() != e.GetReason() && !strings.Contains(m.GetReason(), e.GetReason()) {
						m.Reason += "\t" + e.GetReason()
					}
					if e.GetMetadata() != nil {
						if m.GetMetadata() == nil {
							m.Metadata = e.GetMetadata()
						} else {
							m.Metadata = appendM(m.GetMetadata(), e.GetMetadata())
						}
					}
				}
			}
			m.Reason = removeDuplicatesFromTSVLine(m.GetReason())
			m.Domain = removeDuplicatesFromTSVLine(m.GetDomain())
			msgs = append(msgs, m)
		case errdetails.BadRequestMessageName:
			m := new(errdetails.BadRequest)
			for _, msg := range ds {
				b, ok := msg.(*errdetails.BadRequest)
				if ok && b != nil && b.GetFieldViolations() != nil && !visited[b.String()] {
					visited[b.String()] = true
					if m.GetFieldViolations() == nil {
						m = b
					} else {
						m.FieldViolations = append(m.GetFieldViolations(), b.GetFieldViolations()...)
					}
				}
			}
			slices.SortFunc(m.FieldViolations, func(left, right *errdetails.BadRequestFieldViolation) int {
				return cmp.Compare(left.GetField(), right.GetField())
			})
			m.FieldViolations = slices.CompactFunc(m.GetFieldViolations(), func(left, right *errdetails.BadRequestFieldViolation) bool {
				return left.GetField() == right.GetField()
			})
			msgs = append(msgs, m)
		case errdetails.BadRequestFieldViolationMessageName:
			m := new(errdetails.BadRequestFieldViolation)
			for _, msg := range ds {
				b, ok := msg.(*errdetails.BadRequestFieldViolation)
				if ok && b != nil && !visited[b.String()] {
					visited[b.String()] = true
					if m.GetField() == "" {
						m.Field = b.GetField()
					} else if m.GetField() != b.GetField() && !strings.Contains(m.GetField(), b.GetField()) {
						m.Field += "\t" + b.GetField()
					}
					if m.GetDescription() == "" {
						m.Description = b.GetDescription()
					} else if m.GetDescription() != b.GetDescription() && !strings.Contains(m.GetDescription(), b.GetDescription()) {
						m.Description += "\t" + b.GetDescription()
					}
				}
			}
			msgs = append(msgs, m)
		case errdetails.LocalizedMessageMessageName:
			m := new(errdetails.LocalizedMessage)
			for _, msg := range ds {
				l, ok := msg.(*errdetails.LocalizedMessage)
				if ok && l != nil && !visited[l.String()] {
					visited[l.String()] = true
					if m.GetLocale() == "" {
						m.Locale = l.GetLocale()
					} else if m.GetLocale() != l.GetLocale() && !strings.Contains(m.GetLocale(), l.GetLocale()) {
						m.Locale += "\t" + l.GetLocale()
					}
					if m.GetMessage() == "" {
						m.Message = l.GetMessage()
					} else if m.GetMessage() != l.GetMessage() && !strings.Contains(m.GetMessage(), l.GetMessage()) {
						m.Message += "\t" + l.GetMessage()
					}
				}
			}
			msgs = append(msgs, m)
		case errdetails.PreconditionFailureMessageName:
			m := new(errdetails.PreconditionFailure)
			for _, msg := range ds {
				p, ok := msg.(*errdetails.PreconditionFailure)
				if ok && p != nil && p.GetViolations() != nil && !visited[p.String()] {
					visited[p.String()] = true
					if m.GetViolations() == nil {
						m = p
					} else {
						m.Violations = append(m.GetViolations(), p.GetViolations()...)
					}
				}
			}
			slices.SortFunc(m.Violations, func(left, right *errdetails.PreconditionFailureViolation) int {
				return cmp.Compare(left.GetType(), right.GetType())
			})
			m.Violations = slices.CompactFunc(m.GetViolations(), func(left, right *errdetails.PreconditionFailureViolation) bool {
				return left.GetType() == right.GetType()
			})
			msgs = append(msgs, m)
		case errdetails.PreconditionFailureViolationMessageName:
			m := new(errdetails.PreconditionFailureViolation)
			for _, msg := range ds {
				p, ok := msg.(*errdetails.PreconditionFailureViolation)
				if ok && p != nil && !visited[p.String()] {
					visited[p.String()] = true
					if m.GetType() == "" {
						m.Type = p.GetType()
					} else if m.GetType() != p.GetType() && !strings.Contains(m.GetType(), p.GetType()) {
						m.Type += "\t" + p.GetType()
					}
					if m.GetSubject() == "" {
						m.Subject = p.GetSubject()
					} else if m.GetSubject() != p.GetSubject() && !strings.Contains(m.GetSubject(), p.GetSubject()) {
						m.Subject += "\t" + p.GetSubject()
					}
					if m.GetDescription() == "" {
						m.Description = p.GetDescription()
					} else if m.GetDescription() != p.GetDescription() && !strings.Contains(m.GetDescription(), p.GetDescription()) {
						m.Description += "\t" + p.GetDescription()
					}
				}
			}
			msgs = append(msgs, m)
		case errdetails.HelpMessageName:
			m := new(errdetails.Help)
			for _, msg := range ds {
				h, ok := msg.(*errdetails.Help)
				if ok && h != nil && h.GetLinks() != nil && !visited[h.String()] {
					visited[h.String()] = true
					if m.GetLinks() == nil {
						m = h
					} else {
						m.Links = append(m.GetLinks(), h.GetLinks()...)
					}
				}
			}
			slices.SortFunc(m.Links, func(left, right *errdetails.HelpLink) int {
				return cmp.Compare(left.GetUrl(), right.GetUrl())
			})
			m.Links = slices.CompactFunc(m.GetLinks(), func(left, right *errdetails.HelpLink) bool {
				return left.GetUrl() == right.GetUrl()
			})
			msgs = append(msgs, m)
		case errdetails.HelpLinkMessageName:
			m := new(errdetails.HelpLink)
			for _, msg := range ds {
				h, ok := msg.(*errdetails.HelpLink)
				if ok && h != nil && !visited[h.String()] {
					visited[h.String()] = true
					if m.GetUrl() == "" {
						m.Url = h.GetUrl()
					} else if m.GetUrl() != h.GetUrl() && !strings.Contains(m.GetUrl(), h.GetUrl()) {
						m.Url += "\t" + h.GetUrl()
					}
					if m.GetDescription() == "" {
						m.Description = h.GetDescription()
					} else if m.GetDescription() != h.GetDescription() && !strings.Contains(m.GetDescription(), h.GetDescription()) {
						m.Description += "\t" + h.GetDescription()
					}
				}
			}
			msgs = append(msgs, m)
		case errdetails.QuotaFailureMessageName:
			m := new(errdetails.QuotaFailure)
			for _, msg := range ds {
				q, ok := msg.(*errdetails.QuotaFailure)
				if ok && q != nil && q.GetViolations() != nil && !visited[q.String()] {
					visited[q.String()] = true
					if m.GetViolations() == nil {
						m = q
					} else {
						m.Violations = append(m.GetViolations(), q.GetViolations()...)
					}
				}
			}
			slices.SortFunc(m.Violations, func(left, right *errdetails.QuotaFailureViolation) int {
				return cmp.Compare(left.GetSubject(), right.GetSubject())
			})
			m.Violations = slices.CompactFunc(m.GetViolations(), func(left, right *errdetails.QuotaFailureViolation) bool {
				return left.GetSubject() == right.GetSubject()
			})
			msgs = append(msgs, m)
		case errdetails.QuotaFailureViolationMessageName:
			m := new(errdetails.QuotaFailureViolation)
			for _, msg := range ds {
				q, ok := msg.(*errdetails.QuotaFailureViolation)
				if ok && q != nil && !visited[q.String()] {
					visited[q.String()] = true
					if m.GetSubject() == "" {
						m.Subject = q.GetSubject()
					} else if m.GetSubject() != q.GetSubject() && !strings.Contains(m.GetSubject(), q.GetSubject()) {
						m.Subject += "\t" + q.GetSubject()
					}
					if m.GetDescription() == "" {
						m.Description = q.GetDescription()
					} else if m.GetDescription() != q.GetDescription() && !strings.Contains(m.GetDescription(), q.GetDescription()) {
						m.Description += "\t" + q.GetDescription()
					}
				}
			}
			msgs = append(msgs, m)
		case errdetails.RequestInfoMessageName:
			m := new(errdetails.RequestInfo)
			for _, msg := range ds {
				r, ok := msg.(*errdetails.RequestInfo)
				if ok && r != nil && !visited[r.String()] {
					visited[r.String()] = true
					if m.GetRequestId() == "" {
						m.RequestId = r.GetRequestId()
					} else if m.GetRequestId() != r.GetRequestId() && !strings.Contains(m.GetRequestId(), r.GetRequestId()) {
						m.RequestId += "\t" + r.GetRequestId()
					}
					if m.GetServingData() == "" {
						m.ServingData = r.GetServingData()
					} else if m.GetServingData() != r.GetServingData() && !strings.Contains(m.GetServingData(), r.GetServingData()) {
						m.ServingData += "\t" + r.GetServingData()
					}
				}
			}
			msgs = append(msgs, m)
		case errdetails.ResourceInfoMessageName:
			m := new(errdetails.ResourceInfo)
			for _, msg := range ds {
				r, ok := msg.(*errdetails.ResourceInfo)
				if ok && r != nil && !visited[r.String()] {
					visited[r.String()] = true
					if m.GetResourceType() == "" {
						m.ResourceType = r.GetResourceType()
					} else if m.GetResourceType() != r.GetResourceType() && len(m.GetResourceType()) < len(r.GetResourceType()) {
						m.ResourceType += r.GetResourceType()
					}
					if m.GetResourceName() == "" {
						m.ResourceName = r.GetResourceName()
					} else if m.GetResourceName() != r.GetResourceName() && !strings.Contains(m.GetResourceName(), r.GetResourceName()) {
						m.ResourceName += "\t" + r.GetResourceName()
					}
					if m.GetDescription() == "" {
						m.Description = r.GetDescription()
					} else if m.GetDescription() != r.GetDescription() && !strings.Contains(m.GetDescription(), r.GetDescription()) {
						m.Description += "\t" + r.GetDescription()
					}
				}
			}
			msgs = append(msgs, m)
		case errdetails.RetryInfoMessageName:
			m := new(errdetails.RetryInfo)
			for _, msg := range ds {
				r, ok := msg.(*errdetails.RetryInfo)
				if ok && r != nil && !visited[r.String()] {
					visited[r.String()] = true
					if m.GetRetryDelay() == nil || r.GetRetryDelay().Seconds < m.GetRetryDelay().Seconds {
						m.RetryDelay = r.GetRetryDelay()
					}
				}
			}
			msgs = append(msgs, m)
		}
	}
	if st == nil {
		if err != nil {
			st = New(codes.Unknown, err.Error())
		} else {
			st = New(codes.Unknown, "")
		}
	}
	if msgs != nil {
		sst, err := status.New(st.Code(), st.Message()).WithDetails(msgs...)
		if err == nil {
			st = sst
		}
	}
	Log(st.Code(), st.Err())
	return st
}

func typeURL(msg proto.Message) string {
	switch msg.(type) {
	case *errdetails.DebugInfo:
		return errdetails.DebugInfoMessageName
	case *errdetails.ErrorInfo:
		return errdetails.ErrorInfoMessageName
	case *errdetails.BadRequest:
		return errdetails.BadRequestMessageName
	case *errdetails.BadRequestFieldViolation:
		return errdetails.BadRequestFieldViolationMessageName
	case *errdetails.LocalizedMessage:
		return errdetails.LocalizedMessageMessageName
	case *errdetails.PreconditionFailure:
		return errdetails.PreconditionFailureMessageName
	case *errdetails.PreconditionFailureViolation:
		return errdetails.PreconditionFailureViolationMessageName
	case *errdetails.Help:
		return errdetails.HelpMessageName
	case *errdetails.HelpLink:
		return errdetails.HelpLinkMessageName
	case *errdetails.QuotaFailure:
		return errdetails.QuotaFailureMessageName
	case *errdetails.QuotaFailureViolation:
		return errdetails.QuotaFailureViolationMessageName
	case *errdetails.RequestInfo:
		return errdetails.RequestInfoMessageName
	case *errdetails.ResourceInfo:
		return errdetails.ResourceInfoMessageName
	case *errdetails.RetryInfo:
		return errdetails.RetryInfoMessageName
	}
	return "unknown"
}

func appendM[K comparable](maps ...map[K]string) (result map[K]string) {
	if len(maps) == 0 {
		return nil
	}
	result = maps[0]
	for _, m := range maps[1:] {
		for k, v := range m {
			ev, ok := result[k]
			if ok && v != ev && !strings.Contains(v, ev) {
				v += "\t" + ev
			}
			result[k] = v
		}
	}
	return result
}

func removeDuplicatesFromTSVLine(line string) string {
	fields := strings.Split(line, "\t")
	uniqueFields := make(map[string]bool)
	result := make([]string, 0, len(fields))
	for _, field := range fields {
		if !uniqueFields[field] {
			uniqueFields[field] = true
			result = append(result, field)
		}
	}
	return strings.Join(result, "\t")
}

func Log(code codes.Code, err error) {
	if err != nil {
		switch code {
		case codes.Internal,
			codes.DataLoss:
			log.Error(err.Error())
		case codes.Unavailable,
			codes.ResourceExhausted:
			log.Warn(err.Error())
		case codes.FailedPrecondition,
			codes.InvalidArgument,
			codes.OutOfRange,
			codes.Unauthenticated,
			codes.PermissionDenied,
			codes.Unknown:
			log.Debug(err.Error())
		case codes.Aborted,
			codes.Canceled,
			codes.DeadlineExceeded,
			codes.AlreadyExists,
			codes.NotFound,
			codes.OK,
			codes.Unimplemented:
		default:
			log.Warn(errors.ErrGRPCUnexpectedStatusError(code.String(), err))
		}
	}
}
