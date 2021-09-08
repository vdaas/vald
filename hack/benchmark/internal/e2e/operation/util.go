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

func statusError(tb testing.TB, st *status.Status) {
	tb.Helper()
	if st == nil {
		return
	}
	tb.Errorf("code: %d\tmessage: %s\tdetails: %s",
		st.Code(),
		st.Message(),
		errdetails.Serialize(st.Details()))
}
