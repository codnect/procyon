package health

type Checker interface {
	DoHealthCheck() (Health, error)
}
