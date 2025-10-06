syntax on
set number"

set bg=dark
set colorcolumn=80
set cursorline

set hlsearch
set incsearch

set tabstop=2
set softtabstop=2
set shiftwidth=2
set smartindent
set autoindent
set expandtab
au FileType go,make setlocal noexpandtab

set listchars=eol:¬,extends:›,precedes:‹,space:░,tab:«–»,trail:•
nmap <leader>l <esc>:set list!<cr>
