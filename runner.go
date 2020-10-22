package procyon

type ApplicationRunner interface {
	Run(arguments ApplicationArguments)
}
