package greetings

import (
    "testing"
)

func TestHelloName(t *testing.T) {
    name := "Imam"
    want := "Hi, Imam. Welcome!"
    msg, err := Hello(name)
    if msg != want || err != nil {
        t.Fatalf(`Hello("Imam") = %q, %v, want %q, nil`, msg, err, want)
    }
}

func TestHelloEmpty(t *testing.T) {
    msg, err := Hello("")
    if msg != "" || err == nil {
        t.Fatalf(`Hello("") = %q, %v, want "", error`, msg, err)
    }
}
