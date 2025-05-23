package server

import (
	"net/http"

	"github.com/diegoafg1009/auto-radar-auth-microservice/internal/handlers"
	"github.com/diegoafg1009/auto-radar-auth-microservice/internal/services"
	"github.com/diegoafg1009/auto-radar-auth-microservice/pkg/genproto/auth/v1/authv1connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	path, handler := authv1connect.NewAuthServiceHandler(
		handlers.NewAuthHandler(
			services.NewAuthService(
				s.db.GetUserRepository(),
			),
		),
	)

	mux.Handle(path, handler)

	return h2c.NewHandler(mux, &http2.Server{})
}
