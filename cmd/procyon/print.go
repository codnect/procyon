package main

import (
	"codnect.io/logy"
	"fmt"
)

type consoleColor int

const (
	defaultConsoleColor consoleColor = iota + 1
	greenConsoleColor
	redConsoleColor
	yellowConsoleColor
	blueConsoleColor
)

func (c consoleColor) format(message string) string {
	switch c {
	case defaultConsoleColor:
		return fmt.Sprintf("\u001B[0;1m%s\u001B[0m", message)
	case greenConsoleColor:
		return fmt.Sprintf("\u001B[32;1m%s\u001B[0m", message)
	case redConsoleColor:
		return fmt.Sprintf("\u001B[31;1m%s\u001B[0m", message)
	case yellowConsoleColor:
		return fmt.Sprintf("\u001B[33;1m%s\u001B[0m", message)
	case blueConsoleColor:
		return fmt.Sprintf("\u001B[34;1m%s\u001B[0m", message)
	}

	return message
}

func (c consoleColor) print(message string) {
	if logy.SupportsColor() {
		fmt.Print(c.format(message))
	} else {
		fmt.Print(message)
	}
}

func (c consoleColor) println(message string) {
	c.print(fmt.Sprintf("%s\n", message))
}
