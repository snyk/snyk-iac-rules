package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:     "snyk-iac-custom-rules",
		Usage:    "Use this SDK to write, debug, test, and bundle custom rules for the Snyk IaC CLI",
		Commands: []*cli.Command{},
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}
