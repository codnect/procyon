package procyon

import context "github.com/procyon-projects/procyon-context"

type ApplicationRunner interface {
	OnApplicationRun(context context.Context, arguments ApplicationArguments)
}
