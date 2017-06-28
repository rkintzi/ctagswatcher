package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/rkintzi/ctagswatcher/ctags"
)

func abort(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func main() {
	tmpFile, err := ioutil.TempFile("", "catgs-reindex-file")
	if err != nil {
		abort("Can not open temporary file: %v", err)
	}
	file, err := os.OpenFile(os.Args[1], os.O_RDWR, 0666)
	if err != nil {
		abort("Can not open file '%s': %v", os.Args[1])
	}
	_, err = io.Copy(tmpFile, ctags.NewFilenameFilter(file, os.Args[2]))
	if err != nil {
		abort("Can not merge tags: %v", err)
	}
	_, err = tmpFile.Seek(0, os.SEEK_SET)
	if err != nil {
		abort("Can not seek in tempemporary file: %v", err)
	}
	_, err = file.Seek(0, os.SEEK_SET)
	if err != nil {
		abort("Can not seek in tags file: %v", err)
	}
	_, err = io.Copy(file, tmpFile)
	if err != nil {
		abort("Can not copy file: %v", err)
	}
	tmpFile.Close()
	file.Close()
	err = os.Remove(tmpFile.Name())
	if err != nil {
		abort("Can not remove temporary file: %v", err)
	}
}
