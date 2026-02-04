package http

import (
	"fmt"
	"testing"
)

func TestBasicRouting(t *testing.T) {
	tree := NewRadixTree()

	tree.AddRoute("/api/users", "GET", "list users")
	tree.AddRoute("/api/users/{id}", "GET", "get user")
	tree.AddRoute("/api/users/{id}", "DELETE", "delete user")
	tree.AddRoute("/api/users/{email}/status", "GET", "get status by email")
	tree.AddRoute("/api/posts/{postId}/comments/{commentId}", "GET", "get comment")
	tree.AddRoute("/static/**", "GET", "static files")

	tests := []struct {
		path       string
		method     Method
		wantFound  bool
		wantParams map[string]string
	}{
		{"/api/users", "GET", true, nil},
		{"/api/users/123", "GET", true, map[string]string{"id": "123"}},
		{"/api/users/123", "DELETE", true, map[string]string{"id": "123"}},
		{"/api/users/test@example.com/status", "GET", true, map[string]string{"email": "test@example.com"}},
		{"/api/posts/42/comments/99", "GET", true, map[string]string{"postId": "42", "commentId": "99"}},
		{"/static/css/style.css", "GET", true, map[string]string{"**": "css/style.css"}},
		{"/api/users/123", "POST", false, nil},
		{"/notfound", "GET", false, nil},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s %s", tt.method, tt.path), func(t *testing.T) {
			result := tree.Match(tt.path, tt.method)

			if result.Found != tt.wantFound {
				t.Errorf("found = %v, want %v", result.Found, tt.wantFound)
				return
			}

			if tt.wantFound && tt.wantParams != nil {
				for key, want := range tt.wantParams {
					got := result.Params.Get(key)
					if got != want {
						t.Errorf("param %q = %q, want %q", key, got, want)
					}
				}
			}
		})
	}
}

func TestEndpointStaysClean(t *testing.T) {
	tree := NewRadixTree()

	type MyHandler struct {
		Name string
	}

	handler := &MyHandler{Name: "test"}
	tree.AddRoute("/api/{id}", "GET", handler)

	result := tree.Match("/api/123", "GET")
	if !result.Found {
		t.Fatal("should match")
	}

}

func TestDifferentParamNames(t *testing.T) {
	tree := NewRadixTree()

	tree.AddRoute("/api/users/{id}", "GET", "by id")
	tree.AddRoute("/api/users/{email}/profile", "GET", "profile by email")

	r1 := tree.Match("/api/users/123", "GET")
	if r1.Params.Get("id") != "123" {
		t.Errorf("expected id=123, got %v", r1.Params)
	}

	r2 := tree.Match("/api/users/test@example.com/profile", "GET")
	if r2.Params.Get("email") != "test@example.com" {
		t.Errorf("expected email=test@example.com, got %v", r2.Params)
	}
}

func TestStaticPriority(t *testing.T) {
	tree := NewRadixTree()

	tree.AddRoute("/api/users/me", "GET", "current user")
	tree.AddRoute("/api/users/{id}", "GET", "user by id")

	result := tree.Match("/api/users/me", "GET")
	if !result.Found {
		t.Fatal("should match")
	}

	result = tree.Match("/api/users/123", "GET")
	if result.Params.Get("id") != "123" {
		t.Error("param should still work")
	}
}

func TestDuplicateRoute(t *testing.T) {
	tree := NewRadixTree()

	if err := tree.AddRoute("/api/users", "GET", "h1"); err != nil {
		t.Fatalf("first add failed: %v", err)
	}

	if err := tree.AddRoute("/api/users", "GET", "h2"); err != ErrDuplicateRoute {
		t.Errorf("expected ErrDuplicateRoute, got %v", err)
	}
}

func TestTrailingSlash(t *testing.T) {
	tree := NewRadixTree()
	tree.AddRoute("/api/users/", "GET", "handler")

	if r := tree.Match("/api/users", "GET"); !r.Found {
		t.Error("/api/users should match")
	}
	if r := tree.Match("/api/users/", "GET"); !r.Found {
		t.Error("/api/users/ should match")
	}
}

func TestInvalidPatterns(t *testing.T) {
	tree := NewRadixTree()

	if err := tree.AddRoute("/api/{}", "GET", "h"); err != ErrEmptyParamName {
		t.Errorf("expected ErrEmptyParamName, got %v", err)
	}

	if err := tree.AddRoute("/api/{unclosed", "GET", "h"); err != ErrUnclosedParam {
		t.Errorf("expected ErrUnclosedParam, got %v", err)
	}
}

func TestCatchAllWithPath(t *testing.T) {
	tree := NewRadixTree()
	tree.AddRoute("/files/**", "GET", "files")

	tests := []struct {
		path string
		want string
	}{
		{"/files/a", "a"},
		{"/files/a/b", "a/b"},
		{"/files/a/b/c.txt", "a/b/c.txt"},
	}

	for _, tt := range tests {
		result := tree.Match(tt.path, "GET")
		if !result.Found {
			t.Errorf("%s should match", tt.path)
			continue
		}
		got := result.Params.Get("**")
		if got != tt.want {
			t.Errorf("path=%s: got %q, want %q", tt.path, got, tt.want)
		}
	}
}

func BenchmarkMatch(b *testing.B) {
	tree := NewRadixTree()

	routes := []string{
		"/",
		"/api/v1/users",
		"/api/v1/users/{id}",
		"/api/v1/users/{id}/posts/{postId}/comments/{commentId}",
		"/static/**",
	}
	for _, r := range routes {
		tree.AddRoute(r, "GET", "handler")
	}

	paths := []string{
		"/api/v1/users",
		"/api/v1/users/123",
		"/api/v1/users/123/posts/456/comments/789",
		"/static/css/app.css",
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		tree.Match(paths[i%len(paths)], "GET")
	}
}

func BenchmarkMatchStatic(b *testing.B) {
	tree := NewRadixTree()
	tree.AddRoute("/api/v1/users", "GET", "handler")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		tree.Match("/api/v1/users", "GET")
	}
}

func BenchmarkMatchParam(b *testing.B) {
	tree := NewRadixTree()
	tree.AddRoute("/api/v1/users/{id}", "GET", "handler")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		tree.Match("/api/v1/users/123", "GET")
	}
}

func BenchmarkMatchDeepParam(b *testing.B) {
	tree := NewRadixTree()
	tree.AddRoute("/api/v1/users/{userId}/posts/{postId}/comments/{commentId}", "GET", "handler")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		tree.Match("/api/v1/users/123/posts/456/comments/789", "GET")
	}
}

func BenchmarkMatchCatchAll(b *testing.B) {
	tree := NewRadixTree()
	tree.AddRoute("/static/**", "GET", "handler")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		tree.Match("/static/css/app/style.css", "GET")
	}
}

// Parallel benchmark
func BenchmarkMatchParallel(b *testing.B) {
	tree := NewRadixTree()
	tree.AddRoute("/api/v1/users/{id}", "GET", "handler")

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			tree.Match("/api/v1/users/123", "GET")
		}
	})
}
