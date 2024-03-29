// Code generated by mockery v2.20.0. DO NOT EDIT.

package mockCache

import mock "github.com/stretchr/testify/mock"

// Authentication is an autogenerated mock type for the Authentication type
type Authentication struct {
	mock.Mock
}

type mockConstructorTestingTNewAuthentication interface {
	mock.TestingT
	Cleanup(func())
}

// NewAuthentication creates a new instance of Authentication. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAuthentication(t mockConstructorTestingTNewAuthentication) *Authentication {
	mock := &Authentication{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
