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
