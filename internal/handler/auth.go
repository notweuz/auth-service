package handler

import (
	"auth-service/internal/errs"
	"auth-service/internal/interceptor"
	"auth-service/internal/interfaces"
	"context"

	"github.com/notweuz/auth-proto/pb"
	"github.com/rs/zerolog/log"
)

type authHandler struct {
	pb.UnimplementedAuthServiceServer
	service interfaces.AuthService
}

func NewAuthHandler(service interfaces.AuthService) interfaces.AuthHandler {
	return &authHandler{service: service}
}

func (a authHandler) Login(ctx context.Context, request *pb.AuthRequest) (*pb.AuthResponse, error) {
	log.Debug().Msg("Login request received")
	authResponse, err := a.service.Login(request)
	if err != nil {
		return nil, errs.ToGRPC(err)
	}
	return authResponse, nil
}

func (a authHandler) Register(ctx context.Context, request *pb.AuthRequest) (*pb.AuthResponse, error) {
	log.Debug().Msg("Register request received")
	authResponse, err := a.service.Register(request)
	if err != nil {
		return nil, errs.ToGRPC(err)
	}
	return authResponse, nil
}

func (a authHandler) ChangePassword(ctx context.Context, request *pb.ChangePasswordRequest) (*pb.AuthResponse, error) {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, errs.ToGRPC(errs.Unauthorized("Invalid credentials", "No token provided"))
	}
	authResponse, err := a.service.ChangePassword(userID, request)
	if err != nil {
		return nil, errs.ToGRPC(err)
	}
	return authResponse, nil
}
