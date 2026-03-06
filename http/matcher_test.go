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
		tree := newRadixTreeMatcher()
		tree.addEndpoint(&Endpoint{"GET", "/users", nil})
		tree.addEndpoint(&Endpoint{"GET", "/users/profile", nil})
		tree.addEndpoint(&Endpoint{"POST", "/users", nil})
		tree.addEndpoint(&Endpoint{"GET", "/", nil})
		tree.addEndpoint(&Endpoint{"GET", "/a/b/c/d/e", nil})
		tree.addEndpoint(&Endpoint{"GET", "/abcdef", nil})
		tree.addEndpoint(&Endpoint{"GET", "/abc", nil})
		tree.addEndpoint(&Endpoint{"GET", "/abcxyz", nil})

		testCases := []struct {
			method Method
			path   string
			match  bool
		}{
			{"GET", "/users", true},
			{"GET", "/users/profile", true},
			{"POST", "/users", true},
			{"DELETE", "/users", false},
			{"GET", "/", true},
			{"GET", "/a/b/c/d/e", true},
			{"GET", "/a/b/c/d", false},
			{"GET", "/a/b/c/d/e/f", false},
			{"GET", "/notfound", false},
			{"GET", "", false},
			{"GET", "/abcdef", true},
			{"GET", "/abc", true},
			{"GET", "/abcxyz", true},
			{"GET", "/abcd", false},
			{"GET", "/ab", false},
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

	t.Run("static routes with shared prefixes", func(t *testing.T) {
		tree := newRadixTreeMatcher()
		tree.addEndpoint(&Endpoint{"GET", "/users", nil})
		tree.addEndpoint(&Endpoint{"GET", "/users/profile", nil})
		tree.addEndpoint(&Endpoint{"GET", "/users/password", nil})
		tree.addEndpoint(&Endpoint{"GET", "/us", nil})
		tree.addEndpoint(&Endpoint{"GET", "/api/v1/items", nil})
		tree.addEndpoint(&Endpoint{"GET", "/api/v2/items", nil})
		tree.addEndpoint(&Endpoint{"GET", "/api/v1/orders", nil})
		tree.addEndpoint(&Endpoint{"GET", "/abcdef", nil})
		tree.addEndpoint(&Endpoint{"GET", "/abc", nil})
		tree.addEndpoint(&Endpoint{"GET", "/abcxyz", nil})

		testCases := []struct {
			path  string
			match bool
		}{
			{"/users", true},
			{"/users/profile", true},
			{"/users/password", true},
			{"/us", true},
			{"/u", false},
			{"/user", false},
			{"/users/p", false},
			{"/api/v1/items", true},
			{"/api/v2/items", true},
			{"/api/v1/orders", true},
			{"/api/v1", false},
			{"/api/v3/items", false},
			{"/abcdef", true},
			{"/abc", true},
			{"/abcxyz", true},
			{"/abcd", false},
			{"/ab", false},
			{"/abqz", false},
			{"/usez", false},
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

	t.Run("param routes", func(t *testing.T) {
		tree := newRadixTreeMatcher()
		tree.addEndpoint(&Endpoint{"GET", "/users/{id}", nil})
		tree.addEndpoint(&Endpoint{"GET", "/users/{id}/posts/{postId}", nil})
		tree.addEndpoint(&Endpoint{"GET", "/{root}", nil})

		testCases := []struct {
			method Method
			path   string
			match  bool
			params map[string]string
		}{
			{"GET", "/users/42", true, map[string]string{"id": "42"}},
			{"GET", "/users/42/posts/7", true, map[string]string{"id": "42", "postId": "7"}},
			{"GET", "/users/abc/posts/xyz", true, map[string]string{"id": "abc", "postId": "xyz"}},
			{"GET", "/users", true, map[string]string{"root": "users"}},
			{"GET", "/users/42/posts", false, nil},
			{"GET", "/anything", true, map[string]string{"root": "anything"}},
			{"GET", "/123", true, map[string]string{"root": "123"}},
			{"GET", "/a/b/c", false, nil},
		}

		for _, tc := range testCases {
			req, _ := http.NewRequest(string(tc.method), tc.path, nil)
			ctx := NewContext(req, nil)
			_, ok := tree.Match(ctx)
			if ok != tc.match {
				t.Errorf("%s %s: got match=%v, want %v", tc.method, tc.path, ok, tc.match)
			}
			if ok {
				r := ctx.Request()

				for k := range tc.params {
					if r.PathValue(k) != tc.params[k] {
						t.Errorf("%s %s: param %s: got %q, want %q", tc.method, tc.path, k, r.PathValue(k), tc.params[k])
					}
				}
			}
		}
	})

	t.Run("param value capture", func(t *testing.T) {
		tree := newRadixTreeMatcher()
		tree.addEndpoint(&Endpoint{"GET", "/users/{id}/posts/{postId}", nil})

		req, _ := http.NewRequest("GET", "/users/42/posts/7", nil)
		ctx := NewContext(req, nil)
		_, ok := tree.Match(ctx)
		if !ok {
			t.Fatal("expected match")
		}

		r := ctx.Request()
		id := r.PathValue("id")
		postId := r.PathValue("postId")

		if id != "42" {
			t.Errorf("id: got %q, want %q", id, "42")
		}
		if postId != "7" {
			t.Errorf("postId: got %q, want %q", postId, "7")
		}
	})

	t.Run("wildcard routes", func(t *testing.T) {
		tree := newRadixTreeMatcher()
		tree.addEndpoint(&Endpoint{"GET", "/files/*/download", nil})
		tree.addEndpoint(&Endpoint{"GET", "/a/*/b/*/c", nil})

		testCases := []struct {
			method Method
			path   string
			match  bool
		}{
			{"GET", "/files/image.png/download", true},
			{"GET", "/files/doc.pdf/download", true},
			{"GET", "/files/anything/download", true},
			{"GET", "/files/download", false},
			{"GET", "/files/a/b/download", false},
			{"GET", "/a/x/b/y/c", true},
			{"GET", "/a/x/b/y/z", false},
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
		tree := newRadixTreeMatcher()
		tree.addEndpoint(&Endpoint{"GET", "/static/**", nil})
		tree.addEndpoint(&Endpoint{"GET", "/api/**/health", nil})
		tree.addEndpoint(&Endpoint{"GET", "/assets/**", nil})
		tree.addEndpoint(&Endpoint{"GET", "/x/**", nil})
		tree.addEndpoint(&Endpoint{"GET", "/xy", nil})

		testCases := []struct {
			method Method
			path   string
			match  bool
		}{
			{"GET", "/static", true},
			{"GET", "/static/a", true},
			{"GET", "/static/a/b", true},
			{"GET", "/static/a/b/c", true},
			{"GET", "/static/css/main.css", true},
			{"GET", "/static/js/app.js", true},
			{"GET", "/api/health", true},
			{"GET", "/api/health/health", true},
			{"GET", "/api/v1/health", true},
			{"GET", "/api/v1/v2/health", true},
			{"GET", "/api/v1/v2/v3/health", true},
			{"GET", "/assets", true},
			{"GET", "/assets/a/b/c", true},
			{"GET", "/x", true},
			{"GET", "/x/a/b", true},
			{"GET", "/xy", true},
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

	t.Run("double wildcard at root", func(t *testing.T) {
		tree := newRadixTreeMatcher()
		tree.addEndpoint(&Endpoint{"GET", "/**", nil})

		testCases := []struct {
			path  string
			match bool
		}{
			{"/", true},
			{"/a", true},
			{"/a/b", true},
			{"/a/b/c/d/e", true},
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

	t.Run("double wildcard in middle with suffix", func(t *testing.T) {
		tree := newRadixTreeMatcher()
		tree.addEndpoint(&Endpoint{"GET", "/a/**/b/c", nil})

		testCases := []struct {
			path  string
			match bool
		}{
			{"/a/b/c", true},
			{"/a/x/b/c", true},
			{"/a/x/y/b/c", true},
			{"/a/x/y/z/b/c", true},
			{"/a/b", false},
			{"/a/b/c/d", false},
			{"/a/b/b/c", true},
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

	t.Run("pattern routes", func(t *testing.T) {
		tree := newRadixTreeMatcher()
		tree.addEndpoint(&Endpoint{"GET", "/files/*.json", nil})
		tree.addEndpoint(&Endpoint{"GET", "/images/img?.png", nil})
		tree.addEndpoint(&Endpoint{"GET", "/docs/report-*.pdf", nil})
		tree.addEndpoint(&Endpoint{"GET", "/data/??-*.csv", nil})
		tree.addEndpoint(&Endpoint{"GET", "/logs/app**.log", nil})
		tree.addEndpoint(&Endpoint{"GET", "/files/*.json/meta", nil})
		tree.addEndpoint(&Endpoint{"GET", "/assets/img*", nil})

		testCases := []struct {
			method Method
			path   string
			match  bool
		}{
			{"GET", "/files/data.json", true},
			{"GET", "/files/config.json", true},
			{"GET", "/files/.json", true},
			{"GET", "/files/data.xml", false},
			{"GET", "/files/data.jsonx", false},
			{"GET", "/images/img1.png", true},
			{"GET", "/images/img2.png", true},
			{"GET", "/images/img12.png", false},
			{"GET", "/images/img.png", false},
			{"GET", "/docs/report-2024.pdf", true},
			{"GET", "/docs/report-.pdf", true},
			{"GET", "/docs/report.pdf", false},
			{"GET", "/data/ab-test.csv", true},
			{"GET", "/data/xy-data.csv", true},
			{"GET", "/data/a-test.csv", false},
			{"GET", "/data/abc-test.csv", false},
			{"GET", "/logs/app.log", true},
			{"GET", "/logs/app123.log", true},
			{"GET", "/logs/app.txt", false},
			{"GET", "/files/data.json/meta", true},
			{"GET", "/files/x.json/meta", true},
			{"GET", "/assets/img", true},
			{"GET", "/assets/img123", true},
			{"GET", "/assets/im", false},
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
		tree := newRadixTreeMatcher()
		tree.addEndpoint(&Endpoint{"GET", "/files/exact", nil})
		tree.addEndpoint(&Endpoint{"GET", "/files/*.json", nil})
		tree.addEndpoint(&Endpoint{"GET", "/files/{id}", nil})
		tree.addEndpoint(&Endpoint{"GET", "/files/*", nil})
		tree.addEndpoint(&Endpoint{"GET", "/files/**", nil})

		testCases := []struct {
			path        string
			wantPattern string
		}{
			{"/files/exact", "/files/exact"},
			{"/files/data.json", "/files/*.json"},
			{"/files/42", "/files/{id}"},
			{"/files/anything", "/files/{id}"},
			{"/files/a/b/c", "/files/**"},
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

	t.Run("param vs wildcard priority", func(t *testing.T) {
		tree := newRadixTreeMatcher()
		tree.addEndpoint(&Endpoint{"GET", "/items/{id}/detail", nil})
		tree.addEndpoint(&Endpoint{"GET", "/items/*/detail", nil})

		req, _ := http.NewRequest("GET", "/items/42/detail", nil)
		ctx := NewContext(req, nil)
		ep, ok := tree.Match(ctx)
		if !ok {
			t.Fatal("expected match")
		}
		if ep.path != "/items/{id}/detail" {
			t.Errorf("got %s, want /items/{id}/detail (param has higher priority)", ep.path)
		}
	})

	t.Run("multiple methods same path", func(t *testing.T) {
		tree := newRadixTreeMatcher()
		tree.addEndpoint(&Endpoint{"GET", "/resource", nil})
		tree.addEndpoint(&Endpoint{"POST", "/resource", nil})
		tree.addEndpoint(&Endpoint{"PUT", "/resource", nil})
		tree.addEndpoint(&Endpoint{"DELETE", "/resource", nil})
		tree.addEndpoint(&Endpoint{"PATCH", "/resource", nil})
		tree.addEndpoint(&Endpoint{"HEAD", "/resource", nil})
		tree.addEndpoint(&Endpoint{"OPTIONS", "/resource", nil})
		tree.addEndpoint(&Endpoint{"CONNECT", "/resource", nil})
		tree.addEndpoint(&Endpoint{"TRACE", "/resource", nil})
		tree.addEndpoint(&Endpoint{"", "/resource", nil})

		testCases := []struct {
			method Method
			match  bool
		}{
			{"GET", true},
			{"POST", true},
			{"PUT", true},
			{"DELETE", true},
			{"PATCH", true},
			{"HEAD", true},
			{"OPTIONS", true},
			{"CONNECT", true},
			{"TRACE", true},
		}

		for _, tc := range testCases {
			req, _ := http.NewRequest(string(tc.method), "/resource", nil)
			ctx := NewContext(req, nil)
			ep, ok := tree.Match(ctx)
			if ok != tc.match {
				t.Errorf("%s /resource: got match=%v, want %v", tc.method, ok, tc.match)
				continue
			}
			if ok && ep.method != tc.method {
				t.Errorf("%s /resource: got method=%s", tc.method, ep.method)
			}
		}
	})

	t.Run("mixed route types in same tree", func(t *testing.T) {
		tree := newRadixTreeMatcher()
		tree.addEndpoint(&Endpoint{"GET", "/api/users", nil})
		tree.addEndpoint(&Endpoint{"GET", "/api/users/{id}", nil})
		tree.addEndpoint(&Endpoint{"GET", "/api/users/{id}/roles/*", nil})
		tree.addEndpoint(&Endpoint{"GET", "/api/**/health", nil})
		tree.addEndpoint(&Endpoint{"GET", "/api/files/*.json", nil})

		testCases := []struct {
			path        string
			match       bool
			wantPattern string
		}{
			{"/api/users", true, "/api/users"},
			{"/api/users/42", true, "/api/users/{id}"},
			{"/api/users/42/roles/admin", true, "/api/users/{id}/roles/*"},
			{"/api/health", true, "/api/**/health"},
			{"/api/v1/v2/health", true, "/api/**/health"},
			{"/api/files/data.json", true, "/api/files/*.json"},
			{"/api/files/data.xml", false, ""},
			{"/api/users/42/roles", false, ""},
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

	t.Run("param to wildcard backtracking", func(t *testing.T) {
		tree := newRadixTreeMatcher()
		tree.addEndpoint(&Endpoint{"GET", "/files/{id}/detail", nil})
		tree.addEndpoint(&Endpoint{"GET", "/files/*/other", nil})

		testCases := []struct {
			path  string
			match bool
		}{
			{"/files/test/detail", true},
			{"/files/test/other", true},
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

	t.Run("pattern child backtracking", func(t *testing.T) {
		tree := newRadixTreeMatcher()
		tree.addEndpoint(&Endpoint{"GET", "/files/*.json/download", nil})
		tree.addEndpoint(&Endpoint{"GET", "/files/*.txt/upload", nil})

		testCases := []struct {
			path  string
			match bool
		}{
			{"/files/data.json/download", true},
			{"/files/data.txt/upload", true},
			{"/files/data.txt/download", false},
			{"/files/data.json/upload", false},
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

	t.Run("trailing slash normalization", func(t *testing.T) {
		tree := newRadixTreeMatcher()
		tree.addEndpoint(&Endpoint{"GET", "/users", nil})
		tree.addEndpoint(&Endpoint{"GET", "/a/b/c", nil})

		testCases := []struct {
			path  string
			match bool
		}{
			{"/users/", true},
			{"/a/b/c/", true},
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

	t.Run("path normalization on add", func(t *testing.T) {
		tree := newRadixTreeMatcher()
		tree.addEndpoint(&Endpoint{"GET", "users", nil})
		tree.addEndpoint(&Endpoint{"GET", "/items/", nil})

		testCases := []struct {
			path  string
			match bool
		}{
			{"/users", true},
			{"/items", true},
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

	t.Run("duplicate route returns error", func(t *testing.T) {
		tree := newRadixTreeMatcher()
		err := tree.addEndpoint(&Endpoint{"GET", "/users", nil})
		if err != nil {
			t.Fatalf("first add failed: %v", err)
		}
		err = tree.addEndpoint(&Endpoint{"GET", "/users", nil})
		if err == nil {
			t.Error("duplicate route should return error")
		}
	})

	t.Run("buildRouteEntry errors", func(t *testing.T) {
		testCases := []struct {
			name string
			path string
		}{
			{"unclosed param", "/users/{id"},
			{"empty param name", "/users/{}"},
			{"too many params", "/a/{p1}/{p2}/{p3}/{p4}/{p5}/{p6}/{p7}/{p8}/{p9}/{p10}/{p11}/{p12}/{p13}/{p14}/{p15}/{p16}/{p17}"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				tree := newRadixTreeMatcher()
				err := tree.addEndpoint(&Endpoint{"GET", tc.path, nil})
				if err == nil {
					t.Errorf("expected error for path %q", tc.path)
				}
			})
		}
	})

}

func BenchmarkRadixTreeMatcher_Match(b *testing.B) {
	tree := newRadixTreeMatcher()
	tree.addEndpoint(&Endpoint{"GET", "/users", nil})
	tree.addEndpoint(&Endpoint{"GET", "/users/{id}", nil})
	tree.addEndpoint(&Endpoint{"GET", "/users/{id}/posts/{postId}", nil})
	tree.addEndpoint(&Endpoint{"GET", "/files/*/download", nil})
	tree.addEndpoint(&Endpoint{"GET", "/static/**", nil})
	tree.addEndpoint(&Endpoint{"GET", "/api/**/health", nil})
	tree.addEndpoint(&Endpoint{"GET", "/files/*.json", nil})
	tree.addEndpoint(&Endpoint{"POST", "/users", nil})
	tree.addEndpoint(&Endpoint{"PUT", "/users/{id}", nil})
	tree.addEndpoint(&Endpoint{"DELETE", "/users/{id}", nil})

	benchCases := []struct {
		name   string
		method string
		path   string
	}{
		{"static", "GET", "/users"},
		{"param_1", "GET", "/users/42"},
		{"param_2", "GET", "/users/42/posts/7"},
		{"wildcard", "GET", "/files/image.png/download"},
		{"double_wildcard_short", "GET", "/static/css/main.css"},
		{"double_wildcard_deep", "GET", "/static/a/b/c/d/e/f"},
		{"double_wildcard_middle", "GET", "/api/v1/v2/health"},
		{"pattern", "GET", "/files/data.json"},
		{"not_found", "GET", "/notfound/path"},
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
