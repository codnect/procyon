package app

import (
	"github.com/gookit/color"
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
	blue := color.New(color.FgBlue, color.Bold)

	for _, line := range bannerText {
		color.Fprintf(w, blue.Sprintf(line))
	}

	yellow := color.New(color.FgLightYellow)
	color.Fprintf(w, yellow.Sprintf("%24s%s)\n", "(", Version))

	return nil
}
