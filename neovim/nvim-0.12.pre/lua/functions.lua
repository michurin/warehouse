local M = {}

-- -------------------------------

function M.cword()
  return vim.fn.expand('<cword>')
end

function M.visual_text()
  local m = vim.fn.mode() -- visualmode() предыдущий
  local a = vim.fn.getpos('v')
  local b = vim.fn.getpos('.')
  if a[2] > b[2] or (a[2] == b[2] and a[3] > b[3]) then
    a, b = b, a
  end
  -- vim.fn.setreg('g', vim.inspect(m)..'|'..vim.inspect(a)..'|'..vim.inspect(b)..'\n') -- debug
  if m == 'v' then
    return vim.api.nvim_buf_get_text(0, a[2]-1, a[3]-1, b[2]-1, b[3], {})
  elseif m == 'V' then
    return vim.api.nvim_buf_get_lines(0, a[2]-1, b[2], false)
  end
  return {}
end

function M.visual_line()
  local ln = M.visual_text()
  if #ln == 0 then return '' end
  return ln[1]
end

function M.input(prompt)
  return function()
    local ok, res = pcall(vim.fn.input, prompt..' ') -- prevent ctrl-c error
    if ok then
      return res
    end
    return ''
  end
end

function M.files_from_cmd(cmd)
  return function ()
    local files = vim.fn.systemlist(cmd)
    table.sort(files)
    return files
  end
end

function M.files_from_find_cdir()
  local dir = vim.fs.dirname(vim.api.nvim_buf_get_name(0)) -- vim.fn.expand('%:p:h')
  return M.files_from_cmd('find ' .. vim.fn.shellescape(dir) .. ' -type f')() -- will sort
end

function M.files_buffers()
  local bufs = vim.api.nvim_list_bufs()
  local names = {}
  for _, buf in ipairs(bufs) do
    if vim.api.nvim_buf_is_loaded(buf) then
      local name = vim.api.nvim_buf_get_name(buf)
      if name ~= '' and vim.fn.buflisted(buf) == 1 then
        table.insert(names, name)
      end
    end
  end
  return names
end

-- -------------------------------

function M.grep_in_files(source, inp)
  return function ()
    local files = source()
    if #files == 0 then
      print('No files found')
      return
    end
    local word = inp()
    if word == '' then
      print('No pattern')
      return
    end
    word = vim.fn.escape(word, '\\/.*$^~[]')
    vim.cmd('vimgrep /\\<' .. word .. '\\>/ ' .. table.concat(files, ' '))
    vim.cmd('copen')
    vim.fn.setreg('/', '\\<' .. word .. '\\>')
    vim.opt.hlsearch = true -- kick highlighting
  end
end

-- -------------------------------

local buff_last_pos = {}

M.qf_buffers_events = {'BufLeave', 'BufEnter', 'WinLeave', 'WinEnter'}

function M.qf_buffers_handler()
  local bufnr = vim.api.nvim_get_current_buf()
  local win = vim.api.nvim_get_current_win()
  if vim.api.nvim_win_is_valid(win) then
    buff_last_pos[bufnr] = vim.api.nvim_win_get_cursor(win)
  end
end

function M.qf_buffers()
  local items = {}
  local wins = vim.api.nvim_list_wins()
  local bufs = vim.api.nvim_list_bufs()

  for _, buf in ipairs(bufs) do
    if vim.api.nvim_buf_is_loaded(buf) then
      local name = vim.api.nvim_buf_get_name(buf)
      if name ~= '' and vim.fn.buflisted(buf) == 1 then
        local pos = buff_last_pos[buf]
        local debug = buf -- debugging
        for _, win in ipairs(wins) do
          if vim.api.nvim_win_get_buf(win) == buf then
            pos = vim.api.nvim_win_get_cursor(win)
            debug = debug .. '/' .. win
          end
        end
        local row = pos[1]
        local col = vim.fn.max({pos[2], 1}) -- can be 0
        local line = vim.api.nvim_buf_get_lines(buf, row - 1, row, false)[1]
        table.insert(items, {
          filename = name,
          lnum = row,
          col = col,
          text = line .. ' [' .. debug .. ']',
        })
      end
    end
  end

  table.sort(items, function(a, b)
    return a.filename < b.filename
  end)

  vim.fn.setqflist({}, 'r', {
    title = 'Buffers',
    items = items,
  })

  vim.cmd.copen()
end

-- -------------------------------

local viewing_buffer = -1

