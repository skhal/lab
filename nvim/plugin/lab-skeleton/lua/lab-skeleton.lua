-- Copyright 2026 Samvel Khalatyan. All rights reserved.
--
-- Use of this source code is governed by a BSD-style
-- license that can be found in the LICENSE file.

local M = {
	path = vim.fs.joinpath(vim.fn.stdpath("data"), "lab-skeleton", "skel"),
	find = {
		default = function(_, ft)
			return "new." .. ft
		end,
	},
	subs = {
		default = {},
	},
	augroup = vim.api.nvim_create_augroup("LabSkeleton", { clear = true }),
}

local function table_merge(dst, src)
	for k, v in pairs(src) do
		dst[k] = v
	end
end

function M.setup(opts)
	if opts.path ~= nil then
		M.path = opts.path
	end
	if opts.find ~= nil then
		M.find.default = opts.find
	end
	if opts.subs ~= nil then
		table_merge(M.subs.default, opts.subs)
	end
	vim.api.nvim_create_autocmd("BufNewFile", {
		group = M.augroup,
		desc = "Load template",
		pattern = "*",
		callback = M.load,
	})
end

function M.register(ft, find, subs)
	if ft == "default" then
		vim.api.nvim_echo({
			{ "Use setup() to change defaults.", "ErrorMsg" },
		}, true, {})
	end
	if find ~= nil then
		M.find[ft] = find
	end
	if subs ~= nil then
		if not M.subs[ft] then
			M.subs[ft] = subs
		else
			table_merge(M.subs[ft], subs)
		end
	end
end

function M.load(ev)
	local report = function(skel, subs)
		local msgs = {
			{ ("skel: %s\n"):format(skel.path), "Normal" },
			{ ("ft: %s\n"):format(skel.filetype), "Normal" },
		}
		if next(subs) then
			for k, v in pairs(subs) do
				table.insert(msgs, { (". %s: %s\n"):format(k, v), "Normal" })
			end
		end
		vim.api.nvim_echo(msgs, true, {})
	end
	local load_skeleton = function(file)
		vim.cmd("0r " .. file)
	end
	local run_substitutes = function(subs)
		for key, val in pairs(subs) do
			vim.cmd("silent! %s/{{" .. key .. "}}/" .. val)
		end
	end
	local position_cursor = function()
		for line_num, line in ipairs(vim.api.nvim_buf_get_lines(0, 0, -1, false)) do
			local from, _ = line:find("{{cursor}}")
			if from ~= nil then
				vim.cmd("silent! %s/{{cursor}}//g")
				vim.api.nvim_win_set_cursor(0, { line_num, from - 1 })
				break
			end
		end
	end
	local load = function(e)
		local skel = M.find_skeleton(e.file)
		local subs = M.gen_substitutes({ file = e.file, filetype = skel.filetype })
		load_skeleton(skel.path)
		run_substitutes(subs)
		position_cursor()
		report(skel, subs)
	end
	local ok, err = pcall(load, ev)
	if not ok then
		vim.api.nvim_echo({ { err, "ErrorMsg" } }, true, {})
		return
	end
end

function M.find_skeleton(file)
	local ft = vim.filetype.match({ filename = file })
	local find = M.find[ft] or M.find.default
	local name = find(file, ft)
	local path = vim.fs.joinpath(M.path, name)
	if not (vim.uv or vim.loop).fs_stat(path) then
		error(("skeleton %s: does not exist."):format(name))
	end
	return {
		filetype = ft,
		path = path,
	}
end

function M.gen_substitutes(opts)
	local gen_substitutes = function(gens, o)
		local subs = {}
		for k, f in pairs(gens) do
			local ok, v = pcall(f, o)
			if not ok then
				error(("generate %s\n%s"):format(k, v))
			end
			if v ~= nil then
				subs[k] = v
			end
		end
		return subs
	end
	local ok, subs = pcall(gen_substitutes, M.subs.default or {}, opts)
	if not ok then
		error(("common substitutes\n%s"):format(subs))
	end
	local ftsubs
	ok, ftsubs = pcall(gen_substitutes, M.subs[opts.filetype] or {}, opts)
	if not ok then
		error(("filetype %s\n%s"):format(opts.filetype, ftsubs))
	end
	if next(ftsubs) then
		if next(subs) then
			table_merge(subs, ftsubs)
		else
			subs = ftsubs
		end
	end
	return subs or {}
end

return M
