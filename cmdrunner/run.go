package cmdrunner

import (
	"io"
	"os/exec"
)

type ProcFunc func(r io.Reader) error

func Run(cmd *exec.Cmd, pf ProcFunc) error {
	p, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	err = cmd.Start()
	if err != nil {
		p.Close()
		return err
	}
	err = pf(p)
	if err != nil {
		p.Close()
		return err
	}
	return cmd.Wait()
}
