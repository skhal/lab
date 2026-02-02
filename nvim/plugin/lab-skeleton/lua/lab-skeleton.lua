-- Copyright 2026 Samvel Khalatyan. All rights reserved.
--
-- Use of this source code is governed by a BSD-style
-- license that can be found in the LICENSE file.

local M = {
	path = vim.fs.joinpath(vim.fn.stdpath("data"), "lab-skeleton", "skel"),
	find = {
		default = function(_, ft)
			return "new." .. (ft or "default")
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
	local load_skeleton = function(skel, subs)
		local read_file = function(path, f)
			vim.uv.fs_open(path, "r", tonumber("444", 8), function(err, fd)
				assert(not err, err)
				---@diagnostic disable-next-line:redefined-local
				vim.uv.fs_fstat(fd, function(err, stat) -- luacheck: ignore err
					assert(not err, err)
					---@diagnostic disable-next-line:redefined-local
					vim.uv.fs_read(fd, stat.size, 0, function(err, data) -- luacheck: ignore err
						assert(not err, err)
						---@diagnostic disable-next-line:redefined-local
						vim.uv.fs_close(fd, function(err) -- luacheck: ignore err
							assert(not err, err)
							return f(data)
						end)
					end)
				end)
			end)
		end
		read_file(
			skel.path,
			vim.schedule_wrap(function(data)
				for key, val in pairs(subs) do
					data = data:gsub(("{{%s}}"):format(key), val)
				end
				data = data:gsub("\n$", "") -- remove trailing EOL
				local lines = {}
				local pos = {}
				for row, line in pairs(vim.split(data, "\n")) do
					local col, _ = line:find("{{cursor}}")
					if col ~= nil then
						pos = { row, col - 1 }
						line = line:gsub("{{cursor}}", "")
					end
					table.insert(lines, line)
				end
				vim.api.nvim_buf_set_lines(0, 0, -1, false, lines)
				if next(pos) then
					vim.api.nvim_win_set_cursor(0, pos)
				end
			end)
		)
	end
	local load = function(e)
		local skel = M.find_skeleton(e.file)
		local subs = M.gen_substitutes({ file = e.file, filetype = skel.filetype })
		load_skeleton(skel, subs)
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
	local stat = vim.uv.fs_stat(path)
	if not stat then
		error(("skeleton %s: does not exist."):format(name))
	end
	return {
		filetype = ft,
		path = path,
		stat = stat,
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
