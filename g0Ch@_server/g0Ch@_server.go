package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"sync"
)

type Connection struct {
	listener   net.Listener
	connection net.Conn
	number     int
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
			fmt.Println("CLOSING CONNECTION", conn.number, "ON PORT", conn.listener.Addr().String()[5:])
		}
	}()

	count := 0

	// listen on port
	fmt.Println("[", count, "] WAITING FOR LISTENER ...")
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(10000))
	if err != nil {
		fmt.Println(err)
		return
	}

	for true {
		// connect
		fmt.Println("[", count, "] ACCEPTING CONNECTION ...")
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		conn := Connection{listener: listener, connection: connection, number: count}
		allConnections = append(allConnections, conn)

		go chatter(conn, ch)

		openConn++
		count++
	}

	fmt.Println("STOP SERVER ...")
}

// chatter is a service function that cares about one connection to a client.
// Whenever an user input comes in, it'll be send to all other clients.
func chatter(connection Connection, ch chan Connection) {
	message, err := bufio.NewReader(connection.connection).ReadString('\n')

	for err == nil {
		fmt.Print("INCOMING: ", message)
		notifyAll(message)
		message, err = bufio.NewReader(connection.connection).ReadString('\n')
	}

	ch <- connection
}

// notifyAll sends a message to all connected clients.
func notifyAll(message string) {
	mutex.Lock()
	for _, c := range allConnections {
		c.connection.Write([]byte(message))
	}
	mutex.Unlock()
}
