package order

type NilHandling int

const (
	Native NilHandling = iota + 1
	NilsFirst
	NilsLast
)

type Order interface {
	Direction() Direction
	Property() string
	IsAscending() bool
	IsDescending() bool
	IsIgnoreCase() bool
	NilHandling() NilHandling
}

type Option interface {
	applyOption(*order)
}

type order struct {
	direction   Direction
	property    string
	ignoreCase  bool
	nilHandling NilHandling
}

func (o *order) Direction() Direction {
	return o.direction
}

func (o *order) Property() string {
	return o.property
}

func (o *order) IsAscending() bool {
	return o.direction == Ascending
}

func (o *order) IsDescending() bool {
	return o.direction == Descending
}

func (o *order) IsIgnoreCase() bool {
	return o.ignoreCase
}

func (o *order) NilHandling() NilHandling {
	return o.nilHandling
}

func By(property string, options ...Option) Order {
	o := &order{
		direction:   Default,
		property:    property,
		ignoreCase:  false,
		nilHandling: Native,
	}

	for _, option := range options {
		option.applyOption(o)
	}

	return o
}

type directionOption Direction

func (d directionOption) applyOption(o *order) {
	o.direction = Direction(d)
}

func WithDirection(direction Direction) Option {
	return directionOption(direction)
}

type nilHandlingOption NilHandling

func (n nilHandlingOption) applyOption(o *order) {
	o.nilHandling = NilHandling(n)
}

func WithNilHandling(nilHandling NilHandling) Option {
	return nilHandlingOption(nilHandling)
}

type ignoreCaseOption bool

func (i ignoreCaseOption) applyOption(o *order) {
	o.ignoreCase = bool(i)
}

func IgnoreCase() Option {
	return ignoreCaseOption(true)
}
