package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"pokedex/internal/pokecache"
	"strings"
	"time"
)

type config struct {
	Next          string
	Previous      string
	cache         *pokecache.Cache
	caughtPokemon map[string]Pokemon
}

type ShallowLocationList struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

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

type Pokemon struct {
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
	} `json:"abilities"`
	BaseExperience int `json:"base_experience"`
	Cries          struct {
		Latest string `json:"latest"`
		Legacy string `json:"legacy"`
	} `json:"cries"`
	Forms []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	GameIndices []struct {
		GameIndex int `json:"game_index"`
		Version   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version"`
	} `json:"game_indices"`
	Height    int `json:"height"`
	HeldItems []struct {
		Item struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"item"`
		VersionDetails []struct {
			Rarity  int `json:"rarity"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"held_items"`
	ID                     int    `json:"id"`
	IsDefault              bool   `json:"is_default"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Moves                  []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt  int `json:"level_learned_at"`
			MoveLearnMethod struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"move_learn_method"`
			Order        any `json:"order"`
			VersionGroup struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version_group"`
		} `json:"version_group_details"`
	} `json:"moves"`
	Name          string `json:"name"`
	Order         int    `json:"order"`
	PastAbilities []struct {
		Abilities []struct {
			Ability  any  `json:"ability"`
			IsHidden bool `json:"is_hidden"`
			Slot     int  `json:"slot"`
		} `json:"abilities"`
		Generation struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"generation"`
	} `json:"past_abilities"`
	PastStats []any `json:"past_stats"`
	PastTypes []any `json:"past_types"`
	Species   struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Sprites struct {
		BackDefault      string `json:"back_default"`
		BackFemale       string `json:"back_female"`
		BackShiny        string `json:"back_shiny"`
		BackShinyFemale  string `json:"back_shiny_female"`
		FrontDefault     string `json:"front_default"`
		FrontFemale      string `json:"front_female"`
		FrontShiny       string `json:"front_shiny"`
		FrontShinyFemale string `json:"front_shiny_female"`
		Other            struct {
			DreamWorld struct {
				FrontDefault string `json:"front_default"`
				FrontFemale  any    `json:"front_female"`
			} `json:"dream_world"`
			Home struct {
				FrontDefault     string `json:"front_default"`
				FrontFemale      string `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale string `json:"front_shiny_female"`
			} `json:"home"`
			OfficialArtwork struct {
				FrontDefault string `json:"front_default"`
				FrontShiny   string `json:"front_shiny"`
			} `json:"official-artwork"`
			Showdown struct {
				BackDefault      string `json:"back_default"`
				BackFemale       string `json:"back_female"`
				BackShiny        string `json:"back_shiny"`
				BackShinyFemale  any    `json:"back_shiny_female"`
				FrontDefault     string `json:"front_default"`
				FrontFemale      string `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale string `json:"front_shiny_female"`
			} `json:"showdown"`
		} `json:"other"`
		Versions struct {
			GenerationI struct {
				RedBlue struct {
					BackDefault      any `json:"back_default"`
					BackGray         any `json:"back_gray"`
					BackTransparent  any `json:"back_transparent"`
					FrontDefault     any `json:"front_default"`
					FrontGray        any `json:"front_gray"`
					FrontTransparent any `json:"front_transparent"`
				} `json:"red-blue"`
				Yellow struct {
					BackDefault      any `json:"back_default"`
					BackGray         any `json:"back_gray"`
					BackTransparent  any `json:"back_transparent"`
					FrontDefault     any `json:"front_default"`
					FrontGray        any `json:"front_gray"`
					FrontTransparent any `json:"front_transparent"`
				} `json:"yellow"`
			} `json:"generation-i"`
			GenerationIi struct {
				Crystal struct {
					BackDefault           any `json:"back_default"`
					BackShiny             any `json:"back_shiny"`
					BackShinyTransparent  any `json:"back_shiny_transparent"`
					BackTransparent       any `json:"back_transparent"`
					FrontDefault          any `json:"front_default"`
					FrontShiny            any `json:"front_shiny"`
					FrontShinyTransparent any `json:"front_shiny_transparent"`
					FrontTransparent      any `json:"front_transparent"`
				} `json:"crystal"`
				Gold struct {
					BackDefault      any `json:"back_default"`
					BackShiny        any `json:"back_shiny"`
					FrontDefault     any `json:"front_default"`
					FrontShiny       any `json:"front_shiny"`
					FrontTransparent any `json:"front_transparent"`
				} `json:"gold"`
				Silver struct {
					BackDefault      any `json:"back_default"`
					BackShiny        any `json:"back_shiny"`
					FrontDefault     any `json:"front_default"`
					FrontShiny       any `json:"front_shiny"`
					FrontTransparent any `json:"front_transparent"`
				} `json:"silver"`
			} `json:"generation-ii"`
			GenerationIii struct {
				Emerald struct {
					FrontDefault any `json:"front_default"`
					FrontShiny   any `json:"front_shiny"`
				} `json:"emerald"`
				FireredLeafgreen struct {
					BackDefault  any `json:"back_default"`
					BackShiny    any `json:"back_shiny"`
					FrontDefault any `json:"front_default"`
					FrontShiny   any `json:"front_shiny"`
				} `json:"firered-leafgreen"`
				RubySapphire struct {
					BackDefault  any `json:"back_default"`
					BackShiny    any `json:"back_shiny"`
					FrontDefault any `json:"front_default"`
					FrontShiny   any `json:"front_shiny"`
				} `json:"ruby-sapphire"`
			} `json:"generation-iii"`
			GenerationIv struct {
				DiamondPearl struct {
					BackDefault      string `json:"back_default"`
					BackFemale       string `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"diamond-pearl"`
				HeartgoldSoulsilver struct {
					BackDefault      string `json:"back_default"`
					BackFemale       string `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"heartgold-soulsilver"`
				Platinum struct {
					BackDefault      string `json:"back_default"`
					BackFemale       string `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"platinum"`
			} `json:"generation-iv"`
			GenerationIx struct {
				ScarletViolet struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  any    `json:"front_female"`
				} `json:"scarlet-violet"`
			} `json:"generation-ix"`
			GenerationV struct {
				BlackWhite struct {
					Animated struct {
						BackDefault      string `json:"back_default"`
						BackFemale       string `json:"back_female"`
						BackShiny        string `json:"back_shiny"`
						BackShinyFemale  string `json:"back_shiny_female"`
						FrontDefault     string `json:"front_default"`
						FrontFemale      string `json:"front_female"`
						FrontShiny       string `json:"front_shiny"`
						FrontShinyFemale string `json:"front_shiny_female"`
					} `json:"animated"`
					BackDefault      string `json:"back_default"`
					BackFemale       string `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"black-white"`
			} `json:"generation-v"`
			GenerationVi struct {
				OmegarubyAlphasapphire struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"omegaruby-alphasapphire"`
				XY struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"x-y"`
			} `json:"generation-vi"`
			GenerationVii struct {
				Icons struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  any    `json:"front_female"`
				} `json:"icons"`
				UltraSunUltraMoon struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"ultra-sun-ultra-moon"`
			} `json:"generation-vii"`
			GenerationViii struct {
				BrilliantDiamondShiningPearl struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  any    `json:"front_female"`
				} `json:"brilliant-diamond-shining-pearl"`
				Icons struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  any    `json:"front_female"`
				} `json:"icons"`
			} `json:"generation-viii"`
		} `json:"versions"`
	} `json:"sprites"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

func cleanInput(text string) []string {
	cleaned_output := make([]string, 0)
	text_clean := strings.Fields(text)
	for _, word := range text_clean {
		cleaned_output = append(cleaned_output, strings.ToLower(word))
	}
	return cleaned_output
}

func commandExit(cfg *config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, args ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage: ")
	fmt.Println("help: Displays a help message")
	fmt.Println("exit: Exit the Pokedex")
	fmt.Println("map: Displays the next 20 locations available...")
	fmt.Println("mapb: Displays the previous 20 locations available....")
	fmt.Println("explore: Lists all of the Pokemon available here...")
	fmt.Println("catch: catches a pokemon available in the area")
	fmt.Println("inspect: inspects a pokemon in the pokedex which you have caught")
	return nil
}

func commandMap(cfg *config, args ...string) error {
	location_area_url_base := cfg.Next
	if location_area_url_base == "" {
		location_area_url_base = "https://pokeapi.co/api/v2/location-area/?limit=20&offset=0"
	}
	res, err := http.Get(location_area_url_base)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	var locationList ShallowLocationList
	err = json.Unmarshal(body, &locationList)
	cfg.Next = locationList.Next
	cfg.Previous = locationList.Previous
	for _, loc := range locationList.Results {
		fmt.Println(loc.Name)
	}
	return nil
}

func commandMapb(cfg *config, args ...string) error {
	if cfg.Previous == "" {
		fmt.Println("you are on the first page, please use map to navigate forward!")
		return nil
	}
	location_previous := cfg.Previous
	res, err := http.Get(location_previous)
	if err != nil {
		fmt.Println("Error going backwards")
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	var locationList ShallowLocationList
	err = json.Unmarshal(body, &locationList)
	cfg.Next = locationList.Next
	cfg.Previous = locationList.Previous
	for _, loc := range locationList.Results {
		fmt.Println(loc.Name)
	}
	return nil
}

func commandExplore(cfg *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("No location provided to explore, please provide a location!")
	}
	location := args[0]
	body, ok := cfg.cache.Get(location)
	if !ok {
		//location := args
		url := "https://pokeapi.co/api/v2/location-area/" + location
		res, err := http.Get(url)
		if err != nil {
			fmt.Errorf("Error with location-name please try again!")
		}
		defer res.Body.Close()
		body, err = io.ReadAll(res.Body)
		if err != nil {
			fmt.Errorf("Error Reading into data variable!")
		}
		cfg.cache.Add(location, body)
	}
	var locationData Location
	err := json.Unmarshal(body, &locationData)
	if err != nil {
		fmt.Errorf("Issue with unmarshalling the data")
	}
	for item := range locationData.PokemonEncounters {
		pokemonName := locationData.PokemonEncounters[item].Pokemon.Name
		fmt.Println(pokemonName)
	}
	return nil
}

func commandCatch(cfg *config, args ...string) error {
	if len(args) == 0 {
		fmt.Errorf("Please enter a pokemon's name to catch!")
	}
	pokemonName := args[0]
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)
	url := "https://pokeapi.co/api/v2/pokemon/" + pokemonName
	res, err := http.Get(url)
	if err != nil {
		fmt.Errorf("Error with fetching pokemon data from the PokeAPI")
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Errorf("Error reading into data variable")
	}
	var pokemonData Pokemon
	err = json.Unmarshal(body, &pokemonData)
	if err != nil {
		fmt.Errorf("Issue with unmarshalling the pokemon data")
	}
	catchChance := (float64(rand.Intn(pokemonData.BaseExperience)) / float64(pokemonData.BaseExperience))
	const catchSuccess float64 = 0.5
	if catchChance > catchSuccess {
		cfg.caughtPokemon[pokemonName] = pokemonData
		fmt.Printf("You have successfully caught %s!!!!\n", pokemonName)
	}
	return nil
}

func commandInspect(cfg *config, args ...string) error {
	if len(args) == 0 {
		fmt.Errorf("Please enter the name of a pokemon in your inventory to inspect")
	}
	pokemonName := args[0]

	data, ok := cfg.caughtPokemon[pokemonName]
	if !ok {
		fmt.Errorf("%s is not in your list of caught pokemon, please enter a name from the pokemon you have caught")
	}
	fmt.Printf("Name: %s\n", pokemonName)
	fmt.Printf("Height: %d\n", data.Height)
	fmt.Printf("Weight: %d\n", data.Weight)
	fmt.Println("Stats:")
	fmt.Printf(" -hp: %d\n", data.Stats[0].BaseStat)
	fmt.Printf(" -attack: %d\n", data.Stats[1].BaseStat)
	fmt.Printf(" -defense: %d\n", data.Stats[2].BaseStat)
	fmt.Printf(" -special-attack: %d\n", data.Stats[3].BaseStat)
	fmt.Printf(" -special-defense: %d\n", data.Stats[4].BaseStat)
	fmt.Printf(" -speed: %d\n", data.Stats[5].BaseStat)
	fmt.Println("Types: ")
	fmt.Printf("  - %s\n", data.Types[0].Type.Name)
	fmt.Printf("  - %s\n", data.Types[1].Type.Name)
	return nil
}

func main() {
	cfg := &config{
		Next:          "https://pokeapi.co/api/v2/location-area/?limit=20&offset=0",
		Previous:      "",
		cache:         pokecache.NewCache(5 * time.Minute),
		caughtPokemon: make(map[string]Pokemon),
	}
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
		"explore": {
			name:        "explore",
			description: "Explores all the Pokemon available at the location",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "catches a pokemon in the location",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "inspects a pokemon to display it's stats",
			callback:    commandInspect,
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
		additionalInput := ""
		if len(words) > 1 {
			additionalInput = words[1]
		}
		err := command.callback(cfg, additionalInput)
		if err != nil {
			fmt.Println("Error executing command:", err)
		}
	}
}
