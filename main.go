package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/go-yaml/yaml"
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
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:        "proxy-port",
			EnvVars:     []string{"PROXY_PORT"},
			Value:       "3. default-value-flag-declaration",
			Destination: &cfg.Proxy.Port,
		},
		&cli.StringFlag{
			Name:        "proxy-addr",
			EnvVars:     []string{"PROXY_ADDRESS"},
			Value:       "3. default-value-flag-declaration",
			Destination: &cfg.Proxy.Address,
		},
		&cli.StringSliceFlag{
			Name: "patch",
		},
		&cli.StringFlag{
			Name: "config-file",
		},
	}

	app := &cli.App{
		Action: func(c *cli.Context) error {
			//return prettify(cfg)
			return nil
		},
		Before: func(c *cli.Context) error {
			// by the time we reach Before, flags have already been parsed
			_ = prettify(cfg)
			if c.IsSet("config-file") {
				_ = parseConfig(c.Value("config-file").(string), &cfg)
				if c.IsSet("patch") {
					v := c.Value("patch").(cli.StringSlice)
					for _, v := range v.Value() {
						if len(strings.Split(v, "=")) != 2 {
							break
						}
						left := strings.Split(v, "=")[0]
						right := strings.Split(v, "=")[1]
						if err := c.Set(left, right); err != nil {
							return err
						}
					}
				}
				_ = prettify(cfg)
			}
			return nil
		},
		Flags: flags,
	}

	app.Run(os.Args)
}

func parseConfig(file string, cfg *Config) error {
	contents, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(contents, &cfg)
}

func prettify(cfg Config) error {
	empJSON, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", string(empJSON))
	return nil
}
