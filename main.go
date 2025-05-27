package main

import (
	"fmt"
	"strings"
)

func cleanInput(text string) []string {
	text = strings.ToLower(text)
	text = strings.TrimSpace(text)
	retval := strings.Fields(text)
	return retval
}

func main() {
	fmt.Println("Hello, World!")
}
