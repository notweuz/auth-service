package service

import (
	"auth-service/internal/errs"
	"auth-service/internal/interfaces"
	"auth-service/internal/model"
	"errors"

	"github.com/rs/zerolog/log"
)

type userService struct {
	database interfaces.UserDatabase
}

func NewUserService(database interfaces.UserDatabase) interfaces.UserService {
	return &userService{database: database}
}

func (u *userService) Create(user *model.User) error {
	createdUser, err := u.database.Create(user)
	log.Debug().Str("username", user.Username).Msg("Attempting to create a new user")
	if err != nil {
		if errors.Is(err, errs.RecordAlreadyExists) {
			log.Warn().Str("username", user.Username).Msg("User already exists")
			return errs.Conflict("Couldn't create a new user", "username is taken")
		}
		log.Error().Err(err).Msg("Failed to create a new user")
		return errs.InternalError("Couldn't create a new user", err.Error())
	}
	log.Info().Str("username", createdUser.Username).Uint64("id", createdUser.ID).Msg("User created successfully")
	return nil
}

func (u *userService) FindByID(id uint64) (*model.User, error) {
	user, err := u.database.FindByID(id)
	log.Debug().Uint64("id", id).Msg("Attempting to find a user by id")
	if err != nil {
		if errors.Is(err, errs.RecordNotFound) {
			log.Warn().Uint64("id", id).Msg("User not found")
			return nil, errs.NotFound("User not found", "User with that id doesn't exist")
		}
		log.Error().Err(err).Msg("Failed to find a user by id")
		return nil, errs.InternalError("Couldn't find a user", err.Error())
	}
	log.Info().Str("username", user.Username).Uint64("id", user.ID).Msg("User found successfully")
	return user, nil
}

func (u *userService) Update(user *model.User) (*model.User, error) {
	updatedUser, err := u.database.Update(user)
	log.Debug().Uint64("id", user.ID).Msg("Attempting to update user")
	if err != nil {
		if errors.Is(err, errs.RecordAlreadyExists) {
			log.Warn().Uint64("id", user.ID).Msg("User already exists")
			return nil, errs.Conflict("Couldn't update user", "username is taken")
		}
		log.Error().Err(err).Msg("Failed to update user")
		return nil, err
	}
	log.Info().Uint64("id", updatedUser.ID).Msg("User updated successfully")
	return updatedUser, nil
}

func (u *userService) UpdateUsername(id uint64, newUsername string) (*model.User, error) {
	user, err := u.database.FindByID(id)
	if err != nil {
		return nil, err
	}
	log.Info().Str("old_username", user.Username).Str("new_username", newUsername).Uint64("id", id).
		Msg("Attempting to update user")
	user.Username = newUsername
	return u.Update(user)
}

func (u *userService) FindByUsername(username string) (*model.User, error) {
	user, err := u.database.FindByUsername(username)
	if err != nil {
		if errors.Is(err, errs.RecordNotFound) {
			log.Warn().Str("username", username).Msg("User not found")
			return nil, errs.NotFound("User not found", "User with that username doesn't exist")
		}
		log.Error().Err(err).Str("username", username).Msg("Failed to find a user by username")
		return nil, errs.InternalError("Couldn't find a user", err.Error())
	}
	log.Info().Str("username", user.Username).Uint64("id", user.ID).Msg("User found successfully")
	return user, err
}
