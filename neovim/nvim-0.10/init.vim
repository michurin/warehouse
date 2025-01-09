call plug#begin() " https://github.com/junegunn/vim-plug +PlugInstall
  Plug 'neovim/nvim-lspconfig'
  Plug 'nvim-lua/plenary.nvim'
  Plug 'nvim-telescope/telescope.nvim'
  Plug 'hrsh7th/cmp-nvim-lsp'
  Plug 'hrsh7th/cmp-buffer'
  Plug 'hrsh7th/cmp-path'
  Plug 'hrsh7th/cmp-cmdline'
  Plug 'hrsh7th/nvim-cmp'
  Plug 'hrsh7th/cmp-vsnip'
  Plug 'hrsh7th/vim-vsnip'
  Plug 'nvim-treesitter/nvim-treesitter' ", {'do': ':TSUpdate'} " TSInstall lang
  Plug 'nvim-treesitter/nvim-treesitter-textobjects'
  Plug 'nvim-treesitter/nvim-treesitter-context'
  Plug 'nvim-treesitter/nvim-treesitter-refactor'

  Plug 'mfussenegger/nvim-dap'
  Plug 'nvim-neotest/nvim-nio' " required by nvim-dap-ui
  Plug 'rcarriga/nvim-dap-ui'
  Plug 'leoluz/nvim-dap-go' " go install github.com/go-delve/delve/cmd/dlv@latest
  Plug 'nvim-telescope/telescope-dap.nvim' " require('telescope').load_extension('dap') after require('telescope').setup()

  Plug 'codota/tabnine-nvim', { 'do': './dl_binaries.sh' }
call plug#end()

if filereadable(getcwd() . "/.nogofumpt") " Oh, too hackish. vim.lsp.buf.list_workspace_folders() or util.root_pattern?
  let g:nogofumpt_tweak = 1
endif

lua <<LSP_AND_COMPLETION_SETTINTS
-- Completion

local cmp = require'cmp'

cmp.setup({
  snippet = { -- must be specified
    expand = function(args)
      vim.fn["vsnip#anonymous"](args.body)
    end,
  },
  view = {
    docs = {
      auto_open = true,
    },
  },
  window = {
    completion = cmp.config.window.bordered({
      winhighlight = 'CmpItemKindDefault:Comment,CmpItemAbbrDefault:Normal,CmpItemAbbrMatchDefault:Search',
    }),
    documentation = cmp.config.window.bordered(),
  },
  mapping = cmp.mapping.preset.insert({
    ['<C-b>'] = cmp.mapping(cmp.mapping.scroll_docs(-4), { 'i', 'c' }),
    ['<C-f>'] = cmp.mapping(cmp.mapping.scroll_docs(4), { 'i', 'c' }),
    ['<C-Space>'] = cmp.mapping(cmp.mapping.complete(), { 'i', 'c' }),
--    ['<C-y>'] = cmp.config.disable, -- Specify `cmp.config.disable` if you want to remove the default `<C-y>` mapping.
    ['<C-e>'] = cmp.mapping({
      i = cmp.mapping.abort(),
      c = cmp.mapping.close(),
    }),
    ['<CR>'] = cmp.mapping.confirm({ select = true }), -- Accept currently selected item. Set `select` to `false` to only confirm explicitly selected items.
  }),
  sources = cmp.config.sources({
    { name = 'nvim_lsp', group_index = 1 },
    { name = 'vsnip', group_index = 2 },
    { name = 'buffer', group_index = 3 },
  })
})

-- Set configuration for specific filetype.
cmp.setup.filetype('gitcommit', {
  sources = cmp.config.sources({
    { name = 'cmp_git' }, -- You can specify the `cmp_git` source if you were installed it.
  }, {
    { name = 'buffer' },
  })
})

