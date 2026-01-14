// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mdview

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func listenAndServe(addr string) error {
	mux := http.NewServeMux()
	mux.Handle("/", handle(serveFile))
	s := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	log.Println("start server at", addr)
	return s.ListenAndServe()
}

const readmeName = "README.md"

type handle func(w http.ResponseWriter, req *http.Request) error

func (h handle) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if err := h(w, req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
	}
}

func serveFile(w http.ResponseWriter, req *http.Request) error {
	p := req.URL.Path
	if strings.HasSuffix(p, "/") {
		p += readmeName
	}
	p = filepath.Clean("./" + p) // make local
	html, err := renderMarkdown(p)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(html))
	return nil
}

type Preview struct {
	Info *FileInfo
	HTML template.HTML
}

type FileInfo struct {
	Path    string
	AbsPath string

	Atime time.Time
	Mtime time.Time
}

const extMarkdown = ".md"

func renderMarkdown(path string) ([]byte, error) {
	info, err := stat(path)
	if err != nil {
		return nil, err
	}
	data, err := readAndRender(info)
	if err != nil {
		return nil, err
	}
	html, err := execute(Preview{Info: info, HTML: template.HTML(data)})
	if err != nil {
		return nil, err
	}
	return html, nil
}

func stat(path string) (*FileInfo, error) {
	if extMarkdown != filepath.Ext(path) {
		return nil, fmt.Errorf("%s: not markdown", path)
	}
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	abspath := filepath.Join(pwd, path)
	fi, err := os.Stat(abspath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("%s: does not exist", path)
		}
		return nil, fmt.Errorf("%s: fail to stat", path)
	}
	if !fi.Mode().IsRegular() {
		return nil, fmt.Errorf("%s: not a regular file", path)
	}
	return &FileInfo{
		Path:    path,
		AbsPath: abspath,
		Atime:   time.Now(),
		Mtime:   fi.ModTime(),
	}, nil
}

func readAndRender(info *FileInfo) ([]byte, error) {
	data, err := os.ReadFile(info.AbsPath)
	if err != nil {
		return nil, fmt.Errorf("%s: fail to read", info.Path)
	}
	return Render(data), nil
}

var (
	//go:embed html
	embfs   embed.FS
	funcMap = template.FuncMap{
		"FormatRFC822": func(t time.Time) string {
			return t.Format(time.RFC822)
		},
	}
	tmpls = template.Must(template.New("templates").Funcs(funcMap).ParseFS(embfs, "html/*.html"))
)

func execute(p Preview) ([]byte, error) {
	var buf bytes.Buffer
	if err := tmpls.ExecuteTemplate(&buf, "main.html", p); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
