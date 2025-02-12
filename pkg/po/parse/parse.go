// Package parse provides the functionality to parse PO (Portable Object)
// files from various sources, such as files, byte slices, strings, and readers.
// It integrates with the `ast`, `generator`, and `types`
// packages to normalize and generate structured representations of PO files.
//
// Key Features:
// - Parses PO files from different input sources (file paths, strings, byte slices, and readers).
// - Normalizes PO entries using the `ast.Normalizer`.
// - Generates a structured `types.File` representation of the PO file.
// - Handles errors and warnings during parsing and normalization.
//
// Example Usage:
//
//	cfg := parsers.Config{}
//	p, err := parse.NewParser("example.po", cfg)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	file := p.Parse()
//	if len(p.Errors()) > 0 {
//	    log.Fatal(p.Errors()[0])
//	}
//	fmt.Println(file)
//
// For more details, refer to the individual functions and types.
package parse

import (
	"io"
	"os"

	"github.com/Tom5521/xgotext/pkg/po/parse/ast"
	"github.com/Tom5521/xgotext/pkg/po/parse/generator"
	"github.com/Tom5521/xgotext/pkg/po/types"
)

type Parser struct {
	norm *ast.Normalizer

	warns  []string
	errors []error
}

func NewParser(path string) (*Parser, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return NewParserFromBytes(file, path)
}

func NewParserFromReader(r io.Reader, name string) (*Parser, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return NewParserFromBytes(data, name)
}

func NewParserFromFile(f *os.File) (*Parser, error) {
	return NewParserFromReader(f, f.Name())
}

func NewParserFromString(s, name string) (*Parser, error) {
	return NewParserFromBytes([]byte(s), name)
}

func NewParserFromBytes(data []byte, name string) (*Parser, error) {
	p := &Parser{}
	p.processpath(data, name)
	if len(p.errors) > 0 {
		return nil, p.errors[0]
	}
	return p, nil
}

func (p Parser) Errors() []error {
	return p.errors
}

func (p Parser) Warnings() []string {
	return p.warns
}

func (p *Parser) processpath(content []byte, path string) {
	norm, errs := ast.NewParser(content, path).Normalizer()
	if len(errs) > 0 {
		p.errors = append(p.errors, errs...)
		return
	}

	p.norm = norm
}

func (p *Parser) Parse() *types.File {
	p.norm.Normalize()
	p.warns = append(p.warns, p.norm.Warnings()...)
	if len(p.norm.Errors()) > 0 {
		p.errors = append(p.errors, p.norm.Errors()...)
		return nil
	}

	g := generator.New(p.norm.File())

	file := g.Generate()
	if len(g.Errors()) > 0 {
		p.errors = append(p.errors, g.Errors()...)
		return nil
	}

	return file
}
