package cmd

import "fmt"

func health(name string) string {
	return fmt.Sprintf("health check %s", name)
}
