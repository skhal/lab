-- Copyright 2026 Samvel Khalatyan. All rights reserved.
--
-- Use of this source code is governed by a BSD-style
-- license that can be found in the LICENSE file.

-- claude: a warm dark colorscheme built around the Claude palette: ivory text
-- on a charcoal background, with a clay-orange accent.

vim.cmd("highlight clear")
if vim.fn.exists("syntax_on") == 1 then
	vim.cmd("syntax reset")
end

vim.o.background = "dark"
vim.g.colors_name = "claude"

local c = {
	-- keep-sorted start
	amber = "#e0b265",
	bg = "#262624",
	bg_alt = "#2e2e2b",
	bg_float = "#30302e",
	bg_select = "#434340",
	blue = "#6a9bcc",
	border = "#4a4945",
	clay = "#d97757",
	comment = "#87867f",
	cyan = "#7fb5b0",
	fg = "#e8e6dc",
	fg_bright = "#faf9f5",
	fg_dim = "#b0aea5",
	green = "#8fa86f",
	purple = "#b49fcc",
	red = "#d4605a",
	-- keep-sorted end
}

local function set(groups)
	for name, spec in pairs(groups) do
		vim.api.nvim_set_hl(0, name, spec)
	end
end

-- editor
set({
	-- keep-sorted start
	ColorColumn = { bg = c.bg_alt },
	Conceal = { fg = c.comment },
	CurSearch = { fg = c.bg, bg = c.clay, bold = true },
	Cursor = { fg = c.bg, bg = c.fg_bright },
	CursorColumn = { bg = c.bg_alt },
	CursorLine = { bg = c.bg_alt },
	CursorLineNr = { fg = c.clay, bold = true },
	Directory = { fg = c.blue },
	EndOfBuffer = { fg = c.bg },
	ErrorMsg = { fg = c.red, bold = true },
	FloatBorder = { fg = c.border, bg = "none" },
	FoldColumn = { fg = c.comment },
	Folded = { fg = c.fg_dim, bg = c.bg_alt },
	IncSearch = { fg = c.bg, bg = c.clay },
	LineNr = { fg = c.comment },
	MatchParen = { fg = c.clay, bold = true, underline = true },
	ModeMsg = { fg = c.fg_dim, bold = true },
	MoreMsg = { fg = c.green },
	NonText = { fg = c.border },
	Normal = { fg = c.fg, bg = "none" },
	NormalFloat = { fg = c.fg, bg = "none" },
	Pmenu = { fg = c.fg, bg = c.bg_float },
	PmenuSbar = { bg = c.bg_select },
	PmenuSel = { fg = c.fg_bright, bg = c.bg_select, bold = true },
	PmenuThumb = { bg = c.comment },
	Question = { fg = c.green },
	Search = { fg = c.bg, bg = c.amber },
	SignColumn = { fg = c.comment, bg = "none" },
	SpecialKey = { fg = c.border },
	StatusLine = { fg = c.fg, bg = c.bg_select },
	StatusLineNC = { fg = c.comment, bg = c.bg_alt },
	TabLine = { fg = c.comment, bg = c.bg_alt },
	TabLineFill = { bg = c.bg_alt },
	TabLineSel = { fg = c.fg_bright, bg = c.bg, bold = true },
	Title = { fg = c.clay, bold = true },
	Visual = { bg = c.bg_select },
	WarningMsg = { fg = c.amber },
	WildMenu = { fg = c.fg_bright, bg = c.bg_select },
	WinSeparator = { fg = c.border },
	-- keep-sorted end
})

