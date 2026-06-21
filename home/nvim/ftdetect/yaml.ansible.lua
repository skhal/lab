-- Copyright 2026 Samvel Khalatyan. All rights reserved.
--
-- Use of this source code is governed by a BSD-style
-- license that can be found in the LICENSE file.

---@param path string path to file
local function is_ansible(p)
	p = vim.fs.normalize(vim.fs.abspath(p))
	p = vim.fs.dirname(p)
	while p ~= "/" do
		local f = vim.fs.joinpath(p, "ansible.cfg")
		if vim.uv.fs_stat(f) then
			return true
		end
		p = vim.fs.dirname(p)
	end
	return false
end

vim.api.nvim_create_autocmd({ "BufRead", "BufNewFile" }, {
  pattern = { '*' },
  callback = function(ev)
    if vim.bo.filetype ~= "yaml" then
      return vim.bo.filetype
		end
		if not is_ansible(ev.file) then
			return vim.bo.filetype
		end
		vim.bo.filetype = "yaml.ansible"
  end,
})
