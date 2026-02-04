-- Copyright 2026 Samvel Khalatyan. All rights reserved.
--
-- Use of this source code is governed by a BSD-style
-- license that can be found in the LICENSE file.

local lab_skeleton_path = vim.fs.joinpath(vim.fn.stdpath("data"), "lab-skeleton")
if not vim.uv.fs_stat(lab_skeleton_path) then
	vim.api.nvim_echo({
		{ ("Plugin lab-skeleton is not installed.\n%s"):format(lab_skeleton_path), "ErrorMsg" },
	}, true, {})
	return {}
end

vim.opt.rtp:append(lab_skeleton_path)

local git = require("lab-git")
local skel = require("lab-skeleton")

skel.setup({
	subs = {
		year = function(_)
			return os.date("%Y")
		end,

		holder = function(_)
			return git.config_get("user.name")
		end,
	},
})

local registrations = {
	-- keep-sorted start block=yes
	c = {
		find = function(file, _)
			local f = "new.c"
			if file:find("main%.c$") ~= nil then
				f = "new_main.c"
			end
			return f
		end,
		subs = {
			header = function(opts)
				local relpath = git.relpath(opts.file)
				relpath = relpath:gsub("_test%.c$", ".c")
				relpath = relpath:gsub("%.c$", ".h")
				return vim.fn.escape(relpath, "/")
			end,
		},
	},
	cpp = {
		find = function(file, _)
			local f = "new.cc"
			if file:find("main%.cc$") ~= nil then
				f = "new_main.cc"
			elseif file:find("_test%.cc$") ~= nil then
				f = "new_test.cc"
			elseif file:find("%.h$") ~= nil then
				f = "new.h"
			end
			return f
		end,
		subs = {
			guard = function(opts)
				if not opts.file:find("%.h$") then
					return
				end
				local relpath = git.relpath(opts.file)
				local guard = relpath:gsub("^/", "") -- if outside of a git worktree
				guard = guard:gsub("[/%.]", "_") .. "_"
				return guard:upper()
			end,
			header = function(opts)
				local relpath = git.relpath(opts.file)
				relpath = relpath:gsub("_test%.cc$", ".cc")
				relpath = relpath:gsub("%.cc$", ".h")
				return vim.fn.escape(relpath, "/")
			end,
			namespace = function(opts)
				local relpath = git.relpath(opts.file)
				local ns = vim.fs.dirname(relpath)
				ns = ns:gsub("^/", "") -- if outside of a git worktree
				return ns:gsub("/", "::")
			end,
		},
	},
	go = {
		find = function(file, _)
			local f = "new.go"
			if file:find("_test%.go$") ~= nil then
				f = "new_test.go"
			end
			return f
		end,
		subs = {
			package = function(opts)
				if opts.file:find("main.go$") ~= nil then
					return "main"
				end
				local abspath = vim.fs.abspath(opts.file)
				local dirname = vim.fs.dirname(abspath)
				local pkg = vim.fs.basename(dirname)
				return pkg
			end,
		},
	},
	proto = {
		subs = {
			edition = function(_)
				return 2024
			end,
			go_package = function(opts)
				local go_module = function()
					local cmd = { "go", "list", "-m" }
					local obj = vim.system(cmd, { text = true }):wait()
					if obj.code ~= 0 then
						error("go: can't get module name")
					end
					return (obj.stdout):gsub("%s+$", "")
				end
				local gomod = go_module()
				local relpath = git.relpath(opts.file)
				local dirname = vim.fs.dirname(relpath)
				local pkg = vim.fs.joinpath(gomod, dirname)
				return vim.fn.escape(pkg, "/")
			end,
			package = function(opts)
				local relpath = git.relpath(opts.file)
				local pkg = vim.fs.dirname(relpath)
				pkg = pkg:gsub("/", ".")
				return pkg:gsub("%.pb$", "") -- remove .pb suffix if any
			end,
		},
	},
	-- keep-sorted end
}

for ft, o in pairs(registrations) do
	skel.register(ft, o.find, o.subs)
end
