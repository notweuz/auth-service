package handler

import (
	"auth-service/internal/interceptor"
	"auth-service/internal/interfaces"
	"context"
	"fmt"

	"github.com/notweuz/authentication-proto/pb"
	"github.com/rs/zerolog/log"
)

type userHandler struct {
	pb.UnimplementedUserServiceServer
	service interfaces.UserService
}

func NewUserHandler(userService interfaces.UserService) interfaces.UserHandler {
	return &userHandler{
		service: userService,
	}
}

func (u *userHandler) GetUser(ctx context.Context, request *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	log.Debug().Uint64("user_id", request.UserId).Msg("get user by id")
	user, err := u.service.FindByID(request.UserId)
	if err != nil {
		log.Error().Err(err).Msg("failed to get user by id")
		return nil, err
	}
	log.Debug().Str("username", user.Username).Msg("user found")
	return &pb.GetUserResponse{
		Id:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt.String(),
	}, nil
}

func (u *userHandler) ChangeUsername(ctx context.Context, request *pb.ChangeUsernameRequest) (*pb.GetUserResponse, error) {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("user id not found")
	}
	user, err := u.service.UpdateUsername(userID, request.NewUsername)
	if err != nil {
		return nil, err
	}
	return &pb.GetUserResponse{
		Id:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt.String(),
	}, nil
}
