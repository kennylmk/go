package greetings

import (
	"regexp"
	"testing"
)

func TestHelloName(t *testing.T) {
	name := "Greet"
	want := regexp.MustCompile(`\b` + name + `\b`)
	msg, err := Hello("GREET")
	if !want.MatchString(msg) || err != nil {
		t.Fatalf(`Hello("GREET) = %q, %v, want watch for %#q, nill`, msg, err, want)

	}
}

func TestHelloEmpty(t *testing.T) {
	msg, err := Hello("")
	if msg != "" || err == nil {
		t.Fatalf(`Hello("")=%q, %v, want "",error`, msg, err)
	}
}
