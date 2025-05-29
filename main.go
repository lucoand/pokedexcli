package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
)

type cliCommand struct {
	name string
	description string
	callback func() error
}

var commands = map[string]cliCommand{}

func cleanInput(text string) []string {
	text = strings.ToLower(text)
	text = strings.TrimSpace(text)
	retval := strings.Fields(text)
	return retval
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:\n")
	for _, cmd := range commands {
		fmt.Println(cmd.name + ": " + cmd.description)
	}
	return nil
}

func init() {
	commands["help"] =  cliCommand{
		name: "help",
		description: "Displays a help message",
		callback: commandHelp,
	}
	commands["exit"] =  cliCommand{
		name: "exit",
		description: "Exit the Pokedex",
		callback: commandExit,
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for true {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		args := cleanInput(scanner.Text())
		if cmd, exists := commands[args[0]]; exists {
			cmd.callback()
		} else {
			fmt.Println("Unknown command")
		}
	}
}
