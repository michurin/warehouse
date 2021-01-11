set spell spelllang=ru_ru,en_us
syn match UrlNoSpell 'https\?:\/\/[^[:space:]]\+' contains=@NoSpell

set cryptmethod=blowfish2

set hlsearch
set incsearch
nnoremap ,/ :noh<CR>

set history=500
filetype plugin on
filetype indent on
set wildmenu
set ruler

syntax enable

set encoding=utf8

set nobackup
set nowb
set noswapfile

set lbr
set tw=0 " disable autowrap on paste

set ai "Auto indent
set si "Smart indent
set wrap "Wrap lines

set autoread
autocmd FocusGained * checktime

nnoremap ,r gt
nnoremap ,e gT

nnoremap ,d :bprev<CR>
nnoremap ,f :bnext<CR>

nnoremap ,w :execute "vimgrep /\\<" . expand('<cword>') . "\\>/j **/*.go"
nnoremap ,a :cprev<CR>
nnoremap ,s :cnext<CR>

function! ToggleQuickFix()
    if empty(filter(getwininfo(), 'v:val.quickfix'))
        copen
    else
        cclose
    endif
endfunction

nnoremap ,q :call ToggleQuickFix()<cr>

" ---------

function! InsertTabWrapper()
    let col = col('.') - 1
    if !col || getline('.')[col - 1] !~ '\k'
        return "\<tab>"
    else
        return "\<c-p>"
    endif
endfunction
inoremap <expr> <tab> InsertTabWrapper()
inoremap <s-tab> <c-n>

" Navigate the complete menu items like CTRL+n / CTRL+p would.
inoremap <expr> <Down> pumvisible() ? "<C-n>" :"<Down>"
inoremap <expr> <Up> pumvisible() ? "<C-p>" : "<Up>"
" Select the complete menu item like CTRL+y would.
inoremap <expr> <Right> pumvisible() ? "<C-y>" : "<Right>"
inoremap <expr> <CR> pumvisible() ? "<C-y>" :"<CR>"
" Cancel the complete menu item like CTRL+e would.
inoremap <expr> <Left> pumvisible() ? "<C-e>" : "<Left>"

" ---------

"git clone https://github.com/fatih/vim-go.git ~/.vim/pack/plugins/start/vim-go
"vim +GoInstallBinaries


let g:netrw_banner=0 "Disable annoing banner
"let g:netrw_browse_split=3 "Open files in a new tab (as 't' does)
let g:netrw_altv=1 "Open split on the right side
let g:netrw_liststyle=3 "Tree view


"git clone https://github.com/preservim/nerdtree.git ~/.vim/pack/plugins/start/nerdtree
"map <C-n> :NERDTreeToggle<CR>


"git clone https://github.com/morhetz/gruvbox.git ~/.vim/pack/plugins/start/gruvbox
"let g:gruvbox_guisp_fallback="bg" " workaround for spell checking
"colorscheme gruvbox
"set background=dark
"let g:gruvbox_contrast_dark="hard"

"colorscheme industry
"runtime colors/industry.vim

hi SpellBad term=NONE cterm=underline ctermfg=NONE gui=bold guifg=NONE ctermbg=NONE
hi SpellLocal term=NONE cterm=underline ctermfg=NONE gui=bold guifg=NONE ctermbg=NONE
hi SpecialKey term=none cterm=none ctermfg=DarkGray gui=none guifg=DarkGray ctermbg=NONE

