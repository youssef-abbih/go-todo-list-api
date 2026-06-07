package utils

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/youssef-abbih/go-todo-list/middleware"
)

func TestGetUserID(t *testing.T) {
	// Case 1: No user ID in context
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	_, err := GetUserID(req)
	if err == nil {
		t.Errorf("expected error when no user ID in context")
	}

	// Case 2: Wrong type in context (not a string)
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, 42)
	req = req.WithContext(ctx)
	_, err = GetUserID(req)
	if err == nil {
		t.Errorf("expected error when user ID is not a string")
	}

	// Case 3: Invalid string (not a number)
	ctx = context.WithValue(req.Context(), middleware.UserContextKey, "abc")
	req = req.WithContext(ctx)
	_, err = GetUserID(req)
	if err == nil {
		t.Errorf("expected error when user ID is not a valid number")
	}

	// Case 4: Valid user ID
	ctx = context.WithValue(req.Context(), middleware.UserContextKey, "42")
	req = req.WithContext(ctx)
	id, err := GetUserID(req)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if id != 42 {
		t.Errorf("expected 42, got %d", id)
	}
}
