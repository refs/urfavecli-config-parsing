package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

func main() {
	flags := []cli.Flag{
		altsrc.NewIntFlag(&cli.IntFlag{Name: "test", Value: 0}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "test_two", Value: "nobody"}),
		&cli.StringFlag{Name: "load"},
	}

	app := &cli.App{
		Action: func(c *cli.Context) error {
			fmt.Printf("--test value.*default: %v\n", c.Int("test"))
			fmt.Printf("--test value.*default: %v\n", c.String("test_two"))
			return nil
		},
		Before: altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc("load")),
		Flags: flags,
	}

	app.Run(os.Args)
}