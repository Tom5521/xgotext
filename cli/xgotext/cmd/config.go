package cmd

import (
	goparse "github.com/Tom5521/gotext-tools/v2/pkg/go/parse"
	"github.com/Tom5521/gotext-tools/v2/pkg/po"
	"github.com/Tom5521/gotext-tools/v2/pkg/po/compile"
	poparse "github.com/Tom5521/gotext-tools/v2/pkg/po/parse"
)

var (
	PoParserCfg poparse.PoConfig
	GoParserCfg goparse.Config
	CompilerCfg compile.PoConfig
	HeadersCfg  = po.DefaultTemplateHeaderConfig()
)

func initConfig() {
	HeadersCfg.Nplurals = nplurals
	HeadersCfg.ProjectIDVersion = packageVersion
	HeadersCfg.ReportMsgidBugsTo = msgidBugsAddress
	HeadersCfg.Language = lang

	GoParserCfg = goparse.Config{
		CleanDuplicates: true,
		Exclude:         exclude,
		ExtractAll:      extractAll,
		HeaderConfig:    &HeadersCfg,
		Logger:          logger,
		Verbose:         verbose,
	}
	CompilerCfg = compile.PoConfig{
		Logger:          logger,
		ForcePo:         forcePo,
		OmitHeader:      omitHeader,
		PackageName:     packageName,
		CopyrightHolder: copyrightHolder,
		ForeignUser:     foreignUser,
		Title:           title,
		NoLocation:      noLocation,
		AddLocation:     compile.PoLocationMode(addLocation),
		MsgstrPrefix:    msgstrPrefix,
		MsgstrSuffix:    msgstrSuffix,
		Verbose:         verbose,
		HeaderComments:  true,
		HeaderFields:    true,
		WordWrap:        wordWrap,
	}

	PoParserCfg = poparse.PoConfig{
		Logger: logger,
	}
}
