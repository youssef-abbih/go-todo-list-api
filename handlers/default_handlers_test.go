package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	//"github.com/youssef-abbih/go-todo-list/models"
)

// Test DefaultResponse handler
func TestDefaultResponse(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	DefaultResponse(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	expected := `{"message":"Welcome to my Todo List API"}`
	if strings.TrimSpace(string(body)) != expected {
		t.Errorf("expected body %q, got %q", expected, string(body))
	}
}

func TestHealthCheck(t *testing.T) {
    // Test wrong method first — no DB needed
    req := httptest.NewRequest(http.MethodPost, "/health", nil)
    rec := httptest.NewRecorder()
    HealthCheck(rec, req)
    if rec.Code != http.StatusMethodNotAllowed {
        t.Errorf("expected 405 Method Not Allowed, got %d", rec.Code)
    }
	setup()
    

    req= httptest.NewRequest(http.MethodGet, "/health", nil)
    rec = httptest.NewRecorder()
    HealthCheck(rec, req)
    if rec.Code != http.StatusOK && rec.Code != http.StatusInternalServerError {
        t.Errorf("unexpected status %d", rec.Code)
    }
}