package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/youssef-abbih/go-todo-list/models"
	"github.com/youssef-abbih/go-todo-list/utils"
	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if user.Email == "" || user.Password == "" {
		http.Error(w, "email or password are missing", http.StatusBadRequest)
		return
	}

	if _, found := models.GetUserByEmail(user.Email); found {
		http.Error(w, "User already exist", http.StatusConflict)
		return
	}

	var err error
	user.Password, err = models.HashPassword(user.Password)

	if err != nil {
		http.Error(w, "Error while hashing the password", http.StatusBadRequest)
		return
	}

	createdUser, err := models.AddUser(user)
	if err != nil {
		http.Error(w, "Error saving user to database", http.StatusInternalServerError)
		return
	}
	w.Header().Set(ContentTypeHeader, MimeJSON)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": http.StatusCreated,
		"user": map[string]interface{}{
			"id":    createdUser.ID,
			"email": createdUser.Email,
		},
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	existingUser, found := models.GetUserByEmail(user.Email)

	if !found {
		http.Error(w, "Invalid email or password", http.StatusBadRequest)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password)); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	claims := jwt.MapClaims{
		"sub":   existingUser.ID,
		"email": existingUser.Email,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}
	secretKey := utils.LoadJWTSecretkey()
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedJwtToken, err := jwtToken.SignedString([]byte(secretKey))

	if err != nil {
		http.Error(w, "Error while signing the JWT Token", http.StatusInternalServerError)
		return
	}

	w.Header().Set(ContentTypeHeader, MimeJSON)
	json.NewEncoder(w).Encode(map[string]string{"Token": signedJwtToken})

}
