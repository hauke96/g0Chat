package main

import (
	"../GeneralParser"
)

// parseConsoleArgs takes the arguments and parses them into two categories:
// normal and predefinings arguments.
// It also evaluates the predefining ones.
func parseConsoleArgs(args []string, settings *Settings) {
	args, predefiningArgs := GeneralParser.ParseArgs(args, ":p:", ":port:")
	parsePredefined(args, predefiningArgs, settings)
}

func parsePredefined(args []string, predefiningArgs map[byte]string, settings *Settings) {
	for key, value := range predefiningArgs {
		switch key {
		case 'p':
			settings.port = value
		}
	}
}
