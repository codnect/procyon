package tx

type Propagation int

const (
	PropagationRequired Propagation = iota
	PropagationSupports
	PropagationMandatory
	PropagationNever
	PropagationNotSupported
	PropagationNested
	PropagationRequiredNew
)
