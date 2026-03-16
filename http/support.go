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

package http

import "net/http"

// Cookie represents an HTTP cookie.
type Cookie = http.Cookie

// SameSite represents the SameSite attribute of the cookie.
type SameSite = http.SameSite

const (
	// SameSiteDefaultMode represents the default mode of the SameSite attribute.
	SameSiteDefaultMode SameSite = iota + 1
	// SameSiteLaxMode represents the lax mode of the SameSite attribute.
	SameSiteLaxMode
	// SameSiteStrictMode represents the strict mode of the SameSite attribute.
	SameSiteStrictMode
	// SameSiteNoneMode represents the none mode of the SameSite attribute.
	SameSiteNoneMode
)

// A Header represents the key-value pairs in an HTTP header.
type Header = http.Header
