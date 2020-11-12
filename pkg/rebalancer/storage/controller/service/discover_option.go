package service

type DiscovererOption func(d *discoverer) error

var (
	defaultDiscovererOpts = []DiscovererOption{}
)

func WithJobName(name string) DiscovererOption {
	return func(d *discoverer) error {
		d.jobName = name
		return nil
	}
}

func WithJobNamespace(ns string) DiscovererOption {
	return func(d *discoverer) error {
		d.jobNamespace = ns
		return nil
	}
}
