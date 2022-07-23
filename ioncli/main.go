package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"gitlab.com/ionburst/ionburst-sdk-go"
	"gitlab.com/ionburst/ionburst-sdk-go/ioprogress"
)

var debug bool
var profile string
var config string

var (
	Version    string
	Build      string
	APIVersion string
)

func main() {
	_ = ionburst.IonConfig{}
	var classification string

	cli.AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}
USAGE:
   {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
   {{if len .Authors}}
AUTHOR:
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
COMMANDS:
{{range .Commands}}{{if not .HideHelp}}   {{join .Names ", "}}{{ "\t"}}{{.Usage}}{{ "\n" }}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
GLOBAL OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}{{if .Copyright }}
COPYRIGHT:
   {{.Copyright}}
   {{end}}{{if .Version}}
VERSION:
   {{.Version}}
   {{end}}
`

	cli.VersionFlag = &cli.BoolFlag{
		Name:  "version",
		Usage: "Show Version Information",
	}

	app := &cli.App{
		Name:      "ioncli",
		Version:   fmt.Sprintf("%s [ Build: %s | API: %s ]", Version, Build, APIVersion),
		Copyright: "(C) Ionburst Limited",
		Usage:     "Command Line Utility for Ionburst",

		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name: "debug",
				Aliases: []string{
					"v",
				},
				Value:       false,
				Usage:       "Show debug/verbose output",
				Destination: &debug,
			},
			&cli.StringFlag{
				Name: "profile",
				Aliases: []string{
					"p",
				},
				Value:       ionburst.DefaultIonburstCredentialsProfileName,
				Usage:       "The credentials profile to use",
				Destination: &profile,
			},
			&cli.StringFlag{
				Name: "config-file",
				Aliases: []string{
					"c",
				},
				Value:       ionburst.GetDefaultIonburstConfigPath(),
				Usage:       "The Ionburst credentials filepath to use",
				Destination: &profile,
			},
		},
		Commands: []*cli.Command{
			ClassificationsCmd,
			SecretsCmd,
			{
				Name:      "get",
				Usage:     "Download an object from Ionburst",
				ArgsUsage: "<id> <outputfile>",
				Action: func(c *cli.Context) error {

					id := c.Args().Get(0)
					if id == "" {
						return errors.New("Please specify an id for the object to be downloaded")
					}

					file := c.Args().Get(1)
					if file == "" {
						return errors.New("Please specify an object to download")
					}

					cli, err := ionburst.NewClientPathAndProfile(config, profile, debug)
					if err != nil {
						return err
					}

					rdr, len, err := cli.GetWithLen(id)
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
				Usage:     "Upload an object to Ionburst",
				ArgsUsage: "<id> <objecttoupload>",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: "classification",
						Aliases: []string{
							"t",
						},
						Value:       "",
						Usage:       "Classification to apply to uploaded object",
						Destination: &classification,
					},
				},
				Action: func(c *cli.Context) error {

					id := c.Args().Get(0)
					if id == "" {
						return errors.New("Please specify an id for the object to be uploaded")
					}

					file := c.Args().Get(1)
					if file == "" {
						return errors.New("Please specify an object to upload")
					} else if !ionburst.FileExists(file) {
						return errors.New(fmt.Sprintf("The specififed object doesnt exist: %s", file))
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

					err = cli.Put(id, progressR, c.String("classification"))
					return err
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete an object from Ionburst",
				ArgsUsage: "<id>",
				Action: func(c *cli.Context) error {
					id := c.Args().Get(0)
					if id == "" {
						return errors.New("Please specify an id for the object to be deleted")
					}

					cli, err := ionburst.NewClientPathAndProfile(config, profile, debug)
					if err != nil {
						return err
					}

					err = cli.Delete(id)
					if err != nil {
						return err
					}

					return nil
				},
			},
			{
				Name:      "head",
				Usage:     "Check an object from Ionburst",
				ArgsUsage: "<id>",
				Action: func(c *cli.Context) error {
					id := c.Args().Get(0)
					if id == "" {
						return errors.New("Please specify an id for the object to be deleted")
					}

					cli, err := ionburst.NewClientPathAndProfile(config, profile, debug)
					if err != nil {
						return err
					}

					size, err := cli.HeadWithLen(id)
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

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
