// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package check

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

const (
	ansibleProjectConfig    = "ansible.cfg"
	ansibleCollectionConfig = "galaxy.yml"
)
const ansibleProjectLookupDepth = 5

// Collect groups files into Ansible and non-Ansible blocks. Ansible files are
// grouped by Ansible project or collection. An Ansible project is a folder
// with ansible.cfg whereas a collection is a folder with galaxy.yml file.
func Collect(fsys fs.FS, files Files) (map[Base]Files, Files) {
	collector := newCollector(fsys)

	notCollectedFiles := collector.Collect(files)

	return collector.Groups(), notCollectedFiles
}

// BaseKind identifies the type of Ansible base folder: a project or a
// collection.
//
//go:generate stringer -type BaseKind -linecomment
type BaseKind int

const (
	_ BaseKind = iota
	// keep-sorted start
	BaseCollection // collection
	BaseProject    // project
	// keep-sorted end
)

// Base is the path to either Ansible project or a collection.
type Base struct {
	Path string   // path to the base
	Kind BaseKind // kind of the base: project or collection
}

// String implements [fmt.Stringer] interface.
func (b Base) String() string {
	return fmt.Sprintf("%s %s", b.Path, b.Kind)
}

type collector struct {
	groups map[Base]Files
	fsys   fs.FS
}

func newCollector(fsys fs.FS) *collector {
	return &collector{
		groups: make(map[Base]Files),
		fsys:   fsys,
	}
}

// Collect filters Ansible files and returns a list of left-over, non-Ansible
// files.
func (ctr *collector) Collect(files []string) []string {
	var ff []string
	for _, f := range files {
		if !ctr.collect(f) {
			ff = append(ff, f)
		}
	}
	return ff
}

func (ctr *collector) collect(file string) bool {
	ok, base := ctr.base(file)
	if !ok {
		return false
	}
	f, err := filepath.Rel(base.Path, file)
	if err != nil {
		f = file
	}
	ctr.groups[base] = append(ctr.groups[base], f)
	return true
}

func (ctr *collector) base(file string) (bool, Base) {
	path := file
	for range ansibleProjectLookupDepth {
		path = filepath.Dir(path)
		if ctr.isAnsibleProject(path) {
			return true, Base{Path: path, Kind: BaseProject}
		}
		if ctr.isAnsibleCollection(path) {
			return true, Base{Path: path, Kind: BaseCollection}
		}
		if path == "." {
			break
		}
	}
	return false, Base{}
}

func (ctr *collector) isAnsibleProject(path string) bool {
	_, err := fs.Stat(ctr.fsys, filepath.Join(path, ansibleProjectConfig))
	return err == nil
}

func (ctr *collector) isAnsibleCollection(path string) bool {
	_, err := fs.Stat(ctr.fsys, filepath.Join(path, ansibleCollectionConfig))
	return err == nil
}

// Groups returns a map of Ansible files grouped by either project or
// collection.
func (ctr collector) Groups() map[Base]Files { return ctr.groups }
