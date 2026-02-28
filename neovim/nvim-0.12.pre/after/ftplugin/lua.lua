vim.opt_local.expandtab = true
vim.opt_local.shiftwidth = 2
vim.opt_local.tabstop = 2
vim.opt_local.softtabstop = 2

vim.opt_local.foldmethod = 'expr'
vim.opt_local.foldexpr = 'v:lua.vim.lsp.foldexpr()'
vim.opt_local.foldtext = 'v:lua.vim.lsp.foldtext()'
vim.opt_local.foldlevel = 99
vim.opt_local.foldlevelstart = 99
vim.opt_local.foldenable = true
