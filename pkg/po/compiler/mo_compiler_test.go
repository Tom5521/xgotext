package compiler_test

import (
	"bytes"
	"fmt"
	"os/exec"
	"testing"

	"github.com/Tom5521/xgotext/internal/util"
	"github.com/Tom5521/xgotext/pkg/po"
	"github.com/Tom5521/xgotext/pkg/po/compiler"
	"github.com/Tom5521/xgotext/pkg/po/parse"
	"github.com/kr/pretty"
)

func TestMoCompiler(t *testing.T) {
	input := po.Entries{
		{Context: "My context :3", ID: "id1", Str: "HELLO"},
		{
			ID:     "id2",
			Plural: "helooows",
			Plurals: po.PluralEntries{
				po.PluralEntry{ID: 0, Str: "Holanda"},
				po.PluralEntry{ID: 1, Str: "Holandas"},
			},
		},
		{ID: "id3", Str: "Hello3"},
	}

	c := compiler.NewMo(&po.File{Entries: input})

	parser, err := parse.NewMoFromBytes(c.ToBytes(), "test.mo")
	if err != nil {
		t.Error(err)
		return
	}

	parsedFile := parser.Parse()
	if len(parser.Errors()) > 0 {
		t.Error(parser.Errors()[0])
		return
	}
	parsed := parsedFile.Entries
	if !util.Equal(parsed, input) {
		t.Error("Sended and parsed differ!")
		t.Logf("SENDED:\n%v", input)
		t.Logf("PARSED:\n%v", parsed)
		t.Log("DIFF:")

		for _, d := range pretty.Diff(parsed, input) {
			fmt.Println(d)
		}
		return
	}
}

func TestMoCompilerWithMsgunfmt(t *testing.T) {
	input := po.Entries{
		{Context: "My context :3", ID: "id1", Str: "HELLO"},
		{
			ID:     "id2",
			Plural: "helooows",
			Plurals: po.PluralEntries{
				po.PluralEntry{ID: 0, Str: "Holanda"},
				po.PluralEntry{ID: 1, Str: "Holandas"},
			},
		},
		{ID: "id3", Str: "Hello3"},
	}

	c := compiler.NewMo(
		&po.File{Entries: input},
	)

	_, err := exec.LookPath("msgunfmt")
	if err != nil {
		t.Skip("gettext-tools isn't in the PATH")
		return
	}

	var stdout bytes.Buffer

	cmd := exec.Command("msgunfmt")
	cmd.Stdin = bytes.NewReader(c.ToBytes())
	cmd.Stdout = &stdout
	cmd.Stderr = &stdout

	err = cmd.Run()
	if err != nil {
		t.Error(stdout.String())
		return
	}

	parser, _ := parse.NewPoFromReader(&stdout, "lol.po",
		parse.PoWithSkipHeader(true),
	)
	parsedFile := parser.Parse()
	if err = parser.Error(); err != nil {
		t.Error(err)
		return
	}
	if !util.Equal(parsedFile.Entries, input) {
		t.Fail()
	}
}
