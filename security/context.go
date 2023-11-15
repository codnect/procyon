package security

import "context"

type ctxSecurity struct{}

var ctxSecurityKey = &ctxSecurity{}

type Context interface {
	Authentication() Authentication
	SetAuthentication(authentication Authentication)
}

func FromContext(ctx context.Context) Context {
	securityCtx, ok := ctx.Value(ctxSecurityKey).(Context)
	if !ok {
		return nil
	}

	return securityCtx
}

func ToContext(ctx context.Context, securityCtx Context) context.Context {
	return context.WithValue(ctx, ctxSecurityKey, securityCtx)
}
