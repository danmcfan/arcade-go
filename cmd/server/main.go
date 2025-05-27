package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danmcfan/arcade-go/internal"
)

func main() {
	server := &http.Server{
		Addr:    ":8080",
		Handler: http.FileServer(http.FS(internal.AssetFiles())),
	}

	// Channel to listen for errors coming from the listener.
	serverErrors := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		fmt.Println("Starting server on port 8080...")
		serverErrors <- server.ListenAndServe()
	}()

	// Channel to listen for an interrupt or terminate signal from the OS.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		fmt.Printf("Error starting server: %v\n", err)
		server.Close()

	case sig := <-shutdown:
		fmt.Printf("Server is shutting down due to %v signal\n", sig)

		// Give outstanding requests 5 seconds to complete.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			fmt.Printf("Could not stop server gracefully: %v\n", err)
			server.Close()
		}
	}
}
