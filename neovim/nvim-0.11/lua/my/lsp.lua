vim.lsp.config('gopls', { -- https://github.com/golang/tools/blob/master/gopls/doc/settings.md
  gofumpt = vim.api.nvim_eval('exists("g:nogofumpt_tweak")') == 0, -- true
  experimentalPostfixCompletions = true,
  analyses = { -- https://github.com/golang/tools/blob/master/gopls/doc/analyzers.md
    unusedparams = true,
    unusedwrite = true,
    useany = true,
    shadow = true,
    ST1000 = false,
  },
  staticcheck = true,
})
vim.lsp.enable('gopls')

-- brew install lua-language-server
-- https://luals.github.io/wiki/settings/
vim.lsp.config('lua_ls', {
  diagnostics = {
    disable = {'lowercase-global'},
    globals = {'vim'},
  }
})
vim.lsp.enable('lua_ls')

-- brew install python-lsp-server ?
-- brew install pyright
vim.lsp.config('pyright', {
  analysis={
    useLibraryCodeForTypes = false,
    typeCheckingMode = "off"
  },
  linting = {
    pylintEnabled = true,
    enabled = true
  }
})
vim.lsp.enable('pyright')

vim.lsp.config('ts_ls', {
  cmd = {'typescript-language-server', '--stdio', '--log-level', 'log'},
  single_file_support = true,
})
vim.lsp.enable('ts_ls')

