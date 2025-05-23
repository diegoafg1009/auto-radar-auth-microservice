package handlers

import (
	"context"

	"connectrpc.com/connect"
	"github.com/diegoafg1009/auto-radar-auth-microservice/internal/dtos"
	"github.com/diegoafg1009/auto-radar-auth-microservice/internal/services"
	v1 "github.com/diegoafg1009/auto-radar-auth-microservice/pkg/genproto/auth/v1"
	"github.com/diegoafg1009/auto-radar-auth-microservice/pkg/genproto/auth/v1/authv1connect"
)

type Auth struct {
	authv1connect.UnimplementedAuthServiceHandler
	authService services.Auth
}

func NewAuthHandler(authService services.Auth) *Auth {
	return &Auth{
		authService: authService,
	}
}

func (h *Auth) Login(ctx context.Context, req *connect.Request[v1.LoginRequest]) (res *connect.Response[v1.LoginResponse], err error) {
	request := dtos.LoginUserRequest{
		Email:    req.Msg.Email,
		Password: req.Msg.Password,
	}
	loginResponse, err := h.authService.Login(ctx, request)
	if err != nil {
		return nil, err
	}
	response := v1.LoginResponse{
		Token: loginResponse.Token,
	}
	return connect.NewResponse(&response), nil
}

func (h *Auth) Register(ctx context.Context, req *connect.Request[v1.RegisterRequest]) (res *connect.Response[v1.RegisterResponse], err error) {
	request := dtos.RegisterUserRequest{
		FirstName: req.Msg.FirstName,
		LastName:  req.Msg.LastName,
		Email:     req.Msg.Email,
		Password:  req.Msg.Password,
	}
	registerResponse, err := h.authService.Register(ctx, request)
	if err != nil {
		return nil, err
	}
	response := v1.RegisterResponse{
		Id: registerResponse.Id,
	}
	return connect.NewResponse(&response), nil
}
