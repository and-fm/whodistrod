package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/and-fm/whodistrod/internal/config"
	"github.com/and-fm/whodistrod/internal/core"
	"go.uber.org/dig"
)

type server struct {
	dig.In

	BaseRouter core.BaseRouter
	Config     *config.Config
}

func NewServer(s server) core.Server {
	return &s
}

func (s *server) Start() {
	done := make(chan bool, 1)
	go s.GracefulShutdown(done)

	s.BaseRouter.BaseEcho().HideBanner = true
	s.BaseRouter.BaseEcho().Start(fmt.Sprintf(":%v", s.Config.Port))

	<-done
	log.Println("Graceful shutdown complete.")
}

func (s *server) GracefulShutdown(done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop() // Allow Ctrl+C to force shutdown

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.BaseRouter.BaseEcho().Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}
