package database

import (
	"context"

	"time"

	"github.com/diegoafg1009/auto-radar-auth-microservice/internal/domain"
	"github.com/diegoafg1009/auto-radar-auth-microservice/internal/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type user struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) repositories.User {
	return &user{collection: db.Collection("users")}
}

func (u *user) Create(ctx context.Context, user *domain.User) (id string, err error) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	result, err := u.collection.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}

	return result.InsertedID.(string), nil
}

func (u *user) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	err := u.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	return &user, err
}

func (u *user) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := u.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	return &user, err
}
