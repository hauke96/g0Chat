package GeneralParser

import (
	"testing"
)

func TestAddArgument(t *testing.T) {
	p := NewParser()
	len1 := len(p.args)
	p.RegisterArgument("foo", "f", "")
	len2 := len(p.args)

	if len1 == len2 {
		t.Error("The argument has not been added.")
	}
}

func TestParser(t *testing.T) {
	p := NewParser()

	integerArg := p.RegisterArgument("foo", "f", "A sample argument").Required().Int()
	stringArg := p.RegisterArgument("bar", "b", "A second sample argument").String()
	boolArg := p.RegisterArgument("", "s", "And another sample argument").Bool()
	flagArg := p.RegisterArgument("flag", "l", "A sample flag").Bool()

	args := make([]string, 4)
	args[0] = "-f=123"
	args[1] = "--bar=hello world!"
	args[2] = "-s=true"
	args[3] = "-l"

	p.parseArgs(args)

	if integerArg == nil || *integerArg != 123 {
		t.Error("The integerArg has not been parsed correctly.")
	}
	if stringArg == nil || *stringArg != "hello world!" {
		t.Error("The stringArg has not been parsed correctly.")
	}
	if boolArg == nil || *boolArg != true {
		t.Error("The boolArg has not been parsed correctly.")
	}
	if flagArg == nil || *flagArg != true {
		t.Error("The flagArg has not been parsed correctly.")
	}
}

func TestArgumentSetValue(t *testing.T) {
	a := NewParser().RegisterArgument("foo", "f", "")
	var v1, v2 int

	v1 = *a.intValue
	a.set("42")
	v2 = *a.intValue

	if v1 == v2 {
		t.Error("The settings failed: both values v1(", v1, ",", &v1, ") and v2(", v2, ",", &v2, ") are equal :(")
	}
}
