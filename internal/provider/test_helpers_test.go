package provider

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

type testRoute struct {
	Method  string
	Path    string
	Status  int
	Body    string
	Handler func(http.ResponseWriter, *http.Request)
}

func newTestServer(t *testing.T, routes []testRoute) *httptest.Server {
	t.Helper()

	routeMap := make(map[string]testRoute, len(routes))
	for _, route := range routes {
		key := route.Method + " " + route.Path
		routeMap[key] = route
	}

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if route, ok := routeMap[r.Method+" "+r.URL.Path]; ok {
			if route.Handler != nil {
				route.Handler(w, r)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			status := route.Status
			if status == 0 {
				status = http.StatusOK
			}
			w.WriteHeader(status)
			if route.Body != "" {
				_, _ = w.Write([]byte(route.Body))
			}
			return
		}

		http.Error(w, "Not Found", http.StatusNotFound)
	}))
}

func testAccProviderConfig(serverURL string) string {
	return "provider \"elevenlabs\" {\n" +
		"  api_key  = \"test-key\"\n" +
		"  base_url = \"" + serverURL + "\"\n" +
		"}\n"
}

func writeTempFile(t *testing.T, name string, contents []byte) string {
	t.Helper()

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, name)
	if err := os.WriteFile(path, contents, 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return path
}
