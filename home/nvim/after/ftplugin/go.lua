-- Copyright 2026 Samvel Khalatyan. All rights reserved.
--
-- Use of this source code is governed by a BSD-style
-- license that can be found in the LICENSE file.

-- keep-sorted start
vim.bo.equalprg = "goimports"
vim.bo.expandtab = false
-- keep-sorted end

local function select_and_edit(file, files)
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

local RelatedFile = {
	doc = function()
		local f = vim.fn.expand("%")
		if f:find("doc.go$") then
			return
		end
		f = vim.fs.joinpath(vim.fs.dirname(f), "doc.go")
		vim.cmd.edit(f)
	end,
	source = function()
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
	test = function()
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

vim.keymap.set("n", "<localleader>rd", RelatedFile.doc, { buffer = true })
vim.keymap.set("n", "<localleader>rs", RelatedFile.source, { buffer = true })
vim.keymap.set("n", "<localleader>rt", RelatedFile.test, { buffer = true })
