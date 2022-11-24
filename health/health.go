package health

type Status int

const (
	StatusUnknown Status = iota + 1
	StatusUp
	StatusDown
	StatusOutOfService
)

type Health interface {
	HealthStatus() Status
	Details() map[string]any
}

type Option interface {
	applyOption()
}

func Unknown(options ...Option) Health {
	return nil
}

func Up(options ...Option) Health {
	return nil
}

func Down(options ...Option) Health {
	return nil
}

func OutOfService(options ...Option) Health {
	return nil
}

func WithDetails(map[string]any) Option {
	return nil
}
