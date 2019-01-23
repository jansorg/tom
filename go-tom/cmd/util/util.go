package util

import (
	"fmt"
	"os"
)

func Fatal(err ...interface{}) {
	_, _ = fmt.Fprintln(os.Stderr, append([]interface{}{"Error: "}, err...)...)
	os.Exit(1)
}

func Fatalf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, "Error: "+format+"\n", args...)
	os.Exit(1)
}
