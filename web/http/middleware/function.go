package middleware

import "github.com/procyon-projects/procyon/web/http"

type Function func(ctx http.Context, next http.RequestDelegate) error
