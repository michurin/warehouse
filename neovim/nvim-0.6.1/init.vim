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
call plug#end()

lua <<TELESCOPE_HELPERS
function ft_args()
  local bufnr = vim.api.nvim_get_current_buf()
  local filetype = vim.bo[bufnr].filetype
  return ({ -- rg options
    ["go"]={
      "-g", "*.go",
      "-g", "!vendor",
      "-g", "!mock",
      "-g", "!mocks",
      "-g", "!.git",
      "-g", "!*_test.go",
    },
    ["js"]={"-g", "*.js"},
  })[filetype]
end
function ft_alt_args()
  local bufnr = vim.api.nvim_get_current_buf()
  local filetype = vim.bo[bufnr].filetype
  return ({ -- rg options
    ["go"]={
      "-g", "*_test.go",
      "-g", "!vendor",
      "-g", "!mock",
      "-g", "!mocks",
      "-g", "!.git",
    },
  })[filetype]
end
TELESCOPE_HELPERS

" std
nnoremap <space>ff <cmd>lua require('telescope.builtin').find_files()<cr>
nnoremap <space>fg <cmd>lua require('telescope.builtin').live_grep()<cr>
nnoremap <space>fb <cmd>lua require('telescope.builtin').buffers()<cr>
nnoremap <space>fh <cmd>lua require('telescope.builtin').help_tags()<cr>
" like /
nnoremap <space>f/ <cmd>lua require('telescope.builtin').current_buffer_fuzzy_find()<cr>
" z=
nnoremap z= <cmd>lua require('telescope.builtin').spell_suggest()<cr>
" m — method, d — diagnostics, l — language, t — tests
nnoremap <space>fm <cmd>lua require('telescope.builtin').grep_string()<cr>
nnoremap <space>fd <cmd>lua require('telescope.builtin').diagnostics()<cr>
nnoremap <space>fl <cmd>lua require('telescope.builtin').live_grep({additional_args=ft_args, wrap_results=true})<cr>
nnoremap <space>ft <cmd>lua require('telescope.builtin').live_grep({additional_args=ft_alt_args, wrap_results=true})<cr>
" resume
nnoremap <space>fr <cmd>lua require('telescope.builtin').resume()<cr>
" std LSP
nnoremap gr <cmd>lua require('telescope.builtin').lsp_references()<cr>
nnoremap gi <cmd>lua require('telescope.builtin').lsp_implementations()<cr>
nnoremap gd <cmd>lua require('telescope.builtin').lsp_definitions()<cr>
nnoremap gy <cmd>lua require('telescope.builtin').lsp_type_definitions()<cr>

lua <<TELESCOPE_SETTINGS
local action_layout = require("telescope.actions.layout")
require('telescope').setup{
  defaults = {
    layout_strategy = 'vertical',
    layout_config = {
      height = 0.9,
      width = 0.9,
      preview_cutoff = 3,
    },
    mappings = {
--      n = {
--        ["<C-t>"] = action_layout.toggle_preview,
--      },
--      i = {
--        ["<C-t>"] = action_layout.toggle_preview,
--      },
    },
  },
}
TELESCOPE_SETTINGS

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
  mapping = {
    ['<C-b>'] = cmp.mapping(cmp.mapping.scroll_docs(-4), { 'i', 'c' }),
    ['<C-f>'] = cmp.mapping(cmp.mapping.scroll_docs(4), { 'i', 'c' }),
    ['<C-Space>'] = cmp.mapping(cmp.mapping.complete(), { 'i', 'c' }),
--    ['<C-y>'] = cmp.config.disable, -- Specify `cmp.config.disable` if you want to remove the default `<C-y>` mapping.
    ['<C-e>'] = cmp.mapping({
      i = cmp.mapping.abort(),
      c = cmp.mapping.close(),
    }),
    ['<CR>'] = cmp.mapping.confirm({ select = true }), -- Accept currently selected item. Set `select` to `false` to only confirm explicitly selected items.
  },
  sources = cmp.config.sources({
    { name = 'nvim_lsp' },
    { name = 'vsnip' },
  }, {
    { name = 'buffer' },
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
  sources = {
    { name = 'buffer' },
  }
})

