-- Copyright 2026 Samvel Khalatyan. All rights reserved.
--
-- Use of this source code is governed by a BSD-style
-- license that can be found in the LICENSE file.

-- keep-sorted start
vim.o.autoindent = true
vim.o.colorcolumn = "80"
vim.o.cursorline = true
vim.o.expandtab = true
vim.o.hlsearch = true
vim.o.incsearch = true
vim.o.listchars = "eol:¬,extends:›,precedes:‹,space:░,tab:«–»,trail:•"
vim.o.number = true
vim.o.shiftwidth = 2
vim.o.smartindent = true
vim.o.softtabstop = 2
vim.o.tabstop = 2
-- keep-sorted end

vim.keymap.set({ "n" }, "<leader>l", "<esc>:set list!<cr>")

vim.api.nvim_create_autocmd("Filetype", {
	group = vim.api.nvim_create_augroup("Indents", { clear = true }),
	pattern = { "go", "make" },
	command = "set noexpandtab",
})
