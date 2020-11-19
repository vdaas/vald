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

func WithAgentName(an string) DiscovererOption {
	return func(d *discoverer) error {
		d.agentName = an
		return nil
	}
}

func WithAgentNamespace(ans string) DiscovererOption {
	return func(d *discoverer) error {
		d.agentNamespace = ans
		return nil
	}
}

func WithAgentResouceType(art string) DiscovererOption {
	return func(d *discoverer) error {
		d.agentResourceType = art
		return nil
	}
}
