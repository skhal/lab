-- Copyright 2026 Samvel Khalatyan. All rights reserved.
--
-- Use of this source code is governed by a BSD-style
-- license that can be found in the LICENSE file.

local selector = require("file.selector")

local function selectSource(file)
	local filter = function(item)
		if item:find("doc%.go$") then
			return false
		end
		if item:find("_test%.go$") then
			return false
		end
		if item:find("_string%.go$") then
			return false
		end
		return true
	end
	selector.Select(file, { pattern = "*.go", filter = filter }, vim.cmd.edit)
end

local function selectTest(file)
	selector.Select(file, { pattern = "*_test.go" }, vim.cmd.edit)
end

local relatedFile = {
	Doc = function()
		local f = vim.fn.expand("%")
		if f:find("doc.go$") then
			return
		end
		f = vim.fs.joinpath(vim.fs.dirname(f), "doc.go")
		vim.cmd.edit(f)
	end,
	Source = function()
		local f = vim.fn.expand("%")
		if f:find("doc%.go$") then
			selectSource(f)
			return
		end
		if not f:find("_test%.go$") then
			selectSource(f)
			return
		end
		f = f:gsub("_test%.go$", ".go")
		vim.cmd.edit(f)
	end,
	Test = function()
		local f = vim.fn.expand("%")
		if f:find("doc%.go$") then
			selectTest(f)
			return
		end
		if f:find("_test%.go$") then
			return
		end
		f = f:gsub("%.go$", "_test.go")
		vim.cmd.edit(f)
	end,
}

return relatedFile
