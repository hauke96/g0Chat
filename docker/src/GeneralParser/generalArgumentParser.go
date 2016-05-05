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
	KnownLongArgs  string
	KnownShortArgs string
	description    string
}

// NewParser creates an empty parser with no arguments.
func NewParser() *parser {
	p := parser{
		args:           make(map[string]*argument),
		longToShortArg: make(map[string]string),
		KnownLongArgs:  ":",
		KnownShortArgs: ":",
		description:    "",
	}
	return &p
}

// Description sets a new description for further information at the help-text.
func (p *parser) Description(newDescription string) {
	p.description = newDescription
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

	p.KnownShortArgs += shortKey + ":"
	p.KnownLongArgs += longKey + ":"
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
	if len(args) == 0 {
		return
	}
	// ------------------------------
	// CREATE OUTPUT WRITER
	// ------------------------------
	writer := new(tabwriter.Writer)
	writer.Init(os.Stdout, 0, 4, 2, ' ', 0)
	defer writer.Flush()

	// ------------------------------
	// CHECK FOR -h/--help
	// ------------------------------
	if args[0] == "--help" || args[0] == "-h" {
		p.showHelp()
		os.Exit(1)
	}

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
	invalidArgExists := false

	// ------------------------------
	// SEPARATE INTO NORMAL AND PREDEFINING
	// ------------------------------
	for _, arg := range args {
		splittedArg := strings.Split(arg, "=")
		key := splittedArg[0]

		// notice: only --foo or -f are allowed, -foo and --f are not allwed!
		// This is just to have the normal feeling of arguments in linux, blame me but i like it ;)
		if key[0:2] == "--" && len(key) > 3 { // 3 because of -- and at least 2 other characters
			key = key[2:]
		} else if len(key) == 2 { // - and one other character
			key = key[1:]
		}

		if len(key) == 1 && strings.Contains(p.KnownShortArgs, ":"+key+":") ||
			len(key) > 1 && strings.Contains(p.KnownLongArgs, ":"+key+":") { // is it a valid short or long argument?

			if len(key) > 1 { // long argument like --foo and not -f
				key = p.longToShortArg[key]
			}

			if len(splittedArg) >= 2 { // argument with value
				fmt.Fprintln(writer, "PREDEF:", key, "\t=", splittedArg[1])
				p.args[key].set(splittedArg[1])
			} else { // argument without value (=flag)
				fmt.Fprintln(writer, "PREDEF:", key, "\t=", true)
				p.args[key].set("true")
			}
		} else { // not valid
			if len(key) == 1 { // just to have the - in from of the argument (a bit more pretty ;) )
				key = "-" + key
			}
			fmt.Println("ERROR: Unknown argument", key, "!")
			invalidArgExists = true
		}
	}
	if invalidArgExists {
		fmt.Println()
		p.showHelp()
		os.Exit(1)
	}
}

func (p *parser) showHelp() {
	writer := new(tabwriter.Writer)
	writer.Init(os.Stdout, 0, 4, 2, ' ', 0)
	defer writer.Flush()

	fmt.Fprintln(writer, "Usage: <argument>=<value> <flag>")
	fmt.Fprintln(writer, "")
	fmt.Fprintln(writer, "Here're all command that are available:\n")

	requiredExists := false

	for _, arg := range p.args {
		fmt.Fprint(writer, "\t-"+arg.shortKey+", --"+arg.longKey+"\t")
		if arg.required {
			fmt.Fprint(writer, "*")
			requiredExists = true
		}
		fmt.Fprintln(writer, "\t"+arg.helpText)
	}

	fmt.Fprintln(writer, "\t-h, --help\t\tShows this help message")
	fmt.Fprintln(writer, "")
	if requiredExists {
		fmt.Fprintln(writer, "Commands with a * are required and have to be set.")
		fmt.Fprintln(writer, "")
	}
	if p.description != "" {
		fmt.Fprintln(writer, p.description)
	}
	fmt.Fprintln(writer, "")
}
