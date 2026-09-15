// Copyright 2025 Codnect
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package component

import "math"

const (
	// HighestPrecedence represents the highest possible precedence.
	HighestPrecedence = math.MinInt
	// LowestPrecedence represents the lowest possible precedence.
	LowestPrecedence = math.MaxInt
)

// Ordered defines the precedence of a component. Lower values have higher precedence.
type Ordered interface {
	// Order returns the precedence value of the component.
	Order() int
}
