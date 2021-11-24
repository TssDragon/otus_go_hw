package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	str := stringutil.Reverse("Hello, OTUS!")
	fmt.Println(str)
}
