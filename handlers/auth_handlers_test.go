package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/youssef-abbih/go-todo-list/models"
	"github.com/youssef-abbih/go-todo-list/utils"
)

var testUserID uint
var testToken string
var testEmail string

func setupTestUser() {
	testEmail = fmt.Sprintf("test_%d@test.com", time.Now().UnixNano())
	user := models.User{Email: testEmail, Password: "password"}
	created, _ := models.AddUser(user)
	testUserID = created.ID

	claims := jwt.MapClaims{
		"sub":   created.ID,
		"email": created.Email,
		"exp":   time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	testToken, _ = token.SignedString([]byte(utils.LoadJWTSecretkey()))
}

func TestRegister(t *testing.T) {
	setup()
	setupTestUser()

	// Wrong method
	req := httptest.NewRequest(http.MethodGet, "/users/register", nil)
	rec := httptest.NewRecorder()
	Register(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}

	// Missing email
	req = httptest.NewRequest(http.MethodPost, "/users/register", strings.NewReader(`{"email": "", "password": "secret"}`))
	req.Header.Set(ContentTypeHeader, MimeJSON)
	rec = httptest.NewRecorder()
	Register(rec, req)
	if rec.Code != 400 {
		t.Errorf("expected 400, got %d", rec.Code)
	}

	// Missing password
	req = httptest.NewRequest(http.MethodPost, "/users/register", strings.NewReader(`{"email": "user@email.com", "password": ""}`))
	req.Header.Set(ContentTypeHeader, MimeJSON)
	rec = httptest.NewRecorder()
	Register(rec, req)
	if rec.Code != 400 {
		t.Errorf("expected 400, got %d", rec.Code)
	}

	// Duplicate user
	req = httptest.NewRequest(http.MethodPost, "/users/register", strings.NewReader(`{"email": "`+testEmail+`", "password": "password"}`))
	req.Header.Set(ContentTypeHeader, MimeJSON)
	rec = httptest.NewRecorder()
	Register(rec, req)
	if rec.Code != 409 {
		t.Errorf("expected 409, got %d", rec.Code)
	}

	// Valid registration
	newEmail := fmt.Sprintf("new_%d@test.com", time.Now().UnixNano())
	req = httptest.NewRequest(http.MethodPost, "/users/register", strings.NewReader(`{"email": "`+newEmail+`", "password": "secret123"}`))
	req.Header.Set(ContentTypeHeader, MimeJSON)
	rec = httptest.NewRecorder()
	Register(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestLogin(t *testing.T) {
	setup()
	setupTestUser()

	// Wrong method
	req := httptest.NewRequest(http.MethodGet, "/users/login", nil)
	rec := httptest.NewRecorder()
	Login(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}

	// Invalid JSON
	req = httptest.NewRequest(http.MethodPost, "/users/login", strings.NewReader("bad json"))
	req.Header.Set(ContentTypeHeader, MimeJSON)
	rec = httptest.NewRecorder()
	Login(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}

	// Email not found
	req = httptest.NewRequest(http.MethodPost, "/users/login", strings.NewReader(`{"email": "notfound@test.com", "password": "password"}`))
	req.Header.Set(ContentTypeHeader, MimeJSON)
	rec = httptest.NewRecorder()
	Login(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}

	// Wrong password
	req = httptest.NewRequest(http.MethodPost, "/users/login", strings.NewReader(`{"email": "`+testEmail+`", "password": "wrongpassword"}`))
	req.Header.Set(ContentTypeHeader, MimeJSON)
	rec = httptest.NewRecorder()
	Login(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}

	// Valid login
	req = httptest.NewRequest(http.MethodPost, "/users/login", strings.NewReader(`{"email": "`+testEmail+`", "password": "password"}`))
	req.Header.Set(ContentTypeHeader, MimeJSON)
	rec = httptest.NewRecorder()
	Login(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	// Check token in response
	var resp map[string]string
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp["Token"] == "" {
		t.Errorf("expected token in response, got empty")
	}
}
