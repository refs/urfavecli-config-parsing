package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/imdario/mergo"

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
		&cli.StringFlag{
			Name: "load",
		},
	}

	app := &cli.App{
		Action: func(c *cli.Context) error {
			//return prettify(cfg)
			return nil
		},
		Before: func(c *cli.Context) error {
			// because urfavecli does not provide with hooks before flags are parsed, and flag parsing is not exported,
			// emulating viper config parsing unmarshaling requires some hacking. Here we are hijacking the lifecycle by:
			// 1. remembering the original flags value (after parsing by urfave)
			// 2. loading config from disk into a struct
			// 3. re-apply original values to the context. This ensures the correct order in the sense of the flag has already been parsed, by the point we reach the Before method
			//    the values are meant to be "final", therefore if they have been re-declared in the config file, essentially the config file values are discarded, as the value of the
			//    flag is set to its original; if, on the other hand, the value was NOT parsed from a flag, but we encounter it in the config file, the value will remain untouched and
			//    the config struct
			// preserve the initial values of the flags

			initialFlagValues := map[string]string{}
			for i := range flags {
				if flags[i].IsSet() {
					for j := range flags[i].Names() {
						initialFlagValues[flags[i].Names()[j]] = c.Value(flags[i].Names()[j]).(string)
					}
				}
			}

			merger := Config{}
			if err := parseConfig(c.String("load"), &merger); err != nil {
				return err
			}
			_ = prettify(merger)
			_ = prettify(cfg)
			mergo.Merge(&cfg, merger)
			_ = prettify(cfg)

			for k, v := range initialFlagValues {
				_ = c.Set(k, v)
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
