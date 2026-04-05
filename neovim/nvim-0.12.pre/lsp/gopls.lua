return {
  cmd = { 'gopls' },
  filetypes = { 'go', 'gomod', 'gowork', 'gotmpl' },
  root_markers = { 'go.mod', '.git' },
  gofumpt = true, -- vim.api.nvim_eval('exists("g:nogofumpt_tweak")') == 0
  experimentalPostfixCompletions = true,
  analyses = { -- https://github.com/golang/tools/blob/master/gopls/doc/analyzers.md
    unusedparams = true,
    unusedwrite = true,
    useany = true,
    shadow = true,
    ST1000 = false,
  },
  staticcheck = true,
  gopls = {
    buildFlags = { '-tags=integration_tests' }, -- TODO get from file? Like gofumpt?
  },
}
