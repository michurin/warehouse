local M = {}

function M.select(items, opts, callback)
  opts = opts or {}
  local prompt = opts.prompt or "Select"
  local format = opts.format_item or tostring
  local max_height = opts.max_height or 12 -- TODO hardcoded hard limit
  local title = ' ' .. prompt:gsub('^%s+', ''):gsub('[%s:]+$', '') .. ' '

  -- print(vim.o.filetype) -- TODO: base highlighting and auto focusing on ft=go, title='Code actions'

  local labels = {}
  local jump_to = 1
  for j, item in ipairs(items) do
    local text = format(item)
    if text:lower():match('fill') then -- TODO make this logic based on file type
      jump_to = j
    end
    table.insert(labels, text)
  end

  local buf = vim.api.nvim_create_buf(false, true) -- listed=false, scratch=true (remove on win close)
  local width = #title
  for _, item in ipairs(labels) do
    width = math.max(width, vim.fn.strdisplaywidth(item))
  end
  width = width + 0 -- TODO

  local height = math.min(#items, max_height)
  local win = vim.api.nvim_open_win(buf, true, {
    relative = 'editor',
    width = width,
    height = height,
    row = 0,
    col = 0,
    style = 'minimal',
    border = 'rounded',
    title = title,
    footer = ' ' .. tostring(#items) .. ' ',
  })

  vim.wo[win].cursorline = true
  vim.wo[win].number = false
  vim.wo[win].relativenumber = false
  vim.wo[win].signcolumn = 'no'
  vim.bo[buf].bufhidden = 'wipe' -- important for auto removing
  vim.bo[buf].filetype = 'custom_select'

  vim.fn.matchadd('Comment', [[^Fill]])              -- win level. Use nvim_buf_add_highlight for buffer level
  vim.fn.matchadd('Statement', [[^Fill \zs\S\+\ze]]) -- TODO make highlighting based on file type
  vim.fn.matchadd('Type', [[\zs"[^"]\+"\ze]])
  vim.fn.matchadd('Constant', [[^Browse documentation]])

  vim.bo[buf].modifiable = true
  vim.api.nvim_buf_set_lines(buf, 0, -1, false, labels)
  vim.bo[buf].modifiable = false

  vim.api.nvim_win_set_cursor(win, { jump_to, 0 })

  local function close()
    if vim.api.nvim_win_is_valid(win) then
      vim.api.nvim_win_close(win, true)
    end
  end

  local function confirm()
    local row = vim.api.nvim_win_get_cursor(win)[1]
    close()
    if callback then
      callback(items[row], row)
    end
  end

  local function map(lhs, rhs)
    vim.keymap.set('n', lhs, rhs, { buffer = buf, nowait = true })
  end

  map('<CR>', confirm)
  map('<Esc>', close)
  map('q', close)
end

return M
