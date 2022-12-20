package cmd

import "fmt"

func health(name string) string {
	return fmt.Sprintf("Health check %s", name)
}
