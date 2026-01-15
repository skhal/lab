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

// splitFirstLineFunc splits a block of data into the first line and the rest.
type splitFirstLineFunc func([]byte) ([]byte, []byte)

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

func newInserter(filename string) (*inserter, error) {
	switch filepath.Ext(filename) {
	case "", ".sh": // no extension: default to shell
		return &inserter{
			prefix:         "#",
			splitFirstLine: splitShebang,
		}, nil
	case ".cc", ".go", ".h":
		return &inserter{
			prefix: "//",
		}, nil
	case ".html", ".md":
		return &inserter{
			start:          "<!--",
			prefix:         " ",
			end:            "-->",
			splitFirstLine: splitDoctype,
		}, nil
	}
	return nil, fmt.Errorf("%s: unsupported file type", filename)
}

func (ins *inserter) Insert(b []byte, lic []byte) ([]byte, error) {
	var buf bytes.Buffer
	b = ins.writeFirstLine(&buf, b)
	ins.writeLicense(&buf, lic)
	buf.Write(b)
	return buf.Bytes(), nil
}

func (ins *inserter) writeFirstLine(buf *bytes.Buffer, b []byte) []byte {
	if ins.splitFirstLine == nil {
		return b
	}
	first, b := ins.splitFirstLine(b)
	if first != nil {
		buf.Write(first)
		buf.WriteByte(eol)
	}
	return b
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

func splitShebang(b []byte) ([]byte, []byte) {
	if len(b) < 2 {
		return nil, b
	}
	if !bytes.HasPrefix(b, []byte("#!")) {
		return nil, b
	}
	before, after, _ := bytes.Cut(b, []byte("\n"))
	return before, after
}

func splitDoctype(b []byte) ([]byte, []byte) {
	before, after, _ := bytes.Cut(b, []byte("\n"))
	if !bytes.HasPrefix(bytes.ToLower(before), []byte("<!doctype")) {
		return nil, b
	}
	return before, after
}
