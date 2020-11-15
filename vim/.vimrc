set spell spelllang=ru_ru,en_us
syn match UrlNoSpell 'https\?:\/\/[^[:space:]]\+' contains=@NoSpell

set cryptmethod=blowfish2

set history=500
filetype plugin on
filetype indent on
set wildmenu
set ruler

syntax enable
set background=dark

set encoding=utf8

set nobackup
set nowb
set noswapfile

set lbr
set tw=500

set ai "Auto indent
set si "Smart indent
set wrap "Wrap lines

nnoremap <C-l> gt
nnoremap <C-h> gT

let g:netrw_banner=0 "Disable annoing banner
let g:netrw_browse_split=3 "Open files in a new tab (as 't' does)
let g:netrw_altv=1 "Open split on the right side
let g:netrw_liststyle=3 "Tree view

