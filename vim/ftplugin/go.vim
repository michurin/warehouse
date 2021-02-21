set number
highlight LineNr ctermfg=grey

set tabstop=4
set noexpandtab
set softtabstop=0
set shiftwidth=4
set autoindent

set foldmethod=syntax
set foldlevelstart=1000
set foldlevel=1000

"filetype indent on

set synmaxcol=10000  " default 3000

set list lcs=trail:·,tab:▹·

map <F7> :set list lcs=trail:\ ,tab:\ \ <CR>:set nonumber<CR>
imap <F7> <ESC>:set list lcs=trail:\ ,tab:\ \ <CR>:set nonumber<CR>
map <F8> :set list lcs=trail:·,tab:▹·<CR>:set number<CR>
imap <F8> <ESC>:set list lcs=trail:·,tab:▹·<CR>:set number<CR>
map <F9> :%!gofmt<CR>
imap <F9> <ESC>:%!gofmt<CR>

" setlocal omnifunc=go#complete#Complete
" setlocal omnifunc=go#complete#Complete Last set from ~/.vim/pack/plugins/start/vim-go/ftplugin/go.vim

set autowrite

" inoremap <buffer> . .<C-x><C-o> " autocomplete after dot

" REMAP CTRL-] (ugly)
nmap <silent> <c-]> <Plug>(coc-definition)
nmap <silent> <c-t> <c-o>