-- syntax
set({
	-- keep-sorted start
	Boolean = { fg = c.clay },
	Character = { fg = c.green },
	Comment = { fg = c.comment, italic = true },
	Conditional = { fg = c.clay },
	Constant = { fg = c.amber },
	Debug = { fg = c.red },
	Define = { fg = c.purple },
	Delimiter = { fg = c.fg_dim },
	Error = { fg = c.red },
	Exception = { fg = c.clay },
	Float = { fg = c.amber },
	Function = { fg = c.blue },
	Identifier = { fg = c.fg },
	Include = { fg = c.purple },
	Keyword = { fg = c.clay },
	Label = { fg = c.clay },
	Macro = { fg = c.purple },
	Number = { fg = c.amber },
	Operator = { fg = c.fg_dim },
	PreProc = { fg = c.purple },
	Repeat = { fg = c.clay },
	Special = { fg = c.cyan },
	SpecialChar = { fg = c.cyan },
	SpecialComment = { fg = c.comment, bold = true },
	Statement = { fg = c.clay },
	StorageClass = { fg = c.clay },
	String = { fg = c.green },
	Structure = { fg = c.cyan },
	Tag = { fg = c.blue },
	Todo = { fg = c.bg, bg = c.amber, bold = true },
	Type = { fg = c.cyan },
	Typedef = { fg = c.cyan },
	Underlined = { fg = c.blue, underline = true },
	-- keep-sorted end
})

-- treesitter
set({
	-- keep-sorted start
	["@attribute"] = { link = "PreProc" },
	["@constant"] = { link = "Constant" },
	["@constant.builtin"] = { fg = c.clay },
	["@constructor"] = { fg = c.cyan },
	["@function.builtin"] = { fg = c.blue, italic = true },
	["@keyword.import"] = { link = "Include" },
	["@markup.heading"] = { fg = c.clay, bold = true },
	["@markup.italic"] = { italic = true },
	["@markup.link"] = { fg = c.blue, underline = true },
	["@markup.raw"] = { fg = c.green },
	["@markup.strong"] = { bold = true },
	["@module"] = { fg = c.fg },
	["@parameter"] = { fg = c.fg_bright, italic = true },
	["@property"] = { fg = c.fg },
	["@punctuation.bracket"] = { fg = c.fg_dim },
	["@punctuation.delimiter"] = { fg = c.fg_dim },
	["@tag.attribute"] = { fg = c.amber },
	["@tag.delimiter"] = { fg = c.fg_dim },
	["@variable"] = { fg = c.fg },
	["@variable.builtin"] = { fg = c.clay, italic = true },
	["@variable.parameter"] = { fg = c.fg_bright, italic = true },
	-- keep-sorted end
})

-- diagnostics
set({
	-- keep-sorted start block=yes
	DiagnosticError = { fg = c.red },
	DiagnosticHint = { fg = c.cyan },
	DiagnosticInfo = { fg = c.blue },
	DiagnosticOk = { fg = c.green },
	DiagnosticUnderlineError = {
		fg = c.red,
		sp = c.red,
		underline = true,
		bg = "#3d2926",
	},
	DiagnosticUnderlineHint = { fg = c.cyan, sp = c.cyan, underline = true },
	DiagnosticUnderlineInfo = { fg = c.blue, sp = c.blue, underline = true },
	DiagnosticUnderlineWarn = { fg = c.amber, sp = c.amber, underline = true },
	DiagnosticWarn = { fg = c.amber },
	-- keep-sorted end
})

-- diff and git
set({
	-- keep-sorted start
	Added = { fg = c.green },
	Changed = { fg = c.amber },
	DiffAdd = { bg = "#2f3a2a" },
	DiffChange = { bg = "#3a3526" },
	DiffDelete = { fg = c.red, bg = "#3d2926" },
	DiffText = { bg = "#51472c", bold = true },
	Removed = { fg = c.red },
	-- keep-sorted end
})

-- terminal palette
local terminal = {
	c.bg_select,
	c.red,
	c.green,
	c.amber,
	c.blue,
	c.purple,
	c.cyan,
	c.fg_dim,
	c.comment,
	c.red,
	c.green,
	c.amber,
	c.blue,
	c.purple,
	c.cyan,
	c.fg_bright,
}
for i, color in ipairs(terminal) do
	vim.g["terminal_color_" .. (i - 1)] = color
end
