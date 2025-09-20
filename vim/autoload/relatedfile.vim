" Copyright 2025 Samvel Khalatyan. All rights reserved.

function! relatedfile#OpenHeader(file)
  call s:open(a:file, s:header_by_file, 'header')
endfunction

function! relatedfile#OpenSource(file)
  call s:open(a:file, s:source_by_file, 'source')
endfunction

function! relatedfile#OpenTest(file)
  call s:open(a:file, s:test_by_file, 'test')
endfunction

function! relatedfile#OpenExample(file)
  call s:open(a:file, s:example_by_file, 'example')
endfunction

function! s:open(file, substitutor_by_filetype, reltype)
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
  echow a:reltype . ': ' . relatedfile
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
  \ 'go': s:makeSubstitutor('\(_example\)\?_test\.go$', '.go'),
  \}
let s:test_by_file = {
  \ 'c': s:makeSubstitutor('\(\(_test\)\@<!\.cc\|\.h\)$', '_test.cc'),
  \ 'cpp': s:makeSubstitutor('\(\(_test\)\@<!\.cc\|\.h\)$', '_test.cc'),
  \ 'go': s:makeSubstitutor('\(_example_test\.go\|\(_test\)\@<!\.go\)$', '_test.go'),
  \}
let s:example_by_file = {
  \ 'go': s:makeSubstitutor('\(\(_example\)\@<!_test\.go\|\(example_test\)\@<!\.go\)$', '_example_test.go'),
  \}
