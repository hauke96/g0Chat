package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

// ------------------------------
//
// ------------------------------

type Settings struct {
	port string
}

type Connection struct {
	listener   net.Listener
	connection net.Conn
	number     int
	channel    string
}

var allConnections = make([]Connection, 0)
var mutex = &sync.Mutex{}
var settings = &Settings{port: "10000"}

func main() {
	// ------------------------------
	// PARSE ARGUMENTS
	// ------------------------------
	settings = parseConsoleArgs()

	// ------------------------------
	// PREPARE CLEANUP FOR CTRL+C EVENT
	// ------------------------------
	prepareCleanup()

	fmt.Println("START SERVER ...")

	ch := make(chan (Connection))
	openConn := 0
	count := 0

	// ------------------------------
	// ROUTINE FOR CLOSED CONNECTIONS
	// ------------------------------
	go func() {
		for {
			conn := <-ch
			fmt.Print("[ ", conn.number, " ] CLOSING CONNECTION ", conn.number, " ON PORT ", conn.listener.Addr().String()[5:], "\n")
		}
	}()

	// ------------------------------
	// START LISTENING
	// ------------------------------
	fmt.Print("[ ", count, " ] WAITING FOR LISTENER ...")
	listener, err := net.Listen("tcp", ":"+settings.port)

	if err != nil {
		fmt.Println("FAIL")
		fmt.Println(err)
		return
	}
	defer listener.Close()
	fmt.Println("OK")
	fmt.Println()

	// ------------------------------
	// ROUTINE FOR NEW CONNECTIONS
	// ------------------------------
	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("FAIL")
			fmt.Println(err)
			return
		}
		fmt.Print("[ ", count, " ] ACCEPTED CONNECTION ...")
		defer connection.Close()
		fmt.Println("OK")

		channel, err := bufio.NewReader(connection).ReadString('\n')
		if err == nil {

			fmt.Print("[ ", count, " ] CHANNEL: ", channel)

			conn := Connection{
				listener:   listener,
				connection: connection,
				number:     count,
				channel:    channel[0 : len(channel)-1]}

			allConnections = append(allConnections, conn)

			go chatter(conn, ch)

			openConn++
			count++
		} else {
			fmt.Print("[ ", count, " ] ERROR WHILE ACCEPTING THE CONNECTION: ", err, "\n")
		}

		fmt.Println()
	}
}

// chatter is a service function that cares about one connection to a client.
// Whenever an user input comes in, it'll be send to all other clients.
func chatter(connection Connection, ch chan Connection) {
	message, err := bufio.NewReader(connection.connection).ReadString('\n')

	for err == nil {
		fmt.Print("[ ", connection.number, " ] INCOMING:  ", message)
		splittedMessage := strings.Split(message, "\x02")
		if len(splittedMessage) == 2 {
			fmt.Print("[ ", connection.number, " ] MESSAGE:   ", splittedMessage[1])
			fmt.Print("[ ", connection.number, " ] ON CHANNEL ", splittedMessage[0], "\n")
			notifyAll(splittedMessage[1], splittedMessage[0])
		} else {
			fmt.Print("[ ", connection.number, " ] ERROR: Splitting the message failed.")
		}

		fmt.Println()

		message, err = bufio.NewReader(connection.connection).ReadString('\n')
	}

	ch <- connection
}

// notifyAll sends a message to all connected clients.
func notifyAll(message, channel string) {
	mutex.Lock()
	for _, c := range allConnections {
		if c.channel == channel {
			c.connection.Write([]byte(message))
		}
	}
	mutex.Unlock()
}

// prepareCleanup creates a routine that's executed when either a os.Interrupt or a SIGTERM event is fired. This will call the cleanup() function.
func prepareCleanup() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\nSTOP SERVER ...")
		cleanup()
		fmt.Println("BYE BYE")
		os.Exit(1)
	}()
}

// cleanup closes all connection so that the client knows that the server has been shut down.
func cleanup() {
	mutex.Lock()
	fmt.Println("CLOSE ALL " + strconv.Itoa(len(allConnections)) + " CONNECTIONS...")
	for _, c := range allConnections {
		fmt.Println("  CONN NR. " + strconv.Itoa(c.number) + " ON CHAN " + c.channel)
		c.connection.Write([]byte("\x04SERVER: Server is shutting down. Good bye :)\n"))
		c.connection.Close()
	}
	mutex.Unlock()
	fmt.Println("DONE")
}
