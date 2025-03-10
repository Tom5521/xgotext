package parse_test

import (
	"fmt"
	"testing"

	"github.com/Tom5521/xgotext/internal/util"
	"github.com/Tom5521/xgotext/pkg/po"
	"github.com/Tom5521/xgotext/pkg/po/compiler"
	"github.com/Tom5521/xgotext/pkg/po/parse"
	"github.com/kr/pretty"
)

func TestMoParse(t *testing.T) {
	entries := po.Entries{
		{
			ID:      "Apple",
			Context: "USA",
			Plural:  "Apples",
			Plurals: po.PluralEntries{
				{0, "Manzana"},
				{1, "Manzanas"},
			},
		},
		{ID: "Hi", Str: "Hola", Context: "casual"},
		{ID: "", Str: ""},
		{ID: "How are you?", Str: "Como estás?"},
	}

	com := compiler.NewMo(&po.File{Entries: entries})
	moFile := com.ToBytes()

	parser, err := parse.NewMoFromBytes(moFile, "test.mo")
	if err != nil {
		t.Error(err)
		return
	}

	parsedEntries := parser.Parse().Entries
	if parser.Error() != nil {
		t.Error(parser.Error())
		return
	}

	if !util.Equal(entries, parsedEntries) {
		t.Error("Parsed entries differ!")
		fmt.Println("--- ORIGINAL:", entries)
		fmt.Println("--- PARSED:", parsedEntries)
		fmt.Println("--- DIFF:")
		for _, d := range pretty.Diff(entries, parsedEntries) {
			fmt.Println(d)
		}
		return
	}
}
