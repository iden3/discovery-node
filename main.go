package main

import (
	"os"

	"github.com/fatih/color"
	"github.com/iden3/discovery-node/commands"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "config"},
	}
	app.Commands = []cli.Command{}
	app.Commands = append(app.Commands, commands.NodeCommands...)
	err := app.Run(os.Args)
	if err != nil {
		color.Red(err.Error())
	}
}
