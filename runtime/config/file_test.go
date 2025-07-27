package config

import (
	"github.com/stretchr/testify/mock"
	"os"
)

type MockFileSystem struct {
	mock.Mock
}

func (m *MockFileSystem) Stat(name string) (os.FileInfo, error) {
	result := m.Called(name)

	if result.Get(0) == nil {
		return nil, result.Error(1)
	}

	return result.Get(0).(os.FileInfo), result.Error(1)
}

func (m *MockFileSystem) Open(name string) (*os.File, error) {
	result := m.Called(name)

	if result.Get(0) == nil {
		return nil, result.Error(1)
	}

	return result.Get(0).(*os.File), result.Error(1)
}
