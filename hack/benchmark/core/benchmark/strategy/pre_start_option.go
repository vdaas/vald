package strategy

type defaultPreStartOption func(*defaultPreStart)

func withDefaultPreStartCreateIndex(flag bool) defaultPreStartOption {
	return func(d *defaultPreStart) {
		d.createIndex = flag
	}
}
