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
		err := handler(cctx)
		return err
	}
}

func main() {
	app := &cli.App{
		Name:  "feedr",
		Usage: "Get the latest news!",
		Action: Wrapper(func(ctx *cli.Context) error {
			err := ReadConfigFile(SOURCES_PATH, &sources)
			if err != nil {
				return err
			}

			sources = append(sources, Source{"coca", "cola"})

			err = WriteConfigFile(SOURCES_PATH, sources)
			if err != nil {
				return err
			}

			return nil
		}),
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("feedr: %s\n", err.Error())
	}
}
