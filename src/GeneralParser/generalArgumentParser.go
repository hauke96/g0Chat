package GeneralParser

import (
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

// NewParser creates an empty parser with no arguments.
func NewParser() *parser {
	p := parser{
		args:           make(map[string]*argument),
		longToShortArg: make(map[string]string),
		knownLongArgs:  ":",
		knownShortArgs: ":",
	}
	return &p
}

// RegisterArgument creates a new argument based on the parameters given.
// Now allowed is an empty shortKey and existing long/shortKeys
func (p *parser) RegisterArgument(longKey, shortKey, help string) *argument {
	// the long key is allowed to be empty, the short key not
	if shortKey == "" {
		return nil
	}
	// when the long key and a fitting entry exists, we've a duplicate here
	if longKey != "" && p.longToShortArg[longKey] != "" {
		return nil
	}
	// when there's a short key entry already, we've a duplicate here
	if p.args[shortKey] != nil {
		return nil
	}

	stdString := ""
	stdInt := 0
	stdBool := false

	a := argument{
		longKey:     longKey,
		shortKey:    shortKey,
		helpText:    help,
		stringValue: &stdString,
		intValue:    &stdInt,
		boolValue:   &stdBool,
	}

	p.knownShortArgs += shortKey + ":"
	p.knownLongArgs += longKey + ":"
	p.longToShortArg[longKey] = shortKey

	p.args[shortKey] = &a
	return &a
}

// Parse goes through the arguments (from 1 to n, so the first one is skiped) and sets the values of the arguments
func (p *parser) Parse() {
	p.parseArgs(os.Args[1:])
}

// parseArgs takes the arguments and parses them into two categories:
// normal and predefinings arguments.
// It also evaluates the predefining ones.
func (p *parser) parseArgs(args []string) {
	// ------------------------------
	// CREATE OUTPUT WRITER
	// ------------------------------
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
	//	newArgs = make([]string, 0)
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

			if len(splittedArg[0]) > 1 { // long argument like --foo and not -f
				splittedArg[0] = p.longToShortArg[splittedArg[0]]
			}

			if len(splittedArg) >= 2 { // argument with value
				fmt.Fprintln(writer, "PREDEF:", splittedArg[0], "\t=", splittedArg[1])
				p.args[splittedArg[0]].set(splittedArg[1])
			} else { // argument without value (=flag)
				fmt.Fprintln(writer, "PREDEF:", splittedArg[0], "\t=", true)
				p.args[splittedArg[0]].set("true")
			}
		} else { // not valid
			fmt.Println("ERROR: Unknown argument", splittedArg[0], "but I'll ignore it :/")
		}
	}
}
