// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package copyright

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"regexp"
	"text/template"
)

// ErrNotFound indicates missing or invalid copyright.
var ErrNotFound = errors.New("copyright is not found")

var lineRx = regexp.MustCompile(`^(\s*(/[/\*]|[#"])?) Copyright`)

var (
	//go:embed data
	embedFS      embed.FS
	licenseTmpls = template.Must(template.New("licenses").ParseFS(embedFS, "data/license_*.txt"))
)

// ReadFileFn reads file and returns its content or error.
type ReadFileFn = func(string) ([]byte, error)

// Config configures the runner.
type Config struct {
	ReadFile ReadFileFn
}

// Run checks whether the file include a copyright comment.
func Run(cfg *Config, files []string) error {
	for _, f := range files {
		data, err := cfg.ReadFile(f)
		if err != nil {
			return err
		}
		if err := check(data); err != nil {
			return fmt.Errorf("%s: %w", f, err)
		}
	}
	return nil
}

const eol = '\n'

func check(buf []byte) error {
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
		return
	}
	return true, nil
}

func compileBlockRx(prefix []byte) (*regexp.Regexp, error) {
	data := struct {
		Prefix          string
		EmptyLinePrefix string
	}{
		Prefix: string(prefix),
	}
	// The separator line is empty in comment blocks, e.g. HTML.
	if len(bytes.TrimSpace(prefix)) != 0 {
		data.EmptyLinePrefix = data.Prefix
	}
	var b bytes.Buffer
	if err := licenseTmpls.ExecuteTemplate(&b, "license_bsd.txt", data); err != nil {
		return nil, err
	}
	return regexp.Compile("^" + b.String()) // must match the beginning
}
