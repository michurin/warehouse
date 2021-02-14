" Mac: cd .vim && ln -s /usr/local/opt/fzf .
" Most common commands:
" :FZF
set rtp+=~/.vim/fzf
" brew install ripgrep
" git clone git@github.com:junegunn/fzf.vim.git ~/.vim/bundle/fzf.vim
" Most common commands (https://www.chrisatmachine.com/Neovim/08-fzf/, https://youtu.be/on1AzaZzQ7k):
" :Files - FZF with preview
" :Rg - look inside
" :Buffers
" :BLines - file (enormous)
" :Lines - search all buffers
set rtp+=~/.vim/bundle/fzf.vim

set spell spelllang=ru_ru,en_us
syn match UrlNoSpell 'https\?:\/\/[^[:space:]]\+' contains=@NoSpell

set cryptmethod=blowfish2
set backspace=indent,eol,start " to work on mac

set hlsearch
set incsearch
set laststatus=2 " always
" nnoremap ,/ :noh<CR>

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

" Coc replacement
"function! InsertTabWrapper()
"    let col = col('.') - 1
"    if !col || getline('.')[col - 1] !~ '\k'
"        return "\<tab>"
"    else
"        return "\<c-p>"
"    endif
"endfunction
"inoremap <expr> <tab> InsertTabWrapper()
"inoremap <s-tab> <c-n>

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
" :helptags ALL
" :help vim-go
"vim +GoInstallBinaries


let g:netrw_banner=0 "Disable annoing banner
"let g:netrw_browse_split=3 "Open files in a new tab (as 't' does)
let g:netrw_altv=1 "Open split on the right side
let g:netrw_liststyle=3 "Tree view
autocmd FileType netrw setl bufhidden=delete " or use :qa!
" set compatible " limit search to project
set path+=** " search all subdirs
set wildmenu " file search menu


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

" ---- COC

" https://github.com/neoclide/coc.nvim/wiki/Install-coc.nvim#using-vim8s-native-package-manager
" https://octetz.com/docs/2019/2019-04-24-vim-as-a-go-ide/
" https://github.com/neoclide/coc.nvim/wiki/Language-servers#go

" if hidden is not set, TextEdit might fail.
set hidden

" Some servers have issues with backup files, see #649
set nobackup
set nowritebackup

" =2 Better display for messages
set cmdheight=1

" Smaller updatetime for CursorHold & CursorHoldI
set updatetime=300

" don't give |ins-completion-menu| messages.
set shortmess+=c

" always show signcolumns
set signcolumn=yes
set signcolumn=no " TODO TO THINK

" Use tab for trigger completion with characters ahead and navigate.
" Use command ':verbose imap <tab>' to make sure tab is not mapped by other plugin.
inoremap <silent><expr> <TAB>
      \ pumvisible() ? "\<C-n>" :
      \ <SID>check_back_space() ? "\<TAB>" :
      \ coc#refresh()
inoremap <expr><S-TAB> pumvisible() ? "\<C-p>" : "\<C-h>"

function! s:check_back_space() abort
  let col = col('.') - 1
  return !col || getline('.')[col - 1]  =~# '\s'
endfunction

" Use <c-space> to trigger completion.
inoremap <silent><expr> <c-space> coc#refresh()

" Use <cr> to confirm completion, `<C-g>u` means break undo chain at current position.
" Coc only does snippet and additional edit on confirm.
inoremap <expr> <cr> pumvisible() ? "\<C-y>" : "\<C-g>u\<CR>"

" Use `[c` and `]c` to navigate diagnostics
nmap <silent> [c <Plug>(coc-diagnostic-prev)
nmap <silent> ]c <Plug>(coc-diagnostic-next)

" Remap keys for gotos
nmap <silent> gd <Plug>(coc-definition)
nmap <silent> gy <Plug>(coc-type-definition)
nmap <silent> gi <Plug>(coc-implementation)
nmap <silent> gr <Plug>(coc-references)

" Use K to show documentation in preview window
nnoremap <silent> K :call <SID>show_documentation()<CR>

function! s:show_documentation()
  if (index(['vim','help'], &filetype) >= 0)
    execute 'h '.expand('<cword>')
  else
    call CocAction('doHover')
  endif
endfunction

" Highlight symbol under cursor on CursorHold
autocmd CursorHold * silent call CocActionAsync('highlight')
hi CocHighlightText ctermbg=241 guibg=#665c54

" Remap for rename current word
nmap <leader>rn <Plug>(coc-rename)

" Remap for format selected region
vmap <leader>f  <Plug>(coc-format-selected)
nmap <leader>f  <Plug>(coc-format-selected)

augroup mygroup
  autocmd!
  " Setup formatexpr specified filetype(s).
  autocmd FileType typescript,json setl formatexpr=CocAction('formatSelected')
  " Update signature help on jump placeholder
  autocmd User CocJumpPlaceholder call CocActionAsync('showSignatureHelp')
augroup end

" Remap for do codeAction of selected region, ex: `<leader>aap` for current paragraph
vmap <leader>a  <Plug>(coc-codeaction-selected)
nmap <leader>a  <Plug>(coc-codeaction-selected)

" Remap for do codeAction of current line
nmap <leader>ac  <Plug>(coc-codeaction)
" Fix autofix problem of current line
nmap <leader>qf  <Plug>(coc-fix-current)

" Use `:Format` to format current buffer
command! -nargs=0 Format :call CocAction('format')

" Use `:Fold` to fold current buffer
command! -nargs=? Fold :call     CocAction('fold', <f-args>)


" Add diagnostic info for https://github.com/itchyny/lightline.vim
let g:lightline = {
      \ 'colorscheme': 'wombat',
      \ 'active': {
      \   'left': [ [ 'mode', 'paste' ],
      \             [ 'cocstatus', 'readonly', 'filename', 'modified' ] ]
      \ },
      \ 'component_function': {
      \   'cocstatus': 'coc#status'
      \ },
      \ }



" Using CocList
" Show all diagnostics
nnoremap <silent> <space>a  :<C-u>CocList diagnostics<cr>
" Manage extensions
nnoremap <silent> <space>e  :<C-u>CocList extensions<cr>
" Show commands
nnoremap <silent> <space>c  :<C-u>CocList commands<cr>
" Find symbol of current document
nnoremap <silent> <space>o  :<C-u>CocList outline<cr>
" Search workspace symbols
nnoremap <silent> <space>s  :<C-u>CocList -I symbols<cr>
" Do default action for next item.
nnoremap <silent> <space>j  :<C-u>CocNext<CR>
" Do default action for previous item.
nnoremap <silent> <space>k  :<C-u>CocPrev<CR>
" Resume latest coc list
nnoremap <silent> <space>p  :<C-u>CocListResume<CR>

