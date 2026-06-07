package models

import (
	"gorm.io/gorm"
	"time"
)

const queryTaskWithUser = "id = ? AND user_id = ?"

// Task represents a single to-do item.
type Task struct {
	ID          uint           	`json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time      	`json:"created_at"`
	UpdatedAt   time.Time      	`json:"updated_at"`
	DeletedAt   gorm.DeletedAt 	`gorm:"index" json:"-"`
	Title       string         	`json:"title"`
	Description string         	`json:"description"`
	Completed   bool           	`json:"completed"`
	UserID 		uint 			`json:"user_id"`
	User   		User 			`json:"-" gorm:"foreignKey:UserID"`
}

// GetTasks retrieves all tasks from the database
func GetTasks(userID uint) []Task {
	var tasks []Task
	DB.Where("user_id = ?", userID).Find(&tasks)
	return tasks
}

// AddTask adds a new task and returns it with its ID set by the DB
func AddTask(task Task, userID uint) Task {
	task.UserID = userID

	task.CreatedAt = time.Now()

	DB.Create(&task)
	return task
}

// GetTaskByID retrieves a single task by its ID
func GetTaskByID(id , userID uint) (Task, bool) {
	var task Task
	result := DB.Where(queryTaskWithUser, id, userID).First(&task, id)
	if result.Error != nil {
		return Task{}, false
	}
	return task, true
}

// UpdateTask updates the task with the given ID
func UpdateTask(id, userID uint, updated Task) (Task, bool) {
	var existing Task
	result := DB.Where(queryTaskWithUser, id, userID).First(&existing)
	if result.Error != nil {
		return Task{}, false
	}

	if updated.Title == existing.Title &&
		updated.Description == existing.Description &&
		updated.Completed == existing.Completed {
		return existing, true
	}

	updated.ID = existing.ID
	updated.CreatedAt = existing.CreatedAt
	updated.UserID = existing.UserID
	DB.Save(&updated)
	return updated, true
}

// DeleteTask deletes a task by ID
func DeleteTask(id, userID uint) (Task, bool) {
	var task Task
	result := DB.Where(queryTaskWithUser, id, userID).Unscoped().First(&task)
	if result.Error != nil {
		return Task{}, false
	}
	if task.DeletedAt.Valid {
		return Task{}, false // Already deleted
	}

	DB.Delete(&task)
	return task, true
}
