package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/rkintzi/ctagswatcher/ctags"
)

var cmd *exec.Cmd

func abort(format string, args ...interface{}) {
	if cmd != nil {
		cmd.Wait()
	}
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
	cmd = exec.Command("ctags", "-f", "-", filename)
	cmdPipe, err := cmd.StdoutPipe()
	if err != nil {
		abort("Can not create pipe to subprocess: %v", err)
	}
	err = cmd.Start()
	if err != nil {
		abort("Can not start ctags: %v", err)
	}
	err = ctags.MergeTags(newFile, ctags.NewFilenameFilter(oldFile, filename), cmdPipe)
	if err != nil {
		abort("Can not merge tags: %v", err)
	}
	cmd.Wait()
}
