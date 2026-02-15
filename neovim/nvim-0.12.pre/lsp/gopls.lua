return {
  cmd = { 'gopls' },
  filetypes = { 'go', 'gomod' },
  root_markers = { 'go.mod', '.git' },
  gofumpt = true,
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
