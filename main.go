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
	Address string
	Port    string
}

type Proxy struct {
	Service `yaml:"service"`
}

type OCS struct {
	Service
}

// Config is a global configuration
type Config struct {
	Proxy Proxy
	OCS   OCS
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
	proxyFlags := []cli.Flag{
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

	ocsFlags := []cli.Flag{
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "ocs-port",
			EnvVars:     []string{"OCS_PORT"},
			Value:       "3. default-value-flag-declaration",
			Destination: &cfg.OCS.Port,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "ocs-addr",
			EnvVars:     []string{"OCS_ADDRESS"},
			Value:       "3. default-value-flag-declaration",
			Destination: &cfg.OCS.Address,
		}),
	}

	app := &cli.App{
		Action: func(c *cli.Context) error {
			return prettify(cfg)
		},
		Commands: []*cli.Command{
			{
				Name:   "proxy",
				Flags:  proxyFlags,
				Before: altsrc.InitInputSourceWithContext(proxyFlags, altsrc.NewYamlSourceFromFlagFunc("config-file")),
				Action: func(context *cli.Context) error {
					return prettify(cfg.Proxy)
				},
			},
			{
				Name:   "ocs",
				Flags:  ocsFlags,
				Before: altsrc.InitInputSourceWithContext(ocsFlags, altsrc.NewYamlSourceFromFlagFunc("config-file")),
				Action: func(context *cli.Context) error {
					return prettify(cfg.OCS)
				},
			},
		},
		Flags: rootFlags,
	}

	app.Run(os.Args)
}

func prettify(cfg interface{}) error {
	empJSON, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", string(empJSON))
	return nil
}
