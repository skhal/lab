-- Copyright 2026 Samvel Khalatyan. All rights reserved.
--
-- Use of this source code is governed by a BSD-style
-- license that can be found in the LICENSE file.

-- keep-sorted start
vim.bo.equalprg = "goimports"
vim.bo.expandtab = false
-- keep-sorted end

local function select_and_edit(file, files)
	if #files == 1 then
		vim.cmd.edit(files[1])
		return
	end
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
		if item:find("_string%.go$") then
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

-- Go automates execution of different Go commands such as generate.
local Go = {
	-- generate runs `go generate` on current buffer. It reports and error if
	-- go-generate fails.
	generate = function()
		local f = vim.fn.expand("%")
		local cmd = { "go", "generate", f }
		local obj = vim.system(cmd, { text = true }):wait()
		if obj.code ~= 0 then
			error(("go generate %s: failed: %s"):format(f, obj.stderr))
			return
		end
		vim.notify(("go generate %s: done"):format(f))
	end,
}

vim.keymap.set("n", "<localleader>gen", Go.generate, { buffer = true })

-- LocationTree is a hierarchy of location items with structure fields and
-- methods stored in a structure block. It helps group methods by structure
-- type in method declaration order, see [LocationTree.Add].
local LocationTree = {}

function LocationTree:new()
	local o = {
		-- items is a list of top-level location items excluding fields and methods.
		items = {},

		-- structs is a map of struct type name to the fields and methods.
		-- LocationTree uses it to group methods by receiver.
		structs = {},

		-- interfaces is a map of interface type name to the methods. LocationTree
		-- uses it to group interface methods by interface.
		interfaces = {},

		-- last structure or interface type name.
		lastName = nil,

		-- pending is a list of structure types that had a method declaration but
		-- no structure definition yet. The following code is a valid Go code
		-- with method declared before the receiver type:
		--
		--	func (Foo) MethodA() {} // Foo was not declared yet
		--											    // pending={Foo={methods={{text=MethodA}}}}
		--	...
		--	type Foo struct {}
		pending = {},

		adders = {
			default = LocationTree.addItem,
			kind = {
				-- keep-sorted start
				Class = LocationTree.addStruct,
				Field = LocationTree.addField,
				Interface = LocationTree.addInterface,
				Method = LocationTree.addMethod,
				Struct = LocationTree.addStruct,
				-- keep-sorted end
			},
		},
	}
	setmetatable(o, self)
	self.__index = self
	return o
end

-- Add introduces an item into the LocationTree. It stores field and method
-- items under structures, in the order of declaration. Everything else goes
-- into the items list.
function LocationTree:Add(item)
	local add = self.adders.kind[item.kind] or self.adders.default
	add(self, item)
end

function LocationTree:addItem(item)
	table.insert(self.items, item)
end

-- addStruct registers a structure to receive fields and methods. Keep in mind
-- that the structure may already exist if method declarations precede type
-- declaration.
function LocationTree:addStruct(item)
	table.insert(self.items, item)
	if self.structs[item.ident] then
		self.pending[item.ident] = nil
	else
		self:registerStruct(item.ident)
	end
	self.lastName = item.ident
end

-- addInterface registers an interface to receive methods. It stores the
-- interface type name in [self.lastName] to add interface methods, which do
-- not have a receiver and come right after the interface type declaration.
function LocationTree:addInterface(item)
	table.insert(self.items, item)
	self.interfaces[item.ident] = {
		methods = {},
	}
	self.lastName = item.ident
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
	table.insert(self.structs[self.lastName].fields, item)
end

-- addMethod attaches a method to the receiver type. Keep in mind that a method
-- can be declared before the receiver. Is so, addMethod marks the structure
-- "pending".
function LocationTree:addMethod(item)
	local name = item.ident:match("^%(%*?(%w+)%)")
	if name ~= nil then
		self:addStructMethod(item, name)
	else
		self:addInterfaceMethod(item)
	end
end

-- addStructMethod adds a method item to the structure type name.
function LocationTree:addStructMethod(item, name)
	local struct = self.structs[name]
	if not struct then
		struct = self:registerStruct(name)
		self.pending[name] = true
	end
	table.insert(struct.methods, item)
