" Copyright 2025 Samvel Khalatyan. All rights reserved.

" go#Doc calls `go doc` for the keyword under cursor. It adds dot to the keyword
" allowed charactes to include package name.
function! go#Doc()
  let l:iskeyword_save = &iskeyword
  setlocal iskeyword+=.
  let l:word = expand('<cword>')
  execute printf('!go doc %s', l:word)
  let &iskeyword = l:iskeyword_save
endfunction
