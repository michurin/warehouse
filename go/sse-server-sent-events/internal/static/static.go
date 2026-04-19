package static

import (
	"embed"
	"io/fs"
)

//go:embed public_html
var emgedFS embed.FS

var FS fs.FS

func init() {
	err := (error)(nil)
	FS, err = fs.Sub(emgedFS, "public_html")
	if err != nil {
		panic(err)
	}
}
