// Пакет config представляет инструменты для работы с конфигурациями сервера и агента
package config

import (
	"net"
	"os"

	"github.com/caarlos0/env/v6"
	flag "github.com/spf13/pflag"
)

// NetAddress строка вида [<хост>]:<порт> и реализует интерфейс pflag.Value
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
		// если не указан хост, то используем localhost по умолчанию
		s = "localhost:" + port
	}
	*a = NetAddress(s)
	return nil
}

func (a *NetAddress) Type() string {
	return "string"
}

// Config конфигурация агента
type Config struct {
	// адрес и порт сервиса: переменная окружения ОС `TOOKHOOK_LISTEN` или флаг `-s`
	Listen NetAddress `env:"TOOKHOOK_LISTEN"`
	// адрес подключения к базе данных: переменная окружения ОС `TOOKHOOK_DATABASE_URI` или флаг `-d`
	DarabaseURI string `env:"TOOKHOOK_DATABASE_URI"`
	// уровень логирования: переменная окружения ОС `TOOKHOOK_LOG_LEVEL` или флаг `-l`
	LogLevel string `env:"TOOKHOOK_LOG_LEVEL"`
	Plugins  string `env:"TOOKHOOK_PLUGINS"`
}

func parseFromCmd(c *Config) error {
	cmd := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	listen := NetAddress("localhost:8080")
	cmd.VarP(&listen, "listen", "s", "адрес и порт сервиса")
	databaseURI := cmd.StringP("database-uri", "d", "hooks.yml", "адрес подключения к базе данных")
	logLevel := cmd.StringP("log-level", "l", "info", "уровень логирования")
	plugins := cmd.StringP("plugins", "p", "", "плагины")

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

// Parse разбирает настройки из аргументов командной строки
// и переменных окружения. Переменные окружения имеют более высокий
// приоритет, чем аргументы.
func Parse(c *Config) error {
	if err := parseFromCmd(c); err != nil {
		return err
	}
	if err := parseFromEnv(c); err != nil {
		return err
	}
	return nil
}
