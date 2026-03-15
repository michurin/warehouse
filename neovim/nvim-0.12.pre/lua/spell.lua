local M = {}

function M.act()
  local word = vim.fn.expand('<cword>')
  local suggestions = vim.fn.spellsuggest(word)

  if #suggestions == 0 then
    print('no suggestions')
    return
  end

  vim.ui.select(suggestions, {
    prompt = 'Spell suggest' .. word,
  }, function(choice)
    if choice then
      vim.cmd('normal! ciw' .. choice)
    end
  end)
end

M.opts = { desc = 'Spell suggestions with select UI' }

return M
