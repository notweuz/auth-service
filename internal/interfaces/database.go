package interfaces

import "auth-service/internal/model"

type UserDatabase interface {
	Create(user *model.User) (*model.User, error)
	FindByID(id uint64) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
	Update(user *model.User) (*model.User, error)
}