-- Use cmdline & path source for ':' (if you enabled `native_menu`, this won't work anymore).
cmp.setup.cmdline(':', {
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
  buf_set_keymap('n', '<space>e', '<cmd>lua vim.lsp.diagnostic.show_line_diagnostics()<CR>', opts)
  buf_set_keymap('n', '[d', '<cmd>lua vim.lsp.diagnostic.goto_prev()<CR>', opts)
  buf_set_keymap('n', ']d', '<cmd>lua vim.lsp.diagnostic.goto_next()<CR>', opts)
  buf_set_keymap('n', '<space>q', '<cmd>lua vim.lsp.diagnostic.set_loclist()<CR>', opts)
  buf_set_keymap('n', '<space>f', '<cmd>lua vim.lsp.buf.formatting()<CR>', opts)
end

local servers = {'gopls', 'intelephense', 'pyright', 'tsserver'}
local capabilities = require('cmp_nvim_lsp').update_capabilities(vim.lsp.protocol.make_client_capabilities())
for _, lsp in pairs(servers) do
  require('lspconfig')[lsp].setup {
    on_attach = on_attach,
    capabilities = capabilities,
    flags = {
      debounce_text_changes = 150, -- This will be the default in neovim 0.7+
    },
    settings={
      gopls = { -- https://github.com/golang/tools/blob/master/gopls/doc/settings.md
        gofumpt = vim.api.nvim_eval('exists("g:nogofumpt_tweak")') == 0 -- true
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
    },
  }
end
LSP_AND_COMPLETION_SETTINTS

set completeopt=menu,menuone,noselect

highlight Pmenu ctermfg=153 ctermbg=234
highlight PmenuSel ctermfg=153 ctermbg=240

highlight LspDiagnosticsDefaultHint ctermfg=64 ctermbg=234
highlight LspDiagnosticsDefaultInformation ctermfg=31 ctermbg=234
highlight LspDiagnosticsDefaultWarning ctermfg=137 ctermbg=234
highlight LspDiagnosticsDefaultError ctermfg=124 ctermbg=234

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
        ["]m"] = "@function.outer",
        ["]]"] = "@class.outer",
      },
      goto_next_end = {
        ["]M"] = "@function.outer",
        ["]["] = "@class.outer",
      },
      goto_previous_start = {
        ["[m"] = "@function.outer",
        ["[["] = "@class.outer",
      },
      goto_previous_end = {
        ["[M"] = "@function.outer",
        ["[]"] = "@class.outer",
      },
    },
  },
}
TREESITTER_SETTINGS

" Telescope
" https://github.com/nvim-telescope/telescope.nvim
" https://github.com/nvim-telescope/telescope.nvim/wiki/Configuration-Recipes
" PLS
" - npm config set prefix "${HOME}/.npm-packages"
" - npm install -g typescript typescript-language-server
" Completion
" https://github.com/hrsh7th/nvim-cmp

" VERY COMMON SETTINGS

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
set list lcs=trail:+,tab:▹·
set fillchars=fold:\ " (space)
set foldmethod=syntax
set foldlevelstart=99
set foldlevel=99
set synmaxcol=10000

highlight Whitespace term=none cterm=none ctermfg=DarkGray gui=none guifg=DarkGray ctermbg=none
highlight EndOfBuffer term=none cterm=none ctermfg=DarkGray gui=none guifg=DarkGray ctermbg=none
highlight LineNr ctermfg=grey

autocmd TextYankPost * lua vim.highlight.on_yank {higroup="hlTextYankPost", timeout=400}
highlight link hlTextYankPost Visual

nnoremap <silent> ]c :cnext<CR>
nnoremap <silent> [c :cprevious<CR>
nnoremap <silent> ]l :lnext<CR>
nnoremap <silent> [l :lprevious<CR>
nnoremap <silent> ]b :bnext<CR>
nnoremap <silent> [b :bprevious<CR>

autocmd FileType qf setlocal nobuflisted

" SPELLING

setlocal spell spelllang=en_us,ru_yo spelloptions=camel spellcapcheck=

syntax match UrlNoSpell 'https\?:\/\/[^[:space:]]\+' contains=@NoSpell

highlight SpellBad term=none cterm=underline ctermfg=none gui=bold guifg=none ctermbg=none
highlight SpellCap term=none cterm=underline ctermfg=none gui=bold guifg=none ctermbg=none
highlight SpellRare term=none cterm=underline ctermfg=none gui=bold guifg=none ctermbg=none
highlight SpellLocal term=none cterm=underline ctermfg=none gui=bold guifg=none ctermbg=none

" FOLDING

function! XFoldText()
    let line = getline(v:foldstart)
    let folded_line_num = v:foldend - v:foldstart
    let line_text = substitute(line, '^{\+[0-9]\+\s\+', '', 'g') " hackish way to remove default marker
    let fillcharcount = winwidth(0) - len(line_text) - len(folded_line_num) - 11
    return '⟫ ' . line_text . ' ' . repeat('╶', fillcharcount) . ' (' . folded_line_num . ')'
endfunction
set foldtext=XFoldText()
highlight Folded ctermfg=155 ctermbg=235
" it is useful modeline: vi:fdm=marker:foldlevel=0

" GO STUFF

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

" https://github.com/nvim-treesitter/nvim-treesitter-textobjects
" map <silent> [[ :noh<CR>?^func\><CR>:let @/=''<CR>:set hls<CR>
" map <silent> ]] :noh<CR>/^func\><CR>:let @/=''<CR>:set hls<CR>

autocmd FileType go autocmd BufWritePre *.go lua goimports(1000)

" NON-GO

autocmd FileType javascript autocmd BufWritePre *.js lua vim.lsp.buf.formatting_sync(nil, 500)

autocmd FileType go     setlocal noexpandtab
autocmd FileType vim    setlocal shiftwidth=2 tabstop=2 softtabstop=2
autocmd FileType lua    setlocal shiftwidth=2 tabstop=2 softtabstop=2
autocmd FileType css    setlocal shiftwidth=2 tabstop=8 softtabstop=2
autocmd FileType html   setlocal shiftwidth=2 tabstop=8 softtabstop=2
autocmd FileType javascript setlocal shiftwidth=2 tabstop=8 softtabstop=2
autocmd FileType json   setlocal shiftwidth=2 tabstop=8 softtabstop=2
autocmd FileType yaml   setlocal shiftwidth=2 tabstop=2 softtabstop=2 foldmethod=indent
autocmd FileType make   setlocal tabstop=8
autocmd FileType tcl    setlocal shiftwidth=2 tabstop=8 softtabstop=2 foldmethod=indent
