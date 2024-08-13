package main

import (
	"fmt"
	"os"

	"github.com/Fernando-Dourado/harness-move-project/operation"
	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
	"github.com/urfave/cli"
)

var Version = "development"

func main() {
	app := cli.NewApp()
	app.Name = "harness-move-project"
	app.Version = Version
	app.Usage = "Non-official Harness CLI to move project between organizations"
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:     "api-token",
			Usage:    "Harness API token for authentication.",
			EnvVar:   "HARNESS_API_TOKEN",
			Required: true,
		},
		cli.StringFlag{
			Name:     "account",
			Usage:    "Account identifier within Harness.",
			EnvVar:   "HARNESS_ACCOUNT",
			Required: true,
		},
		cli.StringFlag{
			Name:     "source-org",
			Usage:    "Source organization from where the project will be moved.",
			Required: true,
		},
		cli.StringFlag{
			Name:     "source-project",
			Usage:    "Source project to be moved.",
			Required: true,
		},
		cli.StringFlag{
			Name:     "target-org",
			Usage:    "Target organization to where the project will be moved.",
			Required: true,
		},
		cli.StringFlag{
			Name:     "target-project",
			Usage:    "Target project name in the target organization. Defaults to the source project name if not specified.",
			Required: false,
		},
		cli.StringFlag{
			Name:     "proxy",
			Usage:    "Proxy URL to use for network requests.",
			Required: false,
		},
	}
	app.Run(os.Args)
}

func run(c *cli.Context) {

	// Set proxy if provided
	if proxy := c.String("proxy"); proxy != "" {
		client := resty.New()
		client.SetProxy(proxy)
	}
	mv := operation.Move{
		Config: operation.Config{
			Token:   c.String("api-token"),
			Account: c.String("account"),
		},
		Source: operation.NoName{
			Org:     c.String("source-org"),
			Project: c.String("source-project"),
		},
		Target: operation.NoName{
			Org:     c.String("target-org"),
			Project: c.String("target-project"),
		},
	}

	applyArgumentRules(&mv)

	if err := mv.Exec(); err != nil {
		fmt.Println(color.RedString(fmt.Sprint("Failed:", err.Error())))
		os.Exit(1)
	}
}

func applyArgumentRules(mv *operation.Move) {
	// USE SOURCE PROJECT AS TARGET, WHEN TARGET NOT SET
	if len(mv.Target.Project) == 0 {
		mv.Target.Project = mv.Source.Project
	}
}
