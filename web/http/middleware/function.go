package middleware

import "codnect.io/procyon/web/http"

type Function func(ctx http.Context, next http.RequestDelegate) error
