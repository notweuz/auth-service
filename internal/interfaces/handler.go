package interfaces

import "github.com/notweuz/authentication-proto/pb"

type UserHandler interface {
	pb.UserServiceServer
}

type AuthHandler interface {
	pb.AuthServiceServer
}
