package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func generateTestToken(secret string, expired bool) string {
	exp := time.Now().Add(time.Hour)
	if expired {
		exp = time.Now().Add(-time.Hour)
	}
	claims := jwt.MapClaims{
		"user_id": "42",
		"exp":     exp.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(secret))
	return signed
}

func TestAuthMiddleware(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := AuthMiddleware(next)

	// Define the test table structural schema
	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
	}{
		{
			name:           "Case 1: Missing Authorization header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Case 2: Wrong format",
			authHeader:     "Token abc",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Case 3: Invalid token",
			authHeader:     "Bearer invalidtoken",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Case 4: Expired token",
			authHeader:     "Bearer " + generateTestToken("your-secret-key", true),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Case 5: Valid token",
			authHeader:     "Bearer " + generateTestToken("your-secret-key", false),
			expectedStatus: http.StatusOK,
		},
	}

	// Loop over each test case dynamically
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}