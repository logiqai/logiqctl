package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/logiqai/logiqbox/services"

	"github.com/logiqai/logiqbox/cfg"
	"github.com/urfave/cli/v2"
)

var (
	app            = cli.NewApp()
	tailNamespaces = false
	tailLabels     = false
	tailApps       = false
	tailProcs      = false
	listNamespaces = false
)

func info() {
	app.Name = "Logiq-box"
	app.Usage = "Logiq CLI Tool"
	app.Authors = []*cli.Author{
		{
			Name:  "logiq.ai",
			Email: "cli@logiq.ai",
		},
	}
	app.Version = "1.0.0"
}

func commands() {
	app.Commands = []*cli.Command{
		{
			Name:    "configure",
			Aliases: []string{"c"},
			Usage:   "Configure Logiq-box",
			Action: func(c *cli.Context) error {
				cfg.Configure()
				return nil
			},
		},
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "List of applications that you can tail",
			Action: func(c *cli.Context) error {
				args := c.Args()
				config, err := getConfig()
				if err == nil {
					services.GetApplications(config, listNamespaces, args.Slice())
				}
				return nil
			},
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:        "namespaces",
					Usage:       "-namespaces",
					Hidden:      false,
					Destination: &listNamespaces,
				},
			},
		},
		{
			Name:      "tail",
			Aliases:   []string{"t"},
			Usage:     "tail logs filtered by namespace, application, labels or process / pod name",
			ArgsUsage: "[-apps application names and/or -namespaces K8S namespace names and/or -labels K8S labels - procs process id / pod name]",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:        "namespaces",
					Usage:       "-namespaces",
					Hidden:      false,
					Destination: &tailNamespaces,
				},
				&cli.BoolFlag{
					Name:        "labels",
					Usage:       "-labels",
					Hidden:      false,
					Destination: &tailLabels,
				},
				&cli.BoolFlag{
					Name:        "apps",
					Usage:       "-apps",
					Hidden:      false,
					Destination: &tailApps,
				},
				&cli.BoolFlag{
					Name:        "process",
					Usage:       "-procs",
					Hidden:      false,
					Destination: &tailProcs,
				},
			},
			Action: func(c *cli.Context) error {
				args := c.Args()
				config, err := getConfig()
				if err == nil {
					fmt.Println("Crunching data for you...")
					services.Tail(config, tailApps, tailLabels, tailNamespaces, args.Slice())
				}
				return nil
			},
		},
		{
			Name:        "next",
			Aliases:     []string{"n"},
			Usage:       "query n",
			Description: "Get the next 'n' values for the last query or search",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "output",
					Value:   "column",
					Usage:   "Set output format to be column|json|raw",
					Hidden:  false,
					Aliases: []string{"o"},
				},
			},
			Action: func(context *cli.Context) error {
				config, _ := getConfig()
				services.GetNext(context, config)
				return nil
			},
		},
		{
			Name:      "query",
			Aliases:   []string{"q"},
			Usage:     `query "sudo cron" 2h`,
			ArgsUsage: "[application names, relative time]",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:        "debug",
					Usage:       "-debug true",
					Hidden:      true,
					DefaultText: "false",
					Aliases:     []string{"-debug", "-d", "--debug"},
				},
				&cli.StringFlag{
					Name:    "end_time",
					Value:   "10h",
					Usage:   "Relative end time.",
					Aliases: []string{"et"},
				},
				&cli.StringFlag{
					Name:    "filter",
					Usage:   "-filter 'Hostname=127.0.0.1,10.231.253.255;Message=tito*'",
					Aliases: []string{"f"},
				},
				&cli.StringFlag{
					Name:    "output",
					Value:   "column",
					Usage:   "Set output format to be column|json|raw",
					Hidden:  false,
					Aliases: []string{"o"},
				},
				&cli.StringFlag{
					Name:    "start_time",
					Value:   "10h",
					Usage:   "Relative start time.",
					Aliases: []string{"st"},
				},
				&cli.BoolFlag{
					Name:        "tail",
					Usage:       "Tail the data without paginating",
					Hidden:      false,
					DefaultText: "false",
					Aliases:     []string{"t"},
				},
			},
			Action: func(c *cli.Context) error {
				args := c.Args()
				if !args.Present() && args.Len() != 1 {
					fmt.Println("Incorrect Usage")
				} else {
					config, err := getConfig()
					apps := strings.Split(args.Get(0), " ")
					if err == nil {
						services.Query(c, config, apps, "", "QUERY")
					}
				}
				return nil
			},
			Subcommands: []*cli.Command{
				{
					Name:        "next",
					Aliases:     []string{"n"},
					Usage:       "query n",
					UsageText:   "",
					Description: "Get the next 'n' values for the last query or search",
					Action: func(context *cli.Context) error {
						config, _ := getConfig()
						services.GetNext(context, config)
						return nil
					},
				},
			},
		},
		{
			Name:      "search",
			Aliases:   []string{"s"},
			Usage:     `search sudo`,
			ArgsUsage: "[search_term, relative time]",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "st",
					Value: "10h",
					Usage: "Relative start time.",
				},

				&cli.StringFlag{
					Name:  "et",
					Value: "10h",
					Usage: "Relative end time.",
				},
				&cli.StringFlag{
					Name:  "filter",
					Usage: "--filter 'Hostname=127.0.0.1,10.231.253.255;Message=tito*'",
				},
			},
			Action: func(c *cli.Context) error {
				args := c.Args()
				if !args.Present() && args.Len() != 1 {
					fmt.Println("Incorrect Usage")
				} else {
					config, err := getConfig()
					if err == nil {
						services.Query(c, config, nil, args.Get(0), "SEARCH")
					}
				}
				return nil
			},
			Subcommands: []*cli.Command{
				{
					Name:        "next",
					Aliases:     []string{"n"},
					Usage:       "query n",
					UsageText:   "",
					Description: "Get the next 'n' values for the last query or search",
					Action: func(context *cli.Context) error {
						config, _ := getConfig()
						services.GetNext(context, config)
						return nil
					},
				},
			},
		},
	}
}

func getConfig() (*cfg.Config, error) {
	profiles, err := cfg.LoadConfig()
	if err != nil {
		fmt.Println("Logiq-box is not configured! Loading interactive configuration.")
		profiles, err = cfg.Configure()
		if err != nil {
			fmt.Println("Could not run interactive configuration")
			fmt.Printf("Create a file by copying the following content to %s", cfg.GetConfigFilePath())
			fmt.Println(cfg.GetSampleProfile().String())
			return nil, err
		}
	}
	config, err := profiles.GetDefaultProfile()
	if err != nil {
		fmt.Println("Cannot find default profile")
		fmt.Printf("Please Set a default profile in %s ", cfg.GetConfigFilePath())
		return nil, err
	}
	return config, err
}

func main() {
	info()
	commands()
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "st",
			Value: "10h",
			Usage: "Relative start time.",
		},
		&cli.StringFlag{
			Name:  "et",
			Value: "10h",
			Usage: "Relative end time.",
		},
		&cli.StringFlag{
			Name:  "debug",
			Value: "false",
			Usage: "--debug true",
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
