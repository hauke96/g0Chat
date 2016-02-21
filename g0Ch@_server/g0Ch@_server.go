package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

type Connection struct {
	listener   net.Listener
	connection net.Conn
	number     int
	channel    string
}

var allConnections = make([]Connection, 0)
var mutex = &sync.Mutex{}

func main() {
	fmt.Println("START SERVER ...")

	ch := make(chan (Connection))
	openConn := 0

	go func() {
		for true {
			conn := <-ch
			fmt.Println("[", conn.number, "] CLOSING CONNECTION", conn.number, "ON PORT", conn.listener.Addr().String()[5:])
		}
	}()

	count := 0

	// listen on port
	fmt.Println("[", count, "] WAITING FOR LISTENER ...")
	listener, err := net.Listen("tcp", ":10000")
	defer listener.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println()

	for true {
		// connect
		connection, err := listener.Accept()
		fmt.Println("[", count, "] ACCEPTED CONNECTION ...")
		defer connection.Close()
		if err != nil {
			fmt.Println(err)
			return
		}

		channel, err := bufio.NewReader(connection).ReadString('\n')
		if err == nil {

			fmt.Print("[ ", count, " ] CHANNEL:", channel)

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
			fmt.Println("[", count, "] ERROR WHILE ACCEPTING THE CONNECTION:", err)
		}

		fmt.Println()
	}

	fmt.Println("STOP SERVER ...")
}

// chatter is a service function that cares about one connection to a client.
// Whenever an user input comes in, it'll be send to all other clients.
func chatter(connection Connection, ch chan Connection) {
	message, err := bufio.NewReader(connection.connection).ReadString('\n')

	for err == nil {
		fmt.Print("[ ", connection.number, " ] INCOMING:   ", message)
		splittedMessage := strings.Split(message, "\x02")
		fmt.Print("[ ", connection.number, " ] MESSAGE:    ", splittedMessage[1])
		fmt.Println("[", connection.number, "] ON CHANNEL ", splittedMessage[0])
		notifyAll(splittedMessage[1], splittedMessage[0])

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
