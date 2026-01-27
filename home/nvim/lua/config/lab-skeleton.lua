-- Copyright 2026 Samvel Khalatyan. All rights reserved.
--
-- Use of this source code is governed by a BSD-style
-- license that can be found in the LICENSE file.

local lab_skeleton_path = vim.fn.stdpath("data") .. "/lab-skeleton"
if not (vim.uv or vim.loop).fs_stat(lab_skeleton_path) then
	vim.api.nvim_echo({
		{ "Plugin lab-skeleton is not installed.\n", "ErrorMsg" },
		{ lab_skeleton_path, "ErrorMsg" },
	}, true, {})
	vim.api.finish()
	return {}
end

vim.opt.rtp:append(lab_skeleton_path)
require("lab-skeleton").setup({})
