:set autoindent
:set number

function! Insert_use_dumper()
	:normal Ouse Data::Dumper; $Data::Dumper::Maxdepth=2; $Data::Dumper::Indent=3; $Data::Dumper::Terse=1; $Data::Dumper::Sortkeys=1;
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

:command Usedd :call Insert_use_dumper()
:command Mu :call Mark_up_text()
:command Md :call Mark_down_text()

" let perl_fold=1
" set foldenable foldmethod=syntax foldlevelstart=1
set foldlevelstart=1000
set foldlevel=1000
set foldmarker={,}
set foldmethod=marker
" set foldtext=substitute(getline(v:foldstart),'{.*','{...}',)
" set foldcolumn=4
" set foldlevelstart=1000

:hi Folded ctermfg=2 ctermbg=0
set list lcs=trail:·,tab:▹·

highlight Comment ctermfg=5
highlight Todo ctermbg=1 ctermfg=7

set laststatus=2
au InsertEnter * hi StatusLine term=reverse ctermbg=2 ctermfg=0
au InsertLeave * hi StatusLine term=reverse ctermbg=0 ctermfg=0
highlight StatusLine term=reverse ctermbg=0 ctermfg=0
