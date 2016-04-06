package GeneralParser

import (
	"errors"
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

func (a *argument) String() *string {
	return a.stringValue
}

func (a *argument) Int() *int {
	return a.intValue
}

func (a *argument) Bool() *bool {
	return a.boolValue
}

func (a *argument) Help(text string) *argument {
	a.helpText = text
	return a
}

func (a *argument) Required() *argument {
	a.required = true
	return a
}

func (a *argument) set(value string) error {
	a.stringValue = &value

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return errors.New("\tThe argument " + value + " is not an integer!")
	}
	*a.intValue = intValue

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return errors.New("\tThe argument " + value + " is not a boolean!")
	}
	*a.boolValue = boolValue

	return nil
}
