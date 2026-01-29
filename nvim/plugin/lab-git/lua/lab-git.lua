-- Copyright 2026 Samvel Khalatyan. All rights reserved.
--
-- Use of this source code is governed by a BSD-style
-- license that can be found in the LICENSE file.

local M = {
	cache = {},
}

function M.config_get(key)
	local cache_key = ("config_get.%s"):format(key)
	if M.cache[cache_key] ~= nil then
		return M.cache[cache_key]
	end
	local cmd = { "git", "config", "--get", key }
	local obj = vim.system(cmd, { text = true }):wait()
	if obj.code ~= 0 then
		error(("git-config: can't get %s"):format(key))
	end
	local val = (obj.stdout):gsub("%s+$", "")
	M.cache[cache_key] = val
	return val
end

function M.worktree_path()
	local cache_key = "worktree_path"
	if M.cache[cache_key] ~= nil then
		return M.cache[cache_key]
	end
	local cmd = { "git", "rev-parse", "--show-toplevel" }
	local obj = vim.system(cmd, { text = true }):wait()
	if obj.code ~= 0 then
		error("git-rev-parse: can't get worktree path")
	end
	local val = (obj.stdout):gsub("%s+$", "")
	M.cache[cache_key] = val
	return val
end

function M.relpath(file)
	local worktree = M.worktree_path()
	local abspath = vim.fs.abspath(file)
	local relpath = vim.fs.relpath(worktree, abspath)
	return relpath or file -- file if it is outside of the worktree
end

return M
