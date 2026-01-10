package parser

import (
	"fmt"
	"testing"

	"github.com/0xmukesh/coco/internal/lexer"
)

func TestParser(t *testing.T) {
	source := `let one = 1;
let one_point_two = 1.2;
let two_point = 2.;
let name = something;
let not_name = !name;
return -two_point;
let one_plus_two = 1 + 2;
let is_true = true;
let is_false = false;
let is_not_true = !is_true;
let six = (1 + 2) * 2;`
	l := lexer.New(source)
	tks := l.Lex()

	p := New(tks)
	program := p.ParseProgram()
	errs := p.Error()

	if len(errs) != 0 {
		for _, e := range errs {
			fmt.Println(e)
		}
	} else {
		fmt.Println(program.String())
	}
}
