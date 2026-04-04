-- :source %
-- :lua =XX('example', 'E')
-- :lua =XX('example', 'EE')
-- :lua =XX('example', 'ee')
-- :lua =XX(vim.fn.expand('%:~:.:h'), vim.fn.expand('<cword>'))
function XX(directory, type_name)
  vim.system({ 'go', 'run', '.', directory, type_name }, { text = true }, function(obj)
    local items = vim.json.decode(obj.stdout)

    if #items == 0 then
      print('nothing found')
      vim.schedule(function()
        vim.cmd('cclose')
      end)
      return
    end

    local qf = {}
    for _, item in ipairs(items) do
      table.insert(qf, {
        filename = item.file,
        lnum = item.line,
        col = item.col,
        text = item.message,
      })
    end

    vim.schedule(function()
      vim.fn.setqflist(qf, 'r')
      vim.cmd('copen')
    end)
  end)
end
