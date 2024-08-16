package procyon

import (
	"codnect.io/procyon-core/component"
	"codnect.io/procyon-core/component/filter"
	"codnect.io/procyon-core/runtime"
	"context"
	"os"
	"os/signal"
	goruntime "runtime"
	"syscall"
	"time"
)

type Application struct {
}

func New() *Application {
	return &Application{}
}

func (a *Application) Run(args ...string) error {
	startTime := time.Now()

	banner, err := resolveBanner()
	if err != nil {
		return err
	}

	err = banner.PrintBanner(os.Stdout)
	if err != nil {
		return err
	}

	var arguments *runtime.Arguments
	arguments, err = runtime.ParseArguments(args)
	if err != nil {
		return err
	}

	log.Info("Starting application using Go {} ({}/{})", goruntime.Version()[2:], goruntime.GOOS, goruntime.GOARCH)
	log.Info("Running with Procyon {}", Version)

	ctx := createContext(arguments)
	err = ctx.Start()

	if err != nil {
		return err
	}

	timeTakenToStartup := time.Now().Sub(startTime)
	log.Info("Started application in {} seconds", timeTakenToStartup.Seconds())

	err = callCommandLineRunners(ctx, arguments)

	if err != nil {
		return err
	}

	if isServerApp(ctx, ctx.Container()) {
		waitForShutdown(ctx)
	}

	return nil
}

func callCommandLineRunners(ctx runtime.Context, args *runtime.Arguments) error {
	runners := ctx.Container().ListObjects(ctx, filter.ByTypeOf[runtime.CommandLineRunner]())

	for _, runner := range runners {
		cmdRunner := runner.(runtime.CommandLineRunner)
		err := cmdRunner.Run(ctx, args)

		if err != nil {
			return err
		}
	}

	return nil
}

func isServerApp(ctx context.Context, container component.Container) bool {
	servers := container.ListObjects(ctx, filter.ByTypeOf[runtime.Server]())
	return len(servers) != 0
}

func waitForShutdown(ctx runtime.Context) {
	shutdownChannel := make(chan os.Signal, 1)
	signal.Notify(shutdownChannel, syscall.SIGINT, syscall.SIGTERM)

	shutdown := false

	for {
		select {
		case <-shutdownChannel:
			shutdown = true
			break
		case <-ctx.Done():
			shutdown = true
			break
		}

		if shutdown {
			break
		}
	}
}
