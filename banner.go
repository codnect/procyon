package procyon

import (
	"codnect.io/logy"
	"codnect.io/procyon-core/component"
	"codnect.io/procyon-core/component/filter"
	"codnect.io/procyon-core/runtime"
	"errors"
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
)

type bannerPrinter struct {
}

func newBannerPrinter() *bannerPrinter {
	return &bannerPrinter{}
}

func (p *bannerPrinter) PrintBanner(w io.Writer) error {
	for _, line := range bannerText {
		if logy.SupportsColor() {
			w.Write([]byte(fmt.Sprintf("\u001B[34;1m%s\u001B[0m", line)))
		} else {
			w.Write([]byte(line))
		}
	}

	if logy.SupportsColor() {
		w.Write([]byte(fmt.Sprintf("\u001B[93m%24s%s)\u001B[0m\n", "(", Version)))
	} else {
		w.Write([]byte(fmt.Sprintf("%24s%s)\n", "(", Version)))
	}

	return nil
}

func resolveBanner() (runtime.Banner, error) {
	bannerPrinters := component.List(filter.ByTypeOf[runtime.Banner]())

	if len(bannerPrinters) > 1 {
		return nil, errors.New("banners cannot be distinguished because too many matching found")
	} else if len(bannerPrinters) == 1 {
		constructor := bannerPrinters[0].Definition().Constructor()
		banner, err := constructor.Invoke()

		if err != nil {
			return nil, fmt.Errorf("banner is not initialized, error: %e", err)
		}

		return banner[0].(runtime.Banner), nil
	}

	return newBannerPrinter(), nil
}
