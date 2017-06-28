package ctags

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

func CtagsMaps() (exts []string, err error) {
	cmd := exec.Command("ctags", "--list-maps")
	cmdPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("can't create pipe for command ctags --list-maps: %v", err)
	}
	defer func() { err = cmdPipe.Close() }()
	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("can't start process for command ctags --list-maps: %v", err)
	}
	var exts []string
	r := bufio.NewReader(cmdPipe)
	for lineb, err := r.ReadBytes('\n'); err == nil || err == io.EOF; lineb, err = r.ReadBytes('\n') {
		line := string(lineb)
		i := strings.Index(line, "\t")
		if i == -1 {
			continue
		}
		exts = append(exts, strings.Split(line[i:len(line)], " ")...)
		if err == io.EOF {
			break
		}
	}
	if err != io.EOF {
		return nil, fmt.Errorf("can't read from process: %v", err)
	}
	err = cmd.Wait()
	if err != io.EOF {
		return nil, fmt.Errorf("can't wait for process: %v", err)
	}
	return
}
