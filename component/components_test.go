// Copyright 2026 Codnect
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

import "context"

type AnyComponent interface {
	AnyMethod()
}

type AnySimpleComponent struct {
}

func NewAnySimpleComponent() AnySimpleComponent {
	return AnySimpleComponent{}
}

type AnyPointerComponent struct {
}

func NewAnyPointerComponent() *AnyPointerComponent {
	return &AnyPointerComponent{}
}

func (a *AnyPointerComponent) AnyMethod() {}

type AnyDependentComponent struct {
	dependency AnySimpleComponent
}

func NewAnyDependentComponent(dependency AnySimpleComponent) *AnyDependentComponent {
	return &AnyDependentComponent{
		dependency: dependency,
	}
}

func (a *AnyDependentComponent) AnyMethod() {}

type AnyInitializableComponent struct {
	initError error
}

func NewAnyInitializableComponent() *AnyInitializableComponent {
	return &AnyInitializableComponent{}
}

func (a *AnyInitializableComponent) AnyMethod() {}

func (a *AnyInitializableComponent) Init(ctx context.Context) error {
	return a.initError
}

type AnyIndexedComponent struct {
	components []AnyComponent
}

func NewAnyIndexedComponent(components []AnyComponent) *AnyIndexedComponent {
	return &AnyIndexedComponent{
		components: components,
	}
}

type AnyDisposableComponent struct {
	disposeError error
}

func NewAnyDisposableComponent() *AnyDisposableComponent {
	return &AnyDisposableComponent{}
}

func (a *AnyDisposableComponent) Dispose() error {
	return a.disposeError
}
