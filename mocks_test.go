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

package procyon

import (
	"bytes"
	"context"
	stdio "io"

	"codnect.io/procyon/io"
	"codnect.io/procyon/runtime"
	"github.com/stretchr/testify/mock"
)

type AnyComponent struct {
}

type AnyMockWriter struct {
	mock.Mock
	buf bytes.Buffer
}

func (w *AnyMockWriter) Write(p []byte) (int, error) {
	args := w.Called(p)
	err := args.Error(1)

	if err == nil {
		w.buf.Write(p)
	}

	return args.Int(0), err
}

func (w *AnyMockWriter) String() string {
	return w.buf.String()
}

type AnyMockApplication struct {
	mock.Mock
}

func (a *AnyMockApplication) SetBannerPrinter(printer runtime.BannerPrinter) {
	a.Called(printer)
}

func (a *AnyMockApplication) ResourceResolver() io.ResourceResolver {
	result := a.Called()
	if result.Get(0) == nil {
		return nil
	}

	return result.Get(0).(io.ResourceResolver)
}

func (a *AnyMockApplication) Run(args ...string) error {
	result := a.Called(args)
	return result.Error(0)
}

type AnyMockResource struct {
	mock.Mock
}

func (a *AnyMockResource) Name() string {
	result := a.Called()
	if result.Get(0) == nil {
		return ""
	}

	return result.Get(0).(string)
}

func (a *AnyMockResource) Location() string {
	result := a.Called()
	if result.Get(0) == nil {
		return ""
	}

	return result.Get(0).(string)
}

func (a *AnyMockResource) Exists() bool {
	result := a.Called()
	return result.Bool(0)
}

func (a *AnyMockResource) Reader() (stdio.ReadCloser, error) {
	result := a.Called()
	if result.Get(0) == nil {
		return nil, result.Error(1)
	}

	return result.Get(0).(stdio.ReadCloser), result.Error(1)
}

type AnyMockResourceResolver struct {
	mock.Mock
}

func (a *AnyMockResourceResolver) Resolve(ctx context.Context, location string) (io.Resource, error) {
	result := a.Called(ctx, location)
	if result.Get(0) == nil {
		return nil, result.Error(1)
	}

	return result.Get(0).(io.Resource), nil
}

type AnyMockLifecycle struct {
	mock.Mock
}

func (l *AnyMockLifecycle) Start(ctx context.Context) error {
	results := l.Called(ctx)
	return results.Error(0)
}

func (l *AnyMockLifecycle) Stop(ctx context.Context) error {
	results := l.Called(ctx)
	return results.Error(0)
}

func (l *AnyMockLifecycle) IsRunning() bool {
	results := l.Called()
	return results.Bool(0)
}

type anyMockLifecycleManager struct {
	mock.Mock
}

func (a *anyMockLifecycleManager) Startup(ctx context.Context) error {
	result := a.Called(ctx)
	return result.Error(0)
}

func (a *anyMockLifecycleManager) Shutdown(ctx context.Context) error {
	result := a.Called(ctx)
	return result.Error(0)
}

func (a *anyMockLifecycleManager) IsRunning() bool {
	result := a.Called()
	return result.Bool(0)
}

type AnyMockBannerPrinter struct {
	mock.Mock
}

func (p *AnyMockBannerPrinter) Print(env runtime.Environment, writer stdio.Writer) error {
	result := p.Called(env, writer)
	return result.Error(0)
}

type AnyMockEnvironmentCustomizer struct {
	mock.Mock
}

func (c *AnyMockEnvironmentCustomizer) CustomizeEnvironment(env runtime.Environment, app runtime.Application) error {
	results := c.Called(env, app)
	return results.Error(0)
}

type AnyMockContextInitializer struct {
	mock.Mock
}

func (a *AnyMockContextInitializer) InitializeContext(ctx runtime.Context) error {
	results := a.Called(ctx)
	return results.Error(0)
}

type AnyMockCommandLineRunner struct {
	mock.Mock
}

func newAnyMockCommandLinerRunner(anyComponent AnyComponent) *AnyMockCommandLineRunner {
	return &AnyMockCommandLineRunner{}
}

func (r *AnyMockCommandLineRunner) Run(ctx runtime.Context, args *runtime.Args) error {
	results := r.Called(ctx, args)
	return results.Error(0)
}

type AnyMockServerApp struct {
	mock.Mock
}

func (a *AnyMockServerApp) Start(ctx context.Context) error {
	results := a.Called(ctx)
	return results.Error(0)
}

func (a *AnyMockServerApp) Stop(ctx context.Context) error {
	results := a.Called(ctx)
	return results.Error(0)
}

func (a *AnyMockServerApp) Port() int {
	results := a.Called()
	return results.Int(0)
}
