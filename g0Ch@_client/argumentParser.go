package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

type argumentParser struct {
	args            []string
	predefiningArgs map[byte]string
}

// parseArgs takes the arguments and parses them into two categories:
// normal and predefinings arguments.
// It also evaluates the predefining ones.
func (a *argumentParser) parseArgs(args []string) *settings {
	// ------------------------------
	// INIT
	// ------------------------------
	a.args = args
	a.predefiningArgs = make(map[byte]string, 0) // e.g. port=10000

	// adding known arg with the syntax :argname: to better check if an arg is known or not
	knownShortArgs := ":u:l:i:p:"
	knownLongArgs := ":username:limit:port:ip:"

	a.args = a.args[1:] // leave the first argument, which is the applications path

	writer := new(tabwriter.Writer)
	writer.Init(os.Stdout, 0, 4, 2, ' ', 0)
	defer writer.Flush()

	// ------------------------------
	// SPLIT COMBINED ARGS
	// ------------------------------
	// E.g. [-bgh] -> [-b] [-g] [-h]
	newArgs := make([]string, 0)
	for _, arg := range a.args {
		if !strings.Contains(arg, "=") && arg[0] == '-' && arg[1] != '-' && len(arg) > 2 { // 3 because of - and at least 2 other characters

			arg = arg[1:] // remove the -

			for _, v := range arg {
				newArgs = append(newArgs, "-"+string(v))
			}
		} else {
			newArgs = append(newArgs, arg)
		}
	}

	a.args = newArgs

	// ------------------------------
	// SEPARATE INTO NORMAL AND PREDEFINING
	// ------------------------------
	newArgs = make([]string, 0)
	for _, arg := range a.args {
		splittedArg := strings.Split(arg, "=")

		// noteice: only --foo or -f are allowed, -foo and --f are not allwed!
		// This is just to have the normal feeling of arguments in linux, blame me but i like it ;)
		if splittedArg[0][0:2] == "--" && len(splittedArg[0]) > 3 { // 3 because of -- and at least 2 other characters
			splittedArg[0] = splittedArg[0][2:]
		} else if len(splittedArg[0]) == 2 { // - and another character
			splittedArg[0] = splittedArg[0][1:]
		}

		if len(splittedArg[0]) == 1 && strings.Contains(knownShortArgs, ":"+splittedArg[0]+":") ||
			len(splittedArg[0]) > 1 && strings.Contains(knownLongArgs, ":"+splittedArg[0]+":") { // is it a valid short or long argument?

			if len(splittedArg) >= 2 { // long argument
				fmt.Fprintln(writer, "PREDEF:", splittedArg[0], "\t=", splittedArg[1])
				a.predefiningArgs[splittedArg[0][0]] = splittedArg[1]
			} else { // short argument
				newArgs = append(newArgs, string(splittedArg[0][0]))
			}
		} else { // not valid
			fmt.Println("ERROR: Unknown argument", splittedArg[0], "but I'll ignore it :/")
		}
	}

	return a.parsePredefined()
}

func (a *argumentParser) parsePredefined() *settings {
	settings := settings{args: a.args, predefiningArgs: a.predefiningArgs}

	for key, value := range a.predefiningArgs {
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
		}
	}

	return &settings
}
