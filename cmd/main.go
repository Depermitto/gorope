package main

import (
	"fmt"
	"github.com/Depermitto/gorope"
)

func main() {
	rope := gorope.FromString("Hello world")
	fmt.Println(rope)
}
