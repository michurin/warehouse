call plug#begin() " https://github.com/junegunn/vim-plug +PlugInstall
  Plug 'neovim/nvim-lspconfig'
  Plug 'steelsojka/completion-buffers'
  Plug 'nvim-lua/completion-nvim'
  Plug 'junegunn/fzf' " , { 'do': { -> fzf#install() } }
  Plug 'junegunn/fzf.vim'
call plug#end()

" ---------- LSP

lua << EOF
local nvim_lsp = require('lspconfig')

-- Use an on_attach function to only map the following keys
-- after the language server attaches to the current buffer
local on_attach = function(client, bufnr)
  local function buf_set_keymap(...) vim.api.nvim_buf_set_keymap(bufnr, ...) end
  local function buf_set_option(...) vim.api.nvim_buf_set_option(bufnr, ...) end

  -- Enable completion triggered by <c-x><c-o>
  buf_set_option('omnifunc', 'v:lua.vim.lsp.omnifunc')

  -- Mappings.
  local opts = { noremap=true, silent=true }

  -- See `:help vim.lsp.*` for documentation on any of the below functions
  buf_set_keymap('n', 'gD', '<cmd>lua vim.lsp.buf.declaration()<CR>', opts)
  buf_set_keymap('n', 'gd', '<cmd>lua vim.lsp.buf.definition()<CR>', opts)
  buf_set_keymap('n', 'K', '<cmd>lua vim.lsp.buf.hover()<CR>', opts)
  buf_set_keymap('n', 'gi', '<cmd>lua vim.lsp.buf.implementation()<CR>', opts)
  buf_set_keymap('n', '<C-k>', '<cmd>lua vim.lsp.buf.signature_help()<CR>', opts)
  buf_set_keymap('n', '<space>wa', '<cmd>lua vim.lsp.buf.add_workspace_folder()<CR>', opts)
  buf_set_keymap('n', '<space>wr', '<cmd>lua vim.lsp.buf.remove_workspace_folder()<CR>', opts)
  buf_set_keymap('n', '<space>wl', '<cmd>lua print(vim.inspect(vim.lsp.buf.list_workspace_folders()))<CR>', opts)
  buf_set_keymap('n', 'gy', '<cmd>lua vim.lsp.buf.type_definition()<CR>', opts)
  buf_set_keymap('n', '<space>rn', '<cmd>lua vim.lsp.buf.rename()<CR>', opts)
  buf_set_keymap('n', '<space>ca', '<cmd>lua vim.lsp.buf.code_action()<CR>', opts)
  buf_set_keymap('n', 'gr', '<cmd>lua vim.lsp.buf.references()<CR>', opts)
  buf_set_keymap('n', '<space>e', '<cmd>lua vim.lsp.diagnostic.show_line_diagnostics()<CR>', opts)
  buf_set_keymap('n', '[d', '<cmd>lua vim.lsp.diagnostic.goto_prev()<CR>', opts)
  buf_set_keymap('n', ']d', '<cmd>lua vim.lsp.diagnostic.goto_next()<CR>', opts)
  buf_set_keymap('n', '<space>q', '<cmd>lua vim.lsp.diagnostic.set_loclist()<CR>', opts)
  buf_set_keymap('n', '<space>f', '<cmd>lua vim.lsp.buf.formatting()<CR>', opts)

  require'completion'.on_attach(client, bufnr) -- it has to be done after PLS
end

-- Use a loop to conveniently call 'setup' on multiple servers and
-- map buffer local keybindings when the language server attaches
-- All available servers https://github.com/neovim/nvim-lspconfig/tree/master/lua/lspconfig
-- TODO gofumpt but according extra vars from *exrc*
-- https://neovim.discourse.group/t/gopls-settings-buildflags/790
-- https://www.gitmemory.com/issue/golang/go/44204/781570319
local servers = {
  {name='gopls'},
  {name='intelephense'},
  {name='pyright'}
}
for _, lsp in ipairs(servers) do
  nvim_lsp[lsp.name].setup {
    on_attach = on_attach,
    flags = {
      debounce_text_changes = 150,
    },
    settings={
      python={
        analysis={
          useLibraryCodeForTypes = false,
          typeCheckingMode = "off"
        },
        linting = {
          pylintEnabled = true,
          enabled = true
        }
      },
      gopls = { -- https://github.com/golang/tools/blob/master/gopls/doc/settings.md
        gofumpt = true
      }
    }
  }
end
EOF

highlight LspDiagnosticsDefaultHint ctermfg=64 ctermbg=234
highlight LspDiagnosticsDefaultInformation ctermfg=31 ctermbg=234
highlight LspDiagnosticsDefaultWarning ctermfg=137 ctermbg=234
highlight LspDiagnosticsDefaultError ctermfg=124 ctermbg=234

