package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/0xmukesh/coco/internal/constants"
	"github.com/0xmukesh/coco/internal/repl"
)

func main() {
	mode := flag.String("mode", constants.REPL_MODE_PARSE, "mode to run the repl on")
	flag.Parse()

	if *mode != constants.REPL_MODE_LEX && *mode != constants.REPL_MODE_PARSE && *mode != constants.REPL_MODE_EVAL {
		fmt.Fprintf(os.Stderr, "invalid repl mode")
		os.Exit(1)
	}

	repl.Start(os.Stdin, os.Stdout, "~/ ", *mode)
}
