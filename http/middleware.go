// Copyright 2026 Codnect
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package http

// Middleware represents a component in the request pipeline that operates
// on the raw request/response level. It executes before and after the
// handler chain but has no access to the structured Result.
type Middleware interface {
	Invoke(ctx *Context, next RequestDelegate) error
}

// routingMiddleware is responsible for matching the incoming request
// to a registered endpoint. It runs early in the pipeline so that
// subsequent middleware can inspect the matched endpoint (e.g. for
// authorization or logging) before it is executed.
//
// If a match is found, the endpoint is stored in the context via
// ctx.SetEndpoint. If no match is found, the context will have
// a nil endpoint and downstream middleware can decide how to respond.
type routingMiddleware struct {
	matcher EndpointMatcher
}

// newRoutingMiddleware creates a new routingMiddleware with the given
// EndpointMatcher. This middleware should always be the first in the pipeline.
func newRoutingMiddleware(matcher EndpointMatcher) *routingMiddleware {
	return &routingMiddleware{
		matcher: matcher,
	}
}

// Invoke matches the incoming request against the registered endpoints.
// If a matching endpoint is found, it is set on the context for later use.
// The next delegate is always called regardless of whether a match was found.
func (r *routingMiddleware) Invoke(ctx *Context, next RequestDelegate) error {
	endpoint, ok := r.matcher.Match(ctx)

	if ok {
		ctx.SetEndpoint(endpoint)
	}

	return next(ctx)
}

// endpointMiddleware is responsible for executing the matched endpoint's
// request delegate. It should always be the last middleware in the pipeline.
//
// If no endpoint was matched by the routing middleware, it passes control
// to the next delegate (which is typically the terminal no-op), allowing
// the response to fall through as a 404 or be handled elsewhere.
type endpointMiddleware struct {
}

// newEndpointMiddleware creates a new endpointMiddleware.
// This middleware should always be the last in the pipeline.
func newEndpointMiddleware() *endpointMiddleware {
	return &endpointMiddleware{}
}

// Invoke retrieves the endpoint from the context and executes its
// request delegate. If no endpoint is present, it calls next to allow
// the pipeline to complete normally.
func (r *endpointMiddleware) Invoke(ctx *Context, next RequestDelegate) error {
	endpoint := ctx.Endpoint()

	if endpoint == nil {
		return next(ctx)
	}

	return endpoint.RequestDelegate()(ctx)
}
