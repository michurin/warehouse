require'treesitter-context'.setup{
  enable = true, -- autocmd VimEnter * TSContextEnable
  throttle = true, -- may improve performance
  max_lines = 10, -- 0 — no limit
  mode = 'cursor', -- 'topline', 'cursor',
  -- separator = '┄',
  patterns = { -- lua print(vim.inspect(require'nvim-treesitter.ts_utils'.get_node_at_cursor():type()))
    default = {'class', 'function', 'method', 'for', 'while', 'if', 'switch', 'case'},
    yaml = {'block_mapping_pair'},
    json = {'pair'},
    toml = {'table', 'bare_key'},
    markdown = {'section'},
    go = {'import_declaration', 'assignment_statement', 'short_var_declaration', 'defer_statement', 'func_literal'}, -- func_literal for anonymous functions
  },
}

-- ------------------

require'nvim-treesitter.configs'.setup {
  textobjects = {
    select = {
      enable = true,

      -- Automatically jump forward to textobj, similar to targets.vim
      lookahead = true,

      keymaps = {
        -- You can use the capture groups defined in textobjects.scm
        ["af"] = "@function.outer",
        ["if"] = "@function.inner",
        ["ac"] = "@class.outer",
        ["ic"] = "@class.inner",
      },
    },
    swap = {
      enable = true,
      swap_next = {
        ["<space>a"] = "@parameter.inner",
      },
      swap_previous = {
        ["<space>A"] = "@parameter.inner",
      },
    },
    move = {
      enable = true,
      set_jumps = true, -- whether to set jumps in the jumplist
      goto_next_start = {
        ["]m"] = "@class.outer",
        ["]]"] = "@function.outer",
      },
      goto_next_end = {
        ["]M"] = "@class.outer",
        ["]["] = "@function.outer",
      },
      goto_previous_start = {
        ["[m"] = "@class.outer",
        ["[["] = "@function.outer",
      },
      goto_previous_end = {
        ["[M"] = "@class.outer",
        ["[]"] = "@function.outer",
      },
    },
  },
  refactor = {
    highlight_definitions = {
      enable = false,
      -- Set to false if you have an `updatetime` of ~100.
      clear_on_cursor_move = false,
    },
    highlight_current_scope = {
      -- Works wired for me, even with custom TSCurrentScope
      enable = false,
    },
    smart_rename = {
      enable = true,
      keymaps = {
        smart_rename = "gR",
      },
    },
    navigation = {
      enable = true,
      keymaps = {
        -- goto_definition = "gnd",
        -- list_definitions = "gnD",
        -- list_definitions_toc = "gO",
        goto_next_usage = "]u",
        goto_previous_usage = "[u",
      },
    },
  },
  incremental_selection = {
    enable = true,
    keymaps = {
      init_selection = "gww", -- set to `false` to disable one of the mappings
      node_incremental = "gwi",
      scope_incremental = "gwj",
      node_decremental = "gwd",
    },
  },
}
