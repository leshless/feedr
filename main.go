package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

// HandlerWrapper is a middleware that performs auxilary actions before and after calling a Hanlder.
func Wrapper(handler cli.ActionFunc) cli.ActionFunc {
	return func(cctx *cli.Context) error {
		err := OpenConfigFiles()
		defer CloseConfigFiles()
		if err != nil {
			return err
		}

		err = handler(cctx)
		return err
	}
}

func main() {
	app := &cli.App{
		Name:  "feedr",
		Usage: "Get the latest news!",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "lang",
				Aliases: []string{"l"},
				Value:   "english",
				Usage:   "language for the greeting",
			},
		},
		Action: Wrapper(func(ctx *cli.Context) error {
			fmt.Println("test")
			return nil
		}),
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("feedr: %s\n", err.Error())
	}
}
