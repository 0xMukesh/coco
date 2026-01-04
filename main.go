package main

import (
	"fmt"
	"os"

	"github.com/0xmukesh/coco/cli"
)

func main() {
	if err := cli.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "an error occurred while executing program - %v", err)
	}
}
