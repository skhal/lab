" Copyright 2025 Samvel Khalatyan. All rights reserved.

function! relatedfile#OpenHeader(file)
  call s:open(a:file, s:header_by_file)
endfunction

function! relatedfile#OpenSource(file)
  call s:open(a:file, s:source_by_file)
endfunction

function! relatedfile#OpenTest(file)
  call s:open(a:file, s:test_by_file)
endfunction

function! s:open(file, substitutor_by_filetype)
  let l:file = fnamemodify(a:file, ':p')
  let l:Subsitutor = get(a:substitutor_by_filetype, &filetype, '')
  if l:Subsitutor == ''
    echoerr 'unsupported file' . a:file
    return
  endif
  let l:relatedfile = l:Subsitutor(l:file)
  if l:relatedfile == l:file
    return
  endif
  execute 'edit ' . l:relatedfile
endfunction

function! s:makeSubstitutor(pattern, string)
  return { filename -> substitute(filename, a:pattern, a:string, '') }
endfunction

let s:header_by_file = {
  \ 'cpp': s:makeSubstitutor('\(_test\)\?\.cc$', '.h'),
  \}
let s:source_by_file = {
  \ 'c': s:makeSubstitutor('\(_test\.cc\|\.h\)$', '.cc'),
  \ 'cpp': s:makeSubstitutor('\(_test\.cc\|\.h\)$', '.cc'),
  \}
let s:test_by_file = {
  \ 'c': s:makeSubstitutor('\(\(_test\)\@<!\.cc\|\.h\)$', '_test.cc'),
  \ 'cpp': s:makeSubstitutor('\(\(_test\)\@<!\.cc\|\.h\)$', '_test.cc'),
  \}
