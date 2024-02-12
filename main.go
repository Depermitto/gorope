package main

import (
	"Rope/pkg/gorope"
	"fmt"
)

func main() {
	rope := gorope.New([]byte("Hello world"))
	fmt.Println(rope)
}
