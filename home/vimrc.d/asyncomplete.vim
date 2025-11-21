" Copyright 2025 Samvel Khalatyan. All rights reserved.
"
" Use of this source code is governed by a BSD-style
" license that can be found in the LICENSE file.
"
" -- Plugin https://github.com/prabirshrestha/asyncomplete.vim
let g:asyncomplete_auto_popup = 0

function! s:check_back_space() abort
  let col = col('.') - 1
  return !col || getline('.')[col - 1] =~ '\s'
endfunction

inoremap <silent><expr> <Tab>
  \ pumvisible() ? "\<C-n>" :
  \ <SID>check_back_space() ? "\<TAB>" :
  \ asyncomplete#force_refresh()
inoremap <expr><S-Tab> pumvisible() ? "\<C-p>" : "\<C-h>"
