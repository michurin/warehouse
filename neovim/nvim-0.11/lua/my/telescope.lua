function ft_args() -- TODO do not polute global scope
  local bufnr = vim.api.nvim_get_current_buf()
  local filetype = vim.bo[bufnr].filetype
  return ({ -- rg options
    ["go"]={
      "-g", "*.go",
      "-g", "!vendor",
      "-g", "!mock",
      "-g", "!mocks",
      "-g", "!.git",
      "-g", "!*_test.go",
    },
    ["js"]={"-g", "*.js"},
  })[filetype]
end

function ft_alt_args() -- TODO do not polute global scope
  local bufnr = vim.api.nvim_get_current_buf()
  local filetype = vim.bo[bufnr].filetype
  return ({ -- rg options
    ["go"]={
      "-g", "*_test.go",
      "-g", "!vendor",
      "-g", "!mock",
      "-g", "!mocks",
      "-g", "!.git",
    },
  })[filetype]
end

-- std
vim.keymap.set('n', '<space>ff', function() require('telescope.builtin').find_files(); end)
vim.keymap.set('n', '<space>fg', function() require('telescope.builtin').live_grep(); end)
vim.keymap.set('n', '<space>fb', function() require('telescope.builtin').buffers(); end)
vim.keymap.set('n', '<space>fh', function() require('telescope.builtin').help_tags(); end)
-- like /
vim.keymap.set('n', '<space>f/', function() require('telescope.builtin').current_buffer_fuzzy_find(); end)
-- z=
vim.keymap.set('n', 'z=', function() require('telescope.builtin').spell_suggest(); end)
-- m — method, d — diagnostics, l — language, t — tests, c — relative to current buffer dir
vim.keymap.set('n', '<space>fm', function() require('telescope.builtin').grep_string(); end)
vim.keymap.set('n', '<space>fd', function() require('telescope.builtin').diagnostics(); end)
vim.keymap.set('n', '<space>fl', function() require('telescope.builtin').live_grep({additional_args=ft_args, wrap_results=true}); end)
vim.keymap.set('n', '<space>ft', function() require('telescope.builtin').live_grep({additional_args=ft_alt_args, wrap_results=true}); end)
vim.keymap.set('n', '<space>fc', function() require('telescope.builtin').live_grep({["search_dirs"]={vim.fn.expand("%:p")}, wrap_results=true}); end)
-- resume
vim.keymap.set('n', '<space>fr', function() require('telescope.builtin').resume(); end)
-- std LSP (grr, gri — nvim 0.11)
vim.keymap.set('n', 'grr', function() require('telescope.builtin').lsp_references({show_line=false}); end)
vim.keymap.set('n', 'gri', function() require('telescope.builtin').lsp_implementations({show_line=false}); end)
vim.keymap.set('n', 'gd', function() require('telescope.builtin').lsp_definitions({show_line=false}); end)
vim.keymap.set('n', 'gs', function() require('telescope.builtin').lsp_definitions({show_line=false, jump_type='vsplit'}); end)
vim.keymap.set('n', 'ga', function() require('telescope.builtin').lsp_definitions({show_line=false, jump_type='tab'}); end)
vim.keymap.set('n', 'gy', function() require('telescope.builtin').lsp_type_definitions({show_line=false}); end)
-- treesitter
vim.keymap.set('n', '<space>fs', function() require('telescope.builtin').treesitter(); end)
-- all
vim.keymap.set('n', '<space>fa', function() require('telescope.builtin').builtin(); end)

-- local action_layout = require("telescope.actions.layout")
require('telescope').setup{
  defaults = {
    layout_strategy = 'vertical',
    layout_config = {
      height = 0.9,
      width = 0.9,
      preview_cutoff = 3,
    },
    mappings = {
--      n = {
--        ["<C-t>"] = action_layout.toggle_preview,
--      },
--      i = {
--        ["<C-t>"] = action_layout.toggle_preview,
--      },
    },
  },
  pickers = {
    buffers = {
      show_all_buffers = true,
      sort_lastused = true,
--      theme = "dropdown",
--      previewer = false,
      mappings = {
        i = {
          ["<c-e>"] = "delete_buffer",
        }
      }
    },
  },
}
require('telescope').load_extension('dap')
