package log

// GlgOption represetns option for glglogger.
type GlgOption func(*glglogger)

var (
	defaultGlgOpts = []GlgOption{
		WithLogLevel(INFO.String()),
	}
)

func WithLogLevel(level string) GlgOption {
	return func(g *glglogger) {
		// uppr := strings.ToUpper(level)
		// for _, level := range levels {
		// 	if uppr == level {
		// 		g.level = uppr
		// 		return
		// 	}
		// }
	}
}
