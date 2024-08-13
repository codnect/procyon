package procyon

import (
	"codnect.io/logy"
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
