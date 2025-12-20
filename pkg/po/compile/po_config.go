package compile

import (
	"errors"
	"io"
	"log"

	"github.com/Tom5521/gotext-tools/v2/pkg/po"
)

// NOTE: The PoConfig structure should focus less on meeting the needs of gettext tools implementation and focus more on fulfilling the requirements as a library.
//
// This should be refactored.
//
// By the way, the options like CleanDuplicates must be removed, because
// it's just unnesessary.

// PoConfig holds the settings for the compiler, affecting how translations are processed.
type PoConfig struct {
	// It is used to restore the configuration using the method [PoConfig.RestoreLastCfg]
	// and is saved when using the asd method [PoConfig.ApplyOptions]
	lastCfg any

	// The logger can be nil, otherwise this logger will be used to print all errors by default.
	Logger *log.Logger

	// NOTE: This setting behavior is WRONG; it should be changed.

	// But, currently, it works as it says fortunately.
	ForcePo         bool           // If true, forces the creation of a `.po` file, even if not strictly needed.
	OmitHeader      bool           // If true, omits the header section in the generated `.po` file.
	PackageName     string         // Name of the package associated with the translation.
	CopyrightHolder string         // Name of the entity holding copyright over the translation.
	ForeignUser     bool           // If true, marks the translation as public domain.
	Title           string         // Title to be included in the `.po` file header.
	NoLocation      bool           // If true, suppresses location comments in the `.po` file.
	AddLocation     PoLocationMode // Specifies how location comments should be included ("never", "file", "full").
	MsgstrPrefix    string         // Prefix added to all translation strings.
	MsgstrSuffix    string         // Suffix added to all translation strings.
	IgnoreErrors    bool           // If true, allows compilation to proceed despite non-critical errors.
	Verbose         bool
	CommentFuzzy    bool
	ManageHeader    bool
	HeaderComments  bool
	HeaderFields    bool
	WordWrap        bool
	HeaderConfig    *po.HeaderConfig

	UseCustomObsoletePrefix  bool
	CustomObsoletePrefixRune rune
	NoColor                  bool
}

func NewPoConfigFromOptions(opts ...PoOption) PoConfig {
	var config PoConfig

	config.ApplyOptions(opts...)

	return config
}

// ApplyOptions overwrites the configuration with the options provided,
// saving the previous state so that it can be restored
// later with [PoConfig.RestoreLastCfg] if desired.
func (c *PoConfig) ApplyOptions(opts ...PoOption) {
	c.lastCfg = *c
	for _, po := range opts {
		po(c)
	}
}

// RestoreLastCfg restores the configuration state prior to the last
// [PoConfig.ApplyOptions] if it exists, otherwise it does nothing.
func (c *PoConfig) RestoreLastCfg() {
	*c = c.lastCfg.(PoConfig)
}

type PoLocationMode string

const (
	PoLocationModeFull  PoLocationMode = "full"
	PoLocationModeNever PoLocationMode = "never"
	PoLocationModeFile  PoLocationMode = "file"
)

func DefaultPoConfig(opts ...PoOption) PoConfig {
	c := PoConfig{
		Logger:       log.New(io.Discard, "", 0),
		PackageName:  "PACKAGE NAME",
		AddLocation:  PoLocationModeFull,
		ManageHeader: false,
	}

	c.ApplyOptions(opts...)

	return c
}

type PoOption func(*PoConfig)

func PoWithNoColor(v bool) PoOption {
	return func(pc *PoConfig) {
		pc.NoColor = v
	}
}

func PoWithUseCustomObsoletePrefix(u bool) PoOption {
	return func(pc *PoConfig) {
		pc.UseCustomObsoletePrefix = u
	}
}

func PoWithCustomObsoletePrefixRune(r rune) PoOption {
	return func(pc *PoConfig) {
		pc.CustomObsoletePrefixRune = r
	}
}

func PoWithManageHeader(b bool) PoOption {
	return func(pc *PoConfig) {
		pc.ManageHeader = b
	}
}

func PoWithWordWrap(w bool) PoOption {
	return func(pc *PoConfig) {
		pc.WordWrap = w
	}
}

func PoWithHeaderFields(w bool) PoOption {
	return func(pc *PoConfig) {
		pc.HeaderFields = w
	}
}

func PoWithHeaderComments(hc bool) PoOption {
	return func(pc *PoConfig) {
		pc.HeaderComments = hc
	}
}

func PoWithCommentFuzzy(c bool) PoOption {
	return func(pc *PoConfig) {
		pc.CommentFuzzy = c
	}
}

func PoWithVerbose(v bool) PoOption {
	return func(c *PoConfig) {
		c.Verbose = v
	}
}

func PoWithLogger(logger *log.Logger) PoOption {
	return func(c *PoConfig) {
		c.Logger = logger
	}
}

func PoWithIgnoreErrors(i bool) PoOption {
	return func(c *PoConfig) {
		c.IgnoreErrors = i
	}
}

func PoWithConfig(cfg PoConfig) PoOption {
	return func(c *PoConfig) {
		*c = cfg
	}
}

func PoWithForcePo(f bool) PoOption {
	return func(c *PoConfig) {
		c.ForcePo = f
	}
}

func PoWithOmitHeader(o bool) PoOption {
	return func(c *PoConfig) {
		c.OmitHeader = o
	}
}

func PoWithPackageName(name string) PoOption {
	return func(c *PoConfig) {
		c.PackageName = name
	}
}

func PoWithCopyrightHolder(holder string) PoOption {
	return func(c *PoConfig) {
		c.CopyrightHolder = holder
	}
}

func PoWithForeignUser(f bool) PoOption {
	return func(c *PoConfig) {
		c.ForeignUser = f
	}
}

func PoWithTitle(t string) PoOption {
	return func(c *PoConfig) {
		c.Title = t
	}
}

func PoWithNoLocation(n bool) PoOption {
	return func(c *PoConfig) {
		c.NoLocation = n
	}
}

func PoWithAddLocation(loc PoLocationMode) PoOption {
	return func(c *PoConfig) {
		c.AddLocation = loc
	}
}

func PoWithMsgstrPrefix(prefix string) PoOption {
	return func(c *PoConfig) {
		c.MsgstrPrefix = prefix
	}
}

func PoWithMsgstrSuffix(suffix string) PoOption {
	return func(c *PoConfig) {
		c.MsgstrSuffix = suffix
	}
}

func (c PoConfig) Validate() error {
	if c.NoLocation && c.AddLocation != PoLocationModeNever {
		return errors.New("noLocation and AddLocation are in conflict")
	}

	return nil
}
