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

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"go.codnect.io/procyon"
)

type applicationSample string

const (
	sampleHTTP      applicationSample = "http"
	sampleCLIRunner applicationSample = "cli-runner"
)

// valid reports whether the application sample is supported.
func (s applicationSample) valid() bool {
	switch s {
	case "", sampleHTTP, sampleCLIRunner:
		return true
	default:
		return false
	}
}

const (
	goModTemplate = `module %s

go 1.27

require go.codnect.io/procyon %s
`

	procyonYAMLTemplate = `application:
  name: %s
`

	mainTemplate = `package main

import (
	"os"

	"go.codnect.io/procyon"
)

func main() {
	if err := procyon.Run(); err != nil {
		os.Exit(1)
	}
}
`

	httpMainTemplate = `package main

import (
	"os"

	"%s/internal/controller"

	"go.codnect.io/procyon"
	"go.codnect.io/procyon/component"
)

func init() {
	component.Register(controller.NewWelcomeController)
}

func main() {
	if err := procyon.Run(); err != nil {
		os.Exit(1)
	}
}
`

	welcomeControllerTemplate = `package controller

import "go.codnect.io/procyon/http"

// WelcomeController provides the default welcome endpoint.
type WelcomeController struct {
}

// NewWelcomeController creates a new WelcomeController.
func NewWelcomeController() *WelcomeController {
	return &WelcomeController{}
}

// MapEndpoints maps the controller endpoints.
func (c *WelcomeController) MapEndpoints(endpoints http.Endpoints) {
	endpoints.MapGet("/", http.Handle(c.hello))
}

// hello handles the welcome endpoint.
func (c *WelcomeController) hello(ctx *http.Context) error {
	_, err := ctx.Response().Writer().Write([]byte("Hello, World!"))
	return err
}
`

	cliRunnerMainTemplate = `package main

import (
	"os"

	"%s/internal/runner"

	"go.codnect.io/procyon"
	"go.codnect.io/procyon/component"
)

func init() {
	component.Register(runner.NewWelcomeRunner)
}

func main() {
	if err := procyon.Run(); err != nil {
		os.Exit(1)
	}
}
`

	welcomeRunnerTemplate = `package runner

import (
	"fmt"

	"go.codnect.io/procyon/runtime"
)

// WelcomeRunner prints a welcome message when the application starts.
type WelcomeRunner struct {
}

// NewWelcomeRunner creates a new WelcomeRunner.
func NewWelcomeRunner() *WelcomeRunner {
	return &WelcomeRunner{}
}

// Run executes the runner.
func (r *WelcomeRunner) Run(ctx runtime.Context, args *runtime.Args) error {
	fmt.Println("Hello, World!")
	return nil
}
`
)

// initCmd represents the command to initialize a new Procyon application.
var initCmd = &cobra.Command{
	Use:   "init <name>",
	Short: "Initialize a new Procyon application",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		sampleValue, err := cmd.Flags().GetString("sample")
		if err != nil {
			return err
		}

		sample := applicationSample(sampleValue)
		if !sample.valid() {
			return fmt.Errorf("unknown sample %q", sample)
		}

		return initializeApplication(name, sample)
	},
}

func init() {
	initCmd.Flags().String(
		"sample",
		"",
		"Sample code (http, cli-runner)",
	)
}

// initializeApplication initializes a new Procyon application.
func initializeApplication(name string, sample applicationSample) error {
	if _, err := os.Stat(name); err == nil {
		return fmt.Errorf("directory %q already exists", name)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("check application directory: %w", err)
	}

	if err := os.Mkdir(name, 0755); err != nil {
		return fmt.Errorf("create application directory: %w", err)
	}

	if err := writeFile(
		filepath.Join(name, "go.mod"),
		fmt.Sprintf(goModTemplate, name, procyon.Version),
	); err != nil {
		return err
	}

	if err := initializeResources(name); err != nil {
		return err
	}

	switch sample {
	case sampleHTTP:
		if err := initializeHTTPSample(name); err != nil {
			return err
		}

	case sampleCLIRunner:
		if err := initializeCLIRunnerSample(name); err != nil {
			return err
		}

	default:
		if err := initializeMain(name); err != nil {
			return err
		}
	}

	fmt.Printf("Procyon application %q initialized successfully\n", name)

	return nil
}

// initializeResources initializes the common application resources.
func initializeResources(name string) error {
	resourceDir := filepath.Join(name, "resources")

	if err := os.MkdirAll(resourceDir, 0755); err != nil {
		return fmt.Errorf("create resources directory: %w", err)
	}

	return writeFile(
		filepath.Join(resourceDir, "procyon.yaml"),
		fmt.Sprintf(procyonYAMLTemplate, name),
	)
}

// initializeMain initializes the application entry point.
func initializeMain(name string) error {
	commandDir := filepath.Join(name, "cmd", name)

	if err := os.MkdirAll(commandDir, 0755); err != nil {
		return fmt.Errorf("create command directory: %w", err)
	}

	return writeFile(
		filepath.Join(commandDir, "main.go"),
		mainTemplate,
	)
}

// initializeHTTPSample initializes the HTTP sample.
func initializeHTTPSample(name string) error {
	commandDir := filepath.Join(name, "cmd", name)
	controllerDir := filepath.Join(name, "internal", "controller")

	if err := os.MkdirAll(commandDir, 0755); err != nil {
		return fmt.Errorf("create command directory: %w", err)
	}

	if err := os.MkdirAll(controllerDir, 0755); err != nil {
		return fmt.Errorf("create controller directory: %w", err)
	}

	if err := writeFile(
		filepath.Join(commandDir, "main.go"),
		fmt.Sprintf(httpMainTemplate, name),
	); err != nil {
		return err
	}

	return writeFile(
		filepath.Join(controllerDir, "welcome.go"),
		welcomeControllerTemplate,
	)
}

// initializeCLIRunnerSample initializes the CLI runner sample.
func initializeCLIRunnerSample(name string) error {
	commandDir := filepath.Join(name, "cmd", name)
	runnerDir := filepath.Join(name, "internal", "runner")

	if err := os.MkdirAll(commandDir, 0755); err != nil {
		return fmt.Errorf("create command directory: %w", err)
	}

	if err := os.MkdirAll(runnerDir, 0755); err != nil {
		return fmt.Errorf("create runner directory: %w", err)
	}

	if err := writeFile(
		filepath.Join(commandDir, "main.go"),
		fmt.Sprintf(cliRunnerMainTemplate, name),
	); err != nil {
		return err
	}

	return writeFile(
		filepath.Join(runnerDir, "welcome.go"),
		welcomeRunnerTemplate,
	)
}

// writeFile writes content to the specified file.
func writeFile(path, content string) error {
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}

	return nil
}
