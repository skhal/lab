" Copyright 2025 Samvel Khalatyan. All rights reserved.

" go#Doc calls `go doc` for the keyword under cursor. It allows for dots in
" the keyword name. The result depends on where the cursor is positioned - it
" runs `go doc` for the path up to and including current token in dot-separated
" fully qualified name.
function! go#Doc()
  let l:name = s:getName()
  execute printf('!go doc %s', l:name)
endfunction

function! s:getName()
  let l:fqname = s:getFullyQualifiedName()
  let l:name = expand('<cword>')
  let l:idx = stridx(l:fqname, l:name)
  let l:idx_dot = stridx(l:fqname, '.', l:idx)
  if l:idx_dot == -1
    return l:fqname
  endif
  let l:idx_dot -= 1  " exclude dot itself
  return l:fqname[:l:idx_dot]
endfunction

function! s:getFullyQualifiedName()
  let l:iskeyword_save = &iskeyword
  setlocal iskeyword+=.
  let l:word = expand('<cword>')
  let &iskeyword = l:iskeyword_save
  return l:word
endfunction
