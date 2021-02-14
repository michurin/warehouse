set number

set tabstop=4
set noexpandtab
set softtabstop=0
set shiftwidth=4
set autoindent

set foldmethod=syntax
let g:go_fold_enable = ['block', 'import', 'varconst', 'package_comment']
set foldlevelstart=1000
set foldlevel=1000

" set to zero if vim get luggish
let g:go_highlight_types = 1
let g:go_highlight_fields = 1
let g:go_highlight_functions = 1
let g:go_highlight_function_calls = 1

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

set autowrite

" Go syntax highlighting
let g:go_highlight_fields = 1
let g:go_highlight_functions = 1
let g:go_highlight_function_calls = 1
let g:go_highlight_extra_types = 1
let g:go_highlight_operators = 1

" Auto formatting and importing
let g:go_fmt_autosave = 1
let g:go_fmt_command = "goimports"

" Status line types/signatures
let g:go_auto_type_info = 1

set completeopt-=preview

" Run :GoBuild or :GoTestCompile based on the go file
function! s:build_go_files()
  let l:file = expand('%')
  if l:file =~# '^\f\+_test\.go$'
    call go#test#Test(0, 1)
  elseif l:file =~# '^\f\+\.go$'
    call go#cmd#Build(0)
  endif
endfunction

" Map keys for most used commands.
" Ex: `\b` for building, `\r` for running and `\b` for running test.
autocmd FileType go nmap <leader>b :<C-u>call <SID>build_go_files()<CR>
autocmd FileType go nmap <leader>r  <Plug>(go-run)
autocmd FileType go nmap <leader>t  <Plug>(go-test)

" inoremap <buffer> . .<C-x><C-o> " autocomplete after dot
