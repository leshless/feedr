package main

import (
	"fmt"
	"os"
	"strings"

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

	var items strings.Builder
	var errors strings.Builder

	for _, res := range results {
		if res.Err != nil {
			errors.WriteString(fmt.Sprintf("%s\n", res.Err.Error()))
		} else {
			for _, item := range res.Feed.Items {

			}
		}
	}

	return nil
}

func ListHandler(cctx *cli.Context) error {
	if cctx.Bool("add") {
		args := cctx.Args()
		if args.Len() != 2 {
			return fmt.Errorf("bad args")
		}

		name := args.Get(0)
		url := args.Get(1)
		sources = append(sources, Source{name, url})

		err := WriteConfigFile(SOURCES_PATH, sources)
		if err != nil {
			return err
		}
	} else {
		for _, source := range sources {
			fmt.Printf("%s (url: %s)\n", source.Name, source.Url)
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
					&cli.BoolFlag{
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
