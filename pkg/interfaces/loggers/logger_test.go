package loggers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
)

type loggerMock struct {
	mock.Mock
	t *testing.T
}

func (m *loggerMock) Debug(format string, params ...interface{}) {
	fmt.Sprintf(format, params...) // nolint: vet,megacheck
}
func (m *loggerMock) Info(format string, params ...interface{}) {
	fmt.Sprintf(format, params...) // nolint: vet,megacheck
}
func (m *loggerMock) Warn(format string, params ...interface{}) {
	fmt.Sprintf(format, params...) // nolint: vet,megacheck
}
func (m *loggerMock) Error(format string, params ...interface{}) {
	fmt.Sprintf(format, params...) // nolint: vet,megacheck
}
func (m *loggerMock) Crit(format string, params ...interface{}) {
	fmt.Sprintf(format, params...) // nolint: vet,megacheck
}
func (m *loggerMock) Success(format string, params ...interface{}) {
	fmt.Sprintf(format, params...) // nolint: vet,megacheck
}
