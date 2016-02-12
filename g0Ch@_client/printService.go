package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"text/tabwriter"
	"time"
)

type printService struct {
	settings *settings
}

// welcomeDialog asks for the username, port, ip and some other values that are
// necessary or specified by the console arguments.
// When an property is already predefined via the console arguments, there'll be
// no dialog for this.
func (p printService) welcomeDialog() {
	// ------------------------------
	// CHECK FOR UN-PREDEFINED ARGS
	// ------------------------------
	predefinedArgs := p.settings.predefiningArgs
	if _, b := predefinedArgs['u']; !b {
		p.settings.username = read("Choose username: ")
	}

	if _, b := predefinedArgs['i']; !b {
		p.settings.ip = read("IP: ")
	}

	if _, b := predefinedArgs['p']; !b {
		p.settings.port = read("Port (normally 10000): ")
	}

	// ------------------------------
	// SET DEFAULS
	// ------------------------------
	p.settings.messageLimit = 50

	// ------------------------------
	// CHECK NORMAL ARGS
	// ------------------------------
	for _, arg := range p.settings.args {
		switch arg {
		case "-l":
			p.settings.messageLimit = p.askForLimit()
		}
	}

	p.settings.messageList = make([]string, p.settings.messageLimit)
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
		fmt.Fprintln(writer, "-h, --help\tShows this page.")

		fmt.Fprintln(writer, "\nType  e x i t  as message to leave the chat.")

		fmt.Println()

		os.Stdin.Read([]byte{})
		return true
	}
	return false
}

// askForLimit will show a dialog where the user can enter a size limit for the message list.
func (p printService) askForLimit() int {
	limitInput := read("Message limit:")
	limit, err := strconv.Atoi(limitInput)
	for err != nil {
		if err != nil {
			fmt.Println("ERROR: Maybe", limit, "is not a number?")
		}
		limitInput = read("Message limit:")
		limit, err = strconv.Atoi(limitInput)
	}
	return limit
}

// printAll prints the list of saved messages (length is specified by the
// messageLimit variable) an empty line and the user input.
func (p printService) printAll() {
	for {
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
