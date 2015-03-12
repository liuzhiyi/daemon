package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	flDaemon = flag.Bool("-daemon", false, "enable daemon mode")
	flHost   = flag.String("-host", "", "set host")
)

func init() {
	flag.Usage = func() {
		fmt.Fprint(os.Stdout, "Usage: daemon [OPTIONS] COMMAND [arg...]\n\nOptions:\n")
		flag.PrintDefaults()
		flag.CommandLine.SetOutput(os.Stdout)
		help := `\nCommands:\n
                     install:    register service
                     uninstall:  remove service`
		fmt.Fprint(os.Stdout, help)
	}
}
