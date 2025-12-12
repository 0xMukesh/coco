package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/0xmukesh/coco/internal/lexer"
)

func Start(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)

	for {
		fmt.Printf(">> ")
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == "exit" {
			break
		}

		l := lexer.New(line)
		tks := l.Tokenize()

		for _, t := range tks {
			fmt.Printf("%+v\n", t)
		}
	}
}
