package health

type Status int

const (
	StatusUnknown Status = iota + 1
	StatusUp
	StatusDown
	StatusOutOfService
)
