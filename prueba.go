package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type Transition struct {
	Symbol string `json:"symbol"`
	From   string `json:"from"`
	To     string `json:"to"`
}

type State struct {
    Name        string                      `json:"name"`
    IsInitial   bool                        `json:"is_initial"`
    IsFinal     bool                        `json:"is_final"`
    Transitions map[string]*Transition      `json:"transitions"`
}

type Automaton struct {
    States       []*State          `json:"states"`
    Transitions  map[string]map[string]Transition `json:"transitions"`
    InitialState string            `json:"initial_state"`
    FinalStates  map[string]bool   `json:"final_states"`
}

func (a *Automaton) Run(input string) string {
    currentState := a.InitialState
    for _, c := range input {
        transition, ok := a.Transitions[currentState][string(c)]
        if !ok {
            return "Error: no transition exists for state " + currentState + " and input " + string(c)
        }
        currentState = transition.To
    }
    if a.FinalStates[currentState] {
        return "Accepted"
    } else {
        return "Rejected"
    }
}

func main() {
	// Create a new Fyne application
	myApp := app.New()

	// Create a new window
	myWindow := myApp.NewWindow("Automaton")

	// Create input widgets
	inputLabel := widget.NewLabel("Input string:")
	inputEntry := widget.NewEntry()

	// Create output widgets
	outputLabel := widget.NewLabel("Result:")
	outputText := widget.NewLabel("")

	var loadedAutomaton *Automaton // Define the loaded automaton pointer

	runButton := widget.NewButton("Run", func() {
		// Run the loaded automaton with the current input
		if loadedAutomaton == nil {
			return
		}
		output := loadedAutomaton.Run(inputEntry.Text)
		outputText.SetText(output)
	})

	// Create button widgets
	loadButton := widget.NewButton("Load Automaton", func() {
		fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader != nil {
				defer reader.Close()
				bytes, err := ioutil.ReadAll(reader)
				if err != nil {
					dialog.ShowError(err, myWindow)
					return
				}
				automaton, err := loadAutomatonFromBytes(bytes)
				if err != nil {
					dialog.ShowError(err, myWindow)
					return
				}
				loadedAutomaton = automaton // Save the loaded automaton
				inputEntry.SetText("")
				outputText.SetText("")
				runButton.OnTapped() // Automatically run the loaded automaton
			}
		}, myWindow)
		fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
		fileDialog.Show()
	})

	// Create a container for input widgets
	inputContainer := container.NewVBox(inputLabel, inputEntry)

	// Create a container for output widgets
	outputContainer := container.NewVBox(outputLabel, outputText)

	// Create a container for the input and output containers and the buttons
	buttonContainer := container.NewHBox(loadButton, runButton)
	content := container.NewVBox(inputContainer, outputContainer, buttonContainer)

	// Add the content container to the window
	myWindow.SetContent(content)

	// Show the window and start the app
	myWindow.ShowAndRun()
}


func loadAutomatonFromBytes(data []byte) (*Automaton, error) {
	// fmt.Println(string(data))
    var automaton Automaton
    err := json.Unmarshal(data, &automaton)
    if err != nil {
        return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
    }

    // Convert the transitions map from map[string]Transition to map[string]*Transition
    for _, state := range automaton.States {
        for symbol, transition := range state.Transitions {
            state.Transitions[symbol] = transition
        }
    }

    // Find the initial state of the automaton
    for i, state := range automaton.States {
        if state.IsInitial {
            automaton.InitialState = state.Name
            // Move the initial state to the front of the slice for easier access
            if i != 0 {
                automaton.States[0], automaton.States[i] = automaton.States[i], automaton.States[0]
            }
            break
        }
    }

    return &automaton, nil
}

func loadAutomatonFromFile(filename string) (*Automaton, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	automaton, err := loadAutomatonFromBytes(data)
	if err != nil {
		// Catch any errors thrown by loadAutomatonFromBytes and return a more descriptive error
		return nil, fmt.Errorf("error loading automaton from file: %v", err)
	}

	return automaton, nil
}