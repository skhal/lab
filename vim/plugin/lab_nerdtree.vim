" Copyright 2025 Samvel Khalatyan. All rights reserved.
"
" Use of this source code is governed by a BSD-style
" license that can be found in the LICENSE file.

if exists('g:lab_nerdtree_loaded')
  finish
endif
let g:lab_nerdtree_loaded=1

nnoremap <silent> <plug>(lab-nerdtree-toggle) :<c-u>call lab#nerdtree#Toggle()<cr>
