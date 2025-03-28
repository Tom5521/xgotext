package po_test

import (
	"math/rand"
	"testing"

	"github.com/Tom5521/xgotext/pkg/po"
)

func BenchmarkSortEntries(b *testing.B) {
	entries := po.Entries{
		{
			ID:      "Apple",
			Context: "USA",
			Plural:  "Apples",
			Plurals: po.PluralEntries{
				{0, "Manzana"},
				{1, "Manzanas"},
			},
		},
		{ID: "Hi", Str: "Hola", Context: "casual"},
		{ID: "", Str: ""},
		{ID: "How are you?", Str: "Como estás?"},
	}

	rand.Shuffle(len(entries), func(i, j int) {
		entries[i], entries[j] = entries[j], entries[i]
	})

	tests := []struct {
		name   string
		method func() po.Entries
	}{
		{"Sort", entries.Sort},
		{"SortByFile", entries.SortByFile},
		{"SortByID", entries.SortByID},
		{"SortByLine", entries.SortByLine},
		{"SortByFuzzy", entries.SortByFuzzy},
	}

	for _, t := range tests {
		b.Run(t.name, func(b *testing.B) {
			t.method()
		})
	}
}

// TODO: Finish this.
func BenchmarkEntriesSolve(b *testing.B) {
	b.SkipNow() // This isn't finished yet.
	entries := po.Entries{
		{
			ID:      "Apple",
			Context: "USA",
			Plural:  "Apples",
			Plurals: po.PluralEntries{
				{0, "Manzana"},
				{1, "Manzanas"},
			},
		},
		{ID: "Hi", Str: "Hola", Context: "casual"},
		{ID: "", Str: ""},
		{ID: "How are you?", Str: "Como estás?"},
	}

	rand.Shuffle(len(entries), func(i, j int) {
		entries[i], entries[j] = entries[j], entries[i]
	})

	tests := []struct {
		name   string
		method func() po.Entries
	}{
		// {"Solve", entries.Solve},
	}

	for _, t := range tests {
		b.Run(t.name, func(b *testing.B) {
			t.method()
		})
	}
}
