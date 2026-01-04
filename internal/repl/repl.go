package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/0xmukesh/coco/internal/eval"
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
				io.WriteString(out, fmt.Sprintf("%+v\n", tok))
			}
		case "parse":
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) != 0 {
				for _, msg := range p.Errors() {
					io.WriteString(out, "ERR: "+msg+"\n")
				}
				continue
			}

			io.WriteString(out, program.String()+"\n")
		case "eval":
			p := parser.New(l)
			program := p.ParseProgram()
			res := eval.Eval(program)

			io.WriteString(out, res.Inspect()+"\n")
		default:
			io.WriteString(out, "invalid mode")
		}

	}
}
