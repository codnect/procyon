package cfg

import "context"

type Loader interface {
	IsLoadable(resource Resource) bool
	LoadData(ctx context.Context, resource Resource) (*Data, error)
}
