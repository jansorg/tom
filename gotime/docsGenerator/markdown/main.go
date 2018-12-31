package main

import (
	"log"

	"github.com/spf13/cobra/doc"

	"github.com/jansorg/gotime/gotime/cmd"
)

// generates markdown documentation for the commandline
// to be called from the main directory
func main() {
	c := cmd.RootCmd
	err := doc.GenMarkdownTree(c, "./docs/markdown/")
	if err != nil {
		log.Fatal(err)
	}
}
