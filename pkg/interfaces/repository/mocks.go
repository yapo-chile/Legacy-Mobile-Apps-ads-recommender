package repository

import (
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
