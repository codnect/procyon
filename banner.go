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

package procyon

import (
	"fmt"
	"io"
)

var (
	bannerText = []string{
		"   ___  _______ ______ _____  ___\n",
		"  / _ \\/ __/ _ / __/ // / _ \\/ _ \\\n",
		" / .__/_/  \\___\\__/\\_, /\\___/_//_/\n",
		"/_/               /___/\n",
	}
	versionFormat = "%24s%s)\n"
)

type BannerPrinter struct {
}

func NewBannerPrinter() *BannerPrinter {
	return &BannerPrinter{}
}

func (b *BannerPrinter) PrintBanner(w io.Writer) error {
	for _, line := range bannerText {
		_, err := w.Write([]byte(line))
		if err != nil {
			return err
		}
	}

	_, err := w.Write([]byte(fmt.Sprintf(versionFormat, "(", Version)))
	return err
}
