// Copyright 2025 Codnect
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package runtime

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseArgs(t *testing.T) {
	testCases := []struct {
		name              string
		args              []string
		wantErr           error
		wantOptionArgs    map[string][]string
		wantNonOptionArgs []string
	}{
		{
			name:    "wrong arg format",
			args:    []string{"--argKey3"},
			wantErr: errors.New("wrong argument format '--argKey3'"),
		},
		{
			name: "valid args",
			args: []string{"argKey1", "-argKey2", "--argKey3=arg3Val", "--argKey3=arg3AnotherVal"},
			wantOptionArgs: map[string][]string{
				"argKey3": {
					"arg3Val",
					"arg3AnotherVal",
				},
			},
			wantNonOptionArgs: []string{"argKey1", "-argKey2"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given

			// when
			args, err := ParseArgs(tc.args)

			// then
			if tc.wantErr != nil {
				require.Nil(t, args)

				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantOptionArgs, args.optArgs)
			assert.Equal(t, tc.wantNonOptionArgs, args.nonOptsArgs)
		})
	}
}

func TestArgs_OptionNames(t *testing.T) {
	testCases := []struct {
		name            string
		args            []string
		wantOptionNames []string
	}{
		{
			name:            "option does not exist",
			args:            []string{"argKey1", "-argKey2", "--argKey3=arg3Val"},
			wantOptionNames: []string{"argKey3"},
		},
		{
			name:            "non-option args",
			args:            []string{"argKey1", "-argKey2"},
			wantOptionNames: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			args, err := ParseArgs(tc.args)
			require.NoError(t, err)

			// when
			optionNames := args.OptionNames()

			// then
			assert.Equal(t, tc.wantOptionNames, optionNames)
		})
	}
}

func TestArgs_ContainsOption(t *testing.T) {
	testCases := []struct {
		name       string
		args       []string
		optionName string
		wantResult bool
	}{
		{
			name:       "option does not exist",
			args:       []string{"argKey1", "-argKey2", "--argKey3=arg3Val"},
			optionName: "argKey4",
			wantResult: false,
		},
		{
			name:       "option exists",
			args:       []string{"--argKey3=arg3Val"},
			optionName: "argKey3",
			wantResult: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			args, err := ParseArgs(tc.args)
			require.NoError(t, err)

			// when
			result := args.ContainsOption(tc.optionName)

			// then
			assert.Equal(t, tc.wantResult, result)
		})
	}
}

func TestArgs_OptionValues(t *testing.T) {
	testCases := []struct {
		name       string
		args       []string
		optionName string
		wantValues []string
	}{
		{
			name:       "option does not exist",
			optionName: "argKey4",
			args:       []string{"argKey1", "-argKey2", "--argKey3=arg3Val"},
			wantValues: nil,
		},
		{
			name:       "option with one value",
			optionName: "argKey3",
			args:       []string{"argKey1", "-argKey2", "--argKey3=arg3Val"},
			wantValues: []string{"arg3Val"},
		},
		{
			name:       "option with multiple values",
			optionName: "argKey3",
			args:       []string{"argKey1", "-argKey2", "--argKey3=arg3Val", "--argKey3=arg3AnotherVal"},
			wantValues: []string{"arg3Val", "arg3AnotherVal"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			args, err := ParseArgs(tc.args)
			require.NoError(t, err)

			// when
			values := args.OptionValues(tc.optionName)

			// then
			assert.Equal(t, tc.wantValues, values)
		})
	}
}

func TestArgs_NonOptionArgs(t *testing.T) {
	// given
	args, err := ParseArgs([]string{"argKey1", "-argKey2", "--argKey3=arg3Val", "--argKey3=arg3AnotherVal"})
	require.NoError(t, err)

	// when
	values := args.NonOptionArgs()

	// then
	assert.Equal(t, []string{"argKey1", "-argKey2"}, values)
}

func TestNewArgsPropertySource(t *testing.T) {
	testCases := []struct {
		name      string
		args      *Args
		wantPanic error
	}{
		{
			name:      "nil args",
			wantPanic: errors.New("nil args"),
		},
		{
			name: "valid args",
			args: &Args{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given

			// when
			if tc.wantPanic != nil {
				require.PanicsWithValue(t, tc.wantPanic.Error(), func() {
					NewArgsPropertySource(tc.args)
				})
				return
			}

			argsPropertySource := NewArgsPropertySource(tc.args)

			// then
			require.NotNil(t, argsPropertySource)
		})
	}
}

func TestArgsPropertySource_Name(t *testing.T) {
	// given
	args, err := ParseArgs([]string{})
	require.NoError(t, err)

	argsPropSource := NewArgsPropertySource(args)

	// when
	name := argsPropSource.Name()

	// then
	assert.Equal(t, "commandLineArgs", name)
}

func TestArgsPropertySource_Underlying(t *testing.T) {
	// given
	args, err := ParseArgs([]string{})
	require.NoError(t, err)

	argsPropSource := NewArgsPropertySource(args)

	// when
	underlying := argsPropSource.Underlying()

	// then
	assert.NotNil(t, underlying)
}

