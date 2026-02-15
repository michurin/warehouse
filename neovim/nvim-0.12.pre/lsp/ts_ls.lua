return {
  cmd = {'typescript-language-server', '--stdio', '--log-level', 'log'},
  single_file_support = true,
  filetypes = {
    'javascript',
    'javascriptreact',
    'javascript.jsx',
    'typescript',
    'typescriptreact',
    'typescript.tsx',
  },
}