-- Use buffer source for `/` (if you enabled `native_menu`, this won't work anymore).
cmp.setup.cmdline('/', {
  mapping = cmp.mapping.preset.cmdline(),
  sources = {
    { name = 'buffer' },
  }
})

-- Use cmdline & path source for ':' (if you enabled `native_menu`, this won't work anymore).
cmp.setup.cmdline(':', {
  mapping = cmp.mapping.preset.cmdline(),
  sources = cmp.config.sources({
    { name = 'path' },
  }, {
    { name = 'cmdline' },
  })
})

-- Lsp
local on_attach = function(client, bufnr)
  local function buf_set_keymap(...) vim.api.nvim_buf_set_keymap(bufnr, ...) end
  local function buf_set_option(...) vim.api.nvim_buf_set_option(bufnr, ...) end

  -- Enable completion triggered by <c-x><c-o>
  buf_set_option('omnifunc', 'v:lua.vim.lsp.omnifunc')

  -- Mappings
  local opts = { noremap=true, silent=true }

  -- See `:help vim.lsp.*` for documentation on any of the below functions
  buf_set_keymap('n', 'gD', '<cmd>lua vim.lsp.buf.declaration()<CR>', opts)
--  buf_set_keymap('n', 'gd', '<cmd>lua vim.lsp.buf.definition()<CR>', opts)
  buf_set_keymap('n', 'K', '<cmd>lua vim.lsp.buf.hover()<CR>', opts)
--  buf_set_keymap('n', 'gi', '<cmd>lua vim.lsp.buf.implementation()<CR>', opts)
  buf_set_keymap('n', '<C-k>', '<cmd>lua vim.lsp.buf.signature_help()<CR>', opts)
  buf_set_keymap('n', '<space>wa', '<cmd>lua vim.lsp.buf.add_workspace_folder()<CR>', opts)
  buf_set_keymap('n', '<space>wr', '<cmd>lua vim.lsp.buf.remove_workspace_folder()<CR>', opts)
  buf_set_keymap('n', '<space>wl', '<cmd>lua print(vim.inspect(vim.lsp.buf.list_workspace_folders()))<CR>', opts)
--  buf_set_keymap('n', 'gy', '<cmd>lua vim.lsp.buf.type_definition()<CR>', opts)
  buf_set_keymap('n', '<space>rn', '<cmd>lua vim.lsp.buf.rename()<CR>', opts)
  buf_set_keymap('n', '<space>ca', '<cmd>lua vim.lsp.buf.code_action()<CR>', opts)
--  buf_set_keymap('n', 'gr', '<cmd>lua vim.lsp.buf.references()<CR>', opts)
  buf_set_keymap('n', '[d', '<cmd>lua vim.diagnostic.goto_prev()<CR>', opts)
  buf_set_keymap('n', ']d', '<cmd>lua vim.diagnostic.goto_next()<CR>', opts)
  buf_set_keymap('n', '<space>q', '<cmd>lua vim.lsp.diagnostic.set_loclist()<CR>', opts)
  buf_set_keymap('n', '<space>f', '<cmd>lua vim.lsp.buf.format()<CR>', opts)
-- deprecated vim.lsp.diagnostic.show_line_diagnostics()
  buf_set_keymap('n', '<space>e', '<cmd>lua vim.diagnostic.open_float({source="if_many"})<CR>', opts)
end

local servers = {'gopls', 'intelephense', 'pyright', 'ts_ls', 'ccls', 'lua_ls'}
local capabilities = require('cmp_nvim_lsp').default_capabilities(vim.lsp.protocol.make_client_capabilities())
for _, lsp in pairs(servers) do
  require('lspconfig')[lsp].setup {
    on_attach = on_attach,
    capabilities = capabilities,
    flags = {
      debounce_text_changes = 150, -- This will be the default in neovim 0.7+
    },
    settings={
      gopls = { -- https://github.com/golang/tools/blob/master/gopls/doc/settings.md
        gofumpt = vim.api.nvim_eval('exists("g:nogofumpt_tweak")') == 0, -- true
        experimentalPostfixCompletions = true,
        analyses = {
          unusedparams = true,
          unusedwrite = true,
          useany = true,
          shadow = true,
        },
        staticcheck = true,
      },
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
      Lua = {
        diagnostics = {
          globals = {'vim'},
        }
      },
    },
  }
end
LSP_AND_COMPLETION_SETTINTS

set completeopt=menu,menuone,noselect

lua <<TREESITTER_CONTEXT
require'treesitter-context'.setup{
  enable = true, -- autocmd VimEnter * TSContextEnable
  throttle = true, -- may improve performance
  max_lines = 10, -- 0 — no limit
  mode = 'cursor', -- 'topline', 'cursor',
  -- separator = '┄',
  patterns = { -- lua print(vim.inspect(require'nvim-treesitter.ts_utils'.get_node_at_cursor():type()))
    default = {'class', 'function', 'method', 'for', 'while', 'if', 'switch', 'case'},
    yaml = {'block_mapping_pair'},
    json = {'pair'},
    toml = {'table', 'bare_key'},
    markdown = {'section'},
    go = {'import_declaration', 'assignment_statement', 'short_var_declaration', 'defer_statement', 'func_literal'}, -- func_literal for anonymous functions
  },
}
TREESITTER_CONTEXT

lua <<TREESITTER_SETTINGS
require'nvim-treesitter.configs'.setup {
  textobjects = {
    select = {
      enable = true,

      -- Automatically jump forward to textobj, similar to targets.vim
      lookahead = true,

      keymaps = {
        -- You can use the capture groups defined in textobjects.scm
        ["af"] = "@function.outer",
        ["if"] = "@function.inner",
        ["ac"] = "@class.outer",
        ["ic"] = "@class.inner",
      },
    },
    swap = {
      enable = true,
      swap_next = {
        ["<space>a"] = "@parameter.inner",
      },
      swap_previous = {
        ["<space>A"] = "@parameter.inner",
      },
    },
    move = {
      enable = true,
      set_jumps = true, -- whether to set jumps in the jumplist
      goto_next_start = {
        ["]m"] = "@class.outer",
        ["]]"] = "@function.outer",
      },
      goto_next_end = {
        ["]M"] = "@class.outer",
        ["]["] = "@function.outer",
      },
      goto_previous_start = {
        ["[m"] = "@class.outer",
        ["[["] = "@function.outer",
      },
      goto_previous_end = {
        ["[M"] = "@class.outer",
        ["[]"] = "@function.outer",
      },
    },
  },
  refactor = {
    highlight_definitions = {
      enable = false,
      -- Set to false if you have an `updatetime` of ~100.
      clear_on_cursor_move = false,
    },
    highlight_current_scope = {
      -- Works wired for me, even with custom TSCurrentScope
      enable = false,
    },
    smart_rename = {
      enable = true,
      keymaps = {
        smart_rename = "gR",
      },
    },
    navigation = {
      enable = true,
      keymaps = {
        -- goto_definition = "gnd",
        -- list_definitions = "gnD",
        -- list_definitions_toc = "gO",
        goto_next_usage = "]u",
        goto_previous_usage = "[u",
      },
    },
  },
}
TREESITTER_SETTINGS

" VERY COMMON SETTINGS

set wildmode=list:full,longest
set statusline=%<%f\ %h%m%r%=%-10.(%l,%v%)\ %8.(%B%)
set nofixendofline
set scrolloff=4
set number " relativenumber " /, C-G, C-T instead
set title
set hidden
set shiftwidth=4
set tabstop=4
set softtabstop=4
set expandtab
set autoindent
set list
set listchars=trail:+,tab:▹·,nbsp:␣,extends:▶,precedes:◀
set whichwrap+=<,>,[,]
set fillchars=fold:\ " (space)
" set foldmethod=syntax
set foldlevelstart=99
set foldlevel=99
set synmaxcol=10000

set foldmethod=expr
set foldexpr=nvim_treesitter#foldexpr()

set splitright
set isfname-=# " TODO: do it for YAML only?

set guicursor=n-c-sm:block,i-ci-ve:ver25,r-cr-o-v:hor20

let g:netrw_winsize = 30
let g:netrw_banner = 0
let g:netrw_keepdir = 0

function! ToggleNetrw()
  let f = 1
  let i = bufnr("$")
  while (i >= 1)
    if (getbufvar(i, "&filetype") == "netrw")
      silent exe "bwipeout " . i
      let f = 0
    endif
    let i-=1
  endwhile
  if f
    let g:NetrwIsOpen=1
    silent Lexplore %:p:h
  endif
endfunction

nnoremap <space>dd :call ToggleNetrw()<CR>

nnoremap <silent> ]c :cnext<CR>
nnoremap <silent> [c :cprevious<CR>
nnoremap <silent> ]l :lnext<CR>
nnoremap <silent> [l :lprevious<CR>
nnoremap <silent> ]b :bnext<CR>
nnoremap <silent> [b :bprevious<CR>

autocmd FileType qf setlocal nobuflisted

" SPELLING

set spell spelllang=en_us,ru_yo spelloptions=camel " spellcapcheck=

syntax match UrlNoSpell 'https\?:\/\/[^[:space:]]\+' contains=@NoSpell

" FOLDING

function! XFoldText()
    let line = getline(v:foldstart)
    let folded_line_num = v:foldend - v:foldstart
    let line_text = substitute(substitute(line, '^{\+[0-9]\+\s\+', '', 'g'), '^  ', '', 'g') " hackish way to remove two spaces and default marker
    let fillcharcount = winwidth(0) - strchars(line_text) - len(folded_line_num) - 11
    return '⟫ ' . line_text . ' ' . repeat('╶', fillcharcount) . ' (' . folded_line_num . ')'
endfunction
set foldtext=XFoldText()

" GO STUFF

autocmd bufenter *.go syntax keyword goTodo contained ACHTUNG TODO FIXME XXX BUG todo fixme xxx bug

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

function! s:GoLint()
  cexpr system('golangci-lint run')
endfunction

command! GA call s:GoAlt('e')
command! GAA call s:GoAlt('bo vs')
command! GL call s:GoLint()
" alias pbcopy='xclip -selection clipboard'
" alias pbpaste='xclip -selection clipboard -o'
" alias pbcopy='xsel --clipboard --input'
" alias pbpaste='xsel --clipboard --output'
command! GP :lgetexpr system("pbpaste | sed -n '/^[[:space:]]/ {s/^[[:space:]]*//; s/\\(:[0-9][0-9]*\\)/\\1:>/; p;}'") | lopen " Uh. Ugly

" https://github.com/neovim/nvim-lspconfig/issues/115#issuecomment-1128949874
lua <<EOF
function org_imports()
  local clients = vim.lsp.buf_get_clients()
  for _, client in pairs(clients) do
    local params = vim.lsp.util.make_range_params(nil, client.offset_encoding)
    params.context = {only = {"source.organizeImports"}}
    local result = vim.lsp.buf_request_sync(0, "textDocument/codeAction", params, 5000)
    for _, res in pairs(result or {}) do
      for _, r in pairs(res.result or {}) do
        if r.edit then
          vim.lsp.util.apply_workspace_edit(r.edit, client.offset_encoding)
        else
          vim.lsp.buf.execute_command(r.command)
        end
      end
    end
  end
end
EOF

" looking nice alternative, however won't work in some cases
" autocmd FileType go autocmd BufWritePre *.go lua vim.lsp.buf.code_action({ context = { only = { "source.organizeImports" } }, apply = true })

autocmd FileType go autocmd BufWritePre *.go lua vim.lsp.buf.format({async=false})
autocmd FileType go autocmd BufWritePre *.go lua org_imports()

" NON-GO

autocmd BufRead * let &l:modifiable = !&readonly " Blocking changes to read only files

autocmd FileType javascript autocmd BufWritePre *.js lua vim.lsp.buf.format()

autocmd FileType go     setlocal noexpandtab
autocmd FileType vim    setlocal shiftwidth=2 tabstop=2 softtabstop=2
autocmd FileType markdown setlocal shiftwidth=2 tabstop=2 softtabstop=2
autocmd FileType lua    setlocal shiftwidth=2 tabstop=2 softtabstop=2
autocmd FileType css    setlocal shiftwidth=2 tabstop=8 softtabstop=2
autocmd FileType html   setlocal shiftwidth=2 tabstop=8 softtabstop=2
autocmd FileType javascript setlocal shiftwidth=2 tabstop=8 softtabstop=2
autocmd FileType json   setlocal shiftwidth=2 tabstop=8 softtabstop=2
autocmd FileType xslt   setlocal shiftwidth=2 tabstop=8 softtabstop=2
autocmd FileType yaml   setlocal shiftwidth=2 tabstop=2 softtabstop=2 foldmethod=indent
autocmd FileType make   setlocal tabstop=8
autocmd FileType tcl    setlocal shiftwidth=2 tabstop=8 softtabstop=2 foldmethod=indent
autocmd FileType perl   setlocal noexpandtab
autocmd FileType gitconfig setlocal noexpandtab

" force json. Slightly ugly
command! JSON :setlocal syntax=json foldmethod=syntax buftype=nofile | :echo ''

" Git related hacks

" ugly, however it make you free to search over raw git blame output
function! GitBlame()
  let l:l = line('.')
  execute 'new | 0read ! git blame ' expand('%')
  set buftype=nofile
  set bufhidden=wipe
  execute l:l
endfunction
command! GitBlame :call GitBlame()

" Language related hacks

function! EchoWarning(msg)
  echohl WarningMsg
  echo "Warning"
  echohl None
  echon ': ' a:msg
endfunction
nnoremap Ж :call EchoWarning('RU')<CR>
set langmap=ФИСВУАПРШОЛДЬТЩЗЙКЫЕГМЦЧНЯ;ABCDEFGHIJKLMNOPQRSTUVWXYZ,фисвуапршолдьтщзйкыегмцчня;abcdefghijklmnopqrstuvwxyz

" Hacks

"function SyntaxItemHack()
"  return synIDattr(synID(line("."),col("."),1),"name")
"endfunction
"set statusline=%<%f\ %h%m%r\ %{&filetype}%=%{SyntaxItemHack()}\ %-10.(%l,%v%)\ %8.(%B%)

lua require('my.runsh')
lua require('my.colors')
lua require('my.tabnine')
lua require('my.telescope')
lua require('my.dap')

lua vim.keymap.set('n', '<space>wee', function() vim.system({'av', 'openstash', vim.api.nvim_buf_get_name(0), vim.api.nvim_win_get_cursor(0)[1]}, {stdout=false, stderr=false, tmeout=5000}) end)
