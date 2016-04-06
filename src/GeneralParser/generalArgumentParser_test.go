package GeneralParser

import (
	"fmt"
	"testing"
)

func TestAddArgument(t *testing.T) {
	p := NewParser()
	len1 := len(p.args)
	p.RegisterArgument("foo", "f", "")
	len2 := len(p.args)

	if len1 == len2 {
		t.Fail()
	}
}

func TestParser(t *testing.T) {

}

func TestArgumentSetValue(t *testing.T) {
	a := argument{
		shortKey: "f",
		longKey:  "foo",
	}
	var v1, v2 *int
	fmt.Println(a)
	v1 = a.Int()
	fmt.Println(a)
	a.set("42")
	v2 = a.Int()
	if v1 == v2 {
		t.Error("The settings failed: both values are equal :(")
	}
}

func TestPointerStuff(t *testing.T) {
	p := NewParser()
	a, _ := p.RegisterArgument("foo", "f", "")
	test(a)
	if *a.stringValue != "foo" {
		t.Fail()
	}
}

func test(a *argument) {
	s := "foo"
	a.stringValue = &s
}
