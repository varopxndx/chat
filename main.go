package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/varopxndx/chat/config"
	"github.com/varopxndx/chat/handler"
	"github.com/varopxndx/chat/rabbit"
	"github.com/varopxndx/chat/router"
	"github.com/varopxndx/chat/service"
	"github.com/varopxndx/chat/token"
	"github.com/varopxndx/chat/usecase"
)

// Abnormal exit constants
const (
	ExitAbnormalErrorLoadingConfiguration = iota
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// read config file
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config: %w", err)
		os.Exit(ExitAbnormalErrorLoadingConfiguration)
	}

	// create services
	rabbit := rabbit.New(cfg.RabbitMQ)
	db := service.SetupDB(cfg.DB)
	defer rabbit.Close()
	defer db.Close()

	tokenService := token.New(cfg.JwtSecretKey)

	// create usecase layer
	useCases := usecase.New(db, rabbit, cfg.StooqURL)

	// create webserver
	ws := handler.NewWebsocketServer(useCases, rabbit)
	go ws.Run()

	// create handler layer
	handler := handler.New(useCases, ws, tokenService)

	// create router layer
	router := router.Setup(handler, tokenService)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Sprintf("Fatal error starting server: %v \n", err))
		}
	}()

	fmt.Println("Server ready: http://localhost:8080/v1/")

	<-ctx.Done()
	stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		panic(fmt.Sprintf("Fatal error shutdown server: %v \n", err))
	}
	log.Println("Server has been stopped...")
}
