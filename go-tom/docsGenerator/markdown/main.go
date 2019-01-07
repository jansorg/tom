package main

import (
	"log"
	"os"

	"github.com/spf13/cobra/doc"

	"github.com/jansorg/gotime/go-tom/cmd"
)

// generates markdown documentation for the commandline
// to be called from the main directory
func main() {
	target := "./docs/markdown"
	if len(os.Args) == 1 {
		target = os.Args[0]
	}

	c := cmd.RootCmd
	err := doc.GenMarkdownTree(c, target)
	if err != nil {
		log.Fatal(err)
	}
}
