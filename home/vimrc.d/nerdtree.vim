" Copyright 2026 Samvel Khalatyan. All rights reserved.
"
" Use of this source code is governed by a BSD-style
" license that can be found in the LICENSE file.
"
" -- Plugin https://github.com/preservim/nerdtree
" keep-sorted start
let g:NERDTreeHijackNetrw = 0
let g:NERDTreeMinimalMenu = 1
let g:NERDTreeMinimalUI = 1
let g:NERDTreeWinPos = "right"
" keep-sorted end

augroup LabNERDTree
  autocmd VimEnter * nnoremap <Leader>t <plug>(lab-nerdtree-toggle)
augroup END
