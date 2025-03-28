package compiler

import (
	"bufio"
	"bytes"
	bin "encoding/binary"
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/Tom5521/xgotext/internal/util"
	"github.com/Tom5521/xgotext/pkg/po"
)

// Aliase this bc I'm too lazy to write "uint32" every time I want to use it.
type u32 = uint32

var (
	magicNumber = func() u32 {
		if util.IsBigEndian {
			return util.BigEndianMagicNumber
		}
		return util.LittleEndianMagicNumber
	}()
	order = bin.NativeEndian
)

const (
	eot = "\x04"
	nul = "\x00"
)

var _ po.Compiler = (*MoCompiler)(nil)

type MoCompiler struct {
	File   *po.File
	Config MoConfig
}

func NewMo(file *po.File, opts ...MoOption) MoCompiler {
	c := MoCompiler{
		File:   file,
		Config: DefaultMoConfig(opts...),
	}

	return c
}

func (mc *MoCompiler) SetFile(f *po.File) {
	mc.File = f
}

func (mc *MoCompiler) ToWriterWithOptions(w io.Writer, opts ...MoOption) error {
	mc.Config.ApplyOptions(opts...)
	defer mc.Config.RestoreLastCfg()
	return mc.ToWriter(w)
}

func (mc *MoCompiler) ToBytesWithOptions(options ...MoOption) []byte {
	mc.Config.ApplyOptions(options...)
	defer mc.Config.RestoreLastCfg()
	return mc.ToBytes()
}

func (mc *MoCompiler) ToFileWithOptions(f string, options ...MoOption) error {
	mc.Config.ApplyOptions(options...)
	defer mc.Config.RestoreLastCfg()
	return mc.ToFile(f)
}

// A len() function with fixed-size return.
func flen(value any) u32 {
	return u32(reflect.ValueOf(value).Len())
}

// Code translated from: https://github.com/izimobil/polib/blob/master/polib.py#L553
func (mc *MoCompiler) writeTo(writer io.Writer) error {
	entries := mc.File.Entries.Solve().CleanFuzzy().CleanObsoletes()

	if mc.Config.Sort {
		entries = mc.Config.SortMode.SortMethod(entries)()
	}

	var (
		offsets   []u32
		ids, strs string
	)

	for _, e := range entries {
		msgid := e.UnifiedID()
		msgstr := e.UnifiedStr()

		offsets = append(offsets,
			flen(ids),
			flen(msgid),
			flen(strs),
			flen(msgstr),
		)
		ids += msgid + nul
		strs += msgstr + nul
	}

	keystart := 7*4 + 16*flen(entries)
	valuestart := keystart + flen(ids)

	var koffsets, voffsets []u32

	for i := 0; i < len(offsets); i += 4 {
		o1 := offsets[i]
		l1 := offsets[i+1]
		o2 := offsets[i+2]
		l2 := offsets[i+3]

		koffsets = append(koffsets, l1, o1+keystart)
		voffsets = append(voffsets, l2, o2+valuestart)
	}

	// Write the data to the file
	data := []any{
		magicNumber,
		u32(0),
		flen(entries),
		u32(7 * 4),            // Offset of the original strings table
		7*4 + flen(entries)*8, // Offset of the translated strings table
		u32(0), keystart,
		koffsets,
		voffsets,
		[]byte(ids),
		[]byte(strs),
	}

	for _, v := range data {
		err := bin.Write(writer, order, v)
		if err != nil && !mc.Config.IgnoreErrors {
			return err
		}
	}

	return nil
}

func (mc MoCompiler) ToWriter(w io.Writer) error {
	buf := bufio.NewWriter(w)
	err := mc.writeTo(buf)

	if err != nil && !mc.Config.IgnoreErrors {
		return err
	}

	err = buf.Flush()
	if err != nil && !mc.Config.IgnoreErrors {
		return err
	}

	return nil
}

func (mc MoCompiler) ToFile(f string) error {
	if mc.Config.Verbose {
		mc.Config.Logger.Println("Opening file...")
	}
	// Open the file with the determined flags.
	flags := os.O_WRONLY | os.O_TRUNC | os.O_CREATE
	if !mc.Config.Force {
		flags |= os.O_EXCL
	}
	file, err := os.OpenFile(f, flags, os.ModePerm)
	if err != nil && !mc.Config.IgnoreErrors {
		err = fmt.Errorf("error opening file: %w", err)
		mc.Config.Logger.Println("ERROR:", err)
		return err
	}
	defer file.Close()

	if mc.Config.Verbose {
		mc.Config.Logger.Println("Cleaning file contents...")
	}

	// Write compiled translations to the file.
	return mc.ToWriter(file)
}

func (mc MoCompiler) ToBytes() []byte {
	var b bytes.Buffer

	mc.ToWriter(&b)

	return b.Bytes()
}
