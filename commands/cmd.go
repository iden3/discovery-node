package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/fatih/color"
	"github.com/iden3/discovery-node/config"
	"github.com/iden3/discovery-node/endpoint"
	"github.com/iden3/discovery-node/node"
	"github.com/urfave/cli"
)

// NodeCommands contain the cli commands
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

	nodesrv, err := node.RunNode()
	if err != nil {
		color.Red(err.Error())
		os.Exit(0)
	}

	apiService := endpoint.Serve(config.C, *nodesrv)
	apiService.Run(":" + strconv.Itoa(config.C.Ports.API))

	return nil
}
