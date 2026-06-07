package models

import (
	"os"
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "testpass123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if hash == "" || hash == password {
		t.Errorf("password hashing failed to transform string")
	}
}

func TestInitDB_Environments(t *testing.T) {
	// Keep a backup of the original environment variables
	oldEnv := os.Getenv("ENV")
	oldDevUser := os.Getenv("DEV_DB_USER")
	oldProdUser := os.Getenv("PROD_DB_USER")

	// Restore original variables when this test finishes
	defer func() {
		os.Setenv("ENV", oldEnv)
		os.Setenv("DEV_DB_USER", oldDevUser)
		os.Setenv("PROD_DB_USER", oldProdUser)
	}()

	// Mock minimal strings to run through the branches safely
	os.Setenv("DEV_DB_USER", "user")
	os.Setenv("PROD_DB_USER", "user")

	envs := []string{"DEV", "PROD"}

	for _, env := range envs {
		t.Run(env, func(t *testing.T) {
			os.Setenv("ENV", env)
			SeedTestData(nil) // Hits the branch logic safely with nil
		})
	}
}