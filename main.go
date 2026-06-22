package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	//"encoding/json"
	"net/http"
	//"io"
)

type Location struct {
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

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func cleanInput(text string) []string {
	cleaned_output := make([]string, 0)
	text_clean := strings.Fields(text)
	for _, word := range text_clean {
		cleaned_output = append(cleaned_output, strings.ToLower(word))
	}
	return cleaned_output
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage: ")
	fmt.Println("help: Displays a help message")
	fmt.Println("exit: Exit the Pokedex")
	fmt.Println("map: Displays the next 20 locations available...")
	fmt.Println("mapb: Displays the previous 20 locations available....")
	return nil
}

func commandMap() error {
	location_area_url_base := `https://pokeapi.co/api/v2/location-area/`
	for i := 0; i < 20; i++ {
		location_url := location_area_url_base + strconv.Itoa(i+1)
		res, err := http.Get(location_url)
		if err != nil {
			fmt.Println(err)
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
		}
		var location Location
		err = json.Unmarshal(body, &location)
		fmt.Println(location.Name)
	}
	return nil
}

func commandMapb() error {
	fmt.Println("Do something in reverse.....")
	return nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	supportedCommands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays the next 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 locations",
			callback:    commandMapb,
		},
	}

	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())
		words := strings.Fields(input)
		if len(words) == 0 {
			continue
		}
		commandInput := strings.ToLower(words[0])
		command, ok := supportedCommands[commandInput]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		err := command.callback()
		if err != nil {
			fmt.Println("Error executing command:", err)
		}
	}
}
