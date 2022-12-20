package cmd

import (
	"testing"
)

func TestHealth(t *testing.T) {
	name := "AppName"
	want := "Health check AppName"

	if got := health(name); got != want {
		t.Errorf("hello() = %q, want %q", got, want)
	}
}
