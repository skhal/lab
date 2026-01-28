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
require("lab-skeleton").setup({})
