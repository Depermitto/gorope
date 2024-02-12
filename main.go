package main

import (
	"Rope/gorope"
	"fmt"
)

func main() {
	rope := gorope.New([]byte("Hello world"))
	fmt.Println(rope)
}
