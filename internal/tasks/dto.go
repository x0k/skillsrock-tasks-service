package tasks

import "time"

type CreateTaskDTO struct {
	Title       string
	Description *string
	Status      Status
	Priority    Priority
	DueDate     time.Time
}

type FindTasksDTO struct {
	Status    *Status
	Priority  *Priority
	DueBefore *time.Time
	DueAfter  *time.Time
	Title     *string
}
