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
	"testing"
)

func TestRadixTreeMatcher_Match(t *testing.T) {
	t.Run("static routes", func(t *testing.T) {
		tree := NewRequestEndpointMatcher(nil)

		endpoints := []struct {
			method string
			path   string
		}{
			{"GET", "/"},
			{"GET", "/user"},
			{"GET", "/user/repos"},
			{"GET", "/user/orgs"},
			{"GET", "/user/followers"},
			{"GET", "/user/following"},
			{"GET", "/user/emails"},
			{"GET", "/user/keys"},
			{"GET", "/user/teams"},
			{"GET", "/user/issues"},
			{"POST", "/user/repos"},
			{"GET", "/users"},
			{"GET", "/gists"},
			{"GET", "/gists/public"},
			{"GET", "/gists/starred"},
			{"POST", "/gists"},
			{"GET", "/notifications"},
			{"PUT", "/notifications"},
			{"GET", "/authorizations"},
			{"POST", "/authorizations"},
			{"GET", "/repositories"},
			{"GET", "/repos"},
			{"GET", "/search/repositories"},
			{"GET", "/search/commits"},
			{"GET", "/search/code"},
			{"GET", "/search/issues"},
			{"GET", "/search/users"},
			{"GET", "/search/topics"},
			{"GET", "/search/labels"},
			{"GET", "/gitignore/templates"},
			{"GET", "/licenses"},
			{"GET", "/emojis"},
			{"GET", "/markdown"},
			{"POST", "/markdown"},
			{"GET", "/meta"},
			{"GET", "/rate_limit"},
			{"GET", "/feeds"},
			{"GET", "/events"},
		}

		for _, e := range endpoints {
			err := tree.addEndpoint(&Endpoint{method: Method(e.method), path: e.path})
			if err != nil {
				t.Fatalf("failed to add route %s %s: %v", e.method, e.path, err)
			}
		}

		testCases := []struct {
			method Method
			path   string
			match  bool
		}{
			// direct matches
			{"GET", "/", true},
			{"GET", "/user", true},
			{"GET", "/user/repos", true},
			{"GET", "/user/orgs", true},
			{"GET", "/user/followers", true},
			{"GET", "/user/following", true},
			{"GET", "/user/emails", true},
			{"GET", "/user/keys", true},
			{"GET", "/user/teams", true},
			{"GET", "/user/issues", true},
			{"POST", "/user/repos", true},
			{"GET", "/users", true},
			{"GET", "/gists", true},
			{"GET", "/gists/public", true},
			{"GET", "/gists/starred", true},
			{"POST", "/gists", true},
			{"GET", "/notifications", true},
			{"PUT", "/notifications", true},
			{"GET", "/authorizations", true},
			{"POST", "/authorizations", true},
			{"GET", "/repositories", true},
			{"GET", "/repos", true},
			{"GET", "/search/repositories", true},
			{"GET", "/search/commits", true},
			{"GET", "/search/code", true},
			{"GET", "/search/issues", true},
			{"GET", "/search/users", true},
			{"GET", "/search/topics", true},
			{"GET", "/search/labels", true},
			{"GET", "/gitignore/templates", true},
			{"GET", "/licenses", true},
			{"GET", "/emojis", true},
			{"GET", "/markdown", true},
			{"POST", "/markdown", true},
			{"GET", "/meta", true},
			{"GET", "/rate_limit", true},
			{"GET", "/feeds", true},
			{"GET", "/events", true},

			// wrong method
			{"DELETE", "/user/repos", false},

			// shared prefix partial matches must fail
			{"GET", "/user/follow", false},
			{"GET", "/use", false},
			{"GET", "/repo", false},
			{"GET", "/gists/star", false},
			{"GET", "/gitignore", false},
			{"GET", "/git", false},
			{"GET", "/search", false},
			{"GET", "/search/repo", false},
			{"GET", "/license", false},
			{"GET", "/event", false},
			{"GET", "/emoji", false},
			{"GET", "/rate", false},
			{"GET", "/mark", false},
			{"GET", "/met", false},

			// not found
			{"GET", "/notfound", false},
			{"GET", "", false},
		}

		for _, tc := range testCases {
			req, _ := http.NewRequest(string(tc.method), tc.path, nil)
			ctx := NewContext(req, nil)
			ep, ok := tree.Match(ctx)
			if ok != tc.match {
				t.Errorf("%s %s: got match=%v, want %v", tc.method, tc.path, ok, tc.match)
			}
			if ok && ep.method != tc.method {
				t.Errorf("%s %s: got method=%s, want %s", tc.method, tc.path, ep.method, tc.method)
			}
		}
	})

	t.Run("param routes", func(t *testing.T) {
		tree := NewRequestEndpointMatcher(nil)

		endpoints := []struct {
			method string
			path   string
		}{
			{"GET", "/users/{username}"},
			{"GET", "/users/{username}/repos"},
			{"GET", "/users/{username}/followers"},
			{"GET", "/repos/{owner}/{repo}"},
			{"GET", "/repos/{owner}/{repo}/issues"},
			{"GET", "/repos/{owner}/{repo}/issues/{number}"},
			{"GET", "/repos/{owner}/{repo}/pulls/{number}/comments"},
			{"GET", "/orgs/{org}/members/{username}"},
			{"GET", "/teams/{id}/repos/{owner}/{repo}"},
		}

		for _, e := range endpoints {
			err := tree.addEndpoint(&Endpoint{method: Method(e.method), path: e.path})
			if err != nil {
				t.Fatalf("failed to add route %s %s: %v", e.method, e.path, err)
			}
		}

		testCases := []struct {
			method Method
			path   string
			match  bool
			params map[string]string
		}{
			{"GET", "/users/octocat", true, map[string]string{"username": "octocat"}},
			{"GET", "/users/octocat/repos", true, map[string]string{"username": "octocat"}},
			{"GET", "/users/octocat/followers", true, map[string]string{"username": "octocat"}},
			{"GET", "/repos/octocat/hello-world", true, map[string]string{"owner": "octocat", "repo": "hello-world"}},
			{"GET", "/repos/octocat/hello-world/issues", true, map[string]string{"owner": "octocat", "repo": "hello-world"}},
			{"GET", "/repos/octocat/hello-world/issues/42", true, map[string]string{"owner": "octocat", "repo": "hello-world", "number": "42"}},
			{"GET", "/repos/octocat/hello-world/pulls/7/comments", true, map[string]string{"owner": "octocat", "repo": "hello-world", "number": "7"}},
			{"GET", "/orgs/github/members/octocat", true, map[string]string{"org": "github", "username": "octocat"}},
			{"GET", "/teams/1/repos/octocat/hello-world", true, map[string]string{"id": "1", "owner": "octocat", "repo": "hello-world"}},
			{"GET", "/users", false, nil},
			{"GET", "/repos/octocat", false, nil},
			{"GET", "/repos/octocat/hello-world/issues/42/labels/bug/extra", false, nil},
		}

		for _, tc := range testCases {
			req, _ := http.NewRequest(string(tc.method), tc.path, nil)
			ctx := NewContext(req, nil)
			_, ok := tree.Match(ctx)
			if ok != tc.match {
				t.Errorf("%s %s: got match=%v, want %v", tc.method, tc.path, ok, tc.match)
			}
			if ok && tc.params != nil {
				r := ctx.Request()
				for k, want := range tc.params {
					if got := r.PathValue(k); got != want {
						t.Errorf("%s %s: param %s: got %q, want %q", tc.method, tc.path, k, got, want)
					}
				}
			}
		}
	})

	t.Run("wildcard routes", func(t *testing.T) {
		tree := NewRequestEndpointMatcher(nil)

		endpoints := []struct {
			method string
			path   string
		}{
			{"GET", "/repos/{owner}/{repo}/releases/*/assets"},
			{"GET", "/orgs/{org}/hooks/*/pings"},
			{"GET", "/repos/{owner}/{repo}/deployments/*/statuses"},
		}

		for _, e := range endpoints {
			err := tree.addEndpoint(&Endpoint{method: Method(e.method), path: e.path})
			if err != nil {
				t.Fatalf("failed to add route %s %s: %v", e.method, e.path, err)
			}
		}

		testCases := []struct {
			method Method
			path   string
			match  bool
		}{
			{"GET", "/repos/octocat/hello-world/releases/v1.0/assets", true},
			{"GET", "/repos/octocat/hello-world/releases/latest/assets", true},
			{"GET", "/repos/octocat/hello-world/releases/assets", false},
			{"GET", "/repos/octocat/hello-world/releases/v1/v2/assets", false},
			{"GET", "/orgs/github/hooks/42/pings", true},
			{"GET", "/orgs/github/hooks/pings", false},
			{"GET", "/repos/octocat/hello-world/deployments/123/statuses", true},
			{"GET", "/repos/octocat/hello-world/deployments/statuses", false},
			{"GET", "/repos/octocat/hello-world/deployments/a/b/statuses", false},
		}

		for _, tc := range testCases {
			req, _ := http.NewRequest(string(tc.method), tc.path, nil)
			ctx := NewContext(req, nil)
			_, ok := tree.Match(ctx)
			if ok != tc.match {
				t.Errorf("%s %s: got match=%v, want %v", tc.method, tc.path, ok, tc.match)
			}
		}
	})

	t.Run("double wildcard routes", func(t *testing.T) {
		tree := NewRequestEndpointMatcher(nil)

		endpoints := []struct {
			method string
			path   string
		}{
			{"GET", "/repos/{owner}/{repo}/git/**"},
			{"GET", "/static/**"},
			{"GET", "/api/**/health"},
			{"GET", "/cdn/{vendor}/**"},
			{"GET", "/proxy/**/status"},
			{"GET", "/staticx"},
		}

		for _, e := range endpoints {
			err := tree.addEndpoint(&Endpoint{method: Method(e.method), path: e.path})
			if err != nil {
				t.Fatalf("failed to add route %s %s: %v", e.method, e.path, err)
			}
		}

		testCases := []struct {
			method Method
			path   string
			match  bool
		}{
			// /repos/{owner}/{repo}/git/**
			{"GET", "/repos/octocat/hello-world/git", true},
			{"GET", "/repos/octocat/hello-world/git/refs", true},
			{"GET", "/repos/octocat/hello-world/git/refs/heads/main", true},
			{"GET", "/repos/octocat/hello-world/git/commits/abc123", true},
			{"GET", "/repos/octocat/hello-world/git/trees/def456", true},
			{"GET", "/repos/octocat/hello-world/git/blobs/a1b2c3/raw", true},

			// /static/** zero-match and deep paths
			{"GET", "/static", true},
			{"GET", "/static/css/main.css", true},
			{"GET", "/static/js/vendor/lodash.min.js", true},

			// /static/** vs /staticx — split edge case
			{"GET", "/staticx", true},

			// /api/**/health — zero-match through deep
			{"GET", "/api/health", true},
			{"GET", "/api/v1/health", true},
			{"GET", "/api/v1/v2/v3/health", true},
			{"GET", "/api/health/health", true},
			{"GET", "/api/v1/check", false},

			// /cdn/{vendor}/** — param + double wildcard
			{"GET", "/cdn/cloudflare", true},
			{"GET", "/cdn/cloudflare/js/app.js", true},
			{"GET", "/cdn/aws/images/logo.png", true},

			// /proxy/**/status — double wildcard in middle
			{"GET", "/proxy/status", true},
			{"GET", "/proxy/us-east/status", true},
			{"GET", "/proxy/us-east/zone-1/status", true},
			{"GET", "/proxy/us-east/zone-1/zone-2/status", true},
			{"GET", "/proxy/us-east/health", false},
		}

		for _, tc := range testCases {
			req, _ := http.NewRequest(string(tc.method), tc.path, nil)
			ctx := NewContext(req, nil)
			_, ok := tree.Match(ctx)
			if ok != tc.match {
				t.Errorf("%s %s: got match=%v, want %v", tc.method, tc.path, ok, tc.match)
			}
		}
	})

	t.Run("pattern routes", func(t *testing.T) {
		tree := NewRequestEndpointMatcher(nil)

		endpoints := []struct {
			method string
			path   string
		}{
			{"GET", "/repos/{owner}/{repo}/contents/*.json"},
			{"GET", "/repos/{owner}/{repo}/archive/*.tar.gz"},
			{"GET", "/repos/{owner}/{repo}/contents/*.json/meta"},
			{"GET", "/downloads/release-v?.zip"},
			{"GET", "/exports/report-??-*.csv"},
			{"GET", "/assets/logo*"},
		}

		for _, e := range endpoints {
			err := tree.addEndpoint(&Endpoint{method: Method(e.method), path: e.path})
			if err != nil {
				t.Fatalf("failed to add route %s %s: %v", e.method, e.path, err)
			}
		}

		testCases := []struct {
			method Method
			path   string
			match  bool
		}{
			// *.json
			{"GET", "/repos/octocat/hello-world/contents/config.json", true},
			{"GET", "/repos/octocat/hello-world/contents/package.json", true},
			{"GET", "/repos/octocat/hello-world/contents/.json", true},
			{"GET", "/repos/octocat/hello-world/contents/readme.md", false},
			{"GET", "/repos/octocat/hello-world/contents/data.jsonx", false},

			// *.tar.gz
			{"GET", "/repos/octocat/hello-world/archive/v1.0.tar.gz", true},
			{"GET", "/repos/octocat/hello-world/archive/latest.tar.gz", true},
			{"GET", "/repos/octocat/hello-world/archive/v1.0.zip", false},

			// *.json/meta — pattern child reuse
			{"GET", "/repos/octocat/hello-world/contents/schema.json/meta", true},
			{"GET", "/repos/octocat/hello-world/contents/data.json/meta", true},

			// release-v?.zip — ? matches single char
			{"GET", "/downloads/release-v1.zip", true},
			{"GET", "/downloads/release-v9.zip", true},
			{"GET", "/downloads/release-v12.zip", false},
			{"GET", "/downloads/release-v.zip", false},

			// report-??-*.csv — ?? + * combined
			{"GET", "/exports/report-US-2024.csv", true},
			{"GET", "/exports/report-EU-sales.csv", true},
			{"GET", "/exports/report-A-2024.csv", false},
			{"GET", "/exports/report-USA-2024.csv", false},

			// logo* — trailing star
			{"GET", "/assets/logo", true},
			{"GET", "/assets/logo-dark", true},
			{"GET", "/assets/logo192.png", true},
			{"GET", "/assets/log", false},
		}

		for _, tc := range testCases {
			req, _ := http.NewRequest(string(tc.method), tc.path, nil)
			ctx := NewContext(req, nil)
			_, ok := tree.Match(ctx)
			if ok != tc.match {
				t.Errorf("%s %s: got match=%v, want %v", tc.method, tc.path, ok, tc.match)
			}
		}
	})

	t.Run("priority: static > pattern > param > wildcard > double wildcard", func(t *testing.T) {
		tree := NewRequestEndpointMatcher(nil)

		endpoints := []struct {
			method string
			path   string
		}{
			{"GET", "/repos/octocat/hello-world"},
			{"GET", "/repos/{owner}/{repo}/contents/*.json"},
			{"GET", "/repos/{owner}/{repo}"},
			{"GET", "/repos/{owner}/*"},
			{"GET", "/repos/**"},
			{"GET", "/gists/public"},
			{"GET", "/gists/starred"},
			{"GET", "/gists/{id}"},
			{"GET", "/repos/{owner}/{repo}/releases/latest"},
			{"GET", "/repos/{owner}/{repo}/releases/{id}"},
			{"GET", "/repos/{owner}/{repo}/issues/comments"},
			{"GET", "/repos/{owner}/{repo}/issues/{number}"},
		}

		for _, e := range endpoints {
			err := tree.addEndpoint(&Endpoint{method: Method(e.method), path: e.path})
			if err != nil {
				t.Fatalf("failed to add route %s %s: %v", e.method, e.path, err)
			}
		}

		testCases := []struct {
			path        string
			wantPattern string
		}{
			// static wins over param
			{"/repos/octocat/hello-world", "/repos/octocat/hello-world"},
			{"/gists/public", "/gists/public"},
			{"/gists/starred", "/gists/starred"},
			{"/repos/octocat/hello-world/releases/latest", "/repos/{owner}/{repo}/releases/latest"},
			{"/repos/octocat/hello-world/issues/comments", "/repos/{owner}/{repo}/issues/comments"},

			// pattern wins over param
			{"/repos/octocat/hello-world/contents/config.json", "/repos/{owner}/{repo}/contents/*.json"},

			// param wins over wildcard
			{"/repos/facebook/react", "/repos/{owner}/{repo}"},
			{"/gists/abc123", "/gists/{id}"},
			{"/repos/octocat/hello-world/releases/42", "/repos/{owner}/{repo}/releases/{id}"},
			{"/repos/octocat/hello-world/issues/99", "/repos/{owner}/{repo}/issues/{number}"},

			// double wildcard catches the rest
			{"/repos/octocat/hello-world/deep/nested/path", "/repos/**"},
		}

		for _, tc := range testCases {
			req, _ := http.NewRequest("GET", tc.path, nil)
			ctx := NewContext(req, nil)
			ep, ok := tree.Match(ctx)
			if !ok {
				t.Errorf("GET %s: no match, want %s", tc.path, tc.wantPattern)
				continue
			}
			if ep.path != tc.wantPattern {
				t.Errorf("GET %s: matched %s, want %s", tc.path, ep.path, tc.wantPattern)
			}
		}
	})

	t.Run("backtracking", func(t *testing.T) {
		tree := NewRequestEndpointMatcher(nil)

		endpoints := []struct {
			method string
			path   string
		}{
			// param vs wildcard
			{"GET", "/repos/{owner}/{repo}/issues/{number}"},
			{"GET", "/repos/{owner}/{repo}/issues/*/reactions"},

			// pattern vs pattern
			{"GET", "/repos/{owner}/{repo}/contents/*.json/download"},
			{"GET", "/repos/{owner}/{repo}/contents/*.md/preview"},

			// param vs double wildcard
			{"GET", "/repos/{owner}/{repo}/branches/{branch}"},
			{"GET", "/repos/{owner}/{repo}/git/**"},
		}

		for _, e := range endpoints {
			err := tree.addEndpoint(&Endpoint{method: Method(e.method), path: e.path})
			if err != nil {
				t.Fatalf("failed to add route %s %s: %v", e.method, e.path, err)
			}
		}

		testCases := []struct {
			path        string
			match       bool
			wantPattern string
		}{
			// param matches
			{"/repos/octocat/hello-world/issues/42", true, "/repos/{owner}/{repo}/issues/{number}"},
			{"/repos/octocat/hello-world/branches/main", true, "/repos/{owner}/{repo}/branches/{branch}"},

			// param fails on suffix, wildcard succeeds
			{"/repos/octocat/hello-world/issues/42/reactions", true, "/repos/{owner}/{repo}/issues/*/reactions"},

			// first pattern fails, second pattern succeeds
			{"/repos/octocat/hello-world/contents/readme.md/preview", true, "/repos/{owner}/{repo}/contents/*.md/preview"},
			{"/repos/octocat/hello-world/contents/config.json/download", true, "/repos/{owner}/{repo}/contents/*.json/download"},

			// pattern mismatch — no backtrack target
			{"/repos/octocat/hello-world/contents/readme.md/download", false, ""},
			{"/repos/octocat/hello-world/contents/config.json/preview", false, ""},

			// double wildcard catches deep git paths
			{"/repos/octocat/hello-world/git/refs/heads/main", true, "/repos/{owner}/{repo}/git/**"},
			{"/repos/octocat/hello-world/git/commits/abc123", true, "/repos/{owner}/{repo}/git/**"},
		}

		for _, tc := range testCases {
			req, _ := http.NewRequest("GET", tc.path, nil)
			ctx := NewContext(req, nil)
			ep, ok := tree.Match(ctx)
			if ok != tc.match {
				t.Errorf("GET %s: got match=%v, want %v", tc.path, ok, tc.match)
				continue
			}
			if ok && ep.path != tc.wantPattern {
				t.Errorf("GET %s: matched %s, want %s", tc.path, ep.path, tc.wantPattern)
			}
		}
	})

	t.Run("multiple methods same path", func(t *testing.T) {
		tree := NewRequestEndpointMatcher(nil)

		endpoints := []struct {
			method string
			path   string
		}{
			{"GET", "/repos/{owner}/{repo}"},
			{"PATCH", "/repos/{owner}/{repo}"},
			{"DELETE", "/repos/{owner}/{repo}"},
			{"GET", "/repos/{owner}/{repo}/issues"},
			{"POST", "/repos/{owner}/{repo}/issues"},
			{"GET", "/gists/{id}"},
			{"PATCH", "/gists/{id}"},
			{"DELETE", "/gists/{id}"},
			{"PUT", "/repos/{owner}/{repo}/subscription"},
			{"HEAD", "/repos/{owner}/{repo}/subscription"},
			{"OPTIONS", "/repos/{owner}/{repo}/subscription"},
			{"CONNECT", "/repos/{owner}/{repo}/subscription"},
			{"TRACE", "/repos/{owner}/{repo}/subscription"},
		}

		for _, e := range endpoints {
			err := tree.addEndpoint(&Endpoint{method: Method(e.method), path: e.path})
			if err != nil {
				t.Fatalf("failed to add route %s %s: %v", e.method, e.path, err)
			}
		}

		testCases := []struct {
			method Method
			path   string
			match  bool
		}{
			// repos - GET, PATCH, DELETE
			{"GET", "/repos/octocat/hello-world", true},
			{"PATCH", "/repos/octocat/hello-world", true},
			{"DELETE", "/repos/octocat/hello-world", true},
			{"PUT", "/repos/octocat/hello-world", false},
			{"POST", "/repos/octocat/hello-world", false},

			// issues - GET, POST
			{"GET", "/repos/octocat/hello-world/issues", true},
			{"POST", "/repos/octocat/hello-world/issues", true},
			{"DELETE", "/repos/octocat/hello-world/issues", false},

			// gists - GET, PATCH, DELETE
			{"GET", "/gists/abc123", true},
			{"PATCH", "/gists/abc123", true},
			{"DELETE", "/gists/abc123", true},
			{"POST", "/gists/abc123", false},

			// subscription - PUT, HEAD, OPTIONS, CONNECT, TRACE
			{"PUT", "/repos/octocat/hello-world/subscription", true},
			{"HEAD", "/repos/octocat/hello-world/subscription", true},
			{"OPTIONS", "/repos/octocat/hello-world/subscription", true},
			{"CONNECT", "/repos/octocat/hello-world/subscription", true},
			{"TRACE", "/repos/octocat/hello-world/subscription", true},
			{"GET", "/repos/octocat/hello-world/subscription", false},
		}

		for _, tc := range testCases {
			req, _ := http.NewRequest(string(tc.method), tc.path, nil)
			ctx := NewContext(req, nil)
			ep, ok := tree.Match(ctx)
			if ok != tc.match {
				t.Errorf("%s %s: got match=%v, want %v", tc.method, tc.path, ok, tc.match)
				continue
			}
			if ok && ep.method != tc.method {
				t.Errorf("%s %s: got method=%s", tc.method, tc.path, ep.method)
			}
		}
	})

	t.Run("mixed route types in same tree", func(t *testing.T) {
		tree := NewRequestEndpointMatcher(nil)

		endpoints := []struct {
			method string
			path   string
		}{
			// static
			{"GET", "/"},
			{"GET", "/user"},
			{"GET", "/user/repos"},
			{"GET", "/user/followers"},
			{"GET", "/user/following"},
			{"GET", "/users"},
			{"GET", "/gists"},
			{"GET", "/gists/public"},
			{"GET", "/gists/starred"},
			{"GET", "/notifications"},
			{"GET", "/repositories"},
			{"GET", "/search/repositories"},
			{"GET", "/search/code"},
			{"GET", "/gitignore/templates"},
			{"GET", "/licenses"},
			{"GET", "/events"},

			// param
			{"GET", "/users/{username}"},
			{"GET", "/users/{username}/repos"},
			{"GET", "/users/{username}/followers"},
			{"GET", "/repos/{owner}/{repo}"},
			{"PATCH", "/repos/{owner}/{repo}"},
			{"DELETE", "/repos/{owner}/{repo}"},
			{"GET", "/repos/{owner}/{repo}/issues"},
			{"POST", "/repos/{owner}/{repo}/issues"},
			{"GET", "/repos/{owner}/{repo}/issues/{number}"},
			{"GET", "/repos/{owner}/{repo}/issues/comments"},
			{"GET", "/repos/{owner}/{repo}/issues/comments/{id}"},
			{"GET", "/repos/{owner}/{repo}/pulls"},
			{"GET", "/repos/{owner}/{repo}/pulls/{number}"},
			{"GET", "/repos/{owner}/{repo}/commits"},
			{"GET", "/repos/{owner}/{repo}/commits/{sha}"},
			{"GET", "/repos/{owner}/{repo}/releases"},
			{"GET", "/repos/{owner}/{repo}/releases/{id}"},
			{"GET", "/repos/{owner}/{repo}/releases/latest"},
			{"GET", "/repos/{owner}/{repo}/releases/tags/{tag}"},
			{"GET", "/repos/{owner}/{repo}/branches"},
			{"GET", "/repos/{owner}/{repo}/branches/{branch}"},
			{"GET", "/repos/{owner}/{repo}/contents/{path}"},
			{"PUT", "/repos/{owner}/{repo}/subscription"},
			{"GET", "/gists/{id}"},
			{"PATCH", "/gists/{id}"},
			{"DELETE", "/gists/{id}"},
			{"GET", "/gists/{id}/commits"},
			{"GET", "/gists/{id}/comments"},
			{"GET", "/gists/{id}/comments/{commentId}"},
			{"GET", "/orgs/{org}"},
			{"GET", "/orgs/{org}/repos"},
			{"GET", "/orgs/{org}/members"},
			{"GET", "/orgs/{org}/members/{username}"},
			{"GET", "/teams/{id}"},
			{"GET", "/teams/{id}/members"},
			{"GET", "/teams/{id}/repos/{owner}/{repo}"},
			{"GET", "/authorizations/{id}"},
			{"GET", "/notifications/threads/{id}"},
			{"GET", "/notifications/threads/{id}/subscription"},
			{"GET", "/gitignore/templates/{name}"},
			{"GET", "/licenses/{license}"},

			// wildcard
			{"GET", "/repos/{owner}/{repo}/releases/*/assets"},

			// double wildcard
			{"GET", "/repos/{owner}/{repo}/git/**"},
			{"GET", "/static/**"},

			// pattern
			{"GET", "/repos/{owner}/{repo}/contents/*.json"},
			{"GET", "/repos/{owner}/{repo}/archive/*.tar.gz"},
		}

		for _, e := range endpoints {
			err := tree.addEndpoint(&Endpoint{method: Method(e.method), path: e.path})
			if err != nil {
				t.Fatalf("failed to add route %s %s: %v", e.method, e.path, err)
			}
		}

		testCases := []struct {
			method      Method
			path        string
			match       bool
			wantPattern string
		}{
			// static
			{"GET", "/user", true, "/user"},
			{"GET", "/gists/public", true, "/gists/public"},
			{"GET", "/search/repositories", true, "/search/repositories"},
			{"GET", "/gitignore/templates", true, "/gitignore/templates"},
			{"GET", "/events", true, "/events"},

			// static vs param (static wins)
			{"GET", "/gists/public", true, "/gists/public"},
			{"GET", "/gists/abc123", true, "/gists/{id}"},
			{"GET", "/gists/starred", true, "/gists/starred"},
			{"GET", "/repos/octocat/hello-world/releases/latest", true, "/repos/{owner}/{repo}/releases/latest"},
			{"GET", "/repos/octocat/hello-world/releases/42", true, "/repos/{owner}/{repo}/releases/{id}"},
			{"GET", "/repos/octocat/hello-world/issues/comments", true, "/repos/{owner}/{repo}/issues/comments"},
			{"GET", "/repos/octocat/hello-world/issues/42", true, "/repos/{owner}/{repo}/issues/{number}"},

			// param
			{"GET", "/users/octocat", true, "/users/{username}"},
			{"GET", "/users/octocat/repos", true, "/users/{username}/repos"},
			{"GET", "/repos/octocat/hello-world", true, "/repos/{owner}/{repo}"},
			{"GET", "/repos/octocat/hello-world/issues/42", true, "/repos/{owner}/{repo}/issues/{number}"},
			{"GET", "/repos/octocat/hello-world/pulls/7", true, "/repos/{owner}/{repo}/pulls/{number}"},
			{"GET", "/repos/octocat/hello-world/commits/abc123", true, "/repos/{owner}/{repo}/commits/{sha}"},
			{"GET", "/repos/octocat/hello-world/branches/main", true, "/repos/{owner}/{repo}/branches/{branch}"},
			{"GET", "/repos/octocat/hello-world/releases/tags/v1.0", true, "/repos/{owner}/{repo}/releases/tags/{tag}"},
			{"GET", "/repos/octocat/hello-world/contents/README.md", true, "/repos/{owner}/{repo}/contents/{path}"},
			{"GET", "/orgs/github/members/octocat", true, "/orgs/{org}/members/{username}"},
			{"GET", "/teams/1/repos/octocat/hello-world", true, "/teams/{id}/repos/{owner}/{repo}"},
			{"GET", "/gists/abc123/comments/5", true, "/gists/{id}/comments/{commentId}"},
			{"GET", "/notifications/threads/77/subscription", true, "/notifications/threads/{id}/subscription"},

			// multiple methods
			{"GET", "/repos/octocat/hello-world", true, "/repos/{owner}/{repo}"},
			{"PATCH", "/repos/octocat/hello-world", true, "/repos/{owner}/{repo}"},
			{"DELETE", "/repos/octocat/hello-world", true, "/repos/{owner}/{repo}"},
			{"POST", "/repos/octocat/hello-world/issues", true, "/repos/{owner}/{repo}/issues"},
			{"PUT", "/repos/octocat/hello-world/subscription", true, "/repos/{owner}/{repo}/subscription"},

			// wildcard
			{"GET", "/repos/octocat/hello-world/releases/v1.0/assets", true, "/repos/{owner}/{repo}/releases/*/assets"},
			{"GET", "/repos/octocat/hello-world/releases/v1/v2/assets", false, ""},

			// param catches single segment at wildcard position (releases/{id} matches "assets")
			{"GET", "/repos/octocat/hello-world/releases/assets", true, "/repos/{owner}/{repo}/releases/{id}"},

			// double wildcard
			{"GET", "/repos/octocat/hello-world/git/refs", true, "/repos/{owner}/{repo}/git/**"},
			{"GET", "/repos/octocat/hello-world/git/refs/heads/main", true, "/repos/{owner}/{repo}/git/**"},
			{"GET", "/static/css/main.css", true, "/static/**"},
			{"GET", "/static/js/vendor/lodash.js", true, "/static/**"},

			// pattern vs param (pattern wins for *.json)
			{"GET", "/repos/octocat/hello-world/contents/config.json", true, "/repos/{owner}/{repo}/contents/*.json"},
			// param catches non-json files
			{"GET", "/repos/octocat/hello-world/contents/README.md", true, "/repos/{owner}/{repo}/contents/{path}"},
			{"GET", "/repos/octocat/hello-world/archive/v1.0.tar.gz", true, "/repos/{owner}/{repo}/archive/*.tar.gz"},
			{"GET", "/repos/octocat/hello-world/archive/v1.0.zip", false, ""},

			// not found
			{"GET", "/repos/octocat", false, ""},
			{"GET", "/unknown/path", false, ""},
			{"DELETE", "/users/octocat", false, ""},
			{"POST", "/gists/abc123", false, ""},
		}

		for _, tc := range testCases {
			req, _ := http.NewRequest(string(tc.method), tc.path, nil)
			ctx := NewContext(req, nil)
			ep, ok := tree.Match(ctx)
			if ok != tc.match {
				t.Errorf("%s %s: got match=%v, want %v", tc.method, tc.path, ok, tc.match)
				continue
			}
			if ok && ep.path != tc.wantPattern {
				t.Errorf("%s %s: matched %s, want %s", tc.method, tc.path, ep.path, tc.wantPattern)
			}
		}
	})

	t.Run("normalization", func(t *testing.T) {
		tree := NewRequestEndpointMatcher(nil)

		endpoints := []struct {
			method string
			path   string
		}{
			// trailing slash on match
			{"GET", "/repos/{owner}/{repo}"},
			{"GET", "/users/{username}/repos"},
			{"GET", "/repos/{owner}/{repo}/issues/{number}"},

			// missing leading slash on add
			{"GET", "gists"},
			{"GET", "orgs/{org}/members"},

			// trailing slash on add
			{"GET", "/notifications/"},
			{"GET", "/search/repositories/"},
		}

		for _, e := range endpoints {
			err := tree.addEndpoint(&Endpoint{method: Method(e.method), path: e.path})
			if err != nil {
				t.Fatalf("failed to add route %s %s: %v", e.method, e.path, err)
			}
		}

		testCases := []struct {
			path  string
			match bool
		}{
			// trailing slash stripped on match
			{"/repos/octocat/hello-world/", true},
			{"/users/octocat/repos/", true},
			{"/repos/octocat/hello-world/issues/42/", true},

			// missing leading slash normalized on add
			{"/gists", true},
			{"/orgs/github/members", true},

			// trailing slash normalized on add
			{"/notifications", true},
			{"/search/repositories", true},

			// root path
			{"/", false},
		}

		for _, tc := range testCases {
			req, _ := http.NewRequest("GET", tc.path, nil)
			ctx := NewContext(req, nil)
			_, ok := tree.Match(ctx)
			if ok != tc.match {
				t.Errorf("GET %s: got match=%v, want %v", tc.path, ok, tc.match)
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

			for _, e := range endpoints {
				t.Run(e.name, func(t *testing.T) {
					tree := NewRequestEndpointMatcher(nil)
					err := tree.addEndpoint(&Endpoint{"GET", e.path, nil})
					if err == nil {
						t.Errorf("expected error for path %q", e.path)
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

			for _, e := range endpoints {
				t.Run(e.name, func(t *testing.T) {
					tree := NewRequestEndpointMatcher(nil)
					err := tree.addEndpoint(&Endpoint{"GET", e.path, nil})
					if err != nil {
						t.Errorf("unexpected error for path %q: %v", e.path, err)
					}
				})
			}
		})

		t.Run("duplicate routes", func(t *testing.T) {
			endpoints := []struct {
				name   string
				method string
				path   string
			}{
				{"same static route", "GET", "/gists/public"},
				{"same param route", "GET", "/repos/{owner}/{repo}"},
				{"same deep param route", "POST", "/repos/{owner}/{repo}/issues"},
			}

			for _, e := range endpoints {
				t.Run(e.name, func(t *testing.T) {
					tree := NewRequestEndpointMatcher(nil)
					_ = tree.addEndpoint(&Endpoint{Method(e.method), e.path, nil})
					err := tree.addEndpoint(&Endpoint{Method(e.method), e.path, nil})
					if err == nil {
						t.Errorf("expected error for duplicate %s %s", e.method, e.path)
					}
				})
			}
		})

		t.Run("different method same path is ok", func(t *testing.T) {
			tree := NewRequestEndpointMatcher(nil)

			endpoints := []struct {
				method string
				path   string
			}{
				{"GET", "/repos/{owner}/{repo}"},
				{"DELETE", "/repos/{owner}/{repo}"},
				{"PATCH", "/repos/{owner}/{repo}"},
				{"GET", "/repos/{owner}/{repo}/issues"},
				{"POST", "/repos/{owner}/{repo}/issues"},
				{"GET", "/gists"},
				{"POST", "/gists"},
				{"GET", "/gists/{id}"},
				{"PATCH", "/gists/{id}"},
				{"DELETE", "/gists/{id}"},
			}

			for _, e := range endpoints {
				err := tree.addEndpoint(&Endpoint{Method(e.method), e.path, nil})
				if err != nil {
					t.Errorf("unexpected error for %s %s: %v", e.method, e.path, err)
				}
			}
		})
	})

}

func BenchmarkRadixTreeMatcher_Match(b *testing.B) {
	tree := NewRequestEndpointMatcher(nil)

	endpoints := []struct {
		method string
		path   string
	}{
		// static
		{"GET", "/"},
		{"GET", "/user"},
		{"GET", "/user/repos"},
		{"GET", "/user/followers"},
		{"GET", "/user/following"},
		{"GET", "/users"},
		{"GET", "/gists"},
		{"GET", "/gists/public"},
		{"GET", "/gists/starred"},
		{"GET", "/notifications"},
		{"GET", "/repositories"},
		{"GET", "/search/repositories"},
		{"GET", "/search/code"},
		{"GET", "/gitignore/templates"},
		{"GET", "/licenses"},
		{"GET", "/events"},

		// param
		{"GET", "/users/{username}"},
		{"GET", "/users/{username}/repos"},
		{"GET", "/users/{username}/followers"},
		{"GET", "/repos/{owner}/{repo}"},
		{"PATCH", "/repos/{owner}/{repo}"},
		{"DELETE", "/repos/{owner}/{repo}"},
		{"GET", "/repos/{owner}/{repo}/issues"},
		{"POST", "/repos/{owner}/{repo}/issues"},
		{"GET", "/repos/{owner}/{repo}/issues/{number}"},
		{"GET", "/repos/{owner}/{repo}/issues/comments"},
		{"GET", "/repos/{owner}/{repo}/issues/comments/{id}"},
		{"GET", "/repos/{owner}/{repo}/pulls"},
		{"GET", "/repos/{owner}/{repo}/pulls/{number}"},
		{"GET", "/repos/{owner}/{repo}/commits"},
		{"GET", "/repos/{owner}/{repo}/commits/{sha}"},
		{"GET", "/repos/{owner}/{repo}/releases"},
		{"GET", "/repos/{owner}/{repo}/releases/{id}"},
		{"GET", "/repos/{owner}/{repo}/releases/latest"},
		{"GET", "/repos/{owner}/{repo}/releases/tags/{tag}"},
		{"GET", "/repos/{owner}/{repo}/branches"},
		{"GET", "/repos/{owner}/{repo}/branches/{branch}"},
		{"GET", "/repos/{owner}/{repo}/contents/{path}"},
		{"PUT", "/repos/{owner}/{repo}/subscription"},
		{"GET", "/gists/{id}"},
		{"PATCH", "/gists/{id}"},
		{"DELETE", "/gists/{id}"},
		{"GET", "/gists/{id}/commits"},
		{"GET", "/gists/{id}/comments"},
		{"GET", "/gists/{id}/comments/{commentId}"},
		{"GET", "/orgs/{org}"},
		{"GET", "/orgs/{org}/repos"},
		{"GET", "/orgs/{org}/members"},
		{"GET", "/orgs/{org}/members/{username}"},
		{"GET", "/teams/{id}"},
		{"GET", "/teams/{id}/members"},
		{"GET", "/teams/{id}/repos/{owner}/{repo}"},
		{"GET", "/authorizations/{id}"},
		{"GET", "/notifications/threads/{id}"},
		{"GET", "/notifications/threads/{id}/subscription"},
		{"GET", "/gitignore/templates/{name}"},
		{"GET", "/licenses/{license}"},

		// wildcard
		{"GET", "/repos/{owner}/{repo}/releases/*/assets"},

		// double wildcard
		{"GET", "/repos/{owner}/{repo}/git/**"},
		{"GET", "/static/**"},
		{"GET", "/api/**/health"},

		// pattern
		{"GET", "/repos/{owner}/{repo}/contents/*.json"},
		{"GET", "/repos/{owner}/{repo}/archive/*.tar.gz"},
	}

	for _, e := range endpoints {
		err := tree.addEndpoint(&Endpoint{method: Method(e.method), path: e.path})
		if err != nil {
			b.Fatalf("failed to add route %s %s: %v", e.method, e.path, err)
		}
	}

	benchCases := []struct {
		name   string
		method string
		path   string
	}{
		// static — shallow
		{"static_root", "GET", "/"},
		{"static_shallow", "GET", "/user"},
		{"static_medium", "GET", "/user/repos"},
		{"static_deep", "GET", "/search/repositories"},

		// param — increasing depth
		{"param_1", "GET", "/users/octocat"},
		{"param_2", "GET", "/repos/octocat/hello-world"},
		{"param_3", "GET", "/repos/octocat/hello-world/issues/42"},
		{"param_4", "GET", "/repos/octocat/hello-world/pulls/7"},
		{"param_deep", "GET", "/teams/1/repos/octocat/hello-world"},
		{"param_with_static_suffix", "GET", "/repos/octocat/hello-world/issues"},
		{"param_with_static_child", "GET", "/users/octocat/repos"},

		// static vs param — static should win
		{"static_over_param", "GET", "/gists/public"},
		{"static_over_param_deep", "GET", "/repos/octocat/hello-world/releases/latest"},
		{"static_over_param_comments", "GET", "/repos/octocat/hello-world/issues/comments"},

		// wildcard
		{"wildcard", "GET", "/repos/octocat/hello-world/releases/v1.0/assets"},

		// double wildcard — increasing depth
		{"double_wildcard_zero", "GET", "/static"},
		{"double_wildcard_shallow", "GET", "/static/css/main.css"},
		{"double_wildcard_deep", "GET", "/static/js/vendor/lodash/core.min.js"},
		{"double_wildcard_param_prefix", "GET", "/repos/octocat/hello-world/git/refs/heads/main"},
		{"double_wildcard_middle", "GET", "/api/v1/v2/health"},
		{"double_wildcard_middle_zero", "GET", "/api/health"},

		// pattern
		{"pattern_json", "GET", "/repos/octocat/hello-world/contents/config.json"},
		{"pattern_targz", "GET", "/repos/octocat/hello-world/archive/v1.0.tar.gz"},

		// multiple methods — same path different method
		{"method_get", "GET", "/repos/octocat/hello-world"},
		{"method_patch", "PATCH", "/repos/octocat/hello-world"},
		{"method_delete", "DELETE", "/repos/octocat/hello-world"},
		{"method_post", "POST", "/repos/octocat/hello-world/issues"},
		{"method_put", "PUT", "/repos/octocat/hello-world/subscription"},

		// not found
		{"not_found_shallow", "GET", "/notfound"},
		{"not_found_deep", "GET", "/repos/octocat/hello-world/unknown/path/here"},
		{"not_found_wrong_method", "DELETE", "/users/octocat"},
	}

	for _, bc := range benchCases {
		b.Run(bc.name, func(b *testing.B) {
			req, _ := http.NewRequest(bc.method, bc.path, nil)
			ctx := NewContext(req, nil)

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				ctx.Request().pathValues.reset()
				tree.Match(ctx)
			}
		})
	}
}
