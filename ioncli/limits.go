package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"gitlab.com/ionburst/ionburst-sdk-go"
)

var LimitsCmd = &cli.Command{
	Name:  "limits",
	Usage: "Retrieve Ionburst API upload limits",
	Subcommands: []*cli.Command{
		{
			Name:  "list",
			Usage: "List all Ionburst API upload limits",
			Action: func(c *cli.Context) error {
				cli, err := ionburst.NewClientPathAndProfile(config, profile, debug)
				if err != nil {
					return err
				}
				data, err := cli.GetDataUploadLimit()
				if err != nil {
					return err
				}
				secrets, err := cli.GetSecretsUploadLimit()
				if err != nil {
					return err
				}
				fmt.Printf("Object upload limit: %d\n", data)
				fmt.Printf("Secret upload limit: %d\n", secrets)
				return nil
			},
		},
	},
}
