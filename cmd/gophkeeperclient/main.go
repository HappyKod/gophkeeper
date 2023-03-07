package main

import (
	"github.com/pterm/pterm"
)

const registration = "registration"
const authorization = "authorization"

func main() {
	var options []string
	options = append(options, registration)
	options = append(options, authorization)

	selectedOption, _ := pterm.DefaultInteractiveSelect.WithOptions(options).Show()
	pterm.Info.Printfln("Selected option: %s", pterm.Green(selectedOption))
	if selectedOption == registration {

	} else {

	}

	pterm.Description.Printfln("TODO: create client")
}
