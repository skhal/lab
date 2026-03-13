-- Copyright 2026 Samvel Khalatyan. All rights reserved.
--
-- Use of this source code is governed by a BSD-style
-- license that can be found in the LICENSE file.

-- keep-sorted start
vim.bo.equalprg = "goimports"
vim.bo.expandtab = false
-- keep-sorted end

local function select_and_edit(file, files)
	vim.ui.select(files, {
		prompt = "Select file from " .. vim.fs.dirname(file) .. ":",
		format_item = function(item)
			return vim.fs.basename(item)
		end,
	}, function(choice)
		vim.cmd.edit(choice)
	end)
end

local function select_source(file)
	local opt_nosuf = false
	local opt_list = true
	local files = vim.fn.globpath(vim.fs.dirname(file), "*.go", opt_nosuf, opt_list)
	files = vim.tbl_filter(function(item)
		if vim.fs.basename(file) == vim.fs.basename(item) then
			return false
		end
		if item:find("doc%.go$") then
			return false
		end
		if item:find("_test%.go$") then
			return false
		end
		return true
	end, files)
	if not next(files) then
		vim.api.nvim_echo({
			{ "no source files found" },
		}, true, {})
		return
	end
	select_and_edit(file, files)
end

local function select_test(file)
	local opt_nosuf = false
	local opt_list = true
	local files = vim.fn.globpath(vim.fs.dirname(file), "*_test.go", opt_nosuf, opt_list)
	if not next(files) then
		vim.api.nvim_echo({
			{ "no test files found" },
		}, true, {})
		return
	end
	select_and_edit(file, files)
end

local RelatedFile = {
	doc = function()
		local f = vim.fn.expand("%")
		if f:find("doc.go$") then
			return
		end
		f = vim.fs.joinpath(vim.fs.dirname(f), "doc.go")
		vim.cmd.edit(f)
	end,
	source = function()
		local f = vim.fn.expand("%")
		if f:find("doc%.go$") then
			select_source(f)
			return
		end
		if not f:find("_test%.go$") then
			select_source(f)
			return
		end
		f = f:gsub("_test%.go$", ".go")
		vim.cmd.edit(f)
	end,
	test = function()
		local f = vim.fn.expand("%")
		if f:find("doc%.go$") then
			select_test(f)
			return
		end
		if f:find("_test%.go$") then
			return
		end
		f = f:gsub("%.go$", "_test.go")
		vim.cmd.edit(f)
	end,
}

vim.keymap.set("n", "<localleader>rd", RelatedFile.doc, { buffer = true })
vim.keymap.set("n", "<localleader>rs", RelatedFile.source, { buffer = true })
vim.keymap.set("n", "<localleader>rt", RelatedFile.test, { buffer = true })

-- receiverType extracts structure name from the identifier. It is the
-- receiver type in case of a method identifier, or the identifier.
local function receiverType(ident)
	local name = ident:match("^%(%*?(%w+)%)")
	if name ~= nil then
		return name
	end
	error(("failed to get the receiver type: %s"):format(ident))
end

-- LocationTree is a hierarchy of location items with structure fields and
-- methods stored in a structure block. It helps group methods by structure
-- type in method declaration order, see [LocationTree.Add].
local LocationTree = {}

function LocationTree:new()
	local o = {
		items = {}, -- top-level location items except fields and methods
		structs = {}, -- a map of struct type to a table with fields and methods
		lastStructName = nil, -- last structure type for fields
		-- pending is a list of structure types that had a method declaration but
		-- no structure definition yet. The following code is a valid Go code
		-- with method declared before the receiver type:
		--
		--	func (Foo) MethodA() {} // Foo was not declared yet
		--											    // pending={Foo={methods={{text=MethodA}}}}
		--	...
		--	type Foo struct {}
		pending = {},
	}
	setmetatable(o, self)
	self.__index = self
	return o
end

-- Add introduces an item into the LocationTree. It stores field and method
-- items under structures, in the order of declaration. Everything else goes
-- into the items list.
function LocationTree:Add(item)
	if item.kind == "Struct" then
		self:addStruct(item)
	elseif item.kind == "Field" then
		self:addField(item)
		return
	elseif item.kind == "Method" then
		self:addMethod(item)
		return
	end
	table.insert(self.items, item)
