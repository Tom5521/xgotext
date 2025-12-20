package compile

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/Tom5521/gotext-tools/v2/internal/slices"
	"github.com/Tom5521/gotext-tools/v2/pkg/po"
)

const (
	copyrightFormat = `# Copyright (C) %s
# This file is distributed under the same license as the %s package.`
	foreignCopyrightFormat = `# This file is put in the public domain.`
	headerFormat           = `# %s
%s
#
` + headerEntry
	headerEntry = `msgid ""
msgstr ""
`
	headerFieldFormat = `"%s: %s\n"`
)

func (c PoCompiler) compileEntries(writer io.Writer, eb *entryBuilder, entries po.Entries) error {
	c.info("writing entries...")

	for _, e := range entries {
		eb.Entry = e
		_, err := writer.Write(eb.BuildEntry())
		if err != nil {
			return c.error("error writing entry: %w", err)
		}
	}

	return nil
}

func escapePOString(s string) string {
	var buf strings.Builder
	for _, r := range s {
		switch r {
		case '"':
			buf.WriteString(`\"`)
		case '\\':
			buf.WriteString(`\\`)
		case '\n':
			buf.WriteString(`\n`)
		case '\t':
			buf.WriteString(`\t`)
		case '\r':
			buf.WriteString(`\r`)
		default:
			if strconv.IsPrint(r) {
				buf.WriteRune(r)
			} else {
				fmt.Fprintf(&buf, "\\x%02x", r)
			}
		}
	}
	return buf.String()
}

func (c PoCompiler) writeHeader(writer io.Writer, entries *po.Entries, eb *entryBuilder) error {
	c.info("writing header...")
	i := c.File.Index("", "")
	if i == -1 {
		return nil
	}
	header := c.File.HeaderFromIndex(i)
	*entries = slices.Delete(*entries, i, i+1)

	if c.Config.HeaderConfig != nil {
		header = c.Config.HeaderConfig.ToHeader()
	}

	_, err := writer.Write(eb.BuildHeader(header))
	return err
}
