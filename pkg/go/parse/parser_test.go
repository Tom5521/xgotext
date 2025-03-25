package parse_test

import (
	"fmt"
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
	parser, err := parse.NewParserFromString(input, "test.go", parse.WithNoHeader(true))
	if err != nil {
		t.Error(err)
	}

	file := parser.Parse()
	if err = parser.Error(); err != nil {
		t.Error(err)
	}

	if !file.Entries.Equal(expected) {
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
		parse.WithNoHeader(true),
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

	if !file.Entries.Equal(expected) {
		t.Error("Unexpected translation")
		t.Log("got:", file.Entries)
		t.Log("expected:", expected)
		t.Log("DIFF:")
		for _, d := range pretty.Diff(file.Entries, expected) {
			t.Log(d)
		}
	}
}

func TestExtractAll2(t *testing.T) {
	const input = `package main

import "github.com/leonelquinteros/gotext"

func main(){
	_ = "Hello World"
	a := "Hi world"
	b := "I love onions!"
	
	var eggs string = "sugar"
	const asadasd = "asad"

	for i := range "Hello World From a Loop!"{
		fmt.Println("Hello!")
	}
	for i := "hello";i != "from";i += "another loop"{

	}

	a := struct{x string}{"Hello from a struct!"}

	if "Hello From an if" != "Bye from an if"{
		print("no")
	}

	"hello from an anonymous string!"

	a := make(map[string]string)
	a["Hello from a key"] = "Hello from a value"
}`

	parser, _ := parse.NewParserFromString(
		input,
		"lol",
		parse.WithExtractAll(true),
		parse.WithNoHeader(true),
	)

	file := parser.Parse()
	if err := parser.Error(); err != nil {
		t.Error(err)
		return
	}

	expectedIDs := []string{
		"Hello World",
		"Hi world",
		"I love onions!",
		"sugar",
		"asad",
		"Hello World From a Loop!",
		"Hello!",
		"hello",
		"from",
		"another loop",
		"Hello from a struct!",
		"Hello From an if",
		"Bye from an if",
		"no",
		"hello from an anonymous string!",
		"Hello from a key",
		"Hello from a value",
	}

	var ids []string
	for _, e := range file.Entries {
		ids = append(ids, e.ID)
	}
	if !util.Equal(expectedIDs, ids) {
		t.Error("Unexpected ids slice")
		for _, d := range pretty.Diff(expectedIDs, ids) {
			fmt.Println(d)
		}
	}
}
