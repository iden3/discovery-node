package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/fatih/color"
	"github.com/urfave/cli"
	"github.com/iden3/discovery-research/discovery-node/config"
	"github.com/iden3/discovery-research/discovery-node/node"
	"github.com/iden3/discovery-research/discovery-node/endpoint"
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

	apiService := endpoint.Serve(config.C)
	apiService.Run(":" + strconv.Itoa(config.C.Ports.API))

	return nil
}
