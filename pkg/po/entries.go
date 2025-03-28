package po

import (
	"slices"

	"github.com/Tom5521/xgotext/internal/util"
)

// Entries represents a collection of Entry objects.
type Entries []Entry

func (e Entries) Equal(e2 Entries) bool {
	return util.Equal(e, e2)
}

func (e Entries) Contains(c Entry) bool {
	return slices.ContainsFunc(e, func(e Entry) bool { return util.Equal(e, c) })
}

func (e Entries) ContainsUnifiedID(id string) bool {
	return slices.ContainsFunc(e, func(e Entry) bool { return e.UnifiedID() == id })
}

func (e Entries) IndexByUnifiedID(uid string) int {
	return slices.IndexFunc(e, func(e Entry) bool {
		return e.UnifiedID() == uid
	})
}

func (e Entries) Index(e2 Entry) int {
	return slices.IndexFunc(e, func(e Entry) bool { return util.Equal(e, e2) })
}

func (e Entries) IndexByIDAndCtx(id, context string) int {
	return slices.IndexFunc(e,
		func(e Entry) bool {
			return e.ID == id && e.Context == context
		},
	)
}

func (e Entries) IsSorted() bool {
	return slices.IsSortedFunc(e, CompareEntry)
}

// Sort organizes the entries by grouping them by file and sorting them by line.
func (e Entries) Sort() Entries {
	slices.SortFunc(e, CompareEntry)
	return e
}

func (e Entries) IsSortedByObsolete() bool {
	return slices.IsSortedFunc(e, CompareEntryByObsolete)
}

func (e Entries) SortByObsolete() Entries {
	slices.SortFunc(e, CompareEntryByObsolete)
	return e
}

func (e Entries) IsSortedByFuzzy() bool {
	return slices.IsSortedFunc(e, CompareEntryByFuzzy)
}

func (e Entries) SortByFuzzy() Entries {
	slices.SortFunc(e, CompareEntryByFuzzy)
	return e
}

func (e Entries) IsSortedByFile() bool {
	return slices.IsSortedFunc(e, CompareEntryByFile)
}

// SortByFile sorts the entries by the file name of the first location.
func (e Entries) SortByFile() Entries {
	slices.SortFunc(e, CompareEntryByFile)
	return e
}

func (e Entries) IsSortedByID() bool {
	return slices.IsSortedFunc(e, CompareEntryByID)
}

// SortByID sorts the entries by their ID.
func (e Entries) SortByID() Entries {
	slices.SortFunc(e, CompareEntryByID)
	return e
}

func (e Entries) IsSortedByLine() bool {
	return slices.IsSortedFunc(e, CompareEntryByLine)
}

// SortByLine sorts the entries by line number in their first location.
func (e Entries) SortByLine() Entries {
	slices.SortFunc(e, CompareEntryByLine)
	return e
}

func (e Entries) HasDuplicates() bool {
	seen := make(map[string]bool)

	for _, entry := range e {
		uid := entry.UnifiedID()
		_, seened := seen[uid]
		if seened {
			return true
		}

		seen[uid] = true
	}

	return false
}

func (e Entries) CleanObsoletes() Entries {
	return slices.DeleteFunc(e, func(e Entry) bool {
		return e.Obsolete
	})
}

// CleanDuplicates removes duplicate entries with the same ID and context, merging their locations.
func (e Entries) CleanDuplicates() Entries {
	var cleaned Entries
	seenID := make(map[string]int)

	for _, entry := range e {
		uid := entry.UnifiedID()
		idIndex, ok := seenID[uid]
		if ok {
			cleaned[idIndex].Locations = append(cleaned[idIndex].Locations, entry.Locations...)
			continue
		}
		seenID[uid] = len(cleaned)
		cleaned = append(cleaned, entry)
	}

	return cleaned
}

// Returns the result of the merging of e1+e2, if no merging is desired, nil is returned.
type SolveFunc func(e1, e2 Entry) (merged *Entry)

func (e Entries) SolveFunc(f SolveFunc) Entries {
	var cleaned Entries
	seen := make(map[string]int)

	for _, entry := range e {
		uid := entry.UnifiedID()
		idIndex, ok := seen[uid]
		if ok {
			e2 := cleaned[idIndex]
			merged := f(entry, e2)
			if merged != nil {
				cleaned[idIndex] = *merged
			}

			continue
		}
		seen[uid] = len(cleaned)
		cleaned = append(cleaned, entry)
	}

	return cleaned
}

/*
// TODO: Finish this.

// Solve processes a list of translation entries and merges those with the same ID and context,
// keeping the most complete translation. If two entries have the same ID and context, the one
// with a non-empty translation string is retained. Additionally, if the entries are similar but not.

	func (e Entries) Solve() Entries {
		return e.SolveWithUIDPriority(nil)
	}

	func baseSolveFunc(f func(e1, e2 Entry) bool) SolveFunc {
		return func(e1, e2 Entry) (merged *Entry) {
		}
	}

func (e Entries) Solve() Entries {
}

	func (e Entries) SolveWithPriority(preferred Entries) Entries {
		return e.SolveFunc(func(e1, e2 Entry) (merged *Entry) {
			if e1.IsPlural() {
			} else {
			}
		})
	}

	func (e Entries) SolveWithUIDPriority(preferredUIDs []string) Entries {
		return e.SolveFunc(baseSolveFunc(func(e1, e2 Entry) bool {
			e1IsPreferred := slices.Contains(preferredUIDs, e1.UnifiedID())
			e2IsPreferred := slices.Contains(preferredUIDs, e2.UnifiedID())

			if e1IsPreferred == e2IsPreferred {
				return false
			}
		}))
	}
*/
func (e Entries) CleanFuzzy() Entries {
	e = slices.DeleteFunc(e, func(e Entry) bool {
		return e.IsFuzzy()
	})
	return e
}

func (e Entries) FuzzyFind(id, context string) int {
	return slices.IndexFunc(e, func(e Entry) bool {
		return util.FuzzyEqual(id, e.ID) && e.Context == context
	})
}
