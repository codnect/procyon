package procyon

import (
	"codnect.io/procyon-core/component"
	"codnect.io/procyon-core/component/filter"
	"codnect.io/procyon-core/runtime"
	"context"
	"errors"
	"fmt"
	"os"
	goruntime "runtime"
	"strings"
	"time"
)

type Application struct {
	container component.Container
}

func New() *Application {
	return &Application{
		container: component.NewObjectContainer(),
	}
}

func (a *Application) Run(args ...string) error {
	startTime := time.Now()

	banner, err := a.resolveBanner()
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

	log.Info("Starting application using Go {}", goruntime.Version()[2:])
	log.Debug("Running with Procyon {}", Version)

	var environment runtime.Environment
	environment, err = a.prepareEnvironment(arguments)
	if err != nil {
		panic(err)
	}

	a.logProfileInfo(environment)

	err = a.prepareContext(environment, arguments)
	if err != nil {
		panic(err)
	}

	timeTakenToStartup := time.Now().Sub(startTime)
	log.Info("Started application in {} seconds", timeTakenToStartup.Seconds())

	return nil
}

func (a *Application) Exit() int {
	return 0
}

func (a *Application) resolveBanner() (runtime.Banner, error) {
	bannerPrinters := component.List(filter.ByTypeOf[runtime.Banner]())

	if len(bannerPrinters) > 1 {
		return nil, errors.New("banners cannot be distinguished because too many matching found")
	} else if len(bannerPrinters) == 1 {
		constructor := bannerPrinters[0].Definition().Constructor()
		banner, err := constructor.Invoke()

		if err != nil {
			return nil, fmt.Errorf("banner is not initialized, error: %e", err)
		}

		return banner[0].(runtime.Banner), nil
	}

	return newBannerPrinter(), nil
}

func (a *Application) prepareEnvironment(args *runtime.Arguments) (runtime.Environment, error) {
	environment := runtime.NewDefaultEnvironment()

	propertySources := environment.PropertySources()

	propertySources.AddLast(runtime.NewArgumentsSource(args))
	propertySources.AddLast(runtime.NewEnvironmentSource())

	err := a.configureEnvironment(nil, environment)
	if err != nil {
		return nil, err
	}

	return environment, nil
}

func (a *Application) prepareContext(environment runtime.Environment, arguments *runtime.Arguments) error {
	err := a.configureContext(nil)
	if err != nil {
		return err
	}

	singletons := a.container.Singletons()
	err = singletons.Register("procyonApplicationArguments", arguments)
	if err != nil {
		return err
	}

	return nil
}

func (a *Application) configureEnvironment(ctx context.Context, environment runtime.Environment) error {
	configurerList := a.container.ListObjects(ctx, filter.ByTypeOf[runtime.EnvironmentConfigurer]())

	for _, configurer := range configurerList {
		envConfigurer := configurer.(runtime.EnvironmentConfigurer)
		err := envConfigurer.ConfigureEnvironment(ctx, environment)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *Application) configureContext(ctx runtime.Context) error {
	configurerList := a.container.ListObjects(ctx, filter.ByTypeOf[runtime.ContextConfigurer]())

	for _, configurer := range configurerList {
		contextConfigurer := configurer.(runtime.ContextConfigurer)
		err := contextConfigurer.ConfigureContext(ctx)
		if err != nil {
			return err
		}
	}

	loader := newComponentLoader(a.container)
	err := loader.loadDefinitions(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (a *Application) logProfileInfo(environment runtime.Environment) {
	if log.IsInfoEnabled() {
		activeProfiles := environment.ActiveProfiles()

		if len(activeProfiles) == 0 {
			defaultProfiles := environment.DefaultProfiles()
			log.Info("No active profile, using default profile(s): {}", strings.Join(defaultProfiles, ","))
		} else {
			log.Info("The application is using the following profile(s): {}", strings.Join(activeProfiles, ","))
		}
	}
}
