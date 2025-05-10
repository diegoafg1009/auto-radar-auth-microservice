package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/diegoafg1009/auto-radar-auth-microservice/internal/database"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port      int
	apiServer *http.Server
	db        database.Service
}

func NewServer() *Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	newServer := &Server{
		port: port,
		db:   database.New(),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", newServer.port),
		Handler:      newServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	newServer.apiServer = server

	return newServer
}

func (s *Server) ListenAndServe() error {
	return s.apiServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server on port", s.port)
	return s.apiServer.Shutdown(context.Background())
}
