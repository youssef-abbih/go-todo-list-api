package utils

import (
	"os"
	"log"
	"strings"
	"crypto/rand"
	"encoding/base64"
	"github.com/joho/godotenv"
)

func generateSecretKey() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(b)
}

func StoreSecretKey() {
	key := generateSecretKey()
	envPath := ".env"

	data, _ := os.ReadFile(envPath)
	lines := []string{}
	found := false

	for _, line := range strings.Split(string(data),"\n") {
		if strings.HasPrefix(line, "JWT_SECRET=") {
			line = "JWT_SECRET=" + key
			found = true
		}
		lines = append(lines, line)
	}

	if !found {
		lines = append(lines, "JWT_SECRET=" + key)
	}
	// Write or overwrite the JWT_SECRET line
	err := os.WriteFile(envPath, []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		panic(err)
	}
}

func LoadJWTSecretkey() string{
	if os.Getenv("JWT_SECRET") == "" {
		_ = godotenv.Load("../.env")
	}
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		log.Fatal("Secret Key cannot be empty")
	}
	return secretKey
}