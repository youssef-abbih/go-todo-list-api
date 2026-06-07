// @title Todo List API
// @version 1.0
// @description This is a simple REST API for managing tasks.
// @host localhost:8080
// @BasePath /
// @schemes http

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/swaggo/http-swagger"
	// docs is imported side-effect exclusively to register swagger documentation specifications
	_ "github.com/youssef-abbih/go-todo-list/docs"
	"github.com/youssef-abbih/go-todo-list/handlers"
	"github.com/youssef-abbih/go-todo-list/models"
	"github.com/youssef-abbih/go-todo-list/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	// Initialize DB
	models.InitDB()

	// Set up router
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.SecureHeadersMiddleware)   // Security headers
	r.Use(middleware.LogRequestMiddleware)      // Request logging

	// Public routes
	r.Get("/", handlers.DefaultResponse)
	r.Get("/health", handlers.HealthCheck)
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// Protected /tasks routes
	r.Route("/tasks", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware) // Scoped only to /tasks/*
		r.Get("/", handlers.GetTasks)
		r.Post("/", handlers.PostTask)
		r.Get("/{id}", handlers.GetTask)
		r.Put("/{id}", handlers.PutTask)
		r.Delete("/{id}", handlers.DeleteTask)
	})

	// Server setup
	port := ":8080"
	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}

	// Graceful shutdown setup
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Println("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
		close(idleConnsClosed)
	}()

	// Start server
	log.Printf("Server running on http://localhost%s", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}

	<-idleConnsClosed
	log.Println("Server stopped gracefully")
}
