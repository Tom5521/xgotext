package parse_test

import (
	"testing"

	"github.com/Tom5521/xgotext/internal/util"
	"github.com/Tom5521/xgotext/pkg/go/parse"
	"github.com/Tom5521/xgotext/pkg/po"
	"github.com/kr/pretty"
)

func TestParse(t *testing.T) {
	const input = `package main
import "github.com/leonelquinteros/gotext"

func main(){
	gotext.Get("Hello World!")
}`

	expected := po.Entries{
		{
			ID: "Hello World!",
			Locations: []po.Location{
				{
					Line: 5,
					File: "test.go",
				},
			},
		},
	}
	parser, err := parse.NewParserFromString(input, "test.go")
	if err != nil {
		t.Error(err)
	}

	file := parser.Parse()
	if len(parser.Errors()) > 0 {
		t.Error(parser.Errors()[0])
	}

	if !po.EqualEntries(file.Entries[1:], expected) {
		t.Log("Unexpected entries slice")
		t.Log("got:", file.Entries)
		t.Log("expected:", expected)
		t.FailNow()
	}
}

func TestExtractAll(t *testing.T) {
	const input = `package main

import "github.com/leonelquinteros/gotext"

func main(){
	_ = "Hello World"
	a := "Hi world"
	b := "I love onions!"
	
	var eggs string = "sugar"
}`

	parser, err := parse.NewParserFromString(
		input,
		"test.go",
		parse.WithExtractAll(true),
	)
	if err != nil {
		t.Error(err)
	}

	file := parser.Parse()
	if len(parser.Errors()) > 0 {
		t.Error(parser.Errors()[0])
	}

	expected := po.Entries{
		{
			ID: "Hello World",
			Locations: []po.Location{
				{
					File: "test.go",
					Line: 6,
				},
			},
		},
		{
			ID: "Hi world",
			Locations: []po.Location{
				{
					File: "test.go",
					Line: 7,
				},
			},
		},
		{
			ID: "I love onions!",
			Locations: []po.Location{
				{
					File: "test.go",
					Line: 8,
				},
			},
		},
		{
			ID: "sugar",
			Locations: []po.Location{
				{
					File: "test.go",
					Line: 10,
				},
			},
		},
	}

	if !util.Equal(file.Entries[1:], expected) {
		t.Error("Unexpected translation")
		t.Log("got:", file.Entries)
		t.Log("expected:", expected)
		t.Log("DIFF:")
		for _, d := range pretty.Diff(file.Entries, expected) {
			t.Log(d)
		}
	}
}

func BenchmarkParse(b *testing.B) {
	const input = `package main

import "github.com/leonelquinteros/gotext"

func main(){
	_ = "Hello World"
	a := "Hi world"
	b := "I love onions!"
	
	var eggs string = "sugar"
}`
	bytes := []byte(input)
	parser, err := parse.NewParserFromBytes(bytes, "test.go")
	if err != nil {
		b.Error(err)
	}

	tests := []struct {
		name    string
		options []parse.Option
	}{
		{name: "normal"},
		{name: "extract-all", options: []parse.Option{parse.WithExtractAll(true)}},
	}

	for _, t := range tests {
		b.Run(t.name, func(b *testing.B) {
			for range b.N {
				parser.Parse(t.options...)
				if len(parser.Errors()) > 0 {
					b.Error(parser.Errors()[0])
				}
			}
		})
	}
}
