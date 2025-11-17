" Copyright 2025 Samvel Khalatyan. All rights reserved.
"
" Use of this source code is governed by a BSD-style
" license that can be found in the LICENSE file.
"
" -- Plugin https://github.com/prabirshrestha/vim-lsp
" keep-sorted start
let g:lsp_diagnostics_echo_cursor = 1
let g:lsp_diagnostics_virtual_text_enabled = 0
let g:lsp_use_native_client = 1
" keep-sorted end

" Gopls https://go.dev/gopls/editor/vim#a-hrefvimlsp-idvimlspvim-lspa
augroup LspGo
  au!
  autocmd User lsp_setup call lsp#register_server({
      \ 'name': 'gopls',
      \ 'cmd': {server_info->['gopls']},
      \ 'whitelist': ['go'],
      \ })
  autocmd FileType go setlocal omnifunc=lsp#complete
  autocmd FileType go nmap <buffer> gchi <plug>(lsp-call-hierarchy-incoming)
  autocmd FileType go nmap <buffer> gcho <plug>(lsp-call-hierarchy-outgoing)
  autocmd FileType go nmap <buffer> gca <plug>(lsp-code-action)
  autocmd FileType go nmap <buffer> gcaf <plug>(lsp-code-action-float)
  autocmd FileType go nmap <buffer> gcap <plug>(lsp-code-action-preview)
  autocmd FileType go nmap <buffer> gcl <plug>(lsp-code-lens)
  autocmd FileType go nmap <buffer> gd <plug>(lsp-declaration)
  autocmd FileType go nmap <buffer> gpd <plug>(lsp-peek-declaration)
  autocmd FileType go nmap <buffer> gd <plug>(lsp-definition)
  autocmd FileType go nmap <buffer> gpd <plug>(lsp-peek-definition)
  autocmd FileType go nmap <buffer> gds <plug>(lsp-document-symbol)
  autocmd FileType go nmap <buffer> gdss <plug>(lsp-document-symbol-search)
  autocmd FileType go nmap <buffer> gdd <plug>(lsp-document-diagnostics)
  autocmd FileType go nmap <buffer> gh <plug>(lsp-hover)
  autocmd FileType go nmap <buffer> ghf <plug>(lsp-hover-float)
  autocmd FileType go nmap <buffer> ghp <plug>(lsp-hover-preview)
  autocmd FileType go nmap <buffer> gpc <plug>(lsp-preview-close)
  autocmd FileType go nmap <buffer> gpf <plug>(lsp-preview-focus)
  autocmd FileType go nmap <buffer> gne <plug>(lsp-next-error)
  autocmd FileType go nmap <buffer> gnenw <plug>(lsp-next-error-nowrap)
  autocmd FileType go nmap <buffer> gpe <plug>(lsp-previous-error)
  autocmd FileType go nmap <buffer> gpenw <plug>(lsp-previous-error-nowrap)
  autocmd FileType go nmap <buffer> gnw <plug>(lsp-next-warning)
  autocmd FileType go nmap <buffer> gnwnw <plug>(lsp-next-warning-nowrap)
  autocmd FileType go nmap <buffer> gpw <plug>(lsp-previous-warning)
  autocmd FileType go nmap <buffer> gpwnw <plug>(lsp-previous-warning-nowrap)
  autocmd FileType go nmap <buffer> gnd <plug>(lsp-next-diagnostic)
  autocmd FileType go nmap <buffer> gndnw <plug>(lsp-next-diagnostic-nowrap)
  autocmd FileType go nmap <buffer> gpd <plug>(lsp-previous-diagnostic)
  autocmd FileType go nmap <buffer> gpdnw <plug>(lsp-previous-diagnostic-nowrap)
  autocmd FileType go nmap <buffer> gr <plug>(lsp-reference)
  autocmd FileType go nmap <buffer> gr <plug>(lsp-rename)
  autocmd FileType go nmap <buffer> gtd <plug>(lsp-type-definition)
  autocmd FileType go nmap <buffer> gth <plug>(lsp-type-hierarchy)
  autocmd FileType go nmap <buffer> gptd <plug>(lsp-peek-type-definition)
  autocmd FileType go nmap <buffer> gws <plug>(lsp-workspace-symbol)
  autocmd FileType go nmap <buffer> gwss <plug>(lsp-workspace-symbol-search)
  autocmd FileType go nmap <buffer> gdf <plug>(lsp-document-format)
  autocmd FileType go nmap <buffer> gdrf <plug>(lsp-document-range-format)
  autocmd FileType go nmap <buffer> gi <plug>(lsp-implementation)
  autocmd FileType go nmap <buffer> gpi <plug>(lsp-peek-implementation)
  autocmd FileType go nmap <buffer> gs <plug>(lsp-status)
  autocmd FileType go nmap <buffer> gnr <plug>(lsp-next-reference)
  autocmd FileType go nmap <buffer> gpr <plug>(lsp-previous-reference)
  autocmd FileType go nmap <buffer> gsh <plug>(lsp-signature-help)
