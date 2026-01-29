-- Copyright 2026 Samvel Khalatyan. All rights reserved.
--
-- Use of this source code is governed by a BSD-style
-- license that can be found in the LICENSE file.

vim.bo.equalprg = "clang-format21 -assume-filename=%"

vim.api.nvim_buf_set_keymap(0, "n", "<localleader>rs", "<cmd>LspClangdSwitchSourceHeader<cr>", {})
