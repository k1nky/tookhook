package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/k1nky/tookhook/internal/adapter/database"
	httphandler "github.com/k1nky/tookhook/internal/adapter/http"
	"github.com/k1nky/tookhook/internal/adapter/pluginmanager"
	"github.com/k1nky/tookhook/internal/adapter/taskq"
	"github.com/k1nky/tookhook/internal/entity/tasks"
	"github.com/k1nky/tookhook/internal/service/hooker"
	"github.com/k1nky/tookhook/internal/service/monitor"
	"github.com/k1nky/tookhook/internal/service/ruler"
	"github.com/k1nky/tookhook/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	LoggerName         = "tookhook"
	LoggerDefaultLevel = "debug"
)

func main() {
	log = logger.New(LoggerName)
	log.SetLevel(LoggerDefaultLevel)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		log.Errorf("%s", err)
		os.Exit(1)
	}
}

func runServer(cmd *cobra.Command, args []string) {
	ctx := cmd.Root().Context()

	// load plugins
	pm := pluginmanager.New(log)
	for _, v := range viper.GetStringSlice("plugins") {
		_, name := path.Split(v)
		if err := pm.Load(ctx, name, v); err != nil {
			log.Errorf("plugins: %s", err)
			return
		}
	}
	pm.Run(ctx)

	// open rules store
	store := database.New(viper.GetString("database-uri"), log.Sub("store"))
	if err := store.Open(ctx); err != nil {
		log.Errorf("opening db: %s", err)
		return
	}
	ruleService := ruler.New(pm, store, log.Sub("ruler"))
	if err := ruleService.Load(ctx); err != nil {
		log.Errorf("loading rules: %s", err)
		return
	}
	tq := taskq.New(viper.GetString("queue-uri"), tasks.ParentQueueName, log.Sub("asynq"))
	// hook handler service
	hookService := hooker.New(ruleService, pm, log.Sub("hooker"), tq)
	// monitor service
	monitorService := monitor.New(pm, hookService, store, log.Sub("monitor"))
	hookService.Run(ctx)

	// run http server
	httpServer := httphandler.New(log.Sub("http"), hookService, monitorService, ruleService)
	httpServer.ListenAndServe(ctx, viper.GetString("listen"))

	<-ctx.Done()
	time.Sleep(1 * time.Second)
}

func showVersion(cmd *cobra.Command, args []string) {
	s := strings.Builder{}
	fmt.Fprintf(&s, "Build version: %s\n", buildVersion)
	fmt.Fprintf(&s, "Build date: %s\n", buildDate)
	fmt.Fprintf(&s, "Build commit: %s\n", buildCommit)
	fmt.Println(s.String())
}
