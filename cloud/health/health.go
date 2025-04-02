package health

type Option func(health *Health)

type Health struct {
	status  Status
	details map[string]any
}

func newHealth(status Status, options ...Option) *Health {
	h := &Health{
		status:  status,
		details: map[string]any{},
	}

	for _, opt := range options {
		opt(h)
	}

	return h
}

func (h *Health) HealthStatus() Status {
	return h.status
}

func (h *Health) Details() map[string]any {
	return h.details
}

func Of(status Status, options ...Option) *Health {
	return newHealth(status, options...)
}

func Unknown(options ...Option) *Health {
	return newHealth(StatusUnknown, options...)
}

func Up(options ...Option) *Health {
	return newHealth(StatusUp, options...)
}

func Down(options ...Option) *Health {
	return newHealth(StatusDown, options...)
}

func OutOfService(options ...Option) *Health {
	return newHealth(StatusOutOfService, options...)
}

func WithError(err error) Option {
	return nil
}

func WithDetail(key string, value any) Option {
	return nil
}
func WithDetails(details map[string]any) Option {
	return nil
}
