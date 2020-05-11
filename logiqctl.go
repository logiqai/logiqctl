package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/logiqai/logiqctl/services"

	"github.com/logiqai/logiqctl/cfg"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var (
	app            = cli.NewApp()
	tailNamespaces = []string{}
	tailLabels     = []string{}
	tailApps       = []string{}
	tailProcs      = []string{}
	listNamespaces = false
)

func info() {
	app.Name = "Logiqctl"
	app.Usage = "LOGIQ command line toolkit"
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
			Usage:   "Configure logiqctl",
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
			Name:    "tail",
			Aliases: []string{"t"},
			Usage:   "tail logs filtered by namespace, application, labels or process / pod name",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:   "namespaces",
					Usage:  "Comma separated namespaces from which we tail the data e.g. -namespace foo,bar",
					Hidden: false,
				},
				&cli.StringFlag{
					Name:   "labels",
					Usage:  "Comma separated K8S labels to match. Allowed label K/V separator \":\" OR \"=\". e.g. -labels app:some-app",
					Hidden: false,
				},
				&cli.StringFlag{
					Name:   "apps",
					Usage:  "Comma separated application names",
					Hidden: false,
				},
				&cli.StringFlag{
					Name:   "process",
					Usage:  "Command separated Process/Pod names",
					Hidden: false,
				},
				&cli.StringFlag{
					Name:    "output",
					Value:   "column",
					Usage:   "Set output format to be column|json|raw",
					Hidden:  false,
					Aliases: []string{"o"},
				},
			},
			Action: func(c *cli.Context) error {
				var namespaces, applications, procs, labels = []string{}, []string{}, []string{}, []string{}
				if v := c.Value("namespaces").(string); v != "" {
					namespaces = strings.Split(v, ",")
				}
				if v := c.Value("apps").(string); v != "" {
					applications = strings.Split(v, ",")
				}
				if v := c.Value("labels").(string); v != "" {
					labels = strings.Split(v, ",")
				}
				if v := c.Value("process").(string); v != "" {
					procs = strings.Split(v, ",")
				}

				config, err := getConfig()
				args := c.Args()

				log.Debugln(namespaces, labels, applications, procs, args.Slice())
				log.Debugln(len(namespaces), len(labels), len(applications), len(procs), len(args.Slice()))
				if err == nil {
					fmt.Println("Crunching data for you...")
					services.Tail(c, config, namespaces, labels, applications, procs, args.Slice())
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
			Name:        "query",
			Aliases:     []string{"q"},
			Usage:       `logiqctl query`,
			Description: "Query application logs",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:        "debug",
					Usage:       "-debug true",
					Hidden:      true,
					DefaultText: "false",
					Aliases:     []string{"-debug", "-d", "--debug"},
				},
				&cli.StringFlag{
					Name:    "filter",
					Usage:   "Filter expression e.g. 'Hostname=127.0.0.1,10.231.253.255;Message=tito*'",
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

				config, err := getConfig()
				if err == nil {
					services.Query(c, config, "QUERY")
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
		//{
		//	Name:      "search",
		//	Aliases:   []string{"s"},
		//	Usage:     `search sudo`,
		//	ArgsUsage: "[search_term, relative time]",
		//	Flags: []cli.Flag{
		//		&cli.StringFlag{
		//			Name:  "st",
		//			Value: "10h",
		//			Usage: "Relative start time.",
		//		},
		//
		//		&cli.StringFlag{
		//			Name:  "et",
		//			Value: "10h",
		//			Usage: "Relative end time.",
		//		},
		//		&cli.StringFlag{
		//			Name:  "filter",
		//			Usage: "--filter 'Hostname=127.0.0.1,10.231.253.255;Message=tito*'",
		//		},
		//	},
		//	Action: func(c *cli.Context) error {
		//		args := c.Args()
		//		if !args.Present() && args.Len() != 1 {
		//			fmt.Println("Incorrect Usage")
		//		} else {
		//			config, err := getConfig()
		//			if err == nil {
		//				services.Query(c, config, nil, args.Get(0), "SEARCH")
		//			}
		//		}
		//		return nil
		//	},
		//	Subcommands: []*cli.Command{
		//		{
		//			Name:        "next",
		//			Aliases:     []string{"n"},
		//			Usage:       "query n",
		//			UsageText:   "",
		//			Description: "Get the next 'n' values for the last query or search",
		//			Action: func(context *cli.Context) error {
		//				config, _ := getConfig()
		//				services.GetNext(context, config)
		//				return nil
		//			},
		//		},
		//	},
		//},
	}
}

func getConfig() (*cfg.Config, error) {
	profiles, err := cfg.LoadConfig()
	if err != nil {
		fmt.Println("logiqctl is not configured! Loading interactive configuration.")
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
	log.SetLevel(log.ErrorLevel)
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
