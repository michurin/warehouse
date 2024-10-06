--[[

# example

echo '{"home": "$HOME", "val": "'$HOME'"}' | jq
echo ok2
x="OK"
echo $x

#

separators
- empty string
- prefix
  - #
  - --
  - {{{

bindings
- <space>www — run
- w — toggle text wrapping
- q — close
- <esc><esc> — close too

affected registers
- e

# ideas

set fde=getline(v:lnum)=~'^\\s*$'&&getline(v:lnum+1)=~'\\S'?'<1':1
set fdm=expr

--]]


--[[
vim.keymap.set('v', '<space>www', function()
  local vstart = vim.fn.getpos("'<")
  local vend = vim.fn.getpos("'>")
  local lstart = vstart[2]
  local lend = vend[2]
  local command = vim.fn.getline(lstart, lend)
  print(command)
end)
]]

vim.keymap.set('n', '<space>www', function()
  local check_stop_line = function(r)
    local ln = vim.api.nvim_buf_get_lines(0, r-1, r, false)[1]
    return not ln or ln == "" -- or string.sub(ln, 1, 1) == '#' or string.sub(ln, 1, 2) == '--' or string.sub(ln, 1, 3) == '{{{'
  end
  local pos = vim.api.nvim_win_get_cursor(0)
  local ri = pos[1]
  local rj = pos[1]
  while true do
    if check_stop_line(ri) then break end
    ri = ri - 1
  end
  while true do
    if check_stop_line(rj) then break end
    rj = rj + 1
  end
  local command=table.concat(vim.api.nvim_buf_get_lines(0, ri, rj-1, false), '\n')

  local res = vim.fn.system(command)

  local buf = vim.api.nvim_create_buf(false, true)
  vim.api.nvim_buf_set_option(buf, 'bufhidden', 'wipe')

  vim.fn.setreg('e', res)
  vim.api.nvim_buf_set_text(buf, 0, 0, -1, 0, vim.split(command .. '\n\n' .. res, '\n'))
  vim.api.nvim_buf_set_option(buf, 'modifiable', false) -- race here?

  local win = vim.api.nvim_open_win(buf, true, {
    relative='editor',
    width=vim.api.nvim_get_option('columns')-2,
    height=vim.api.nvim_get_option('lines')-3,
    col=1,
    row=1,
    style='minimal',
    border='rounded',
    noautocmd=1,
  })

  local opts = {buffer=buf}

  vim.keymap.set('n', '<space>www', function() vim.api.nvim_win_hide(win) end, opts) -- reset www
  vim.keymap.set('n', 'q',          function() vim.api.nvim_win_hide(win) end, opts)
  vim.keymap.set('n', '<esc><esc>', function() vim.api.nvim_win_hide(win) end, opts)

  vim.keymap.set('n', 'w', function() vim.api.nvim_win_set_option(0, "wrap", not vim.api.nvim_win_get_option(0, 'wrap')) end, opts)

end)
