package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/SintroSecurity/go-libraries/logger"
	"github.com/SintroSecurity/go-libraries/router/response"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-jwt/jwt/v4"
)

type errorStruct struct {
	Error string `json:"error"`
}

// AuthenticateWithClaims validates ownership and validity of JWT tokens with custom claims type
func AuthenticateWithClaims[T jwt.Claims](ctx context.Context, config *Config, newClaims func() T) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := extractTokenFromRequest(r, config.CookieName)
			if err != nil {
				response.UnauthorizedJson(r, w, &errorStruct{Error: err.Error()})
				return
			}

			l := logger.GetLoggerFromContext(ctx)
			l.Debug("request token", l.String("request_id", middleware.GetReqID(ctx)), l.String("token", tokenString))

			deserializedToken, err := jwt.ParseWithClaims(tokenString, newClaims(), func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(config.Key), nil
			})

			if err != nil {
				response.UnauthorizedJson(r, w, &errorStruct{Error: "invalid token provided"})
				return
			}

			if !deserializedToken.Valid {
				response.UnauthorizedJson(r, w, &errorStruct{Error: "invalid token provided"})
				return
			}

			if claims, ok := deserializedToken.Claims.(T); !ok {
				response.UnauthorizedJson(r, w, &errorStruct{Error: "invalid token provided"})
				return
			} else {
				r = r.WithContext(NewContextWithClaims(r.Context(), claims))
				r = r.WithContext(NewContextWithRawToken(r.Context(), tokenString))
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Authenticate validates ownership and validity of JWT tokens using the default Token type
// This function is kept for backward compatibility
func Authenticate(ctx context.Context, config *Config) func(next http.Handler) http.Handler {
	return AuthenticateWithClaims(ctx, config, func() *Token {
		return &Token{}
	})
}

func extractTokenFromRequest(r *http.Request, cookieName string) (string, error) {
	var token string

	value := r.Header.Get("Authorization") //Bearer: xxxxx
	if value != "" {
		token = strings.ReplaceAll(strings.TrimPrefix(value, "Bearer "), " ", "")
	} else {
		cookie, err := r.Cookie(cookieName)
		if err != nil {
			return "", errors.New("token not provided")
		}
		token = cookie.Value
	}

	if token != "" {
		return token, nil
	} else {
		return "", errors.New("invalid token provided")
	}
}
