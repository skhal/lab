" Copyright 2026 Samvel Khalatyan. All rights reserved.
"
" Use of this source code is governed by a BSD-style
" license that can be found in the LICENSE file.

augroup Format
  au!
  " keep-sorted start
  au FileType bzl setlocal equalprg=buildifier
  au FileType c,cpp,proto setlocal equalprg=clang-format21\ -assume-filename=%
  au FileType go setlocal equalprg=goimports
  au FileType json,html setlocal equalprg=prettier\ --stdin-filepath=%
  au FileType markdown setlocal equalprg=markdownfmt
  au FileType pbtxt setlocal equalprg=txtpbfmt
  au FileType sh setlocal equalprg=shfmt\ -i\ 2\ -ci\ -bn
  au FileType yaml setlocal equalprg=yamlfmt\ -in
  " keep-sorted end
  " Format entire buffer `=G` or selection `=` with shortcut `<leader>fc`
  au FileType bzl,c,cpp,proto,go,json,html,markdown,pbtxt,sh,yaml map <leader>fc gg=G<c-o><c-o>
  au FileType bzl,c,cpp,proto,go,json,html,markdown,pbtxt,sh,yaml imap <leader>fc <esc>gg=G<c-o><c-o>
augroup END