augroup END

augroup LspCpp
  au!
  autocmd User lsp_setup call lsp#register_server({
      \ 'name': 'clangd',
      \ 'cmd': {server_info->['clangd']},
      \ 'whitelist': ['c', 'cpp'],
      \ })
  autocmd FileType c,cpp setlocal omnifunc=lsp#complete
  autocmd FileType c,cpp nmap <buffer> gchi <plug>(lsp-call-hierarchy-incoming)
  autocmd FileType c,cpp nmap <buffer> gcho <plug>(lsp-call-hierarchy-outgoing)
  autocmd FileType c,cpp nmap <buffer> gca <plug>(lsp-code-action)
  autocmd FileType c,cpp nmap <buffer> gcaf <plug>(lsp-code-action-float)
  autocmd FileType c,cpp nmap <buffer> gcap <plug>(lsp-code-action-preview)
  autocmd FileType c,cpp nmap <buffer> gcl <plug>(lsp-code-lens)
  autocmd FileType c,cpp nmap <buffer> gd <plug>(lsp-declaration)
  autocmd FileType c,cpp nmap <buffer> gpd <plug>(lsp-peek-declaration)
  autocmd FileType c,cpp nmap <buffer> gd <plug>(lsp-definition)
  autocmd FileType c,cpp nmap <buffer> gpd <plug>(lsp-peek-definition)
  autocmd FileType c,cpp nmap <buffer> gds <plug>(lsp-document-symbol)
  autocmd FileType c,cpp nmap <buffer> gdss <plug>(lsp-document-symbol-search)
  autocmd FileType c,cpp nmap <buffer> gdd <plug>(lsp-document-diagnostics)
  autocmd FileType c,cpp nmap <buffer> gh <plug>(lsp-hover)
  autocmd FileType c,cpp nmap <buffer> ghf <plug>(lsp-hover-float)
  autocmd FileType c,cpp nmap <buffer> ghp <plug>(lsp-hover-preview)
  autocmd FileType c,cpp nmap <buffer> gpc <plug>(lsp-preview-close)
  autocmd FileType c,cpp nmap <buffer> gpf <plug>(lsp-preview-focus)
  autocmd FileType c,cpp nmap <buffer> gne <plug>(lsp-next-error)
  autocmd FileType c,cpp nmap <buffer> gnenw <plug>(lsp-next-error-nowrap)
  autocmd FileType c,cpp nmap <buffer> gpe <plug>(lsp-previous-error)
  autocmd FileType c,cpp nmap <buffer> gpenw <plug>(lsp-previous-error-nowrap)
  autocmd FileType c,cpp nmap <buffer> gnw <plug>(lsp-next-warning)
  autocmd FileType c,cpp nmap <buffer> gnwnw <plug>(lsp-next-warning-nowrap)
  autocmd FileType c,cpp nmap <buffer> gpw <plug>(lsp-previous-warning)
  autocmd FileType c,cpp nmap <buffer> gpwnw <plug>(lsp-previous-warning-nowrap)
  autocmd FileType c,cpp nmap <buffer> gnd <plug>(lsp-next-diagnostic)
  autocmd FileType c,cpp nmap <buffer> gndnw <plug>(lsp-next-diagnostic-nowrap)
  autocmd FileType c,cpp nmap <buffer> gpd <plug>(lsp-previous-diagnostic)
  autocmd FileType c,cpp nmap <buffer> gpdnw <plug>(lsp-previous-diagnostic-nowrap)
  autocmd FileType c,cpp nmap <buffer> gr <plug>(lsp-reference)
  autocmd FileType c,cpp nmap <buffer> gr <plug>(lsp-rename)
  autocmd FileType c,cpp nmap <buffer> gtd <plug>(lsp-type-definition)
  autocmd FileType c,cpp nmap <buffer> gth <plug>(lsp-type-hierarchy)
  autocmd FileType c,cpp nmap <buffer> gptd <plug>(lsp-peek-type-definition)
  autocmd FileType c,cpp nmap <buffer> gws <plug>(lsp-workspace-symbol)
  autocmd FileType c,cpp nmap <buffer> gwss <plug>(lsp-workspace-symbol-search)
  autocmd FileType c,cpp nmap <buffer> gdf <plug>(lsp-document-format)
  autocmd FileType c,cpp nmap <buffer> gdrf <plug>(lsp-document-range-format)
  autocmd FileType c,cpp nmap <buffer> gi <plug>(lsp-implementation)
  autocmd FileType c,cpp nmap <buffer> gpi <plug>(lsp-peek-implementation)
  autocmd FileType c,cpp nmap <buffer> gs <plug>(lsp-status)
  autocmd FileType c,cpp nmap <buffer> gnr <plug>(lsp-next-reference)
  autocmd FileType c,cpp nmap <buffer> gpr <plug>(lsp-previous-reference)
  autocmd FileType c,cpp nmap <buffer> gsh <plug>(lsp-signature-help)
augroup END
