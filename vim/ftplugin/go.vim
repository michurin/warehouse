set number

set tabstop=4
set noexpandtab
set softtabstop=0
set shiftwidth=4
set autoindent

filetype indent on

set synmaxcol=10000  " default 3000

set list lcs=trail:·,tab:▹·
set fillchars+=stl:─,stlnc:─,vert:│,fold:-
highlight VertSplit     ctermfg=black ctermbg=blue
highlight StatusLine    ctermfg=black ctermbg=blue
highlight StatusLineNC  ctermfg=black ctermbg=green

map <F7> :set list lcs=trail:\ ,tab:\ \ <CR>:set nonumber<CR>
imap <F7> <ESC>:set list lcs=trail:\ ,tab:\ \ <CR>:set nonumber<CR>
map <F8> :set list lcs=trail:·,tab:▹·<CR>:set number<CR>
imap <F8> <ESC>:set list lcs=trail:·,tab:▹·<CR>:set number<CR>
map <F9> :%!gofmt<CR>
imap <F9> <ESC>:%!gofmt<CR>

setlocal omnifunc=go#complete#Complete
" setlocal omnifunc=go#complete#Complete Last set from ~/.vim/pack/plugins/start/vim-go/ftplugin/go.vim
