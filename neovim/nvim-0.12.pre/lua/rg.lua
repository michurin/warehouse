local M = {}

function M.search(...)
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

    vim.fn.setqflist({}, ' ', {
      title = 'ripgrep search (' .. pattern .. ')',
      lines = lines,
      efm = '%f:%l:%c:%m',
    })

    vim.cmd('copen')
  end
end

M.opts = { nargs = 1 }

return M
