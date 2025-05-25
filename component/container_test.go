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
	"github.com/stretchr/testify/mock"
	"reflect"
)

type MockContainer struct {
	mock.Mock
}

func (m *MockContainer) RegisterDefinition(def *Definition) error {
	args := m.Called(def)
	return args.Error(0)
}

func (m *MockContainer) UnregisterDefinition(name string) error {
	args := m.Called(name)
	return args.Error(0)
}

func (m *MockContainer) Definition(name string) (*Definition, bool) {
	args := m.Called(name)
	return args.Get(0).(*Definition), args.Bool(1)
}

func (m *MockContainer) ContainsDefinition(name string) bool {
	args := m.Called(name)
	return args.Bool(0)
}

func (m *MockContainer) Definitions() []*Definition {
	args := m.Called()
	return args.Get(0).([]*Definition)
}

func (m *MockContainer) DefinitionsOf(typ reflect.Type) []*Definition {
	args := m.Called(typ)
	return args.Get(0).([]*Definition)
}

func (m *MockContainer) RegisterSingleton(name string, instance any) error {
	args := m.Called(name, instance)
	return args.Error(0)
}

func (m *MockContainer) ContainsSingleton(name string) bool {
	args := m.Called(name)
	return args.Bool(0)
}

func (m *MockContainer) Singleton(name string) (any, bool) {
	args := m.Called(name)
	return args.Get(0), args.Bool(1)
}

func (m *MockContainer) RemoveSingleton(name string) {
	m.Called(name)
}

func (m *MockContainer) CanResolve(ctx context.Context, name string) bool {
	args := m.Called(ctx, name)
	return args.Bool(0)
}

func (m *MockContainer) Resolve(ctx context.Context, typ reflect.Type) (any, error) {
	args := m.Called(ctx, typ)
	return args.Get(0), args.Error(1)
}

func (m *MockContainer) ResolveAll(ctx context.Context, typ reflect.Type) ([]any, error) {
	args := m.Called(ctx, typ)
	return args.Get(0).([]any), args.Error(1)
}

func (m *MockContainer) ResolveNamed(ctx context.Context, name string) (any, error) {
	args := m.Called(ctx, name)
	return args.Get(0), args.Error(1)
}

func (m *MockContainer) ResolveNamedType(ctx context.Context, typ reflect.Type, name string) (any, error) {
	args := m.Called(ctx, typ, name)
	return args.Get(0), args.Error(1)
}

func (m *MockContainer) Bind(typ reflect.Type, instance any) error {
	args := m.Called(typ, instance)
	return args.Error(0)
}

func (m *MockContainer) RegisterScope(name string, scope Scope) error {
	args := m.Called(name, scope)
	return args.Error(0)
}

func (m *MockContainer) Scope(name string) (Scope, bool) {
	args := m.Called(name)
	return args.Get(0).(Scope), args.Bool(1)
}

func (m *MockContainer) UsePreProcessor(initializer PreProcessor) error {
	args := m.Called(initializer)
	return args.Error(0)
}

func (m *MockContainer) UsePostProcessor(initializer PostProcessor) error {
	args := m.Called(initializer)
	return args.Error(0)
}
