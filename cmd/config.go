package main

import (
	"strings"

	"github.com/k1nky/tookhook/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use: "tookhok",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			log.SetLevel(viper.GetString("log-level"))
		},
	}
	log *logger.Logger
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func init() {
	cobra.OnInitialize(initConfig)
	initCommands()
}

func initConfig() {
	viper.SetEnvPrefix("TOOKHOOK")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.SetTypeByDefaultValue(true)
	viper.AutomaticEnv()
}

func initCommands() {
	rootCmd.PersistentFlags().String("log-level", "info", "log level")
	runCmd := &cobra.Command{
		Use:   "run",
		Run:   runServer,
		Short: "Run as server",
	}
	runCmd.Flags().StringP("listen", "s", "0.0.0.0:8080", "listen on address")
	runCmd.Flags().String("database-uri", "file://hooks.yml", "database connection string")
	runCmd.Flags().String("queue-uri", "127.0.0.1:6379", "queue connection string")
	runCmd.Flags().StringSliceP("plugins", "p", []string{}, "list of plugins")
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Run:   showVersion,
		Short: "Show version",
	})
	viper.BindPFlags(rootCmd.PersistentFlags())
	viper.BindPFlags(runCmd.Flags())
}
