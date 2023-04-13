// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package operation

import (
	"testing"

	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/status"
)

func grpcError(tb testing.TB, err error) {
	tb.Helper()
	if err == nil {
		return
	}
	st, serr := status.FromError(err)
	tb.Errorf(
		"error: %v\tserror: %v\tcode: %s\tdetails: %s\tmessage: %s\tstatus-error: %s\tproto: %s",
		err,
		serr,
		st.Code().String(),
		errdetails.Serialize(st.Details()),
		st.Message(),
		st.Err().Error(),
		errdetails.Serialize(st.Proto()),
	)
}

func statusError(tb testing.TB, code int32, message string, details ...interface{}) {
	tb.Helper()
	tb.Errorf("code: %d\tmessage: %s\tdetails: %s",
		code,
		message,
		errdetails.Serialize(details...))
}
