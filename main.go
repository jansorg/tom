package main

import (
	_ "golang.org/x/text/message/catalog"

	"github.com/jansorg/gotime/go-tom/cmd"
)

//go:generate gotext -srclang=en update -out=catalog.go -lang=en,de
func main() {
	cmd.Execute()
}
