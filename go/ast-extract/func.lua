function M(directory, type_name)
  vim.system({ 'go', 'run', '.', 'methods', directory, type_name }, { text = true }, function(obj)
    local items = vim.json.decode(obj.stdout)

    if #items == 0 then
      print('nothing found')
      vim.schedule(function()
        vim.cmd('cclose')
      end)
      return
    end

    vim.schedule(function()
      vim.fn.setqflist(items, 'r')
      vim.cmd('copen')
    end)
  end)
end

function S(directory, pattern)
  vim.system({ 'go', 'run', '.', 'strings', directory }, { text = true }, function(obj)
    local itemsX = vim.json.decode(obj.stdout)

    vim.schedule(function()
      local items = vim.fn.matchfuzzy(itemsX, pattern, { key = 'text' })

      if #items == 0 then
        print('nothing found')
        vim.schedule(function()
          vim.cmd('cclose')
        end)
        return
      end

      vim.fn.setqflist(items, 'r')
      vim.cmd('copen')
    end)
  end)
end
