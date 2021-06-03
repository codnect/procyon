/*
Copyright © 2021 Procyon Framework Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"errors"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strings"
)

// File Names and Paths
const (
	resourcesPath = "resources/"

	gitignoreFileName = ".gitignore"
	mainFileName      = "main.go"
	appYamlFileName   = "procyon.yaml"
	goModuleFileName  = "go.mod"
)

var projectFiles = []string{
	gitignoreFileName,
	mainFileName,
	resourcesPath + appYamlFileName,
	goModuleFileName,
}

// Contents
const (
	gitignoreFileContent = `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with 'go test -c'
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below to include it)
# vendor/`
	mainFileContent = `package main

import (
	"github.com/procyon-projects/procyon"
)

func main() {
	procyon.NewProcyonApplication().Run()
}
`
	appYamlFileContent = `procyon:
  application:
    name: "$APP_NAME"

server:
  port: 8080

logging:
  level: DEBUG
`
)

var (
	defaultColor = color.New(color.FgWhite)
	successColor = color.New(color.FgGreen)
	failColor    = color.New(color.FgRed)
)

var module string

var initCmd = &cobra.Command{
	Use:   "init [application-name]",
	Short: "Initialize a new project",
	Long:  `The init command lets you create a new project.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) == 0 {
			return errors.New("application name is required")
		}

		if len(args) > 1 {
			return errors.New("too many arguments")
		}

		err := checkIfGoInstalled()

		if err != nil {
			return err
		}

		if checkIfProjectIsAlreadyInitialized() {
			color.Blue("Project already initialized.")
			return nil
		}

		err = initializeProject(args[0])

		if err != nil {
			color.Red("Failed to initialize the project!")
		} else {
			color.Green("Completed successfully.")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVarP(&module, "module", "m", "", "Module Name (required)")

	err := initCmd.MarkFlagRequired("module")

	if err != nil {
		panic(err)
	}
}

func initializeProject(applicationName string) error {
	color.Yellow("Initializing project...")

	err := createGitIgnoreFile()

	if err != nil {
		return err
	}

	err = createApplicationPropertyFile(applicationName)

	if err != nil {
		return err
	}

	err = createMainFile()

	if err != nil {
		return err
	}

	err = initGoModAndGetDependencies()

	if err != nil {
		return err
	}

	return nil
}

func createGitIgnoreFile() error {
	defaultColor.Print(gitignoreFileName)

	err := createFile(gitignoreFileName, gitignoreFileContent)

	if err != nil {
		printFailStep()
	} else {
		printSuccessStep()
	}

	return err
}

func createApplicationPropertyFile(applicationName string) error {
	defaultColor.Print(resourcesPath + appYamlFileName)

	if !checkIfExist(resourcesPath) {
		err := os.Mkdir(resourcesPath, os.ModePerm)

		if err != nil {
			printFailStep()
			return err
		}
	}

	content := strings.ReplaceAll(appYamlFileContent, "$APP_NAME", applicationName)
	err := createFile(resourcesPath+appYamlFileName, content)

	if err != nil {
		printFailStep()
	} else {
		printSuccessStep()
	}

	return err
}

func createMainFile() error {
	defaultColor.Print(mainFileName)

	err := createFile(mainFileName, mainFileContent)

	if err != nil {
		printFailStep()
	} else {
		printSuccessStep()
	}

	return err
}

func initGoModAndGetDependencies() error {
	defaultColor.Print(goModuleFileName)

	err := exec.Command("go", "mod", "init", module).Run()

	if err != nil {
		printFailStep()
		return err
	}

	err = exec.Command("go", "get", "-t", "-v", "./...").Run()

	if err != nil {
		printFailStep()
	} else {
		printSuccessStep()
	}

	return err
}

func checkIfProjectIsAlreadyInitialized() bool {
	for _, projectFile := range projectFiles {

		if checkIfExist(projectFile) {
			return true
		}

	}

	return false
}

func printSuccessStep() {
	defaultColor.Print(" [")
	successColor.Print("ok")
	defaultColor.Print("]\n")
}

func printFailStep() {
	defaultColor.Print(" [")
	failColor.Print("failed")
	defaultColor.Print("]\n")
}
