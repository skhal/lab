-- Copyright 2026 Samvel Khalatyan. All rights reserved.
--
-- Use of this source code is governed by a BSD-style
-- license that can be found in the LICENSE file.

local M = {}

local function git_config_username()
	local cmd = vim.system({ "git", "config", "--get", "user.name" }, { text = true }):wait()
	if cmd.code ~= 0 then
		error("git-config: can't get user.name")
	end
	return (cmd.stdout):gsub("+%s+", "")
end

local function go_package(opts)
	local abspath = vim.fs.abspath(opts.file)
	local dirname = vim.fs.dirname(abspath)
	local pkg = vim.fs.basename(dirname)
	return pkg
end

local default_opts = {
	skel_path = vim.fn.stdpath("data") .. "/lab-skeleton/skel",
	ftgens = {
		["go"] = {
			package = go_package,
		},
		[""] = {
			year = function(_)
				return os.date("%Y")
			end,
			holder = function(_)
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
		pattern = { "*.lua", "*.go" },
		callback = function(ev)
			M.load(ev)
		end,
	})
end

local function load_skeleton(file, subs)
	vim.cmd("0r " .. file)
	for key, val in pairs(subs) do
		vim.cmd("silent! %s/{{" .. key .. "}}/" .. val)
	end
end

local function position_cursor()
	for line_num, line in ipairs(vim.api.nvim_buf_get_lines(0, 0, -1, false)) do
		local from, _ = line:find("{{cursor}}")
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
	local opts = { file = ev.file, filetype = skel.filetype }
	local subs
	ok, subs = pcall(M.gen_substitutes, opts)
	if not ok then
		vim.api.nvim_echo({ { subs, "ErrorMsg" } }, true, {})
		return
	end
	local err
	ok, err = pcall(load_skeleton, skel.path, subs)
	if not ok then
		vim.api.nvim_echo({ { err, "ErrorMsg" } }, true, {})
		return
	end
	ok, err = pcall(position_cursor)
	if not ok then
		vim.api.nvim_echo({ { err, "ErrorMsg" } }, true, {})
		return
	end
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

local function gen_substitutes(gens, opts)
	local subs = {}
	for k, f in pairs(gens) do
		local ok, v = pcall(f, opts)
		if not ok then
			error(("failed generate %s\n%s"):format(opts.ft, k, v))
		end
		subs[k] = v
	end
	return subs
end

local function table_merge(dst, src)
	for k, v in pairs(src) do
		dst[k] = v
	end
end

function M.gen_substitutes(opts)
	local ok, subs = pcall(gen_substitutes, M.ftgens[""] or {}, opts)
	if not ok then
		error(("common substitutes: %s"):format(subs))
	end
	if not next(subs) then
		subs = {}
	end
	local ft_subs
	ok, ft_subs = pcall(gen_substitutes, M.ftgens[opts.filetype] or {}, opts)
	if not ok then
		error(("filetype %s: %s"):format(opts.filetype, ft_subs))
	end
	if next(ft_subs) then
		table_merge(subs, ft_subs)
	end
	if not next(subs) then
		error("failed to generate substitutes")
	end
	return subs
end

return M
