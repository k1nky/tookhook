// Package config provides the server configuration.
package config

import (
	"net"
	"os"

	"github.com/caarlos0/env/v6"
	flag "github.com/spf13/pflag"
)

// NetAddress is string [<host>]:<port> and implements interface pflag.Value
type NetAddress string

func (a NetAddress) String() string {
	return string(a)
}

func (a *NetAddress) Set(s string) error {
	host, port, err := net.SplitHostPort(s)
	if err != nil {
		return err
	}
	if len(host) == 0 {
		// use "localhost" by default
		s = "localhost:" + port
	}
	*a = NetAddress(s)
	return nil
}

func (a *NetAddress) Type() string {
	return "string"
}

// Config is the server configuration.
type Config struct {
	// Listen address and port [<host>]:<port>: environment variable `TOOKHOOK_LISTEN` or flag `-s`
	Listen NetAddress `env:"TOOKHOOK_LISTEN"`
	// DarabaseURI is database connection string: environment variable `TOOKHOOK_DATABASE_URI` or flag `-d`
	DarabaseURI string `env:"TOOKHOOK_DATABASE_URI"`
	// LogLevel is log level: environment variable `TOOKHOOK_LOG_LEVEL` or flag `-l`
	LogLevel string `env:"TOOKHOOK_LOG_LEVEL"`
	// Plugins is comma separated list of plugins: environment variable `TOOKHOOK_PLUGINS` or flag `-p`
	Plugins string `env:"TOOKHOOK_PLUGINS"`
}

func parseFromCmd(c *Config) error {
	cmd := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	listen := NetAddress("localhost:8080")
	cmd.VarP(&listen, "listen", "s", "listen address and port [<host>]:<port>")
	databaseURI := cmd.StringP("database-uri", "d", "hooks.yml", "database connection string")
	logLevel := cmd.StringP("log-level", "l", "info", "log level")
	plugins := cmd.StringP("plugins", "p", "", "comma separated list of plugins")

	if err := cmd.Parse(os.Args[1:]); err != nil {
		return err
	}

	*c = Config{
		Listen:      listen,
		DarabaseURI: *databaseURI,
		LogLevel:    *logLevel,
		Plugins:     *plugins,
	}
	return nil
}

func parseFromEnv(c *Config) error {
	if err := env.Parse(c); err != nil {
		return err
	}
	if len(c.Listen) != 0 {
		if err := c.Listen.Set(c.Listen.String()); err != nil {
			return err
		}
	}
	return nil
}

// Parse parses the server configuration.
// Environment variables take precedence over command line arguments.
func Parse(c *Config) error {
	if err := parseFromCmd(c); err != nil {
		return err
	}
	if err := parseFromEnv(c); err != nil {
		return err
	}
	return nil
}
