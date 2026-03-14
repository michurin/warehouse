require('configs')
local F = require('functions')
local S = require('sync_imports')

--

vim.lsp.enable({
  'gopls',
  'lua_ls',  -- sudo pacman -Suy lua-language-server
  'bashls',  -- sudo pacman -Suy bash-language-server
  'pyright', -- sudo pacman -Suy pyright
  'ts_ls',   -- sudp pacman -Suy typescript-language-server
  'protols', -- go install github.com/lasorda/protobuf-language-server@master
})

vim.diagnostic.config({ virtual_text = true })

vim.keymap.set('n', 'grs', function() -- like gra but filling structures only
  vim.lsp.buf.code_action({ apply = true, filter = function(action) return action.kind == "refactor.rewrite.fillStruct" end })
end, { noremap = true })

vim.keymap.set('n', 'gd', function()
  vim.lsp.buf.definition()
end, { noremap = true })
vim.keymap.set('n', '<C-w>gd', function()
  vim.cmd('tab split')
  vim.lsp.buf.definition()
end, { noremap = true })

vim.api.nvim_create_autocmd('BufWritePre', { -- TODO move to ft
  callback = function()
    if vim.bo.filetype == 'go' then
      S.organize_imports_sync(1000)
    end
    vim.lsp.buf.format()
  end
})

--

vim.keymap.set('n', ']j', F.qf_do('cnewer'), { noremap = true })
vim.keymap.set('n', '[j', F.qf_do('colder'), { noremap = true })

--

vim.keymap.set('n', '<space>fl', F.grep_in_files(F.files_from_cmd('git ls-files "*.go" ":!*_test.go"'), F.cword),
  { noremap = true })
vim.keymap.set('v', '<space>fl',
  F.grep_in_files(F.files_from_cmd('git ls-files "*.go" ":!*_test.go"'), F.visual_scalar),
  { noremap = true })
vim.keymap.set('n', '<space>fL',
  F.grep_in_files(F.files_from_cmd('git ls-files "*.go" ":!*_test.go"'), F.input('PATTERN>')), { noremap = true })
vim.keymap.set('n', '<space>ft', F.grep_in_files(F.files_from_cmd('git ls-files "*_test.go"'), F.cword),
  { noremap = true })
vim.keymap.set('v', '<space>ft',
  F.grep_in_files(F.files_from_cmd('git ls-files "*_test.go"'), F.visual_scalar),
  { noremap = true })
vim.keymap.set('n', '<space>fT', F.grep_in_files(F.files_from_cmd('git ls-files "*_test.go"'), F.input('PATTERN>')),
  { noremap = true })
vim.keymap.set('n', '<space>fb', F.grep_in_files(F.files_buffers, F.cword), { noremap = true })
vim.keymap.set('v', '<space>fb', F.grep_in_files(F.files_buffers, F.visual_scalar), { noremap = true })
vim.keymap.set('n', '<space>fB', F.grep_in_files(F.files_buffers, F.input('PATTERN>')), { noremap = true })
vim.keymap.set('n', '<space>fa', F.grep_in_files(F.files_from_cmd('find . -type f'), F.cword), { noremap = true })
vim.keymap.set('v', '<space>fa', F.grep_in_files(F.files_from_cmd('find . -type f'), F.visual_scalar),
  { noremap = true })
vim.keymap.set('n', '<space>fA', F.grep_in_files(F.files_from_cmd('find . -type f'), F.input('PATTERN>')),
  { noremap = true })
vim.keymap.set('n', '<space>fd', F.grep_in_files(F.files_from_find_cdir, F.cword), { noremap = true })
vim.keymap.set('v', '<space>fd', F.grep_in_files(F.files_from_find_cdir, F.visual_scalar),
  { noremap = true })
vim.keymap.set('n', '<space>fD', F.grep_in_files(F.files_from_find_cdir, F.input('PATTERN>')), { noremap = true })

--

vim.keymap.set('n', '<space>ff', F.show_files, { noremap = true })
vim.keymap.set('n', '<space>fF', F.show_files_by_pattern, { noremap = true })

-- TODO move to go.lua? or to lsp settings?
-- outgoing_calls looks buggy
vim.keymap.set('n', '<space>fi', function()
  local view = vim.fn.winsaveview()
  vim.cmd('normal [[0lllll')
  vim.lsp.buf.incoming_calls()
  vim.fn.winrestview(view)
end, { noremap = true })

--

