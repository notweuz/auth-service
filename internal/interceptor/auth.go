package interceptor

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type ctxKey string

const userIDKey ctxKey = "user_id"

func UserIDFromContext(ctx context.Context) (uint64, bool) {
	v, ok := ctx.Value(userIDKey).(uint64)
	return v, ok
}

func isPublicMethod(fullMethod string) bool {
	return strings.HasSuffix(fullMethod, "/Login") || strings.HasSuffix(fullMethod, "/Register")
}

func AuthInterceptor(secret string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if isPublicMethod(info.FullMethod) {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, fmt.Errorf("metadata not found")
		}

		auth := md.Get("authorization")
		if len(auth) == 0 {
			return nil, fmt.Errorf("authorization header not found")
		}

		tokenString := strings.TrimPrefix(auth[0], "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			return nil, fmt.Errorf("invalid token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, fmt.Errorf("invalid claims")
		}

		sub, ok := claims["user_id"].(string)
		if !ok {
			return nil, fmt.Errorf("user_id not found")
		}

		userID, err := strconv.ParseUint(sub, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid user_id")
		}

		ctx = context.WithValue(ctx, userIDKey, userID)
		return handler(ctx, req)
	}
}
