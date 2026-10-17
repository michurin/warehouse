-- Highlighting of search pattern in quickfix window

local M = {}

-- Replaces previous highlighting (if any) in the current window with the new one.
-- Call it right after copen: copen makes quickfix window current.
function M.set(regex)
  if vim.w.qfhl_match_id then
    pcall(vim.fn.matchdelete, vim.w.qfhl_match_id)
    vim.w.qfhl_match_id = nil
  end

  if regex then
    vim.w.qfhl_match_id = vim.fn.matchadd('Title', regex)
  end
end

-- Regex: string literally, case-sensitive
function M.literal(s)
  if s == '' then
    return nil
  end
  return [[\V]] .. vim.fn.escape(s, [[\]])
end

-- Regex: string literally, case-sensitive only if it has uppercase (like rg --smart-case)
function M.smartcase(s)
  if s == '' then
    return nil
  end
  local case = s:match('%u') and [[\C]] or [[\c]]
  return case .. [[\V]] .. vim.fn.escape(s, [[\]])
end

-- Regex: any of words literally, case-insensitive (like matchfuzzy filtering)
function M.words(s)
  local words = vim.tbl_map(function(w) return vim.fn.escape(w, [[\]]) end, vim.split(s, '%s+', { trimempty = true }))
  if #words == 0 then
    return nil
  end
  return [[\c\V]] .. table.concat(words, [[\|]])
end

return M