vim.api.nvim_create_autocmd(F.qf_buffers_events, { callback = F.qf_buffers_handler })
vim.api.nvim_create_user_command('BUF', F.qf_buffers, {})
vim.api.nvim_create_user_command('E', F.smart_file_locate, { nargs = 1 })
vim.api.nvim_create_user_command('D', F.exec_git_diff_all.act, F.exec_git_diff_all.opts)
vim.api.nvim_create_user_command('L', F.exec_lua_command.act, F.exec_lua_command.opts)
vim.api.nvim_create_user_command('SH', F.exec_shell_command.act, F.exec_shell_command.opts)
vim.api.nvim_create_user_command('CC', function() vim.opt.colorcolumn = { 120 } end, {})
vim.api.nvim_create_user_command('J', function() -- TODO idea: put jumplist to QF
  local jumplist, idx = unpack(vim.fn.getjumplist())
  local items = {}
  for i = #jumplist, 1, -1 do
    local jmp = jumplist[i]
    local buf = jmp.bufnr
    if not (vim.api.nvim_buf_is_loaded(buf) and vim.api.nvim_buf_is_valid(buf) and vim.bo[buf].buflisted) then
      break
    end
    table.insert(items, {
      filename = vim.api.nvim_buf_get_name(buf),
      lnum = jmp.lnum,
      col = jmp.col, -- TODO +jmp.addcol?
      text = vim.api.nvim_buf_get_lines(buf, jmp.lnum - 1, jmp.lnum, false)[1],
    })
  end
  vim.fn.setqflist({}, 'r', {
    title = 'Jumps',
    items = items,
  })
  -- TODO vim.cmd('cc ' .. tostring(#jumplist - idx + 1))
  print(#jumplist, idx)
  vim.cmd.copen()
end, {})

--

vim.keymap.set('n', '<space>gf', F.copy_bookmark_to_f, { noremap = true })

--

vim.keymap.set('n', '<space>www', F.exec(F.paragraph_text_block), { noremap = true })
vim.keymap.set('v', '<space>www', F.exec(F.visual_text_block), { noremap = true })

vim.keymap.set('n', '<space>wwd', F.exec_git_diff, { noremap = true })
vim.keymap.set('n', '<space>wwg', F.exec_git_blame, { noremap = true })

vim.keymap.set('n', '<space>wee', function()
  vim.system(
    {
      vim.fn.stdpath('config') .. '/bin/vim-helper-open-git',
      vim.api.nvim_buf_get_name(0),
      tostring(vim.api.nvim_win_get_cursor(0)[1]),
    },
    {
      stdout = false,
      stderr = false,
      tmeout = 5000,
    }
  )
end, { noremap = true })

vim.keymap.set('n', '<space>bb', F.show_keys, { noremap = true })

vim.keymap.set("n", "<space>hi", function() vim.cmd("Inspect") end, { noremap = true })

-- experimental

-- vim.api.nvim_create_user_command('QFSave', function()
--   local result = {}
--   local info = vim.fn.getqflist({ nr = '$' })
--   for i = 1, info.nr do
--     local items = vim.fn.getqflist({ id = i, all = 1 })
--     for _, item in ipairs(items) do
--       if item.bufnr and item.bufnr ~= 0 then
--         item.filename = vim.api.nvim_buf_get_name(item.bufnr)
--         item.bufnr = nil
--       end
--     end
--     table.insert(result, items)
--   end
--   vim.fn.writefile({ vim.fn.json_encode(result) }, 'qf.json')
-- end, {})
--
-- vim.api.nvim_create_user_command("QFLoad", function()
--   local data = vim.fn.json_decode(table.concat(vim.fn.readfile("qf.json"), "\n"))
--   for _, qf in ipairs(data) do
--     vim.fn.setqflist({}, 'a', qf)
--   end
-- end, {})

--

vim.cmd([[
noremap <A-C-S-Up>   :-tabmove<cr>
noremap <A-C-S-Down> :+tabmove<cr>
set spell spelllang=en_us,ru_yo,el spelloptions=camel
]])

-- ideas:

-- TODO MOVE TO SOMEWHERE
function _G.x(items, opts, callback)
  -- print(vim.inspect(opts))
  -- if true then return end
  -- local items = {'one', 'two', 'three', 'x', 'xx', 'xxx'}
  -- local opts = {prompt = 'PROMPT>'}
  -- local callback = function (c, i)
  --    print(c, i)
  -- end
  -- vim.ui.select(items, opts, callback)
  -- print('ok')
  --
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

-- TODO MOVE TO CONFIG
vim.api.nvim_set_hl(0, 'FloatTitle', { fg = '#00cccc' })
vim.api.nvim_set_hl(0, 'FloatFooter', { fg = '#00cccc' })
vim.api.nvim_set_hl(0, 'FloatBorder', { fg = '#00cccc' })
vim.api.nvim_set_hl(0, 'NormalFloat', { fg = '#cccccc' })

-- TODO SETUP IT HERE
vim.ui.select = _G.x

-- ONE MORE IDEA
vim.keymap.set('n', 'z=', function()
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
end, { desc = 'Spell suggestions with select UI' })