" ---------- Go Stuff

" How to debug https://www.getman.io/posts/gopls/
lua << EOF
  -- https://github.com/golang/tools/blob/master/gopls/doc/vim.md#neovim-imports
  function goimports(timeout_ms)
    -- https://github.com/neovim/neovim/blob/23fe6dba138859c1c22850b9ce76219141f546a0/runtime/doc/lsp.txt#L135
    -- https://github.com/neovim/neovim/blob/c1f573fbc94aecd0f5841f7eb671be1a0a29758c/runtime/lua/vim/lsp/buf.lua#L174
    vim.lsp.buf.formatting_sync() -- hackish

    local context = { only = { "source.organizeImports" } }
    vim.validate { context = { context, "t", true } }

    local params = vim.lsp.util.make_range_params()
    params.context = context

    -- See the implementation of the textDocument/codeAction callback
    -- (lua/vim/lsp/handler.lua) for how to do this properly.
    local result = vim.lsp.buf_request_sync(0, "textDocument/codeAction", params, timeout_ms)
    -- Todo:
    -- add formatting
    -- https://github.com/neovim/nvim-lspconfig/issues/115#issuecomment-616844477
    -- local result = vim.lsp.buf_request_sync(0, "textDocument/formatting", params, timeout)
    if not result or next(result) == nil then return end
    local actions = result[1].result
    if not actions then return end
    local action = actions[1]

    -- textDocument/codeAction can return either Command[] or CodeAction[]. If it
    -- is a CodeAction, it can have either an edit, a command or both. Edits
    -- should be executed first.
    if action.edit or type(action.command) == "table" then
      if action.edit then
        vim.lsp.util.apply_workspace_edit(action.edit)
      end
      if type(action.command) == "table" then
        vim.lsp.buf.execute_command(action.command)
      end
    else
      vim.lsp.buf.execute_command(action)
    end
  end
EOF

function! s:GoAlt(cmd)
  let l:cf = expand('%:p')
  if l:cf[-8:] == '_test.go'
    let l:alt = l:cf[:-9].'.go'
  elseif l:cf[-3:] == '.go'
    let l:alt = l:cf[:-4].'_test.go'
  else
    echo 'Not a golang file'
    return
  endif
  if filereadable(l:alt)
    exec a:cmd.' '.l:alt
  else
    echo 'File '.l:alt.' not found'
    return
  endif
endfunction

command! GA call s:GoAlt('e')
command! GAA call s:GoAlt('bo vs')

" ---------- Misc

autocmd FileType go autocmd BufWritePre *.go lua goimports(1000)
autocmd FileType go setlocal number shiftwidth=4 tabstop=4 softtabstop=4 autoindent list lcs=trail:+,tab:▹· foldmethod=syntax foldlevelstart=99 foldlevel=99 synmaxcol=10000
autocmd FileType sh setlocal number shiftwidth=4 tabstop=4 softtabstop=4 expandtab autoindent list lcs=trail:+,tab:▹·
autocmd FileType zsh setlocal number shiftwidth=4 tabstop=4 softtabstop=4 expandtab autoindent list lcs=trail:+,tab:▹·
autocmd FileType python setlocal number shiftwidth=4 tabstop=4 softtabstop=4 expandtab autoindent list lcs=trail:+,tab:▹·
autocmd FileType vim setlocal number shiftwidth=2 tabstop=2 softtabstop=2 expandtab autoindent list lcs=trail:+,tab:▹·
autocmd FileType lua setlocal number shiftwidth=2 tabstop=2 softtabstop=2 expandtab autoindent list lcs=trail:+,tab:▹·
autocmd FileType css setlocal number shiftwidth=2 tabstop=8 softtabstop=2 expandtab autoindent list lcs=trail:+,tab:▹·
autocmd FileType html setlocal number shiftwidth=2 tabstop=8 softtabstop=2 expandtab autoindent list lcs=trail:+,tab:▹·
autocmd FileType javascript setlocal number shiftwidth=2 tabstop=8 softtabstop=2 expandtab autoindent list lcs=trail:+,tab:▹·
autocmd FileType json setlocal number shiftwidth=2 tabstop=8 softtabstop=2 expandtab autoindent list lcs=trail:+,tab:▹· foldmethod=syntax foldlevelstart=99 foldlevel=99
autocmd FileType make setlocal number tabstop=8 autoindent list lcs=trail:+,tab:▹· foldmethod=syntax foldlevelstart=99 foldlevel=99
autocmd FileType tcl setlocal number shiftwidth=2 tabstop=8 softtabstop=2 expandtab autoindent list lcs=trail:+,tab:▹· foldmethod=indent foldlevelstart=99 foldlevel=99

