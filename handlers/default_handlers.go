package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/youssef-abbih/go-todo-list/models"
)

// DefaultResponse godoc
// @Summary Welcome message
// @Description Returns a welcome message from the Todo List API
// @Tags Default
// @Produce plain
// @Success 200 {string} string "Welcome to my Todo List API"
// @Router / [get]
func DefaultResponse(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"message": "Welcome to my Todo List API"})

}

// HealthCheck godoc
// @Summary Health check
// @Description Checks if the database connection is alive
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 500 {string} string "DB unreachable"
// @Router /health [get]
func HealthCheck(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	sqlDB, _ := models.DB.DB()
	err := sqlDB.Ping()
	if err != nil {
		http.Error(w, "DB unreachable", http.StatusInternalServerError)
		return
	}

	w.Header().Set(ContentTypeHeader, MimeJSON)
	json.NewEncoder(w).Encode(map[string]string{"Status": "Ok"})
}
