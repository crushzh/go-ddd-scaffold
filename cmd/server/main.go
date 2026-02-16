package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-ddd-scaffold/internal/container"
	"go-ddd-scaffold/internal/interfaces/http/router"
	"go-ddd-scaffold/pkg/config"
	"go-ddd-scaffold/pkg/logger"
)

// Build-time variables
var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func main() {
	// Command-line flags
	configFile := flag.String("c", "configs/config.yaml", "config file path")
	showVersion := flag.Bool("v", false, "show version")
	flag.Parse()

	if *showVersion {
		fmt.Printf("Version: %s\nBuild:   %s\nCommit:  %s\n", Version, BuildTime, GitCommit)
		os.Exit(0)
	}

	// 1. Load config
	cfg, err := config.Load(*configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 2. Init logger
	if err := logger.Init(&cfg.Log); err != nil {
		fmt.Fprintf(os.Stderr, "failed to init logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()
	logger.Infof("starting service version=%s config=%s", Version, *configFile)

	// 3. Init container (dependency injection)
	c, err := container.New(cfg)
	if err != nil {
		logger.Fatalf("failed to init container: %v", err)
	}
	defer c.Close()

	// 4. Setup router
	r := router.Setup(c)

	// 5. Start HTTP server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	go func() {
		logger.Infof("HTTP server started: %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("HTTP server error: %v", err)
		}
	}()

	// 6. Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	logger.Infof("received signal: %s, shutting down...", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("HTTP server shutdown error: %v", err)
	}

	logger.Info("service exited")
}
