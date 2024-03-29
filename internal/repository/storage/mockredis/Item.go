// Code generated by mockery v2.20.0. DO NOT EDIT.

package mockCache

import (
	item "github.com/evgeniy-dammer/marketplace-api/internal/domain/item"
	context "github.com/evgeniy-dammer/marketplace-api/pkg/context"

	mock "github.com/stretchr/testify/mock"

	query "github.com/evgeniy-dammer/marketplace-api/pkg/query"

	queryparameter "github.com/evgeniy-dammer/marketplace-api/pkg/queryparameter"
)

// Item is an autogenerated mock type for the Item type
type Item struct {
	mock.Mock
}

// ItemCreate provides a mock function with given fields: ctx, _a1
func (_m *Item) ItemCreate(ctx context.Context, _a1 item.Item) error {
	ret := _m.Called(ctx, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, item.Item) error); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ItemDelete provides a mock function with given fields: ctx, itemID
func (_m *Item) ItemDelete(ctx context.Context, itemID string) error {
	ret := _m.Called(ctx, itemID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, itemID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ItemGetAll provides a mock function with given fields: ctx, meta, params
func (_m *Item) ItemGetAll(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]item.Item, error) {
	ret := _m.Called(ctx, meta, params)

	var r0 []item.Item
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, query.MetaData, queryparameter.QueryParameter) ([]item.Item, error)); ok {
		return rf(ctx, meta, params)
	}
	if rf, ok := ret.Get(0).(func(context.Context, query.MetaData, queryparameter.QueryParameter) []item.Item); ok {
		r0 = rf(ctx, meta, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]item.Item)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, query.MetaData, queryparameter.QueryParameter) error); ok {
		r1 = rf(ctx, meta, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ItemGetOne provides a mock function with given fields: ctx, itemID
func (_m *Item) ItemGetOne(ctx context.Context, itemID string) (item.Item, error) {
	ret := _m.Called(ctx, itemID)

	var r0 item.Item
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (item.Item, error)); ok {
		return rf(ctx, itemID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) item.Item); ok {
		r0 = rf(ctx, itemID)
	} else {
		r0 = ret.Get(0).(item.Item)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, itemID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ItemInvalidate provides a mock function with given fields: ctx
func (_m *Item) ItemInvalidate(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ItemSetAll provides a mock function with given fields: ctx, meta, params, items
func (_m *Item) ItemSetAll(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter, items []item.Item) error {
	ret := _m.Called(ctx, meta, params, items)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, query.MetaData, queryparameter.QueryParameter, []item.Item) error); ok {
		r0 = rf(ctx, meta, params, items)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ItemUpdate provides a mock function with given fields: ctx, _a1
func (_m *Item) ItemUpdate(ctx context.Context, _a1 item.Item) error {
	ret := _m.Called(ctx, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, item.Item) error); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewItem interface {
	mock.TestingT
	Cleanup(func())
}

// NewItem creates a new instance of Item. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewItem(t mockConstructorTestingTNewItem) *Item {
	mock := &Item{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
