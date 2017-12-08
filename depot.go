package main

import (
	"github.com/teris-io/cli"
	"github.com/teris-io/log"
	"github.com/teris-io/log/std"
	"os"
)

func init() {
	std.Use(os.Stderr, log.UnsetLevel, std.DefaultFmtFun)
}

func main() {
	start := cli.NewCommand("start", "start the depot server").
		WithAction(func(args []string, options map[string]string) int {
			log.Level(log.InfoLevel).Log("test")
			// FIXME
			return 0
		})

	app := cli.New("depot server").
		WithCommand(start)

	os.Exit(app.Run(os.Args, os.Stdout))
}
