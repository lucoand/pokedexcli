package pokeapi

import (
	"encoding/json"
	"fmt"
	"github.com/lucoand/pokedexcli/internal/pokecache"
	"io"
	"log"
	"net/http"
	"reflect"
)

type Config struct {
	Previous string
	Next     string
}

type LocationArea struct {
	Count    int `json:"count"`
	Next     any `json:"next"`
	Previous any `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type LocationInfo struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func printLocationPokemon(locationInfo LocationInfo) {
	if len(locationInfo.PokemonEncounters) < 1 {
		fmt.Println("No pokemon found!")
		return
	}
	fmt.Println("Found Pokemon:")
	for _, encounter := range locationInfo.PokemonEncounters {
		fmt.Printf(" - %v\n", encounter.Pokemon.Name)
	}
}

func printLocationArea(locationArea LocationArea) {
	for _, result := range locationArea.Results {
		fmt.Println(result.Name)
	}
}

func checkCache(url string, cache *pokecache.Cache, dest any) {
	if reflect.TypeOf(dest).Kind() != reflect.Ptr {
		log.Fatal("ERROR - checkCache(): dest must be a pointer.")
	}
	data, ok := cache.Get(url)
	if ok {
		err := json.Unmarshal(data, dest)
		if err != nil {
			log.Fatal("ERROR - Unmarshalling json:", err)
		}
	} else {
		response, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			log.Fatal("ERROR - Bad response:", response.StatusCode)
		}
		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal("ERROR - Unable to read response body:", err)
		}
		cache.Add(url, body)
		err = json.Unmarshal(body, &dest)
		if err != nil {
			log.Fatal("ERROR - Unmarshalling json (response):", err)
		}
	}
}

func GetLocationInfo(location string, cache *pokecache.Cache) LocationInfo {
	url := "https://pokeapi.co/api/v2/location-area/" + location
	var locationInfo LocationInfo
	checkCache(url, cache, &locationInfo)
	return locationInfo
}

func Explore(location string, cache *pokecache.Cache) {
	fmt.Printf("Exploring %v...\n", location)
	locationInfo := GetLocationInfo(location, cache)
	printLocationPokemon(locationInfo)
}

func GetMapData(url string, cache *pokecache.Cache) Config {
	var area LocationArea
	config := Config{}
	checkCache(url, cache, &area)

	printLocationArea(area)

	if area.Previous == nil {
		config.Previous = ""
	} else {
		config.Previous = area.Previous.(string)
	}
	if area.Next == nil {
		config.Next = ""
	} else {
		config.Next = area.Next.(string)
	}

	return config
}
