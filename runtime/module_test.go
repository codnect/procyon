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

package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type AnyModule struct {
}

func (a AnyModule) Init() {
}

type AnyPanicModule struct {
}

func (a AnyPanicModule) Init() {
	panic("any error")
}

func TestUse_NoPanic(t *testing.T) {
	// given

	// when
	assert.NotPanics(t, func() {
		Use[AnyModule]()
	})

	// then
}

func TestUse_Panic(t *testing.T) {
	// given

	// when
	assert.PanicsWithError(t, "failed to initialize module: any error", func() {
		Use[AnyPanicModule]()
	})

	// then
}
