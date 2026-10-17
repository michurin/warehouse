local QS = require('qfsort')
local QH = require('qfhl')

local M = {}

-- sort_opts: extra qfsort options (by_location is always on); ...: extra rg args
function M.search(sort_opts, ...)
  sort_opts = vim.tbl_extend('force', { by_location = true }, sort_opts or {})
  local cmdPfx = {
    'rg',
    '--vimgrep',
    '--hidden',
    '--smart-case',
    '--glob', '!.git',
    '--glob', '!node_modules',
    ...,
  }
  return function(opts)
    local pattern = opts.args
    local cmd = vim.tbl_extend('force', {}, cmdPfx)
    cmd[#cmd + 1] = pattern

    local lines = vim.fn.systemlist(cmd) -- TODO pcall for no-executable error
    if vim.v.shell_error ~= 0 then
      print('no files? (rc=' .. tostring(vim.v.shell_error) .. ')')
      return
    end

    local items = vim.fn.getqflist({ lines = lines, efm = '%f:%l:%c:%m' }).items
    QS.sort(items, sort_opts) -- mocks last; rg output order is not stable

    vim.fn.setqflist({}, ' ', {
      title = 'ripgrep search (' .. pattern .. ')',
      items = items,
    })

    vim.cmd('copen')

    QH.set(QH.smartcase(pattern)) -- literally: rg regex is not vim regex
  end
end

M.opts = { nargs = 1 }

return M
