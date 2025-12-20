package compile

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/Tom5521/gotext-tools/v2/pkg/po"
)

type entryBuilder struct {
	po.Entry
	Config PoConfig
}

func (eb *entryBuilder) BuildEntry() []byte {
	var buf bytes.Buffer
	buf.WriteString(eb.comment())
	buf.WriteString(eb.msgid())
	buf.WriteString(eb.msgstr())

	commentFuzzy := eb.IsFuzzy() && eb.Config.CommentFuzzy

	if eb.Obsolete || commentFuzzy {
		entry := buf.String()
		buf.Reset()

		prefix := "#"
		if eb.Obsolete {
			if eb.Config.UseCustomObsoletePrefix {
				prefix += string(eb.Config.CustomObsoletePrefixRune)
			} else {
				prefix += "~"
			}
		}
		prefix += " "

		for _, line := range strings.Split(entry, "\n") {
			if strings.HasPrefix(line, "#") || line == "" {
				buf.WriteString(line + "\n")
				continue
			}

			buf.WriteString(prefix + line + "\n")
		}
	}

	buf.WriteByte('\n')
	return buf.Bytes()
}

func (eb *entryBuilder) BuildHeader(header po.Header) []byte {
	var buf bytes.Buffer

	if eb.Config.HeaderComments {
		copyright := fmt.Sprintf(copyrightFormat, eb.Config.CopyrightHolder, eb.Config.PackageName)
		if eb.Config.ForeignUser {
			copyright = foreignCopyrightFormat
		}

		fmt.Fprintf(&buf, headerFormat, eb.Config.Title, copyright)
	} else {
		fmt.Fprint(&buf, headerEntry)
	}

	if eb.Config.HeaderFields {
		for i, field := range header.Fields {
			fmt.Fprintf(&buf, headerFieldFormat, field.Key, field.Value)

			if i != len(header.Fields) {
				fmt.Fprint(&buf, "\n")
			}
		}
	}

	fmt.Fprintln(&buf)

	buf.WriteString(eb.msgid())
	buf.WriteString(eb.msgstr())

	return buf.Bytes()
}

func (eb *entryBuilder) msgid() string {
	var b strings.Builder
	if eb.HasContext() {
		b.WriteString(eb.keyword("msgctxt"))
		b.WriteString(eb.string(eb.Context))
	}
	b.WriteString(eb.keyword("msgid"))
	b.WriteString(eb.string(eb.ID))

	if eb.IsPlural() {
		b.WriteString(eb.keyword("msgid_plural"))
		b.WriteString(eb.string(eb.Plural))
	}

	return b.String()
}

func (eb *entryBuilder) msgstr() string {
	var msgstr strings.Builder
	const format = "msgstr[%d]"
	if eb.IsPlural() {
		if len(eb.Plurals) == 0 {
			id := eb.string(eb.ID)
			for i := 0; i < 2; i++ {
				fmt.Fprint(&msgstr, eb.keyword(fmt.Sprintf(format, i)))
				fmt.Fprint(&msgstr, id)
			}
			return msgstr.String()
		}
		for _, pe := range eb.Plurals {
			fmt.Fprint(&msgstr, eb.keyword(fmt.Sprintf(format, pe.ID)))
			fmt.Fprint(&msgstr,
				eb.string(
					eb.Config.MsgstrPrefix+pe.Str+eb.Config.MsgstrSuffix,
				),
			)
		}

		return msgstr.String()
	}

	fmt.Fprint(&msgstr, eb.keyword("msgstr"))
	fmt.Fprint(&msgstr, eb.string(
		eb.Config.MsgstrPrefix+eb.Str+eb.Config.MsgstrSuffix,
	))

	return msgstr.String()
}

func (eb *entryBuilder) comment() string {
	var b strings.Builder
	b.WriteString(eb.translatorComment())
	b.WriteString(eb.extractedComment())
	b.WriteString(eb.referenceComment())
	b.WriteString(eb.flagComment())
	b.WriteString(eb.previousComment())

	return b.String()
}

func (eb *entryBuilder) translatorComment() string {
	var b strings.Builder
	for _, comment := range eb.Comments {
		fmt.Fprintf(&b, "# %s\n", comment)
	}
	return b.String()
}

func (eb *entryBuilder) extractedComment() string {
	var b strings.Builder
	for _, comment := range eb.ExtractedComments {
		fmt.Fprintf(&b, "#. %s\n", comment)
	}
	return b.String()
}

func (eb *entryBuilder) referenceComment() string {
	if eb.Config.NoLocation || eb.Config.AddLocation == PoLocationModeNever {
		return ""
	}
	var b strings.Builder

	var writeRef func(id int)
	switch eb.Config.AddLocation {
	case PoLocationModeFull:
		writeRef = func(id int) {
			l := eb.Locations[id]
			fmt.Fprintf(&b, "%s:%d\n", l.File, l.Line)
		}
	case PoLocationModeFile:
		writeRef = func(id int) {
			l := eb.Locations[id]
			fmt.Fprintf(&b, "%s\n", l.File)
		}
	}

	for i := range eb.Locations {
		fmt.Fprint(&b, "#: ")
		writeRef(i)
	}

	return b.String()
}

func (eb *entryBuilder) flagComment() string {
	var comments string
	for _, f := range eb.Flags {
		comment := fmt.Sprintf("#, %s\n", f)
		comments += comment
	}

	return comments
}

func (eb *entryBuilder) previousComment() string {
	var b strings.Builder
	for _, p := range eb.Previous {
		fmt.Fprintf(&b, "#| %s\n", p)
	}
	return b.String()
}

func (eb *entryBuilder) string(str string) string {
	var builder strings.Builder
	if eb.Config.WordWrap {
		lines := strings.Split(str, "\n")
		for i, line := range lines {
			if i != len(lines)-1 {
				line += "\n"
			}
			fmt.Fprintf(&builder, "\"%s\"\n", escapePOString(line))
		}
		return builder.String()
	}

	fmt.Fprintf(&builder, "\"%s\"\n", escapePOString(str))

	return builder.String()
}

func (eb *entryBuilder) text(str string) string {
	return str
}

func (eb *entryBuilder) keyword(kw string) string {
	return kw + " "
}
