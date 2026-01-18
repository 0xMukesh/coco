package driver

import (
	"errors"
	"os"
	"path/filepath"
)

type Source struct {
	Name   string
	Code   []byte
	Exists bool
}

func NewSourceFromFile(file string) (*Source, error) {
	path, err := filepath.Abs(file)
	if err != nil {
		return nil, err
	}

	if filepath.Ext(path) != ".coco" {
		return nil, errors.New("only .coco files are accepted")
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return &Source{
		Name:   path,
		Code:   b,
		Exists: true,
	}, nil
}

func NewDummySource(code string) *Source {
	return &Source{
		Name:   "<dummy>",
		Code:   []byte(code),
		Exists: false,
	}
}
