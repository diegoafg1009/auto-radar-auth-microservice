package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/diegoafg1009/auto-radar-auth-microservice/internal/repositories"
	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
	Health() map[string]string
	GetUserRepository() repositories.User
}

type service struct {
	db   *mongo.Client
	user repositories.User
}

var (
	host     = os.Getenv("MONGO_DB_HOST")
	port     = os.Getenv("MONGO_DB_PORT")
	username = os.Getenv("MONGO_DB_USERNAME")
	password = os.Getenv("MONGO_DB_PASSWORD")
	database = os.Getenv("MONGO_DB_DATABASE")
)

func New() Service {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)))

	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatal("Database connection failed")
	}

	log.Println("Database connected successfully")

	db := client.Database(database)

	return &service{
		db:   db.Client(),
		user: NewUserRepository(db),
	}
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := s.db.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("db down: %v", err)
	}

	return map[string]string{
		"message": "It's healthy",
	}
}

func (s *service) GetUserRepository() repositories.User {
	return s.user
}
