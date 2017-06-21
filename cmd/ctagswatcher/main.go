package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rkintzi/ctagswatcher/cmd/ctagwatcher/cfg"
)

func abort(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func main() {
	home := os.Getenv("HOME")
	path := filepath.Join(home, ".ctagswatcher.toml")
	cfg, err := cfg.ReadConfig(path)
	if err != nil {
		abort("Can't open config file: %v", err)
	}
	fmt.Println(cfg)
}
