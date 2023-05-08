package target

import (
	"reflect"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
)

func NewMirrorTargetTemplate(opts ...MirrorTargetOption) (*MirrorTarget, error) {
	mt := new(MirrorTarget)
	for _, opt := range append(defaultMirrorTargetOptions, opts...) {
		if err := opt(mt); err != nil {
			oerr := errors.ErrOptionFailed(err, reflect.ValueOf(opt))
			e := &errors.ErrCriticalOption{}
			if errors.As(err, &e) {
				log.Error(oerr)
				return nil, oerr
			}
			log.Warn(oerr)
		}
	}
	return mt, nil
}
