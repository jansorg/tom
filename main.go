package main

import (
	_ "time/tzdata"

	_ "golang.org/x/text/message/catalog"

	"github.com/jansorg/tom/go-tom/cmd"
)

var version = "0.13.0"
var commit = "unknown"
var date = "unknown"

func main() {
	cmd.Execute(version, commit, date)
}
