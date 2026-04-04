package database

import (
	"auth-service/internal/errs"
	"auth-service/internal/interfaces"
	"auth-service/internal/model"
	"errors"

	"gorm.io/gorm"
)

type userDatabase struct {
	db *gorm.DB
}

func NewUserDatabase(db *gorm.DB) interfaces.UserDatabase {
	return &userDatabase{db: db}
}

func (u *userDatabase) Create(user *model.User) (*model.User, error) {
	if err := u.db.Create(user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, errs.RecordAlreadyExists
		}
		return nil, err
	}
	return user, nil
}

func (u *userDatabase) FindByID(id uint64) (*model.User, error) {
	var user model.User
	if err := u.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.RecordNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (u *userDatabase) Update(user *model.User) (*model.User, error) {
	if err := u.db.Save(user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, errs.RecordAlreadyExists
		}
		return nil, err
	}
	return user, nil
}

func (u *userDatabase) FindByUsername(username string) (*model.User, error) {
	var user model.User
	if err := u.db.First(&user.Username, username).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.RecordNotFound
		}
		return nil, err
	}
	return &user, nil
}
