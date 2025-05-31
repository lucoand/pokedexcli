package main

import (
	"bufio"
	"fmt"
	"github.com/lucoand/pokedexcli/internal/pokeapi"
	"github.com/lucoand/pokedexcli/internal/pokecache"
	"os"
	"strings"
	"time"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*pokeapi.Config, string) error
}

var commands = map[string]cliCommand{}
var config = pokeapi.Config{}
var interval, _ = time.ParseDuration("5s")
var cache = pokecache.NewCache(interval)

func cleanInput(text string) []string {
	text = strings.ToLower(text)
	text = strings.TrimSpace(text)
	retval := strings.Fields(text)
	return retval
}

func commandExit(_ *pokeapi.Config, _ string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(_ *pokeapi.Config, _ string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Usage:\n\n")
	for _, cmd := range commands {
		fmt.Println(cmd.name + ": " + cmd.description)
	}
	return nil
}

func commandMap(config *pokeapi.Config, _ string) error {
	if config.Next == "" {
		fmt.Println("You're on the last page.")
		return nil
	}
	*config = pokeapi.GetMapData(config.Next, cache)
	return nil
}

func commandMapB(config *pokeapi.Config, _ string) error {
	if config.Previous == "" {
		fmt.Println("You're on the first page.")
		return nil
	}
	*config = pokeapi.GetMapData(config.Previous, cache)
	return nil
}

func commandExplore(_ *pokeapi.Config, location string) error {
	pokeapi.Explore(location, cache)
	return nil
}

func init() {
	commands["help"] = cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp,
	}
	commands["exit"] = cliCommand{
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	}
	commands["map"] = cliCommand{
		name:        "map",
		description: "Display the next 20 location areas",
		callback:    commandMap,
	}
	commands["mapb"] = cliCommand{
		name:        "mapb",
		description: "Display the previous 20 location areas",
		callback:    commandMapB,
	}
	commands["explore"] = cliCommand{
		name:        "explore",
		description: "Usage: \"explore <area>\". Lists pokemon available in target <area>",
		callback:    commandExplore,
	}
	config.Previous = ""
	config.Next = "https://pokeapi.co/api/v2/location-area/"
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for true {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		args := cleanInput(scanner.Text())
		if cmd, exists := commands[args[0]]; exists {
			arg := ""
			if len(args) > 1 {
				arg = args[1]
			}
			cmd.callback(&config, arg)
		} else {
			fmt.Println("Unknown command")
		}
	}
}
