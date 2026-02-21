local opt = vim.opt
opt.guicursor = 'n-c-sm:block,i-ci-ve:ver25,r-cr-o-v:hor20'
-- opt.colorcolumn = '80'
opt.signcolumn = 'number' -- Always show sign column
opt.termguicolors = true -- Enable true colors
opt.ignorecase = true -- Ignore case in search
opt.swapfile = false -- Disable swap files
opt.autoindent = true -- Enable auto indentation
opt.expandtab = false
opt.tabstop = 4 -- Number of spaces for a tab
opt.softtabstop = 4 -- Number of spaces for a tab when editing
opt.shiftwidth = 4 -- Number of spaces for autoindent
opt.shiftround = true -- Round indent to multiple of shiftwidth
opt.listchars = 'trail:+,tab:▹·,nbsp:␣,extends:▶,precedes:◀' -- Characters to show for tabs, spaces, and end of line
opt.list = true -- Show whitespace characters
opt.number = true -- Show line numbers
opt.relativenumber = false
opt.numberwidth = 2 -- Width of the line number column
opt.wrap = false -- Disable line wrapping
-- opt.cursorline = true
opt.scrolloff = 8 -- Keep 8 lines above and below the cursor
opt.inccommand = 'nosplit' -- Shows the effects of a command incrementally in the buffer
opt.undodir = vim.fn.stdpath('state') .. '/undodir' -- Directory for undo files
opt.undofile = true -- Enable persistent undo
opt.completeopt = { 'menuone', 'popup', 'noinsert' } -- Options for completion menu
opt.winborder = 'rounded' -- Use rounded borders for windows
opt.hlsearch = true

opt.foldmethod = 'syntax'
opt.foldlevelstart = 99

--

function _G.custom_fold_text()
  local line = vim.fn.getline(vim.v.foldstart) .. ' '
  local line_count = vim.v.foldend - vim.v.foldstart + 1
  local count_str = ' ' .. line_count
  local win_width = vim.api.nvim_win_get_width(0)
  local gutter_width = vim.fn.getwininfo(vim.api.nvim_get_current_win())[1].textoff
  local available_width = win_width - gutter_width
  local line_display_width = vim.fn.strdisplaywidth(line)
  local padding = available_width - line_display_width - vim.fn.strdisplaywidth(count_str)
  if padding < 0 then padding = 1 end
  return line .. string.rep('·', padding) .. count_str
end

opt.foldtext = "v:lua.custom_fold_text()"

--

opt.wildmenu = true -- <Tab>
opt.wildmode = 'full,longest,noselect'
opt.wildoptions = 'pum,tagfile,fuzzy'

vim.opt.whichwrap = 'b,s,<,>,[,],h,l' -- Cursor left/right to move to the previous/next line

opt.modeline = true

vim.cmd.filetype('plugin indent on') -- Enable filetype detection, plugins, and indentation

opt.langmap =
    'ФИСВУАПРШОЛДЬТЩЗЙКЫЕГМЦЧНЯ;' ..
    'ABCDEFGHIJKLMNOPQRSTUVWXYZ,' ..
    'фисвуапршолдьтщзйкыегмцчня;' ..
    'abcdefghijklmnopqrstuvwxyz,' ..
    'ΑA,ΒB,ΨC,ΔD,ΕE,ΦF,ΓG,ΗH,ΙI,ΞJ,ΚK,ΛL,ΜM,ΝN,ΟO,ΠP,QQ,ΡR,ΣS,ΤT,ΘU,ΩV,WW,ΧX,ΥY,ΖZ,' ..
    'αa,βb,ψc,δd,εe,φf,γg,ηh,ιi,ξj,κk,λl,μm,νn,οo,πp,qq,ρr,σs,τt,θu,ωv,ςw,χx,υy,ζz'
vim.opt.iskeyword = '@,48-57,_,192-255,.,-' -- extra `.` and `-`

-- COLORS

vim.cmd.colorscheme('vim')

vim.api.nvim_create_autocmd('TextYankPost', {
  group    = vim.api.nvim_create_augroup('highlight_yank', {}),
  desc     = 'Hightlight selection on yank',
  pattern  = '*',
  callback = function()
    vim.highlight.on_yank { higroup = 'Visual', timeout = 500 }
  end,
})

