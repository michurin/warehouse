require('dap-go').setup()
require('dapui').setup({
  layouts = {
    {
      elements = {'repl', 'scopes'},
      size = 0.25,
      position = 'bottom',
    }
  },
  controls = {
    element = "repl",
    enabled = false,
    icons = {pause="Pause",play="Play",step_into="Info",step_over="Over",step_out="Out",step_back="Back",run_last="Run",terminate="Kill"},
  }
})

vim.keymap.set('n', '<space>dp', function() require'dap'.toggle_breakpoint(); end)
vim.keymap.set('n', '<space>dP', function() require'dap'.set_breakpoint(vim.fn.input('Breakpoint condition: ')); end)
vim.keymap.set('n', '<space>dt', function() require'dap-go'.debug_test(); end)
vim.keymap.set('n', '<space>dc', function() require'dap'.continue(); end)
vim.keymap.set('n', '<space>dn', function() require'dap'.step_over(); end)
vim.keymap.set('n', '<space>di', function() require'dap'.step_into(); end)
vim.keymap.set('n', '<space>do', function() require'dap'.step_out(); end)
vim.keymap.set('n', '<space>dv', function() require'dapui'.float_element('scopes', {enter=1}); end)
vim.keymap.set('n', '<space>dr', function() require'dapui'.float_element('repl', {enter=1}); end)
vim.keymap.set('n', '<space>du', function() require'dapui'.toggle(); end)
vim.keymap.set('n', '<space>sc', function() require'telescope'.extensions.dap.commands(); end)
vim.keymap.set('n', '<space>sC', function() require'telescope'.extensions.dap.configurations(); end)
vim.keymap.set('n', '<space>sp', function() require'telescope'.extensions.dap.list_breakpoints({show_line=false}); end)
vim.keymap.set('n', '<space>sv', function() require'telescope'.extensions.dap.variables(); end)
vim.keymap.set('n', '<space>sf', function() require'telescope'.extensions.dap.frames(); end)

