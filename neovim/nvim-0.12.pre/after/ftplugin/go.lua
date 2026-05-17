vim.opt_local.expandtab = false

vim.keymap.set('i', '.', '.<c-x><c-o>', { noremap = true, buffer = 0 })

vim.opt_local.foldmethod = 'expr'
vim.opt_local.foldexpr = 'v:lua.vim.lsp.foldexpr()'
vim.opt_local.foldtext = 'v:lua.vim.lsp.foldtext()'
vim.opt_local.foldlevel = 99
vim.opt_local.foldlevelstart = 99
vim.opt_local.foldenable = true

vim.keymap.set('n', '<space>cwv', function()
  local line = 'ctx = context.WithValue(ctx, "", nil)'
  local row, _ = unpack(vim.api.nvim_win_get_cursor(0))
  vim.api.nvim_buf_set_lines(0, row, row, false, { line })
  vim.lsp.buf.format({ async = false })
  vim.api.nvim_win_set_cursor(0, { row + 1, 0 })
  vim.cmd.normal({ 'f"f"', bang = true })
  vim.cmd.startinsert()
end)
