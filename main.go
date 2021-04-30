package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/urfave/cli/v2/altsrc"

	"github.com/urfave/cli/v2"
)

// Service configures a single service
type Service struct {
	Address string `yaml:"address"`
	Port    string `yaml:"port"`
}

type Proxy struct {
	Service `yaml:"service"`
}

// Config is a global configuration
type Config struct {
	Proxy Proxy `yaml:"proxy"`
}

/*
   0. Command line flag value from user
   1. Environment variable (if specified)
   2. Configuration file (if specified)
   3. Default defined on the flag
*/

func main() {
	cfg := Config{}
	rootFlags := []cli.Flag{
		&cli.StringFlag{
			Name: "config-file",
		},
	}
	flags := []cli.Flag{
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "proxy-port",
			EnvVars:     []string{"PROXY_PORT"},
			Value:       "3. default-value-flag-declaration",
			Destination: &cfg.Proxy.Port,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "proxy-addr",
			EnvVars:     []string{"PROXY_ADDRESS"},
			Value:       "3. default-value-flag-declaration",
			Destination: &cfg.Proxy.Address,
		}),
	}

	app := &cli.App{
		Action: func(c *cli.Context) error {
			return prettify(cfg)
		},
		Commands: []*cli.Command{
			{
				Name:   "server",
				Flags:  flags,
				Before: altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc("config-file")),
				Action: func(context *cli.Context) error {
					return prettify(cfg)
				},
			},
		},
		Flags: rootFlags,
	}

	app.Run(os.Args)
}

func prettify(cfg Config) error {
	empJSON, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", string(empJSON))
	return nil
}
