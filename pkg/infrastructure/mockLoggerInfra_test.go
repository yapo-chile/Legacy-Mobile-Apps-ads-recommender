package infrastructure

import (
	"github.com/stretchr/testify/mock"
)

// MockLoggerRepository simulate Logger Repo
type MockLoggerInfrastructure struct {
	mock.Mock
}

// Info simulate Info
func (m *MockLoggerInfrastructure) Info(message string, params ...interface{}) {
	m.Called()
}

// Debug simulate Debug Logger
func (m *MockLoggerInfrastructure) Debug(message string, params ...interface{}) {
	m.Called()
}

// Crit simulate Crit Logger
func (m *MockLoggerInfrastructure) Crit(message string, params ...interface{}) {
	m.Called()
}

// Error simulate Error Logger
func (m *MockLoggerInfrastructure) Error(message string, params ...interface{}) {
	m.Called()
}

// Warn simulate Warn Logger
func (m *MockLoggerInfrastructure) Warn(message string, params ...interface{}) {
	m.Called()
}

// Success simulate Warn Logger
func (m *MockLoggerInfrastructure) Success(message string, params ...interface{}) {
	m.Called()
}
