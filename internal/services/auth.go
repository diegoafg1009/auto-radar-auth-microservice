package services

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/diegoafg1009/auto-radar-auth-microservice/internal/domain"
	"github.com/diegoafg1009/auto-radar-auth-microservice/internal/dtos"
	"github.com/diegoafg1009/auto-radar-auth-microservice/internal/repositories"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/crypto/bcrypt"
)

type Auth interface {
	Login(ctx context.Context, req dtos.LoginUserRequest) (res dtos.LoginUserResponse, err error)
	Register(ctx context.Context, req dtos.RegisterUserRequest) (res dtos.RegisterUserResponse, err error)
}

type auth struct {
	userRepository repositories.User
}

func NewAuthService(userRepository repositories.User) Auth {
	return &auth{userRepository: userRepository}
}

func (a *auth) Login(ctx context.Context, request dtos.LoginUserRequest) (res dtos.LoginUserResponse, err error) {
	user, err := a.userRepository.GetByEmail(ctx, request.Email)
	if err != nil {
		return dtos.LoginUserResponse{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return dtos.LoginUserResponse{}, err
	}

	token, err := generateToken(user)
	if err != nil {
		log.Println("Error generating token:", err)
		return dtos.LoginUserResponse{}, err
	}

	return dtos.LoginUserResponse{
		Token: token,
	}, nil
}

func (a *auth) Register(ctx context.Context, request dtos.RegisterUserRequest) (res dtos.RegisterUserResponse, err error) {
	user, err := a.userRepository.GetByEmail(ctx, request.Email)
	if err != nil {
		return dtos.RegisterUserResponse{}, err
	}

	if user != nil {
		return dtos.RegisterUserResponse{}, errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return dtos.RegisterUserResponse{}, err
	}

	newUser := &domain.User{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		Password:  string(hashedPassword),
	}

	id, err := a.userRepository.Create(ctx, newUser)
	if err != nil {
		return dtos.RegisterUserResponse{}, err
	}

	return dtos.RegisterUserResponse{
		Id: id,
	}, nil
}

func generateToken(user *domain.User) (string, error) {
	var jwtSecret = os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", errors.New("JWT_SECRET is not set")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": user.ID.Hex(),
			"exp": time.Now().Add(time.Hour * 24).Unix(),
			"iat": time.Now().Unix(),
		})

	stringifiedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return stringifiedToken, nil
}
