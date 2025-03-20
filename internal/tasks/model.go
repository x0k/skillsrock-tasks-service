package tasks

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrInvalidStatus = errors.New("invalid status")
var ErrInvalidPriority = errors.New("invalid priority")
var ErrTaskNotFound = errors.New("task not found")
var ErrTaskIsAlreadyDone = errors.New("task is already done")
var ErrInvalidTasksTitle = errors.New("invalid task title")
var ErrInvalidDueDate = errors.New("invalid due data")

type Status string

func (s Status) String() string {
	return string(s)
}

func (s Status) IsValid() bool {
	_, ok := statuses[string(s)]
	return ok
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

func NewStatus(st string) (Status, error) {
	s, ok := statuses[st]
	if !ok {
		return Pending, ErrInvalidStatus
	}
	return s, nil
}

type Priority string

func (p Priority) String() string {
	return string(p)
}

func (p Priority) IsValid() bool {
	_, ok := priorities[string(p)]
	return ok
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

func NewPriority(pr string) (Priority, error) {
	p, ok := priorities[pr]
	if !ok {
		return Low, ErrInvalidPriority
	}
	return p, nil
}

type TaskId uuid.UUID

func (id TaskId) String() string {
	return uuid.UUID(id).String()
}

func NewTaskId(id string) (TaskId, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return TaskId(uuid.Nil), err
	}
	return TaskId(uid), nil
}

type Task struct {
	Id          TaskId
	Title       string
	Description *string
	Status      Status
	Priority    Priority
	DueDate     time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type TaskParams struct {
	Title       string
	Description *string
	Status      Status
	Priority    Priority
	DueDate     time.Time
}

func NewTask(params TaskParams) (Task, error) {
	now := time.Now()
	if len(params.Title) == 0 {
		return Task{}, ErrInvalidTasksTitle
	}
	if params.DueDate.Before(now) {
		return Task{}, ErrInvalidDueDate
	}
	if !params.Status.IsValid() {
		return Task{}, ErrInvalidStatus
	}
	if !params.Priority.IsValid() {
		return Task{}, ErrInvalidPriority
	}
	return Task{
		Id:          TaskId(uuid.New()),
		Title:       params.Title,
		Description: params.Description,
		Status:      params.Status,
		Priority:    params.Priority,
		DueDate:     params.DueDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

type TasksFilter struct {
	Title     *string
	Status    *Status
	Priority  *Priority
	DueBefore *time.Time
	DueAfter  *time.Time
}

func (f TasksFilter) IsEmpty() bool {
	return f.Title == nil && f.Status == nil && f.Priority == nil && f.DueBefore == nil && f.DueAfter == nil
}
