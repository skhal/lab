-- Copyright 2026 Samvel Khalatyan. All rights reserved.
--
-- Use of this source code is governed by a BSD-style
-- license that can be found in the LICENSE file.

local function open(path)
	vim.cmd.Lexplore({ args = { path }, range = { 25 } })
	vim.notify(("open %s"):format(path))
end

local M = {
	Open = function()
		local f = vim.fn.expand("%")
		open(vim.fs.dirname(f))
	end,
}

return M
