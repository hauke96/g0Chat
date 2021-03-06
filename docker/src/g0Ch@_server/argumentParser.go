package main

import (
	"../GeneralParser"
)

// parseConsoleArgs takes the arguments and parses them into two categories:
// normal and predefinings arguments.
// It also evaluates the predefining ones.
func parseConsoleArgs() *Settings {
	p := GeneralParser.NewParser()

	port := p.RegisterArgument("port", "p", "The port of the g0Ch@ server (usually 44494)").Default("44494").String()

	p.Parse()

	settings := Settings{}
	settings.port = *port

	return &settings
}
