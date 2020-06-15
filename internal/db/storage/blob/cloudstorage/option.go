package cloudstorage

// Option configures client of google cloud storage.
type Option func(*client) error

var (
	defaultOpts = []Option{}
)
