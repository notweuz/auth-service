package handler

import (
	"auth-service/internal/errs"
	"auth-service/internal/interceptor"
	"auth-service/internal/interfaces"
	"context"

	"github.com/notweuz/authentication-proto/pb"
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
	token, err := a.service.Login(request)
	if err != nil {
		return nil, err
	}
	return &pb.AuthResponse{Token: *token}, nil
}

func (a authHandler) Register(ctx context.Context, request *pb.AuthRequest) (*pb.AuthResponse, error) {
	log.Debug().Msg("Register request received")
	token, err := a.service.Register(request)
	if err != nil {
		return nil, err
	}
	return &pb.AuthResponse{Token: *token}, nil
}

func (a authHandler) ChangePassword(ctx context.Context, request *pb.ChangePasswordRequest) (*pb.AuthResponse, error) {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, errs.Unauthorized("Invalid credentials", "No token provided")
	}
	token, err := a.service.ChangePassword(userID, request)
	if err != nil {
		return nil, err
	}
	return &pb.AuthResponse{Token: *token}, nil
}
