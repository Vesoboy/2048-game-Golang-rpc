package main

import (
	"game-2048/internal/server"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	// Create router
	router := mux.NewRouter()

	// Create game server
	gameServer := server.NewGameServer(logger)
	gameServer.SetupRoutes(router)

	// Serve static files
	staticDir := filepath.Join(".", "static")
	staticHandler := http.FileServer(http.Dir(staticDir))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticHandler))

	// Serve index.html
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
	})

	// Start server
	addr := ":8081"
	logger.Info("Starting server", zap.String("addr", addr))
	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Fatal("Server failed", zap.Error(err))
	}
}
