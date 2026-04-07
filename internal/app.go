package internal

import (
	"auth-service/internal/config"
	"auth-service/internal/database"
	"auth-service/internal/handler"
	"auth-service/internal/interceptor"
	"auth-service/internal/service"
	"net"
	"strconv"

	"github.com/notweuz/auth-proto/pb"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type App struct {
	Server   *grpc.Server
	Listener net.Listener
}

func SetupApp(cfg *config.Config, db *gorm.DB) (*App, error) {
	list, err := net.Listen("tcp", ":"+strconv.Itoa(cfg.GRPCPort))
	if err != nil {
		return nil, err
	}

	userDatabase := database.NewUserDatabase(db)
	userService := service.NewUserService(userDatabase)
	userHandler := handler.NewUserHandler(userService)
	authService := service.NewAuthService(userService, cfg)
	authHandler := handler.NewAuthHandler(authService)

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.AuthInterceptor(cfg.JwtSecret, authService),
			interceptor.LoggerInterceptor(),
			interceptor.ValidationInterceptor(),
		),
	)

	pb.RegisterUserServiceServer(server, userHandler)
	pb.RegisterAuthServiceServer(server, authHandler)

	return &App{
		Server:   server,
		Listener: list,
	}, nil
}

func (a *App) Close() {
	if a.Server != nil {
		a.Server.GracefulStop()
	}
	if a.Listener != nil {
		err := a.Listener.Close()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to close listener")
		}
	}
}
