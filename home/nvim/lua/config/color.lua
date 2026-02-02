-- Copyright 2026 Samvel Khalatyan. All rights reserved.
--
-- Use of this source code is governed by a BSD-style
-- license that can be found in the LICENSE file.

vim.o.colorcolumn = "80"

for _, group in pairs({ "Normal", "NormalFloat" }) do
	vim.api.nvim_set_hl(0, group, { bg = "none" })
end
