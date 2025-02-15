package lexer_test

import (
	"fmt"
	"testing"

	"github.com/Tom5521/xgotext/internal/util"
	. "github.com/Tom5521/xgotext/pkg/po/parse/lexer"
	. "github.com/Tom5521/xgotext/pkg/po/parse/token"
)

func TestLexer(t *testing.T) {
	input := `#: file:32
msgid "MEOW!"
msgstr "LOL"
msgctxt "WOAS"
msgid "MEOW!"
msgstr "MIAU!"
msgstr[1234] "apples"
"1234"
msgid_plural "a"`

	l := NewFromString(input)

	expectedTokens := []Token{
		{COMMENT, "#: file:32", 0},
		{MSGID, "msgid", 11},
		{STRING, "MEOW!", 17},
		{MSGSTR, "msgstr", 25},
		{STRING, "LOL", 32},
		{MSGCTXT, "msgctxt", 38},
		{STRING, "WOAS", 46},
		{MSGID, "msgid", 53},
		{STRING, "MEOW!", 59},
		{MSGSTR, "msgstr", 67},
		{STRING, "MIAU!", 74},
		{PluralMsgstr, "msgstr[1234]", 82},
		{STRING, "apples", 95},
		{STRING, "1234", 104},
		{PluralMsgid, "msgid_plural", 111},
		{STRING, "a", 124},
	}
	var tokens []Token
	for i, etok := range expectedTokens {
		ctok := l.NextToken()
		tokens = append(tokens, ctok)

		if etok.Type == STRING {
			etok.Literal = `"` + etok.Literal + `"`
		}

		if !util.Equal(etok, ctok) {
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

func BenchmarkLexer(b *testing.B) {
	l := NewFromString(`#: file:32
msgid "MEOW!"
msgstr "LOL"
msgctxt "WOAS"
msgid "MEOW!"
msgstr "MIAU!"
msgstr[1234] "apples"
"1234"
msgid_plural "a"`)

	for range b.N {
		tok := l.NextToken()
		for tok.Type != EOF {
			tok = l.NextToken()
		}
	}
}
