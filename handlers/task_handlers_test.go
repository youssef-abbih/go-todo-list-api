package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/youssef-abbih/go-todo-list/middleware"
	"github.com/youssef-abbih/go-todo-list/models"
)

func setParam(r *http.Request, key, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

func addAuth(req *http.Request) *http.Request {
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, fmt.Sprintf("%d", testUserID))
	return req.WithContext(ctx)
}

func setup() {
	models.InitDB()
	setupTestUser()
}

// Test POST /tasks
func TestPostTask(t *testing.T) {
	setup()

	validTask := models.Task{Title: "Test", Description: "Test desc", Completed: false}
	body, _ := json.Marshal(validTask)
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
	req.Header.Set(ContentTypeHeader, MimeJSON)
	req = addAuth(req)
	rec := httptest.NewRecorder()
	PostTask(rec, req)
	if rec.Code != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d", rec.Code)
	}

	malformedReq := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader("invalid json"))
	malformedReq.Header.Set(ContentTypeHeader, MimeJSON)
	malformedReq = addAuth(malformedReq)
	malformedRec := httptest.NewRecorder()
	PostTask(malformedRec, malformedReq)
	if malformedRec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request for malformed JSON, got %d", malformedRec.Code)
	}
}

// Test GET /tasks
func TestGetTasks(t *testing.T) {
	setup()
	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	req = addAuth(req)
	rec := httptest.NewRecorder()
	GetTasks(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", rec.Code)
	}
}

// Test GET /tasks/{id}
func TestGetTask(t *testing.T) {
	setup()

	task := models.AddTask(models.Task{Title: "Test", Description: "desc", Completed: false}, testUserID)

	req := httptest.NewRequest(http.MethodGet, "/tasks/"+fmt.Sprintf("%d", task.ID), nil)
	req = setParam(req, "id", fmt.Sprintf("%d", task.ID))
	req = addAuth(req)
	rec := httptest.NewRecorder()
	GetTask(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", rec.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/tasks/9999", nil)
	req = setParam(req, "id", "9999")
	req = addAuth(req)
	rec = httptest.NewRecorder()
	GetTask(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404 Not Found, got %d", rec.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/tasks/abc", nil)
	req = setParam(req, "id", "abc")
	req = addAuth(req)
	rec = httptest.NewRecorder()
	GetTask(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d", rec.Code)
	}
}

// Test PUT /tasks/{id}
func TestPutTask(t *testing.T) {
	setup()

	task := models.AddTask(models.Task{Title: "Test", Description: "desc", Completed: false}, testUserID)

	updated := models.Task{Title: "Updated", Description: "Updated desc", Completed: true}
	body, _ := json.Marshal(updated)

	req := httptest.NewRequest(http.MethodPut, "/tasks/"+fmt.Sprintf("%d", task.ID), bytes.NewReader(body))
	req = setParam(req, "id", fmt.Sprintf("%d", task.ID))
	req.Header.Set(ContentTypeHeader, MimeJSON)
	req = addAuth(req)
	rec := httptest.NewRecorder()
	PutTask(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", rec.Code)
	}

	body, _ = json.Marshal(updated)
	nonexistent := httptest.NewRequest(http.MethodPut, "/tasks/9999", bytes.NewReader(body))
	nonexistent = setParam(nonexistent, "id", "9999")
	nonexistent.Header.Set(ContentTypeHeader, MimeJSON)
	nonexistent = addAuth(nonexistent)
	rec = httptest.NewRecorder()
	PutTask(rec, nonexistent)
	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404 Not Found, got %d", rec.Code)
	}

	malformed := httptest.NewRequest(http.MethodPut, "/tasks/"+fmt.Sprintf("%d", task.ID), strings.NewReader("bad json"))
	malformed = setParam(malformed, "id", fmt.Sprintf("%d", task.ID))
	malformed.Header.Set(ContentTypeHeader, MimeJSON)
	malformed = addAuth(malformed)
	rec = httptest.NewRecorder()
	PutTask(rec, malformed)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d", rec.Code)
	}
}

// Test DELETE /tasks/{id}
func TestDeleteTask(t *testing.T) {
	setup()

	task := models.AddTask(models.Task{Title: "Test", Description: "desc", Completed: false}, testUserID)

	req := httptest.NewRequest(http.MethodDelete, "/tasks/"+fmt.Sprintf("%d", task.ID), nil)
	req = setParam(req, "id", fmt.Sprintf("%d", task.ID))
	req = addAuth(req)
	rec := httptest.NewRecorder()
	DeleteTask(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", rec.Code)
	}

	nonexistent := httptest.NewRequest(http.MethodDelete, "/tasks/9999", nil)
	nonexistent = setParam(nonexistent, "id", "9999")
	nonexistent = addAuth(nonexistent)
	rec = httptest.NewRecorder()
	DeleteTask(rec, nonexistent)
	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404 Not Found, got %d", rec.Code)
	}

	invalid := httptest.NewRequest(http.MethodDelete, "/tasks/abc", nil)
	invalid = setParam(invalid, "id", "abc")
	invalid = addAuth(invalid)
	rec = httptest.NewRecorder()
	DeleteTask(rec, invalid)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d", rec.Code)
	}
}
