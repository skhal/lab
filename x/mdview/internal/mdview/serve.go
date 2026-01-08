// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mdview

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func listenAndServe(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", serveFile)
	s := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	log.Println("start server at", addr)
	return s.ListenAndServe()
}

func serveFile(w http.ResponseWriter, req *http.Request) {
	html, err := renderMarkdown(req.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(html))
}

type Preview struct {
	Info *FileInfo
	HTML template.HTML
}

type FileInfo struct {
	Path    string
	AbsPath string
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
		return nil, fmt.Errorf("%s: fail to stat", path)
	}
	if !fi.Mode().IsRegular() {
		return nil, fmt.Errorf("%s: not a regular file", path)
	}
	return &FileInfo{
		Path:    path,
		AbsPath: abspath,
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
	embfs embed.FS
	tmpls = template.Must(template.New("templates").ParseFS(embfs, "html/*.html"))
)

func execute(p Preview) ([]byte, error) {
	var buf bytes.Buffer
	if err := tmpls.ExecuteTemplate(&buf, "main.html", p); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
