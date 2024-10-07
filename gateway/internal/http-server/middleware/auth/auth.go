package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"prtf-gateway/internal/lib/jwt"
	"prtf-gateway/internal/lib/logger/sl"
	"prtf-gateway/internal/sso/grpc"
	"strings"
)

var (
	ErrInvalidToken       = errors.New("invalid token")
	ErrFailedIsAdminCheck = errors.New("failed to check if user is admin")
)

const (
	errorKey   = "error"
	uidKey     = "user_id"
	isAdminKey = "is_admin"
)

func extractAccessToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 {
		return ""
	}

	return splitToken[1]
}

func NewAuthMiddleware(log *slog.Logger, permProvider grpc.Client) func(next http.Handler) http.Handler {
	const op = "middleware.auth.NewAuthMiddleware"

	log = log.With(slog.String("op", op))

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := extractAccessToken(r)
			if tokenStr == "" {
				next.ServeHTTP(w, r)
				return
			}

			claims, err := jwt.ParseToken(tokenStr)
			if err != nil {
				log.Warn("failed to parse token", sl.Err(err))

				ctx := context.WithValue(r.Context(), errorKey, ErrInvalidToken)
				next.ServeHTTP(w, r.WithContext(ctx))

				return
			}

			log.Info("user authorized", slog.Any("claims", claims))

			isAdmin, err := permProvider.IsAdmin(r.Context(), claims.UserId)
			if err != nil {
				log.Error("failed to check if user is admin", sl.Err(err))

				ctx := context.WithValue(r.Context(), errorKey, ErrFailedIsAdminCheck)
				next.ServeHTTP(w, r.WithContext(ctx))

				return
			}

			ctx := context.WithValue(r.Context(), uidKey, claims.UserId)
			ctx = context.WithValue(r.Context(), isAdminKey, isAdmin)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UIDFromContext(ctx context.Context) (int64, bool) {
	uid, ok := ctx.Value(uidKey).(int64)
	return uid, ok
}

func ErrorFromContext(ctx context.Context) (error, bool) {
	err, ok := ctx.Value(errorKey).(error)
	return err, ok
}
