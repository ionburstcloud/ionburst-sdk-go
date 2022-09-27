package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/urfave/cli/v2"
	"gitlab.com/ionburst/ionburst-sdk-go"
	"gitlab.com/ionburst/ionburst-sdk-go/ioprogress"
)

var classification string

var SecretsCmd = &cli.Command{
	Name:    "secrets",
	Aliases: []string{"sec"},
	Usage:   "Manage Ionburst Secrets",
	Subcommands: []*cli.Command{
		{
			Name:      "get",
			Usage:     "Download a secret from Ionburst",
			ArgsUsage: "<id> <outputfile>",
			Action: func(c *cli.Context) error {

				id := c.Args().Get(0)
				if id == "" {
					return errors.New("Please specify an id for the secret to be downloaded")
				}

				file := c.Args().Get(1)
				if file == "" {
					return errors.New("Please specify a secret to download")
				}

				cli, err := ionburst.NewClientPathAndProfile(config, profile, debug)
				if err != nil {
					return err
				}

				rdr, len, err := cli.GetSecretsWithLen(id)
				if err != nil {
					return err
				}

				progressR := &ioprogress.Reader{
					Reader: rdr,
					Size:   len,
				}

				wr, err := os.Create(file)
				if err != nil {
					return err
				}

				_, err = io.Copy(wr, progressR)
				return err

			},
		},
		{
			Name:      "put",
			Usage:     "Upload a secret to Ionburst",
			ArgsUsage: "<id> <secrettoupload>",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name: "classification",
					Aliases: []string{
						"t",
					},
					Value:       "",
					Usage:       "Classification to apply to uploaded secret",
					Destination: &classification,
				},
			},
			Action: func(c *cli.Context) error {

				id := c.Args().Get(0)
				if id == "" {
					return errors.New("Please specify an id for the secret to be uploaded")
				}

				file := c.Args().Get(1)
				if file == "" {
					return errors.New("Please specify a secret to upload")
				} else if !ionburst.FileExists(file) {
					return errors.New(fmt.Sprintf("The specififed secret doesnt exist: %s", file))
				}
				stat, err := os.Stat(file)
				if err != nil {
					return err
				}
				rdr, err := os.Open(file)
				if err != nil {
					return err
				}
				progressR := &ioprogress.Reader{
					Reader: rdr,
					Size:   stat.Size(),
				}

				cli, err := ionburst.NewClientPathAndProfile(config, profile, debug)
				if err != nil {
					return err
				}

				err = cli.PutSecrets(id, progressR, c.String("classification"))
				return err
			},
		},
		{
			Name:      "delete",
			Usage:     "Delete a secret from Ionburst",
			ArgsUsage: "<id>",
			Action: func(c *cli.Context) error {
				id := c.Args().Get(0)
				if id == "" {
					return errors.New("Please specify an id for the secret to be deleted")
				}

				cli, err := ionburst.NewClientPathAndProfile(config, profile, debug)
				if err != nil {
					return err
				}

				err = cli.DeleteSecrets(id)
				if err != nil {
					return err
				}

				return nil
			},
		},
		{
			Name:      "head",
			Usage:     "Check a secret from Ionburst",
			ArgsUsage: "<id>",
			Action: func(c *cli.Context) error {
				id := c.Args().Get(0)
				if id == "" {
					return errors.New("Please specify an id for the object to be checked")
				}

				cli, err := ionburst.NewClientPathAndProfile(config, profile, debug)
				if err != nil {
					return err
				}

				size, err := cli.HeadSecretsWithLen(id)
				if err != nil {
					return err
				} else {
					fmt.Printf("Size: %d\n", size)
				}

				return nil
			},
		},
	},
}
