package database

import (
	"auth-service/internal/errs"
	"auth-service/internal/interfaces"
	"auth-service/internal/model"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type userDatabase struct {
	db *gorm.DB
}

func NewUserDatabase(db *gorm.DB) interfaces.UserDatabase {
	return &userDatabase{db: db}
}

func isDuplicateKeyError(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

func (u *userDatabase) Create(user *model.User) (*model.User, error) {
	if err := u.db.Create(user).Error; err != nil {
		if isDuplicateKeyError(err) {
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
		if isDuplicateKeyError(err) {
			return nil, errs.RecordAlreadyExists
		}
		return nil, err
	}
	return user, nil
}

func (u *userDatabase) FindByUsername(username string) (*model.User, error) {
	var user model.User
	if err := u.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.RecordNotFound
		}
		return nil, err
	}
	return &user, nil
}
