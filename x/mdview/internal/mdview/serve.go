// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mdview

import (
	"fmt"
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
	preview, err := renderMarkdown(req.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(preview.HTML))
}

type Preview struct {
	Info *FileInfo
	HTML []byte
}

type FileInfo struct {
	Path    string
	AbsPath string
}

const extMarkdown = ".md"

func renderMarkdown(path string) (*Preview, error) {
	info, err := stat(path)
	if err != nil {
		return nil, err
	}
	html, err := readAndRender(info)
	if err != nil {
		return nil, err
	}
	return &Preview{
		Info: info,
		HTML: html,
	}, nil
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
