package parse_test

import (
	"fmt"
	"testing"

	. "github.com/Tom5521/xgotext/pkg/po/parse"
)

func TestParse(t *testing.T) {
	input := `#: file:32
msgid "MEOW!"
msgstr "LOL"
msgctxt "WOAS"
msgid "MEOW!"
msgstr "MIAU!"
msgstr[1234] "apples"
"1234"
msgid_plural "a"`

	l := NewLexerFromString(input)

	expectedTokens := []Token{
		{COMMENT, ": file:32"},
		{MSGID, "msgid"},
		{STRING, "MEOW!"},
		{MSGSTR, "msgstr"},
		{STRING, "LOL"},
		{MSGCTXT, "msgctxt"},
		{STRING, "WOAS"},
		{MSGID, "msgid"},
		{STRING, "MEOW!"},
		{MSGSTR, "msgstr"},
		{STRING, "MIAU!"},
		{PluralMsgstr, "msgstr[1234]"},
		{STRING, "apples"},
		{STRING, "1234"},
		{PluralMsgid, "msgid_plural"},
		{STRING, "a"},
	}
	var tokens []Token
	for i, etok := range expectedTokens {
		ctok := l.NextToken()
		tokens = append(tokens, ctok)

		if etok.Literal != ctok.Literal || etok.Type != ctok.Type {
			t.Errorf("unexpected token [%d]:", i)
			t.Error("got:", ctok)
			t.Error("expected:", etok)
			break
		}
	}
	if t.Failed() {
		fmt.Println(tokens)
	}
}