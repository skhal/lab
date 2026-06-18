-- Copyright 2026 Samvel Khalatyan. All rights reserved.
--
-- Use of this source code is governed by a BSD-style
-- license that can be found in the LICENSE file.

local macos = vim.uv.os_uname().sysname:lower():find("darwin")
local clang_cmd = macos and "/opt/homebrew/opt/llvm@21/bin/clangd" or "clangd21"

return {
	cmd = { clang_cmd },
}
