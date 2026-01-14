// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package license

import (
	"bytes"
	"embed"
	"errors"
	"regexp"
	"strings"
	"text/template"
)

var lineRx = regexp.MustCompile(`^([\t ]*(/[/\*]|[#"])?) Copyright`)

var (
	//go:embed data
	embedFS      embed.FS
	licenseTmpls = template.Must(template.New("licenses").ParseFS(embedFS, "data/license_*.txt"))
)

const eol = '\n'

func Check(buf []byte) (err error) {
	ln := 1
	defer func() {
		if errors.Is(err, ErrInvalid) {
			err = NewInvalidError(ln)
		}
	}()
	for len(buf) > 0 {
		if ok, err := match(buf); ok {
			return nil
		} else if err != nil {
			return err
		}
		idx := bytes.IndexByte(buf, eol)
		if idx == -1 {
			break
		}
		idx += 1 // skip eol
		buf = buf[idx:]
		ln += 1
	}
	return ErrNotFound
}

func match(buf []byte) (ok bool, err error) {
	matches := lineRx.FindSubmatch(buf)
	if matches == nil {
		return
	}
	blockRx, err := compileBlockRx(matches[1])
	if err != nil {
		return
	}
	if !blockRx.Match(buf) {
		return false, ErrInvalid
	}
	return true, nil
}

func compileBlockRx(prefix []byte) (*regexp.Regexp, error) {
	b, err := genLicenseBlock(prefix)
	if err != nil {
		return nil, err
	}
	return regexp.Compile("^" + string(b)) // must match the beginning
}

// LicenseData is input to the license template.
type LicenseData struct {
	Prefix string
	Year   string
	Holder string
}

// EmptyLinePrefix returns prefix for empty lines in the license. It is empty
// if the prefix is space-only, e.g. HTML comments.
func (ld *LicenseData) EmptyLinePrefix() string {
	if len(strings.TrimSpace(ld.Prefix)) == 0 {
		return ""
	}
	return ld.Prefix
}

func genLicenseBlock(prefix []byte) ([]byte, error) {
	data := &LicenseData{
		Prefix: string(prefix),
		Year:   `\d{4}`,
		Holder: `\w+( \w+)?`,
	}
	var b bytes.Buffer
	if err := licenseTmpls.ExecuteTemplate(&b, "license_bsd.txt", data); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
