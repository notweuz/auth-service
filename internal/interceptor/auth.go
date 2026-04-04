package interceptor

import (
	"auth-service/internal/errs"
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
			return nil, errs.ToGRPC(errs.Unauthorized("Invalid credentials", "Metadata not found"))
		}

		auth := md.Get("authorization")
		if len(auth) == 0 {
			return nil, errs.ToGRPC(errs.Unauthorized("Invalid credentials", "No token provided"))
		}

		tokenString := strings.TrimPrefix(auth[0], "Bearer ")
		if tokenString == auth[0] {
			return nil, errs.ToGRPC(errs.Unauthorized("Invalid credentials", "Invalid authorization header"))
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			return nil, errs.ToGRPC(errs.Unauthorized("Invalid credentials", "Invalid token"))
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, errs.ToGRPC(errs.Unauthorized("Invalid credentials", "Invalid token claims"))
		}

		userIDValue, ok := claims["user_id"]
		if !ok {
			return nil, errs.ToGRPC(errs.Unauthorized("Invalid credentials", "user_id not found"))
		}

		var userID uint64
		switch value := userIDValue.(type) {
		case string:
			userID, err = strconv.ParseUint(value, 10, 64)
		case float64:
			userID = uint64(value)
		default:
			err = fmt.Errorf("unsupported user_id type")
		}
		if err != nil {
			return nil, errs.ToGRPC(errs.Unauthorized("Invalid credentials", "Invalid user_id"))
		}

		ctx = context.WithValue(ctx, userIDKey, userID)
		return handler(ctx, req)
	}
}
