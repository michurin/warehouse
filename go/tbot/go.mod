module github.com/michurin/cnbot

go 1.20

replace github.com/michurin/minlog => ../minlog-next

replace github.com/michurin/systemd-env-file => ../systemd-env-file

require (
	github.com/michurin/minlog v0.0.0-00010101000000-000000000000
	github.com/michurin/systemd-env-file v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.8.3
	golang.org/x/sync v0.1.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
