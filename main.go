package main

import (
	gorope "Rope/pkg"
	"fmt"
)

func main() {
	rope := gorope.New([]byte("Hello world"))
	fmt.Println(rope)
}
