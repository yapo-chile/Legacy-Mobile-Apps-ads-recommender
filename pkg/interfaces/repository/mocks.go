package repository

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type MockElasticSearchHandler struct {
	mock.Mock
}

func (m *MockElasticSearchHandler) Info() (interface{}, error) {
	args := m.Called()
	return args.Get(0), args.Error(1)
}
func (m *MockElasticSearchHandler) Create(index string) error {
	args := m.Called(index)
	return args.Error(0)
}
func (m *MockElasticSearchHandler) PutMapping(mapping []byte, index string) error {
	args := m.Called(mapping, index)
	return args.Error(0)
}
func (m *MockElasticSearchHandler) Search(index, query string, size, from int) (string, error) {
	args := m.Called(index, query, size, from)
	return args.Get(0).(string), args.Error(1)
}

type MockDataMapping struct {
	mock.Mock
}

func (m *MockDataMapping) Get(s string) string {
	args := m.Called(s)
	return args.Get(0).(string)
}

type mockHTTPHandler struct { // nolint: deadcode
	mock.Mock
}

func (m *mockHTTPHandler) Send(request HTTPRequest) (interface{}, error) {
	args := m.Called(request)
	return args.Get(0), args.Error(1)
}

func (m *mockHTTPHandler) NewRequest() HTTPRequest {
	args := m.Called()
	return args.Get(0).(HTTPRequest)
}

type mockRequest struct { // nolint: deadcode
	mock.Mock
}

func (m *mockRequest) SetMethod(method string) HTTPRequest {
	args := m.Called(method)
	return args.Get(0).(HTTPRequest)
}

func (m *mockRequest) GetMethod() string {
	args := m.Called()
	return args.String(0)
}

func (m *mockRequest) SetPath(path string) HTTPRequest {
	args := m.Called(path)
	return args.Get(0).(HTTPRequest)
}

func (m *mockRequest) GetPath() string {
	args := m.Called()
	return args.String(0)
}

func (m *mockRequest) SetHeaders(headers map[string]string) HTTPRequest {
	args := m.Called(headers)
	return args.Get(0).(HTTPRequest)
}

func (m *mockRequest) GetHeaders() map[string][]string {
	args := m.Called()
	return args.Get(0).(map[string][]string)
}

func (m *mockRequest) SetBody(body interface{}) HTTPRequest {
	args := m.Called()
	return args.Get(0).(HTTPRequest)
}

func (m *mockRequest) GetBody() interface{} {
	args := m.Called()
	return args.Get(0)
}

func (m *mockRequest) SetQueryParams(queryParams map[string]string) HTTPRequest {
	args := m.Called(queryParams)
	return args.Get(0).(HTTPRequest)
}

func (m *mockRequest) GetQueryParams() map[string][]string {
	args := m.Called()
	return args.Get(0).(map[string][]string)
}

func (m *mockRequest) SetTimeOut(t int) HTTPRequest {
	args := m.Called(t)
	return args.Get(0).(HTTPRequest)
}

func (m *mockRequest) GetTimeOut() time.Duration {
	args := m.Called()
	return args.Get(0).(time.Duration)
}
