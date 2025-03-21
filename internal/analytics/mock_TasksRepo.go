// Code generated by mockery. DO NOT EDIT.

package analytics

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	tasks "github.com/x0k/skillrock-tasks-service/internal/tasks"

	time "time"
)

// MockTasksRepo is an autogenerated mock type for the TasksRepo type
type MockTasksRepo struct {
	mock.Mock
}

type MockTasksRepo_Expecter struct {
	mock *mock.Mock
}

func (_m *MockTasksRepo) EXPECT() *MockTasksRepo_Expecter {
	return &MockTasksRepo_Expecter{mock: &_m.Mock}
}

// AverageCompletionTime provides a mock function with given fields: ctx
func (_m *MockTasksRepo) AverageCompletionTime(ctx context.Context) (float64, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for AverageCompletionTime")
	}

	var r0 float64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (float64, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) float64); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(float64)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockTasksRepo_AverageCompletionTime_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AverageCompletionTime'
type MockTasksRepo_AverageCompletionTime_Call struct {
	*mock.Call
}

// AverageCompletionTime is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockTasksRepo_Expecter) AverageCompletionTime(ctx interface{}) *MockTasksRepo_AverageCompletionTime_Call {
	return &MockTasksRepo_AverageCompletionTime_Call{Call: _e.mock.On("AverageCompletionTime", ctx)}
}

func (_c *MockTasksRepo_AverageCompletionTime_Call) Run(run func(ctx context.Context)) *MockTasksRepo_AverageCompletionTime_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockTasksRepo_AverageCompletionTime_Call) Return(_a0 float64, _a1 error) *MockTasksRepo_AverageCompletionTime_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockTasksRepo_AverageCompletionTime_Call) RunAndReturn(run func(context.Context) (float64, error)) *MockTasksRepo_AverageCompletionTime_Call {
	_c.Call.Return(run)
	return _c
}

// CountCompletedAndOverdueTasks provides a mock function with given fields: ctx, date
func (_m *MockTasksRepo) CountCompletedAndOverdueTasks(ctx context.Context, date time.Time) (int64, int64, error) {
	ret := _m.Called(ctx, date)

	if len(ret) == 0 {
		panic("no return value specified for CountCompletedAndOverdueTasks")
	}

	var r0 int64
	var r1 int64
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, time.Time) (int64, int64, error)); ok {
		return rf(ctx, date)
	}
	if rf, ok := ret.Get(0).(func(context.Context, time.Time) int64); ok {
		r0 = rf(ctx, date)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, time.Time) int64); ok {
		r1 = rf(ctx, date)
	} else {
		r1 = ret.Get(1).(int64)
	}

	if rf, ok := ret.Get(2).(func(context.Context, time.Time) error); ok {
		r2 = rf(ctx, date)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockTasksRepo_CountCompletedAndOverdueTasks_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CountCompletedAndOverdueTasks'
type MockTasksRepo_CountCompletedAndOverdueTasks_Call struct {
	*mock.Call
}

// CountCompletedAndOverdueTasks is a helper method to define mock.On call
//   - ctx context.Context
//   - date time.Time
func (_e *MockTasksRepo_Expecter) CountCompletedAndOverdueTasks(ctx interface{}, date interface{}) *MockTasksRepo_CountCompletedAndOverdueTasks_Call {
	return &MockTasksRepo_CountCompletedAndOverdueTasks_Call{Call: _e.mock.On("CountCompletedAndOverdueTasks", ctx, date)}
}

func (_c *MockTasksRepo_CountCompletedAndOverdueTasks_Call) Run(run func(ctx context.Context, date time.Time)) *MockTasksRepo_CountCompletedAndOverdueTasks_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(time.Time))
	})
	return _c
}

func (_c *MockTasksRepo_CountCompletedAndOverdueTasks_Call) Return(_a0 int64, _a1 int64, _a2 error) *MockTasksRepo_CountCompletedAndOverdueTasks_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MockTasksRepo_CountCompletedAndOverdueTasks_Call) RunAndReturn(run func(context.Context, time.Time) (int64, int64, error)) *MockTasksRepo_CountCompletedAndOverdueTasks_Call {
	_c.Call.Return(run)
	return _c
}

// TasksCountByStatus provides a mock function with given fields: ctx
func (_m *MockTasksRepo) TasksCountByStatus(ctx context.Context) (map[tasks.Status]int64, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for TasksCountByStatus")
	}

	var r0 map[tasks.Status]int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (map[tasks.Status]int64, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) map[tasks.Status]int64); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[tasks.Status]int64)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockTasksRepo_TasksCountByStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'TasksCountByStatus'
type MockTasksRepo_TasksCountByStatus_Call struct {
	*mock.Call
}

// TasksCountByStatus is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockTasksRepo_Expecter) TasksCountByStatus(ctx interface{}) *MockTasksRepo_TasksCountByStatus_Call {
	return &MockTasksRepo_TasksCountByStatus_Call{Call: _e.mock.On("TasksCountByStatus", ctx)}
}

func (_c *MockTasksRepo_TasksCountByStatus_Call) Run(run func(ctx context.Context)) *MockTasksRepo_TasksCountByStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockTasksRepo_TasksCountByStatus_Call) Return(_a0 map[tasks.Status]int64, _a1 error) *MockTasksRepo_TasksCountByStatus_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockTasksRepo_TasksCountByStatus_Call) RunAndReturn(run func(context.Context) (map[tasks.Status]int64, error)) *MockTasksRepo_TasksCountByStatus_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockTasksRepo creates a new instance of MockTasksRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockTasksRepo(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockTasksRepo {
	mock := &MockTasksRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
