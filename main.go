package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Transition struct {
	Symbol string `json:"symbol"`
	From   string `json:"from"`
	To     string `json:"to"`
}

type State struct {
	Name       string        `json:"name"`
	IsInitial  bool          `json:"is_initial"`
	IsFinal    bool          `json:"is_final"`
	Transitions []Transition `json:"transitions"`
}

type Automaton struct {
	States        []State  `json:"states"`
	Alphabet      []string `json:"alphabet"`
	InitialState  string   `json:"initial_state"`
}

func main() {
	// Load the automaton from a file
	automaton, err := loadAutomatonFromFile("automaton.json")
	if err != nil {
		log.Fatal(err)
	}

	// Load the input string from a file
	input, err := loadInputStringFromFile("input.json")
	if err != nil {
		log.Fatal(err)
	}

	// Create a new Fyne application
	myApp := app.New()

	// Create a new window
	myWindow := myApp.NewWindow("Automaton")

	// Create output widgets
	outputLabel := widget.NewLabel("Result:")
	outputText := widget.NewLabel("")

	// Create button widget
	runButton := widget.NewButton("Run", func() {
		// Test whether the input string is accepted by the automaton
		if automaton.Accepts(input) {
			outputText.SetText("Accepted")
		} else {
			outputText.SetText("Rejected")
		}
	})

	// Create a container for output widgets
	outputContainer := container.NewVBox(outputLabel, outputText)

	// Create a container for the output container and the button
	content := container.NewVBox(outputContainer, runButton)

	// Add the content container to the window
	myWindow.SetContent(content)

	// Show the window and start the app
	myWindow.ShowAndRun()
}

func (a Automaton) Accepts(input string) bool {
	currentState := a.InitialState
	for _, c := range input {
		symbol := string(c)
		nextState := ""
		for _, state := range a.States {
			if state.Name == currentState {
				for _, transition := range state.Transitions {
					if transition.Symbol == symbol {
						nextState = transition.To
						break
					}
				}
				break
			}
		}
		if nextState == "" {
			return false
		}
		currentState = nextState
	}
	for _, state := range a.States {
		if state.Name == currentState {
			return state.IsFinal
		}
	}
	return false
}

func loadAutomatonFromFile(filename string) (Automaton, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Automaton{}, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return Automaton{}, err
	}

	var automaton Automaton
	err = json.Unmarshal(bytes, &automaton)
	if err != nil {
		return Automaton{}, err
	}

	for _, state := range automaton.States {
		if state.IsInitial {
			automaton.InitialState = state.Name
			break
		}
	}

	return automaton, nil
}

func loadInputStringFromFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var input struct {
		Input string `json:"input"`
	}

	err = json.NewDecoder(file).Decode(&input)
	if err != nil {
		return "", err
	}

	return input.Input, nil
}