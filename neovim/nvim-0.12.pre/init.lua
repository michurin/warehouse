require('configs')
local F = require('functions')

--

vim.lsp.enable({
  'gopls',
  'lua_ls',
  'bashls',
  'pyright',
  'ts_ls',
  'protols',
})

vim.diagnostic.config({ virtual_text = true })

vim.keymap.set('n', 'gd', function() vim.lsp.buf.definition() end, { noremap = true })

vim.api.nvim_create_autocmd('BufWritePre', { callback = function() vim.lsp.buf.format() end })

--

vim.keymap.set('n', '<space>fl', F.grep_in_files(F.files_from_cmd('git ls-files "*.go" ":!*_test.go"'), F.cword),
  { noremap = true })
vim.keymap.set('v', '<space>fl',
  F.grep_in_files(F.files_from_cmd('git ls-files "*.go" ":!*_test.go"'), F.visual_scalar),
  { noremap = true })
vim.keymap.set('n', '<space>fL',
  F.grep_in_files(F.files_from_cmd('git ls-files "*.go" ":!*_test.go"'), F.input('PATTERN>')), { noremap = true })
vim.keymap.set('n', '<space>ft', F.grep_in_files(F.files_from_cmd('git ls-files "*_test.go"'), F.cword),
  { noremap = true })
vim.keymap.set('v', '<space>ft',
  F.grep_in_files(F.files_from_cmd('git ls-files "*_test.go"'), F.visual_scalar),
  { noremap = true })
vim.keymap.set('n', '<space>fT', F.grep_in_files(F.files_from_cmd('git ls-files "*_test.go"'), F.input('PATTERN>')),
  { noremap = true })
vim.keymap.set('n', '<space>fb', F.grep_in_files(F.files_buffers, F.cword), { noremap = true })
vim.keymap.set('v', '<space>fb', F.grep_in_files(F.files_buffers, F.visual_scalar), { noremap = true })
vim.keymap.set('n', '<space>fB', F.grep_in_files(F.files_buffers, F.input('PATTERN>')), { noremap = true })
vim.keymap.set('n', '<space>fa', F.grep_in_files(F.files_from_cmd('find . -type f'), F.cword), { noremap = true })
vim.keymap.set('v', '<space>fa', F.grep_in_files(F.files_from_cmd('find . -type f'), F.visual_scalar),
  { noremap = true })
vim.keymap.set('n', '<space>fA', F.grep_in_files(F.files_from_cmd('find . -type f'), F.input('PATTERN>')),
  { noremap = true })
vim.keymap.set('n', '<space>fd', F.grep_in_files(F.files_from_find_cdir, F.cword), { noremap = true })
vim.keymap.set('v', '<space>fd', F.grep_in_files(F.files_from_find_cdir, F.visual_scalar),
  { noremap = true })
vim.keymap.set('n', '<space>fD', F.grep_in_files(F.files_from_find_cdir, F.input('PATTERN>')), { noremap = true })

--

vim.keymap.set('n', '<space>ff', F.show_files, { noremap = true })
vim.keymap.set('n', '<space>fF', F.show_files_by_pattern, { noremap = true })

-- TODO move to go.lua? or to lsp settings?
-- outgoing_calls looks buggy
vim.keymap.set('n', '<space>fi', function()
  local view = vim.fn.winsaveview()
  vim.cmd('normal [[0lllll')
  vim.lsp.buf.incoming_calls()
  vim.fn.winrestview(view)
end, { noremap = true })
--

vim.api.nvim_create_autocmd(F.qf_buffers_events, { callback = F.qf_buffers_handler })
vim.api.nvim_create_user_command('BUF', F.qf_buffers, {})
vim.api.nvim_create_user_command('E', F.smart_file_locate, { nargs = 1 })
vim.api.nvim_create_user_command('D', F.exec_git_diff_all, { nargs = 1 })
vim.api.nvim_create_user_command('CC', function() vim.opt.colorcolumn = { 120 } end, {})

--

vim.keymap.set('n', '<space>www', F.exec(F.paragraph_block), { noremap = true })
vim.keymap.set('v', '<space>www', F.exec(F.visual_text), { noremap = true })

vim.keymap.set('n', '<space>wwd', F.exec_git_diff, { noremap = true })
vim.keymap.set('n', '<space>wwg', F.exec_git_blame, { noremap = true })
vim.keymap.set('n', '<space>wwc', F.exec_command, { noremap = true })

vim.keymap.set('n', '<space>wee', function()
  vim.system(
    {
      vim.fn.stdpath('config') .. '/bin/vim-helper-open-git',
      vim.api.nvim_buf_get_name(0),
      tostring(vim.api.nvim_win_get_cursor(0)[1]),
    },
    {
      stdout = false,
      stderr = false,
      tmeout = 5000,
    }
  )
end, { noremap = true })

vim.keymap.set('n', '<space>bb', F.show_keys, { noremap = true })

vim.keymap.set("n", "<leader>hi", function() vim.cmd("Inspect") end, { noremap = true })
