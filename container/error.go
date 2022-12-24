package container

type notFoundError struct {
	ErrorString string
}

func (ne *notFoundError) Error() string {
	return ne.ErrorString
}
