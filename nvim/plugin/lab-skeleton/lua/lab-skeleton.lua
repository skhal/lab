-- Copyright 2026 Samvel Khalatyan. All rights reserved.
--
-- Use of this source code is governed by a BSD-style
-- license that can be found in the LICENSE file.

local M = {}

default_opts = {
	skel_path = vim.fn.stdpath("data") .. "/lab-skeleton/skel",
}

function M.setup(opts)
	M.skel_path = opts.skel_path or default_opts.skel_path
	M.augroup = vim.api.nvim_create_augroup("LabSkeleton", { clear = true })

	vim.api.nvim_create_autocmd("BufNewFile", {
		group = M.augroup,
		desc = "Load template",
		pattern = { "*.lua" },
		callback = function(ev)
			M.load(ev)
		end,
	})
end

function M.load(ev)
	local ext = vim.fn.fnamemodify(ev.file, ":e")
	local skel_path = M.skel_path .. "/new." .. ext
	if not (vim.uv or vim.loop).fs_stat(skel_path) then
		vim.api.nvim_echo({
			{ ("skeleton %s: does not exist.\n"):format(skel_path), "ErrorMsg" },
		}, true, {})
		return
	end
	local obj = vim.system({ "git", "config", "--get", "user.name" }, { text = true }):wait()
	if obj.code ~= 0 then
		vim.api.nvim_echo({
			{ "Can't get user name", "ErrorMsg" },
		}, true, {})
		return
	end
	local year = os.date("%Y")
	local holder = (obj.stdout):gsub("+%s+", "")
	local substitutes = {
		["{{year}}"] = year,
		["{{holder}}"] = holder,
	}
	vim.cmd("0r " .. skel_path)
	for key, val in pairs(substitutes) do
		vim.cmd("silent! %s/" .. key .. "/" .. val)
	end
	for line_num, line in ipairs(vim.api.nvim_buf_get_lines(0, 0, -1, false)) do
		local from, to = line:find("{{cursor}}")
		if from ~= nil then
			vim.cmd("silent! %s/{{cursor}}//g")
			vim.api.nvim_win_set_cursor(0, { line_num, from - 1 })
			break
		end
	end
	vim.api.nvim_echo({
		{ ("skel: %s\nyear: %d\nholder: %s"):format(skel_path, year, holder), "Normal" },
	}, true, {})
end

return M
