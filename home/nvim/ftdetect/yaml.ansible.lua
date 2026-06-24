-- Copyright 2026 Samvel Khalatyan. All rights reserved.
--
-- Use of this source code is governed by a BSD-style
-- license that can be found in the LICENSE file.

local function detect_filetype(file)
	if not vim.fs.root(vim.fs.normalize(file), { "ansible.cfg" }) then
		return "yaml"
	end
	return "yaml.ansible"
end

vim.filetype.add({
	extension = {
		-- keep-sorted start
		["yaml.ansible"] = "yaml.ansible",
		yaml = detect_filetype,
		yml = detect_filetype,
		-- keep-sorted end
	},
})
