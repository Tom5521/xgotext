package cmd

import (
	"errors"
	"math"

	"github.com/spf13/cobra"

	"github.com/Tom5521/gotext-tools/v2/pkg/po"
	"github.com/Tom5521/gotext-tools/v2/pkg/po/compile"
)

var (
	filesFrom string
	directory string
	output    string
	lessThan  uint
	moreThan  uint
	unique    bool
	useFirst  bool
	lang      string
	color     string

	// forcePo     bool

	noLocation  bool
	addLocation string
	noWrap      bool
	sortOutput  bool
	sortByFile  bool
)

func init() {
	flags := root.Flags()

	flags.StringVarP(&filesFrom, "files-from", "f", "",
		`get list of input files from`,
	)
	flags.StringVarP(&directory, "directory", "D", "",
		`add DIRECTORY to list for input files search
If input file is -, standard input is read.`,
	)
	flags.StringVarP(&output, "output-file", "o", "-",
		`write output to specified file
The results are written to standard output if no output file is specified
or if it is -.`,
	)
	flags.UintVarP(
		&lessThan,
		"less-than",
		"<",
		math.MaxUint,
		`print messages with less than this many
definitions, defaults to infinite if not set`,
	)
	flags.UintVarP(&moreThan, "more-than", ">", 0,
		`print messages with more than this many
definitions, defaults to 0 if not set`,
	)

	flags.BoolVarP(&unique, "unique", "u", false, `shorthand for --less-than=2, requests
that only unique messages be printed`)

	flags.BoolVar(&useFirst, "use-first", false, `use first available translation for each
message, don't merge several translations`)

	flags.StringVar(&lang, "lang", "", `set 'Language' field in the header entry`)
	flags.StringVar(
		&color,
		"color",
		"auto",
		`use colors and other text attributes, may be 'always', 'never', 'auto'`,
	)
	// flags.BoolVar(&forcePo, "force-po", false, `write PO file even if empty`)
	flags.BoolVar(&noLocation, "no-location", false, `do not write '#: filename:line' lines`)
	flags.StringVar(
		&addLocation,
		"add-location",
		"full",
		`generate '#: filename:line' lines (default)`,
	)
	flags.BoolVar(&noWrap, "no-wrap", false, `do not break long message lines, longer than
the output page width, into several lines`)
	flags.BoolVarP(&sortOutput, "sort-output", "s", false, `generate sorted output`)
	flags.BoolVarP(&sortByFile, "sort-by-file", "F", false, `sort output by file location`)

	root.MarkFlagsMutuallyExclusive("unique", "less-than")
	root.MarkFlagsMutuallyExclusive("no-location", "add-location")
	root.MarkFlagsMutuallyExclusive("sort-output", "sort-by-file")
}

var (
	mergeCfg    po.MsgcatMergeConfig
	compilerCfg = compile.DefaultPoConfig()
)

func initCfg(cmd *cobra.Command, args []string) error {
	if lessThan <= 1 {
		return errors.New("impossible selection criteria specified")
	}

	mergeCfg = po.MsgcatMergeConfig{
		LessThan: lessThan,
		MoreThan: moreThan,
		UseFirst: useFirst,
	}

	if unique {
		mergeCfg.LessThan = 2
	}

	compilerCfg.ApplyOptions(
		compile.PoWithWordWrap(!noWrap),
		compile.PoWithNoLocation(noLocation),
	)

	switch addLocation {
	case string(compile.PoLocationModeFile),
		string(compile.PoLocationModeFull),
		string(compile.PoLocationModeNever):
		compilerCfg.AddLocation = compile.PoLocationMode(addLocation)
	}

	return nil
}
