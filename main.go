package main

import (
	"fmt"
	"os"

	"github.com/mmcdole/gofeed"
	"github.com/urfave/cli/v2"
)

// Representaion of RSS feed source
type Source struct {
	Name string `yaml:"name"`
	Url  string `yaml:"url"`
}

// User settings
type Config struct {
	Dummy bool
}

var config Config
var sources []Source

// Results of the single call to RSS feed. Err may be a net error or a parsing error.
type ParseResult struct {
	Name string
	Feed *gofeed.Feed
	Err  error
}

// HandlerWrapper is a middleware that performs auxilary actions before and after calling a Hanlder.
func Wrapper(handler cli.ActionFunc) cli.ActionFunc {
	return func(cctx *cli.Context) error {
		err := ReadConfigFile(CONFIG_PATH, &config)
		if err != nil {
			return err
		}
		err = ReadConfigFile(SOURCES_PATH, &sources)
		if err != nil {
			return err
		}
		err = handler(cctx)
		return err
	}
}

func MainHandler(cctx *cli.Context) error {
	results := FetchAndParse(sources)
	for _, res := range results {
		if res.Err == nil {
			fmt.Printf("%v:%v\n", res.Name, res.Feed)
		}
	}
	return nil
}

func ListHandler(cctx *cli.Context) error {
	results := FetchAndParse(sources)
	for _, res := range results {
		if res.Err == nil {
			fmt.Printf("%v:%v\n", res.Name, res.Feed)
		}
	}
	return nil
}

func main() {
	app := &cli.App{
		Name:   "feedr",
		Usage:  "Get the latest news.",
		Action: Wrapper(MainHandler),
		Commands: []*cli.Command{
			{
				Name:   "list",
				Usage:  "Manage the list of your feeds.",
				Action: Wrapper(ListHandler),
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:    "add",
						Aliases: []string{"a"},
						Usage:   "Adds feed to a list.",
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("feedr: %s\n", err.Error())
	}
}
