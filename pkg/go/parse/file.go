package parse

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"

	"github.com/Tom5521/xgotext/pkg/po"
)

// Constants.
const (
	DefaultPackageName = "gotext"
	// WantedImport specifies the import path of the "gotext" library
	// that the parser is looking for in the Go files.
	WantedImport = `"github.com/leonelquinteros/gotext"`
)

// File represents a Go source file that is being processed by the parser.
//
// ### Attributes:
// - `config`: The parser configuration options (e.g., to enable verbose logging or extract all strings).
// - `file`: The parsed AST (abstract syntax tree) representation of the Go source file.
// - `content`: The raw content of the source file as a byte slice.
// - `path`: The file path of the Go source file.
// - `pkgName`: The name of the package declared in the file. Defaults to "gotext" unless overridden.
// - `hasGotext`: Indicates whether the file imports the "gotext" package (required for translation extraction).
// - `seenTokens`: Tracks processed AST nodes to avoid duplicate entries.
//
// ### Responsibilities:
// - Parse the Go source file into an AST.
// - Extract translation entries and their metadata (e.g., locations in the source file).
// - Check if the file imports the "gotext" library and determine the package alias if used.
//
// ### Methods:
// - `NewFileFromReader`: Creates a `File` instance from an `io.Reader`.
// - `NewFileFromPath`: Creates a `File` instance from a file path.
// - `NewFileFromBytes`: Creates a `File` instance from raw byte data.
// - `Entries`: Extracts all translation entries from the file.
// - `parse`: Parses the file content into an AST.
// - `determinePackageInfo`: Extracts package-related metadata, such as the package name and `gotext` import.
type File struct {
	config     Config
	options    []Option
	seenTokens map[ast.Node]bool
	file       *ast.File // The parsed abstract syntax tree (AST) of the file.
	reader     *bytes.Reader
	name       string // The path to the file.
	pkgName    string // The name of the package declared in the file.
	hasGotext  bool   // Indicates if the file imports the desired "gotext" package.
}

func (f *File) Reset(d io.Reader, name string, options ...Option) error {
	f.seenTokens = nil
	f.options = options

	b, err := io.ReadAll(d)
	if err != nil {
		return err
	}
	f.reader = bytes.NewReader(b)
	f.name = name

	if err = f.parse(); err != nil {
		return err
	}

	f.determinePackageInfo()

	return nil
}

func NewFileFromBytes(b []byte, name string, options ...Option) (*File, error) {
	file := &File{
		reader:  bytes.NewReader(b),
		name:    name,
		options: options,
		pkgName: DefaultPackageName,
	}

	if err := file.parse(); err != nil {
		return nil, err
	}

	file.determinePackageInfo()
	return file, nil
}

// NewFileFromPath creates a new File instance by reading content from a file on disk.
func NewFileFromPath(path string, options ...Option) (*File, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	return NewFile(file, path, options...)
}

// NewFile creates a new File instance from raw byte data.
func NewFile(b io.Reader, name string, options ...Option) (*File, error) {
	bytedata, err := io.ReadAll(b)
	if err != nil {
		return nil, err
	}

	return NewFileFromBytes(bytedata, name, options...)
}

// parse parses the file content into an AST.
func (f *File) parse() error {
	var err error
	f.file, err = parser.ParseFile(token.NewFileSet(), f.name, f.reader, 0)
	if err != nil {
		return fmt.Errorf("failed to parse the file: %w", err)
	}
	return nil
}

// determinePackageInfo analyzes the file's AST to extract package-related information.
// It determines the package name and checks if the desired "gotext" package is imported.
func (f *File) determinePackageInfo() {
	for _, imp := range f.file.Imports {
		if imp.Path.Value == WantedImport {
			f.hasGotext = true
			if imp.Name != nil {
				f.pkgName = imp.Name.String()
			}
			break
		}
	}
}

// Entries returns all translations found in the file.
func (f *File) Entries() (po.Entries, []error) {
	f.seenTokens = make(map[ast.Node]bool)
	for _, opt := range f.options {
		opt(&f.config)
	}

	defer func() {
		f.seenTokens = nil
	}()

	var entries po.Entries
	var errors []error

	if !f.hasGotext && !f.config.ExtractAll {
		return entries, errors
	}

	ast.Inspect(f.file, func(n ast.Node) bool {
		t, e := f.processNode(n)
		entries = append(entries, t...)
		errors = append(errors, e...)
		return true
	})

	return entries, errors
}
