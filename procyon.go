package procyon

type AppType int

const (
	None AppType = -1
	Web  AppType = 0
)

type AppBuilder interface {
	Type(appType AppType) AppBuilder
	Run(args ...string)
}

func New() AppBuilder {
	return nil
}
