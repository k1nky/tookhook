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
	"github.com/k1nky/tookhook/internal/logger"
	"github.com/k1nky/tookhook/internal/service/hooker"
)

func main() {
	log := logger.New()
	log.SetLevel("debug")

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	cfg := config.Config{}
	if err := config.Parse(&cfg); err != nil {
		panic(err)
	}
	log.SetLevel(cfg.LogLevel)
	log.Debugf("config: %+v", cfg)

	run(ctx, cfg, log)

	<-ctx.Done()
	time.Sleep(1 * time.Second)
}

func run(ctx context.Context, cfg config.Config, log *logger.Logger) {
	pm := pluginmanager.New()
	for _, v := range strings.Split(cfg.Plugins, ",") {
		_, name := path.Split(v)
		pm.LoadForwarder(ctx, name, v)
	}

	store := database.New(cfg.DarabaseURI, log)
	if err := store.Open(ctx); err != nil {
		log.Errorf("failed opening db: %v", err)
		return
	}
	hookService := hooker.New(store, pm, log)

	httpServer := httphandler.New(log, hookService)
	httpServer.ListenAndServe(ctx, string(cfg.Listen))
}
