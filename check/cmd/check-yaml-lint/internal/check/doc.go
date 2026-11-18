// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package check applies different formatters and linters depending on whether
// the file belongs to an Ansible project (or collection) or not.
//
// It groups files by Ansible project (a folder with ansible.cfg file), Ansible
// collection (a folder with galaxy.yml), or everything else. Then it applies
// ansible-lint on the first two groups and yamlfmt on the rest.
package check
