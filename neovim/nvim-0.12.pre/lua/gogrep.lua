-- go install github.com/michurin/warehouse/go/gogrep@latest

local M = {}

M.gogrep_methods = {
  opts = { nargs = '*' },
  act = function()
    local directory = vim.fn.expand('%:~:.:h')
    local type_name = vim.fn.expand('<cword>')
    vim.system({ 'gogrep', 'methods', directory, type_name }, { text = true }, function(obj)
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
}

M.gogrep_strings = {
  opts = { nargs = 1 }, -- all args as one including spaces
  act = function(opts)
    local directory, pattern
    if #opts.fargs == 0 then
      print 'no args: nothing to do' -- consider vim.fn.expand('<cword>')?
      return
    end
    if #opts.fargs == 1 then
      directory = '.'
      pattern = opts.fargs[1]
    else
      directory, pattern = unpack(opts.fargs)
    end
    -- TODO if directory=@ -> vim.fn.expand('%:~:.:h')
    vim.system({ 'gogrep', 'strings', directory }, { text = true }, function(obj)
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
}

return M
