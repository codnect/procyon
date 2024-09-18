package runtime

import "context"

// CommandLineRunner interface allows to create a command-line application.
type CommandLineRunner interface {
	// Run method runs the command-line application with the given arguments.
	Run(ctx context.Context, args *Arguments) error
}
