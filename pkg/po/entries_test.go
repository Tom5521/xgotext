package po_test

import (
	"testing"

	"github.com/Tom5521/xgotext/internal/util"
	"github.com/Tom5521/xgotext/pkg/po"
)

func TestEntries_HasDuplicates(t *testing.T) {
	tests := []struct {
		name string
		e    po.Entries
		want bool
	}{
		{
			"Have",
			po.Entries{
				{ID: "Hello"},
				{ID: "Hello"},
			},
			true,
		},
		{
			"NotHave",
			po.Entries{
				{ID: "Hi"},
				{ID: "Hello"},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.HasDuplicates(); got != tt.want {
				t.Errorf("Entries.HasDuplicates() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntries_CleanDuplicates(t *testing.T) {
	tests := []struct {
		name string
		e    po.Entries
		want po.Entries
	}{
		{
			"NoDuplicates",
			po.Entries{{ID: "Hello1"}, {ID: "Hello2"}},
			po.Entries{{ID: "Hello1"}, {ID: "Hello2"}},
		},
		{
			"HasDuplicates",
			po.Entries{{ID: "Hello"}, {ID: "Hello"}, {ID: "Hello2"}},
			po.Entries{{ID: "Hello"}, {ID: "Hello2"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.CleanDuplicates(); !util.Equal(got, tt.want) {
				t.Errorf("Entries.CleanDuplicates() = %v, want %v", got, tt.want)
			}
		})
	}
}
