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

import (
	"errors"
)

var (
	// ErrNotFound is returned when no matching component or instance is found.
	ErrNotFound = errors.New("not found")
	// ErrTypeMismatch is returned when a resolved instance cannot be assigned or converted to the requested type.
	ErrTypeMismatch = errors.New("type mismatch")
	// ErrAmbiguousMatch is returned when multiple candidates match and a single result cannot be determined.
	ErrAmbiguousMatch = errors.New("ambiguous match")
)
