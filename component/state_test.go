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

package component

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithCreationState(t *testing.T) {
	testCases := []struct {
		name      string
		parentCtx context.Context
	}{
		{
			name:      "nil context",
			parentCtx: nil,
		},
		{
			name:      "any context",
			parentCtx: context.Background(),
		},
		{
			name:      "context with creation state",
			parentCtx: context.WithValue(context.Background(), ctxCreationStateContextKey, newCreationState()),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given

			// when
			result := withCreationState(tc.parentCtx)

			// then
			require.NotNil(t, result)

			state := creationStateFromContext(result)
			assert.NotNil(t, state)
		})
	}
}

func TestCreationStateFromContext(t *testing.T) {
	// given
	state := newCreationState()
	ctx := context.WithValue(context.Background(), ctxCreationStateContextKey, state)

	// when
	result := creationStateFromContext(ctx)

	// then
	assert.Equal(t, state, result)
}

func TestCreationState_PutToPreparation(t *testing.T) {
	testCases := []struct {
		name         string
		preCondition func(state *creationState)
		instanceName string

		wantErr error
	}{
		{
			name: "circular dependency",
			preCondition: func(state *creationState) {
				state.currentlyInCreation["anyInstanceName"] = struct{}{}
			},
			instanceName: "anyInstanceName",
			wantErr:      errors.New("instance is in preparation, maybe it has got circular dependency cycle"),
		},
		{
			name:         "no circular dependency cycle",
			instanceName: "anyInstanceName",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			state := newCreationState()

			if tc.preCondition != nil {
				tc.preCondition(state)
			}

			// when
			err := state.putToPreparation(tc.instanceName)

			// then
			if tc.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestCreationState_RemoveFromPreparation(t *testing.T) {
	// given
	state := newCreationState()
	state.currentlyInCreation["anyInstanceName"] = struct{}{}

	// when
	state.removeFromPreparation("anyInstanceName")

	// then
	assert.NotContains(t, state.currentlyInCreation, "anyInstanceName")
}
