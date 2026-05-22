-- Copyright 2026 Samvel Khalatyan. All rights reserved.
--
-- Use of this source code is governed by a BSD-style
-- license that can be found in the LICENSE file.

local function select_and_edit(file, files)
	if #files == 1 then
		vim.cmd.edit(files[1])
		return
	end
	vim.ui.select(files, {
		prompt = "Select file from " .. vim.fs.dirname(file) .. ":",
		format_item = function(item)
			return vim.fs.basename(item)
		end,
	}, function(choice)
		vim.cmd.edit(choice)
	end)
end

local function select_source(file)
	local opt_nosuf = false
	local opt_list = true
	local files = vim.fn.globpath(vim.fs.dirname(file), "*.go", opt_nosuf, opt_list)
	files = vim.tbl_filter(function(item)
		if vim.fs.basename(file) == vim.fs.basename(item) then
			return false
		end
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
	end, files)
	if not next(files) then
		vim.api.nvim_echo({
			{ "no source files found" },
		}, true, {})
		return
	end
	select_and_edit(file, files)
end

local function select_test(file)
	local opt_nosuf = false
	local opt_list = true
	local files = vim.fn.globpath(vim.fs.dirname(file), "*_test.go", opt_nosuf, opt_list)
	if not next(files) then
		vim.api.nvim_echo({
			{ "no test files found" },
		}, true, {})
		return
	end
	select_and_edit(file, files)
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
			select_source(f)
			return
		end
		if not f:find("_test%.go$") then
			select_source(f)
			return
		end
		f = f:gsub("_test%.go$", ".go")
		vim.cmd.edit(f)
	end,
	Test = function()
		local f = vim.fn.expand("%")
		if f:find("doc%.go$") then
			select_test(f)
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
