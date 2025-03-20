// Code generated by mockery. DO NOT EDIT.

package analytics

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockAnalyticsRepo is an autogenerated mock type for the AnalyticsRepo type
type MockAnalyticsRepo struct {
	mock.Mock
}

type MockAnalyticsRepo_Expecter struct {
	mock *mock.Mock
}

func (_m *MockAnalyticsRepo) EXPECT() *MockAnalyticsRepo_Expecter {
	return &MockAnalyticsRepo_Expecter{mock: &_m.Mock}
}

// Report provides a mock function with given fields: ctx
func (_m *MockAnalyticsRepo) Report(ctx context.Context) (Report, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Report")
	}

	var r0 Report
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (Report, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) Report); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(Report)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAnalyticsRepo_Report_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Report'
type MockAnalyticsRepo_Report_Call struct {
	*mock.Call
}

// Report is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockAnalyticsRepo_Expecter) Report(ctx interface{}) *MockAnalyticsRepo_Report_Call {
	return &MockAnalyticsRepo_Report_Call{Call: _e.mock.On("Report", ctx)}
}

func (_c *MockAnalyticsRepo_Report_Call) Run(run func(ctx context.Context)) *MockAnalyticsRepo_Report_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockAnalyticsRepo_Report_Call) Return(_a0 Report, _a1 error) *MockAnalyticsRepo_Report_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAnalyticsRepo_Report_Call) RunAndReturn(run func(context.Context) (Report, error)) *MockAnalyticsRepo_Report_Call {
	_c.Call.Return(run)
	return _c
}

// SaveReport provides a mock function with given fields: ctx, report
func (_m *MockAnalyticsRepo) SaveReport(ctx context.Context, report Report) error {
	ret := _m.Called(ctx, report)

	if len(ret) == 0 {
		panic("no return value specified for SaveReport")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, Report) error); ok {
		r0 = rf(ctx, report)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockAnalyticsRepo_SaveReport_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SaveReport'
type MockAnalyticsRepo_SaveReport_Call struct {
	*mock.Call
}

// SaveReport is a helper method to define mock.On call
//   - ctx context.Context
//   - report Report
func (_e *MockAnalyticsRepo_Expecter) SaveReport(ctx interface{}, report interface{}) *MockAnalyticsRepo_SaveReport_Call {
	return &MockAnalyticsRepo_SaveReport_Call{Call: _e.mock.On("SaveReport", ctx, report)}
}

func (_c *MockAnalyticsRepo_SaveReport_Call) Run(run func(ctx context.Context, report Report)) *MockAnalyticsRepo_SaveReport_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(Report))
	})
	return _c
}

func (_c *MockAnalyticsRepo_SaveReport_Call) Return(_a0 error) *MockAnalyticsRepo_SaveReport_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAnalyticsRepo_SaveReport_Call) RunAndReturn(run func(context.Context, Report) error) *MockAnalyticsRepo_SaveReport_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockAnalyticsRepo creates a new instance of MockAnalyticsRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockAnalyticsRepo(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockAnalyticsRepo {
	mock := &MockAnalyticsRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
