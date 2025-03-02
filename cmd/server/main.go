package main

import (
	"Server-Monitoring-System/internal/config"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	ctx := context.Background()

	cfg, err := config.NewConfigFromEnv(ctx)
	if err != nil {
		fmt.Println("failed to load config: %w", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, World!")
	})

	server := &http.Server{
		Addr:    cfg.ServerPort,
		Handler: mux,
	}

	go func() {
		err = server.ListenAndServe()
		if err != nil {
			fmt.Println("failed to start server: %w", err)
		}
	}()

	fmt.Println("server started")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	err = server.Shutdown(nil)

}
