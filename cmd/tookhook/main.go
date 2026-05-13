package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"log/slog"

	"github.com/k1nky/tookhook/internal/application/service"
	"github.com/k1nky/tookhook/internal/builtin"
	"github.com/k1nky/tookhook/internal/domain/entity"
	"github.com/k1nky/tookhook/internal/infrastructure/queue"
	"github.com/k1nky/tookhook/internal/infrastructure/repository"
	"github.com/k1nky/tookhook/internal/interface/http"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	buildVersion = "dev"
	buildDate    = "unknown"
	buildCommit  = "unknown"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func showVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}

func runServer(cmd *cobra.Command, args []string) {
	// Setup logger
	logLevel := parseLogLevel(viper.GetString("log.level"))
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)

	logger.Info("starting tookhook",
		"version", buildVersion,
		"date", buildDate,
		"commit", buildCommit,
	)

	// Load endpoints from config file
	cfgFile := viper.GetString("config")
	if cfgFile == "" {
		logger.Error("config file not specified")
		os.Exit(1)
	}

	// Create endpoint repository
	endpointRepo := repository.NewFileRepository(cfgFile)
	if err := endpointRepo.Reload(cmd.Context()); err != nil {
		logger.Error("failed to convert config to entities", "error", err)
		os.Exit(1)
	}
	logger.Info("loaded endpoints")

	// Get server and queue settings from viper (can be overridden by env/flags)
	serverListen := viper.GetString("server.listen")
	queueAddr := viper.GetString("queue.addr")
	queueDB := viper.GetInt("queue.db")
	queueConcurrency := viper.GetInt("queue.concurrency")

	logger.Info("configuration",
		"listen", serverListen,
		"redis_addr", queueAddr,
		"redis_db", queueDB,
		"concurrency", queueConcurrency,
	)

	// Create context with cancellation
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Initialize task queue (Asynq)
	taskQueue := queue.NewAsynqQueue(queueAddr, queueDB, queueConcurrency, logger.With("component", "asynq"))
	defer taskQueue.Close()

	logger.Info("task queue initialized")

	// Initialize handler registry and processor
	handlerRegistry := builtin.NewRegistry(logger)
	processor := service.NewProcessor(handlerRegistry, logger.With("component", "processor"))

	logger.Info("handler registry initialized", "handlers", handlerRegistry.Names())

	// Initialize HTTP server
	httpServer := http.NewServer(serverListen, endpointRepo, taskQueue, logger.With("component", "http"))
	if err := httpServer.Start(ctx); err != nil {
		logger.Error("failed to start HTTP server", "error", err)
		os.Exit(1)
	}

	logger.Info("HTTP server started", "addr", serverListen)

	// Start processing tasks
	go func() {
		if err := taskQueue.Start(ctx, func(ctx context.Context, task *entity.WebhookTask) error {
			// Get endpoint configuration
			endpoint, err := endpointRepo.GetByName(ctx, task.EndpointName)
			if err != nil {
				logger.Error("endpoint not found", "endpoint", task.EndpointName)
				return nil // Don't retry
			}
			// Process webhook through pipeline
			return processor.ProcessWebhook(ctx, endpoint, task)
		}); err != nil {
			logger.Error("queue processing error", "error", err)
		}
	}()

	logger.Info("tookhook started")

	<-ctx.Done()
	logger.Info("shutting down tookhook")

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("HTTP server shutdown error", "error", err)
	}
}

func parseLogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
