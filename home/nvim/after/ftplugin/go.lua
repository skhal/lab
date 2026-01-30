-- Copyright 2026 Samvel Khalatyan. All rights reserved.
--
-- Use of this source code is governed by a BSD-style
-- license that can be found in the LICENSE file.

-- keep-sorted start
vim.bo.equalprg = "goimports"
vim.bo.expandtab = false
-- keep-sorted end
--
local M = {
	source = function()
		local f = vim.fn.expand("%")
		if not f:find("_test%.go$") then
			return
		end
		f = f:gsub("_test%.go$", ".go")
		vim.cmd.edit(f)
	end,
	test = function()
		local f = vim.fn.expand("%")
		if f:find("_test%.go$") then
			return
		end
		f = f:gsub("%.go$", "_test.go")
		vim.cmd.edit(f)
	end,
}

vim.keymap.set("n", "<localleader>rs", M.source, { buffer = true })
vim.keymap.set("n", "<localleader>rt", M.test, { buffer = true })
