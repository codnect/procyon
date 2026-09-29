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

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDefaultEndpointMatcher_Match(t *testing.T) {
	t.Run("static routes", func(t *testing.T) {
		registry := NewDefaultEndpointRegistry()

		endpoints := []struct {
			method Method
			path   string
		}{
			{MethodGet, "/"},
			{MethodGet, "/user"},
			{MethodGet, "/user/repos"},
			{MethodGet, "/user/orgs"},
			{MethodGet, "/user/followers"},
			{MethodGet, "/user/following"},
			{MethodGet, "/user/emails"},
			{MethodGet, "/user/keys"},
			{MethodGet, "/user/teams"},
			{MethodGet, "/user/issues"},
			{MethodPost, "/user/repos"},
			{MethodGet, "/users"},
			{MethodGet, "/gists"},
			{MethodGet, "/gists/public"},
			{MethodGet, "/gists/starred"},
			{MethodPost, "/gists"},
			{MethodGet, "/notifications"},
			{MethodPut, "/notifications"},
			{MethodGet, "/authorizations"},
			{MethodPost, "/authorizations"},
			{MethodGet, "/repositories"},
			{MethodGet, "/repos"},
			{MethodGet, "/search/repositories"},
			{MethodGet, "/search/commits"},
			{MethodGet, "/search/code"},
			{MethodGet, "/search/issues"},
			{MethodGet, "/search/users"},
			{MethodGet, "/search/topics"},
			{MethodGet, "/search/labels"},
			{MethodGet, "/gitignore/templates"},
			{MethodGet, "/licenses"},
			{MethodGet, "/emojis"},
			{MethodGet, "/markdown"},
			{MethodPost, "/markdown"},
			{MethodGet, "/meta"},
			{MethodGet, "/rate_limit"},
			{MethodGet, "/feeds"},
			{MethodGet, "/events"},
		}

		for _, endpoint := range endpoints {
			err := registry.Register(&Endpoint{
				method: endpoint.method,
				path:   endpoint.path,
			})
			if err != nil {
				t.Fatalf("failed to register endpoint %s %s: %v",
					endpoint.method, endpoint.path, err)
			}
		}

		matcher := NewDefaultEndpointMatcher(registry)

		testCases := []struct {
			method Method
			path   string
			match  bool
		}{
			{MethodGet, "/", true},
			{MethodGet, "/user", true},
			{MethodGet, "/user/repos", true},
			{MethodGet, "/user/orgs", true},
			{MethodGet, "/user/followers", true},
			{MethodGet, "/user/following", true},
			{MethodGet, "/user/emails", true},
			{MethodGet, "/user/keys", true},
			{MethodGet, "/user/teams", true},
			{MethodGet, "/user/issues", true},
			{MethodPost, "/user/repos", true},
			{MethodGet, "/users", true},
			{MethodGet, "/gists", true},
			{MethodGet, "/gists/public", true},
			{MethodGet, "/gists/starred", true},
			{MethodPost, "/gists", true},
			{MethodGet, "/notifications", true},
			{MethodPut, "/notifications", true},
			{MethodGet, "/authorizations", true},
			{MethodPost, "/authorizations", true},
			{MethodGet, "/repositories", true},
			{MethodGet, "/repos", true},
			{MethodGet, "/search/repositories", true},
			{MethodGet, "/search/commits", true},
			{MethodGet, "/search/code", true},
			{MethodGet, "/search/issues", true},
			{MethodGet, "/search/users", true},
			{MethodGet, "/search/topics", true},
			{MethodGet, "/search/labels", true},
			{MethodGet, "/gitignore/templates", true},
			{MethodGet, "/licenses", true},
			{MethodGet, "/emojis", true},
			{MethodGet, "/markdown", true},
			{MethodPost, "/markdown", true},
			{MethodGet, "/meta", true},
			{MethodGet, "/rate_limit", true},
			{MethodGet, "/feeds", true},
			{MethodGet, "/events", true},

			{MethodDelete, "/user/repos", false},

			{MethodGet, "/user/follow", false},
			{MethodGet, "/use", false},
			{MethodGet, "/repo", false},
			{MethodGet, "/gists/star", false},
			{MethodGet, "/gitignore", false},
			{MethodGet, "/git", false},
			{MethodGet, "/search", false},
			{MethodGet, "/search/repo", false},
			{MethodGet, "/license", false},
			{MethodGet, "/event", false},
			{MethodGet, "/emoji", false},
			{MethodGet, "/rate", false},
			{MethodGet, "/mark", false},
			{MethodGet, "/met", false},

			{MethodGet, "/notfound", false},
			{MethodGet, "", false},
		}

		for _, tc := range testCases {
			req, _ := http.NewRequest(string(tc.method), tc.path, nil)
			recorder := httptest.NewRecorder()

			ctx := CreateContext(req, recorder)
			endpoint, ok := matcher.Match(ctx)

			if ok != tc.match {
				t.Errorf("%s %s: got match=%v, want %v",
					tc.method, tc.path, ok, tc.match)
			}

			if ok && endpoint.method != tc.method {
				t.Errorf("%s %s: got method=%s, want %s",
					tc.method, tc.path, endpoint.method, tc.method)
			}
		}
	})

	t.Run("param routes", func(t *testing.T) {
		registry := NewDefaultEndpointRegistry()

		endpoints := []struct {
			method Method
			path   string
		}{
			{MethodGet, "/users/{username}"},
			{MethodGet, "/users/{username}/repos"},
			{MethodGet, "/users/{username}/followers"},
			{MethodGet, "/repos/{owner}/{repo}"},
			{MethodGet, "/repos/{owner}/{repo}/issues"},
			{MethodGet, "/repos/{owner}/{repo}/issues/{number}"},
			{MethodGet, "/repos/{owner}/{repo}/pulls/{number}/comments"},
			{MethodGet, "/orgs/{org}/members/{username}"},
			{MethodGet, "/teams/{id}/repos/{owner}/{repo}"},
		}

		for _, endpoint := range endpoints {
			err := registry.Register(&Endpoint{
				method: endpoint.method,
				path:   endpoint.path,
			})
			if err != nil {
				t.Fatalf("failed to register endpoint %s %s: %v",
					endpoint.method, endpoint.path, err)
			}
		}

		matcher := NewDefaultEndpointMatcher(registry)

		testCases := []struct {
			method Method
			path   string
			match  bool
			params map[string]string
		}{
			{MethodGet, "/users/octocat", true, map[string]string{"username": "octocat"}},
			{MethodGet, "/users/octocat/repos", true, map[string]string{"username": "octocat"}},
			{MethodGet, "/users/octocat/followers", true, map[string]string{"username": "octocat"}},
			{MethodGet, "/repos/octocat/hello-world", true, map[string]string{"owner": "octocat", "repo": "hello-world"}},
			{MethodGet, "/repos/octocat/hello-world/issues", true, map[string]string{"owner": "octocat", "repo": "hello-world"}},
			{MethodGet, "/repos/octocat/hello-world/issues/42", true, map[string]string{"owner": "octocat", "repo": "hello-world", "number": "42"}},
			{MethodGet, "/repos/octocat/hello-world/pulls/7/comments", true, map[string]string{"owner": "octocat", "repo": "hello-world", "number": "7"}},
			{MethodGet, "/orgs/github/members/octocat", true, map[string]string{"org": "github", "username": "octocat"}},
			{MethodGet, "/teams/1/repos/octocat/hello-world", true, map[string]string{"id": "1", "owner": "octocat", "repo": "hello-world"}},
			{MethodGet, "/users", false, nil},
			{MethodGet, "/repos/octocat", false, nil},
			{MethodGet, "/repos/octocat/hello-world/issues/42/labels/bug/extra", false, nil},
		}

		for _, tc := range testCases {
			req, _ := http.NewRequest(string(tc.method), tc.path, nil)
			recorder := httptest.NewRecorder()

			ctx := CreateContext(req, recorder)
			_, ok := matcher.Match(ctx)

			if ok != tc.match {
				t.Errorf("%s %s: got match=%v, want %v",
					tc.method, tc.path, ok, tc.match)
			}

			if ok && tc.params != nil {
				for name, want := range tc.params {
					if got := ctx.Request().PathValue(name); got != want {
						t.Errorf("%s %s: param %s: got %q, want %q",
							tc.method, tc.path, name, got, want)
					}
				}
			}
		}
	})

	t.Run("wildcard routes", func(t *testing.T) {
		registry := NewDefaultEndpointRegistry()

		endpoints := []struct {
			method Method
			path   string
		}{
			{MethodGet, "/repos/{owner}/{repo}/releases/*/assets"},
			{MethodGet, "/orgs/{org}/hooks/*/pings"},
			{MethodGet, "/repos/{owner}/{repo}/deployments/*/statuses"},
		}

		for _, endpoint := range endpoints {
			err := registry.Register(&Endpoint{
				method: endpoint.method,
				path:   endpoint.path,
			})
			if err != nil {
				t.Fatalf("failed to register endpoint %s %s: %v",
					endpoint.method, endpoint.path, err)
			}
		}

		matcher := NewDefaultEndpointMatcher(registry)

		testCases := []struct {
			method Method
			path   string
			match  bool
		}{
			{MethodGet, "/repos/octocat/hello-world/releases/v1.0/assets", true},
			{MethodGet, "/repos/octocat/hello-world/releases/latest/assets", true},
			{MethodGet, "/repos/octocat/hello-world/releases/assets", false},
			{MethodGet, "/repos/octocat/hello-world/releases/v1/v2/assets", false},
			{MethodGet, "/orgs/github/hooks/42/pings", true},
			{MethodGet, "/orgs/github/hooks/pings", false},
			{MethodGet, "/repos/octocat/hello-world/deployments/123/statuses", true},
			{MethodGet, "/repos/octocat/hello-world/deployments/statuses", false},
			{MethodGet, "/repos/octocat/hello-world/deployments/a/b/statuses", false},
		}

		for _, tc := range testCases {
			req, _ := http.NewRequest(string(tc.method), tc.path, nil)
			recorder := httptest.NewRecorder()

			ctx := CreateContext(req, recorder)
			_, ok := matcher.Match(ctx)

			if ok != tc.match {
				t.Errorf("%s %s: got match=%v, want %v",
					tc.method, tc.path, ok, tc.match)
			}
		}
	})

	t.Run("double wildcard routes", func(t *testing.T) {
		registry := NewDefaultEndpointRegistry()

		endpoints := []struct {
			method Method
			path   string
		}{
			{MethodGet, "/repos/{owner}/{repo}/git/**"},
			{MethodGet, "/static/**"},
			{MethodGet, "/api/**/health"},
			{MethodGet, "/cdn/{vendor}/**"},
			{MethodGet, "/proxy/**/status"},
			{MethodGet, "/staticx"},
		}

		for _, endpoint := range endpoints {
			err := registry.Register(&Endpoint{
				method: endpoint.method,
				path:   endpoint.path,
			})
			if err != nil {
				t.Fatalf("failed to register endpoint %s %s: %v",
					endpoint.method, endpoint.path, err)
			}
		}

		matcher := NewDefaultEndpointMatcher(registry)

		testCases := []struct {
			method Method
			path   string
			match  bool
		}{
			{MethodGet, "/repos/octocat/hello-world/git", true},
			{MethodGet, "/repos/octocat/hello-world/git/refs", true},
			{MethodGet, "/repos/octocat/hello-world/git/refs/heads/main", true},
			{MethodGet, "/repos/octocat/hello-world/git/commits/abc123", true},
			{MethodGet, "/repos/octocat/hello-world/git/trees/def456", true},
			{MethodGet, "/repos/octocat/hello-world/git/blobs/a1b2c3/raw", true},

			{MethodGet, "/static", true},
			{MethodGet, "/static/css/main.css", true},
			{MethodGet, "/static/js/vendor/lodash.min.js", true},
			{MethodGet, "/staticx", true},

			{MethodGet, "/api/health", true},
			{MethodGet, "/api/v1/health", true},
			{MethodGet, "/api/v1/v2/v3/health", true},
			{MethodGet, "/api/health/health", true},
			{MethodGet, "/api/v1/check", false},

			{MethodGet, "/cdn/cloudflare", true},
			{MethodGet, "/cdn/cloudflare/js/app.js", true},
			{MethodGet, "/cdn/aws/images/logo.png", true},

			{MethodGet, "/proxy/status", true},
			{MethodGet, "/proxy/us-east/status", true},
			{MethodGet, "/proxy/us-east/zone-1/status", true},
			{MethodGet, "/proxy/us-east/zone-1/zone-2/status", true},
			{MethodGet, "/proxy/us-east/health", false},
		}

		for _, tc := range testCases {
			req, _ := http.NewRequest(string(tc.method), tc.path, nil)
			recorder := httptest.NewRecorder()

			ctx := CreateContext(req, recorder)
			_, ok := matcher.Match(ctx)

			if ok != tc.match {
				t.Errorf("%s %s: got match=%v, want %v",
					tc.method, tc.path, ok, tc.match)
			}
		}
	})

	t.Run("pattern routes", func(t *testing.T) {
		registry := NewDefaultEndpointRegistry()

		endpoints := []struct {
			method Method
			path   string
		}{
			{MethodGet, "/repos/{owner}/{repo}/contents/*.json"},
			{MethodGet, "/repos/{owner}/{repo}/archive/*.tar.gz"},
			{MethodGet, "/repos/{owner}/{repo}/contents/*.json/meta"},
			{MethodGet, "/downloads/release-v?.zip"},
			{MethodGet, "/exports/report-??-*.csv"},
			{MethodGet, "/assets/logo*"},
		}

		for _, endpoint := range endpoints {
			err := registry.Register(&Endpoint{
				method: endpoint.method,
				path:   endpoint.path,
			})
			if err != nil {
				t.Fatalf("failed to register endpoint %s %s: %v",
					endpoint.method, endpoint.path, err)
			}
		}

		matcher := NewDefaultEndpointMatcher(registry)

		testCases := []struct {
			method Method
			path   string
			match  bool
		}{
			{MethodGet, "/repos/octocat/hello-world/contents/config.json", true},
			{MethodGet, "/repos/octocat/hello-world/contents/package.json", true},
			{MethodGet, "/repos/octocat/hello-world/contents/.json", true},
			{MethodGet, "/repos/octocat/hello-world/contents/readme.md", false},
			{MethodGet, "/repos/octocat/hello-world/contents/data.jsonx", false},

			{MethodGet, "/repos/octocat/hello-world/archive/v1.0.tar.gz", true},
			{MethodGet, "/repos/octocat/hello-world/archive/latest.tar.gz", true},
			{MethodGet, "/repos/octocat/hello-world/archive/v1.0.zip", false},

			{MethodGet, "/repos/octocat/hello-world/contents/schema.json/meta", true},
			{MethodGet, "/repos/octocat/hello-world/contents/data.json/meta", true},

			{MethodGet, "/downloads/release-v1.zip", true},
			{MethodGet, "/downloads/release-v9.zip", true},
			{MethodGet, "/downloads/release-v12.zip", false},
			{MethodGet, "/downloads/release-v.zip", false},

			{MethodGet, "/exports/report-US-2024.csv", true},
			{MethodGet, "/exports/report-EU-sales.csv", true},
			{MethodGet, "/exports/report-A-2024.csv", false},
			{MethodGet, "/exports/report-USA-2024.csv", false},

			{MethodGet, "/assets/logo", true},
			{MethodGet, "/assets/logo-dark", true},
			{MethodGet, "/assets/logo192.png", true},
			{MethodGet, "/assets/log", false},
		}

		for _, tc := range testCases {
			req, _ := http.NewRequest(string(tc.method), tc.path, nil)
			recorder := httptest.NewRecorder()

			ctx := CreateContext(req, recorder)
			_, ok := matcher.Match(ctx)

			if ok != tc.match {
				t.Errorf("%s %s: got match=%v, want %v",
					tc.method, tc.path, ok, tc.match)
			}
		}
	})

	t.Run("priority: static > pattern > param > wildcard > double wildcard", func(t *testing.T) {
		registry := NewDefaultEndpointRegistry()

		endpoints := []struct {
			method Method
			path   string
		}{
			{MethodGet, "/repos/octocat/hello-world"},
			{MethodGet, "/repos/{owner}/{repo}/contents/*.json"},
			{MethodGet, "/repos/{owner}/{repo}"},
			{MethodGet, "/repos/{owner}/*"},
			{MethodGet, "/repos/**"},
			{MethodGet, "/gists/public"},
			{MethodGet, "/gists/starred"},
			{MethodGet, "/gists/{id}"},
			{MethodGet, "/repos/{owner}/{repo}/releases/latest"},
			{MethodGet, "/repos/{owner}/{repo}/releases/{id}"},
			{MethodGet, "/repos/{owner}/{repo}/issues/comments"},
			{MethodGet, "/repos/{owner}/{repo}/issues/{number}"},
		}

		for _, endpoint := range endpoints {
			err := registry.Register(&Endpoint{
				method: endpoint.method,
				path:   endpoint.path,
			})
			if err != nil {
				t.Fatalf("failed to register endpoint %s %s: %v",
					endpoint.method, endpoint.path, err)
			}
		}

		matcher := NewDefaultEndpointMatcher(registry)

		testCases := []struct {
			path        string
			wantPattern string
		}{
			{"/repos/octocat/hello-world", "/repos/octocat/hello-world"},
			{"/gists/public", "/gists/public"},
			{"/gists/starred", "/gists/starred"},
			{"/repos/octocat/hello-world/releases/latest", "/repos/{owner}/{repo}/releases/latest"},
			{"/repos/octocat/hello-world/issues/comments", "/repos/{owner}/{repo}/issues/comments"},

			{"/repos/octocat/hello-world/contents/config.json", "/repos/{owner}/{repo}/contents/*.json"},

			{"/repos/facebook/react", "/repos/{owner}/{repo}"},
			{"/gists/abc123", "/gists/{id}"},
			{"/repos/octocat/hello-world/releases/42", "/repos/{owner}/{repo}/releases/{id}"},
			{"/repos/octocat/hello-world/issues/99", "/repos/{owner}/{repo}/issues/{number}"},

			{"/repos/octocat/hello-world/deep/nested/path", "/repos/**"},
		}

		for _, tc := range testCases {
			req, _ := http.NewRequest("GET", tc.path, nil)
			recorder := httptest.NewRecorder()

			ctx := CreateContext(req, recorder)
			endpoint, ok := matcher.Match(ctx)

			if !ok {
				t.Errorf("GET %s: no match, want %s", tc.path, tc.wantPattern)
				continue
			}

			if endpoint.path != tc.wantPattern {
				t.Errorf("GET %s: matched %s, want %s",
					tc.path, endpoint.path, tc.wantPattern)
			}
		}
	})

	t.Run("backtracking", func(t *testing.T) {
		registry := NewDefaultEndpointRegistry()

		endpoints := []struct {
			method Method
			path   string
		}{
			{MethodGet, "/repos/{owner}/{repo}/issues/{number}"},
			{MethodGet, "/repos/{owner}/{repo}/issues/*/reactions"},

			{MethodGet, "/repos/{owner}/{repo}/contents/*.json/download"},
			{MethodGet, "/repos/{owner}/{repo}/contents/*.md/preview"},

			{MethodGet, "/repos/{owner}/{repo}/branches/{branch}"},
			{MethodGet, "/repos/{owner}/{repo}/git/**"},
		}

		for _, endpoint := range endpoints {
			err := registry.Register(&Endpoint{
				method: endpoint.method,
				path:   endpoint.path,
			})
			if err != nil {
				t.Fatalf("failed to register endpoint %s %s: %v",
					endpoint.method, endpoint.path, err)
			}
		}

		matcher := NewDefaultEndpointMatcher(registry)

		testCases := []struct {
			path        string
			match       bool
			wantPattern string
		}{
			{"/repos/octocat/hello-world/issues/42", true, "/repos/{owner}/{repo}/issues/{number}"},
			{"/repos/octocat/hello-world/branches/main", true, "/repos/{owner}/{repo}/branches/{branch}"},
			{"/repos/octocat/hello-world/issues/42/reactions", true, "/repos/{owner}/{repo}/issues/*/reactions"},
			{"/repos/octocat/hello-world/contents/readme.md/preview", true, "/repos/{owner}/{repo}/contents/*.md/preview"},
			{"/repos/octocat/hello-world/contents/config.json/download", true, "/repos/{owner}/{repo}/contents/*.json/download"},
			{"/repos/octocat/hello-world/contents/readme.md/download", false, ""},
			{"/repos/octocat/hello-world/contents/config.json/preview", false, ""},
			{"/repos/octocat/hello-world/git/refs/heads/main", true, "/repos/{owner}/{repo}/git/**"},
			{"/repos/octocat/hello-world/git/commits/abc123", true, "/repos/{owner}/{repo}/git/**"},
		}

		for _, tc := range testCases {
			req, _ := http.NewRequest("GET", tc.path, nil)
			recorder := httptest.NewRecorder()

			ctx := CreateContext(req, recorder)
			endpoint, ok := matcher.Match(ctx)

			if ok != tc.match {
				t.Errorf("GET %s: got match=%v, want %v",
					tc.path, ok, tc.match)
				continue
			}

			if ok && endpoint.path != tc.wantPattern {
				t.Errorf("GET %s: matched %s, want %s",
					tc.path, endpoint.path, tc.wantPattern)
			}
		}
	})

	t.Run("multiple methods same path", func(t *testing.T) {
		registry := NewDefaultEndpointRegistry()

		endpoints := []struct {
			method Method
			path   string
		}{
			{MethodGet, "/repos/{owner}/{repo}"},
			{MethodPatch, "/repos/{owner}/{repo}"},
			{MethodDelete, "/repos/{owner}/{repo}"},
			{MethodGet, "/repos/{owner}/{repo}/issues"},
			{MethodPost, "/repos/{owner}/{repo}/issues"},
			{MethodGet, "/gists/{id}"},
			{MethodPatch, "/gists/{id}"},
			{MethodDelete, "/gists/{id}"},
			{MethodPut, "/repos/{owner}/{repo}/subscription"},
			{MethodHead, "/repos/{owner}/{repo}/subscription"},
			{MethodOptions, "/repos/{owner}/{repo}/subscription"},
			{MethodTrace, "/repos/{owner}/{repo}/subscription"},
		}

		for _, endpoint := range endpoints {
			err := registry.Register(&Endpoint{
				method: endpoint.method,
				path:   endpoint.path,
			})
			if err != nil {
				t.Fatalf("failed to register endpoint %s %s: %v",
					endpoint.method, endpoint.path, err)
			}
		}

		matcher := NewDefaultEndpointMatcher(registry)

		testCases := []struct {
			method Method
			path   string
			match  bool
		}{
			{MethodGet, "/repos/octocat/hello-world", true},
			{MethodPatch, "/repos/octocat/hello-world", true},
			{MethodDelete, "/repos/octocat/hello-world", true},
			{MethodPut, "/repos/octocat/hello-world", false},
			{MethodPost, "/repos/octocat/hello-world", false},

			{MethodGet, "/repos/octocat/hello-world/issues", true},
			{MethodPost, "/repos/octocat/hello-world/issues", true},
			{MethodDelete, "/repos/octocat/hello-world/issues", false},

			{MethodGet, "/gists/abc123", true},
			{MethodPatch, "/gists/abc123", true},
			{MethodDelete, "/gists/abc123", true},
			{MethodPost, "/gists/abc123", false},

			{MethodPut, "/repos/octocat/hello-world/subscription", true},
			{MethodHead, "/repos/octocat/hello-world/subscription", true},
			{MethodOptions, "/repos/octocat/hello-world/subscription", true},
			{MethodTrace, "/repos/octocat/hello-world/subscription", true},
			{MethodGet, "/repos/octocat/hello-world/subscription", false},
		}

		for _, tc := range testCases {
			req, _ := http.NewRequest(string(tc.method), tc.path, nil)
			recorder := httptest.NewRecorder()

			ctx := CreateContext(req, recorder)
			endpoint, ok := matcher.Match(ctx)

			if ok != tc.match {
				t.Errorf("%s %s: got match=%v, want %v",
					tc.method, tc.path, ok, tc.match)
				continue
			}

			if ok && endpoint.method != tc.method {
				t.Errorf("%s %s: got method=%s, want %s",
					tc.method, tc.path, endpoint.method, tc.method)
			}
		}
	})

	t.Run("mixed route types in same tree", func(t *testing.T) {
		registry := NewDefaultEndpointRegistry()

		endpoints := []struct {
			method Method
			path   string
		}{
			{MethodGet, "/"},
			{MethodGet, "/user"},
			{MethodGet, "/user/repos"},
			{MethodGet, "/user/followers"},
			{MethodGet, "/user/following"},
			{MethodGet, "/users"},
			{MethodGet, "/gists"},
			{MethodGet, "/gists/public"},
			{MethodGet, "/gists/starred"},
			{MethodGet, "/notifications"},
			{MethodGet, "/repositories"},
			{MethodGet, "/search/repositories"},
			{MethodGet, "/search/code"},
			{MethodGet, "/gitignore/templates"},
			{MethodGet, "/licenses"},
			{MethodGet, "/events"},

			{MethodGet, "/users/{username}"},
			{MethodGet, "/users/{username}/repos"},
			{MethodGet, "/users/{username}/followers"},
			{MethodGet, "/repos/{owner}/{repo}"},
			{MethodPatch, "/repos/{owner}/{repo}"},
			{MethodDelete, "/repos/{owner}/{repo}"},
			{MethodGet, "/repos/{owner}/{repo}/issues"},
			{MethodPost, "/repos/{owner}/{repo}/issues"},
			{MethodGet, "/repos/{owner}/{repo}/issues/{number}"},
			{MethodGet, "/repos/{owner}/{repo}/issues/comments"},
			{MethodGet, "/repos/{owner}/{repo}/issues/comments/{id}"},
			{MethodGet, "/repos/{owner}/{repo}/pulls"},
			{MethodGet, "/repos/{owner}/{repo}/pulls/{number}"},
			{MethodGet, "/repos/{owner}/{repo}/commits"},
			{MethodGet, "/repos/{owner}/{repo}/commits/{sha}"},
			{MethodGet, "/repos/{owner}/{repo}/releases"},
			{MethodGet, "/repos/{owner}/{repo}/releases/{id}"},
			{MethodGet, "/repos/{owner}/{repo}/releases/latest"},
			{MethodGet, "/repos/{owner}/{repo}/releases/tags/{tag}"},
			{MethodGet, "/repos/{owner}/{repo}/branches"},
			{MethodGet, "/repos/{owner}/{repo}/branches/{branch}"},
			{MethodGet, "/repos/{owner}/{repo}/contents/{path}"},
			{MethodPut, "/repos/{owner}/{repo}/subscription"},
			{MethodGet, "/gists/{id}"},
			{MethodPatch, "/gists/{id}"},
			{MethodDelete, "/gists/{id}"},
			{MethodGet, "/gists/{id}/commits"},
			{MethodGet, "/gists/{id}/comments"},
			{MethodGet, "/gists/{id}/comments/{commentId}"},
			{MethodGet, "/orgs/{org}"},
			{MethodGet, "/orgs/{org}/repos"},
			{MethodGet, "/orgs/{org}/members"},
			{MethodGet, "/orgs/{org}/members/{username}"},
			{MethodGet, "/teams/{id}"},
			{MethodGet, "/teams/{id}/members"},
			{MethodGet, "/teams/{id}/repos/{owner}/{repo}"},
			{MethodGet, "/authorizations/{id}"},
			{MethodGet, "/notifications/threads/{id}"},
			{MethodGet, "/notifications/threads/{id}/subscription"},
			{MethodGet, "/gitignore/templates/{name}"},
			{MethodGet, "/licenses/{license}"},

			{MethodGet, "/repos/{owner}/{repo}/releases/*/assets"},

			{MethodGet, "/repos/{owner}/{repo}/git/**"},
			{MethodGet, "/static/**"},

			{MethodGet, "/repos/{owner}/{repo}/contents/*.json"},
			{MethodGet, "/repos/{owner}/{repo}/archive/*.tar.gz"},
		}

		for _, endpoint := range endpoints {
			err := registry.Register(&Endpoint{
				method: endpoint.method,
				path:   endpoint.path,
			})
			if err != nil {
				t.Fatalf("failed to register endpoint %s %s: %v",
					endpoint.method, endpoint.path, err)
			}
		}

		matcher := NewDefaultEndpointMatcher(registry)

		testCases := []struct {
			method      Method
			path        string
			match       bool
			wantPattern string
		}{
			{MethodGet, "/user", true, "/user"},
			{MethodGet, "/gists/public", true, "/gists/public"},
			{MethodGet, "/search/repositories", true, "/search/repositories"},
			{MethodGet, "/gitignore/templates", true, "/gitignore/templates"},
			{MethodGet, "/events", true, "/events"},

			{MethodGet, "/gists/public", true, "/gists/public"},
			{MethodGet, "/gists/abc123", true, "/gists/{id}"},
			{MethodGet, "/gists/starred", true, "/gists/starred"},
			{MethodGet, "/repos/octocat/hello-world/releases/latest", true, "/repos/{owner}/{repo}/releases/latest"},
			{MethodGet, "/repos/octocat/hello-world/releases/42", true, "/repos/{owner}/{repo}/releases/{id}"},
			{MethodGet, "/repos/octocat/hello-world/issues/comments", true, "/repos/{owner}/{repo}/issues/comments"},
			{MethodGet, "/repos/octocat/hello-world/issues/42", true, "/repos/{owner}/{repo}/issues/{number}"},

			{MethodGet, "/users/octocat", true, "/users/{username}"},
			{MethodGet, "/users/octocat/repos", true, "/users/{username}/repos"},
			{MethodGet, "/repos/octocat/hello-world", true, "/repos/{owner}/{repo}"},
			{MethodGet, "/repos/octocat/hello-world/issues/42", true, "/repos/{owner}/{repo}/issues/{number}"},
			{MethodGet, "/repos/octocat/hello-world/pulls/7", true, "/repos/{owner}/{repo}/pulls/{number}"},
			{MethodGet, "/repos/octocat/hello-world/commits/abc123", true, "/repos/{owner}/{repo}/commits/{sha}"},
			{MethodGet, "/repos/octocat/hello-world/branches/main", true, "/repos/{owner}/{repo}/branches/{branch}"},
			{MethodGet, "/repos/octocat/hello-world/releases/tags/v1.0", true, "/repos/{owner}/{repo}/releases/tags/{tag}"},
			{MethodGet, "/repos/octocat/hello-world/contents/README.md", true, "/repos/{owner}/{repo}/contents/{path}"},
			{MethodGet, "/orgs/github/members/octocat", true, "/orgs/{org}/members/{username}"},
			{MethodGet, "/teams/1/repos/octocat/hello-world", true, "/teams/{id}/repos/{owner}/{repo}"},
			{MethodGet, "/gists/abc123/comments/5", true, "/gists/{id}/comments/{commentId}"},
			{MethodGet, "/notifications/threads/77/subscription", true, "/notifications/threads/{id}/subscription"},

			{MethodGet, "/repos/octocat/hello-world", true, "/repos/{owner}/{repo}"},
			{MethodPatch, "/repos/octocat/hello-world", true, "/repos/{owner}/{repo}"},
			{MethodDelete, "/repos/octocat/hello-world", true, "/repos/{owner}/{repo}"},
			{MethodPost, "/repos/octocat/hello-world/issues", true, "/repos/{owner}/{repo}/issues"},
			{MethodPut, "/repos/octocat/hello-world/subscription", true, "/repos/{owner}/{repo}/subscription"},

			{MethodGet, "/repos/octocat/hello-world/releases/v1.0/assets", true, "/repos/{owner}/{repo}/releases/*/assets"},
			{MethodGet, "/repos/octocat/hello-world/releases/v1/v2/assets", false, ""},

			{MethodGet, "/repos/octocat/hello-world/releases/assets", true, "/repos/{owner}/{repo}/releases/{id}"},

			{MethodGet, "/repos/octocat/hello-world/git/refs", true, "/repos/{owner}/{repo}/git/**"},
			{MethodGet, "/repos/octocat/hello-world/git/refs/heads/main", true, "/repos/{owner}/{repo}/git/**"},
			{MethodGet, "/static/css/main.css", true, "/static/**"},
			{MethodGet, "/static/js/vendor/lodash.js", true, "/static/**"},

			{MethodGet, "/repos/octocat/hello-world/contents/config.json", true, "/repos/{owner}/{repo}/contents/*.json"},
			{MethodGet, "/repos/octocat/hello-world/contents/README.md", true, "/repos/{owner}/{repo}/contents/{path}"},
			{MethodGet, "/repos/octocat/hello-world/archive/v1.0.tar.gz", true, "/repos/{owner}/{repo}/archive/*.tar.gz"},
			{MethodGet, "/repos/octocat/hello-world/archive/v1.0.zip", false, ""},

			{MethodGet, "/repos/octocat", false, ""},
			{MethodGet, "/unknown/path", false, ""},
			{MethodDelete, "/users/octocat", false, ""},
			{MethodPost, "/gists/abc123", false, ""},
		}

		for _, tc := range testCases {
			req, _ := http.NewRequest(string(tc.method), tc.path, nil)
			recorder := httptest.NewRecorder()

			ctx := CreateContext(req, recorder)
			endpoint, ok := matcher.Match(ctx)

			if ok != tc.match {
				t.Errorf("%s %s: got match=%v, want %v",
					tc.method, tc.path, ok, tc.match)
				continue
			}

			if ok && endpoint.path != tc.wantPattern {
				t.Errorf("%s %s: matched %s, want %s",
					tc.method, tc.path, endpoint.path, tc.wantPattern)
			}
		}
	})

	t.Run("normalization", func(t *testing.T) {
		registry := NewDefaultEndpointRegistry()

		endpoints := []struct {
			method Method
			path   string
		}{
			{MethodGet, "/repos/{owner}/{repo}"},
			{MethodGet, "/users/{username}/repos"},
			{MethodGet, "/repos/{owner}/{repo}/issues/{number}"},

			{MethodGet, "gists"},
			{MethodGet, "orgs/{org}/members"},

			{MethodGet, "/notifications/"},
			{MethodGet, "/search/repositories/"},
		}

		for _, endpoint := range endpoints {
			err := registry.Register(&Endpoint{
				method: endpoint.method,
				path:   endpoint.path,
			})
			if err != nil {
				t.Fatalf("failed to register endpoint %s %s: %v",
					endpoint.method, endpoint.path, err)
			}
		}

		matcher := NewDefaultEndpointMatcher(registry)

		testCases := []struct {
			path  string
			match bool
		}{
			{"/repos/octocat/hello-world/", true},
			{"/users/octocat/repos/", true},
			{"/repos/octocat/hello-world/issues/42/", true},

			{"/gists", true},
			{"/orgs/github/members", true},

			{"/notifications", true},
			{"/search/repositories", true},

			{"/", false},
		}

		for _, tc := range testCases {
			req, _ := http.NewRequest("GET", tc.path, nil)
			recorder := httptest.NewRecorder()

			ctx := CreateContext(req, recorder)
			_, ok := matcher.Match(ctx)

			if ok != tc.match {
				t.Errorf("GET %s: got match=%v, want %v",
					tc.path, ok, tc.match)
			}
		}
	})

	t.Run("error handling", func(t *testing.T) {
		t.Run("invalid patterns", func(t *testing.T) {
			endpoints := []struct {
				name string
				path string
			}{
				{"unclosed param at end", "/repos/{owner"},
				{"unclosed param in middle", "/repos/{owner/issues"},
				{"empty param name", "/repos/{}"},
				{"empty param name in middle", "/repos/{}/issues"},
				{"too many params", "/repos/{p1}/{p2}/{p3}/{p4}/{p5}/{p6}/{p7}/{p8}/{p9}/{p10}/{p11}/{p12}/{p13}/{p14}/{p15}/{p16}/{p17}"},
			}

			for _, endpoint := range endpoints {
				t.Run(endpoint.name, func(t *testing.T) {
					matcher := NewDefaultEndpointMatcher(NewDefaultEndpointRegistry())

					err := matcher.addEndpoint(&Endpoint{
						method: MethodGet,
						path:   endpoint.path,
					})
					if err == nil {
						t.Errorf("expected error for path %q", endpoint.path)
					}
				})
			}
		})

		t.Run("valid patterns", func(t *testing.T) {
			endpoints := []struct {
				name string
				path string
			}{
				{"single param", "/repos/{owner}"},
				{"max params", "/repos/{p1}/{p2}/{p3}/{p4}/{p5}/{p6}/{p7}/{p8}/{p9}/{p10}/{p11}/{p12}/{p13}/{p14}/{p15}/{p16}"},
				{"param with static suffix", "/repos/{owner}/{repo}/issues"},
				{"double wildcard", "/repos/{owner}/{repo}/git/**"},
				{"wildcard segment", "/repos/{owner}/{repo}/releases/*/assets"},
				{"pattern segment", "/repos/{owner}/{repo}/contents/*.json"},
			}

			for _, endpoint := range endpoints {
				t.Run(endpoint.name, func(t *testing.T) {
					matcher := NewDefaultEndpointMatcher(NewDefaultEndpointRegistry())

					err := matcher.addEndpoint(&Endpoint{
						method: MethodGet,
						path:   endpoint.path,
					})
					if err != nil {
						t.Errorf("unexpected error for path %q: %v",
							endpoint.path, err)
					}
				})
			}
		})

		t.Run("duplicate routes", func(t *testing.T) {
			endpoints := []struct {
				name   string
				method Method
				path   string
			}{
				{"same static route", MethodGet, "/gists/public"},
				{"same param route", MethodGet, "/repos/{owner}/{repo}"},
				{"same deep param route", MethodPost, "/repos/{owner}/{repo}/issues"},
			}

			for _, endpoint := range endpoints {
				t.Run(endpoint.name, func(t *testing.T) {
					matcher := NewDefaultEndpointMatcher(NewDefaultEndpointRegistry())

					_ = matcher.addEndpoint(&Endpoint{
						method: endpoint.method,
						path:   endpoint.path,
					})

					err := matcher.addEndpoint(&Endpoint{
						method: endpoint.method,
						path:   endpoint.path,
					})
					if err == nil {
						t.Errorf("expected error for duplicate %s %s",
							endpoint.method, endpoint.path)
					}
				})
			}
		})

		t.Run("different method same path is ok", func(t *testing.T) {
			matcher := NewDefaultEndpointMatcher(NewDefaultEndpointRegistry())

			endpoints := []struct {
				method Method
				path   string
			}{
				{MethodGet, "/repos/{owner}/{repo}"},
				{MethodDelete, "/repos/{owner}/{repo}"},
				{MethodPatch, "/repos/{owner}/{repo}"},
				{MethodGet, "/repos/{owner}/{repo}/issues"},
				{MethodPost, "/repos/{owner}/{repo}/issues"},
				{MethodGet, "/gists"},
				{MethodPost, "/gists"},
				{MethodGet, "/gists/{id}"},
				{MethodPatch, "/gists/{id}"},
				{MethodDelete, "/gists/{id}"},
			}

			for _, endpoint := range endpoints {
				err := matcher.addEndpoint(&Endpoint{
					method: endpoint.method,
					path:   endpoint.path,
				})
				if err != nil {
					t.Errorf("unexpected error for %s %s: %v",
						endpoint.method, endpoint.path, err)
				}
			}
		})
	})
}

