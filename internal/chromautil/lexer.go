//go:build ignore

package chromautil

import (
	. "github.com/alecthomas/chroma/v2"
	. "github.com/alecthomas/chroma/v2/lexers"
)

var PoLexer = Register(MustNewLexer(
	&Config{
		Name:      "Po",
		Aliases:   []string{"Portable Object", "Gettext", "po", "pot", "PO", "POT", "Pot"},
		Filenames: []string{"*.po", "*.pot"},
	},
	poRules,
),
)

// TODO: Finish this.
func poRules() Rules {
	return Rules{
		"root": {
			{`\n`, TextWhitespace, nil},
			{`\s+`, TextWhitespace, nil},
			{`#\s+[^\n\r]*`, CommentSingle, nil},
			{`(msgid|msgstr|msgid_plural)\b`, KeywordDeclaration, nil},
			{`"(\\\\|\\"|[^"])*"`, LiteralString, nil},
			{`[\[\]]`, Punctuation, nil},
			{`\d+`, LiteralNumber, nil},
		},
	}
}
