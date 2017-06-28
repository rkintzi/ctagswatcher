package main

import (
	"fmt"
	"io"
	"os"

	"github.com/rkintzi/ctagswatcher/ctags"
)

func abort(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func main() {
	filename := os.Args[2]
	newFile, err := os.OpenFile(os.Args[3], os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		abort("Can not open file '%s': %v", os.Args[3])
	}
	oldFile, err := os.Open(os.Args[1])
	if err != nil {
		abort("Can not open file '%s': %v", os.Args[1])
	}
	_, err = io.Copy(newFile, ctags.NewFilenameFilter(oldFile, filename))
	if err != nil {
		abort("Can not merge tags: %v", err)
	}
}
