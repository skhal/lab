-- Copyright 2026 Samvel Khalatyan. All rights reserved.
--
-- Use of this source code is governed by a BSD-style
-- license that can be found in the LICENSE file.

-- keep-sorted start
vim.o.colorcolumn = "80"
vim.o.cursorline = true
vim.o.number = true
-- keep-sorted end

vim.o.listchars = "eol:¬,extends:›,precedes:‹,space:░,tab:«–»,trail:•"
vim.keymap.set({ "n" }, "<leader>l", "<esc>:set list!<cr>")

-- keep-sorted start
vim.o.expandtab = true
vim.o.shiftwidth = 2
vim.o.smartindent = true
vim.o.softtabstop = 2
vim.o.tabstop = 2
-- keep-sorted end
vim.api.nvim_create_autocmd("Filetype", {
	group = vim.api.nvim_create_augroup("Indents", { clear = true }),
	pattern = { "go", "make" },
	command = "set noexpandtab",
})

-- Plugin manager
require("config.lazy")
