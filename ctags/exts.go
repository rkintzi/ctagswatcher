package ctags

import (
	"bufio"
	"io"
	"os/exec"
	"strings"

	"github.com/rkintzi/ctagswatcher/cmdrunner"
)

func CtagsMaps() (exts []string, err error) {
	cmd := exec.Command("ctags", "--list-maps")
	err = cmdrunner.Run(cmd, func(cmdPipe io.Reader) error {
		r := bufio.NewReader(cmdPipe)
		for lineb, err := r.ReadBytes('\n'); err == nil || err == io.EOF; lineb, err = r.ReadBytes('\n') {
			line := string(lineb)
			i := strings.Index(line, "\t")
			if i == -1 {
				continue
			}
			exts = append(exts, strings.Split(line[i:len(line)], " ")...)
			if err == io.EOF {
				return nil
			} else if err != nil {
				return err
			}
		}
		return nil
	})
	return
}
