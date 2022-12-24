package cloud

import "net/url"

type ServiceInstance interface {
	InstanceId() string
	ServiceId() string
	URL() url.URL
	Scheme() string
	Host() string
	Port() int
	IsSecure() bool
	Metadata() map[string]string
}
