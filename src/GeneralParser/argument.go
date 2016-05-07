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
	if a.stringValue == nil {
		a.stringValue = new(string)
	}
	return a.stringValue
}

// String defines this argument as an argument that contains an integer.
// After calling 'Parse' on the parser this argument belongs to, the value will be set
func (a *argument) Int() *int {
	if a.intValue == nil {
		a.intValue = new(int)
	}
	return a.intValue
}

// String defines this argument as an argument that contains a boolen.
// After calling 'Parse' on the parser this argument belongs to, the value will be set
func (a *argument) Bool() *bool {
	if a.boolValue == nil {
		a.boolValue = new(bool)
	}
	return a.boolValue
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

// Default sets a default value to this argument
func (a *argument) Default(value string) *argument {
	a.set(value)
	return a
}

// set sets all fitting fields (the one for int, bool and string) to the given value.
// If it's a string, the string field is set and if the string contains a number only, the int field is also set
// If it's a bool, only the bool field is set
// If it's an int, the int and string field is set.
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