local function show_viewing_buffer (content, line)
  if not vim.api.nvim_buf_is_valid(viewing_buffer) then
    viewing_buffer = vim.api.nvim_create_buf(false, true)

    vim.api.nvim_buf_set_option(viewing_buffer, 'buftype', 'nofile') -- nvim_buf_set_option is legacy! TODO
    vim.api.nvim_buf_set_option(viewing_buffer, 'bufhidden', 'hide') -- важно
    vim.api.nvim_buf_set_option(viewing_buffer, 'swapfile', false)
    vim.api.nvim_buf_set_option(viewing_buffer, 'modifiable', false)

    vim.api.nvim_buf_set_name(viewing_buffer, '[OUTPUT]') -- nvim_buf_get_name всегда добавляет путь
  end
  vim.api.nvim_buf_set_option(viewing_buffer, 'modifiable', true)
  vim.api.nvim_buf_set_lines(viewing_buffer, 0, -1, false, content)
  vim.api.nvim_buf_set_option(viewing_buffer, 'modifiable', false)
  vim.api.nvim_win_set_buf(0, viewing_buffer) -- TODO in other window?

  vim.api.nvim_win_set_cursor(0, {line, 0})
  vim.fn.setpos("'<", {0, line, 0, 0})
  vim.fn.setpos("'>", {0, #content, 0, 0})
end

function M.paragraph_block ()
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
  return vim.api.nvim_buf_get_lines(0, ri, rj-1, false)
end

local function exec_bash(cmd_fetcherer)
  local cmd_lines = cmd_fetcherer()
  local cmd_result = vim.fn.systemlist(table.concat(cmd_lines, '\n'))
  local content = cmd_lines
  table.insert(content, '')
  table.insert(content, '')
  local line = #content + 1
  for i = 1, #cmd_result do
    table.insert(content, cmd_result[i])
  end
  table.insert(content, '')
  show_viewing_buffer(content, line)
end

local function exec_sql(cmd_fetcherer)
  local cmd_lines = cmd_fetcherer()
  local cmd = {}
  if vim.fn.match(cmd_lines[1], '^--') >= 0 then
    local x = cmd_lines[1]:gsub('^..', '')
    cmd = vim.fn.split(x)
  end
  table.insert(cmd, table.concat(cmd_lines, '\n'))
  local cmd_result = vim.fn.systemlist(cmd)
  local content = {}
  for i = 1, #cmd_result do
    if not (i == 1 and cmd_result[i] == 'mysql: [Warning] Using a password on the command line interface can be insecure.') then
      table.insert(content, cmd_result[i])
    end
  end
  show_viewing_buffer(content, 1)
end

function M.exec(cmd_fetcherer)
  return function ()
    local ft = vim.o.filetype
    if ft == 'sql' then
      exec_sql(cmd_fetcherer)
      return
    end
    if ft == 'sh' or ft == 'bash' then
      exec_bash(cmd_fetcherer)
      return
    end
    show_viewing_buffer({'Unknown type: '..ft}, 1)
  end
end

function M.exec_git_diff()
  local result = vim.fn.systemlist('git diff --no-prefix ' .. vim.fn.shellescape(vim.api.nvim_buf_get_name(0)))
  show_viewing_buffer(result, 1)
  vim.opt_local.filetype="diff"
end

function M.exec_git_diff_all(opts)
  local command = opts.args -- local command = vim.fn.input('git diff> ')
  local result = vim.fn.systemlist('git diff --no-prefix ' .. command)
  show_viewing_buffer(result, 1)
  vim.opt_local.filetype="diff"
end

function M.exec_git_blame()
  local result = vim.fn.systemlist('git blame ' .. vim.fn.shellescape(vim.api.nvim_buf_get_name(0)))
  local line = vim.api.nvim_win_get_cursor(0)[1]
  show_viewing_buffer(result, line)
end

function M.exec_command()
  local command = vim.fn.input('SH> ')
  local result = vim.fn.systemlist(command)
  show_viewing_buffer(result, 1)
end

function M.show_keys()
  local text = {}
  local mappings = vim.api.nvim_get_keymap('n')
  for _, m in ipairs(mappings) do
    table.insert(text, '{{{ '..m.mode..' "'..m.lhs..'" '..(m.desc or "<nil>"))
    local x = vim.split(vim.inspect(m), '\n')
    for i = 1, #x do
      table.insert(text, x[i])
    end
    table.insert(text, '}}}')
  end
  show_viewing_buffer(text, 1)
  vim.api.nvim_buf_call(0, function()
    vim.opt_local.foldmethod='marker'
    vim.opt_local.foldlevel=0
  end)
end

-- -------------------------------

function M.show_files()
  local content = vim.fn.systemlist('git ls-files "*.go"')
  show_viewing_buffer(content, 1)
end

function M.show_files_by_pattern()
  local pattern = vim.fn.input('*> ')
  local content = vim.fn.systemlist('find . -type f -name ' .. vim.fn.shellescape(pattern))
  show_viewing_buffer(content, 1)
end

-- -------------------------------

local function smart_find_and_fill(file, line)
  local stat = vim.loop.fs_stat(file)
  if stat == nil then
    local files = vim.fn.systemlist('find . -type f | grep '..vim.fn.shellescape(file))
    if #files == 0 then
      print('nofiles fallback')
      local f = file:gsub('^[^/]+/', '')
      if f == file then
        print('totally no files')
        return
      end
      smart_find_and_fill(f, line)
      return
    end
    if #files == 1 then
      file = files[1]
    else
      table.sort(files)
      local items = {}
      for _, f in ipairs(files) do
        f = f:gsub('^%./', '')
        local bufnr = vim.fn.bufnr(vim.fn.fnamemodify(f, ":p"))
        local text = '-'
        if bufnr > -1 then
          local lines = vim.api.nvim_buf_get_lines(bufnr, line-1, line, false)
          if #lines > 0 then
            text = lines[1]..' (b='..tostring(bufnr)..')'
          end
        end
        table.insert(items, {filename = f, lnum = line, col = 0, text = text})
      end
      vim.fn.setqflist({}, 'r', {title = 'FN', items = items})
      vim.cmd.copen()
      return
    end
  end
  vim.api.nvim_cmd({cmd = "edit", args = {file}}, {})
  vim.api.nvim_win_set_cursor(0, {tonumber(line), 0})
end

function M.smart_file_locate(opts) -- USAGE: vim.api.nvim_create_user_command('E', F, { nargs = 1 })
  local file, line = opts.args:match('^(.*):([0-9]+)$')
  if file == nil then
    print('nofile')
    return
  end
  if line == nil then
    print('noline')
    return
  end
  smart_find_and_fill(file, tonumber(line))
end

-- -------------------------------

vim.cmd([[
noremap <A-C-S-Up>  :-tabmove<cr>
noremap <A-C-S-Down> :+tabmove<cr>
]])

return M