package interfaces

import (
	"auth-service/internal/model"

	"github.com/notweuz/authentication-proto/pb"
)

type AuthService interface {
	Register(request *pb.AuthRequest) (*string, error)
	Login(request *pb.AuthRequest) (*string, error)
	ChangePassword(id uint64, request *pb.ChangePasswordRequest) (*string, error)
}

type UserService interface {
	Create(user *model.User) error
	FindByID(id uint64) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
	Update(user *model.User) (*model.User, error)
	UpdateUsername(id uint64, newUsername string) (*model.User, error)
}
