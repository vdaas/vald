package ngtd

type Option func(*server)

var (
	defaultOptions = []Option{
		WithDimentaion(128),
		WithIndexDir("/tmp/ngtd/"),
		WithServerType(HTTP),
		WithPort(8200),
	}
)

func WithDimentaion(dim int) Option {
	return func(n *server) {
		if dim > 0 {
			n.dim = dim
		}
	}
}

func WithServerType(t ServerType) Option {
	return func(n *server) {
		n.srvType = t
	}
}

func WithIndexDir(path string) Option {
	return func(n *server) {
		if len(path) != 0 {
			n.indexDir = path
		}
	}
}

func WithPort(port int) Option {
	return func(n *server) {
		if port > 0 {
			n.port = port
		}
	}
}
