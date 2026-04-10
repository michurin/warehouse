local M = {}

-- -------------------------------

function M.cword()
  return vim.fn.expand('<cword>')
end

function M.visual_scalar()
  return table.concat(M.visual_text_block(), '\n')
end

function M.input(prompt)
  return function()
    local ok, res = pcall(vim.fn.input, prompt .. ' ') -- prevent ctrl-c error
    if ok then
      return res
    end
    return ''
  end
end

function M.files_from_cmd(cmd)
  return function()
    local files = vim.fn.systemlist(cmd)
    table.sort(files)
    return files
  end
end

function M.files_from_find_cdir()
  local dir = vim.fs.dirname(vim.api.nvim_buf_get_name(0))                    -- vim.fn.expand('%:p:h')
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
  return function()
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
    -- vim.api.nvim_cmd({cmd = 'vimgrep', args = {'cmd', '%'}}, {})
    local args = files
    table.insert(args, 1, '/\\<' .. word .. '\\>/')
    local ok, _ = pcall(vim.api.nvim_cmd, { cmd = 'vimgrep', args = args }, {})
    if not ok then
      vim.notify('No matche') -- print won't work
      return
    end
    vim.cmd.copen()
    vim.fn.setreg('/', '\\<' .. word .. '\\>')
    vim.opt.hlsearch = true -- kick highlighting
  end
end

-- -------------------------------

local buff_last_pos = {}

M.qf_buffers_events = { 'BufLeave', 'BufEnter', 'WinLeave', 'WinEnter' }

function M.qf_buffers_handler()
  local bufnr = vim.api.nvim_get_current_buf()
  local win = vim.api.nvim_get_current_win()
  if vim.api.nvim_win_is_valid(win) then
    buff_last_pos[bufnr] = vim.api.nvim_win_get_cursor(win)
  end
end

M.qf_buffers = {
  opts = {},
  act = function()
    local items = {}
    local wins = vim.api.nvim_list_wins()
    local bufs = vim.api.nvim_list_bufs()

    for _, buf in ipairs(bufs) do
      if vim.api.nvim_buf_is_loaded(buf) then
        local name = vim.api.nvim_buf_get_name(buf)
        if name ~= '' and vim.fn.buflisted(buf) == 1 then
          local pos = buff_last_pos[buf]
          local debug = tostring(buf) -- debugging
          for _, win in ipairs(wins) do
            if vim.api.nvim_win_get_buf(win) == buf then
              pos = vim.api.nvim_win_get_cursor(win)
              debug = debug .. '/' .. win
            end
          end
          local row = pos[1]
          local col = vim.fn.max({ pos[2], 1 }) -- can be 0
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

    vim.fn.setqflist({}, ' ', {
      title = 'Buffers',
      items = items,
    })

    vim.cmd.copen()
  end
}

-- -------------------------------

local viewing_buffer = -1

