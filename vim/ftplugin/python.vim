set number

set tabstop=8
set expandtab
set softtabstop=4
set shiftwidth=4
set autoindent

filetype indent on

" set background=dark

set foldmethod=indent
set foldlevel=99

set list lcs=trail:·,tab:▹·
" setlocal foldmethod=indent

function! Insert_import_pdb()
	:normal oimport pdb
	:normal 0
endfunction

function! Insert_set_pdb()
	:normal opdb.set_trace()
	:normal 0
endfunction

function! Mark_up_text()
	:set list lcs=trail:·,tab:▹·
	:set number
endfunction

function! Mark_down_text()
	:set list!
	:set number!
endfunction

:command! Pi :call Insert_import_pdb()
:command! Ps :call Insert_set_pdb()
:command! Mu :call Mark_up_text()
:command! Md :call Mark_down_text()

