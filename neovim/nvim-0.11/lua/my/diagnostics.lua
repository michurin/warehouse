-- https://neovim.io/doc/user/diagnostic.html
-- consider vim.diagnostic.open_float()
vim.keymap.set('n', '<space>ds', function() vim.diagnostic.open_float(); end)

vim.diagnostic.config({
  virtual_text = true,
  severity_sort = true,
  signs = {
    text = {
      [vim.diagnostic.severity.ERROR] = '',
      [vim.diagnostic.severity.WARN] = '',
      [vim.diagnostic.severity.INFO] = '',
      [vim.diagnostic.severity.HINT] = '',
    },
    -- linehl = {},
    numhl = {
      [vim.diagnostic.severity.ERROR] = 'DiagnosticError',
      [vim.diagnostic.severity.WARN] = 'DiagnosticWarn',
      [vim.diagnostic.severity.INFO] = 'DiagnosticInfo',
      [vim.diagnostic.severity.HINT] = 'DiagnosticHint',
    },
  },
-- TODO
--      signs = {
--        text = {
--            [vim.diagnostic.severity.ERROR] = '',
--            [vim.diagnostic.severity.WARN] = '',
--        },
--        linehl = {
--            [vim.diagnostic.severity.ERROR] = 'ErrorMsg',
--        },
--        numhl = {
--            [vim.diagnostic.severity.WARN] = 'WarningMsg',
--        },
--    },
})

-- move to my.colors
vim.api.nvim_set_hl(0, 'DiagnosticError', {fg='#ff0000', bg='#330000'})
vim.api.nvim_set_hl(0, 'DiagnosticWarn', {fg='#ffff00', bg='#330000'})
vim.api.nvim_set_hl(0, 'DiagnosticInfo', {fg='#00ff00', bg='#003300'})
vim.api.nvim_set_hl(0, 'DiagnosticHint', {fg='#00ffff', bg='#003333'})
vim.api.nvim_set_hl(0, 'DiagnosticOk', {fg='#000000', bg='#ffffff'})
-- DiagnosticUnderlineError/Warn/Info/Hint/Ok ?

