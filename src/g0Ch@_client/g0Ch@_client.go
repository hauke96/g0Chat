package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
)

// ------------------------------
//
// ------------------------------

type Settings struct {
	username, ip, port, channel string
	messageLimit                int
	messageList                 []string
}

var mutex = sync.Mutex{}
var currentInput string
var clientSettings *Settings
var printer printService

func main() {
	// ------------------------------
	// CREATE PARSER, SETTINGS AND PRINTER
	// ------------------------------
	clientSettings = parseConsoleArgs()

	// ------------------------------
	// PREPARE CLEANUP FOR CTRL+C EVENT
	// ------------------------------
	prepareCleanup()

	printer = printService{settings: clientSettings, run: true}

	// ------------------------------
	// CREATE CONNECTION
	// ------------------------------
	connection, err := net.Dial("tcp", clientSettings.ip+":"+clientSettings.port)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	connection.Write([]byte(clientSettings.channel + "\n"))
	connection.Write([]byte(string("(" + clientSettings.username + " says hi)\n")))
	defer connection.Close()

	// ------------------------------
	// START WAITER FOR INPUT
	// ------------------------------
	go waiter(connection)

	// ------------------------------
	// START PRINTER
	// ------------------------------
	go printer.printAll()

	chat(connection)
}

// send a message to the server. This function adds the channel prefix
// and the \n suffix, so that you can't forget it ;)
func send(message string, connection net.Conn) {
	connection.Write([]byte(string(clientSettings.channel + "\x02" + message + "\n")))
}

// waiter is a function that waits for any messages received by the connection to the server.
// When a message comes in, the screen will be re-rendered.
func waiter(connection net.Conn) {
	message, err := bufio.NewReader(connection).ReadString('\n')

	for err == nil {

		mutex.Lock()
		if message[0] == '\x04' {
			printer.run = false
			clientSettings.messageList = append(clientSettings.messageList, message[1:])
			printer.printMessages()
			cleanup()
		} else {
			clientSettings.messageList = append(clientSettings.messageList, message)
		}
		mutex.Unlock()

		message, err = bufio.NewReader(connection).ReadString('\n')
	}

	fmt.Println("ERROR:", err)
}

// chat waits for user input and sends the given string to the server.
func chat(connection net.Conn) {
	// ------------------------------
	// PREPARE SCREEN
	// ------------------------------
	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

	var b []byte = make([]byte, 1)
	for {
		os.Stdin.Read(b) // wait for input

		switch b[0] {
		case '\n':
			// ------------------------------
			// ENTER PRESSED
			// ------------------------------
			if currentInput == "exit" { // when exit --> restore printing of characters
				send("("+clientSettings.username+" says bye)", connection)
				cleanup()
			} else {
				send(clientSettings.username+": "+currentInput, connection)
				currentInput = ""
			}

		case '\u007F':
			// ------------------------------
			// BACKSPACE PRESSED
			// ------------------------------

			// \u007F is a backspace to delete last character. The stty settings do not
			// allow a normal character deletion in the input. Maybe there's a better
			// solution without this hack, but it works for now.
			if currentInput != "" { // dont put this into outer if, because unfancy character will be printed :/
				currentInput = currentInput[0 : len(currentInput)-1]
			}

		default:
			// ------------------------------
			// CHARACTER TYPED
			// ------------------------------
			currentInput += string(b)
		}
	}
}

// prepareCleanup creates a routine that's executed when either a os.Interrupt or a SIGTERM event is fired. This will call the cleanup() function.
func prepareCleanup() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(1)
	}()
}

func cleanup() {
	exec.Command("stty", "-F", "/dev/tty", "sane", "echo").Run()
	fmt.Println("\n\nHope you had fun, see you soon :)\n\n")
	os.Exit(1)
}
