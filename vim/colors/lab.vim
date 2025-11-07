" Copyright 2025 Samvel Khalatyan. All rights reserved.

hi clear
if exists ("syntax_on")
  syntax reset
endif
let g:colors_name = 'lab'
let s:none = ['NONE', 'NONE']

" s:h configures highlights for the scope.
function! s:h(scope, fg, ...) " bg, attr_list
  let l:fg = copy(a:fg)
  let l:bg = get(a:, 1, s:none)

  let l:attr_list = filter(get(a:, 2, ['NONE']), 'type(v:val) == 1')
  let l:attrs = len(l:attr_list) > 0 ? join(l:attr_list, ',') : 'NONE'

  let l:hl_string = [
    \ 'highlight', a:scope,
    \ 'guifg=' . l:fg[0], 'ctermfg=' . l:fg[1],
    \ 'guibg=' . l:bg[0], 'ctermbg=' . l:bg[1],
    \ 'gui=' . l:attrs, 'cterm=' . l:attrs,
    \]

  execute join(l:hl_string, ' ')
endfunction

if &background ==# 'dark'
  " keep-sorted start
  call s:h('Boolean', ['#cf7ea9', '175'])
  call s:h('Character', ['#d8884b', '173'])
  call s:h('ColorColumn', s:none, ['#2b2c2f', '236'])
  call s:h('Comment', ['#8c9292', '243'])
  call s:h('Constant', ['#79a9c4', '74'])
  call s:h('CursorLine', s:none, ['#333437', '236'])
  call s:h('CursorLineNR', ['#858585', '102'])
  call s:h('DiffAdd', s:none, ['#173e1c', '22'])
  call s:h('DiffChange' , s:none, ['#870087', '90'])
  call s:h('DiffDelete', s:none, ['#2c1a1e', '234'])
  call s:h('DiffText', s:none, ['#af00af', '127'], ['bold'])
  call s:h('Error', ['#f14c4c' , '203'])
  call s:h('Folded', s:none, ['#1f2c3b', '235'])
  call s:h('Function', ['#75b3db', '74'])
  call s:h('Identifier', ['#75b3db', '74'])
  call s:h('Keyword', ['#d8884b', '173'])
  call s:h('LineNr', ['#858585', '102'])
  call s:h('Normal', ['#b5bcc5', '250'])
  call s:h('Number', ['#79a9c4', '74'])
  call s:h('Operator', ['#acb4be', '249'])
  call s:h('Pmenu', ['#acb4be', '249'], ['#3a4659', '238'])
  call s:h('PmenuSel', s:none, s:none, ['reverse', 'bold'])
  call s:h('Search', s:none, ['#603216', '52'])
  call s:h('StatusLine', ['#ffffff', '231'], ['#363636', '237'], ['bold'])
  call s:h('StatusLineNC', ['#ffffff', '231'], ['#363636', '237'])
  call s:h('StorageClass', ['#d8884b', '173'])
  call s:h('String', ['#97c67b', '150'])
  call s:h('TabLine', ['#dadde2', '253'], ['#2a2a2b', '235'])
  call s:h('TabLineFill', s:none, ['#1c1d21', '234'])
  call s:h('TabLineSel', ['#ffffff', '231'], ['#1c1d21', '234'], ['bold'])
  call s:h('Type', ['#eede7b', '173'])
  call s:h('VertSplit', ['#a2a7ac', '248'])
  call s:h('Visual', s:none, ['#264f78', '24'])
  call s:h('javaAnnotation', ['#9dbed9', '110'])
  " keep-sorted end
else
  " keep-sorted start
  call s:h('Boolean', ['#221199', '19'])
  call s:h('Character', ['#770088', '90'])
  call s:h('ColorColumn', s:none, ['#f3f3f3', '255'])
  call s:h('Comment', ['#880000', '88'])
  call s:h('Constant', ['#116644', '29'])
  call s:h('CursorLine', s:none, ['#ffffe0', '230'])
  call s:h('CursorLineNR', ['#999999', '247'])
  call s:h('DiffAdd', s:none, ['#d9ffd9', '194'])
  call s:h('DiffChange', s:none, ['#ffd7ff', '225'])
  call s:h('DiffDelete', s:none, ['#ffecec', '224'])
  call s:h('DiffText', s:none, ['#ffafff', '219'], ['bold'])
  call s:h('Error', ['#e51400', '160'])
  call s:h('Folded', s:none, ['#e6f2ff', '255'])
  call s:h('Function' , ['#0000ff', '21'])
  call s:h('Identifier', ['#0000ff', '21'])
  call s:h('Keyword', ['#770088', '90'])
  call s:h('LineNr', ['#999999', '247'])
  call s:h('Normal', ['#3c4043', '238'], ['#ffffff', '231'])
  call s:h('Number', ['#116644', '29'])
  call s:h('Operator', ['#000000', '16'])
  call s:h('Pmenu', ['#3c4043', '238'], ['#eeeeee', '255'])
  call s:h('PmenuSel', ['#1967d2', '26'], ['#e4e4e4', '254'])
  call s:h('Search', s:none, ['#f8c8aa', '223'])
  call s:h('StatusLine', ['#666666', '241'], ['#f5f5f5', '255'], ['bold'])
  call s:h('StatusLineNC', ['#666666', '241'], ['#f5f5f5', '255'])
  call s:h('StorageClass', ['#770088', '90'])
  call s:h('String', ['#008800', '28'])
  call s:h('TabLine', ['#353637', '237'], ['#f5f5f5', '255'])
  call s:h('TabLineFill', s:none, ['#f5f5f5', '255'])
  call s:h('TabLineSel', ['#333333', '236'], ['#ffffff', '231'], ['bold'])
  call s:h('Type', ['#008855', '29'])
  call s:h('VertSplit', ['#b1b1b1', '249'])
  call s:h('Visual', s:none, ['#add6ff', '153'])
  call s:h('javaAnnotation', ['#555555', '240'])
  " keep-sorted end
endif

" keep-sorted start
hi! link FoldColumn Folded
hi! link NonText Normal
hi! link SpecialChar Character
hi! link Statement Keyword
hi! link cCustomClass Identifier
hi! link javaCommentTitle Comment
hi! link javaDocTags Comment
hi! link javaExternal Statement
hi! link javaScopeDecl Statement
hi! link javaStorageClass Statement
hi! link pbStructure Keyword
" keep-sorted end
