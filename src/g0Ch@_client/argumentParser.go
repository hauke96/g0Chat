package main

import (
	"../GeneralParser"
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// parseConsoleArgs takes the arguments and parses them into two categories:
// normal and predefinings arguments.
// It also evaluates the predefining ones.
func parseConsoleArgs() *Settings {
	p := GeneralParser.NewParser()
	p.Description("Type  e x i t  as message to leave the chat.")

	username := p.RegisterArgument("username", "u", "The username/nickname").String()
	limit := p.RegisterArgument("limit", "l", "Value for the size of the message buffer (how many messages are stored)").Default("50").Int()
	port := p.RegisterArgument("port", "p", "The port of the g0Ch@ server (usually 44494)").Default("44494").String()
	ip := p.RegisterArgument("ip", "i", "The IP of the g0Ch@ server").String()
	channel := p.RegisterArgument("channel", "c", "The channel you want to talk in").String()

	p.Parse()

	settings := Settings{}

	if *username != "" {
		settings.username = *username
	} else {
		settings.username = read("Chooose username: ")
	}

	if *limit != 0 {
		settings.messageLimit = *limit
	} else {
		limit := read("Message limit: ")
		limitInt, err := strconv.Atoi(limit)
		for err != nil {
			fmt.Println("ERROR: Maybe", limit, "is not a number?")
			limit = read("Message limit: ")
			limitInt, err = strconv.Atoi(limit)
		}
		settings.messageLimit = limitInt
	}
	settings.messageList = make([]string, settings.messageLimit)

	settings.port = *port

	if *ip != "" {
		settings.ip = *ip
	} else {
		settings.ip = read("IP: ")
	}

	if *channel != "" {
		settings.channel = *channel
	} else {
		settings.channel = read("Channel: ")
	}

	return &settings
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
	return text
}
