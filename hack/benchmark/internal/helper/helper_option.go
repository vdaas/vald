package helper

type HelperOption func(*helper)

var (
	defaultHelperOption = []HelperOption{}
)

func WithTargets(targets []string) HelperOption {
	return func(h *helper) {
		if len(targets) != 0 {
		}
	}
}

func WithParallel() HelperOption {
	return func(h *helper) {
		h.parallel = true
	}
}

func WithSequential() HelperOption {
	return func(h *helper) {

	}
}
