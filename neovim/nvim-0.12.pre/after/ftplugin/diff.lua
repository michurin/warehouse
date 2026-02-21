vim.opt_local.foldmethod = 'expr'
vim.opt_local.foldexpr = 'v:lua.diff_fold()'

function _G.diff_fold()
  local lnum = vim.v.lnum
  local line = vim.fn.getline(lnum)
  if vim.startswith(line, 'diff ') then return 1 end
  if vim.startswith(line, '--- ') then return 1 end
  if vim.startswith(line, '+++ ') then return 1 end
  if vim.startswith(line, '@@ ') then return 1 end
  if vim.startswith(line, ' ') then return 2 end
  if vim.startswith(line, '-') then return 2 end
  if vim.startswith(line, '+') then return 2 end
  return 0
end

vim.api.nvim_set_hl(0, 'Folded', { bg = '#262626', fg = '#ffffff' })
