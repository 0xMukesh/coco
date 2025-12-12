package main

import (
	"os"

	"github.com/0xmukesh/coco/internal/repl"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
