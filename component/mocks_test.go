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
	"reflect"

	"github.com/stretchr/testify/mock"
)

type AnyMockScope struct {
	mock.Mock

	useFactory bool
}

func (a *AnyMockScope) Resolve(ctx context.Context, name string, fn FactoryFunc) (any, error) {
	if a.useFactory {
		return fn(ctx)
	}

	result := a.Called(ctx, name, fn)

	if result.Get(0) == nil {
		return nil, result.Error(1)
	}

	return result.Get(0).(any), result.Error(1)
}

func (a *AnyMockScope) Remove(ctx context.Context, name string) error {
	result := a.Called(ctx, name)
	return result.Error(0)
}

type AnyMockBeforeInitProcessor struct {
	mock.Mock
}

func (a *AnyMockBeforeInitProcessor) ProcessBeforeInit(ctx context.Context, instance any) (any, error) {
	result := a.Called(ctx, instance)

	if result.Get(0) == nil {
		return nil, result.Error(1)
	}

	return result.Get(0).(any), result.Error(1)
}

type AnyMockAfterInitProcessor struct {
	mock.Mock
}

func (a *AnyMockAfterInitProcessor) ProcessAfterInit(ctx context.Context, instance any) (any, error) {
	result := a.Called(ctx, instance)

	if result.Get(0) == nil {
		return nil, result.Error(1)
	}

	return result.Get(0).(any), result.Error(1)
}

type AnyMockContainer struct {
	mock.Mock
}

func (a *AnyMockContainer) RegisterDefinition(def *Definition) error {
	result := a.Called(def)
	return result.Error(0)
}

func (a *AnyMockContainer) UnregisterDefinition(name string) error {
	result := a.Called(name)
	return result.Error(0)
}

func (a *AnyMockContainer) Definition(name string) (*Definition, bool) {
	result := a.Called(name)
	def := result.Get(0)
	if def == nil {
		return nil, result.Bool(1)
	}

	return def.(*Definition), result.Bool(1)
}

func (a *AnyMockContainer) ContainsDefinition(name string) bool {
	result := a.Called(name)
	return result.Bool(0)
}

func (a *AnyMockContainer) Definitions() []*Definition {
	result := a.Called()
	items := result.Get(0)

	if items == nil {
		return nil
	}

	return items.([]*Definition)
}

func (a *AnyMockContainer) DefinitionsOf(typ reflect.Type) []*Definition {
	result := a.Called(typ)
	items := result.Get(0)

	if items == nil {
		return nil
	}

	return items.([]*Definition)
}

func (a *AnyMockContainer) DefinitionNames() []string {
	result := a.Called()
	items := result.Get(0)

	if items == nil {
		return nil
	}

	return items.([]string)
}

func (a *AnyMockContainer) DefinitionNamesOf(typ reflect.Type) []string {
	result := a.Called(typ)
	items := result.Get(0)

	if items == nil {
		return nil
	}

	return items.([]string)
}

func (a *AnyMockContainer) RegisterSingleton(name string, instance any) error {
	result := a.Called(name, instance)
	return result.Error(0)
}

func (a *AnyMockContainer) ContainsSingleton(name string) bool {
	result := a.Called(name)
	return result.Bool(0)
}

func (a *AnyMockContainer) Singleton(name string) (any, bool) {
	result := a.Called(name)
	return result.Get(0), result.Bool(1)
}

func (a *AnyMockContainer) RemoveSingleton(name string) error {
	result := a.Called(name)
	return result.Error(0)
}

func (a *AnyMockContainer) DestroySingletons() {
	a.Called()
}

func (a *AnyMockContainer) SingletonNames() []string {
	result := a.Called()
	items := result.Get(0)

	if items == nil {
		return nil
	}

	return items.([]string)
}

func (a *AnyMockContainer) CanResolve(name string) bool {
	result := a.Called(name)
	return result.Bool(0)
}

func (a *AnyMockContainer) CanResolveType(typ reflect.Type) bool {
	result := a.Called(typ)
	return result.Bool(0)
}

func (a *AnyMockContainer) Resolve(ctx context.Context, name string) (any, error) {
	result := a.Called(ctx, name)
	return result.Get(0), result.Error(1)
}

func (a *AnyMockContainer) ResolveType(ctx context.Context, typ reflect.Type) (any, error) {
	result := a.Called(ctx, typ)
	return result.Get(0), result.Error(1)
}

func (a *AnyMockContainer) ResolveAs(ctx context.Context, name string, typ reflect.Type) (any, error) {
	result := a.Called(ctx, name, typ)
	return result.Get(0), result.Error(1)
}

func (a *AnyMockContainer) ResolveAll(ctx context.Context, typ reflect.Type) ([]any, error) {
	result := a.Called(ctx, typ)
	items := result.Get(0)
	if items == nil {
		return nil, result.Error(1)
	}

	return items.([]any), result.Error(1)
}

func (a *AnyMockContainer) RegisterDependency(typ reflect.Type, val any) error {
	result := a.Called(typ, val)
	return result.Error(0)
}

func (a *AnyMockContainer) RegisterScope(name string, scope Scope) error {
	result := a.Called(name, scope)
	return result.Error(0)
}

func (a *AnyMockContainer) Scope(name string) (Scope, bool) {
	result := a.Called(name)
	scope := result.Get(0)
	if scope == nil {
		return nil, result.Bool(1)
	}

	return scope.(Scope), result.Bool(1)
}

func (a *AnyMockContainer) UseBeforeInitProcessor(processor BeforeInitProcessor) error {
	result := a.Called(processor)
	return result.Error(0)
}

func (a *AnyMockContainer) UseAfterInitProcessor(processor AfterInitProcessor) error {
	result := a.Called(processor)
	return result.Error(0)
}