func BenchmarkDefaultEndpointMatcher_Match(b *testing.B) {
	registry := NewDefaultEndpointRegistry()

	endpoints := []struct {
		method Method
		path   string
	}{
		{MethodGet, "/"},
		{MethodGet, "/user"},
		{MethodGet, "/user/repos"},
		{MethodGet, "/user/followers"},
		{MethodGet, "/user/following"},
		{MethodGet, "/users"},
		{MethodGet, "/gists"},
		{MethodGet, "/gists/public"},
		{MethodGet, "/gists/starred"},
		{MethodGet, "/notifications"},
		{MethodGet, "/repositories"},
		{MethodGet, "/search/repositories"},
		{MethodGet, "/search/code"},
		{MethodGet, "/gitignore/templates"},
		{MethodGet, "/licenses"},
		{MethodGet, "/events"},

		{MethodGet, "/users/{username}"},
		{MethodGet, "/users/{username}/repos"},
		{MethodGet, "/users/{username}/followers"},
		{MethodGet, "/repos/{owner}/{repo}"},
		{MethodPatch, "/repos/{owner}/{repo}"},
		{MethodDelete, "/repos/{owner}/{repo}"},
		{MethodGet, "/repos/{owner}/{repo}/issues"},
		{MethodPost, "/repos/{owner}/{repo}/issues"},
		{MethodGet, "/repos/{owner}/{repo}/issues/{number}"},
		{MethodGet, "/repos/{owner}/{repo}/issues/comments"},
		{MethodGet, "/repos/{owner}/{repo}/issues/comments/{id}"},
		{MethodGet, "/repos/{owner}/{repo}/pulls"},
		{MethodGet, "/repos/{owner}/{repo}/pulls/{number}"},
		{MethodGet, "/repos/{owner}/{repo}/commits"},
		{MethodGet, "/repos/{owner}/{repo}/commits/{sha}"},
		{MethodGet, "/repos/{owner}/{repo}/releases"},
		{MethodGet, "/repos/{owner}/{repo}/releases/{id}"},
		{MethodGet, "/repos/{owner}/{repo}/releases/latest"},
		{MethodGet, "/repos/{owner}/{repo}/releases/tags/{tag}"},
		{MethodGet, "/repos/{owner}/{repo}/branches"},
		{MethodGet, "/repos/{owner}/{repo}/branches/{branch}"},
		{MethodGet, "/repos/{owner}/{repo}/contents/{path}"},
		{MethodPut, "/repos/{owner}/{repo}/subscription"},
		{MethodGet, "/gists/{id}"},
		{MethodPatch, "/gists/{id}"},
		{MethodDelete, "/gists/{id}"},
		{MethodGet, "/gists/{id}/commits"},
		{MethodGet, "/gists/{id}/comments"},
		{MethodGet, "/gists/{id}/comments/{commentId}"},
		{MethodGet, "/orgs/{org}"},
		{MethodGet, "/orgs/{org}/repos"},
		{MethodGet, "/orgs/{org}/members"},
		{MethodGet, "/orgs/{org}/members/{username}"},
		{MethodGet, "/teams/{id}"},
		{MethodGet, "/teams/{id}/members"},
		{MethodGet, "/teams/{id}/repos/{owner}/{repo}"},
		{MethodGet, "/authorizations/{id}"},
		{MethodGet, "/notifications/threads/{id}"},
		{MethodGet, "/notifications/threads/{id}/subscription"},
		{MethodGet, "/gitignore/templates/{name}"},
		{MethodGet, "/licenses/{license}"},

		{MethodGet, "/repos/{owner}/{repo}/releases/*/assets"},

		{MethodGet, "/repos/{owner}/{repo}/git/**"},
		{MethodGet, "/static/**"},
		{MethodGet, "/api/**/health"},

		{MethodGet, "/repos/{owner}/{repo}/contents/*.json"},
		{MethodGet, "/repos/{owner}/{repo}/archive/*.tar.gz"},
	}

	for _, endpoint := range endpoints {
		err := registry.Register(&Endpoint{
			method: endpoint.method,
			path:   endpoint.path,
		})
		if err != nil {
			b.Fatalf("failed to register endpoint %s %s: %v",
				endpoint.method, endpoint.path, err)
		}
	}

	matcher := NewDefaultEndpointMatcher(registry)

	benchCases := []struct {
		name   string
		method Method
		path   string
	}{
		{"static_root", MethodGet, "/"},
		{"static_shallow", MethodGet, "/user"},
		{"static_medium", MethodGet, "/user/repos"},
		{"static_deep", MethodGet, "/search/repositories"},

		{"param_1", MethodGet, "/users/octocat"},
		{"param_2", MethodGet, "/repos/octocat/hello-world"},
		{"param_3", MethodGet, "/repos/octocat/hello-world/issues/42"},
		{"param_4", MethodGet, "/repos/octocat/hello-world/pulls/7"},
		{"param_deep", MethodGet, "/teams/1/repos/octocat/hello-world"},
		{"param_with_static_suffix", MethodGet, "/repos/octocat/hello-world/issues"},
		{"param_with_static_child", MethodGet, "/users/octocat/repos"},

		{"static_over_param", MethodGet, "/gists/public"},
		{"static_over_param_deep", MethodGet, "/repos/octocat/hello-world/releases/latest"},
		{"static_over_param_comments", MethodGet, "/repos/octocat/hello-world/issues/comments"},

		{"wildcard", MethodGet, "/repos/octocat/hello-world/releases/v1.0/assets"},

		{"double_wildcard_zero", MethodGet, "/static"},
		{"double_wildcard_shallow", MethodGet, "/static/css/main.css"},
		{"double_wildcard_deep", MethodGet, "/static/js/vendor/lodash/core.min.js"},
		{"double_wildcard_param_prefix", MethodGet, "/repos/octocat/hello-world/git/refs/heads/main"},
		{"double_wildcard_middle", MethodGet, "/api/v1/v2/health"},
		{"double_wildcard_middle_zero", MethodGet, "/api/health"},

		{"pattern_json", MethodGet, "/repos/octocat/hello-world/contents/config.json"},
		{"pattern_targz", MethodGet, "/repos/octocat/hello-world/archive/v1.0.tar.gz"},

		{"method_get", MethodGet, "/repos/octocat/hello-world"},
		{"method_patch", MethodPatch, "/repos/octocat/hello-world"},
		{"method_delete", MethodDelete, "/repos/octocat/hello-world"},
		{"method_post", MethodPost, "/repos/octocat/hello-world/issues"},
		{"method_put", MethodPut, "/repos/octocat/hello-world/subscription"},

		{"not_found_shallow", MethodGet, "/notfound"},
		{"not_found_deep", MethodGet, "/repos/octocat/hello-world/unknown/path/here"},
		{"not_found_wrong_method", MethodDelete, "/users/octocat"},
	}

	for _, bc := range benchCases {
		b.Run(bc.name, func(b *testing.B) {
			req, _ := http.NewRequest(string(bc.method), bc.path, nil)
			recorder := httptest.NewRecorder()
			ctx := CreateContext(req, recorder)

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				ctx.Request().pathValues.reset()
				matcher.Match(ctx)
			}
		})
	}
}
