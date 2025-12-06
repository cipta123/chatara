package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/umess/backend/internal/handler"
	"github.com/umess/backend/internal/middleware"
	"github.com/umess/backend/internal/websocket"
	"github.com/umess/backend/pkg/auth"
	"github.com/umess/backend/pkg/config"
	"github.com/umess/backend/pkg/database"
)

func main() {
	// Load configuration
	if err := config.Load(); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize JWT
	auth.Init()

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	// Create WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()

	// Initialize handlers
	authHandler := handler.NewAuthHandler()
	conversationHandler := handler.NewConversationHandler()
	messageHandler := handler.NewMessageHandler(hub)
	mediaHandler := handler.NewMediaHandler()
	wsHandler := handler.NewWebSocketHandler(hub)

	// Setup router
	router := mux.NewRouter()

	// CORS middleware
	router.Use(middleware.CORSMiddleware)

	// Public routes
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/auth/register", authHandler.Register).Methods("POST", "OPTIONS")
	api.HandleFunc("/auth/login", authHandler.Login).Methods("POST", "OPTIONS")

	// Protected routes
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	protected.HandleFunc("/auth/me", authHandler.GetMe).Methods("GET", "OPTIONS")
	protected.HandleFunc("/users/search", authHandler.SearchUsers).Methods("GET", "OPTIONS")
	protected.HandleFunc("/conversations", conversationHandler.GetConversations).Methods("GET", "OPTIONS")
	protected.HandleFunc("/conversations", conversationHandler.CreateDirectConversation).Methods("POST", "OPTIONS")
	protected.HandleFunc("/conversations/{conversationId}", conversationHandler.GetConversation).Methods("GET", "OPTIONS")
	protected.HandleFunc("/conversations/{conversationId}/messages", messageHandler.GetMessages).Methods("GET", "OPTIONS")
	protected.HandleFunc("/conversations/{conversationId}/messages", messageHandler.SendMessage).Methods("POST", "OPTIONS")
	protected.HandleFunc("/conversations/{conversationId}/read", messageHandler.MarkAsRead).Methods("POST", "OPTIONS")

	// Media routes
	protected.HandleFunc("/media/upload", mediaHandler.UploadImage).Methods("POST", "OPTIONS")
	protected.HandleFunc("/media/{filename}", mediaHandler.ServeMedia).Methods("GET", "OPTIONS")
	protected.HandleFunc("/media/send", mediaHandler.SendMessageWithImage).Methods("POST", "OPTIONS")

	// WebSocket route
	router.HandleFunc("/ws", wsHandler.HandleConnection)

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// Start server
	port := config.AppConfig.Port
	log.Printf("Server starting on port %s", port)
	log.Printf("WebSocket endpoint: ws://localhost:%s/ws", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
