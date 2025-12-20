package cmd

import (
	"log"

	"github.com/Tom5521/gotext-tools/v2/pkg/po"
	"github.com/Tom5521/gotext-tools/v2/pkg/po/compile"

	"github.com/spf13/cobra"
)

var (
	mergeCfg    po.MergeConfig
	compilerCfg compile.PoConfig
)

func initConfig(cmd *cobra.Command, args []string) {
	compilerCfg = compile.PoConfig{
		NoLocation:  noLocation,
		AddLocation: compile.PoLocationMode(addLocation),
		WordWrap:    !noWrap,
		ForcePo:     forcePo,
		OmitHeader:  true,
		Verbose:     verbose,
		Logger:      log.Default(),
	}

	mergeCfg = po.MergeConfig{
		FuzzyMatch: !noFuzzyMatching,
		Sort:       true,
	}
}
