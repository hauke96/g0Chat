package GeneralParser

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

type parser struct {
	args           map[string]*argument
	longToShortArg map[string]string
	knownLongArgs  string
	knownShortArgs string
}

func NewParser() *parser {
	p := parser{
		args:           make(map[string]*argument),
		longToShortArg: make(map[string]string),
		knownLongArgs:  ":",
		knownShortArgs: ":",
	}
	return &p
}

func (p *parser) RegisterArgument(longKey, shortKey, help string) (*argument, error) {
	// the long key is allowed to be empty, the short key not
	if shortKey == "" {
		return nil, errors.New("Only the long key can be empty!")
	}
	// when the long key and a fitting entry exists, we've a duplicate here
	if longKey != "" && p.longToShortArg[longKey] != "" {
		return nil, errors.New("The long key " + longKey + " already exists!")
	}
	// when there's a short key entry already, we've a duplicate here
	if p.args[shortKey] != nil {
		return nil, errors.New("The short key " + shortKey + " already exists!")
	}

	a := argument{
		longKey:  longKey,
		shortKey: shortKey,
		helpText: help,
	}

	p.knownShortArgs += shortKey + ":"
	p.knownLongArgs += longKey + ":"
	p.longToShortArg[longKey] = shortKey

	p.args[shortKey] = &a
	return &a, nil
}

func (p *parser) Parse() {
	p.parseArgs(os.Args[1:])
}

// parseArgs takes the arguments and parses them into two categories:
// normal and predefinings arguments.
// It also evaluates the predefining ones.
func (p *parser) parseArgs(args []string) {
	// ------------------------------
	// INIT
	// ------------------------------
	//	predefiningArgs := make(map[byte]string, 0) // e.g. --port=10000

	// adding known arg with the syntax :argname: to better check if an arg is known or not
	//	knownShortArgs := ":p:c:"
	//	knownLongArgs := ":port:channels:"

	writer := new(tabwriter.Writer)
	writer.Init(os.Stdout, 0, 4, 2, ' ', 0)
	defer writer.Flush()

	// ------------------------------
	// SPLIT COMBINED ARGS
	// ------------------------------
	// E.g. [-bgh] -> [-b] [-g] [-h]
	newArgs := make([]string, 0)
	for _, arg := range args {
		if !strings.Contains(arg, "=") && arg[0] == '-' && arg[1] != '-' && len(arg) > 2 { // 3 because of - and at least 2 other characters

			arg = arg[1:] // remove the -

			for _, v := range arg {
				newArgs = append(newArgs, "-"+string(v))
			}
		} else {
			newArgs = append(newArgs, arg)
		}
	}

	args = newArgs

	// ------------------------------
	// SEPARATE INTO NORMAL AND PREDEFINING
	// ------------------------------
	newArgs = make([]string, 0)
	for _, arg := range args {
		splittedArg := strings.Split(arg, "=")

		// notice: only --foo or -f are allowed, -foo and --f are not allwed!
		// This is just to have the normal feeling of arguments in linux, blame me but i like it ;)
		if splittedArg[0][0:2] == "--" && len(splittedArg[0]) > 3 { // 3 because of -- and at least 2 other characters
			splittedArg[0] = splittedArg[0][2:]
		} else if len(splittedArg[0]) == 2 { // - and another character
			splittedArg[0] = splittedArg[0][1:]
		}

		if len(splittedArg[0]) == 1 && strings.Contains(p.knownShortArgs, ":"+splittedArg[0]+":") ||
			len(splittedArg[0]) > 1 && strings.Contains(p.knownLongArgs, ":"+splittedArg[0]+":") { // is it a valid short or long argument?

			if len(splittedArg) >= 2 { // long argument
				fmt.Fprintln(writer, "PREDEF:", splittedArg[0], "\t=", splittedArg[1])
				//				predefiningArgs[splittedArg[0][0]] = splittedArg[1]
			} else { // short argument
				newArgs = append(newArgs, string(splittedArg[0][0]))
			}
		} else { // not valid
			fmt.Println("ERROR: Unknown argument", splittedArg[0], "but I'll ignore it :/")
		}
	}
}
