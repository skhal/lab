// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package license

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Add generates a license with current date and holder, and injets it into text
// block b. It uses filename extension to detect comment syntax for license,
// defaulting to shell if the extension misses.
func Add(b []byte, filename, holder string) ([]byte, error) {
	lic, err := genLicenseBlock(LicenseData{
		Year:   strconv.FormatInt(int64(time.Now().Year()), 10),
		Holder: holder,
	})
	if err != nil {
		return nil, err
	}
	ins, err := newInserter(filename)
	if err != nil {
		return nil, err
	}
	return ins.Insert(b, lic)
}

// splitFirstLineFunc splits a block of data into the first line, the rest, and
// optional separator line.
type splitFirstLineFunc func([]byte) splitBlock

// inserter injects a license into a block of text.
//
//			first-line
//	   [start]
//	   [prefix] Copyright ...
//	   [end]
//	   text
type inserter struct {
	start  string // start of the comment, e.g. "<!--" for HTML
	prefix string // comment prefix, e.g. "#" for shell
	end    string // end of the comment, e.g. "-->" for HTML

	// split first line if Shell shebang, HTML doctype, etc.
	splitFirstLine splitFirstLineFunc
}

var (
	// keep-sorted start block=yes
	insC = inserter{
		prefix: "//",
	}
	insHTML = inserter{
		start:          "<!--",
		prefix:         " ",
		end:            "-->",
		splitFirstLine: splitDoctype,
	}
	insLua = inserter{
		prefix: "--",
	}
	insShell = inserter{
		prefix:         "#",
		splitFirstLine: splitShebang,
	}
	insShellNoSplit = inserter{
		prefix: "#",
	}
	insVim = inserter{
		prefix: "\"",
	}
	// keep-sorted end
)

func newInserter(filename string) (*inserter, error) {
	switch filepath.Ext(filename) {
	// keep-sorted start
	case "", ".sh": // no extension: default to shell
		return &insShell, nil
	case ".bazel", ".conf", ".ctags", ".txt", ".txtpb", ".yaml":
		return &insShellNoSplit, nil
	case ".cc", ".go", ".h", ".proto":
		return &insC, nil
	case ".html", ".md":
		return &insHTML, nil
	case ".lua":
		return &insLua, nil
	case ".vim":
		return &insVim, nil
		// keep-sorted end
	}
	base := filepath.Base(filename)
	switch base {
	case ".clangd":
		return &insShellNoSplit, nil
	}
	switch {
	case strings.HasPrefix(base, ".bazel"):
		return &insShellNoSplit, nil
	}
	return nil, fmt.Errorf("%s: unsupported file type", filename)
}

// Insert injects a licence block lic into a data block b.
func (ins *inserter) Insert(b []byte, lic []byte) ([]byte, error) {
	var buf bytes.Buffer
	b = ins.writeFirstLine(&buf, b)
	ins.writeLicense(&buf, lic)
	if len(b) > 0 {
		buf.WriteByte(eol) // separate license
		buf.Write(b)
	}
	return buf.Bytes(), nil
}

func (ins *inserter) writeFirstLine(buf *bytes.Buffer, b []byte) []byte {
	if ins.splitFirstLine == nil {
		return b
	}
	sb := ins.splitFirstLine(b)
	if sb.first != nil {
		buf.Write(sb.first)
		buf.WriteByte(eol)
	}
	if sb.separator != nil {
		buf.Write(sb.separator)
		buf.WriteByte(eol)
	}
	return sb.rest
}

func (ins *inserter) writeLicense(buf *bytes.Buffer, lic []byte) {
	if ins.start != "" {
		buf.WriteString(ins.start)
		buf.WriteByte(eol)
	}
	for ln := range bytes.Lines(lic) {
		if len(bytes.TrimSpace(ln)) == 0 {
			ins.writeLicenseEmptyLine(buf)
			continue
		}
		buf.WriteString(ins.prefix)
		buf.WriteByte(' ')
		buf.Write(ln)
	}
	if ins.end != "" {
		buf.WriteString(ins.end)
		buf.WriteByte(eol)
	}
}

func (ins *inserter) writeLicenseEmptyLine(buf *bytes.Buffer) {
	if strings.TrimSpace(ins.prefix) == "" {
		buf.WriteByte(eol)
		return
	}
	buf.WriteString(ins.prefix)
	buf.WriteByte(eol)
}

type splitBlock struct {
	first     []byte
	rest      []byte
	separator []byte // separator to insert b/w first and rest
}

func splitShebang(b []byte) splitBlock {
	if len(b) < 2 {
		return splitBlock{
			rest: b,
		}
	}
	if !bytes.HasPrefix(b, []byte("#!")) {
		return splitBlock{
			rest: b,
		}
	}
	before, after, _ := bytes.Cut(b, []byte("\n"))
	return splitBlock{
		first:     before,
		rest:      after,
		separator: []byte("#"),
	}
}

func splitDoctype(b []byte) splitBlock {
	before, after, _ := bytes.Cut(b, []byte("\n"))
	if !bytes.HasPrefix(bytes.ToLower(before), []byte("<!doctype")) {
		return splitBlock{
			rest: b,
		}
	}
	return splitBlock{
		first: before,
		rest:  after,
	}
}
