-- Copyright 2026 Samvel Khalatyan. All rights reserved.
--
-- Use of this source code is governed by a BSD-style
-- license that can be found in the LICENSE file.

local conform = require("conform")

local macos = vim.uv.os_uname().sysname:lower():find("darwin")
local clang_format_cmd = macos and "clang-format" or "clang-format21"

conform.setup({
	formatters = {
		-- keep-sorted start block=yes
		["clang-format"] = {
			command = clang_format_cmd,
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
		["yaml.ansible"] = { "keep-sorted" },
		bzl = { "buildifier", "keep-sorted" },
		c = { "clang-format", "keep-sorted" },
		cfg = { "keep-sorted" },
		conf = { "keep-sorted" },
		cpp = { "clang-format", "keep-sorted" },
		go = { "goimports", "keep-sorted" },
		html = { "superhtml", "keep-sorted" },
		json = { "deno" },
		lua = { "stylua", "keep-sorted" },
		markdown = { "markdownfmt", "keep-sorted" },
		pbtxt = { "txtpbfmt" },
		proto = { "clang-format", "keep-sorted" },
		python = { "black", "keep-sorted" },
		sh = { "shfmt", "keep-sorted" },
		svg = { "superhtml", "keep-sorted" },
		text = { "keep-sorted" },
		toml = { "tombi", "keep-sorted" },
		typescript = { "deno_fmt", "keep-sorted" },
		yaml = { "yamlfmt", "keep-sorted" },
		-- keep-sorted end
	},
})

vim.api.nvim_create_autocmd("LspAttach", {
	callback = function(_)
		vim.keymap.set("n", "<localleader>fc", conform.format, { buffer = true })
	end,
})

vim.bo.formatexpr = "v:lua.require'conform'.formatexpr()"
