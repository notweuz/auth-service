package interceptor

import (
	"auth-service/internal/errs"
	"auth-service/internal/interfaces"
	"auth-service/internal/model"
	"context"
	"strings"

	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UserIDFromContext(ctx context.Context) (uint64, bool) {
	v, ok := ctx.Value("sub").(uint64)
	return v, ok
}

func isPublicMethod(fullMethod string) bool {
	return strings.HasSuffix(fullMethod, "/Login") || strings.HasSuffix(fullMethod, "/Register") || strings.HasSuffix(fullMethod, "/ValidateToken")
}

func AuthInterceptor(secret string, authService interfaces.AuthService) grpc.UnaryServerInterceptor {
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

		token := md.Get("authorization")
		if len(token) == 0 {
			return nil, errs.ToGRPC(errs.Unauthorized("Invalid credentials", "No token provided"))
		}

		tokenString := strings.TrimPrefix(token[0], "Bearer ")
		if tokenString == token[0] {
			return nil, errs.ToGRPC(errs.Unauthorized("Invalid credentials", "Invalid token"))
		}

		if active, err := authService.ValidateToken(tokenString); err != nil && active {
			return nil, errs.ToGRPC(errs.Unauthorized("Invalid credentials", "Invalid token"))
		}

		var claims model.Claims
		parsedToken, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil {
			return nil, errs.ToGRPC(errs.Unauthorized("Invalid credentials", "Invalid token"))
		}

		if !parsedToken.Valid {
			return nil, errs.ToGRPC(errs.Unauthorized("Invalid credentials", "Invalid token"))
		}

		ctx = context.WithValue(ctx, "sub", claims.Subject)

		return handler(ctx, req)
	}
}
