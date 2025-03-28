package po_test

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Tom5521/xgotext/internal/util"
	"github.com/Tom5521/xgotext/pkg/po"
	"github.com/Tom5521/xgotext/pkg/po/compiler"
	"github.com/Tom5521/xgotext/pkg/po/parse"
	"github.com/kr/pretty"
)

func TestMergeWithMsgmerge(t *testing.T) {
	t.SkipNow() // this isn't finished yet.
	var err error
	if _, err = exec.LookPath("msgmerge"); err != nil {
		t.Skip("msgmerge isn't in the PATH")
		return
	}

	dir := t.TempDir()
	defPath := filepath.Join(dir, "def.po")
	refPath := filepath.Join(dir, "ref.po")
	outputPath := filepath.Join(dir, "out.po")

	defFile := &po.File{
		Name: defPath,
		Entries: po.Entries{
			{
				ID:  "id1",
				Str: "My translated string",
			},
			{
				ID:  "obsolete string",
				Str: "this is obsolete",
			},
		},
	}
	refFile := &po.File{
		Name: refPath,
		Entries: po.Entries{
			{ID: "id1"},
			{ID: "id2"},
			{ID: "id3"},
		},
	}

	opts := []compiler.PoOption{
		compiler.PoWithForcePo(true),
		compiler.PoWithOmitHeader(true),
	}

	if err = compiler.NewPo(defFile, opts...).ToFile(defPath); err != nil {
		t.Error(err)
		return
	}
	if err = compiler.NewPo(refFile, opts...).ToFile(refPath); err != nil {
		t.Error(err)
		return
	}

	var stderr strings.Builder
	cmd := exec.Command("msgmerge", "-o", outputPath, defPath, refPath)
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		t.Error(stderr.String())
		return
	}

	var parser *parse.PoParser
	parser, err = parse.NewPo(outputPath,
		parse.PoWithSkipHeader(true),
	)
	if err != nil {
		t.Error(err)
		return
	}

	expectedOutput := parser.Parse()
	if err = parser.Error(); err != nil {
		t.Error(err)
		return
	}

	currentOutput := refFile.MergeWithConfig(
		po.NewMergeConfig(po.MergeWithFuzzyMatch(true)),
		// defFile,
	)

	if !util.Equal(expectedOutput.Entries, currentOutput.Entries) {
		for _, d := range pretty.Diff(currentOutput.Entries, expectedOutput.Entries) {
			fmt.Println(d)
		}
		fmt.Println(
			"CURRENT:\n",
			compiler.NewPo(currentOutput, compiler.PoWithOmitHeader(true)).ToString(),
		)
		fmt.Println(
			"EXPECTED:\n",
			compiler.NewPo(expectedOutput, compiler.PoWithIgnoreErrors(true)).ToString(),
		)
		t.Fail()
	}
}
