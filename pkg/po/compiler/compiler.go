package compiler

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Tom5521/xgotext/pkg/po/types"
)

// Compiler is responsible for compiling translations from a `types.File`
// into different output formats, such as strings, byte slices, or files.
type Compiler struct {
	File   *types.File // The source file containing translation entries.
	Config Config      // Configuration settings for compilation.
}

// applyOptions applies a set of options to modify the compiler's configuration.
func (c *Compiler) applyOptions(opts ...Option) {
	for _, opt := range opts {
		opt(&c.Config)
	}
}

// New creates a new Compiler instance with the given translation file and options.
// The provided options override the default configuration.
func New(file *types.File, options ...Option) Compiler {
	return Compiler{
		File:   file,
		Config: DefaultConfig(options...),
	}
}

// ToWriter writes the compiled translations to an `io.Writer` in the PO file format.
// The provided options override the instance's configuration.
func (c Compiler) ToWriter(w io.Writer, options ...Option) error {
	// Apply the provided options, which take precedence over the instance's configuration.
	c.applyOptions(options...)
	var err error

	if c.Config.Verbose {
		c.Config.Logger.Println("Writing header...")
	}
	// Write the PO file header.
	_, err = fmt.Fprintln(w, c.formatHeader())
	if err != nil && !c.Config.IgnoreErrors {
		err = fmt.Errorf("error writing header format: %w", err)
		c.Config.Logger.Println("ERROR:", err)
		return err
	}
	if c.Config.Verbose {
		c.Config.Logger.Println("Cleaning duplicates...")
	}
	// Remove duplicate entries and write each entry to the writer.
	entries := c.File.Entries.CleanDuplicates()
	if c.Config.Verbose {
		c.Config.Logger.Println("Writing entries...")
	}
	for i, e := range entries {
		if e.ID == "" {
			continue
		}

		_, err = fmt.Fprintln(w, c.formatEntry(e))
		if err != nil && !c.Config.IgnoreErrors {
			err = fmt.Errorf("error writing entry[%d]: %w", i, err)
			c.Config.Logger.Println("ERROR:", err)
			return err
		}
	}
	return nil
}

// ToFile writes the compiled translations to a specified file.
// If `ForcePo` is enabled, the file is created or truncated before writing.
// The provided options override the instance's configuration.
func (c Compiler) ToFile(f string, options ...Option) error {
	flags := os.O_RDWR
	if c.Config.ForcePo {
		flags |= os.O_CREATE
	}

	if c.Config.Verbose {
		c.Config.Logger.Println("Opening file...")
	}
	// Open the file with the determined flags.
	file, err := os.OpenFile(f, flags, os.ModePerm)
	if err != nil && !c.Config.IgnoreErrors {
		err = fmt.Errorf("error opening file: %w", err)
		c.Config.Logger.Println("ERROR:", err)
		return err
	}
	defer file.Close()

	// If `ForcePo` is enabled, truncate and reset the file position.
	if c.Config.ForcePo {
		if c.Config.Verbose {
			c.Config.Logger.Println("Cleaning file contents...")
		}
		err = file.Truncate(0)
		if err != nil && !c.Config.IgnoreErrors {
			err = fmt.Errorf("error truncating file: %w", err)
			c.Config.Logger.Println("ERROR:", err)
			return err
		}

		// Move the file pointer back to the beginning.
		_, err = file.Seek(0, 0)
		if err != nil && !c.Config.IgnoreErrors {
			err = fmt.Errorf("error moving the file pointer back to the beginning: %w", err)
			c.Config.Logger.Println("ERROR:", err)
			return err
		}
	}

	// Write compiled translations to the file.
	return c.ToWriter(file, options...)
}

// ToString compiles the translations and returns the result as a string.
// The provided options override the instance's configuration.
func (c Compiler) ToString(options ...Option) string {
	var b strings.Builder

	// Write the compiled content to the string builder.
	c.ToWriter(&b, options...)

	return b.String()
}

// ToBytes compiles the translations and returns the result as a byte slice.
// The provided options override the instance's configuration.
func (c Compiler) ToBytes(options ...Option) []byte {
	var b bytes.Buffer

	// Write the compiled content to the byte buffer.
	c.ToWriter(&b, options...)

	return b.Bytes()
}
