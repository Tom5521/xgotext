package po

import (
	"errors"
	"fmt"

	"github.com/Tom5521/xgotext/internal/util"
)

type File struct {
	Entries Entries
	Name    string
}

func NewFile(name string, entries ...Entry) *File {
	f := &File{Name: name, Entries: entries}

	return f
}

func (f File) Validate() error {
	if f.Entries.HasDuplicates() {
		return errors.New("there are duplicate entries")
	}
	for i, entry := range f.Entries {
		if err := entry.Validate(); err != nil {
			return fmt.Errorf("entry nº%d is invalid: %w", i, err)
		}
	}

	return nil
}

func (f File) Equal(f2 File) bool {
	return util.Equal(f, f2)
}

func (f File) Header() Header {
	return f.Entries.Header()
}

func (f *File) Set(id, context string, e Entry) {
	index := f.Entries.IndexByIDAndCtx(id, context)
	if index == -1 {
		f.Entries = append(f.Entries, e)
		return
	}
	f.Entries[index] = e
}

func (f File) LoadByUnifiedID(uid string) string {
	i := f.Entries.IndexByUnifiedID(uid)
	if i == -1 {
		return ""
	}
	return f.Entries[i].Str
}

func (f File) Load(id string, context string) string {
	i := f.Entries.IndexByIDAndCtx(id, context)
	if i == -1 {
		return ""
	}

	return f.Entries[i].Str
}
