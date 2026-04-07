package interceptor

import (
	"auth-service/internal/errs"
	"context"

	"buf.build/go/protovalidate"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func ValidationInterceptor() grpc.UnaryServerInterceptor {
	validator, err := protovalidate.New()
	if err != nil {
		log.Panic().Err(err).Msg("failed to initialize validator")
	}

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if msg, ok := req.(proto.Message); ok {
			if err := validator.Validate(msg); err != nil {
				return nil, errs.ToGRPC(errs.InvalidArgument("Validation failed!", err.Error()))
			}
		}
		return handler(ctx, req)
	}
}
