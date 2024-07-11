package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

type contextKey int

const userIDKey contextKey = iota
const tokenExp = time.Hour * 24
const secretKey = "supersecretkey"

func CreateAccessToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to create token: %w", err)
	}

	return tokenString, nil
}

func GetUserIDFromContext(ctx context.Context) (string, error) {
	value := ctx.Value(userIDKey)
	if value == nil {
		return "", errors.New("user ID not found in context")
	}

	userID, ok := value.(string)
	if !ok {
		return "", errors.New("context user ID is not a string")
	}

	return userID, nil
}

func GetUserID(tokenString string) (userID string, err error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}

			return []byte(secretKey), nil
		})
	if err != nil {
		return "", fmt.Errorf("unable to parse token: %w", err)
	}
	if !token.Valid {
		return "", fmt.Errorf("token is not valid: %w", err)
	}

	return claims.UserID, nil
}

func WithAuth(logger *zap.Logger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var userID string
			cookie, err := r.Cookie("auth_token")
			if err != nil {
				if !errors.Is(err, http.ErrNoCookie) {
					logger.Error("Cannot get cookie", zap.Error(err))
					http.Error(w, "", http.StatusInternalServerError)
					return
				}

				// Если нет куки с токеном
				userID = uuid.New().String()
				token, createTokenErr := CreateAccessToken(userID)
				if createTokenErr != nil {
					logger.Error("Cannot create access token", zap.Error(createTokenErr))
					http.Error(w, "", http.StatusInternalServerError)
					return
				}

				// Сетим куку с токеном
				http.SetCookie(w, &http.Cookie{
					Name:     "auth_token",
					Value:    token,
					Expires:  time.Now().Add(tokenExp),
					HttpOnly: true,
				})
			}

			if userID == "" {
				userID, err = GetUserID(cookie.Value)
				if err != nil {
					logger.Error("Cannot get userID", zap.Error(err))
					http.Error(w, "", http.StatusUnauthorized)
					return
				}
			}

			// Прокидываем userID из куки в контекст
			ctx := context.WithValue(r.Context(), userIDKey, userID)
			r = r.WithContext(ctx)

			h.ServeHTTP(w, r)
		})
	}
}
