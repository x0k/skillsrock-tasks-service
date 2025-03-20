package analytics

import (
	"errors"

	"github.com/x0k/skillrock-tasks-service/internal/tasks"
)

var ErrReportNotFound = errors.New("report not found")

type Report struct {
	TasksCountByStatus        map[tasks.Status]int64
	AverageTaskCompletionTime float64
	AmountOfCompletedTasks    int64
	AmountOfOverdueTasks      int64
}
