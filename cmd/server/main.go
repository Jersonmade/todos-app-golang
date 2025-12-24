package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Jersonmade/todos-app-golang/internal/handlers"
	"github.com/Jersonmade/todos-app-golang/internal/storage"
)

func main() {
	taskStorage := storage.NewMemoryStorage()
	taskHandler := handlers.NewTaskHandler(taskStorage)

	mux := http.NewServeMux()

	mux.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			taskHandler.Create(w, r)
		case http.MethodGet:
			taskHandler.GetAll(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/todos/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			taskHandler.Get(w, r)
		case http.MethodPut:
			taskHandler.Update(w, r)
		case http.MethodDelete:
			taskHandler.Delete(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	signalCh := make(chan os.Signal, 1)

	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("Server up...")
		err := server.ListenAndServe()

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Println("Error up server", err)
		}
	}()

	<-signalCh
	log.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)

	if err != nil {
		log.Println("Server is closed forcibly: ", err)
	}

	log.Println("Server gracefully stopped")
}
