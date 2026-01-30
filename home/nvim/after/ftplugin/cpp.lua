-- Copyright 2026 Samvel Khalatyan. All rights reserved.
--
-- Use of this source code is governed by a BSD-style
-- license that can be found in the LICENSE file.

vim.bo.equalprg = "clang-format21 -assume-filename=%"

local M = {
	header = function()
		local f = vim.fn.expand("%")
		if not f:find("%.cc$") then
			return
		end
		f = f:gsub("_test%.cc$", ".h")
		f = f:gsub("%.cc$", ".h")
		vim.cmd.edit(f)
	end,
	source = function()
		local f = vim.fn.expand("%")
		if f:find("%.h$") then
			f = f:gsub("%.h$", ".cc")
		elseif f:find("_test%.cc$") then
			f = f:gsub("_test%.cc$", ".cc")
		else
			return
		end
		vim.cmd.edit(f)
	end,
	test = function()
		local f = vim.fn.expand("%")
		if f:find("_test%.cc$") then
			return
		end
		if f:find("%.h$") then
			f = f:gsub("%.h$", "_test.cc")
		elseif f:find("%.cc") then
			f = f:gsub("%.cc$", "_test.cc")
		else
			return
		end
		vim.cmd.edit(f)
	end,
}

vim.keymap.set("n", "<localleader>rh", M.header, { buffer = true })
vim.keymap.set("n", "<localleader>rs", M.source, { buffer = true })
vim.keymap.set("n", "<localleader>rt", M.test, { buffer = true })
