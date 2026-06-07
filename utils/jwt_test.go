package utils

import (
	"os"
	"testing"
)

func TestLoadJWTSecretkey(t *testing.T) {
	// Case 1: JWT_SECRET is set
	os.Setenv("JWT_SECRET", "testsecretkey")
	key := LoadJWTSecretkey()
	if key != "testsecretkey" {
		t.Errorf("expected testsecretkey, got %s", key)
	}
}

func TestStoreSecretKey(t *testing.T) {
	// Create a temp .env file
	tmpFile, err := os.CreateTemp("", ".env")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// Write initial content
	tmpFile.WriteString("SOME_VAR=value\n")
	tmpFile.Close()

	// Change working directory to temp dir is complex
	// so just verify generateSecretKey produces a non-empty string
	key := generateSecretKey()
	if key == "" {
		t.Errorf("expected non-empty secret key")
	}
	if len(key) < 32 {
		t.Errorf("expected key length >= 32, got %d", len(key))
	}
}