Command line:

```
go run . methods example E | jq
go run . strings example | jq
go run . constructors example E | jq
```

NeoVim:

```
:source %
:lua =M('example', 'E')
:lua =M('example', 'EE')
:lua =M('example', 'ee')
:lua =M(vim.fn.expand('%:~:.:h'), vim.fn.expand('<cword>'))
:lua =S('example', 'O2')
```
