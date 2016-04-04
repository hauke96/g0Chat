package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"sync"
)

// ------------------------------
//
// ------------------------------

type Settings struct {
	username, ip, port, channel string
	messageLimit                int
	args, messageList           []string
	predefiningArgs             map[byte]string
}

var mutex = sync.Mutex{}
var currentInput string
var clientSettings *Settings
var printer printService

func main() {
	// ------------------------------
	// SHOW HELP IF WANTED
	// ------------------------------
	if printer.showHelp() { // true --> page was shown
		return
	}

	// ------------------------------
	// CREATE PARSER, SETTINGS AND PRINTER
	// ------------------------------
	clientSettings = parseConsoleArgs(os.Args)

	printer = printService{clientSettings}
	printer.welcomeDialog()

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

// read prints the display string onto the console and waits for a user intput.
// The input will be put into the return value.
// There'll be no additional/empty line when the display string is empty.
func read(display string) string {
	scanner := bufio.NewScanner(os.Stdin)
	if display != "" {
		fmt.Print(display)
		if display[len(display)-1] != ' ' {
			fmt.Print(" ")
		}
	}
	scanner.Scan()
	text := scanner.Text()
	return text //[0 : len(text)-1]
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
		clientSettings.messageList = append(clientSettings.messageList, message)
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
				exec.Command("stty", "-F", "/dev/tty", "sane", "echo").Run()
				send("("+clientSettings.username+" says bye)", connection)
				fmt.Println("\n\nHope you had fun, see you soon :)\n\n")
				return
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
