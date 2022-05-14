package embed

import (
	"embed"
	"github.com/leaanthony/debme"
)

//go:embed resources
//go:embed public
var Content embed.FS

var ViewContent debme.Debme

func init() {
	rootfs, _ := debme.FS(Content, ".")
	ViewContent, _ = rootfs.FS("resources/views")
}
