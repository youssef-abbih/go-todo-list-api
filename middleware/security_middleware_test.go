package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecureHeadersMiddleware(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := SecureHeadersMiddleware(next)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	headers := map[string]string{
		"X-XSS-Protection":       "1; mode-block",
		"X-Frame-Options":         "deny",
		"Content-Security-Policy": "default-src 'self'",
		"Referrer-Policy":         "no-referrer",
	}

	for header, expected := range headers {
		if rec.Header().Get(header) != expected {
			t.Errorf("expected %s: %s, got %s", header, expected, rec.Header().Get(header))
		}
	}
}