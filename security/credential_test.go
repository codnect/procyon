// Copyright 2026 Codnect
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

package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsernamePasswordCredential_Identity(t *testing.T) {
	// given
	credential := UsernamePasswordCredential{
		Username: "anyUsername",
		Password: "anyPassword",
	}

	// when
	identity := credential.Identity()

	// then
	assert.Equal(t, "anyUsername", identity)
}

func TestUsernamePasswordCredential_Secret(t *testing.T) {
	// given
	credential := UsernamePasswordCredential{
		Username: "anyUsername",
		Password: "anyPassword",
	}

	// when
	secret := credential.Secret()

	// then
	assert.Equal(t, "anyPassword", secret)
}
