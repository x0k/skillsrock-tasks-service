package tasks

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrInvalidStatus = errors.New("invalid status")
var ErrInvalidPriority = errors.New("invalid priority")

type Status string

func (s Status) String() string {
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

func NewTask(
	title string,
	description *string,
	status Status,
	priority Priority,
	dueDate time.Time,
) Task {
	now := time.Now()
	return Task{
		Id:          TaskId(uuid.New()),
		Title:       title,
		Description: description,
		Status:      status,
		Priority:    priority,
		DueDate:     dueDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
