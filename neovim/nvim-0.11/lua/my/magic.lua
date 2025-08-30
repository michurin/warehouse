-- stolen from https://yobibyte.github.io/vim.html

vim.keymap.set("n", "<space>c", function()
  vim.ui.input({}, function(c)
    if c and c~="" then
      vim.cmd("noswapfile vnew")
      vim.bo.buftype = "nofile"
      vim.bo.bufhidden = "wipe"
      vim.api.nvim_buf_set_lines(0, 0, -1, false, vim.fn.systemlist(c))
    end
  end)
end)
