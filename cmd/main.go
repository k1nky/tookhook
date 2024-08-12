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

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
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
	// set log level from config
	log.SetLevel(cfg.LogLevel)
	log.Debugf("config: %+v", cfg)
	// version info was requested
	if cfg.Version {
		showVersion()
		return
	}

	// run the service
	if err := run(ctx, cfg, log); err != nil {
		log.Errorf("%s", err)
		os.Exit(1)
		return
	}

	<-ctx.Done()
	time.Sleep(1 * time.Second)
}

func run(ctx context.Context, cfg config.Config, log *logger.Logger) error {
	// load plugins
	pm := pluginmanager.New(log)
	if len(cfg.Plugins) > 0 {
		for _, v := range strings.Split(cfg.Plugins, ",") {
			_, name := path.Split(v)
			if err := pm.Load(ctx, name, v); err != nil {
				return fmt.Errorf("plugins: %s", err)
			}
		}
	}
	pm.Run(ctx)

	// open rules store
	store := database.New(cfg.DarabaseURI, log)
	if err := store.Open(ctx); err != nil {
		return fmt.Errorf("failed opening db: %s", err)
	}
	ruleService := ruler.New(pm, store, log)
	if err := ruleService.Load(ctx); err != nil {
		return fmt.Errorf("failed loading rules: %s", err)
	}
	// hook handler service
	hookService := hooker.New(ruleService, pm, log)
	// monitor service
	monitorService := monitor.New(pm, hookService, store, log)

	// run http server
	httpServer := httphandler.New(log, hookService, monitorService, ruleService)
	httpServer.ListenAndServe(ctx, string(cfg.Listen))
	return nil
}

func showVersion() {
	s := strings.Builder{}
	fmt.Fprintf(&s, "Build version: %s\n", buildVersion)
	fmt.Fprintf(&s, "Build date: %s\n", buildDate)
	fmt.Fprintf(&s, "Build commit: %s\n", buildCommit)
	fmt.Println(s.String())
}
