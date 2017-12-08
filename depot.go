package main

import (
	"os"

	"github.com/teris-io/cli"
	"github.com/teris-io/depot/server"
	"github.com/teris-io/log"
	"github.com/teris-io/log/std"
)

func init() {
	std.Use(os.Stderr, log.UnsetLevel, std.DefaultFmtFun)
}

func main() {
	app := cli.New("depot server").
		WithCommand(cli.NewCommand("server", "starts the depot server").WithAction(server.Start))

	os.Exit(app.Run(os.Args, os.Stdout))
}
