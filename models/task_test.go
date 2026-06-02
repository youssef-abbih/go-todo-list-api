package models

import (
	"testing"
	"time"
	"fmt"
)

func setupTestData() (userID uint, taskID uint) {
    InitDB()
    
    // Create a user first
    user := User{Email: fmt.Sprintf("test_%d@test.com", time.Now().UnixNano())}
    user.Password, _ = HashPassword("password")
    created, _ := AddUser(user)
    userID = created.ID
    
    // Create a task for that user
    task := Task{Title: "Test Task", Description: "This is a test", Completed: false}
    added := AddTask(task, userID)
    taskID = added.ID
    
    return userID, taskID
}

// TestAddTask verifies that a task is correctly added and given an ID
func TestAddTask(t *testing.T) {
	// Reset the DB (in-memory or test DB setup is better, but this is simple)
	InitDB()
	userID, _ := setupTestData()
	
	task := Task{
		Title:       "Test Task",
		Description: "This is a test",
		Completed:   false,
	}

	added := AddTask(task, userID)

	// Check if an ID is assigned
	if added.ID == 0 {
		t.Errorf("expected task to have a non-zero ID, got %d", added.ID)
	}

	// Check if a UserID is assigned
	if added.UserID == 0 {
		t.Errorf("expected task to have a non-zero UserID, got %d", added.UserID)
	}

	// Check if a UserID is assigned
	if added.UserID != userID {
		t.Errorf("expected task to have a UserID %d, got %d", userID, added.UserID)
	}

	// Check if CreatedAt was set
	if time.Since(added.CreatedAt) > time.Second {
		t.Errorf("expected CreatedAt to be recent, got %v", added.CreatedAt)
	}

	nonExistingUser := "example@example.com"
	if _, exist := GetUserByEmail(nonExistingUser); exist {
		
		t.Errorf("User unexpectedly exists")
	}
}

func TestGetTasks(t *testing.T) {
	InitDB()
	existingUserID, _ := setupTestData()

	tasks := GetTasks(existingUserID)

	if len(tasks) == 0 {
		t.Errorf("Expected tasks for user %d but got none", existingUserID)
	}

	nonExistingUser := "example@example.com"
	if _, exist := GetUserByEmail(nonExistingUser); exist {
		
		t.Errorf("User unexpectedly exists")
	}
}

func TestGetTaskByID(t *testing.T) {

	InitDB()
	existingUserID , existingID := setupTestData()
	task, returned := GetTaskByID(existingID, existingUserID)

	if !returned {
		t.Errorf("Expected task with ID %d to be returned, but GetTaskByID returned false", existingID)
	}

	if task.ID != existingID {
		t.Errorf("returned task ID mismatch: got %d, want %d", task.ID, existingID)
	}

	// Test non-existing task
	var nonExistingID uint = 9999
	_, returned = GetTaskByID(nonExistingID, existingUserID)

	if returned {
		t.Errorf("Expected no task to be returned for ID %d, but got true", nonExistingID)
	}

	nonExistingUser := "example@example.com"
	if _, exist := GetUserByEmail(nonExistingUser); exist {
		
		t.Errorf("User unexpectedly exists")
	}
}

func TestDeleteTask(t *testing.T) {
	// Initialize DB and seed some tasks
	InitDB()

	existingUserID , existingID := setupTestData()
	// Try deleting a task with ID 1 (assuming it exists after InitDB)
	deletedTask, deleted := DeleteTask(existingID, existingUserID)

	if !deleted {
		t.Errorf("Expected task with ID 1 to be deleted, but DeleteTask returned false")
	}

	// Confirm the deleted task ID is 1
	if deletedTask.ID != existingID {
		t.Errorf("Deleted task ID mismatch: got %d, want %d", deletedTask.ID, 1)
	}

	_, found := GetTaskByID(existingID, existingUserID)
	if found {
		t.Errorf("Task with ID 1 should not exist after deletion")
	}

	// Try deleting a task that does not exist
	var nonExistingID uint = 9999
	_, deleted = DeleteTask(nonExistingID, existingUserID)
	if deleted {
		t.Errorf("Expected deletion of non-existent task to fail, but it succeeded")
	}

	nonExistingUser := "example@example.com"
	if _, exist := GetUserByEmail(nonExistingUser); exist {
		
		t.Errorf("User unexpectedly exists")
	}
}

func TestUpdateTask(t *testing.T) {
	InitDB()
	updatedTask := Task{
		Title:       "Test Task",
		Description: "This is a test",
		Completed:   false,
	}

	existingUserID , existingID := setupTestData()

	task, updated := UpdateTask(existingID, existingUserID, updatedTask)

	if !updated {
		t.Errorf("Expected task with ID %d to be updated, but UpdateTask returned false", existingID)
	}

	if task.ID != existingID {
		t.Errorf("updated task ID mismatch: got %d, want %d", task.ID, existingID)
	}

	// Test non-existing task
	var nonExistingID uint = 9999
	_, updated = UpdateTask(nonExistingID,existingUserID, updatedTask)

	if updated {
		t.Errorf("Expected no task to be updated for ID %d, but got true", nonExistingID)
	}

	nonExistingUser := "example@example.com"
	if _, exist := GetUserByEmail(nonExistingUser); exist {
		
		t.Errorf("User unexpectedly exists")
	}

}
