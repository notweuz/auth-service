package service

import (
	"auth-service/internal/config"
	"auth-service/internal/errs"
	"auth-service/internal/interfaces"
	"auth-service/internal/model"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/notweuz/authentication-proto/pb"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	cfg         *config.Config
	userService interfaces.UserService
}

func NewAuthService(userService interfaces.UserService, cfg *config.Config) interfaces.AuthService {
	return &authService{userService: userService, cfg: cfg}
}

func (a *authService) Register(request *pb.AuthRequest) (*string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("failed to hash password")
		return nil, errs.InternalError("Failed to register user", err.Error())
	}

	user := &model.User{
		ID:       0,
		Username: request.Username,
		Password: string(hash),
	}

	if err = a.userService.Create(user); err != nil {
		return nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(30 * 24 * time.Hour).Unix(),
	})

	log.Info().Str("username", user.Username).Uint64("id", user.ID).Msg("User registered successfully")

	signedToken, err := token.SignedString([]byte(a.cfg.JwtSecret))

	if err != nil {
		log.Error().Err(err).Msg("failed to sign token")
		return nil, errs.InternalError("Failed to register user", err.Error())
	}

	return &signedToken, nil
}

func (a *authService) Login(request *pb.AuthRequest) (*string, error) {
	user, err := a.userService.FindByUsername(request.Username)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return nil, errs.Unauthorized("Invalid credentials", "Invalid username or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(30 * 24 * time.Hour).Unix(),
	})

	log.Info().Uint64("id", user.ID).Str("username", user.Username).Msg("User successfully logged in")

	signedToken, err := token.SignedString([]byte(a.cfg.JwtSecret))

	if err != nil {
		log.Error().Err(err).Msg("failed to sign token")
		return nil, errs.InternalError("Failed to login", err.Error())
	}

	return &signedToken, nil
}

func (a *authService) ChangePassword(userID uint64, request *pb.ChangePasswordRequest) (*string, error) {
	user, err := a.userService.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.OldPassword)); err != nil {
		return nil, errs.Unauthorized("Invalid credentials", "Invalid username or password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("failed to hash password")
		return nil, errs.InternalError("Failed to change password", err.Error())
	}

	user.Password = string(hashedPassword)
	if _, err = a.userService.Update(user); err != nil {
		log.Error().Err(err).Msg("failed to update user")
		return nil, errs.InternalError("Failed to change password", err.Error())
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(30 * 24 * time.Hour).Unix(),
	})

	log.Info().Uint64("id", user.ID).Str("username", user.Username).Msg("User password successfully changed")

	signedToken, err := token.SignedString([]byte(a.cfg.JwtSecret))
	if err != nil {
		log.Error().Err(err).Msg("failed to sign token")
		return nil, errs.InternalError("Failed to change password", err.Error())
	}

	return &signedToken, nil
}
