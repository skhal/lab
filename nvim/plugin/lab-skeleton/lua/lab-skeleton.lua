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

function load_skeleton(file, subs)
	vim.cmd("0r " .. file)
	for key, val in pairs(subs) do
		vim.cmd("silent! %s/" .. key .. "/" .. val)
	end
end

function position_cursor()
	for line_num, line in ipairs(vim.api.nvim_buf_get_lines(0, 0, -1, false)) do
		local from, to = line:find("{{cursor}}")
		if from ~= nil then
			vim.cmd("silent! %s/{{cursor}}//g")
			vim.api.nvim_win_set_cursor(0, { line_num, from - 1 })
			break
		end
	end
end

function M.load(ev)
	local ok, skel_path = pcall(M.pick_skeleton, ev.file)
	if not ok then
		vim.api.nvim_echo({ { skel_path, "ErrorMsg" } }, true, {})
		return
	end
	local ok, subs = pcall(M.gen_substitutes)
	if not ok then
		vim.api.nvim_echo({ { subs, "ErrorMsg" } }, true, {})
		return
	end
	local ok, err = pcall(load_skeleton, skel_path, subs)
	if not ok then
		vim.api.nvim_echo({ { err, "ErrorMsg" } }, true, {})
		return
	end
	local ok, err = pcall(position_cursor)
end

function M.pick_skeleton(file)
	local ext = vim.fn.fnamemodify(file, ":e")
	local skel_path = M.skel_path .. "/new." .. ext
	if not (vim.uv or vim.loop).fs_stat(skel_path) then
		error(("skeleton %s: does not exist."):format(skel_path))
	end
	return skel_path
end

function M.gen_substitutes()
	local cmd = vim.system({ "git", "config", "--get", "user.name" }, { text = true }):wait()
	if cmd.code ~= 0 then
		error("git-config: can't get user.name")
	end
	return {
		["{{year}}"] = os.date("%Y"),
		["{{holder}}"] = (cmd.stdout):gsub("+%s+", ""),
	}
end

return M
