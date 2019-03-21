package commands

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/iden3/discovery-node/config"
	"github.com/urfave/cli"
)

var NodeCommands = []cli.Command{
	{
		Name:    "start",
		Aliases: []string{},
		Usage:   "start the server",
		Action:  cmdStart,
	},
}

func cmdStart(c *cli.Context) error {
	if err := config.MustRead(c); err != nil {
		color.Red(err.Error())
		os.Exit(0)
	}
	fmt.Println("c", config.C)

	core.RunNode()

	return nil
}
