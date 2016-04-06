package GeneralParser

import (
	"strconv"
)

type argument struct {
	longKey     string
	shortKey    string
	helpText    string
	required    bool
	stringValue *string
	intValue    *int
	boolValue   *bool
}

// String defines this argument as an argument that contains a string.
// After calling 'Parse' on the parser this argument belongs to, the value will be set
func (a *argument) String() *string {
	v := new(string)
	a.stringValue = v
	return v
}

// String defines this argument as an argument that contains an integer.
// After calling 'Parse' on the parser this argument belongs to, the value will be set
func (a *argument) Int() *int {
	v := new(int)
	a.intValue = v
	return v
}

// String defines this argument as an argument that contains a boolen.
// After calling 'Parse' on the parser this argument belongs to, the value will be set
func (a *argument) Bool() *bool {
	v := new(bool)
	a.boolValue = v
	return v
}

// Help sets/redefines the help-message.
func (a *argument) Help(text string) *argument {
	a.helpText = text
	return a
}

// Required defines this argument as a required one. Skipping this argument by executing the programm will cause an error and a user notification.
func (a *argument) Required() *argument {
	a.required = true
	return a
}

func (a *argument) set(value string) {
	intValue, err := strconv.Atoi(value)
	if err == nil {
		*a.intValue = intValue
	}

	boolValue, err := strconv.ParseBool(value)
	if err == nil {
		*a.boolValue = boolValue
		return
	}

	*a.stringValue = value
}
