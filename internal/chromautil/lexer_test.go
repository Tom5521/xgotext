//go:build ignore

package chromautil_test

import (
	"os"
	"testing"

	"github.com/Tom5521/gotext-tools/v2/internal/chromautil"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/quick"
	"github.com/alecthomas/chroma/v2/styles"
)

func TestLOL(t *testing.T) {
	code := `msgid "Lol"
msgstr "waos"
`

	l := chromautil.PoLexer

	s := styles.Get("monokai")
	it, err := l.Tokenise(nil, code)
	if err != nil {
		t.Log(err)
		return
	}

	f := formatters.Fallback
	err = f.Format(os.Stdout, s, it)
	if err != nil {
		t.Log(err)
		return
	}
}

func Test2(t *testing.T) {
	code := `package main
import "fmt"
func main() {
    fmt.Println("Hello, World!")
}`

	// Highlight to terminal with a specific style (e.g., "monokai")
	err := quick.Highlight(os.Stdout, code, "go", "terminal256", "monokai")
	if err != nil {
		// handle error
	}
}
