-- Copyright 2026 Samvel Khalatyan. All rights reserved.
--
-- Use of this source code is governed by a BSD-style
-- license that can be found in the LICENSE file.

local selector = {
	---Select searches for files matching the pattern at file's location,
	---optionally filtered by the filter, and invokes the callback on the
	---selected item.
	---
	---It does nothing of no files are found and skips the selection dialog if
	---there is only one file available.
	---
	---@param file string file with path for which to search for related files
	---@param opts { pattern: string, filter?: function } define selection file pattern and optional filter function.
	---@param callback function called on selection choice.
	Select = function(file, opts, callback)
		local path = vim.fs.dirname(file)
		local opt_nosuf = false
		local opt_list = true
		local files = vim.fn.globpath(path, opts.pattern, opt_nosuf, opt_list)
		files = vim.tbl_filter(function(item)
			if vim.fs.basename(file) == vim.fs.basename(item) then
				return false
			end
			if opts.filter then
				return opts.filter(item)
			end
			return true
		end, files)
		if not next(files) then
			vim.api.nvim_echo({ { "no files found" } }, true, {})
			return
		end
		if #files == 1 then
			callback(files[1])
			return
		end
		vim.ui.select(files, {
			prompt = "Select file from " .. path .. ":",
			format_item = function(item)
				return vim.fs.basename(item)
			end,
		}, function(choice)
			callback(choice)
		end)
	end,
}

return selector
