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
		command = "0r " .. M.skel_path .. "/new.lua",
	})
end

return M
