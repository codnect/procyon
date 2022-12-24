package cloud

type DiscoveryClient interface {
	Description() string
	ServiceInstances(serviceId string) []ServiceInstance
	Services() []string
}
