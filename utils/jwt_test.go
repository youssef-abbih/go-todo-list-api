package utils

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestLoadJWTSecretkey_Success(t *testing.T) {
	// Clean environment context
	os.Unsetenv(jwtSecretKey)
	defer os.Unsetenv(jwtSecretKey)

	// Case 1: Fetching when environment variable is directly set
	os.Setenv(jwtSecretKey, "direct-environment-secret-value")
	key := LoadJWTSecretkey()
	if key != "direct-environment-secret-value" {
		t.Errorf("expected direct-environment-secret-value, got %s", key)
	}

	// Case 2: Falling back to loading from file
	os.Unsetenv(jwtSecretKey)
	tmpFile, err := os.CreateTemp("", "test_env_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	_, _ = tmpFile.WriteString(jwtSecretKey + "=file-loaded-secret-value\n")
	tmpFile.Close()

	// Divert logic to read from temporary test mock file
	oldPath := envFilePath
	envFilePath = tmpFile.Name()
	defer func() { envFilePath = oldPath }()

	keyFromEnvFile := LoadJWTSecretkey()
	if keyFromEnvFile != "file-loaded-secret-value" {
		t.Errorf("expected file-loaded-secret-value, got %s", keyFromEnvFile)
	}
}

// TestLoadJWTSecretkey_Crash verifies that the function terminates on empty strings safely
func TestLoadJWTSecretkey_Crash(t *testing.T) {
	if os.Getenv("BE_CRASH_TEST") == "1" {
		os.Unsetenv(jwtSecretKey)
		oldPath := envFilePath
		envFilePath = "non_existent_file_to_force_failure.env"
		defer func() { envFilePath = oldPath }()
		
		LoadJWTSecretkey() // This should invoke log.Fatal
		return
	}

	// Execute this test as a isolated subprocess check
	cmd := exec.Command(os.Args[0], "-test.run=TestLoadJWTSecretkey_Crash")
	cmd.Env = append(os.Environ(), "BE_CRASH_TEST=1")
	err := cmd.Run()

	// Assert that it cleanly exited with an error exit code status
	if err == nil {
		t.Fatalf("expected process to crash with exit status via log.Fatal, but it exited successfully")
	}
}

func TestStoreSecretKey(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_store_env_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	
	// Pre-populate mock configuration content
	_, _ = tmpFile.WriteString("EXISTING_VAR=keep_me\n" + jwtSecretKey + "=old_token_to_replace\n")
	tmpFile.Close()

	oldPath := envFilePath
	envFilePath = tmpFile.Name()
	defer func() { envFilePath = oldPath }()

	// Execute execution loop (replaces variable)
	StoreSecretKey()

	// Execute append loop loop execution
	updatedData, err := os.ReadFile(envFilePath)
	if err != nil {
		t.Fatal(err)
	}

	content := string(updatedData)
	if !strings.Contains(content, "EXISTING_VAR=keep_me") {
		t.Errorf("clobber protection failed; missing existing components")
	}
	if strings.Contains(content, "old_token_to_replace") {
		t.Errorf("key translation failure; old string remains present")
	}
	if !strings.Contains(content, jwtSecretKey+"=") {
		t.Errorf("expected targeted key assignments to be generated")
	}
}

func TestGenerateSecretKeyProperties(t *testing.T) {
	key := generateSecretKey()
	if key == "" {
		t.Errorf("expected generated secret key string contents to be populated")
	}
	if len(key) < 32 {
		t.Errorf("expected key length matching cryptographic size targets, got %d", len(key))
	}
}