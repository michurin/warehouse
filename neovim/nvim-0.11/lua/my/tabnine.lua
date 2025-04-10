require('tabnine').setup({
  disable_auto_comment=false, -- true,
  accept_keymap="<Tab>",
  dismiss_keymap = "<C-]>",
  debounce_ms = 100, -- 800,
  suggestion_color = {gui = "#80ff80", cterm = 120},
  codelens_color = {gui = "#80ff80", cterm = 120},
  codelens_enabled = true,
  exclude_filetypes = {"TelescopePrompt", "NvimTree"},
  log_file_path = nil, -- absolute path to Tabnine log file
})
