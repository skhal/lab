-- Copyright 2026 Samvel Khalatyan. All rights reserved.
--
-- Use of this source code is governed by a BSD-style
-- license that can be found in the LICENSE file.

local function run(path)
	local cmd = { "go", "generate", path }
	local obj = vim.system(cmd, { text = true }):wait()
	if obj.code ~= 0 then
		error(("go generate %s error: %s"):format(path, obj.stderr))
		return
	end
	vim.notify(("go generate %s done"):format(path))
end

-- generate runs go-generate for current buffer.
local generate = {
	-- RunFile executes go-generate command on current buffer. It reports an
	-- error if the tool fails.
	RunFile = function()
		local f = vim.fn.expand("%")
		run(f)
	end,
	-- RunPackage runs go-generate on the package holding the current buffer.
	-- It reports an error if tool fails.
	RunPackage = function()
		local f = vim.fn.expand("%")
		run(vim.fs.dirname(f))
	end,
}

return generate
