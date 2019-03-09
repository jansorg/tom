package main

import (
	"log"
	"os"

	"github.com/spf13/cobra/doc"

	"github.com/jansorg/tom/go-tom/cmd"
)

// generates markdown documentation for the commandline
// to be called from the main directory
func main() {
	target := "./docs/man"
	if len(os.Args) == 1 {
		target = os.Args[0]
	}

	cmd := cmd.RootCmd
	header := &doc.GenManHeader{
		Title:   "tom",
		Section: "3",
	}
	err := doc.GenManTree(cmd, header, target)
	if err != nil {
		log.Fatal(err)
	}
}
