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
		tree.addEndpoint(&Endpoint{"/users", "GET", nil})
		tree.addEndpoint(&Endpoint{"/users/profile", "GET", nil})
		tree.addEndpoint(&Endpoint{"/users", "POST", nil})

		testCases := []struct {
			method Method
			path   string
			match  bool
		}{
			{"GET", "/users", true},
			{"GET", "/users/profile", true},
			{"POST", "/users", true},
			{"DELETE", "/users", false},
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
		tree := newRadixTreeMatcher()
		tree.addEndpoint(&Endpoint{"GET", "/users/{id}", nil})
		tree.addEndpoint(&Endpoint{"GET", "/users/{id}/posts/{postId}", nil})

		testCases := []struct {
			method     Method
			path       string
			match      bool
			paramCount int
		}{
			{"GET", "/users/42", true, 1},
			{"GET", "/users/42/posts/7", true, 2},
			{"GET", "/users", false, 0},
		}

		for _, tc := range testCases {
			req, _ := http.NewRequest(string(tc.method), tc.path, nil)
			ctx := NewContext(req, nil)
			_, ok := tree.Match(ctx)
			if ok != tc.match {
				t.Errorf("%s %s: got match=%v, want %v", tc.method, tc.path, ok, tc.match)
			}
			if ok {
				req := ctx.Request()
				if req.pathValues.count != tc.paramCount {
					t.Errorf("%s %s: got paramCount=%d, want %d",
						tc.method, tc.path, req.pathValues.count, tc.paramCount)
				}
			}
		}
	})

	t.Run("wildcard routes", func(t *testing.T) {
		tree := newRadixTreeMatcher()
		tree.addEndpoint(&Endpoint{"GET", "/files/*/download", nil})

		testCases := []struct {
			method Method
			path   string
			match  bool
		}{
			{"GET", "/files/image.png/download", true},
			{"GET", "/files/doc.pdf/download", true},
			{"GET", "/files/download", false},     // * boş segment eşlemez
			{"GET", "/files/a/b/download", false}, // * tek segment
		}

		for _, testCase := range testCases {
			req, _ := http.NewRequest(string(testCase.method), testCase.path, nil)
			ctx := NewContext(req, nil)
			_, ok := tree.Match(ctx)
			if ok != testCase.match {
				t.Errorf("%s %s: got match=%v, want %v", testCase.method, testCase.path, ok, testCase.match)
			}
		}
	})

	t.Run("double wildcard routes", func(t *testing.T) {
		tree := newRadixTreeMatcher()
		tree.addEndpoint(&Endpoint{"GET", "/static/**", nil})
		tree.addEndpoint(&Endpoint{"GET", "/api/**/health", nil})

		tests := []struct {
			method Method
			path   string
			match  bool
		}{
			{"GET", "/static", true},              // ** zero-match
			{"GET", "/static/css/main.css", true}, // ** çoklu segment
			{"GET", "/static/js/app.js", true},
			{"GET", "/api/health", true},       // ** zero-match
			{"GET", "/api/v1/health", true},    // ** tek segment
			{"GET", "/api/v1/v2/health", true}, // ** çoklu segment
		}

		for _, tc := range tests {
			req, _ := http.NewRequest(string(tc.method), tc.path, nil)
			ctx := NewContext(req, nil)
			_, ok := tree.Match(ctx)
			if ok != tc.match {
				t.Errorf("%s %s: got match=%v, want %v", tc.method, tc.path, ok, tc.match)
			}
		}
	})

	t.Run("pattern routes", func(t *testing.T) {
		tree := newRadixTreeMatcher()
		tree.addEndpoint(&Endpoint{"GET", "/files/*.json", nil})
		tree.addEndpoint(&Endpoint{"GET", "/images/img?.png", nil})

		tests := []struct {
			method Method
			path   string
			match  bool
		}{
			{"GET", "/files/data.json", true},
			{"GET", "/files/config.json", true},
			{"GET", "/files/data.xml", false},
			{"GET", "/images/img1.png", true},
			{"GET", "/images/img2.png", true},
			{"GET", "/images/img12.png", false}, // ? tek karakter
		}

		for _, tc := range tests {
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

		tests := []struct {
			path        string
			wantPattern string
		}{
			{"/files/exact", "/files/exact"},
			{"/files/data.json", "/files/*.json"},
			{"/files/42", "/files/{id}"},
			{"/files/anything", "/files/{id}"}, // param önce gelir
			{"/files/a/b/c", "/files/**"},      // ** çoklu segment
		}

		for _, tc := range tests {
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

	t.Run("trailing slash normalization", func(t *testing.T) {
		tree := newRadixTreeMatcher()
		tree.addEndpoint(&Endpoint{"GET", "/users", nil})

		req, _ := http.NewRequest("GET", "/users/", nil)
		ctx := NewContext(req, nil)
		_, ok := tree.Match(ctx)
		if !ok {
			t.Error("GET /users/ should match /users (trailing slash stripped)")
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

}