end

-- addStruct registers a structure to receive fields and methods. Keep in mind
-- that the structure may already exist if method declarations precede type
-- declaration.
function LocationTree:addStruct(item)
	if self.structs[item.ident] then
		self.pending[item.ident] = nil
	else
		self:registerStruct(item.ident)
	end
	self.lastStructName = item.ident
end

function LocationTree:registerStruct(name)
	local struct = {
		fields = {},
		methods = {},
	}
	self.structs[name] = struct
	return struct
end

-- addField attaches a field to the structure.
function LocationTree:addField(item)
	table.insert(self.structs[self.lastStructName].fields, item)
end

-- addMethod attaches a method to the receiver type. Keep in mind that a method
-- can be declared before the receiver. Is so, addMethod marks the structure
-- "pending".
function LocationTree:addMethod(item)
	local name = receiverType(item.ident)
	local struct = self.structs[name]
	if not struct then
		struct = self:registerStruct(name)
		self.pending[name] = true
	end
	table.insert(struct.methods, item)
end

-- Items flattens the items list with every structure expanded with field and
-- then methods. Any pending structures, go to the end of the location list.
function LocationTree:Items()
	local items = {}
	for _, v in ipairs(self.items) do
		if v.kind == "Struct" then
			table.insert(items, v)
			local struct = self.structs[v.ident]
			for _, f in ipairs(struct.fields) do
				table.insert(items, f)
			end
			for _, m in ipairs(struct.methods) do
				table.insert(items, m)
			end
		else
			table.insert(items, v)
		end
	end
	for _, name in ipairs(self.pending) do
		local struct = self.structs[name]
		for _, m in ipairs(struct.methods) do
			table.insert(items, m)
		end
	end
	return items
end

local LocationList = {
	indent = {
		kind = {
			-- keep-sorted start
			Field = true,
			Method = true,
			-- keep-sorted end
		},
	},
	renames = {
		kind = {
			-- keep-sorted start
			Constant = "D",
			Field = "F",
			Function = "F",
			Method = "M",
			Struct = "S",
			Variable = "V",
			-- keep-sorted end
		},
	},
}

function LocationList.on_list(opts)
	opts.title = "Document symbols"
	opts.items = LocationList.group(opts.items)
	opts.quickfixtextfunc = LocationList.quickfixtextfunc
	vim.fn.setloclist(0, {}, " ", opts)
	vim.cmd.lopen()
end

function LocationList.group(items)
	local lt = LocationTree:new()
	for _, item in ipairs(items) do
		item.ident = item.text:gsub("^%[%w+%] ", "")
		lt:Add(item)
	end
	return vim.iter(lt:Items())
		:map(function(o)
			return LocationList.rename(o)
		end)
		:totable()
end

-- rename replaces the "[Kind]" prefix in the location list items text with
-- abbreviated kind letter from [LocationList.renames.kind].
function LocationList.rename(opt)
	local kind = LocationList.renames.kind[opt.kind] or opt.kind
	local name = opt.ident
	if opt.kind == "Method" then
		name = name:gsub("^%(%*?(%w+)%)%.", "")
	end
	opt.text = ("%s %s"):format(kind, name)
	if LocationList.indent.kind[opt.kind] then
		opt.text = " " .. opt.text
	end
	return opt
end

-- quickfixtextfunc shows only item text in the location list window.
function LocationList.quickfixtextfunc(opts)
	local formatted = {}
	local items = vim.fn.getloclist(opts.winid, { id = opts.id, items = 0 }).items
	for i = opts.start_idx, opts.end_idx do
		local v = items[i]
		table.insert(formatted, v.text)
	end
	return formatted
end

vim.api.nvim_create_autocmd("LspAttach", {
	callback = function(_)
		vim.keymap.set("n", "gO", function()
			vim.lsp.buf.document_symbol({
				on_list = LocationList.on_list,
				loclist = true,
			})
		end, { buffer = true })
	end,
})
