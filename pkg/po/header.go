package po

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/Tom5521/xgotext/internal/util"
)

// HeaderField represents a single key-value pair in a header.
type HeaderField struct {
	Key   string // The name of the header field.
	Value string // The value associated with the header field.
}

// Header represents a collection of header fields.
type Header struct {
	Fields []HeaderField // A slice storing all registered header fields.
}

func NewHeader(options ...HeaderOption) Header {
	h := HeaderConfig{}
	for _, opt := range options {
		opt(&h)
	}

	return h.ToHeader()
}

func (h Header) Nplurals() (nplurals uint) {
	nplurals = 2
	value := h.Load("Plural-Forms")
	if value == "" {
		return
	}
	if !npluralsRegex.MatchString(value) {
		return
	}
	matches := npluralsRegex.FindStringSubmatch(value)
	n, err := strconv.ParseUint(util.SafeSliceAccess(matches, 1), 10, 64)
	if err != nil {
		return
	}

	nplurals = uint(n)

	return
}

func (h Header) ToEntry() (e Entry) {
	var b strings.Builder

	for _, field := range h.Fields {
		fmt.Fprintf(&b, "\n%s: %s\n", field.Key, field.Value)
	}

	e.Str = b.String()
	return
}

type HeaderConfig struct {
	Nplurals          uint
	ProjectIDVersion  string
	ReportMsgidBugsTo string
	Language          string
	LanguageTeam      string
	LastTranslator    string
}

func (cfg HeaderConfig) ToHeaderWithDefaults() (h Header) {
	h = DefaultTemplateHeader()

	h.Set("Project-Id-Version", cfg.ProjectIDVersion)
	h.Set("Report-Msgid-Bugs-To", cfg.ReportMsgidBugsTo)
	h.Set("Language-Team", cfg.LanguageTeam)
	h.Set("Language", cfg.Language)
	h.Set(
		"Plural-Forms",
		fmt.Sprintf("nplurals=%d; plural=(n != 1);", cfg.Nplurals),
	)

	return
}

func (cfg HeaderConfig) ToHeader() (h Header) {
	h.Set("Project-Id-Version", cfg.ProjectIDVersion)
	h.Set("Report-Msgid-Bugs-To", cfg.ReportMsgidBugsTo)
	h.Set("Language-Team", cfg.LanguageTeam)
	h.Set("Language", cfg.Language)
	h.Set(
		"Plural-Forms",
		fmt.Sprintf("nplurals=%d; plural=(n != 1);", cfg.Nplurals),
	)

	return
}

func HeaderConfigFromOptions(options ...HeaderOption) HeaderConfig {
	var h HeaderConfig
	for _, opt := range options {
		opt(&h)
	}

	return h
}

func DefaultHeaderConfig() HeaderConfig {
	return HeaderConfig{
		Nplurals:         2,
		ProjectIDVersion: "PACKAGE VERSION",
		Language:         "en",
	}
}

var (
	npluralsRegex = regexp.MustCompile(`nplurals=(\d*)`)
	headerRegex   = regexp.MustCompile(`(.*)\s*:\s*(.*)`)
)

func (e Entries) Header() (h Header) {
	i := e.IndexByIDAndCtx("", "")
	if i == -1 {
		return
	}

	header := e[i].Str
	lines := strings.Split(header, "\n")
	for _, line := range lines {
		if !headerRegex.MatchString(line) {
			continue
		}
		matches := headerRegex.FindStringSubmatch(line)
		h.Fields = append(h.Fields,
			HeaderField{
				Key:   util.SafeSliceAccess(matches, 1),
				Value: util.SafeSliceAccess(matches, 2),
			},
		)
	}
	return
}

// DefaultTemplateHeader initializes a Header object with commonly used default fields.
// These fields are typically found in .po files for localization.
func DefaultTemplateHeader() (h Header) {
	// Register standard header fields with optional default values.
	h.Register("Project-Id-Version")                                  // No default value.
	h.Register("Report-Msgid-Bugs-To")                                // No default value.
	h.Register("POT-Creation-Date", time.Now().Format(time.DateTime)) // Current date and time.
	h.Register("PO-Revision-Date")                                    // No default value.
	h.Register("Last-Translator")                                     // No default value.
	h.Register("Language-Team")                                       // No default value.
	h.Register("Language")                                            // No default value.
	h.Register("MIME-Version", "1.0")                                 // MIME version.
	h.Register(
		"Content-Type",
		"text/plain; charset=CHARSET",
	) // Content type with placeholder charset.
	h.Register("Content-Transfer-Encoding", "8bit") // Encoding type.
	h.Register(
		"Plural-Forms",
		"nplurals=2; plural=(n != 1);",
	) // Placeholder plural form formula.
	return h
}

// Register adds a new header field to the Header object if the key does not already exist.
// Parameters:
//   - key: The name of the header field to register.
//   - d: Optional variadic arguments representing the value(s) to associate with the key.
//     If provided, they are concatenated into a single string using fmt.Sprint.
func (h *Header) Register(key string, d ...string) {
	// Check if the key already exists in the Values slice.
	i := slices.IndexFunc(h.Fields, func(f HeaderField) bool {
		return f.Key == key
	})
	if i != -1 {
		// Key already exists; do nothing.
		return
	}

	var values []any
	for _, b := range d {
		values = append(values, b)
	}

	// Append a new HeaderField to the Values slice.
	h.Fields = append(h.Fields,
		HeaderField{
			Key:   key,
			Value: fmt.Sprint(values...), // Concatenate variadic arguments into a single string.
		},
	)
}

// Load retrieves the value associated with a given key from the Header object.
// Parameters:
// - key: The name of the header field to retrieve.
// Returns:
// - The value associated with the key if found; otherwise, an empty string ("").
func (h *Header) Load(key string) string {
	// Search for the key in the Values slice.
	i := slices.IndexFunc(h.Fields, func(f HeaderField) bool {
		return f.Key == key
	})
	if i >= 0 {
		// Key found; return its value.
		return h.Fields[i].Value
	}
	// Key not found; return an empty string.
	return ""
}

// Set updates the value of an existing header field or adds a new field if the key does not exist.
// Parameters:
// - key: The name of the header field to update or add.
// - value: The new value to associate with the key.
func (h *Header) Set(key, value string) {
	// Search for the key in the Values slice.
	i := slices.IndexFunc(h.Fields, func(f HeaderField) bool {
		return f.Key == key
	})
	if i >= 0 {
		// Key found; update its value.
		h.Fields[i].Value = value
		return
	}
	// Key not found; append a new HeaderField to the Values slice.
	h.Fields = append(h.Fields,
		HeaderField{
			Key:   key,
			Value: value,
		},
	)
}
