package main

import "testing"

func Test_sayHello(t *testing.T) {
	name := "Davy"
	want := "Hello Davy"

	if got := sayHello(name); got != want {
		t.Errorf("hello() = %q, want %q", got, want)
	}
}
