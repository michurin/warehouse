call plug#begin() " https://github.com/junegunn/vim-plug +PlugInstall
  Plug 'neovim/nvim-lspconfig'
  Plug 'nvim-lua/completion-nvim'
  Plug 'junegunn/fzf' " , { 'do': { -> fzf#install() } }
  Plug 'junegunn/fzf.vim'
call plug#end()

" ---------- Go Stuff

lua << EOF
  -- https://github.com/golang/tools/blob/master/gopls/doc/vim.md#neovim-imports
  function goimports(timeout_ms)
    local context = { only = { "source.organizeImports" } }
    vim.validate { context = { context, "t", true } }

    local params = vim.lsp.util.make_range_params()
    params.context = context

    -- See the implementation of the textDocument/codeAction callback
    -- (lua/vim/lsp/handler.lua) for how to do this properly.
    local result = vim.lsp.buf_request_sync(0, "textDocument/codeAction", params, timeout_ms)
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
  buf_set_keymap('n', '<space>D', '<cmd>lua vim.lsp.buf.type_definition()<CR>', opts)
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
local servers = {'gopls', 'intelephense', 'pyright'}
for _, lsp in ipairs(servers) do
  nvim_lsp[lsp].setup {
    on_attach = on_attach,
    flags = {
      debounce_text_changes = 150,
    }
  }
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
autocmd FileType go setlocal number shiftwidth=4 tabstop=4 softtabstop=4 autoindent list lcs=trail:+,tab:▹·
autocmd FileType sh setlocal number shiftwidth=4 tabstop=4 softtabstop=4 expandtab autoindent list lcs=trail:+,tab:▹·
autocmd FileType python setlocal number shiftwidth=4 tabstop=4 softtabstop=4 expandtab autoindent list lcs=trail:+,tab:▹·
autocmd FileType vim setlocal number shiftwidth=2 tabstop=2 softtabstop=2 expandtab autoindent list lcs=trail:+,tab:▹·

highlight Whitespace term=none cterm=none ctermfg=DarkGray gui=none guifg=DarkGray ctermbg=none
highlight EndOfBuffer term=none cterm=none ctermfg=DarkGray gui=none guifg=DarkGray ctermbg=none

set nofixendofline

" ---------- Spell

setlocal spell spelllang=en_us,ru_yo

syn match UrlNoSpell 'https\?:\/\/[^[:space:]]\+' contains=@NoSpell

highlight SpellBad term=none cterm=underline ctermfg=none gui=bold guifg=none ctermbg=none
highlight SpellCap term=none cterm=underline ctermfg=none gui=bold guifg=none ctermbg=none
highlight SpellRare term=none cterm=underline ctermfg=none gui=bold guifg=none ctermbg=none
highlight SpellLocal term=none cterm=underline ctermfg=none gui=bold guifg=none ctermbg=none

" ---------- Completion (https://github.com/nvim-lua/completion-nvim)

set completeopt=menuone,noinsert,noselect
highlight Pmenu ctermfg=153 ctermbg=234
highlight PmenuSel ctermfg=153 ctermbg=240

" ---------- Grepping (https://github.com/junegunn/fzf/blob/master/README-VIM.md)

let g:fzf_layout = { 'window': { 'width': 0.95, 'height': 0.95 } }

function! s:MornGG(pat)
  let l:root = luaeval("vim.lsp.buf.list_workspace_folders()[1] or '.'")
  let l:cmd = "rg --column --line-number --no-heading --color=always -g '!vendor' -g '!.git' -g '*.go' -g '!*_test.go' --smart-case -- ".shellescape(a:pat)." ".shellescape(l:root)
  call fzf#vim#grep(l:cmd, 1, fzf#vim#with_preview({'options': ['--prompt', 'GG> ']}), 0)
endfunction

command! -nargs=* GG call s:MornGG(<q-args>)
command! GM :execute 'lvimgrep /func[^()]*([^()]*\<'.expand('<cword>').'\>)/ '.expand('%:p:h').  '/*' | lopen

" Helpful links:
" https://github.com/junegunn/fzf.vim/blob/master/plugin/fzf.vim
" https://github.com/nanotee/nvim-lua-guide
" https://github.com/neovim/neovim/blob/master/runtime/doc/lsp.txt
" Memo:
" source $MYVIMRC
" :lcl closes it
" Todo:
" group settings by filetype
