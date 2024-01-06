package main

import (
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

var appArgs = []string{"run", "main.go"}

var runCmd = &cobra.Command{
	Use:   "run [flags] [--] [args]",
	Short: "Run a Procyon Application",
	Long:  `The run command lets you run a Procyon Application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !checkIfProjectIsAlreadyInitialized() {
			//color.Blue("Please init command to initialize a project first.")
			return nil
		}

		err := checkIfGoInstalled()

		if err != nil {
			return err
		}

		err = runApplication(args)

		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func runApplication(args []string) error {
	appArgs = append(appArgs, args...)
	cmd := exec.Command("go", appArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
