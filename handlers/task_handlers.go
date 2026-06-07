package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/youssef-abbih/go-todo-list/models"
	"github.com/youssef-abbih/go-todo-list/utils"
)

// GetTasks godoc
// @Summary Retrieve all tasks for the authenticated user
// @Description Get a list of all tasks for the currently authenticated user, based on the JWT token.
// @Tags tasks
// @Produce json
// @Success 200 {array} models.Task "List of tasks"
// @Failure 401 {string} string "Unauthorized"
// @Security BearerAuth
// @Router /tasks [get]
func GetTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDUint, err := utils.GetUserID(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// 3. Fetch tasks for this user only
	tasks := models.GetTasks(userIDUint)

	// 4. Return tasks as JSON
	w.Header().Set(ContentTypeHeader, MimeJSON)
	if tasks == nil {
		tasks = []models.Task{}
	}
	json.NewEncoder(w).Encode(tasks)
}

// PostTask godoc
// @Summary Create a new task
// @Description Create and store a new task for the authenticated user.
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body models.Task true "Task to be created"
// @Success 201 {object} models.Task "Created task"
// @Failure 400 {string} string "Invalid input"
// @Failure 401 {string} string "Unauthorized"
// @Security BearerAuth
// @Router /tasks [post]
func PostTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var newTask models.Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(newTask.Title) == "" || strings.TrimSpace(newTask.Description) == "" {
		http.Error(w, "Title and description are required", http.StatusBadRequest)
		return
	}

	userIDUint, err := utils.GetUserID(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	created := models.AddTask(newTask, userIDUint)

	w.Header().Set(ContentTypeHeader, MimeJSON)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

// GetTask godoc
// @Summary Get task by ID
// @Description Retrieve a specific task by ID, if it belongs to the authenticated user.
// @Tags tasks
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} models.Task "The requested task"
// @Failure 400 {string} string "Invalid Task ID"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string ResponseTaskNotFound
// @Security BearerAuth
// @Router /tasks/{id} [get]
func GetTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)

	if err != nil || id <= 0 {
		http.Error(w, "Invalid Task ID", http.StatusBadRequest)
		return
	}

	idUint := uint(id)
	userIDUint, err := utils.GetUserID(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	task, found := models.GetTaskByID(idUint, userIDUint)
	if !found {
		http.Error(w, ResponseTaskNotFound, http.StatusNotFound)
		return
	}

	w.Header().Set(ContentTypeHeader, MimeJSON)
	json.NewEncoder(w).Encode(task)
}

// / DeleteTask godoc
// @Summary Delete task by ID
// @Description Delete a task by ID, if it belongs to the authenticated user.
// @Tags tasks
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} models.Task "Deleted task"
// @Failure 400 {string} string "Invalid Task ID"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string ResponseTaskNotFound
// @Security BearerAuth
// @Router /tasks/{id} [delete]
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid Task ID", http.StatusBadRequest)
		return
	}

	idUint := uint(id)
	userIDUint, err := utils.GetUserID(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	deleted, found := models.DeleteTask(idUint, userIDUint)
	if !found {
		http.Error(w, ResponseTaskNotFound, http.StatusNotFound)
		return
	}

	w.Header().Set(ContentTypeHeader, MimeJSON)
	json.NewEncoder(w).Encode(deleted)
}

// PutTask godoc
// @Summary Update task by ID
// @Description Update an existing task by ID, if it belongs to the authenticated user.
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param task body models.Task true "Updated task data"
// @Success 200 {object} models.Task "Updated task"
// @Failure 400 {string} string "Invalid Task ID or JSON"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string ResponseTaskNotFound
// @Security BearerAuth
// @Router /tasks/{id} [put]
func PutTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Task ID", http.StatusBadRequest)
		return
	}

	var updatedTask models.Task
	err = json.NewDecoder(r.Body).Decode(&updatedTask)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(updatedTask.Title) == "" || strings.TrimSpace(updatedTask.Description) == "" {
		http.Error(w, "Title and description are required", http.StatusBadRequest)
		return
	}

	idUint := uint(id)
	userIDUint, err := utils.GetUserID(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	result, ok := models.UpdateTask(idUint, userIDUint, updatedTask)
	if !ok {
		http.Error(w, ResponseTaskNotFound, http.StatusNotFound)
		return
	}

	w.Header().Set(ContentTypeHeader, MimeJSON)
	json.NewEncoder(w).Encode(result)
}
