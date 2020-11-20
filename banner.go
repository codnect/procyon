package procyon

import (
	"fmt"
)

var bannerText = []string{"",
	"   ___",
	"  / _ \\  _ __   ___    ___  _   _   ___   _ __",
	" / /_)/ | '__| / _ \\  / __|| | | | / _ \\ | '_ \\",
	"/ ___/  | |   | (_) || (__ | |_| || (_) || | | |",
	"\\/     |_|     \\___/  \\___| \\__, | \\___/ |_| |_|",
	"                            |___/",
}

func printBanner() {
	for _, line := range bannerText {
		fmt.Println(line)
	}
}
