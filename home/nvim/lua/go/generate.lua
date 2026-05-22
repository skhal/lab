-- Copyright 2026 Samvel Khalatyan. All rights reserved.
--
-- Use of this source code is governed by a BSD-style
-- license that can be found in the LICENSE file.

-- generate runs go-generate for current buffer.
local generate = {
	-- Run executes go-generate command on current buffer. It reports an error if
	-- the tool fails.
	Run = function()
		local f = vim.fn.expand("%")
		local cmd = { "go", "generate", f }
		local obj = vim.system(cmd, { text = true }):wait()
		if obj.code ~= 0 then
			error(("go generate %s error: %s"):format(f, obj.stderr))
			return
		end
		vim.notify(("go generate %s done"):format(f))
	end,
}

return generate
