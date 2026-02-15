vim.cmd([[
setlocal foldmethod=expr foldexpr=DiffFold(v:lnum)
function! DiffFold(lnum)
  let line = getline(a:lnum)
  if line =~ '^\(diff\|---\|+++\|@@\) '
    return 1
  elseif line[0] =~ '[-+ ]'
    return 2
  else
    return 0
  endif
endfunction
]])

--[[
vim.opt_local.foldmethod = "expr"
vim.opt_local.foldexpr = "v:lua.DiffFold()"

_G.DiffFold = function()
  local lnum = vim.v.lnum
  local line = vim.fn.getline(lnum)
  if line:match("^(diff|%-%-%-|%+%+%+|@@) ") then return 1 end
  if line:match("^[-+ ]") then return 2 end
  return 0
end
]]--
