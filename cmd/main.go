package main

import (
	"context"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/k1nky/tookhook/internal/adapter/database"
	httphandler "github.com/k1nky/tookhook/internal/adapter/http"
	"github.com/k1nky/tookhook/internal/adapter/pluginmanager"
	"github.com/k1nky/tookhook/internal/config"
	"github.com/k1nky/tookhook/internal/service/hooker"
	"github.com/k1nky/tookhook/internal/service/monitor"
	"github.com/k1nky/tookhook/internal/service/ruler"
	"github.com/k1nky/tookhook/pkg/logger"
)

const (
	LoggerName         = "tookhook"
	LoggerDefaultLevel = "debug"
)

func main() {
	log := logger.New(LoggerName)
	log.SetLevel(LoggerDefaultLevel)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	// load the service config
	cfg := config.Config{}
	if err := config.Parse(&cfg); err != nil {
		log.Errorf("config: %s", err)
		return
	}
	log.SetLevel(cfg.LogLevel)
	log.Debugf("config: %+v", cfg)

	// run the service
	run(ctx, cfg, log)

	<-ctx.Done()
	time.Sleep(1 * time.Second)
}

func run(ctx context.Context, cfg config.Config, log *logger.Logger) {
	// load plugins
	pm := pluginmanager.New(log)
	for _, v := range strings.Split(cfg.Plugins, ",") {
		_, name := path.Split(v)
		if err := pm.Load(ctx, name, v); err != nil {
			log.Errorf("plugins: %s", err)
			return
		}
	}
	pm.Run(ctx)

	// open rules store
	store := database.New(cfg.DarabaseURI, log)
	if err := store.Open(ctx); err != nil {
		log.Errorf("failed opening db: %s", err)
		return
	}
	ruleService := ruler.New(pm, store, log)
	if err := ruleService.Load(ctx); err != nil {
		log.Errorf("failed loading rules: %s", err)
		return
	}
	// hook handler service
	hookService := hooker.New(ruleService, pm, log)
	// monitor service
	monitorService := monitor.New(pm, hookService, store, log)

	// run http server
	httpServer := httphandler.New(log, hookService, monitorService, ruleService)
	httpServer.ListenAndServe(ctx, string(cfg.Listen))
}
