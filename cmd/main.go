package main

import "fmt"

func main() {
	msg := sayHello("Ramsey")
	fmt.Println(msg)
}

func sayHello(name string) string {
	return fmt.Sprintf("Hello %s", name)
}
