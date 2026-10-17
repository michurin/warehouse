-- Sorting quickfix items: mocks/tests to the bottom

local M = {}

-- Items can have filename (setqflist input) or bufnr only (parsed by getqflist)
local function filename(item)
  if item.filename then
    return item.filename
  end
  if item.bufnr and item.bufnr > 0 then
    return vim.api.nvim_buf_get_name(item.bufnr)
  end
  return ""
end

function M.is_mock(item)
  local path = filename(item):lower()

  return path:match("/mocks?/")
      or path:match("/tests?/")
      or path:match("/integration/")
      or path:match("/mock_")
      or path:match("_mock%.go$")
      or path:match("_test%.go$")
end

function M.is_func(item)
  return (item.text or ""):match("^func") ~= nil
end

-- Sorts items in place. Mocks/tests always go last. Stable: ties keep the original order.
-- opts.by_func: functions first
-- opts.by_location: then by filename and lnum
function M.sort(items, opts)
  opts = opts or {}

  local index = {}
  for i, item in ipairs(items) do
    index[item] = i
  end

  table.sort(items, function(a, b)
    local am = M.is_mock(a)
    local bm = M.is_mock(b)

    if am ~= bm then
      return not am
    end

    if opts.by_func then
      local af = M.is_func(a)
      local bf = M.is_func(b)

      if af ~= bf then
        return af
      end
    end

    if opts.by_location then
      local afn = filename(a)
      local bfn = filename(b)

      if afn ~= bfn then
        return afn < bfn
      end

      if a.lnum ~= b.lnum then
        return (a.lnum or 0) < (b.lnum or 0)
      end
    end

    return index[a] < index[b]
  end)

  return items
end

return M
