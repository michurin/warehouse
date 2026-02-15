-- IDEAS:
-- ALL BINDINGS
-- lua =vim.api.nvim_get_current_buf()
-- lua vim.api.nvim_buf_set_lines(0, 0, -1, false, vim.split(vim.inspect(vim.fn.maplist()), "\n"))
--
-- search word under cursor:
-- vim.cmd('...')
-- noremap eee :execute "vimgrep /" . expand('<cword>') . "/j " . " `git ls-files`"<CR><bar> <C-W><C-W>
--
-- set search reg:
-- =vim.fn.setreg("/", "Id")

vim.api.nvim_win_set_height(0, vim.fn.min({5, vim.api.nvim_buf_line_count(0)}))

--vim.keymap.set('n', '<Enter>', function()
--  vim.cmd.cclose() -- vim.api.nvim_win_close(0, false)
--end, {noremap=true, buffer=true})

-- vim.keymap.set('n', '<Enter>', function()
  -- vim.api.nvim_win_close(0, false)
  -- vim.cmd.cclose()
-- end, {noremap=true})

-- nnoremap <buffer> <Down> <Down><CR><C-w>p
-- vim.keymap.set('n', '<S-Down>', '<Down><CR><C-w>p', {noremap=true, buffer=true})
vim.keymap.set('n', '<S-j>', function ()
  local count = #vim.fn.getqflist()
  if count == 0 then return end
  local idx = vim.fn.getqflist({ idx = 0 }).idx
  if idx == count then vim.cmd('cfirst') else vim.cmd('cnext') end
  vim.api.nvim_feedkeys('<CR>', 'n', false)
  vim.cmd('wincmd p')
end, {noremap=true, buffer=true})

-- nnoremap <buffer> <Up> <Up><CR><C-w>p
-- vim.keymap.set('n', '<S-Up>', '<Up><CR><C-w>p', {noremap=true, buffer=true})
vim.keymap.set('n', '<S-k>', function ()
  local count = #vim.fn.getqflist()
  if count == 0 then return end
  local idx = vim.fn.getqflist({ idx = 0 }).idx
  if idx == 1 then vim.cmd('clast') else vim.cmd('cprev') end
  vim.api.nvim_feedkeys('<CR>', 'n', false)
  vim.cmd('wincmd p')
end, {noremap=true, buffer=true})
