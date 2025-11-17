" Copyright 2025 Samvel Khalatyan. All rights reserved.
"
" Use of this source code is governed by a BSD-style
" license that can be found in the LICENSE file.
"
" -- Keyword https://vi.stackexchange.com/questions/6228/how-can-i-get-vim-to-show-documentation-of-a-c-c-function
augroup Keyword
  au!
  autocmd FileType c,cpp setlocal keywordprg=cppman
  autocmd FileType go    setlocal keywordprg=:GoDoc
augroup END
