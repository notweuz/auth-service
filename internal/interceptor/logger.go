package interceptor

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func LoggerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		st, _ := status.FromError(err)
		event := log.Info()

		if err != nil {
			event = log.Error().Err(err)
		}

		if p, ok := peer.FromContext(ctx); ok {
			event = event.Str("peer", p.Addr.String())
		}

		if userID, ok := UserIDFromContext(ctx); ok {
			event = event.Uint64("sub", userID)
		}

		event.
			Str("method", info.FullMethod).
			Str("status", st.Code().String()).
			Dur("duration", time.Since(start)).
			Msg("gRPC request")

		return resp, err
	}
}
