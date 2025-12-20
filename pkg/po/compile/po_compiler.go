package compile

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Tom5521/gotext-tools/v2/internal/slices"
	"github.com/Tom5521/gotext-tools/v2/pkg/po"
)

var _ po.Compiler = (*PoCompiler)(nil)

type PoCompiler struct {
	// File is the source PO file containing translation entries
	File *po.File
	// Config contains compilation settings and options
	Config PoConfig
}

// NewPo creates a new PoCompiler instance with the given PO file and options.
// The options are applied to configure the compiler behavior.
func NewPo(file *po.File, options ...PoOption) PoCompiler {
	return PoCompiler{
		File:   file,
		Config: DefaultPoConfig(options...),
	}
}

// error creates and logs an error message if error reporting is enabled.
// Returns nil if IgnoreErrors is true.
func (c PoCompiler) error(format string, a ...any) error {
	if c.Config.IgnoreErrors {
		return nil
	}

	format = "compile: " + format
	err := fmt.Errorf(format, a...)
	if c.Config.Logger != nil {
		c.Config.Logger.Println("ERROR:", err)
	}

	return err
}

// info logs an informational message if verbose logging is enabled.
func (c PoCompiler) info(format string, a ...any) {
	if c.Config.Logger != nil && c.Config.Verbose {
		c.Config.Logger.Println("INFO:", fmt.Sprintf(format, a...))
	}
}

// SetFile updates the PO file reference in the compiler.
func (c *PoCompiler) SetFile(f *po.File) {
	c.File = f
}

// ToWriterWithOptions writes compiled output to an io.Writer with temporary options.
// The options are only applied for this operation and then reverted.
func (c *PoCompiler) ToWriterWithOptions(w io.Writer, opts ...PoOption) error {
	c.Config.ApplyOptions(opts...)
	defer c.Config.RestoreLastCfg()
	return c.ToWriter(w)
}

// ToStringWithOptions returns the compiled output as a string with temporary options.
// The options are only applied for this operation and then reverted.
func (c *PoCompiler) ToStringWithOptions(opts ...PoOption) string {
	c.Config.ApplyOptions(opts...)
	defer c.Config.RestoreLastCfg()
	return c.ToString()
}

// ToFileWithOptions writes compiled output to a file with temporary options.
// The options are only applied for this operation and then reverted.
func (c *PoCompiler) ToFileWithOptions(f string, opts ...PoOption) error {
	c.Config.ApplyOptions(opts...)
	defer c.Config.RestoreLastCfg()
	return c.ToFile(f)
}

// ToBytesWithOptions returns the compiled output as bytes with temporary options.
// The options are only applied for this operation and then reverted.
func (c *PoCompiler) ToBytesWithOptions(opts ...PoOption) []byte {
	c.Config.ApplyOptions(opts...)
	defer c.Config.RestoreLastCfg()
	return c.ToBytes()
}

// ToWriter writes compiled PO content to an io.Writer.
// Handles header writing, duplicate cleaning, and optional syntax highlighting.
func (c PoCompiler) ToWriter(outputWriter io.Writer) error {
	entries := c.File.Entries

	buffer := bufio.NewWriter(outputWriter)

	var writer io.Writer = buffer

	if c.Config.OmitHeader {
		i := c.File.Index("", "")
		if i != -1 {
			entries = slices.Delete(entries, i, i+1)
		}
	}

	eb := entryBuilder{
		Config: c.Config,
	}

	if c.Config.ManageHeader && !c.Config.OmitHeader {
		err := c.writeHeader(writer, &entries, &eb)
		if err != nil {
			return c.error("error writing header: %w", err)
		}
	}

	err := c.compileEntries(writer, &eb, entries)
	if err != nil {
		return c.error("error compiling entries: %w", err)
	}

	err = buffer.Flush()
	if err != nil {
		return c.error("error flushing buffer: %w", err)
	}

	return nil
}

// ToFile writes compiled output to the specified file path.
// By default fails if file exists (unless ForcePo is enabled).
func (c PoCompiler) ToFile(f string) error {
	flags := os.O_WRONLY | os.O_TRUNC | os.O_CREATE
	if !c.Config.ForcePo {
		flags |= os.O_EXCL
	}

	c.info("opening file...")
	file, err := os.OpenFile(f, flags, 0o600)
	if err != nil {
		return c.error("error opening file: %w", err)
	}
	defer file.Close()

	return c.ToWriter(file)
}

// ToString returns the compiled output as a string.
func (c PoCompiler) ToString() string {
	var b strings.Builder
	c.ToWriter(&b)
	return b.String()
}

// ToBytes returns the compiled output as a byte slice.
func (c PoCompiler) ToBytes() []byte {
	var b bytes.Buffer
	c.ToWriter(&b)
	return b.Bytes()
}
