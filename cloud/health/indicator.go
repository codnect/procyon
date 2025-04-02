package health

import "context"

type Indicator interface {
	DoHealthCheck(ctx context.Context) Health
}

type PingIndicator struct {
}

func NewPingIndicator() *PingIndicator {
	return &PingIndicator{}
}

func (p *PingIndicator) DoHealthCheck(ctx context.Context) Health {
	return Up()
}
