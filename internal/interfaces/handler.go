package interfaces

import "github.com/notweuz/auth-proto/pb"

type UserHandler interface {
	pb.UserServiceServer
}

type AuthHandler interface {
	pb.AuthServiceServer
}
