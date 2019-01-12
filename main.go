package main

import (
	_ "golang.org/x/text/message/catalog"

	"github.com/jansorg/tom/go-tom/cmd"
)

//go:generate gotext -srclang=en update -out=catalog.go -lang=en,de
//go:generate go-bindata -pkg tom -prefix "templates/" -o go-tom/templates.go templates/...
func main() {
	cmd.Execute()
}
