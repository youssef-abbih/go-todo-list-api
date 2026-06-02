package handlers
import (
	//"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/youssef-abbih/go-todo-list/utils"
	"github.com/youssef-abbih/go-todo-list/models"
)

var testUserID uint
var testToken string

func setupTestUser() {
    user := models.User{Email: fmt.Sprintf("test_%d@test.com", time.Now().UnixNano()), Password: "password"}
    user.Password, _ = models.HashPassword(user.Password)
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

    // Test wrong method
    req := httptest.NewRequest(http.MethodGet, "/users/register", nil)
    rec := httptest.NewRecorder()
    Register(rec, req)
    if rec.Code != http.StatusMethodNotAllowed {
        t.Errorf("expected 405, got %d", rec.Code)
    }

	body_missing_email := `{"email": "", "password": "secret"}`
	req = httptest.NewRequest(http.MethodPost, "/users/register", strings.NewReader(body_missing_email))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	Register(rec, req)
	if rec.Code != 400 {
		t.Errorf("expected 400, got %d", rec.Code)
	}

	body_missing_password := `{"email": "user@email.com", "password": ""}`
	req = httptest.NewRequest(http.MethodPost, "/users/register", strings.NewReader(body_missing_password))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	Register(rec, req)
	if rec.Code != 400 {
		t.Errorf("expected 400, got %d", rec.Code)
	}

	body_duplicate_user := `{"email": "test@test.com", "password": "password"}`
	req = httptest.NewRequest(http.MethodPost, "/users/register", strings.NewReader(body_duplicate_user))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	Register(rec, req)
	if rec.Code != 409 {
		t.Errorf("expected 400, got %d", rec.Code)
	}

}