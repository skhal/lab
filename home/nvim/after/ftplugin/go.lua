-- Copyright 2026 Samvel Khalatyan. All rights reserved.
--
-- Use of this source code is governed by a BSD-style
-- license that can be found in the LICENSE file.

-- keep-sorted start
vim.bo.equalprg = "goimports"
vim.bo.expandtab = false
-- keep-sorted end

local relatedFile = require("go.related_file")

vim.keymap.set("n", "<localleader>rd", relatedFile.Doc, { buffer = true })
vim.keymap.set("n", "<localleader>rs", relatedFile.Source, { buffer = true })
vim.keymap.set("n", "<localleader>rt", relatedFile.Test, { buffer = true })

local generate = require("go.generate")

vim.keymap.set("n", "<localleader>gen", generate.Run, { buffer = true })

local locationList = require("go.location_list")

vim.api.nvim_create_autocmd("LspAttach", {
	callback = function(_)
		vim.keymap.set("n", "gO", function()
			vim.lsp.buf.document_symbol({
				on_list = locationList.on_list,
				loclist = true,
			})
		end, { buffer = true })
	end,
})
