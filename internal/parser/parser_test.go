package parser

import (
	"fmt"
	"testing"

	"github.com/0xmukesh/coco/internal/lexer"
)

// FIXME: setup proper unit tests for parser
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
let six = (1 + 2) * 2;
let seven = ++six;
six += 3;

{
	let six = 6;
	{
	let seven = 7;
	}
}

let if_true = if (is_true) {
	return 1 + 1;
} else if (2 + 2 == 4) {
	return 4;
} else {
	return 2 + 2;
}

while (is_true) {
	1 + 1;
}

for (let a = 1; a >= 4; a) {}
for (let a = 1; ;a) {}
for (let a = 1;;) {}
for (; a >= 4; a) {}
for (; a >= 4;) {}
for (; ; a) {}
2 + 3 + 4
2 ** 3 ** 4

let something = fn (x, y, z) {
	return x + y + z;
}

something(1, 2, 3);

exit 1;`

	l := lexer.New(source)
	tks := l.Lex()

	p := New(tks)
	program := p.ParseProgram()

	if p.HasErrors() {
		for _, e := range p.Errors() {
			fmt.Println(e)
		}
	}

	fmt.Println(program.String())
}
