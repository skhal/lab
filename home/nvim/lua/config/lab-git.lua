-- Copyright 2026 Samvel Khalatyan. All rights reserved.
--
-- Use of this source code is governed by a BSD-style
-- license that can be found in the LICENSE file.

local lab_git_path = vim.fs.joinpath(vim.fn.stdpath("data"), "lab-git")
if not vim.uv.fs_stat(lab_git_path) then
	vim.api.nvim_echo({
		{ ("Plugin lab-git is not installed.\n%s"):format(lab_git_path), "ErrorMsg" },
	}, true, {})
	return {}
end

vim.opt.rtp:append(lab_git_path)
require("lab-git")
