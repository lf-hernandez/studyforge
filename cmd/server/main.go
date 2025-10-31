package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"studyforge/internal/api/handlers"
	"studyforge/internal/config"
	"studyforge/internal/repository"
	"studyforge/internal/services"
	"studyforge/pkg/ai"
	"studyforge/pkg/utils"
)

func main() {
	log.Println("Starting StudyForge server...")

	// Load configuration
	cfg := config.Load()
	log.Printf("Configuration loaded: server=%s:%s, database=%s", cfg.ServerHost, cfg.ServerPort, cfg.DatabasePath)

	// Initialize database
	db, err := repository.NewDatabase(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := db.RunMigrations("./migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	sessionRepo := repository.NewSessionRepository(db.DB)
	docRepo := repository.NewDocumentRepository(db.DB)
	contentRepo := repository.NewContentRepository(db.DB)

	// Initialize services
	pdfService := services.NewPDFService(contentRepo)
	aiClient := ai.NewHuggingFaceClient(cfg.HuggingFaceKey, cfg.HuggingFaceURL)
	studyService := services.NewStudyService(aiClient, pdfService, contentRepo, docRepo)

	// Initialize handlers
	pdfHandler := handlers.NewPDFHandler(cfg, docRepo, pdfService)
	studyHandler := handlers.NewStudyHandler(studyService)

	// Initialize session manager
	sessionManager := utils.NewSessionManager(sessionRepo)

	// Setup router
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/api/health", handlers.HandleHealth)
	mux.HandleFunc("/api/documents/upload", pdfHandler.HandleUpload)
	mux.HandleFunc("/api/documents", pdfHandler.HandleGetDocument)
	mux.HandleFunc("/api/study/generate", studyHandler.HandleGenerate)
	mux.HandleFunc("/api/study/content", studyHandler.HandleGetContent)

	// Serve static files
	fs := http.FileServer(http.Dir("./web"))
	mux.Handle("/", fs)

	// Wrap with session middleware
	handler := sessionManager.Middleware(corsMiddleware(mux))

	// Start server
	addr := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)
	log.Printf("Server listening on http://%s", addr)

	// Setup graceful shutdown
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	// Start server in goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}

// corsMiddleware adds CORS headers
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
