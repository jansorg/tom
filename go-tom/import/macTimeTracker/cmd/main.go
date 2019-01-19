package main

import (
	"log"
	"os"

	"github.com/jansorg/tom/go-tom/import/macTimeTracker"
)

func main() {
	e := macTimeTracker.Import(os.Args[1], nil)
	if e != nil {
		log.Fatal(e)
	}
}
