package main

import (
	//"fmt"
	//"log"
	"time"
	"log/slog"
	"net/http"
	"os/signal"
	"context"
	// "net/http"
	"os"
	"syscall"

	"github.com/asthap08/Golang_notes/internal/config"
)

func main() {
	//load config
	cfg := config.MustLoad()

	//setup router
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request){
		w.Write([]byte("Hello, World!"))
	})

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