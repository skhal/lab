-- Copyright 2026 Samvel Khalatyan. All rights reserved.
--
-- Use of this source code is governed by a BSD-style
-- license that can be found in the LICENSE file.

vim.o.winborder = "rounded"
if vim.fn.exists("+pumborder") == 1 then
	vim.o.pumborder = "rounded"
end

vim.cmd([[set completeopt+=menuone,preinsert,popup]])

vim.lsp.config("*", {
	-- ref: https://github.com/hrsh7th/cmp-nvim-lsp/blob/cbc7b02bb99fae35cb42f514762b89b5126651ef/lua/cmp_nvim_lsp/init.lua
	capabilities = {
		textDocument = {
			completion = {
				dynamicRegistration = false,
				completionItem = {
					snippetSupport = true,
					commitCharactersSupport = true,
					deprecatedSupport = true,
					preselectSupport = true,
					tagSupport = {
						valueSet = {
							1, -- Deprecated
						},
					},
					insertReplaceSupport = true,
					resolveSupport = {
						properties = {
							"documentation",
							"additionalTextEdits",
							"insertTextFormat",
							"insertTextMode",
							"command",
						},
					},
					insertTextModeSupport = {
						valueSet = {
							1, -- asIs
							2, -- adjustIndentation
						},
					},
					labelDetailsSupport = true,
				},
				contextSupport = true,
				insertTextMode = 1,
				completionList = {
					itemDefaults = {
						"commitCharacters",
						"editRange",
						"insertTextFormat",
						"insertTextMode",
						"data",
					},
				},
			},
		},
		workspace = {
			didChangeWatchedFiles = {
				dynamicRegistration = true,
			},
		},
	},
})

vim.api.nvim_create_autocmd("LspAttach", {
	callback = function(args)
		local client = assert(vim.lsp.get_client_by_id(args.data.client_id))
		vim.lsp.completion.enable(true, client.id, args.buf, { autotrigger = true })
	end,
})
