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
	applyOption(h *health)
}

type health struct {
	status  Status
	details map[string]any
}

func newHealth(status Status, options ...Option) *health {
	h := &health{
		status:  status,
		details: map[string]any{},
	}

	for _, option := range options {
		option.applyOption(h)
	}

	return h
}

func (h *health) HealthStatus() Status {
	return h.status
}

func (h *health) Details() map[string]any {
	return h.details
}

func Unknown(options ...Option) Health {
	return newHealth(StatusUnknown, options...)
}

func Up(options ...Option) Health {
	return newHealth(StatusUp, options...)
}

func Down(options ...Option) Health {
	return newHealth(StatusDown, options...)
}

func OutOfService(options ...Option) Health {
	return newHealth(StatusOutOfService, options...)
}

type detailsOption map[string]any

func (d detailsOption) applyOption(h *health) {
	h.details = d
}

func WithDetails(details map[string]any) Option {
	return detailsOption(details)
}
