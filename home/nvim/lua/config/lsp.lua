-- Copyright 2026 Samvel Khalatyan. All rights reserved.
--
-- Use of this source code is governed by a BSD-style
-- license that can be found in the LICENSE file.

-- keep-sorted start
vim.lsp.enable("clangd")
vim.lsp.enable("gopls")
vim.lsp.enable("lua_ls")
vim.lsp.enable("protols")
vim.lsp.enable("typos_lsp")
-- keep-sorted end

vim.keymap.set("n", "gK", function()
	local new_config = not vim.diagnostic.config().virtual_lines
	vim.diagnostic.config({ virtual_lines = new_config })
end, { desc = "Toggle diagnostic virtual_lines" })
