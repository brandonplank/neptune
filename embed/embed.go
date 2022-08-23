package embed

import (
	"embed"
	"github.com/leaanthony/debme"
	"log"
)

//go:embed resources
//go:embed public
var Content embed.FS

var ViewContent debme.Debme

func init() {
	log.Println("EmbedFS Setting ViewContent virtual fs")
	rootfs, err := debme.FS(Content, ".")
	if err != nil {
		log.Fatalf("EmbedFS Wrapper error:%s", err)
	}
	ViewContent, err = rootfs.FS("resources/views")
	if err != nil {
		log.Fatalf("EmbedFS Wrapper error:%s", err)
	}
}
