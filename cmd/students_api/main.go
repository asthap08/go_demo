package main

import (
	//"fmt"
	"context"
	"log"
	"log/slog"
	"net/http"
	"os/signal"
	"time"
	"os"
	"syscall"
	

	"github.com/asthap08/Golang_notes/internal/config"
	"github.com/asthap08/Golang_notes/internal/http/handlers/student"
	"github.com/asthap08/Golang_notes/internal/storage/sqlite"
)

func main() {
	//load config
	cfg := config.MustLoad()

	//database setup
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("storage initialized",slog.String("env", cfg.Env), slog.String("version", "1.8.9"))

	//setup router
	router := http.NewServeMux()


	router.HandleFunc("POST /api/students", student.New(storage))

	//setup server

	server := http.Server{
		Addr:    cfg.HTTPServer.Address,
		Handler: router,
	}

	slog.Info("Starting server", slog.String("address", cfg.HTTPServer.Address))
	//fmt.Printf("Server started %s", cfg.HTTPServer.Address)

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			slog.Error("failed to start server", slog.String("error", err.Error()))
		}
	}()

<-done

slog.Info("Server stopped gracefully")

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

if err := server.Shutdown(ctx); err != nil {
	slog.Error("Server shutdown failed", slog.String("error", err.Error()))
}

slog.Info("Server exited properly")

}