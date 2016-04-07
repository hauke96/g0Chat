package main

import (
	"fmt"
	"os"
	"os/exec"
	"text/tabwriter"
	"time"
)

type printService struct {
	settings *Settings
	run      bool
}

// showHelp simply shows the help message with all available parameters.
func (p printService) showHelp() bool {
	if len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h") {

		writer := new(tabwriter.Writer)
		writer.Init(os.Stdout, 0, 4, 2, ' ', 0)
		defer writer.Flush()

		fmt.Fprintln(writer, "Usage\t: option=value")
		fmt.Fprintln(writer, "Example\t: --username=Hugo")
		fmt.Fprintln(writer, "or simply\t  -u=Hugo")

		fmt.Fprintln(writer, "")

		fmt.Fprintln(writer, "Here're all command that are available:\n")

		fmt.Fprintln(writer, "-u, --username\tThe username/nickname.")
		fmt.Fprintln(writer, "-i, --ip\tThe IP of the g0Ch@ server.")
		fmt.Fprintln(writer, "-p, --port\tThe port of the g0Ch@ server (usually 10000).")
		fmt.Fprintln(writer, "-l, --limit\tValue for the size of the message buffer.\n\t(how many messages are stored)")
		fmt.Fprintln(writer, "-c, --channel\tThe channel you want to talk in.")
		fmt.Fprintln(writer, "-h, --help\tShows this page.")

		fmt.Fprintln(writer, "\nType  e x i t  as message to leave the chat.")

		fmt.Println()

		os.Stdin.Read([]byte{})
		return true
	}
	return false
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
	fmt.Println(len(p.settings.messageList))
	for _, v := range p.settings.messageList[len(p.settings.messageList)-p.settings.messageLimit:] {
		fmt.Print(v)
	}
}
