package handlers

import (
	"github.com/stretchr/testify/mock"
)

// MockInputRequest is a mock class
type MockInputRequest struct {
	mock.Mock
}

// Set is a mocked method
func (m *MockInputRequest) Set(input interface{}) TargetRequest {
	args := m.Called(input)
	return args.Get(0).(TargetRequest)
}

// MockTargetRequest is a mock class
type MockTargetRequest struct {
	mock.Mock
}

// FromJSONBody is a mocked method
func (m *MockTargetRequest) FromJSONBody() TargetRequest {
	m.Called()
	return m
}

// FromRawBody is a mocked method
func (m *MockTargetRequest) FromRawBody() TargetRequest {
	m.Called()
	return m
}

// FromPath is a mocked method
func (m *MockTargetRequest) FromPath() TargetRequest {
	m.Called()
	return m
}

// FromQuery is a mocked method
func (m *MockTargetRequest) FromQuery() TargetRequest {
	m.Called()
	return m
}

// FromHeaders is a mocked method
func (m *MockTargetRequest) FromHeaders() TargetRequest {
	m.Called()
	return m
}

// FromCookies is a mocked method
func (m *MockTargetRequest) FromCookies() TargetRequest {
	m.Called()
	return m
}

// FromForm is a mocked method
func (m *MockTargetRequest) FromForm() TargetRequest {
	m.Called()
	return m
}
