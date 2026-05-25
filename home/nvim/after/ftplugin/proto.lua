-- Copyright 2026 Samvel Khalatyan. All rights reserved.
--
-- Use of this source code is governed by a BSD-style
-- license that can be found in the LICENSE file.

vim.bo.equalprg = "clang-format21 -assume-filename=%"

local selector = require("file.selector")

local function selectProto()
	local file = vim.fn.expand("%")
	selector.Select(file, { pattern = "*.proto" }, vim.cmd.edit)
end

vim.keymap.set("n", "<localleader>rs", selectProto, { buffer = true })

local generate = require("go.generate")

vim.keymap.set("n", "<localleader>gen", generate.RunPackage, { buffer = true })
