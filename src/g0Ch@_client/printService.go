package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

type printService struct {
	settings *Settings
	run      bool
}

// printAll prints the list of saved messages (length is specified by the
// messageLimit variable) an empty line and the user input.
func (p printService) printAll() {
	for p.run {
		// ------------------------------
		// PRINT MESSAGES
		// ------------------------------
		p.printMessages()
		fmt.Println()

		// ------------------------------
		// PRINT USER INPUT-FIELD
		// ------------------------------
		fmt.Print("> ", currentInput)

		// ------------------------------
		// DEALY IN RENDERING
		// ------------------------------
		time.Sleep(time.Millisecond * 100) // very complex anti-flicker-technique, hard to explain
	}
}

func (p printService) printMessages() {
	// ------------------------------
	// CLEAR CONSOLE
	// ------------------------------
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	// ------------------------------
	// PRINT MESSAGES
	// ------------------------------
	for _, v := range p.settings.messageList[len(p.settings.messageList)-p.settings.messageLimit:] {
		fmt.Print(v)
	}
}
