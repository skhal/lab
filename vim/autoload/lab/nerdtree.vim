" Copyright 2025 Samvel Khalatyan. All rights reserved.
"
" Use of this source code is governed by a BSD-style
" license that can be found in the LICENSE file.

function! lab#nerdtree#Toggle()
  NERDTreeToggle
  if &filetype ==# 'nerdtree'
    " Focus back from NERDTree window if open to reveal current file path
    wincmd p
    if !&diff && strlen(expand('%')) > 0
      NERDTreeFind
    endif
  endif
endfunction

" s:isNERDTreeOpen checks whether NERDTre is open.
" https://codeyarns.com/tech/2014-05-05-how-to-highlight-current-file-in-nerdtree.html#gsc.tab=0
function! s:isNERDTreeOpen()
  return exists("t:NERDTreeBufName") && (bufwinnr(t:NERDTreeBufName) != -1)
endfunction
