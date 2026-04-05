package interfaces

import (
	"auth-service/internal/model"

	"github.com/notweuz/auth-proto/pb"
)

type AuthService interface {
	Register(request *pb.AuthRequest) (*pb.AuthResponse, error)
	Login(request *pb.AuthRequest) (*pb.AuthResponse, error)
	ChangePassword(id uint64, request *pb.ChangePasswordRequest) (*pb.AuthResponse, error)
	ValidateToken(token string) (bool, error)
}

type UserService interface {
	Create(user *model.User) error
	FindByID(id uint64) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
	Update(user *model.User) (*model.User, error)
	UpdateUsername(id uint64, newUsername string) (*model.User, error)
}
