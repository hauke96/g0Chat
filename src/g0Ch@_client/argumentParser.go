package main

import (
	"../GeneralParser"
	"strconv"
)

// parseConsoleArgs takes the arguments and parses them into two categories:
// normal and predefinings arguments.
// It also evaluates the predefining ones.
func parseConsoleArgs(args []string) *Settings {
	args, predefiningArgs := GeneralParser.parseArgs(args, ":u:l:i:p:c:", ":username:limit:port:ip:channel:")
	return parsePredefined(args, predefiningArgs)
}

func parsePredefined(args []string, predefiningArgs map[byte]string) *Settings {
	settings := Settings{args: args, predefiningArgs: predefiningArgs}

	for key, value := range predefiningArgs {
		switch key {
		case 'u':
			settings.username = value
		case 'i':
			settings.ip = value
		case 'p':
			settings.port = value
		case 'l':
			i, err := strconv.Atoi(value)
			if err == nil {
				settings.messageLimit = i
			}
		case 'c':
			settings.channel = value
		}
	}

	return &settings
}
