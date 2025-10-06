" Copyright 2025 Samvel Khalatyan. All rights reserved.

if exists('g:loaded_relatedfile')
  finish
endif
let g:loaded_relatedfile=1

augroup relatedfile
  au!
  au FileType cpp    nmap <leader>rh <esc>:call relatedfile#OpenHeader(expand('%'))<cr>
  au FileType cpp,go nmap <leader>rc <esc>:call relatedfile#OpenSource(expand('%'))<cr>
  au FileType cpp,go nmap <leader>rt <esc>:call relatedfile#OpenTest(expand('%'))<cr>
  au FileType cpp    nmap <leader>rb <esc>:call relatedfile#OpenBuild(expand('%'))<cr>
  au FileType go     nmap <leader>re <esc>:call relatedfile#OpenExample(expand('%'))<cr>
augroup END
