// Code generated by mockery v2.20.0. DO NOT EDIT.

package mockCache

import (
	context "github.com/evgeniy-dammer/marketplace-api/pkg/context"
	mock "github.com/stretchr/testify/mock"

	query "github.com/evgeniy-dammer/marketplace-api/pkg/query"

	queryparameter "github.com/evgeniy-dammer/marketplace-api/pkg/queryparameter"

	table "github.com/evgeniy-dammer/marketplace-api/internal/domain/table"
)

// Table is an autogenerated mock type for the Table type
type Table struct {
	mock.Mock
}

// TableCreate provides a mock function with given fields: ctx, _a1
func (_m *Table) TableCreate(ctx context.Context, _a1 table.Table) error {
	ret := _m.Called(ctx, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, table.Table) error); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TableDelete provides a mock function with given fields: ctx, tableID
func (_m *Table) TableDelete(ctx context.Context, tableID string) error {
	ret := _m.Called(ctx, tableID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, tableID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TableGetAll provides a mock function with given fields: ctx, meta, params
func (_m *Table) TableGetAll(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]table.Table, error) {
	ret := _m.Called(ctx, meta, params)

	var r0 []table.Table
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, query.MetaData, queryparameter.QueryParameter) ([]table.Table, error)); ok {
		return rf(ctx, meta, params)
	}
	if rf, ok := ret.Get(0).(func(context.Context, query.MetaData, queryparameter.QueryParameter) []table.Table); ok {
		r0 = rf(ctx, meta, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]table.Table)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, query.MetaData, queryparameter.QueryParameter) error); ok {
		r1 = rf(ctx, meta, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TableGetOne provides a mock function with given fields: ctx, tableID
func (_m *Table) TableGetOne(ctx context.Context, tableID string) (table.Table, error) {
	ret := _m.Called(ctx, tableID)

	var r0 table.Table
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (table.Table, error)); ok {
		return rf(ctx, tableID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) table.Table); ok {
		r0 = rf(ctx, tableID)
	} else {
		r0 = ret.Get(0).(table.Table)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, tableID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TableInvalidate provides a mock function with given fields: ctx
func (_m *Table) TableInvalidate(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TableSetAll provides a mock function with given fields: ctx, meta, params, tables
func (_m *Table) TableSetAll(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter, tables []table.Table) error {
	ret := _m.Called(ctx, meta, params, tables)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, query.MetaData, queryparameter.QueryParameter, []table.Table) error); ok {
		r0 = rf(ctx, meta, params, tables)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TableUpdate provides a mock function with given fields: ctx, _a1
func (_m *Table) TableUpdate(ctx context.Context, _a1 table.Table) error {
	ret := _m.Called(ctx, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, table.Table) error); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewTable interface {
	mock.TestingT
	Cleanup(func())
}

// NewTable creates a new instance of Table. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTable(t mockConstructorTestingTNewTable) *Table {
	mock := &Table{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