end

-- addInterfaceMethod adds a method to the the interface type [self.lastName].
function LocationTree:addInterfaceMethod(item)
	local interface = self.interfaces[self.lastName]
	table.insert(interface.methods, item)
end

local LocationItemsCollector = {}

function LocationItemsCollector:new(lt)
	local o = {
		loctree = lt,
		items = {},

		collectors = {
			default = LocationItemsCollector.collectItem,
			kind = {
				-- keep-sorted start
				Class = LocationItemsCollector.collectStruct,
				Interface = LocationItemsCollector.collectInterface,
				Struct = LocationItemsCollector.collectStruct,
				-- keep-sorted end
			},
		},
	}
	setmetatable(o, self)
	self.__index = self
	return o
end

function LocationItemsCollector:Collect()
	self.items = {}
	for _, item in ipairs(self.loctree.items) do
		local c = self.collectors.kind[item.kind] or self.collectors.default
		c(self, item)
	end
	self:collectPending()
	return self.items
end

function LocationItemsCollector:collectItem(item)
	table.insert(self.items, item)
end

function LocationItemsCollector:collectStruct(item)
	table.insert(self.items, item)
	local struct = self.loctree.structs[item.ident]
	for _, f in ipairs(struct.fields) do
		table.insert(self.items, f)
	end
	for _, m in ipairs(struct.methods) do
		table.insert(self.items, m)
	end
end

function LocationItemsCollector:collectInterface(item)
	table.insert(self.items, item)
	local interface = self.loctree.interfaces[item.ident]
	for _, m in ipairs(interface.methods) do
		table.insert(self.items, m)
	end
end

function LocationItemsCollector:collectPending()
	for _, name in ipairs(self.loctree.pending) do
		local struct = self.loctree.structs[name]
		for _, m in ipairs(struct.methods) do
			table.insert(self.items, m)
		end
	end
end

-- Items flattens the items list with every structure expanded with field and
-- then methods. Any pending structures, go to the end of the location list.
function LocationTree:Items()
	local collector = LocationItemsCollector:new(self)
	return collector:Collect()
end

local LocationList = {
	namespace = vim.api.nvim_create_namespace("go-loclist"),
	indentation = "  ",
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
			Class = "T",
			Constant = "D",
			Field = "F",
			Function = "F",
			Interface = "I",
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
	vim.fn.setloclist(opts.winid, {}, " ", opts)
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
	if not name:match("^%u") then
		kind = kind:lower()
	end
	opt.text = ("%s %s"):format(kind, name)
	if LocationList.indent.kind[opt.kind] then
		opt.text = LocationList.indentation .. opt.text
	end
	return opt
end

-- quickfixtextfunc shows only item text in the location list window.
function LocationList.quickfixtextfunc(opts)
	local loclist = vim.fn.getloclist(opts.winid, { id = opts.id, items = 0, qfbufnr = 1 })
	if opts.start_idx == 1 then
		vim.api.nvim_buf_clear_namespace(loclist.qfbufnr, LocationList.namespace, 0, -1)
	end
	local formatted = {}
	for i = opts.start_idx, opts.end_idx do
		local item = loclist.items[i]
		table.insert(formatted, item.text)
	end
	vim.schedule(function()
		LocationList.highlight(loclist.qfbufnr)
	end)
	return formatted
end

function LocationList.highlight(bufnr)
	vim.api.nvim_buf_call(bufnr, function()
		vim.cmd([[
			syn clear

			syn match		qfTop			/^\w \l\w*$/
			syn match		qfTopExp	/^\w \u\w*$/

			syn match		qfSub			/^\s\+\w \l\w*$/
			syn match		qfSubExp	/^\s\+\w \u\w*$/

			hi def link qfUnexported NonText
			hi def link qfExported Normal

			hi def link qfTop	qfUnexported
			hi def link qfSub	qfUnexported

			hi def link qfTopExp qfExported
			hi def link qfSubExp qfExported
		]])
	end)
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
