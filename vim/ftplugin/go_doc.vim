" Copyright 2025 Samvel Khalatyan. All rights reserved.

if exists('g:go_doc_loaded')
  finish
endif
let g:go_doc_loaded=1

command! -nargs=* GoDoc call go#Doc()
