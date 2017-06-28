package main

import (
	"fmt"
	"io"
	"io/ioutil"
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
	tagsfilename := os.Args[1]
	filename := os.Args[2]
	tmpFile, err := ioutil.TempFile("", "catgs-reindex-file")
	if err != nil {
		abort("Can not open temporary file: %v", err)
	}
	tagsfile, err := os.OpenFile(tagsfilename, os.O_RDWR, 0666)
	if err != nil {
		abort("Can not open file '%s': %v", tagsfilename)
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
	err = ctags.MergeTags(tmpFile, ctags.NewFilenameFilter(tagsfile, filename), cmdPipe)
	if err != nil {
		abort("Can not merge tags: %v", err)
	}
	cmd.Wait()

	_, err = tmpFile.Seek(0, os.SEEK_SET)
	if err != nil {
		abort("Can not seek in tempemporary file: %v", err)
	}
	_, err = tagsfile.Seek(0, os.SEEK_SET)
	if err != nil {
		abort("Can not seek in tags file: %v", err)
	}
	_, err = io.Copy(tagsfile, tmpFile)
	if err != nil {
		abort("Can not copy file: %v", err)
	}
	tagsfile.Close()
	tmpFile.Close()
	err = os.Remove(tmpFile.Name())
	if err != nil {
		abort("Can not remove temporary file: %v", err)
	}
}
