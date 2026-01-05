package cli

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/0xmukesh/coco/internal/driver"
	"github.com/chzyer/readline"
)

func Start() error {
	filePath := flag.String("file", "", "path of source file")
	flag.Parse()

	if *filePath != "" {
		return executeFromFile(*filePath)
	} else {
		return startRepl()
	}
}

func executeFromFile(filePath string) error {
	source, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file at %s path", filePath)
	}

	d := driver.NewDriver(string(source))
	res, err := d.Process()
	if err != nil {
		return err
	}

	fmt.Println(res.Inspect())
	return nil
}

func startRepl() error {
	fmt.Println("coco repl")
	rl, err := readline.New("~/ ")
	if err != nil {
		return fmt.Errorf("failed to start repl - %v", err)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return fmt.Errorf("failed to scan line - %v", err)
		}

		if line == "exit" {
			return nil
		}

		d := driver.NewDriver(line)
		res, err := d.Process()
		if err != nil {
			fmt.Printf("error: %v\n", err)
			continue
		}

		fmt.Println(res.Inspect())
		continue
	}
}