func TestArgsPropertySource_ContainsProperty(t *testing.T) {
	testCases := []struct {
		name       string
		args       []string
		propName   string
		wantResult bool
	}{
		{
			name:       "option does not exist",
			args:       []string{"argKey1", "-argKey2", "--argKey3=arg3Val"},
			propName:   "argKey4",
			wantResult: false,
		},
		{
			name:       "non-option args",
			args:       []string{"argKey1", "-argKey2", "--argKey3=arg3Val"},
			propName:   NonOptionArgs,
			wantResult: true,
		},
		{
			name:       "no non-option args",
			args:       []string{"--argKey3=arg3Val"},
			propName:   NonOptionArgs,
			wantResult: false,
		},
		{
			name:       "option exists",
			args:       []string{"--argKey3=arg3Val"},
			propName:   "argKey3",
			wantResult: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			args, err := ParseArgs(tc.args)
			require.NoError(t, err)

			argsPropSource := NewArgsPropertySource(args)

			// when
			result := argsPropSource.ContainsProperty(tc.propName)

			// then
			assert.Equal(t, tc.wantResult, result)
		})
	}
}

func TestArgsPropertySource_Property(t *testing.T) {
	testCases := []struct {
		name       string
		args       []string
		propName   string
		wantExists bool
		wantValue  any
	}{
		{
			name:       "option does not exist",
			args:       []string{"argKey1", "-argKey2", "--argKey3=arg3Val"},
			propName:   "argKey4",
			wantExists: false,
		},
		{
			name:       "non-option args",
			args:       []string{"argKey1", "-argKey2", "--argKey3=arg3Val"},
			propName:   NonOptionArgs,
			wantExists: true,
			wantValue:  "argKey1,-argKey2",
		},
		{
			name:       "no non-option args",
			args:       []string{"--argKey3=arg3Val"},
			propName:   NonOptionArgs,
			wantExists: false,
		},
		{
			name:       "option exists",
			args:       []string{"--argKey3=arg3Val"},
			propName:   "argKey3",
			wantExists: true,
			wantValue:  "arg3Val",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			args, err := ParseArgs(tc.args)
			require.NoError(t, err)

			argsPropSource := NewArgsPropertySource(args)

			// when
			val, exists := argsPropSource.Property(tc.propName)

			// then
			assert.Equal(t, tc.wantExists, exists)
			assert.Equal(t, tc.wantValue, val)
		})
	}
}

func TestArgsPropertySource_PropertyOrDefault(t *testing.T) {
	testCases := []struct {
		name         string
		args         []string
		propName     string
		defaultValue string
		wantValue    any
	}{
		{
			name:         "option does not exist",
			args:         []string{"argKey1", "-argKey2", "--argKey3=arg3Val"},
			propName:     "argKey4",
			defaultValue: "anyDefaultValue",
			wantValue:    "anyDefaultValue",
		},
		{
			name:         "non-option args",
			args:         []string{"argKey1", "-argKey2", "--argKey3=arg3Val"},
			propName:     NonOptionArgs,
			defaultValue: "anyDefaultValue",
			wantValue:    "argKey1,-argKey2",
		},
		{
			name:         "no non-option args",
			args:         []string{"--argKey3=arg3Val"},
			propName:     NonOptionArgs,
			defaultValue: "anyDefaultValue",
			wantValue:    "anyDefaultValue",
		},
		{
			name:         "option exists",
			args:         []string{"--argKey3=arg3Val"},
			propName:     "argKey3",
			defaultValue: "anyDefaultValue",
			wantValue:    "arg3Val",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			args, err := ParseArgs(tc.args)
			require.NoError(t, err)

			argsPropSource := NewArgsPropertySource(args)

			// when
			val := argsPropSource.PropertyOrDefault(tc.propName, tc.defaultValue)

			// then
			assert.Equal(t, tc.wantValue, val)
		})
	}
}

func TestArgsPropertySource_PropertyNames(t *testing.T) {
	testCases := []struct {
		name          string
		args          []string
		wantPropNames []string
	}{
		{
			name:          "option does not exist",
			args:          []string{"argKey1", "-argKey2", "--argKey3=arg3Val"},
			wantPropNames: []string{"argKey3"},
		},
		{
			name:          "non-option args",
			args:          []string{"argKey1", "-argKey2"},
			wantPropNames: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			args, err := ParseArgs(tc.args)
			require.NoError(t, err)

			argsPropSource := NewArgsPropertySource(args)

			// when
			propNames := argsPropSource.PropertyNames()

			// then
			assert.Equal(t, tc.wantPropNames, propNames)
		})
	}
}

func TestArgsPropertySource_OptionValues(t *testing.T) {
	testCases := []struct {
		name       string
		args       []string
		propName   string
		wantValues []string
	}{
		{
			name:       "option does not exist",
			propName:   "argKey4",
			args:       []string{"argKey1", "-argKey2", "--argKey3=arg3Val"},
			wantValues: nil,
		},
		{
			name:       "option with one value",
			propName:   "argKey3",
			args:       []string{"argKey1", "-argKey2", "--argKey3=arg3Val"},
			wantValues: []string{"arg3Val"},
		},
		{
			name:       "option with multiple values",
			propName:   "argKey3",
			args:       []string{"argKey1", "-argKey2", "--argKey3=arg3Val", "--argKey3=arg3AnotherVal"},
			wantValues: []string{"arg3Val", "arg3AnotherVal"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			args, err := ParseArgs(tc.args)
			require.NoError(t, err)

			argsPropSource := NewArgsPropertySource(args)

			// when
			values := argsPropSource.OptionValues(tc.propName)

			// then
			assert.Equal(t, tc.wantValues, values)
		})
	}
}

func TestArgsPropertySource_NonOptionArgs(t *testing.T) {
	// given
	args, err := ParseArgs([]string{"argKey1", "-argKey2", "--argKey3=arg3Val", "--argKey3=arg3AnotherVal"})
	require.NoError(t, err)

	argsPropSource := NewArgsPropertySource(args)

	// when
	values := argsPropSource.NonOptionArgs()

	// then
	assert.Equal(t, []string{"argKey1", "-argKey2"}, values)
}
