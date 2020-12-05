package procyon

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type simpleCommandLinePropertySourceMock struct {
	mock.Mock
}

func (propertySource simpleCommandLinePropertySourceMock) ContainsOption(name string) bool {
	result := propertySource.Called(name)
	return result.Bool(0)
}

func (propertySource simpleCommandLinePropertySourceMock) GetOptionValues(name string) []string {
	result := propertySource.Called(name)
	val := result.Get(0)
	if val == nil {
		return nil
	}
	return result.Get(0).([]string)
}

func (propertySource simpleCommandLinePropertySourceMock) GetNonOptionArgs() []string {
	result := propertySource.Called()
	val := result.Get(0)
	if val == nil {
		return nil
	}
	return result.Get(0).([]string)
}

func (propertySource simpleCommandLinePropertySourceMock) GetName() string {
	result := propertySource.Called()
	return result.String(0)
}

func (propertySource simpleCommandLinePropertySourceMock) GetSource() interface{} {
	result := propertySource.Called()
	return result.Get(0)
}

func (propertySource simpleCommandLinePropertySourceMock) GetProperty(name string) interface{} {
	result := propertySource.Called(name)
	return result.Get(0)
}

func (propertySource simpleCommandLinePropertySourceMock) ContainsProperty(name string) bool {
	result := propertySource.Called(name)
	return result.Bool(0)
}

func (propertySource simpleCommandLinePropertySourceMock) GetPropertyNames() []string {
	result := propertySource.Called()
	val := result.Get(0)
	if val == nil {
		return nil
	}
	return result.Get(0).([]string)
}

func TestDefaultApplicationArguments(t *testing.T) {
	var args []string
	appArguments := getApplicationArguments(args)
	commandLinePropertySourceMock := &simpleCommandLinePropertySourceMock{}
	appArguments.source = commandLinePropertySourceMock

	commandLinePropertySourceMock.On("ContainsOption", "exist-argument").Return(true)
	commandLinePropertySourceMock.On("ContainsOption", "non-exist-argument").Return(false)
	assert.True(t, appArguments.ContainsOption("exist-argument"))
	assert.False(t, appArguments.ContainsOption("non-exist-argument"))
	commandLinePropertySourceMock.AssertExpectations(t)

	var arguments = make([]string, 0)
	commandLinePropertySourceMock.On("GetOptionValues", "exist-argument").Return(arguments)
	commandLinePropertySourceMock.On("GetOptionValues", "non-exist-argument").Return(nil)
	assert.NotNil(t, appArguments.GetOptionValues("exist-argument"))
	assert.Nil(t, appArguments.GetOptionValues("non-exist-argument"))
	commandLinePropertySourceMock.AssertExpectations(t)

	var nonOptionArgs = make([]string, 0)
	commandLinePropertySourceMock.On("GetNonOptionArgs").Return(nonOptionArgs)
	assert.Equal(t, nonOptionArgs, appArguments.GetNonOptionArgs())
	commandLinePropertySourceMock.AssertExpectations(t)

	var optionNames = make([]string, 0)
	commandLinePropertySourceMock.On("GetPropertyNames").Return(optionNames)
	assert.Equal(t, optionNames, appArguments.GetOptionNames())
	commandLinePropertySourceMock.AssertExpectations(t)

	var sourceArgs = make([]string, 0)
	appArguments.args = sourceArgs
	assert.Equal(t, sourceArgs, appArguments.GetSourceArgs())
}
