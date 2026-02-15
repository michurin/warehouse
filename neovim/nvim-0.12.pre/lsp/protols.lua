-- NOT https://github.com/neovim/nvim-lspconfig
-- go install github.com/lasorda/protobuf-language-server@master
return {
  cmd = { 'protobuf-language-server' },
  filetypes = { 'proto' },
  root_markers = { '.git' },
  staticcheck = true,
  single_file_support = true,
  settings = {
    ['additional-proto-dirs'] = {},
  },
}

