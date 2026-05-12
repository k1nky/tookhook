package main

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:     "tookhook",
		Short:   "Webhook processing service",
		Version: buildVersion,
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run:   showVersion,
	}

	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run the webhook processing service",
		Run:   runServer,
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	// Persistent flags (available for all commands)
	rootCmd.PersistentFlags().String("log-level", "info", "Log level (debug, info, warn, error)")
	rootCmd.PersistentFlags().String("config", "config.yaml", "Path to configuration file")

	// Run command flags
	runCmd.Flags().StringP("listen", "l", ":8080", "HTTP server listen address")
	runCmd.Flags().String("redis-addr", "localhost:6379", "Redis server address")
	runCmd.Flags().Int("redis-db", 0, "Redis database number")
	runCmd.Flags().Int("concurrency", 10, "Number of concurrent workers")

	// Bind flags to viper
	viper.BindPFlag("log.level", rootCmd.PersistentFlags().Lookup("log-level"))
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	viper.BindPFlag("server.listen", runCmd.Flags().Lookup("listen"))
	viper.BindPFlag("queue.addr", runCmd.Flags().Lookup("redis-addr"))
	viper.BindPFlag("queue.db", runCmd.Flags().Lookup("redis-db"))
	viper.BindPFlag("queue.concurrency", runCmd.Flags().Lookup("concurrency"))

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(runCmd)
}

func initConfig() {
	// Environment variable prefix
	viper.SetEnvPrefix("TOOKHOOK")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("log.level", "info")
	viper.SetDefault("server.listen", ":8080")
	viper.SetDefault("queue.addr", "localhost:6379")
	viper.SetDefault("queue.db", 0)
	viper.SetDefault("queue.concurrency", 10)

	// Read config file if specified
	cfgFile := viper.GetString("config")
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err == nil {
			// Config file found and successfully parsed
		}
	}
}
