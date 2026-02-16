-- Copyright 2026 Samvel Khalatyan. All rights reserved.
--
-- Use of this source code is governed by a BSD-style
-- license that can be found in the LICENSE file.

local conform = require("conform")

conform.setup({
	formatters = {
		-- keep-sorted start block=yes
		["clang-format"] = {
			command = "clang-format21",
		},
		["markdownfmt"] = {
			append_args = { "-soft-wraps" },
		},
		["shfmt"] = {
			append_args = { "-bn", "-ci", "-p" },
		},
		["txtpbfmt"] = {
			meta = {
				url = "https://github.com/protocolbuffers/txtpbfmt",
				description = "txtpbfmt parses, edits and formats text proto files in a way that preserves comments.",
			},
			command = "txtpbfmt",
			args = { "-skip_all_colons", "-stdin_display_path", "$FILENAME" },
		},
		-- keep-sorted end
	},
	formatters_by_ft = {
		-- keep-sorted start
		bzl = { "buildifier" },
		c = { "clang-format" },
		cpp = { "clang-format" },
		go = { "goimports" },
		html = { "prettier" },
		json = { "prettier" },
		lua = { "stylua" },
		markdown = { "markdownfmt" },
		pbtxt = { "txtpbfmt" },
		proto = { "clang-format" },
		sh = { "shfmt" },
		yaml = { "yamlfmt" },
		-- keep-sorted end
	},
})

vim.api.nvim_create_autocmd("LspAttach", {
	callback = function(_)
		vim.keymap.set("n", "<localleader>fc", conform.format, { buffer = true })
	end,
})

vim.bo.formatexpr = "v:lua.require'conform'.formatexpr()"
