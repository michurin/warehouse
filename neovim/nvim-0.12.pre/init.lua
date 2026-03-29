require('configs')
local F = require('functions')
local S = require('sync_imports')
local RG = require('rg')
local CustomSelect = require('select')
local CustomSpelling = require('spell')

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
vim.keymap.set('n', '<space>ft', F.grep_in_files(F.files_from_cmd('git ls-files "*_test.go"'), F.cword),
  { noremap = true })
vim.keymap.set('v', '<space>ft',
  F.grep_in_files(F.files_from_cmd('git ls-files "*_test.go"'), F.visual_scalar),
  { noremap = true })
vim.keymap.set('n', '<space>fb', F.grep_in_files(F.files_buffers, F.cword), { noremap = true })
vim.keymap.set('v', '<space>fb', F.grep_in_files(F.files_buffers, F.visual_scalar), { noremap = true })
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
vim.api.nvim_create_user_command('BUF', F.qf_buffers.act, F.qf_buffers.opts)
vim.api.nvim_create_user_command('E', F.smart_open.act('edit'), F.smart_open.opts)
vim.api.nvim_create_user_command('ET', F.smart_open.act('tabnew'), F.smart_open.opts)
vim.api.nvim_create_user_command('L', F.exec_lua_command.act, F.exec_lua_command.opts)
vim.api.nvim_create_user_command('SH', F.exec_shell_command.act, F.exec_shell_command.opts)
vim.api.nvim_create_user_command('CC', function() vim.opt.colorcolumn = { 120 } end, {})

-- git
vim.api.nvim_create_user_command('GD', F.exec_git_diff_all.act, F.exec_git_diff_all.opts)
vim.api.nvim_create_user_command('GBL', F.exec_git_blame.act, F.exec_git_blame.opts)
vim.api.nvim_create_user_command('GLO', F.exec_git_log.act, F.exec_git_log.opts)

-- find by file name
vim.api.nvim_create_user_command('FF', F.file_search_command.act, F.file_search_command.opts)
vim.api.nvim_create_user_command('FG', F.gitls_search_command.act, F.gitls_search_command.opts)
vim.api.nvim_create_user_command('FB', function() print("TODO FIND BUFFER") end, {})

-- grep: find by content
vim.api.nvim_create_user_command('G', RG.search(), RG.opts)
vim.api.nvim_create_user_command('GL', RG.search('--glob', '*.go', '--glob', '!*_test.go'), RG.opts) -- TODO dedup
vim.api.nvim_create_user_command('GT', RG.search('--glob', '*_test.go'), RG.opts)
vim.api.nvim_create_user_command('GB', function(opts)
  local pattern = opts.args
  local bufs = vim.api.nvim_list_bufs()
  local items = {}
  for _, buf in ipairs(bufs) do
    if vim.api.nvim_buf_is_loaded(buf) then
      local name = vim.api.nvim_buf_get_name(buf)
      if name ~= '' and vim.fn.buflisted(buf) == 1 then
        local lines = vim.api.nvim_buf_get_lines(buf, 0, -1, false)
        for j, l in ipairs(lines) do
          table.insert(items, { filename = name, lnum = j, col = 0, text = l })
        end
      end
    end
  end
  items = vim.fn.matchfuzzy(items, pattern, { key = 'text' }) -- filter and sort as well
  vim.fn.setqflist({}, ' ', { title = 'Files (' .. pattern .. ')', items = items })
  vim.cmd.copen()
end, { nargs = '+' })

-- misc
vim.api.nvim_create_user_command('U', function()
  vim.cmd.edit(vim.fn.fnamemodify(vim.api.nvim_buf_get_name(0), ':h'))
end, {})

vim.api.nvim_create_user_command('UT', function()
  vim.cmd.tabnew(vim.fn.fnamemodify(vim.api.nvim_buf_get_name(0), ':h'))
end, {})

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
  vim.fn.setqflist({}, ' ', {
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

-- custom select
vim.ui.select = CustomSelect.select

-- custom spelling suggestion
vim.keymap.set('n', 'z=', CustomSpelling.act, CustomSpelling.opts)

-- idea
vim.api.nvim_create_autocmd('BufWritePre', {
  callback = function(event)
    local file = event.match
    local dir = vim.fn.fnamemodify(file, ':p:h')

    if vim.fn.isdirectory(dir) == 0 then
      vim.fn.mkdir(dir, 'p')
    end
  end,
})

vim.keymap.set('n', 'gf', function()
  local file = vim.fn.expand('<cfile>')
  local found = vim.fn.findfile(file, vim.o.path)
  if found == '' then
    found = vim.fn.fnamemodify(vim.fn.expand('%:p:h') .. '/' .. file, ':p')
  end
  vim.cmd.edit(found)
end)

-- idea

-- bad thing, choice.cmd is not so good as :tag
vim.keymap.set('n', '<C-]>', function()
  local tags = vim.fn.taglist(vim.fn.expand('<cword>'))
  print(#tags)
  if #tags == 0 then
    print('no tags')
    return
  end
  if #tags == 1 then
    local choice = tags[1]
    vim.cmd.edit(choice.filename)
    vim.cmd(choice.cmd)
    return
  end
  vim.ui.select(tags, {
    prompt = 'Select tag',
    format_item = function(item)
      return item.name .. ' -> ' .. item.filename
    end,
  }, function(choice)
    if choice then
      vim.cmd.edit(choice.filename)
      vim.cmd(choice.cmd)
    end
  end)
end)

-- idea
-- do something like that
-- nvim -o $(git diff --name-only --diff-filter=U --relative)
--
-- idea for git log
-- git log --graph --name-status
