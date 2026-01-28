-- Copyright 2026 Samvel Khalatyan. All rights reserved.
--
-- Use of this source code is governed by a BSD-style
-- license that can be found in the LICENSE file.

local M = {}

function M.config_get()
	if M.config_username_ ~= nil then
		return M.config_username_
	end
	local cmd = { "git", "config", "--get", "user.name" }
	local obj = vim.system(cmd, { text = true }):wait()
	if obj.code ~= 0 then
		error("git: can't get user.name")
	end
	M.config_username_ = (obj.stdout):gsub("%s+$", "")
	return M.config_username_
end

function M.worktree_path()
	if M.worktree_path_ ~= nil then
		return M.worktree_path_
	end
	local cmd = { "git", "rev-parse", "--show-toplevel" }
	local obj = vim.system(cmd, { text = true }):wait()
	if obj.code ~= 0 then
		error("git: can't get worktree path")
	end
	M.worktree_path_ = (obj.stdout):gsub("%s+$", "")
	return M.worktree_path_
end

function M.relpath(file)
	local worktree = M.wortree_path()
	local abspath = vim.fs.abspath(file)
	local relpath = vim.fs.relpath(worktree, abspath)
	return relpath or file -- file if it is outside of the worktree
end

return M
