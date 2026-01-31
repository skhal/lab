-- Copyright 2026 Samvel Khalatyan. All rights reserved.
--
-- Use of this source code is governed by a BSD-style
-- license that can be found in the LICENSE file.

-- keep-sorted start
vim.lsp.enable("clangd")
vim.lsp.enable("gopls")
vim.lsp.enable("lua_ls")
vim.lsp.enable("protols") -- NOLINT
vim.lsp.enable("typos_lsp")
-- keep-sorted end

vim.cmd([[set completeopt+=menuone,noselect,popup]])

vim.api.nvim_create_autocmd("LspAttach", {
	callback = function(args)
		local client = vim.lsp.get_client_by_id(args.data.client_id)
		vim.lsp.completion.enable(true, client.id, args.buf, { autotrigger = true })
		vim.keymap.set("n", "gK", function()
			local new_config = not vim.diagnostic.config().virtual_lines
			vim.diagnostic.config({ virtual_lines = new_config })
		end, { desc = "Toggle diagnostic virtual_lines" })
		vim.keymap.set("n", "<localleader>csi", vim.lsp.buf.incoming_calls, { desc = "Incoming calls" })
		vim.keymap.set("n", "<localleader>cso", vim.lsp.buf.outgoing_calls, { desc = "Outgoing calls" })
	end,
})
