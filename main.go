package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "feedr",
		Usage: "Get the latest news!",
		Action: func(*cli.Context) error {
			fmt.Println("Here is your feed!")
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
