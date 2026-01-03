package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/0xmukesh/coco/internal/lexer"
	"github.com/0xmukesh/coco/internal/parser"
	"github.com/0xmukesh/coco/internal/tokens"
)

func Start(in io.Reader, out io.Writer, prompt string, mode string) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(prompt)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == "exit" {
			break
		}

		l := lexer.New(line)

		switch mode {
		case "lex":
			for tok := l.NextToken(); tok.Type != tokens.EOF; tok = l.NextToken() {
				fmt.Printf("%+v\n", tok)
			}
		case "parse":
			p := parser.New(l)
			program := p.ParseProgram()
			fmt.Printf("%s\n", program.String())
		default:
			fmt.Printf("invalid mode")
		}

	}
}
