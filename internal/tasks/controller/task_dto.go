package tasks_controller

import (
	"time"

	"github.com/google/uuid"
	"github.com/x0k/skillrock-tasks-service/internal/tasks"
)

type TaskDTO struct {
	Id          string  `json:"id" validate:"required"`
	Title       string  `json:"title" validate:"required"`
	Description *string `json:"description,omitempty"`
	Status      string  `json:"status" validate:"required"`
	Priority    string  `json:"priority" validate:"required"`
	DueDate     string  `json:"due_date" validate:"required"`
	CreatedAt   string  `json:"created_at" validate:"required"`
	UpdatedAt   string  `json:"updated_at" validate:"required"`
}

func taskToDTO(task tasks.Task) TaskDTO {
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

func taskFromDTO(dto TaskDTO) (tasks.Task, error) {
	task := tasks.Task{
		Title:       dto.Title,
		Description: dto.Description,
	}
	var err error
	if task.Id, err = tasks.NewTaskId(dto.Id); err != nil {
		return task, err
	}
	if task.Status, err = tasks.NewStatus(dto.Status); err != nil {
		return task, err
	}
	if task.Priority, err = tasks.NewPriority(dto.Priority); err != nil {
		return task, err
	}
	if task.DueDate, err = time.Parse(time.DateOnly, dto.DueDate); err != nil {
		return task, err
	}
	if task.CreatedAt, err = time.Parse(time.DateOnly, dto.CreatedAt); err != nil {
		return task, err
	}
	if task.UpdatedAt, err = time.Parse(time.DateOnly, dto.UpdatedAt); err != nil {
		return task, err
	}
	return task, nil
}
