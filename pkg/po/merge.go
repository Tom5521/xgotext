package po

import (
	"slices"

	"github.com/Tom5521/xgotext/internal/util"
)

type SortMode int

const (
	SortByAll SortMode = iota
	SortByID
	SortByFile
	SortByLine
	SortByFuzzy
	SortByObsolete
	NoSort
)

func (mode SortMode) SortMethod(entries Entries) func() Entries {
	method, ok := map[SortMode]func() Entries{
		SortByAll:      entries.Sort,
		SortByID:       entries.SortByID,
		SortByFile:     entries.SortByFile,
		SortByLine:     entries.SortByLine,
		SortByFuzzy:    entries.SortByFuzzy,
		SortByObsolete: entries.SortByObsolete,
		NoSort:         func() Entries { return entries },
	}[mode]

	if !ok {
		return entries.Sort
	}

	return method
}

type MergeConfig struct {
	FuzzyMatch      bool
	KeepPreviousIDs bool
	Sort            bool
	SortMode        SortMode
}

func NewMergeConfig(opts ...MergeOption) MergeConfig {
	var m MergeConfig
	m.ApplyOption(opts...)
	return m
}

func DefaultMergeConfig() MergeConfig {
	return MergeConfig{
		FuzzyMatch: true,
		Sort:       true,
		SortMode:   SortByAll,
	}
}

func (m *MergeConfig) ApplyOption(opts ...MergeOption) {
	for _, mo := range opts {
		mo(m)
	}
}

type MergeOption func(mc *MergeConfig)

func MergeWithMergeConfig(n MergeConfig) MergeOption {
	return func(mc *MergeConfig) {
		*mc = n
	}
}

func MergeWithFuzzyMatch(f bool) MergeOption {
	return func(mc *MergeConfig) {
		mc.FuzzyMatch = f
	}
}

func MergeWithKeepPreviousIDs(k bool) MergeOption {
	return func(mc *MergeConfig) {
		mc.KeepPreviousIDs = k
	}
}

func (f File) MergeWithConfig(config MergeConfig, files ...Entries) *File {
	ref := f // this is only an alias for "f"
	def := File{slices.Clone(f.Entries), f.Name}

	for _, f2 := range files {
		def.Entries = append(def.Entries, f2...)
	}

	// TODO: Finish this.
	// def.Entries = def.Entries.Solve()

	for i, e := range def.Entries {
		for _, e2 := range ref.Entries {
			if e.Context != e2.Context {
				continue
			}
			if util.IsSimilarButNotIdentical(e.ID, e2.ID) {
				e.Flags = append(e.Flags, "fuzzy")
			}
		}
		def.Entries[i] = e
	}

	for i, e := range def.Entries {
		if e.IsFuzzy() {
			continue
		}
		if !ref.Entries.ContainsUnifiedID(e.UnifiedID()) {
			e.Obsolete = true
			def.Entries[i] = e
		}
	}

	return &def
}

func (f File) MergeWithOptions(files []Entries, options ...MergeOption) *File {
	var cfg MergeConfig
	cfg.ApplyOption(options...)
	return f.MergeWithConfig(cfg, files...)
}

func (f File) Merge(files ...Entries) *File {
	return f.MergeWithConfig(DefaultMergeConfig(), files...)
}
