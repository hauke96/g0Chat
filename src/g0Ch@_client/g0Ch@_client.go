package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
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

	printer = printService{settings: clientSettings, run: true}

	// ------------------------------
	// CREATE CONNECTION
	// ------------------------------
	connection, err := net.Dial("tcp", clientSettings.ip+":"+clientSettings.port)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	// sign in
	connection.Write([]byte(string("\x01" + clientSettings.username + "\x1f" + clientSettings.channel + "\x04")))
	defer connection.Close()

	// ------------------------------
	// PREPARE CLEANUP FOR CTRL+C EVENT
	// ------------------------------
	prepareCleanup(connection)

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
	if message[len(message)-1:] != "\x04" {
		message += "\x04"
	}
	//	fmt.Println("\n")
	//	fmt.Println(message)
	//	fmt.Println([]byte(message))
	//	cleanup()
	connection.Write([]byte(message))
}

// waiter is a function that waits for any messages received by the connection to the server.
// When a message comes in, the screen will be re-rendered.
func waiter(connection net.Conn) {
	rawMessage, err := bufio.NewReader(connection).ReadString('\x04')

loop:
	for err == nil {

		mutex.Lock()
		message := rawMessage[1 : len(rawMessage)-1]
		switch rawMessage[0] {
		case '\x00':
			printer.run = false
			message = strings.Replace(message, "\x1f", ": ", 2)
			clientSettings.messageList = append(clientSettings.messageList, message)
			printer.printMessages()
		case '\x01':
			clientSettings.messageList = append(clientSettings.messageList, "[ "+message+" says hi :) ]")
		case '\x02':
			clientSettings.messageList = append(clientSettings.messageList, "[ "+message+" says bye :( ]")
		case '\x03':
			clientSettings.messageList = append(clientSettings.messageList, "[ SERVER: "+message+" ]")
			break loop
		}
		mutex.Unlock()

		rawMessage, err = bufio.NewReader(connection).ReadString('\x04')
	}

	printer.printMessages()
	printer.run = false
	if err != nil {
		fmt.Println("ERROR: Connection to server lost.")
	}
	cleanup(connection)
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
				cleanup(connection)
			} else {
				send("\x00"+currentInput+"\x04", connection)
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
func prepareCleanup(connection net.Conn) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGKILL)
	go func() {
		<-c
		cleanup(connection)
	}()
}

func cleanup(connection net.Conn) {
	send("\x02"+clientSettings.username+"\x04", connection)
	exec.Command("stty", "-F", "/dev/tty", "sane", "echo").Run()
	fmt.Println("\n\nHope you had fun, see you soon :)\n\n")
	os.Exit(1)
}
