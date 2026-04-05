package service

import (
	"auth-service/internal/config"
	"auth-service/internal/errs"
	"auth-service/internal/interfaces"
	"auth-service/internal/model"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/notweuz/auth-proto/pb"
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

func (a *authService) Register(request *pb.AuthRequest) (*pb.AuthResponse, error) {
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

	token, err := a.generateToken(user)
	if err != nil {
		return nil, err
	}

	log.Info().Str("username", user.Username).Uint64("id", user.ID).Msg("User registered successfully")

	return &pb.AuthResponse{
		Token:  *token,
		UserId: user.ID,
	}, nil
}

func (a *authService) Login(request *pb.AuthRequest) (*pb.AuthResponse, error) {
	user, err := a.userService.FindByUsername(request.Username)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return nil, errs.Unauthorized("Invalid credentials", "Invalid username or password")
	}

	token, err := a.generateToken(user)
	if err != nil {
		return nil, err
	}

	log.Info().Uint64("id", user.ID).Str("username", user.Username).Msg("User successfully logged in")

	return &pb.AuthResponse{
		Token:  *token,
		UserId: user.ID,
	}, nil
}

func (a *authService) ChangePassword(userID uint64, request *pb.ChangePasswordRequest) (*pb.AuthResponse, error) {
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
	user.PasswordVersion++
	if _, err = a.userService.Update(user); err != nil {
		log.Error().Err(err).Msg("failed to update user")
		return nil, errs.InternalError("Failed to change password", err.Error())
	}

	token, err := a.generateToken(user)
	if err != nil {
		return nil, err
	}

	log.Info().Uint64("id", user.ID).Str("username", user.Username).Msg("User password successfully changed")

	return &pb.AuthResponse{
		Token:  *token,
		UserId: userID,
	}, nil
}

func (a *authService) ValidateToken(tokenString string) (bool, error) {
	var claims model.Claims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(a.cfg.JwtSecret), nil
	})
	if err != nil {
		return false, errs.Unauthorized("Invalid credentials", "Invalid token")
	}

	if !token.Valid {
		return false, nil
	}

	userID := claims.Subject
	user, err := a.userService.FindByID(userID)
	if err != nil {
		return false, err
	}

	if claims.PasswordVersion != user.PasswordVersion {
		return false, nil
	}

	return true, nil
}

func (a *authService) generateToken(user *model.User) (*string, error) {
	claims := &model.Claims{
		Subject:         user.ID,
		Exp:             time.Now().Add(30 * 24 * time.Hour).Unix(),
		PasswordVersion: user.PasswordVersion,
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(a.cfg.JwtSecret))
	if err != nil {
		return nil, errs.Unauthorized("Failed to generate token", err.Error())
	}

	return &token, nil
}
