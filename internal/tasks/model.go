package tasks

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrInvalidStatus = errors.New("invalid status")
var ErrInvalidPriority = errors.New("invalid priority")

type Status string

func (s Status) ToString() string {
	return string(s)
}

const (
	Pending    Status = "pending"
	InProgress Status = "in_progress"
	Done       Status = "done"
)

var statuses = map[string]Status{
	string(Pending):    Pending,
	string(InProgress): InProgress,
	string(Done):       Done,
}

func newStatus(st string) (Status, error) {
	status, ok := statuses[st]
	if !ok {
		return Pending, ErrInvalidStatus
	}
	return status, nil
}

type Priority string

func (p Priority) ToString() string {
	return string(p)
}

const (
	Low    Priority = "low"
	Medium Priority = "medium"
	High   Priority = "high"
)

var priorities = map[string]Priority{
	string(Low):    Low,
	string(Medium): Medium,
	string(High):   High,
}

func newPriority(pr string) (Priority, error) {
	priority, ok := priorities[pr]
	if !ok {
		return Low, ErrInvalidPriority
	}
	return priority, nil
}

type taskId uuid.UUID

type task struct {
	Id          taskId
	Title       string
	Description *string
	Status      Status
	Priority    Priority
	DueDate     time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func newTask(dto CreateTaskDTO) task {
	now := time.Now()
	return task{
		Id:          taskId(uuid.New()),
		Title:       dto.Title,
		Description: dto.Description,
		Status:      dto.Status,
		Priority:    dto.Priority,
		DueDate:     dto.DueDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
