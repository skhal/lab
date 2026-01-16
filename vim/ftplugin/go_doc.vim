" Copyright 2026 Samvel Khalatyan. All rights reserved.
"
" Use of this source code is governed by a BSD-style
" license that can be found in the LICENSE file.


if exists('g:go_doc_loaded')
  finish
endif
let g:go_doc_loaded=1

command! -nargs=* GoDoc call go#Doc()
