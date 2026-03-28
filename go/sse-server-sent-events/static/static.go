package static

import (
	"embed"
	"io/fs"
	"log"
)

//go:embed public_html
var emgedFS embed.FS

var FS fs.FS

func init() {
	err := (error)(nil)
	FS, err = fs.Sub(emgedFS, "public_html")
	if err != nil {
		log.Panic(err)
	}
}
