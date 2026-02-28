vim.opt_local.expandtab = false

vim.keymap.set('i', '.', '.<c-x><c-o>', { noremap = true, buffer = 0 })

vim.opt_local.foldmethod = "expr"
vim.opt_local.foldexpr = "v:lua.vim.lsp.foldexpr()"
vim.opt_local.foldtext = "v:lua.vim.lsp.foldtext()"
vim.opt_local.foldlevel = 99
vim.opt_local.foldlevelstart = 99
vim.opt_local.foldenable = true