vim.api.nvim_set_hl(0, 'TelescopeNormal', { fg = '#c0c0c0' })
vim.api.nvim_set_hl(0, 'TelescopeMatching', { bg = '#005f5f', fg = 'none' })
vim.api.nvim_set_hl(0, 'Pmenu', { bg = '#1c1c1c', fg = '#afd7ff' })
vim.api.nvim_set_hl(0, 'PmenuSel', { bg = '#585858', fg = '#afd7ff' })
vim.api.nvim_set_hl(0, 'LspDiagnosticsDefaultHint', { bg = '#1c1c1c', fg = '#5f8700' })
vim.api.nvim_set_hl(0, 'LspDiagnosticsDefaultInformation', { bg = '#1c1c1c', fg = '#0087af' })
vim.api.nvim_set_hl(0, 'LspDiagnosticsDefaultWarning', { bg = '#1c1c1c', fg = '#af875f' })
vim.api.nvim_set_hl(0, 'LspDiagnosticsDefaultError', { bg = '#1c1c1c', fg = '#af0000' })
vim.api.nvim_set_hl(0, 'TreesitterContext', { bg = '#444444' })
vim.api.nvim_set_hl(0, 'TreesitterContextLineNumber', { bg = '#444444', fg = '#ff00d7' })
vim.api.nvim_set_hl(0, 'Whitespace', { fg = '#555555' })
vim.api.nvim_set_hl(0, 'NonText', { fg = '#555555' })
vim.api.nvim_set_hl(0, 'EndOfBuffer', { fg = '#888888' })
vim.api.nvim_set_hl(0, 'LineNr', { fg = '#557755' })
vim.api.nvim_set_hl(0, 'FoldColumn', { bg = '#000000', fg = '#557755' })
vim.api.nvim_set_hl(0, 'StatusLineNC', { bg = '#444444', fg = '#000000' })
vim.api.nvim_set_hl(0, 'StatusLine', { bg = '#444444', fg = '#ffffff' })
vim.api.nvim_set_hl(0, 'ColorColumn', { bg = '#444444' }) -- vim.opt.colorcolumn = { 120 }
vim.api.nvim_set_hl(0, 'VertSplit', { fg = '#444444' })
vim.api.nvim_set_hl(0, 'TabLine', { bg = '#444444', fg = '#000000' })
vim.api.nvim_set_hl(0, 'TabLineSel', { bg = '#444444', fg = '#ffffff' })
vim.api.nvim_set_hl(0, 'TabLineFill', { bg = '#444444' })
vim.api.nvim_set_hl(0, 'CursorLine', { bg = '#6c6c6c' })
vim.api.nvim_set_hl(0, 'CursorLineNr', { bg = '#6c6c6c' })
vim.api.nvim_set_hl(0, 'CursorColumn', { bg = '#6c6c6c' })
vim.api.nvim_set_hl(0, 'Normal', {})
vim.api.nvim_set_hl(0, 'NormalFloat', {})
vim.api.nvim_set_hl(0, 'FloatBorder', {})
vim.api.nvim_set_hl(0, 'Search', { bg = '#005f5f' })
vim.api.nvim_set_hl(0, 'IncSearch', { bg = '#5f5f00' })
vim.api.nvim_set_hl(0, 'Todo', { fg = '#00fafa', bg = '#ffff00' })
vim.api.nvim_set_hl(0, 'SpellBad', { underline = true })
vim.api.nvim_set_hl(0, 'SpellCap', { underline = true })
vim.api.nvim_set_hl(0, 'SpellRare', { underline = true })
vim.api.nvim_set_hl(0, 'SpellLocal', { underline = true })
vim.api.nvim_set_hl(0, 'Folded', { bg = '#262626', fg = '#afff5f' })
vim.api.nvim_set_hl(0, 'htmlH1', { bg = '#303030', fg = '#ffffff' })
vim.api.nvim_set_hl(0, 'markdownH1Delimiter', { bg = '#303030', fg = '#ffffff' })
vim.api.nvim_set_hl(0, 'htmlH2', { bg = '#303030', fg = '#ffff00' })
vim.api.nvim_set_hl(0, 'markdownH2Delimiter', { bg = '#303030', fg = '#ffff00' })
vim.api.nvim_set_hl(0, 'htmlH3', { bg = '#303030', fg = '#5fff00' })
vim.api.nvim_set_hl(0, 'markdownH3Delimiter', { bg = '#303030', fg = '#5fff00' })
vim.api.nvim_set_hl(0, 'htmlH4', { bg = '#303030', fg = '#00d7ff' })
vim.api.nvim_set_hl(0, 'markdownH4Delimiter', { bg = '#303030', fg = '#00d7ff' })
vim.api.nvim_set_hl(0, 'htmlH5', { bg = '#303030', fg = '#ff87ff' })
vim.api.nvim_set_hl(0, 'markdownH5Delimiter', { bg = '#303030', fg = '#ff87ff' })
vim.api.nvim_set_hl(0, 'htmlH6', { bg = '#303030', fg = '#ffffff' })
vim.api.nvim_set_hl(0, 'markdownH6Delimiter', { bg = '#303030', fg = '#ffffff' })
vim.api.nvim_set_hl(0, 'htmlLink', { fg = '#5fd7ff' })
vim.api.nvim_set_hl(0, 'markdownCodeBlock', { fg = '#5fafaf' })
vim.api.nvim_set_hl(0, 'markdownCode', { fg = '#5fafaf' })
vim.api.nvim_set_hl(0, 'markdownStrike', { fg = '#5f8787' })
vim.api.nvim_set_hl(0, 'markdownItalic', { fg = '#ffffff' })
vim.api.nvim_set_hl(0, 'SignColumn', { bg = '#1c1c1c', fg = '#ff0077' })
vim.api.nvim_set_hl(0, 'diffRemoved', { fg = '#ee7777' })
vim.api.nvim_set_hl(0, 'diffAdded', { fg = '#55cc55' })
vim.api.nvim_set_hl(0, 'diffLine', { fg = '#ffff55' })
vim.api.nvim_set_hl(0, 'diffFile', { fg = '#ffff55' })
vim.api.nvim_set_hl(0, 'diffOldFile', { fg = '#ffff55' })
vim.api.nvim_set_hl(0, 'diffNewFile', { fg = '#ffff55' })
vim.api.nvim_set_hl(0, 'diffIndexLine', { fg = '#cc55cc' })
