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
