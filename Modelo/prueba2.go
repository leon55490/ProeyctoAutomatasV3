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

var automaton *Automaton // Define the loaded automaton pointer

type Automaton struct {
	States       []string                     `json:"states"`
	Alphabet     []string                     `json:"alphabet"`
	Transitions  map[string]map[string]string `json:"transitions"`
	InitialState string                       `json:"initialState"`
	FinalStates  []string                     `json:"finalStates"`
}

func (a *Automaton) Run(input string) string {
	fmt.Println("Ultima funcion")
	currentState := a.InitialState
	for _, c := range input {
		transition, ok := a.Transitions[currentState][string(c)]
		if !ok {
			return "Error: no transition exists for state " + currentState + " and input " + string(c)
		}
		currentState = transition
	}
	fmt.Println(string(input))
	for _, finalState := range a.FinalStates {
		if currentState == finalState {
			return "ACEPTADO"
		}
	}
	return "Rechazado"
}

func main() {
	fmt.Println("Corre el main")
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

	runButton := widget.NewButton("Run", func() {
		// Run the loaded automaton with the current input
		if automaton == nil {
			return
		}
		fmt.Println("aqui esta el nuevo mensaje")
		fmt.Println(string(inputEntry.Text))
		output := automaton.Run(inputEntry.Text)
		outputText.SetText(output)
	})

	// Create button widgets
	loadButton := widget.NewButton("Load Automaton", func() {
		fmt.Println("Muestra la opcion de cargar archivos")
		fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader != nil {
				defer reader.Close()
				bytes, err := ioutil.ReadAll(reader)
				if err != nil {
					dialog.ShowError(err, myWindow)
					return
				}
				err = loadAutomatonFromBytes(bytes)
				if err != nil {
					dialog.ShowError(err, myWindow)
					return
				}
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

func loadAutomatonFromBytes(data []byte) error {
	fmt.Println("Lee y carga archivos")
	err := json.Unmarshal(data, &automaton)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON: %v", err)
	}
	return nil
}