highlight Whitespace term=none cterm=none ctermfg=DarkGray gui=none guifg=DarkGray ctermbg=none
highlight EndOfBuffer term=none cterm=none ctermfg=DarkGray gui=none guifg=DarkGray ctermbg=none
highlight LineNr ctermfg=grey

set nofixendofline

" ---------- Spell

setlocal spell spelllang=en_us,ru_yo

syntax match UrlNoSpell 'https\?:\/\/[^[:space:]]\+' contains=@NoSpell

highlight SpellBad term=none cterm=underline ctermfg=none gui=bold guifg=none ctermbg=none
highlight SpellCap term=none cterm=underline ctermfg=none gui=bold guifg=none ctermbg=none
highlight SpellRare term=none cterm=underline ctermfg=none gui=bold guifg=none ctermbg=none
highlight SpellLocal term=none cterm=underline ctermfg=none gui=bold guifg=none ctermbg=none

" ---------- Completion (https://github.com/nvim-lua/completion-nvim)

let g:completion_chain_complete_list = {
  \'default': [
  \ {'complete_items': ['lsp']},
  \ {'complete_items': ['buffers']},
  \ {'mode': '<c-p>'},
  \ {'mode': '<c-n>'}]}

let g:completion_matching_strategy_list = ['exact', 'substring', 'fuzzy', 'all']
let g:completion_matching_smart_case = 1
let g:completion_trigger_on_delete = 1
let g:completion_auto_change_source = 1 " <Plug>(completion_next_source) / *_prev_*

set completeopt=menuone,noinsert,noselect

highlight Pmenu ctermfg=153 ctermbg=234
highlight PmenuSel ctermfg=153 ctermbg=240

" ---------- Grepping (https://github.com/junegunn/fzf/blob/master/README-VIM.md)

let g:fzf_layout = { 'window': { 'width': 0.95, 'height': 0.95 } }

function! s:MornGG(pat)
  let l:root = luaeval("vim.lsp.buf.list_workspace_folders()[1] or '.'") " Todo multi folders?
  let l:cmd = "rg --column --line-number --no-heading --color=always -g '!vendor' -g '!.git' -g '*.go' -g '!*_test.go' --smart-case -- ".shellescape(a:pat)." ".shellescape(l:root)
  call fzf#vim#grep(l:cmd, 1, fzf#vim#with_preview({'options': ['--prompt', 'GG> ']}), 0)
endfunction

command! -nargs=* GG call s:MornGG(<q-args>)
command! GM :execute 'lvimgrep /func[^()]*([^()]*\<'.expand('<cword>').'\>)/ '.expand('%:p:h').  '/*' | lopen

" ---------- Splash

set shortmess=atI " all abbreviations and truncate on CTRL-G, don't give intro

function! MornHelp()
  enew
  setlocal bufhidden=wipe buftype=nofile nobuflisted nocursorcolumn nocursorline nolist nonumber norelativenumber filetype=help noswapfile nospell
  syntax region helpNote start=":[A-Za-z]"hs=s+1 end=" "he=s-1
  syntax region helpVim start="^  " end="\n"
  syntax region helpUnderlined start="^##* *"hs=e+1 end="\n"
  syntax region helpOption start="`"hs=e+1 end="`"he=s-1
  let l:msg =<< EOF
   _      ____  ___   _      _   _
  | |\ | | |_  / / \ \ \  / | | | |\/|
  |_| \| |_|__ \_\_/  \_\/  |_| |_|  |

# Customization

## Go commands

:GM — grep method
:GG — fuzzy find string — `*.go` excluding `vendor/` and `*_test.go`
:GA — go alternate
:GAA — :GA vsplit
EOF
  call append(0, l:msg)
  setlocal nomodifiable nomodified
  nnoremap <buffer><silent> e :enew<CR>
  nnoremap <buffer><silent> i :enew <bar> startinsert<CR>
  nnoremap <buffer><silent> o :enew <bar> startinsert<CR>
endfunction

if argc() == 0
  autocmd VimEnter * call MornHelp()
endif

" Helpful links:
" https://github.com/junegunn/fzf.vim/blob/master/plugin/fzf.vim
" https://github.com/nanotee/nvim-lua-guide
" https://github.com/neovim/neovim/blob/master/runtime/doc/lsp.txt
" :h ins-completion-menu
" Memo:
" source $MYVIMRC
" :lcl closes it
" Todo:
" group settings by filetype
" https://github.com/airblade/vim-gitgutter
