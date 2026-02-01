package driver

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/0xmukesh/coco/internal/ast"
	"github.com/0xmukesh/coco/internal/codegen"
	"github.com/0xmukesh/coco/internal/lexer"
	"github.com/0xmukesh/coco/internal/parser"
	"github.com/0xmukesh/coco/internal/tokens"
	"github.com/0xmukesh/coco/internal/typechecker"
)

type Driver struct {
	source *Source
}

func NewDriver(src *Source) *Driver {
	return &Driver{
		source: src,
	}
}

func NewDriverFromFile(file string) (*Driver, error) {
	src, err := NewSourceFromFile(file)
	if err != nil {
		return nil, err
	}

	return &Driver{
		source: src,
	}, nil
}

func (d *Driver) Lex() ([]tokens.Token, error) {
	l := lexer.New(string(d.source.Code))
	tks := l.Lex()

	for _, t := range tks {
		if t.Type == tokens.ILLEGAL {
			return nil, fmt.Errorf("failed to lex source - %s", t.Literal)
		}
	}

	return tks, nil
}

func (d *Driver) Parse(tks []tokens.Token) (*ast.Program, error) {
	p := parser.New(tks)
	program := p.ParseProgram()
	if p.HasErrors() {
		return nil, errors.New(strings.Join(p.Errors(), "\n"))
	}

	return program, nil
}

func (d *Driver) TypeCheck(program *ast.Program) error {
	tc := typechecker.New()
	tc.Transform(program)
	if tc.HasErrors() {
		return errors.Join(tc.Errors()...)
	}

	return nil
}

func (d *Driver) Codegen(program *ast.Program) (string, error) {
	cg := codegen.New()
	cg.Generate(program)
	if cg.HasErrors() {
		return "", errors.Join(cg.Errors()...)
	}

	ir := cg.EmitIR()
	return ir, nil
}

func (d *Driver) IrFileToBinary(irFilePath string, outFilePath string, deleteIrAtEnd bool) error {
	irFilePath, err := filepath.Abs(irFilePath)
	if err != nil {
		return err
	}

	if deleteIrAtEnd {
		defer os.Remove(irFilePath)
	}

	if filepath.Ext(irFilePath) != ".ll" {
		return errors.New("only .ll files are accepted")
	}

	if outFilePath == "" {
		outFilePath = strings.Replace(irFilePath, ".ll", "", 1)
	}

	cmd := exec.Command("clang", "-O2", irFilePath, "-o", outFilePath)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (d *Driver) Pipeline(outFilePath string, emitIr bool) error {
	outFilePath, err := filepath.Abs(outFilePath)
	if err != nil {
		return err
	}

	tks, err := d.Lex()
	if err != nil {
		return err
	}

	ast, err := d.Parse(tks)
	if err != nil {
		return err
	}

	if err := d.TypeCheck(ast); err != nil {
		return err
	}

	ir, err := d.Codegen(ast)
	if err != nil {
		return err
	}

	var irFilePath string
	if outFilePath == "" {
		irFilePath = filepath.Join(os.TempDir(), "tmp.ll")
	} else {
		if strings.Contains(outFilePath, ".") {
			irFilePath = strings.Split(outFilePath, ".")[0] + ".ll"
		} else {
			irFilePath = outFilePath + ".ll"
		}
	}

	if err := os.WriteFile(irFilePath, []byte(ir), 0777); err != nil {
		return err
	}

	if err := d.IrFileToBinary(irFilePath, outFilePath, !emitIr); err != nil {
		return err
	}

	return nil
}
