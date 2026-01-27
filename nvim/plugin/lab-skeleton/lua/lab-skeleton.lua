-- Copyright 2026 Samvel Khalatyan. All rights reserved.
--
-- Use of this source code is governed by a BSD-style
-- license that can be found in the LICENSE file.

local Skeleton = vim.api.nvim_create_augroup("LabSkeleton", { clear = true })

vim.api.nvim_create_autocmd("BufNewFile", {
	group = Skeleton,
	pattern = { "*.lua" },
	command = "0r ~/.local/share/nvim/lab-skeleton/skel/new.lua",
})
