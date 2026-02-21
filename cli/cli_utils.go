package cli

import (
	"fmt"
	"os"
)

func printErrAndExit(s string) {
	fmt.Fprintf(os.Stderr, "%s\n", s)
	os.Exit(1)
}
