-- Copyright 2026 Samvel Khalatyan. All rights reserved.
--
-- Use of this source code is governed by a BSD-style
-- license that can be found in the LICENSE file.

local M = {}

local function find_skeleton(_, ft)
	return "new." .. ft
end

local function find_go_skeleton(file, _)
	local skel = "new.go"
	if file:find("_test%.go$") ~= nil then
		skel = "new_test.go"
	end
	return skel
end

local function find_c_skeleton(file, _)
	local skel = "new.c"
	if file:find("main%.c$") ~= nil then
		skel = "new_main.c"
	end
	return skel
end

local function find_cpp_skeleton(file, _)
	local skel = "new.cc"
	if file:find("_test%.cc$") ~= nil then
		skel = "new_test.cc"
	elseif file:find("%.h$") ~= nil then
		skel = "new.h"
	end
	return skel
end

local function git_config_username()
	local cmd = { "git", "config", "--get", "user.name" }
	local obj = vim.system(cmd, { text = true }):wait()
	if obj.code ~= 0 then
		error("git: can't get user.name")
	end
	return (obj.stdout):gsub("%s+$", "")
end

local function git_worktree_path()
	local cmd = { "git", "rev-parse", "--show-toplevel" }
	local obj = vim.system(cmd, { text = true }):wait()
	if obj.code ~= 0 then
		error("git: can't get worktree path")
	end
	return (obj.stdout):gsub("%s+$", "")
end

local function git_relpath(file)
	local worktree = git_worktree_path()
	local abspath = vim.fs.abspath(file)
	local relpath = vim.fs.relpath(worktree, abspath)
	if not relpath then
		-- file is outside the worktree
		relpath = file
	end
	return relpath
end

local function c_header(opts)
	local relpath = git_relpath(opts.file)
	relpath = relpath:gsub("_test%.c$", ".c")
	relpath = relpath:gsub("%.c$", ".h")
	return vim.fn.escape(relpath, "/")
end

local function cpp_guard(opts)
	local relpath = git_relpath(opts.file)
	local guard = relpath:gsub("^/", "") -- if outside of a git worktree
	guard = guard:gsub("[/%.]", "_") .. "_"
	return guard:upper()
end

local function cpp_header(opts)
	local relpath = git_relpath(opts.file)
	relpath = relpath:gsub("_test%.cc$", ".cc")
	relpath = relpath:gsub("%.cc$", ".h")
	return vim.fn.escape(relpath, "/")
end

local function cpp_namespace(opts)
	local relpath = git_relpath(opts.file)
	local ns = vim.fs.dirname(relpath)
	ns = ns:gsub("^/", "") -- if outside of a git worktree
	return ns:gsub("/", "::")
end

local function go_package(opts)
	local abspath = vim.fs.abspath(opts.file)
	local dirname = vim.fs.dirname(abspath)
	local pkg = vim.fs.basename(dirname)
	return pkg
end

local default_opts = {
	skel_path = vim.fn.stdpath("data") .. "/lab-skeleton/skel",
	find = {
		c = find_c_skeleton,
		cpp = find_cpp_skeleton,
		go = find_go_skeleton,
		[""] = find_skeleton,
	},
	ftgens = {
		c = {
			header = c_header,
		},
		cpp = {
			guard = cpp_guard,
			header = cpp_header,
			namespace = cpp_namespace,
		},
		go = {
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
	M.find = default_opts.find
	M.ftgens = default_opts.ftgens
	M.augroup = vim.api.nvim_create_augroup("LabSkeleton", { clear = true })
	vim.api.nvim_create_autocmd("BufNewFile", {
		group = M.augroup,
		desc = "Load template",
		pattern = { "*.c", "*.cc", "*.go", "*.h", "*.lua" },
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
	local ok, skel = pcall(M.find_skeleton, ev.file)
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
	local msgs = {
		{ ("skel: %s\n"):format(skel.path), "Normal" },
		{ ("ft: %s\n"):format(skel.filetype), "Normal" },
	}
	if next(subs) then
		for k, v in pairs(subs) do
			table.insert(msgs, { (". %s: %s\n"):format(k, v), "Normal" })
		end
	end
	vim.api.nvim_echo(msgs, true, {})
end

function M.find_skeleton(file)
	local ft = vim.filetype.match({ filename = file })
	local find = M.find[ft] or M.find[""]
	local name = find(file, ft)
	local path = M.skel_path .. "/" .. name
	if not (vim.uv or vim.loop).fs_stat(path) then
		error(("skeleton %s: does not exist."):format(name))
	end
	return {
		filetype = ft,
		path = path,
	}
end

local function gen_substitutes(gens, opts)
	local subs = {}
	for k, f in pairs(gens) do
		local ok, v = pcall(f, opts)
		if not ok then
			error(("generate %s\n%s"):format(k, v))
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
		error(("filetype %s\n%s"):format(opts.filetype, ft_subs))
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
