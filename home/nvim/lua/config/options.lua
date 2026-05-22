-- Copyright 2026 Samvel Khalatyan. All rights reserved.
--
-- Use of this source code is governed by a BSD-style
-- license that can be found in the LICENSE file.

local opt = vim.o

-- General
-- keep-sorted start
opt.mouse = "" -- disable to copy-on-select
opt.number = true
opt.termguicolors = true
-- keep-sorted end

-- Indentation
-- keep-sorted start
opt.expandtab = true
opt.shiftwidth = 2
opt.smartindent = true
opt.softtabstop = 2
opt.tabstop = 2
-- keep-sorted end

-- Invisibles
opt.listchars = "eol:¬,extends:›,precedes:‹,space:░,tab:«–»,trail:•"
vim.keymap.set({ "n" }, "<leader>l", "<esc>:set list!<cr>")

-- fix +q4D73 shown at Neovim startup.
-- Ref: https://github.com/neovim/neovim/issues/28776
local termfeatures = vim.g.termfeatures or {}
termfeatures.osc52 = false
vim.g.termfeatures = termfeatures
