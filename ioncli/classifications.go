package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"gitlab.com/ionburst/ionburst-sdk-go"
)

var ClassificationsCmd = &cli.Command{
	Name:    "classifications",
	Aliases: []string{"class"},
	Usage:   "Manage Ionburst Classifications",
	Subcommands: []*cli.Command{
		{
			Name:  "add",
			Usage: "Add a classification for use",
			Action: func(c *cli.Context) error {
				fmt.Println("This is not currently available")
				return nil
			},
		},
		{
			Name:  "list",
			Usage: "List all available classifications",
			Action: func(c *cli.Context) error {
				cli, err := ionburst.NewClientPathAndProfile(config, profile, debug)
				if err != nil {
					return err
				}
				cls, err := cli.GetClassifications()
				if err != nil {
					return err
				}
				fmt.Printf("Classifications: %d\n", len(cls))
				if len(cls) > 0 {
					fmt.Println("\nClassification\n-----------------")
					for item := range cls {
						fmt.Println(cls[item])
					}
				}
				return nil
			},
		},
		{
			Name:  "remove",
			Usage: "Remove a classification",
			Action: func(c *cli.Context) error {
				fmt.Println("This is not currently available")
				return nil
			},
		},
	},
}
