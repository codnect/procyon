package component

import "math"

const (
	HighestPriority = math.MinInt
	LowestPriority  = math.MaxInt
)

type Prioritized interface {
	Priority() int
}
