package parse_test

import (
	"fmt"
	"testing"

	"github.com/Tom5521/gotext-tools/internal/util"
	"github.com/Tom5521/gotext-tools/pkg/po"
	"github.com/Tom5521/gotext-tools/pkg/po/compiler"
	"github.com/Tom5521/gotext-tools/pkg/po/parse"
	"github.com/kr/pretty"
)

func TestPoParser(t *testing.T) {
	input := &po.File{
		Entries: po.Entries{
			{
				Flags:    []string{"my-flag lol"},
				Comments: []string{"Hello World"},
				ID:       "Hello", Str: "Hola",
			},
			{Context: "CTX", ID: "MEOW", Str: "MIAU"},
			{
				ID:      "Apple",
				Plural:  "Apples",
				Plurals: po.PluralEntries{{ID: 0, Str: "Manzana"}, {ID: 1, Str: "Manzanas"}},
			},
		},
	}

	compiled := compiler.NewPo(input, compiler.PoWithOmitHeader(true)).ToString()

	parser := parse.NewPoFromString(compiled, "test.po")
	parsed := parser.Parse()

	if parser.Error() != nil {
		t.Error(parser.Error())
		return
	}

	if !util.Equal(parsed.Entries, input.Entries) {
		t.Error("Compiled and parsed differ!")
		fmt.Println("ORIGINAL:\n", compiled)
		fmt.Println("PARSED:\n", compiler.NewPo(parsed, compiler.PoWithOmitHeader(true)).ToString())

		fmt.Println("DIFF:")
		for _, d := range pretty.Diff(parsed.Entries, input.Entries) {
			fmt.Println(d)
		}
	}
}
