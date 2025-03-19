package tasks

import (
	"time"

	"github.com/google/uuid"
)

type CreateTaskDTO struct {
	Title       string  `json:"title" validate:"required"`
	Description *string `json:"description,omitempty"`
	Status      string  `json:"status" validate:"required"`
	Priority    string  `json:"priority" validate:"required"`
	DueDate     string  `json:"due_date" validate:"required"`
}

type FindTasksDTO struct {
	Title     string
	Status    string
	Priority  string
	DueBefore string
	DueAfter  string
}

type TaskDTO struct {
	Id          string  `json:"id"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	Status      string  `json:"status"`
	Priority    string  `json:"priority"`
	DueDate     string  `json:"due_date"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

func taskToDTO(task Task) TaskDTO {
	return TaskDTO{
		Id:          uuid.UUID(task.Id).String(),
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status.String(),
		Priority:    task.Priority.String(),
		DueDate:     task.DueDate.Format(time.DateOnly),
		CreatedAt:   task.CreatedAt.Format(time.DateOnly),
		UpdatedAt:   task.UpdatedAt.Format(time.DateOnly),
	}
}
