package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"text/tabwriter"
	"time"
)

var mutex = sync.Mutex{}
var messageList []string
var currentInput, username, ip, port string
var messageLimit int

func main() {
	if showHelp() { // true --> page was shown
		return
	}

	welcomeDialog()

	connection, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return
	}
	connection.Write([]byte(string("(" + username + " says hi)\n")))

	go waiter(connection)

	go printAll()

	chat(connection)
}

// showHelp simply shows the help message with all available parameters.
func showHelp() bool {
	if len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h") {
		fmt.Println("Here're all command that are available:\n")

		writer := new(tabwriter.Writer)
		writer.Init(os.Stdout, 0, 4, 2, ' ', 0)
		defer writer.Flush()

		fmt.Fprintln(writer, "-b\tAdditional field for the size of the message buffer.\n\t(how many messages are stored)")
		fmt.Fprintln(writer, "-h, --help\tShows this page.")

		fmt.Fprintln(writer, "\nType e x i t as message to leave the chat.")

		fmt.Println()

		os.Stdin.Read([]byte{})
		return true
	}
	return false
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

// welcomeDialog asks for the username, port, ip and some other values that are
// necessary or specified by the console arguments.
func welcomeDialog() {
	username = read("Choose username: ")

	ip = read("IP: ")

	port = read("Port (normally 10000): ")

	// some defaults:
	messageLimit = 50

	for _, arg := range os.Args {
		switch arg {
		case "-b":
			limitInput := read("Message limit:")
			limit, err := strconv.Atoi(limitInput)
			for err != nil {
				if err != nil {
					fmt.Println("ERROR: Maybe", limit, "is not a number?")
				}
				limitInput = read("Message limit:")
				limit, err = strconv.Atoi(limitInput)
			}
			messageLimit = limit
		}
	}

	messageList = make([]string, messageLimit)
}

// waiter is a function that waits for any messages received by the connection to the server.
// When a message comes in, the screen will be re-rendered.
func waiter(connection net.Conn) {
	message, err := bufio.NewReader(connection).ReadString('\n')

	for err == nil {
		mutex.Lock()
		messageList = append(messageList, message)
		mutex.Unlock()
		message, err = bufio.NewReader(connection).ReadString('\n')
	}

	fmt.Println("ERROR: ", err)
}

// chat waits for user input and sends the given string to the server.
func chat(connection net.Conn) {
	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

	var b []byte = make([]byte, 1)
	for {
		os.Stdin.Read(b)
		// enter --> send string
		if b[0] == '\n' {
			// when exit --> restore printing of characters
			if currentInput == "exit" {
				exec.Command("stty", "-F", "/dev/tty", "sane", "echo").Run()
				connection.Write([]byte(string("(" + username + " says bye)\n")))
				fmt.Println("\n\nHope you had fun, see you soon :)\n\n")
				connection.Close()
				return
			} else {
				connection.Write([]byte(string(username + ": " + currentInput + "\n")))
				currentInput = ""
			}
		} else if b[0] == '\u007F' {
			// \u007F is a backspace to delete last character. The stty settings do not
			// allow a normal character deletion in the input. Maybe there's a better
			// solution without this hack, but it works for now.
			if currentInput != "" { // dont put this into outer if, because unfancy character will be printed :/
				currentInput = currentInput[0 : len(currentInput)-1]
			}
		} else {
			currentInput += string(b)
		}
	}
}

// printAll prints the list of saved messages (length is specified by the
// messageLimit variable) an empty line and the user input.
func printAll() {
	for {
		// clearing the console
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
		// writing stuff
		for _, v := range messageList[len(messageList)-messageLimit:] {
			fmt.Print(v)
		}
		fmt.Println()
		fmt.Print("> ", currentInput)
		time.Sleep(time.Millisecond * 100)
	}
}