local function show_viewing_buffer(content, line)
  if not vim.api.nvim_buf_is_valid(viewing_buffer) then
    viewing_buffer = vim.api.nvim_create_buf(false, true)

    vim.api.nvim_set_option_value('buftype', 'nofile', { buf = viewing_buffer })
    vim.api.nvim_set_option_value('bufhidden', 'hide', { buf = viewing_buffer }) -- важно
    vim.api.nvim_set_option_value('swapfile', false, { buf = viewing_buffer })
    vim.api.nvim_set_option_value('modifiable', false, { buf = viewing_buffer })

    vim.api.nvim_buf_set_name(viewing_buffer, '[OUTPUT]') -- nvim_buf_get_name всегда добавляет путь
  end
  vim.api.nvim_set_option_value('modifiable', true, { buf = viewing_buffer })
  vim.api.nvim_buf_set_lines(viewing_buffer, 0, -1, false, content)
  vim.api.nvim_set_option_value('modifiable', false, { buf = viewing_buffer })
  vim.api.nvim_win_set_buf(0, viewing_buffer) -- TODO in other window?

  vim.api.nvim_win_set_cursor(0, { line, 0 })
  vim.fn.setpos("'<", { 0, line, 0, 0 })
  vim.fn.setpos("'>", { 0, #content, 0, 0 })
end

function M.paragraph_text_block()
  local check_stop_line = function(r)
    local ln = vim.api.nvim_buf_get_lines(0, r - 1, r, false)[1]
    return not ln or ln == ''
    -- or string.sub(ln, 1, 1) == '#' or string.sub(ln, 1, 2) == '--' or string.sub(ln, 1, 3) == '{{{'
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
  return vim.api.nvim_buf_get_lines(0, ri, rj - 1, false)
end

function M.visual_text_block()
  local m = vim.fn.mode() -- visualmode() предыдущий
  local a = vim.fn.getpos('v')
  local b = vim.fn.getpos('.')
  if a[2] > b[2] or (a[2] == b[2] and a[3] > b[3]) then
    a, b = b, a
  end
  if m == 'v' then
    local ll = vim.api.nvim_buf_get_lines(0, a[2] - 1, b[2], false)
    local s = vim.str_byteindex(ll[1], 'utf-8', a[3] - 1)
    local e = vim.str_byteindex(ll[#ll], 'utf-8', b[3])
    ll[#ll] = string.sub(ll[#ll], 0, e + 1)
    ll[1] = string.sub(ll[1], s + 1)
    -- vim.fn.setreg('g', vim.inspect(ll) .. '|' .. vim.inspect(m) .. '|' .. vim.inspect(a) .. '|' .. vim.inspect(b) .. '\n') -- debug
    return ll
  elseif m == 'V' then
    return vim.api.nvim_buf_get_lines(0, a[2] - 1, b[2], false)
  end
  return {}
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
  -- assume comment like:
  -- mysql -h 127.0.0.1 --port 10150 -u user -ppass databese -t -e
  -- do not forget -e
  -- -E for vertical layout
  local cmd_lines = cmd_fetcherer()
  local cmd = {}
  for i = 1, #cmd_lines do
    if vim.fn.match(cmd_lines[i], '^--') >= 0 then
      local s = cmd_lines[i]:gsub('^..', '')
      local x = vim.fn.split(s)
      if #x > 0 and x[1]:sub(1, 3) == '{{{' then
        table.remove(x, 1)
      end
      if #x > 0 and x[1] == 'mysql' then
        cmd = x
        break
      end
    end
  end
  if #cmd == 0 then
    print('no command')
    return
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
  return function()
    local ft = vim.o.filetype
    if ft == 'sql' then
      exec_sql(cmd_fetcherer)
      return
    end
    if ft == 'sh' or ft == 'bash' then
      exec_bash(cmd_fetcherer)
      return
    end
    show_viewing_buffer({ 'Unknown type: ' .. ft }, 1)
  end
end

M.exec_shell_command = {
  opts = { nargs = '+', complete = 'shellcmdline' },
  act = function(opts) show_viewing_buffer(vim.fn.systemlist(opts.args), 1) end,
}

M.exec_git_diff_all = {
  opts = {
    nargs = '*',
    complete = function()
      return vim.fn.systemlist({
        'git', 'branch', '-q', '-a', '--sort=-committerdate', '--format=%(refname:short)', '--no-color' })
    end
  },
  act = function(opts)
    local result = vim.fn.systemlist({ 'git', 'diff', '--no-prefix', '--no-color', unpack(opts.fargs) })
    show_viewing_buffer(result, 1)
    vim.opt_local.filetype = 'diff'
  end
}

M.exec_git_log = {
  opts = {},
  act = function()
    local result = vim.fn.systemlist({
      'git', 'log', '--graph', '--reflog', '--branches', '--remotes', '--tags', '--decorate', '-n', '120',
      '--date=format:%y-%m-%d %H:%M',
      '--format=format:%C(auto)%h%C(reset)%C(auto)%d%C(reset) %C(green)%ad%C(reset) %C(cyan)%an%C(reset)%n%C(auto)%s%C(reset)%n',
      '--name-status',
    })
    show_viewing_buffer(result, 1)
    vim.opt_local.filetype = 'gitlog' -- TODO
    vim.opt_local.listchars = 'trail: ,tab:  ,nbsp: ,extends:▶,precedes:◀'
  end
}

M.exec_git_blame = {
  opts = {},
  act = function()
    local result = vim.fn.systemlist({
      'git', 'blame', '--date=format:%y-%m-%d %H:%M', '--',
      vim.api.nvim_buf_get_name(0) })
    local line = vim.api.nvim_win_get_cursor(0)[1]
    show_viewing_buffer(result, line)
  end
}

M.exec_lua_command = {
  opts = { nargs = '+', complete = 'lua' },
  act = function(opts)
    local command = opts.args
    local chunk = load('return ' .. command)
    if not chunk then
      print('no chunk')
      return
    end
    local ok, result = pcall(chunk)
    if not ok then
      print('not ok: ' .. result)
      return
    end
    local lines = vim.split(vim.inspect(result), '\n')
    if #lines == 1 and #lines[1] < 40 then -- the result is one short line
      print(lines[1])
      return
    end
    show_viewing_buffer(lines, 1)
    vim.opt_local.filetype = 'lua'
  end,
}

function M.show_keys()
  local text = {}
  local mappings = vim.api.nvim_get_keymap('n')
  for _, m in ipairs(mappings) do
    table.insert(text, '{{{ ' .. m.mode .. ' "' .. m.lhs .. '" ' .. (m.desc or '<nil>'))
    local x = vim.split(vim.inspect(m), '\n')
    for i = 1, #x do
      table.insert(text, x[i])
    end
    table.insert(text, '}}}')
  end
  show_viewing_buffer(text, 1)
  vim.api.nvim_buf_call(0, function()
    vim.opt_local.foldmethod = 'marker'
    vim.opt_local.foldlevel = 0
  end)
end

-- -------------------------------

function M.show_files()
  local content = vim.fn.systemlist('git ls-files "*.go"') -- TODO split
  show_viewing_buffer(content, 1)
end

function M.show_files_by_pattern()
  local pattern = vim.fn.input('*> ')
  local content = vim.fn.systemlist('find . -type f -name ' .. vim.fn.shellescape(pattern)) -- TODO split
  show_viewing_buffer(content, 1)
end

-- -------------------------------

function M.fuzzy_search() -- TODO just idea; lua require('functions').fuzzy_search()
  local s = vim.fn.input('>')
  -- print(s)
  local bufs = vim.api.nvim_list_bufs()
  local items = {}
  for _, buf in ipairs(bufs) do
    if vim.api.nvim_buf_is_loaded(buf) then
      local name = vim.api.nvim_buf_get_name(buf)
      if name ~= '' and vim.fn.buflisted(buf) == 1 then
        local l = vim.api.nvim_buf_get_lines(buf, 0, -1, false)
        local rx = vim.fn.matchfuzzypos(l, s, { matchseq = 1 }) -- matchfuzzypos? sort by score?
        local r = rx[1]
        local scores = rx[3]
        if #r > 0 then
          local filename = vim.fn.fnamemodify(name, ':~:.')
          for i = 1, #l do
            for j = 1, #r do -- TODO nested loop
              if l[i] == r[j] then
                table.insert(items, { filename = filename, lnum = i, col = 0, text = tostring(scores[j]) .. '|' .. l[i] })
                break
              end
            end
          end
        end
      end
    end
  end
  vim.fn.setqflist({}, ' ', { title = 'FN', items = items }) -- TODO items? or list argument?
  vim.cmd.copen()
end

-- -------------------------------

function M.qf_do(command)
  return function()
    local ok, _ = pcall(function() vim.cmd(command) end)
    if not ok then
      print('no list')
    end
  end
end

-- -------------------------------

function M.copy_bookmark_to_f()
  if vim.bo.filetype == 'qf' then
    local it = vim.fn.getqflist()
    local s = {}
    for _, i in ipairs(it) do
      table.insert(s,
        vim.fn.fnamemodify(vim.api.nvim_buf_get_name(i.bufnr), ':~:.') ..
        ':' .. tostring(i.lnum) .. ' ' .. i.text:gsub('^%s+', ''):gsub('%s+$', ''))
    end
    local m = table.concat(s, '\n')
    vim.fn.setreg('f', m)
    print('f:' .. tostring(#s) .. ' lines (qf)')
    return
  end
  local l = vim.api.nvim_win_get_cursor(0)[1]
  local t = vim.api.nvim_buf_get_lines(0, l - 1, l, false)[1]
  local p = vim.api.nvim_buf_get_name(0)
  local m = vim.fn.fnamemodify(p, ':~:.') .. ':' .. tostring(l) .. ' ' .. t:gsub('^%s+', ''):gsub('%s+$', '')
  vim.fn.setreg('f', m)
  print('f:' .. m)
end

-- -------------------------------

local function list_of_files_to_qf_items(ff)
  local items = {}
  for _, f in ipairs(ff) do
    local lnum = 1
    local text = '-'
    local bufnr = vim.fn.bufnr(vim.fn.fnamemodify(f, ':p'))
    if bufnr >= 0 then
      local pos = vim.api.nvim_buf_get_mark(bufnr, '"')
      lnum = pos[1]
      text = vim.api.nvim_buf_get_lines(bufnr, lnum - 1, lnum, false)[1]
    end
    table.insert(items, { filename = f, lnum = lnum, col = 0, text = text })
  end
  return items
end

M.gitls_search_command = {
  opts = { nargs = '+' },
  act = function(opts)
    local pattern = opts.args
    local ff = vim.fn.systemlist({ 'git', 'ls-tree', '--name-only', '-r', 'HEAD' })
    ff = vim.fn.matchfuzzy(ff, pattern) -- filter and sort as well
    local items = list_of_files_to_qf_items(ff)
    vim.fn.setqflist({}, ' ', { title = 'Files (' .. pattern .. ')', items = items })
    vim.cmd.copen()
  end,
}

M.file_search_command = {
  opts = { nargs = '+' },
  act = function(opts)
    local pattern = opts.args
    local ff = vim.fn.systemlist({ 'find', '.', '-type', 'f' })
    ff = vim.fn.matchfuzzy(ff, pattern) -- filter and sort as well
    local items = list_of_files_to_qf_items(ff)
    vim.fn.setqflist({}, ' ', { title = 'Files (' .. pattern .. ')', items = items })
    vim.cmd.copen()
  end,
}

-- -------------------------------

M.smart_open = {
  opts = { nargs = '+' },
  act = function(cmd)
    return function(opts)
      -- :lua x('edit', 'a.go')
      -- :lua x('edit', 'a.go:22')
      -- :lua x('edit', 'a.go#L22')
      -- :lua x('newtab', 'a.go#L22')
      -- :lua x('vsplit', 'a.go#L22')
      local pat = opts.args
      local file = pat
      local line = 1
      local p, q = pat:match("(.+)[:#]L?([0-9]+)")
      if p ~= nil then
        local x = tonumber(q)
        if x ~= nil then
          file = p
          line = x
        end
      end
      local stat = vim.loop.fs_stat(file)
      if stat ~= nil then
        vim.api.nvim_cmd({ cmd = cmd, args = { file } }, {})
        return
      end
      local files = vim.fn.systemlist('find . -type f -print0 | grep -z --color=never ' ..
        vim.fn.shellescape(file) .. ' | tr \'\\0\' \'\\n\'')
      if #files == 0 then
        print('no files')
        return
      end
      if #files == 1 then
        vim.api.nvim_cmd({ cmd = cmd, args = { files[1] } }, {})
        return
      end
      table.sort(files)
      local items = {}
      for _, f in ipairs(files) do
        f = f:gsub('^%./', '')
        local bufnr = vim.fn.bufnr(vim.fn.fnamemodify(f, ':p'))
        local text = '-'
        if bufnr > -1 then
          local lines = vim.api.nvim_buf_get_lines(bufnr, line - 1, line, false)
          if #lines > 0 then
            text = lines[1] .. ' (b=' .. tostring(bufnr) .. ')'
          end
        end
        table.insert(items, { filename = vim.fn.fnamemodify(f, ':~:.'), lnum = line, col = 0, text = text })
      end
      vim.fn.setqflist({}, ' ', { title = 'F: ' .. pat, items = items })
      vim.cmd.copen()
    end
  end
}

M.go_alt = {
  opts = {},
  act = function(cmd)
    return function()
      local file = vim.api.nvim_buf_get_name(0)
      if file == "" then
        print("no file")
        return
      end
      local new_file
      if file:match("_test%.go$") then
        new_file = file:gsub("_test%.go$", ".go")
      elseif file:match("%.go$") then
        new_file = file:gsub("%.go$", "_test.go")
      else
        print("no go file")
        return
      end
      vim.cmd({ cmd = cmd, args = { new_file } }) -- local escaped = vim.fn.fnameescape(new_file) -- TODO?
    end
  end,
}

return M
