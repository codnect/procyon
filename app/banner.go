package app

import (
	"github.com/fatih/color"
	"github.com/procyon-projects/procyon/app/env"
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

func defaultBannerPrinter() *bannerPrinter {
	return &bannerPrinter{}
}

func (p *bannerPrinter) PrintBanner(environment env.Environment, w io.Writer) error {
	blue := color.New(color.FgHiBlue).Add(color.Bold)

	for _, line := range bannerText {
		_, err := blue.Fprint(w, line)
		if err != nil {
			return nil
		}
	}

	yellow := color.New(color.FgHiYellow)
	_, err := yellow.Fprintf(w, "%24s%s)", "(", Version)
	if err != nil {
		return nil
	}

	return nil
}
