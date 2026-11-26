-- Copyright 2026 Samvel Khalatyan. All rights reserved.
--
-- Use of this source code is governed by a BSD-style
-- license that can be found in the LICENSE file.

local M = {
	source = function()
		local f = vim.fn.expand("%")
		if f:find("tests/.*_test%.py$") then
			f = f:gsub("tests/(.*)_test%.py$", "src/%1.py")
		else
			return
		end
		vim.cmd.edit(f)
	end,
	test = function()
		local f = vim.fn.expand("%")
		if f:find("src/.*%.py$") then
			f = f:gsub("src/(.*)%.py$", "tests/%1_test.py")
		else
			return
		end
		vim.cmd.edit(f)
	end,
}

vim.keymap.set("n", "<localleader>rs", M.source, { buffer = true })
vim.keymap.set("n", "<localleader>rt", M.test, { buffer = true })
