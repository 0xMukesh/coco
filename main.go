package main

import "fmt"

func main() {
	fmt.Println(len("%d\n\x00"))
}
