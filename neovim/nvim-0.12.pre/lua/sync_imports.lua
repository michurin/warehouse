local M = {}

function M.organize_imports_sync(timeout_ms)
  -- one line async version:
  -- vim.lsp.buf.code_action({ context = { only = { "source.organizeImports" } }, apply = true })
  local client = vim.lsp.get_clients({ bufnr = 0 })[1]
  local encoding = client and client.offset_encoding or 'utf-16'
  local params = {
    textDocument = vim.lsp.util.make_text_document_params(),
    range = vim.lsp.util.make_range_params(0, encoding).range,
    context = {
      diagnostics = vim.diagnostic.get(0),
      only = { 'source.organizeImports' },
    },
  }
  local tmout = timeout_ms or 1000

  local responses = vim.lsp.buf_request_sync(0, 'textDocument/codeAction', params, tmout)

  if not responses then return end

  for _, resp in pairs(responses) do
    for _, action in ipairs(resp.result or {}) do
      if action.edit then vim.lsp.util.apply_workspace_edit(action.edit, encoding) end
      if action.command then vim.exec_cmd(action.command, { bufnr = 0 }) end -- legacy: vim.lsp.buf.execute_command(action.command)
    end
  end
end

return M
