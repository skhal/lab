-- Copyright 2026 Samvel Khalatyan. All rights reserved.
--
-- Use of this source code is governed by a BSD-style
-- license that can be found in the LICENSE file.

local M = {}

default_opts = {
	skel_path = vim.fn.stdpath("data") .. "/lab-skeleton/skel",
	ftgens = {
		[""] = {
			year = function()
				return os.date("%Y")
			end,
			holder = function()
				return git_config_username()
			end,
		},
	},
}

function M.setup(opts)
	M.skel_path = opts.skel_path or default_opts.skel_path
	M.ftgens = default_opts.ftgens
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
		vim.cmd("silent! %s/{{" .. key .. "}}/" .. val)
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
	local ok, skel = pcall(M.pick_skeleton, ev.file)
	if not ok then
		vim.api.nvim_echo({ { skel.path, "ErrorMsg" } }, true, {})
		return
	end
	local ok, subs = pcall(M.gen_substitutes, skel.filetype)
	if not ok then
		vim.api.nvim_echo({ { subs, "ErrorMsg" } }, true, {})
		return
	end
	local ok, err = pcall(load_skeleton, skel.path, subs)
	if not ok then
		vim.api.nvim_echo({ { err, "ErrorMsg" } }, true, {})
		return
	end
	local ok, err = pcall(position_cursor)
end

function M.pick_skeleton(file)
	local ft = vim.filetype.match({ filename = file })
	local skel_path = M.skel_path .. "/new." .. ft
	if not (vim.uv or vim.loop).fs_stat(skel_path) then
		error(("skeleton %s: does not exist."):format(skel_path))
	end
	return {
		path = skel_path,
		filetype = ft,
	}
end

function git_config_username()
	local cmd = vim.system({ "git", "config", "--get", "user.name" }, { text = true }):wait()
	if cmd.code ~= 0 then
		error("git-config: can't get user.name")
	end
	return (cmd.stdout):gsub("+%s+", "")
end

function M.gen_substitutes(filetype)
	local subs = {}
	for ft, gens in pairs(M.ftgens) do
		for k, f in pairs(gens) do
			local ok, v = pcall(f)
			if not ok then
				error(("filetype %s: failed generate a substitute %s\n%s"):format(ft, k, v))
			end
			subs[k] = v
		end
	end
	if not next(subs) then
		error("failed to generate substitutes")
	end
	return subs
end

return M
