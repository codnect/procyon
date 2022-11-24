package cloud

type ServiceRegistry interface {
	Register(instance ServiceInstance)
	Deregister(instance ServiceInstance)
	SetStatus(instance ServiceInstance, status string)
	GetStatus(instance ServiceInstance) any
}
