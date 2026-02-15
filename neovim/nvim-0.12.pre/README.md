# Just list of files (gf to open)

```
README.md
after/ftplugin/go.lua
after/ftplugin/markdown.lua
after/ftplugin/qf.lua
after/syntax/README
after/syntax/brief.vim
after/syntax/go.vim
after/syntax/plantuml.vim
bin/vim-helper-open-git
init.lua
lsp/bashls.lua
lsp/gopls.lua
lsp/lua_ls.lua
lsp/protols.lua
lsp/pyright.lua
lsp/ts_ls.lua
lua/configs.lua
lua/functions.lua
```

# Install and run

```
git clone https://github.com/neovim/neovim.git
make CMAKE_INSTALL_PREFIX=$HOME/nvim12 CMAKE_BUILD_TYPE=Release
make install
```

```
NVIM_APPNAME=nvim-12 ~/nvim12/bin/nvim "$@"
```

# Folding

Debug: `:set foldcolumn=5`

# First aid

```
:helpgrep
```

# Ideas

## Manual folding

```
vim.opt_local.foldmethod='manual'
vim.api.nvim_cmd({cmd='fold', range={10, 20}}, {})
```

## Limit file size

```
local function git_files_under_size(max_bytes)
  local cmd = [[
git ls-files |
git cat-file --batch-check='%(objectname) %(objectsize) %(rest)' |
awk '$2 <= ]] .. max_bytes .. [[ {print $3}']]
  return vim.fn.systemlist(cmd)
end
```

## Keys

```
lua =vim.fn.mapcheck('gu')
:echo mapcheck("gu", "n")
```

## Colors

```
:Inspect
```

```
vim.api.nvim_set_hl(0, "QuickFixLine", {bg='#ff0000'})
```

## Links

- <https://vieitesss.github.io/posts/Neovim-new-config/>
- <https://github.com/neovim/nvim-lspconfig>

## Garbage

### Go

```
vim.keymap.set('n', 'gu',
function ()
  vim.lsp.buf.references(nil, { on_list = function (o) vim.fn.setqflist({}, ' ', o) vim.cmd.cfirst() end })
end,
{noremap = true})
```

### Markdown folding (expr-based)

```
function _G.HashFold()
  local line = vim.fn.getline(vim.v.lnum)

  local level = line:match("^(#+)")
  if level then
    return #level
  end

  return "="
end

function _G.SectionFold()
  local line = vim.fn.getline(vim.v.lnum)

  -- заголовок?
  local hashes = line:match("^(#+)%s")
  if not hashes then
    return "="
  end

  local level = #hashes

  -- ищем следующий заголовок
  local next_lnum = vim.v.lnum + 1
  local last = vim.fn.line("$")

  while next_lnum <= last do
    local next_line = vim.fn.getline(next_lnum)
    local next_hashes = next_line:match("^(#+)%s")

    if next_hashes then
      local next_level = #next_hashes
      if next_level <= level then
        return ">"  -- начать фолд здесь
      end
      break
    end

    next_lnum = next_lnum + 1
  end

  return ">"
end

function _G.SimpleSectionFold()
  local line = vim.fn.getline(vim.v.lnum)
  if line:match("^#+%s") then
    return ">"
  end
  return "="
end

function _G.HashSectionFold()
  local line = vim.fn.getline(vim.v.lnum)

  -- заголовок?
  local hashes = line:match("^(#+)%s")
  if not hashes then
    return "="
  end

  return #hashes
end


vim.opt_local.foldmethod = "expr"
vim.opt_local.foldexpr = "v:lua.HashSectionFold()"
vim.opt_local.foldlevel = 99
```
