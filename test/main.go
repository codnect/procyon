package main

import "codnect.io/procyon"

func main() {
	err := procyon.New().Run()
	if err != nil {
		panic(err)
	}
}
